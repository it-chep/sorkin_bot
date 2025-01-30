package worker_pool

import (
	"context"
	"time"
)

type Task interface {
	Process(ctx context.Context) error
	NextSchedule(now time.Time) time.Time
}
