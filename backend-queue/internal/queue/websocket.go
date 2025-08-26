package queue

// import (
// 	"fmt"

// 	"github.com/gofiber/fiber/v2"
// 	"github.com/gofiber/websocket/v2"
// )

// func RegisterQueueWS(app *fiber.App) {
// 	// WebSocket route
// 	app.Get("/queue", websocket.New(func(c *websocket.Conn) {
// 		defer c.Close()

// 		// อ่าน query param
// 		queueNo := c.Query("queue_no")
// 		fmt.Println("New WebSocket connection for queue:", queueNo)

// 		// ตัวอย่างส่ง message ทุก 5 วินาที
// 		for i := 0; i <= 10; i++ {
// 			message := fmt.Sprintf(`{"current_queue":%d,"total_queue":20}`, i)
// 			if err := c.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
// 				fmt.Println("Write error:", err)
// 				break
// 			}
// 		}
// 	}))
// }
