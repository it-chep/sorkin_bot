package worker_pool

import (
	"context"
	"log"
	"time"
)

type Worker struct {
	task Task
}

func NewWorker(task Task) Worker {
	return Worker{
		task: task,
	}
}

func (w Worker) Start(ctx context.Context) {
LOOP:
	for {
		now := time.Now().UTC()
		select {
		case <-ctx.Done():
			break LOOP
		case <-time.After(w.task.NextSchedule(now).Sub(now)):
			if err := w.task.Process(ctx); err != nil {
				log.Printf("task process error: %v", err) // todo: log
			}
		}
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

func (wp WorkerPool) Run(ctx context.Context) {
	for _, w := range wp.workers {
		go w.Start(ctx)
	}
}
