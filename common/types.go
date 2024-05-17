package common

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nathanroberts55/beatbattle/twitch"
)

type App struct {
	Twitch *twitch.TwitchService
	Server *fiber.App
}

type Ctx struct {
	*fiber.Ctx
	App *App
}
