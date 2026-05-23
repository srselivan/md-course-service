package http

import (
	"github.com/gofiber/contrib/v3/swaggo"
	"github.com/gofiber/fiber/v3"

	_ "course-service/docs"
)

func registerSwagger(app *fiber.App) {
	app.Get("/swagger/*", swaggo.HandlerDefault)
}
