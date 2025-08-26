package router

import (
	"backend-queue/config"
	"backend-queue/database"
	"backend-queue/infra/postgres"
	redisDB "backend-queue/infra/postgres/redis"
	"backend-queue/internal/queue"

	"github.com/redis/go-redis/v9"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"gorm.io/gorm"
)

type RouterConfig struct {
	GormDB  *gorm.DB
	RedisDB *redis.Client
	Config  *config.Config
}

// สร้าง Fiber app พร้อม DI
func (rc *RouterConfig) SetupRouter() *fiber.App {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // หรือใส่ URL ของ frontend เช่น "http://localhost:5173"
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// --- DI Setup ---
	var db *gorm.DB
	if rc.GormDB != nil {
		db = rc.GormDB
	} else {
		db = database.InitPostgres(rc.Config)
	}

	var redisdb *redis.Client
	if rc.RedisDB != nil {
		redisdb = rc.RedisDB
	} else {
		redisdb = database.InitRedis(rc.Config)
	}

	// Repository
	queueRepo := postgres.NewQueueRepository(db)
	redisRepo := redisDB.NewQueueRedisRepo(redisdb)

	// UseCase
	queueService := queue.NewQueueUseCase(queueRepo, redisRepo)

	// Handler
	queueHandler := queue.NewQueueHandler(queueService)

	queue.RegisterQueueWS(app, redisdb)

	// --- Routes ---
	api := app.Group("/api")
	queueGroup := api.Group("/queue")

	queueGroup.Get("/", queueHandler.GetQueues)
	queueGroup.Get("/:id", queueHandler.GetQueueByID)
	queueGroup.Get("/user/:name", queueHandler.GetQueueByName)
	queueGroup.Post("/", queueHandler.CreateQueue)
	queueGroup.Post("/qr", queueHandler.CreateQueueWithQR)
	queueGroup.Put("/next", queueHandler.CallNextQueue)
	queueGroup.Put("/:id", queueHandler.UpdateQueue)
	queueGroup.Delete("/:id", queueHandler.DeleteQueue)

	return app
}
