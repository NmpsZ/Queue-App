package queue

import (
	"backend-queue/entity"
	"context"
)

type Repository interface {
	Find(ctx context.Context) ([]entity.Queue, error)
	FindByID(ctx context.Context, id uint) (*entity.Queue, error)
	Create(ctx context.Context, queue *entity.Queue) error
	Update(ctx context.Context, queue *entity.Queue) error
	Delete(ctx context.Context, id uint) error
	FindByQueueNo(ctx context.Context, queueNo string) (*entity.Queue, error)
	FindQueueByName(ctx context.Context, name string) ([]entity.Queue, error)
	FindActiveQueueByName(ctx context.Context, name string) (*entity.Queue, error)
	FindNextWaiting(ctx context.Context) (*entity.Queue, error)
}
