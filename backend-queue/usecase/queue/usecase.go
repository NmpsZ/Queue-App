package queue

import (
	"backend-queue/entity"
	"context"
)

type UseCase interface {
	FindQueue(ctx context.Context) ([]entity.Queue, error)
	GetQueueByID(ctx context.Context, id uint) (*entity.Queue, error)
	AddQueue(ctx context.Context, q *entity.Queue) error
	UpdateQueue(ctx context.Context, q *entity.Queue) error
	DeleteQueue(ctx context.Context, id uint) error
	CreateQueueWithQR(ctx context.Context, q *entity.Queue) (*entity.Queue, string, error)
	FindQueueByName(ctx context.Context, name string) ([]entity.Queue, error)
	CallNextQueue(ctx context.Context) (*entity.Queue, error)
}
