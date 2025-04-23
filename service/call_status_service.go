package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"cdr/config"
	"cdr/models"
)

// CallStatusService 处理呼叫状态推送的业务逻辑
type CallStatusService struct {
	config       *config.Config
	cdrService   *CDRService          // 用于生成号码等功能
	currentCalls map[string]*callInfo // 记录当前进行中的通话
	logger       *Logger              // 日志记录器
	workerPool   *WorkerPool          // 工作池
	mutex        sync.Mutex           // 用于保护 currentCalls map 的并发访问
}

type callInfo struct {
	status     *models.CallStatus
	eventIndex int
	startTime  time.Time
}

// NewCallStatusService 创建呼叫状态服务实例
func NewCallStatusService(cfg *config.Config, cdrService *CDRService) *CallStatusService {
	logger, err := NewLogger("status")
	if err != nil {
		log.Printf("初始化日志记录器失败: %v", err)
	}

	return &CallStatusService{
		config:       cfg,
		cdrService:   cdrService,
		currentCalls: make(map[string]*callInfo),
		logger:       logger,
		workerPool:   NewWorkerPool(cfg.Push.Workers),
	}
}

// StartNewCall 开始一个新的呼叫并推送第一个状态
func (s *CallStatusService) StartNewCall() error {
	// 生成新的呼叫信息
	status := &models.CallStatus{
		AccountID:      s.config.Account.ID,
		CallID:         s.cdrService.GenerateCallID(),
		ServiceType:    s.config.Account.ServiceType,
		Caller:         s.cdrService.GeneratePhoneNumber(),
		Callee:         s.cdrService.GeneratePhoneNumber(),
		EventTime:      fmt.Sprintf("%d", time.Now().Unix()),
		EventType:      models.EventTypeCalling,
		AllEventType:   []int{models.EventTypeCalling},
		MessageType:    1,
		Party:          1,
		SubscriptionID: "sim_" + time.Now().Format("20060102150405"),
		UserData:       fmt.Sprintf("{\"startTime\":\"%d\"}", time.Now().Unix()),
	}

	// 保存呼叫信息
	s.mutex.Lock()
	s.currentCalls[status.CallID] = &callInfo{
		status:     status,
		eventIndex: 0,
		startTime:  time.Now(),
	}
	s.mutex.Unlock()

	// 推送第一个状态
	return s.pushStatus(status)
}

// UpdateCallStatus 更新现有呼叫的状态
func (s *CallStatusService) UpdateCallStatus() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for callID, info := range s.currentCalls {
		// 如果通话已经结束，从map中删除
		if info.status.EventType == models.EventTypeEnded {
			delete(s.currentCalls, callID)
			continue
		}

		// 获取下一个状态
		nextEventType := s.getNextEventType(info.status.EventType)
		info.status.EventType = nextEventType
		info.status.AllEventType = append(info.status.AllEventType, nextEventType)
		info.status.EventTime = fmt.Sprintf("%d", time.Now().Unix())
		info.eventIndex++

		// 推送状态
		if err := s.pushStatus(info.status); err != nil {
			log.Printf("推送状态失败 CallID:%s, Error:%v", callID, err)
		}
	}
	return nil
}

// getNextEventType 获取下一个状态
func (s *CallStatusService) getNextEventType(currentType int) int {
	switch currentType {
	case models.EventTypeCalling:
		return models.EventTypeRinging
	case models.EventTypeRinging:
		return models.EventTypeAnswered
	case models.EventTypeAnswered:
		return models.EventTypeEnded
	default:
		return models.EventTypeEnded
	}
}

// pushStatus 推送状态（带重试机制）
func (s *CallStatusService) pushStatus(status *models.CallStatus) error {
	jsonData, err := json.Marshal(status)
	if err != nil {
		return fmt.Errorf("JSON序列化失败: %v", err)
	}

	// 记录推送的URL
	pushURL := s.config.Push.StatusURL

	// 创建错误通道用于收集推送结果
	errChan := make(chan error, 1)

	// 提交推送任务到工作池
	s.workerPool.Submit(func() {
		var lastErr error
		for i := 0; i < s.config.Retry.Times; i++ {
			if i > 0 {
				delay := s.config.Retry.Delays[i]
				log.Printf("状态推送重试 CallID:%s, 第%d次, 等待%d秒...", status.CallID, i, delay)
				time.Sleep(time.Duration(delay) * time.Second)
			}

			resp, err := http.Post(pushURL, "application/json", bytes.NewBuffer(jsonData))
			if err != nil {
				lastErr = fmt.Errorf("HTTP请求失败: %v", err)
				log.Printf("状态推送失败: %v", lastErr)
				// 记录失败日志
				if s.logger != nil {
					s.logger.LogPushStatus(status.CallID, pushURL, jsonData, 0, lastErr)
				}
				continue
			}
			defer resp.Body.Close()

			// 记录推送结果
			if s.logger != nil {
				s.logger.LogPushStatus(status.CallID, pushURL, jsonData, resp.StatusCode, nil)
			}

			if resp.StatusCode == http.StatusOK {
				log.Printf("状态推送成功 CallID:%s, EventType:%d", status.CallID, status.EventType)
				errChan <- nil
				return
			}

			lastErr = fmt.Errorf("推送失败，状态码: %d", resp.StatusCode)
			log.Printf("状态推送失败: %v", lastErr)
		}

		errChan <- fmt.Errorf("状态推送重试%d次失败，最后错误: %v", s.config.Retry.Times, lastErr)
	})

	return <-errChan
}
