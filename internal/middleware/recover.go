package middleware

import (
	"log/slog"

	"backend/internal/response"

	"github.com/gofiber/fiber/v2"
)

func Recover(log *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				log.Error(
					"panic recovered",
					"requestId", GetRequestID(c),
					"error", r,
				)

				_ = response.InternalServerError(c)
			}
		}()

		return c.Next()
	}
}
