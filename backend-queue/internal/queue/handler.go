package queue

import (
	"backend-queue/entity"
	"backend-queue/usecase/queue"
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/redis/go-redis/v9"
)

type QueueHandler struct {
	service    queue.UseCase
	timeNowUTC func() time.Time
	RedisDB    *redis.Client
}

func NewQueueHandler(service queue.UseCase) QueueHandler {
	return QueueHandler{
		service:    service,
		timeNowUTC: func() time.Time { return time.Now().UTC() },
	}
}

func RegisterQueueWS(app *fiber.App, redisClient *redis.Client) {
	app.Get("/queue", websocket.New(func(c *websocket.Conn) {
		defer c.Close()

		queueNo := c.Query("queue_no")
		log.Printf("WebSocket connected for queue: %s", queueNo)

		ctx := context.Background()
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				// ดึง current queue จาก Redis
				currentID, err := redisClient.LIndex(ctx, "queue:waiting", -1).Result()
				if err != nil && err != redis.Nil {
					log.Printf("Redis error: %v", err)
					return
				}

				totalQueue, _ := redisClient.LLen(ctx, "queue:waiting").Result()

				// ถ้า currentID เป็น string ต้องแปลงเป็น int
				currentQueue := 0
				if currentID != "" {
					currentQueue, _ = strconv.Atoi(currentID)
				}

				msg := fmt.Sprintf(`{"current_queue":%d,"total_queue":%d}`, currentQueue, totalQueue)
				if err := c.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
					log.Printf("WebSocket write error: %v", err)
					return
				}
			}
		}
	}))
}

// Get all queues
func (h *QueueHandler) GetQueues(c *fiber.Ctx) error {
	log.Printf("[QueueHandler] GetQueues called")
	data, err := h.service.FindQueue(c.Context())
	if err != nil {
		log.Printf("[QueueHandler] GetQueues error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	log.Printf("[QueueHandler] GetQueues passed, returned %d queues", len(data))
	return c.Status(fiber.StatusOK).JSON(data)
}

// Get queue by ID
func (h *QueueHandler) GetQueueByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		log.Printf("[QueueHandler] GetQueueByID invalid ID: %s", idParam)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	data, err := h.service.GetQueueByID(c.Context(), uint(id))
	if err != nil {
		log.Printf("[QueueHandler] GetQueueByID not found ID: %d", id)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Queue not found"})
	}
	log.Printf("[QueueHandler] GetQueueByID passed, returned ID: %d", id)
	return c.Status(fiber.StatusOK).JSON(data)
}

// Create a new queue
func (h *QueueHandler) CreateQueue(c *fiber.Ctx) error {
	var q entity.Queue
	if err := c.BodyParser(&q); err != nil {
		log.Printf("[QueueHandler] CreateQueue body parse error: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.service.AddQueue(c.Context(), &q); err != nil {
		log.Printf("[QueueHandler] CreateQueue error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	log.Printf("[QueueHandler] CreateQueue passed, created ID: %d", q.ID)
	return c.Status(fiber.StatusCreated).JSON(q)
}

// Update an existing queue
func (h *QueueHandler) UpdateQueue(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		log.Printf("[QueueHandler] UpdateQueue invalid ID: %s", idParam)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var q entity.Queue
	if err := c.BodyParser(&q); err != nil {
		log.Printf("[QueueHandler] UpdateQueue body parse error: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	q.ID = uint(id)
	if err := h.service.UpdateQueue(c.Context(), &q); err != nil {
		log.Printf("[QueueHandler] UpdateQueue error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	log.Printf("[QueueHandler] UpdateQueue passed, updated ID: %d", q.ID)
	return c.Status(fiber.StatusOK).JSON(q)
}

// Delete a queue
func (h *QueueHandler) DeleteQueue(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		log.Printf("[QueueHandler] DeleteQueue invalid ID: %s", idParam)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	if err := h.service.DeleteQueue(c.Context(), uint(id)); err != nil {
		log.Printf("[QueueHandler] DeleteQueue error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	log.Printf("[QueueHandler] DeleteQueue passed, deleted ID: %d", id)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Deleted successfully"})
}

func (h *QueueHandler) CreateQueueWithQR(c *fiber.Ctx) error {
	var q entity.Queue

	log.Println("[QueueHandler] CreateQueueWithQR called")

	if err := c.BodyParser(&q); err != nil {
		log.Printf("[QueueHandler] Error parsing request body: %v\n", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	queue, qrBase64, err := h.service.CreateQueueWithQR(c.Context(), &q)
	if err != nil {
		log.Printf("[QueueHandler] Error creating queue: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	log.Printf("[QueueHandler] Queue ID %d created successfully with QR Code\n", queue.ID)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"queue": queue,
		"qr":    qrBase64,
	})
}

func (h *QueueHandler) GetQueueByName(c *fiber.Ctx) error {
	log.Println("[QueueHandler] GetQueueByName called")

	nameParam := c.Params("name")
	if nameParam == "" {
		log.Println("[QueueHandler] GetQueueByName missing name param")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Name is required"})
	}

	data, err := h.service.FindQueueByName(c.Context(), nameParam)
	if err != nil {
		log.Printf("[QueueHandler] GetQueueByName error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Server error"})
	}

	if len(data) == 0 {
		log.Printf("[QueueHandler] No queue found for name: %s", nameParam)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "No queues found"})
	}

	log.Printf("[QueueHandler] GetQueueByName passed, found %d queues for name: %s", len(data), nameParam)
	return c.Status(fiber.StatusOK).JSON(data)
}

func (h *QueueHandler) CallNextQueue(c *fiber.Ctx) error {
	log.Println("[QueueHandler] CallNextQueue called")

	queue, err := h.service.CallNextQueue(c.Context())
	if err != nil {
		if err.Error() == "no waiting queue found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "next queue called",
		"queue":   queue,
	})
}
