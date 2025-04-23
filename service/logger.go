package service

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

// Logger 处理日志记录的结构体
type Logger struct {
	logDir      string
	statusFile  *os.File // 状态推送日志文件
	cdrFile     *os.File // CDR推送日志文件
	maxFileSize int64
}

// NewLogger 创建日志记录器实例
func NewLogger(logTypes ...string) (*Logger, error) {
	// 创建logs目录
	logDir := "logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("创建日志目录失败: %v", err)
	}

	// 创建日志文件
	logger := &Logger{
		logDir:      logDir,
		maxFileSize: 10 * 1024 * 1024, // 10MB
	}

	// 根据传入的日志类型创建对应的日志文件
	for _, logType := range logTypes {
		if err := logger.rotateLogFile(logType); err != nil {
			return nil, err
		}
	}

	return logger, nil
}

// rotateLogFile 创建或轮转日志文件
func (l *Logger) rotateLogFile(fileType string) error {
	// 生成新的日志文件名
	timestamp := time.Now().Format("20060102150405")
	var fileName string
	if fileType == "status" {
		fileName = fmt.Sprintf("push_status_%s.log", timestamp)
	} else {
		fileName = fmt.Sprintf("push_cdr_%s.log", timestamp)
	}
	filePath := filepath.Join(l.logDir, fileName)

	// 创建新的日志文件
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("创建日志文件失败: %v", err)
	}

	// 根据文件类型关闭和更新相应的文件指针
	if fileType == "status" {
		if l.statusFile != nil {
			l.statusFile.Close()
		}
		l.statusFile = file
	} else {
		if l.cdrFile != nil {
			l.cdrFile.Close()
		}
		l.cdrFile = file
	}
	return nil
}

// LogPushStatus 记录推送状态的日志
func (l *Logger) LogPushStatus(callID string, url string, requestData []byte, statusCode int, responseErr error) {
	// 检查文件大小是否需要轮转
	if info, err := l.statusFile.Stat(); err == nil && info.Size() > l.maxFileSize {
		if err := l.rotateLogFile("status"); err != nil {
			log.Printf("轮转状态日志文件失败: %v", err)
			return
		}
	}

	// 格式化日志内容
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logContent := fmt.Sprintf("[%s] CallID: %s\nURL: %s\nRequest: %s\nStatusCode: %d\n",
		timestamp, callID, url, string(requestData), statusCode)

	// 如果有错误，添加错误信息
	if responseErr != nil {
		logContent += fmt.Sprintf("Error: %v\n", responseErr)
	}

	logContent += "----------------------------------------\n"

	// 写入日志文件
	if _, err := l.statusFile.WriteString(logContent); err != nil {
		log.Printf("写入状态日志失败: %v", err)
	}
}

// LogPushCDR 记录CDR推送的日志
func (l *Logger) LogPushCDR(callID string, url string, requestData []byte, statusCode int, responseErr error) {
	// 检查文件大小是否需要轮转
	if info, err := l.cdrFile.Stat(); err == nil && info.Size() > l.maxFileSize {
		if err := l.rotateLogFile("cdr"); err != nil {
			log.Printf("轮转CDR日志文件失败: %v", err)
			return
		}
	}

	// 格式化日志内容
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logContent := fmt.Sprintf("[%s] CallID: %s\nURL: %s\nRequest: %s\nStatusCode: %d\n",
		timestamp, callID, url, string(requestData), statusCode)

	// 如果有错误，添加错误信息
	if responseErr != nil {
		logContent += fmt.Sprintf("Error: %v\n", responseErr)
	}

	logContent += "----------------------------------------\n"

	// 写入日志文件
	if _, err := l.cdrFile.WriteString(logContent); err != nil {
		log.Printf("写入CDR日志失败: %v", err)
	}
}

// Close 关闭日志文件
func (l *Logger) Close() error {
	// 关闭状态推送日志文件
	if l.statusFile != nil {
		if err := l.statusFile.Close(); err != nil {
			return err
		}
	}
	// 关闭CDR推送日志文件
	if l.cdrFile != nil {
		if err := l.cdrFile.Close(); err != nil {
			return err
		}
	}
	return nil
}
