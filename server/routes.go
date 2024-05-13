package server

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/nathanroberts55/beatbattle/common"
	"github.com/nathanroberts55/beatbattle/controllers"
)

func RegisterRoutes(a *common.App) {
	app := a.Server

	// Routes
	app.Get("/", controllers.HomeIndex)
	app.Get("/watch", controllers.WatchIndex)
	app.Get("/watch/:streamer", controllers.WatchShow)

	// Configure App
	app.Static("/", "./public")

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
		controllers.WatchStream(a, c)
	}))
}
