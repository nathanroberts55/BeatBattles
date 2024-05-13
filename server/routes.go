package server

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/nathanroberts55/beatbattle/common"
	"github.com/nathanroberts55/beatbattle/controllers"
	"github.com/nathanroberts55/beatbattle/controllers/watch"
)

func buildHandler(app *common.App, controller func(*common.Ctx) error) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return controller(&common.Ctx{
			c,
			app,
		})
	}
}

func RegisterRoutes(a *common.App) {
	app := a.Server

	// Routes
	app.Get("/", buildHandler(a, controllers.HomeIndex))
	app.Get("/watch", buildHandler(a, watch.Index))
	app.Get("/watch/:streamer", buildHandler(a, watch.Show))
	app.Get("/ws/watch/:streamer", websocket.New(func(c *websocket.Conn) {
		watch.Watch(a, c)
	}))

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
}
