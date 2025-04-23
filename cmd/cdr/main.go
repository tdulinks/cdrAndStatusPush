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
	log.Println("话单推送系统启动...")

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

	// 持续批量生成并推送CDR记录
	workers := make(chan struct{}, cfg.Push.Workers)
	var wg sync.WaitGroup
	wg.Add(cfg.Push.Workers)

	for i := 0; i < cfg.Push.Workers; i++ {
		go func() {
			defer wg.Done()
			for range workers {
				if err := cdrService.PushCDR(nil); err != nil {
					log.Printf("推送CDR记录失败: %v", err)
				}
			}
		}()
	}

	// 持续填充工作通道
	for {
		for i := 0; i < cfg.Push.Workers; i++ {
			workers <- struct{}{}
		}
	}
}