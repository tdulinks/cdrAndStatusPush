package main

import (
	"log"
	"math/rand"
	"os"
	"sync"
	"time"

	"cdr/config"
	"cdr/service"
)

func main() {
	log.Println("呼叫状态推送系统启动...")

	// 初始化随机数种子
	rand.Seed(time.Now().UnixNano())

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

	// 持续创建新的呼叫，使用更短的时间间隔
	go func() {
		ticker := time.NewTicker(100 * time.Millisecond) // 减少到100毫秒
		defer ticker.Stop()
		for range ticker.C {
			if err := callStatusService.StartNewCall(); err != nil {
				log.Printf("创建新呼叫失败: %v", err)
			}
		}
	}()

	// 持续批量更新现有呼叫的状态
	var wg sync.WaitGroup
	wg.Add(cfg.Push.Workers * 2)

	// 启动更多的工作协程来处理状态更新
	for i := 0; i < cfg.Push.Workers*2; i++ {
		go func() {
			defer wg.Done()
			ticker := time.NewTicker(50 * time.Millisecond) // 更频繁地检查和更新状态
			defer ticker.Stop()
			for range ticker.C {
				if err := callStatusService.UpdateCallStatus(); err != nil {
					log.Printf("更新呼叫状态失败: %v", err)
				}
			}
		}()
	}

	// 等待所有工作协程完成
	wg.Wait()
}
