package worker_pool

import (
	"context"
	"time"
)

type Worker struct {
	interval time.Duration
	task     Task
}

func NewWorker(task Task) Worker {
	return Worker{
		task:     task,
		interval: 24 * time.Hour,
	}
}

func (w Worker) Run() {
	for {
		go func() {
			_ = w.task.Process(context.Background())
		}()
		time.Sleep(w.interval)
	}
}
