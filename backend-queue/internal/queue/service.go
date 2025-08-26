package queue

import (
	"backend-queue/entity"
	"backend-queue/usecase/queue"
	redisrepo "backend-queue/usecase/queue/redis"
	"backend-queue/utils"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type queueUseCase struct {
	repo       queue.Repository
	redisrepo  redisrepo.Repository
	timeNowUTC func() time.Time
}

func NewQueueUseCase(repo queue.Repository, redisrepo redisrepo.Repository) queue.UseCase {
	return &queueUseCase{
		repo:       repo,
		redisrepo:  redisrepo,
		timeNowUTC: func() time.Time { return time.Now().UTC() },
	}
}

func (u *queueUseCase) FindQueue(ctx context.Context) ([]entity.Queue, error) {
	return u.repo.Find(ctx)
}

func (u *queueUseCase) GetQueueByID(ctx context.Context, id uint) (*entity.Queue, error) {
	return u.repo.FindByID(ctx, id)
}

func (u *queueUseCase) AddQueue(ctx context.Context, q *entity.Queue) error {

	if q.QueueNo != "" {
		existing, _ := u.repo.FindByQueueNo(ctx, q.QueueNo)
		if existing != nil {
			return fmt.Errorf("queue number %s already exists", q.QueueNo)
		}
	} else {
		// ถ้าไม่มี QueueNo ให้ auto-generate
		q.QueueNo = generateQueueNo()
	}
	q.CreatedAt = u.timeNowUTC()
	q.UpdatedAt = u.timeNowUTC()
	return u.repo.Create(ctx, q)
}

func (u *queueUseCase) UpdateQueue(ctx context.Context, q *entity.Queue) error {
	q.UpdatedAt = u.timeNowUTC()
	return u.repo.Update(ctx, q)
}

func (u *queueUseCase) DeleteQueue(ctx context.Context, id uint) error {
	return u.repo.Delete(ctx, id)
}

func (u *queueUseCase) CreateQueueWithQR(ctx context.Context, q *entity.Queue) (*entity.Queue, string, error) {
	// 1. Validation - ชื่อต้องไม่ว่าง
	if strings.TrimSpace(q.Name) == "" {
		return nil, "", errors.New("name is required")
	}

	// 2. Validation - ตรวจสอบว่ามีคิวที่ชื่อเดียวกันและยัง active อยู่หรือไม่
	existing, err := u.repo.FindActiveQueueByName(ctx, q.Name)
	if err != nil {
		return nil, "", err
	}
	if existing != nil {
		return nil, "", fmt.Errorf("queue already exists for name: %s", q.Name)
	}

	// 3. Generate QueueNo ถ้ายังไม่มี
	if q.QueueNo == "" {
		q.QueueNo = generateQueueNo()
	}

	q.CreatedAt = u.timeNowUTC()
	q.UpdatedAt = u.timeNowUTC()

	// 4. Insert DB
	if err := u.repo.Create(ctx, q); err != nil {
		return nil, "", err
	}

	// 5. Generate QR code จาก QueueNo
	qrBase64, err := utils.GenerateQueueQRCodeBase64(q.QueueNo)
	if err != nil {
		return nil, "", err
	}

	// 6. Update QR code field ใน entity
	q.QRCode = qrBase64
	if err := u.repo.Update(ctx, q); err != nil {
		return nil, "", err
	}

	// หลัง Update QRCode
	if err := u.redisrepo.PushWaitingQueue(ctx, q.ID); err != nil {
		return nil, "", fmt.Errorf("failed to push queue to Redis: %v", err)
	}

	return q, qrBase64, nil
}

// func (s *queueUseCase) CallNextQueue(ctx context.Context) (*entity.Queue, error) {
// 	queue, err := s.repo.FindNextWaiting(ctx)
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return nil, errors.New("no waiting queue found")
// 		}
// 		return nil, err
// 	}

// 	queue.Status = "called"
// 	if err := s.repo.Update(ctx, queue); err != nil {
// 		return nil, err
// 	}

// 	return queue, nil
// }

func (s *queueUseCase) CallNextQueue(ctx context.Context) (*entity.Queue, error) {
	// 1. Pop next queue ID จาก Redis
	queueID, err := s.redisrepo.PopNextWaiting(ctx)
	if err != nil {
		if err == redis.Nil {
			return nil, errors.New("no waiting queue found")
		}
		return nil, err
	}

	// 2. ดึงข้อมูล queue จาก Postgres
	queue, err := s.repo.FindByID(ctx, queueID)
	if err != nil {
		return nil, err
	}

	// 3. อัพเดทสถานะ
	queue.Status = "called"
	if err := s.repo.Update(ctx, queue); err != nil {
		return nil, err
	}

	return queue, nil
}

func (u *queueUseCase) FindQueueByName(ctx context.Context, name string) ([]entity.Queue, error) {
	return u.repo.FindQueueByName(ctx, name)
}

func generateQueueNo() string {
	return fmt.Sprintf("Q-%d", time.Now().UnixNano())
}
