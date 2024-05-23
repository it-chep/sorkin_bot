package worker_pool

import (
	"context"
	"time"
)

type Worker struct {
	interval time.Duration
	task     Task
}

func NewWorker(task Task, interval time.Duration) Worker {
	return Worker{
		task:     task,
		interval: interval,
	}
}

func (w Worker) Start() {
	for {
		go func() {
			_ = w.task.Process(context.Background())
		}()
		time.Sleep(w.interval)
	}
}

type WorkerPool struct {
	workers []Worker
	limit   int
}

func NewWorkerPool(workers []Worker) WorkerPool {
	return WorkerPool{
		workers: workers,
	}
}

// todo ctx
func (wp WorkerPool) Run() {
	for _, w := range wp.workers {
		go w.Start()
	}
}
