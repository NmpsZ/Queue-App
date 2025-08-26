package redis

import "context"

type Repository interface {
	PushWaitingQueue(ctx context.Context, queueID uint) error
	PopNextWaiting(ctx context.Context) (uint, error)
}
