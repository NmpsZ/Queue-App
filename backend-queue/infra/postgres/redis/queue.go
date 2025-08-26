package redis

import (
	"context"
	"strconv"
	"time"

	redisrepo "backend-queue/usecase/queue/redis"

	"github.com/redis/go-redis/v9"
)

const waitingQueueKey = "queue:waiting"

type QueueRedisRepo struct {
	client *redis.Client
}

func NewQueueRedisRepo(client *redis.Client) redisrepo.Repository {
	return &QueueRedisRepo{
		client: client,
	}
}

// Push queueID เข้า waiting list
func (r *QueueRedisRepo) PushWaitingQueue(ctx context.Context, queueID uint) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	return r.client.LPush(ctx, waitingQueueKey, strconv.Itoa(int(queueID))).Err()
}

// Pop queueID ถัดไป
func (r *QueueRedisRepo) PopNextWaiting(ctx context.Context) (uint, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	idStr, err := r.client.RPop(ctx, waitingQueueKey).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, nil // ไม่มี queue ใน list
		}
		return 0, err
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, err
	}

	return uint(id), nil
}
