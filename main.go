package main

import (
	"log"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/nathanroberts55/beatbattle/controllers"
	"github.com/nathanroberts55/beatbattle/database"
	"github.com/nathanroberts55/beatbattle/initializers"
	"github.com/nathanroberts55/beatbattle/services"
)

func init() {
	initializers.LoadEnvVariables()
	database.ConnectToDatabase()
	database.SyncDB()
}

type Listener struct {
	conn     *websocket.Conn
	streamer string
}

var listeners []Listener // this would store all active connections

func main() {
	// Create a new engine
	engine := html.New("./views", ".html")

	// Setup App
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// Configure App
	app.Static("/", "./public")

	// Routes
	app.Get("/", controllers.HomeIndex)

	app.Use("/ws", func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws/:id", websocket.New(func(c *websocket.Conn) {

		streamer := c.Params("id")
		listeners = append(listeners, Listener{conn: c, streamer: streamer}) // add the new connection to the list of listeners
		defer func() {
			// remove the connection from listeners when it's closed
			for i, listener := range listeners {
				if listener.conn == c {
					listeners = append(listeners[:i], listeners[i+1:]...)
					break
				}
			}
		}()

		// c.Locals is added to the *websocket.Conn
		log.Println(c.Locals("allowed"))  // true
		log.Println(c.Params("id"))       // 123
		log.Println(c.Query("v"))         // 1.0
		log.Println(c.Cookies("session")) // ""

		// websocket.Conn bindings https://pkg.go.dev/github.com/fasthttp/websocket?tab=doc#pkg-index
		var (
			mt  int
			msg []byte
			err error
		)
		for {
			if mt, msg, err = c.ReadMessage(); err != nil {
				log.Println("read:", err)
				break
			}
			log.Printf("recv: %s", msg)

			if err = c.WriteMessage(mt, msg); err != nil {
				log.Println("write:", err)
				break
			}
		}

	}))

	// Start Twitch client in a separate goroutine
	go services.ConnectToTwitch("pointcrow")

	// Start App
	// app.Listen(":8080")
}
