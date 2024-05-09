package main

import (
	"github.com/nathanroberts55/beatbattle/common"
	"github.com/nathanroberts55/beatbattle/initializers"
	"github.com/nathanroberts55/beatbattle/server"
	"github.com/nathanroberts55/beatbattle/twitch"
)

func init() {
	initializers.LoadEnvVariables()
}

func main() {
	app := common.App{
		Server: server.New(),
		Twitch: twitch.New(),
	}

	// Start Twitch client in a separate goroutine
	go app.Twitch.Start()

	// Setup App
	server.RegisterRoutes(&app)
	app.Server.Listen(":8080")
}
