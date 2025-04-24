package common

import (
	"sync"
)

// StartWorkerPool 启动工作池，处理通用的工作任务
func StartWorkerPool(workers int, handler func() error) {
	workerChan := make(chan struct{}, workers)
	var wg sync.WaitGroup
	wg.Add(workers)

	// 启动工作协程
	for range make([]struct{}, workers) {
		go func() {
			defer wg.Done()
			for range workerChan {
				if err := handler(); err != nil {
					// 错误处理由调用者负责
					continue
				}
			}
		}()
	}

	// 持续填充工作通道
	for {
		for range make([]struct{}, workers) {
			workerChan <- struct{}{}
		}
	}
}
