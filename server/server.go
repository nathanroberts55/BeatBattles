package server

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
)

func New() *fiber.App {
	app := fiber.New(fiber.Config{
		Views:       html.New("./views", ".html"),
		ViewsLayout: "layouts/base",
	})

	// Logging remote IP and Port
	app.Use(logger.New(logger.Config{
		Format:     "${time} | [${ip}:${port}] ${status} - ${method} ${path}\n",
		TimeFormat: "02-Jan-2006 15:04:05",
		TimeZone:   "UTC",
		Output:     os.Stdout,
	}))

	return app
}
