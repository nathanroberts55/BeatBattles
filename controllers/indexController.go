package controllers

import "github.com/gofiber/fiber/v2"

type indexProps struct {
	Params
}

func HomeIndex(c *fiber.Ctx) error {
	props := indexProps{
		defaultParams,
	}
	return c.Render("index", props)
}
