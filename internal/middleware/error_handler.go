package middleware

import (
	"errors"
	"log/slog"

	"backend/internal/apperror"
	"backend/internal/response"

	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(log *slog.Logger) fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		statusCode := fiber.StatusInternalServerError

		var appErr *apperror.Error
		if errors.As(err, &appErr) {
			statusCode = appErr.StatusCode
		}

		if statusCode >= 500 {
			log.Error(
				"request failed",
				"requestId", GetRequestID(c),
				"method", c.Method(),
				"path", c.Path(),
				"status", statusCode,
				"error", err,
			)
		}

		return response.Error(c, err)
	}
}
