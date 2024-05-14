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
		Server:   server.New(),
		Twitch:   twitch.New(),
		Sessions: server.NewSessionStore(),
	}

	// Start Twitch client in a separate goroutine
	go app.Twitch.Start()

	// Setup App
	server.RegisterRoutes(&app)
	err := app.Server.Listen(":8080")
	if err != nil {
		panic(err)
	}
}
