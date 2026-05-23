package middleware

import "github.com/gofiber/fiber/v3"

func JSONContentType() fiber.Handler {
	return func(c fiber.Ctx) error {
		c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
		return c.Next()
	}
}
