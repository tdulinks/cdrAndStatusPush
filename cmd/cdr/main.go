package main

import (
	"log"
	"os"

	"cdr/cmd/common"
	"cdr/config"
	"cdr/service"
)

func main() {
	log.Println("话单推送系统启动...")

	// 加载配置
	cfg, err := config.LoadConfig(config.GetConfigPath())
	if err != nil {
		log.Printf("加载配置失败: %v", err)
		os.Exit(1)
	}

	// 初始化服务
	cdrService, err := service.NewCDRService(cfg)
	if err != nil {
		log.Printf("初始化CDR服务失败: %v", err)
		os.Exit(1)
	}

	// 使用通用工作池处理CDR推送
	common.StartWorkerPool(cfg.Push.Workers, func() error {
		if err := cdrService.PushCDR(nil); err != nil {
			log.Printf("推送CDR记录失败: %v", err)
			return err
		}
		return nil
	})
}
