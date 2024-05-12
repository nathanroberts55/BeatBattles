package controllers

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

type indexProps struct {
	Params
	Year  string
	Title string
}

func HomeIndex(c *fiber.Ctx) error {
	props := indexProps{
		defaultParams,
		fmt.Sprint(time.Now().Year()),
		"Beat Battle",
	}
	return c.Render("index", props)
}
