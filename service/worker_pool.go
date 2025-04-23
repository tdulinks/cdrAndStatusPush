package service

import (
	"log"
	"sync"
)

// WorkerPool 工作池结构体
type WorkerPool struct {
	workerCount int
	jobQueue   chan func()
	wg         sync.WaitGroup
}

// NewWorkerPool 创建新的工作池
func NewWorkerPool(workerCount int) *WorkerPool {
	pool := &WorkerPool{
		workerCount: workerCount,
		jobQueue:   make(chan func(), workerCount*2), // 任务队列容量设为工作协程数的2倍
	}

	// 启动工作协程
	pool.wg.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		go pool.worker()
	}

	return pool
}

// worker 工作协程
func (p *WorkerPool) worker() {
	defer p.wg.Done()

	for job := range p.jobQueue {
		job()
	}
}

// Submit 提交任务到工作池
func (p *WorkerPool) Submit(job func()) {
	p.jobQueue <- job
}

// Close 关闭工作池
func (p *WorkerPool) Close() {
	close(p.jobQueue)
	p.wg.Wait()
	log.Println("工作池已关闭")
}