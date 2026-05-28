package http

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/adaptor"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	_ "course-service/docs"
)

func registerSwagger(app *fiber.App) {
	app.Get("/swagger/*", adaptor.HTTPHandler(
		httpSwagger.Handler(
			httpSwagger.URL("/swagger/doc.json"),
		),
	))
}
