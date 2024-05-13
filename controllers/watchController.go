package controllers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type watchProps struct {
	Params
	Streamer string
}

func WatchIndex(c *fiber.Ctx) error {
	streamer := c.Query("streamer", "")
	if len(streamer) > 0 {
		c.Redirect(fmt.Sprintf("/watch/%s", streamer))
		return nil
	}

	return c.Render("watch/index", defaultParams)
}

func WatchShow(c *fiber.Ctx) error {
	streamer := c.Params("streamer", "ttlnow")
	props := watchProps{
		defaultParams,
		streamer,
	}

	return c.Render("watch/show", props, "layouts/watch")
}
