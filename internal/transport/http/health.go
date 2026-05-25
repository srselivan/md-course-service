package http

import (
	"context"

	"github.com/gofiber/fiber/v3"
)

type Pinger interface {
	Ping(ctx context.Context) error
}

func (s *Server) health(ctx fiber.Ctx) error {
	for _, p := range s.pingers {
		if err := p.Ping(ctx.Context()); err != nil {
			s.logger.Error().Err(err).Msg("health check failed")
			return fiber.NewError(fiber.StatusServiceUnavailable, err.Error())
		}
	}

	return ctx.SendStatus(fiber.StatusOK)
}
