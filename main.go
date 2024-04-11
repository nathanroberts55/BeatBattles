package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/nathanroberts55/beatbattle/controllers"
	"github.com/nathanroberts55/beatbattle/initializers"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDatabase()
	initializers.SyncDB()
}

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

	// Start Twitch client in a separate goroutine
	go initializers.ConnectToTwitch("aspecticor")

	// Start App
	app.Listen(":8080")
}
