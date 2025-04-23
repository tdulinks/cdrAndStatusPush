package service

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"cdr/config"
)

// HealthService 提供系统健康检查功能
type HealthService struct {
	config          *config.Config
	callStatusSvc   *CallStatusService
	lastCheckTime   time.Time
	lastCheckResult *HealthStatus
	mux             sync.RWMutex
}

// HealthStatus 表示系统健康状态
type HealthStatus struct {
	Status           string    `json:"status"`            // 整体状态："healthy" 或 "unhealthy"
	Timestamp        time.Time `json:"timestamp"`         // 检查时间
	ConfigStatus     string    `json:"configStatus"`      // 配置状态
	CallServiceState string    `json:"callServiceState"`  // 呼叫服务状态
	Details          string    `json:"details,omitempty"` // 详细信息（如果有错误）
}

// NewHealthService 创建健康检查服务实例
func NewHealthService(cfg *config.Config, callStatusSvc *CallStatusService) *HealthService {
	return &HealthService{
		config:        cfg,
		callStatusSvc: callStatusSvc,
	}
}

// StartHealthServer 启动健康检查HTTP服务
func (h *HealthService) StartHealthServer(port string) error {
	http.HandleFunc("/health", h.handleHealthCheck)
	return http.ListenAndServe(":"+port, nil)
}

// handleHealthCheck 处理健康检查HTTP请求
func (h *HealthService) handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	h.mux.Lock()
	defer h.mux.Unlock()

	// 如果距离上次检查不到5秒，直接返回缓存的结果
	if time.Since(h.lastCheckTime) < 5*time.Second && h.lastCheckResult != nil {
		h.writeResponse(w, h.lastCheckResult)
		return
	}

	// 执行健康检查
	status := h.checkHealth()
	h.lastCheckTime = time.Now()
	h.lastCheckResult = status

	h.writeResponse(w, status)
}

// checkHealth 执行健康检查
func (h *HealthService) checkHealth() *HealthStatus {
	status := &HealthStatus{
		Timestamp: time.Now(),
	}

	// 检查配置状态
	if h.config == nil {
		status.ConfigStatus = "unhealthy"
		status.Details = "配置未加载"
	} else {
		status.ConfigStatus = "healthy"
	}

	// 检查呼叫服务状态
	if h.callStatusSvc == nil {
		status.CallServiceState = "unhealthy"
		status.Details = "呼叫服务未初始化"
	} else {
		status.CallServiceState = "healthy"
	}

	// 设置整体状态
	if status.ConfigStatus == "healthy" && status.CallServiceState == "healthy" {
		status.Status = "healthy"
	} else {
		status.Status = "unhealthy"
	}

	return status
}

// writeResponse 写入HTTP响应
func (h *HealthService) writeResponse(w http.ResponseWriter, status *HealthStatus) {
	w.Header().Set("Content-Type", "application/json")
	if status.Status != "healthy" {
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	json.NewEncoder(w).Encode(status)
}
