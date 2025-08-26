package postgres

import (
	"backend-queue/entity"
	"backend-queue/usecase/queue"
	"context"

	"gorm.io/gorm"
)

type QueueRepository struct {
	db *gorm.DB
}

func NewQueueRepository(db *gorm.DB) queue.Repository {
	return &QueueRepository{
		db: db,
	}
}

func (r *QueueRepository) Find(ctx context.Context) ([]entity.Queue, error) {
	var queues []entity.Queue
	if err := r.db.WithContext(ctx).Find(&queues).Error; err != nil {
		return nil, err
	}
	return queues, nil
}

func (r *QueueRepository) FindByID(ctx context.Context, id uint) (*entity.Queue, error) {
	var queue entity.Queue
	if err := r.db.WithContext(ctx).First(&queue, id).Error; err != nil {
		return nil, err
	}
	return &queue, nil
}

func (r *QueueRepository) Create(ctx context.Context, queue *entity.Queue) error {
	return r.db.WithContext(ctx).Create(queue).Error
}

func (r *QueueRepository) Update(ctx context.Context, queue *entity.Queue) error {
	return r.db.WithContext(ctx).Save(queue).Error
}

func (r *QueueRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entity.Queue{}, id).Error
}

func (r *QueueRepository) FindByQueueNo(ctx context.Context, queueNo string) (*entity.Queue, error) {
	var queues entity.Queue
	err := r.db.WithContext(ctx).Where("queue_no = ?", queueNo).First(&queues).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &queues, nil
}

func (r *QueueRepository) FindQueueByName(ctx context.Context, name string) ([]entity.Queue, error) {
	var queues []entity.Queue
	if err := r.db.WithContext(ctx).
		Where("name = ?", name).
		Find(&queues).Error; err != nil {
		return nil, err
	}
	return queues, nil
}

func (r *QueueRepository) FindActiveQueueByName(ctx context.Context, name string) (*entity.Queue, error) {
	var queue entity.Queue
	err := r.db.WithContext(ctx).
		Where("name = ? AND status = ?", name, "waiting").
		First(&queue).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &queue, nil
}

func (r *QueueRepository) FindNextWaiting(ctx context.Context) (*entity.Queue, error) {
	var queue entity.Queue
	if err := r.db.WithContext(ctx).
		Where("status = ?", "waiting").
		Order("created_at ASC").
		First(&queue).Error; err != nil {
		return nil, err
	}
	return &queue, nil
}
