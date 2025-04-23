package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"cdr/config"
	"cdr/models"

	"github.com/google/uuid"
)

// CDRService 处理CDR相关的业务逻辑
type CDRService struct {
	config *config.Config
	logger *Logger
	workerPool *WorkerPool
}

// NewCDRService 创建CDR服务实例
func NewCDRService(cfg *config.Config) (*CDRService, error) {
	logger, err := NewLogger("cdr")
	if err != nil {
		return nil, fmt.Errorf("创建日志记录器失败: %v", err)
	}
	return &CDRService{
		config:     cfg,
		logger:     logger,
		workerPool: NewWorkerPool(cfg.Push.Workers),
	}, nil
}

// GenerateCallID 生成唯一的通话ID
func (s *CDRService) GenerateCallID() string {
	// 格式：NM + 时间戳 + uuid前8位
	timestamp := time.Now().Format("200601021504051150")
	uid := strings.Replace(uuid.New().String(), "-", "", -1)[:8]
	return fmt.Sprintf("NM%s%s", timestamp, uid)
}

// GeneratePhoneNumber 生成随机手机号
func (s *CDRService) GeneratePhoneNumber() string {
	prefixes := []string{"131", "132", "133", "134", "135", "136", "137", "138", "139"}
	prefix := prefixes[rand.Intn(len(prefixes))]
	number := rand.Intn(100000000)
	return fmt.Sprintf("%s%08d", prefix, number)
}

// GenerateCDR 生成模拟CDR记录
func (s *CDRService) GenerateCDR() *models.CDR {
	now := time.Now()
	beginTime := now.Add(-time.Duration(rand.Intn(3600)) * time.Second)
	duration := rand.Intn(600) // 最长通话10分钟

	return &models.CDR{
		AccountID:     s.config.Account.ID,
		CallID:        s.GenerateCallID(),
		ServiceType:   s.config.Account.ServiceType,
		Caller:        s.GeneratePhoneNumber(),
		Callee:        s.GeneratePhoneNumber(),
		BeginCallTime: beginTime.UnixNano() / 1e6,
		StartTime:     beginTime.Add(5*time.Second).UnixNano() / 1e6,
		EndTime:       beginTime.Add(time.Duration(duration)*time.Second).UnixNano() / 1e6,
		CallDuration:  duration,
		CallResult:    1,
		CDRCreateTime: now.UnixNano() / 1e6,
		UserData:      fmt.Sprintf("{\"simulateTime\":\"%s\"}", now.Format(time.RFC3339)),
		MessageType:   1,
		CDRType:       1,
	}
}

// PushCDR 推送CDR记录
func (s *CDRService) PushCDR(cdr *models.CDR) error {
	if cdr == nil {
		cdr = s.GenerateCDR()
	}

	jsonData, err := json.Marshal(cdr)
	if err != nil {
		return fmt.Errorf("JSON序列化失败: %v", err)
	}

	// 创建错误通道用于收集推送结果
	errChan := make(chan error, 1)

	// 提交推送任务到工作池
	s.workerPool.Submit(func() {
		var lastErr error
		for i := 0; i < s.config.Retry.Times; i++ {
			if i > 0 {
				delay := s.config.Retry.Delays[i]
				log.Printf("CDR推送重试 CallID:%s, 第%d次, 等待%d秒...", cdr.CallID, i, delay)
				time.Sleep(time.Duration(delay) * time.Second)
			}

			resp, err := http.Post(s.config.Push.CdrURL, "application/json", bytes.NewBuffer(jsonData))
			if err != nil {
				lastErr = fmt.Errorf("HTTP请求失败: %v", err)
				s.logger.LogPushCDR(cdr.CallID, s.config.Push.CdrURL, jsonData, 0, err)
				log.Printf("推送失败: %v", lastErr)
				continue
			}
			defer resp.Body.Close()

			s.logger.LogPushCDR(cdr.CallID, s.config.Push.CdrURL, jsonData, resp.StatusCode, nil)
			if resp.StatusCode == http.StatusOK {
				log.Printf("CDR推送成功，CallID: %v", cdr.CallID)
				errChan <- nil
				return
			}

			lastErr = fmt.Errorf("推送失败，状态码: %d", resp.StatusCode)
			log.Printf("推送失败: %v", lastErr)
		}

		errChan <- fmt.Errorf("CDR推送重试%d次失败，最后错误: %v", s.config.Retry.Times, lastErr)
	})

	return <-errChan
}
