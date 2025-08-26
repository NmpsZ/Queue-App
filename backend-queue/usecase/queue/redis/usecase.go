package redis

import (
	"backend-queue/entity"
	"context"
)

type UseCase interface {
	FindQueueByName(ctx context.Context, name string) ([]entity.Queue, error)
	CallNextQueue(ctx context.Context) (*entity.Queue, error)
}
