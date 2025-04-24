package main

import (
	"log"
	"os"

	"cdr/cmd/common"
	"cdr/config"
	"cdr/service"
)

func main() {
	log.Println("呼叫状态推送系统启动...")

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
	callStatusService := service.NewCallStatusService(cfg, cdrService)

	// 初始化并启动健康检查服务
	healthService := service.NewHealthService(cfg, callStatusService)
	go func() {
		if err := healthService.StartHealthServer("9090"); err != nil {
			log.Printf("启动健康检查服务失败: %v", err)
		}
	}()

	// 使用通用工作池处理呼叫状态更新
	common.StartWorkerPool(cfg.Push.Workers, func() error {
		// 先创建新呼叫
		if err := callStatusService.StartNewCall(); err != nil {
			log.Printf("创建新呼叫失败: %v", err)
			return err
		}

		// 然后更新现有呼叫的状态
		if err := callStatusService.UpdateCallStatus(); err != nil {
			log.Printf("更新呼叫状态失败: %v", err)
			return err
		}
		return nil
	})
}

