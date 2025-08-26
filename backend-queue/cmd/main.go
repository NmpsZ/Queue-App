package main

import (
	"backend-queue/config"
	"backend-queue/database"
	router "backend-queue/routers"
	"fmt"
)

func main() {
	cfg := config.LoadConfig()

	db := database.InitPostgres(cfg)
	database.InitRedis(cfg)

	// // Auto migrate
	// database.DB.AutoMigrate(&queue.Queue{})

	port := ":3000"
	// app := fiber.New()

	// สร้าง RouterConfig แล้วเรียก SetupRouter
	rc := &router.RouterConfig{
		Config: cfg,
		GormDB: db,
	}
	app := rc.SetupRouter()

	// DI setup
	// repo := queue.NewRepository()
	// service := queue.NewService(repo)
	// handler := queue.NewHandler(service)

	// Routes
	// app.Get("/queue", handler.GetQueues)
	// app.Post("/queue", handler.AddQueue)
	// app.Post("/queue/call", handler.CallNextQueue)

	fmt.Println("Server running in port", port)
	if err := app.Listen(port); err != nil {
		fmt.Print("Error Starting Server", port)
	}
}
