package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const RequestIDKey = "request_id"
const RequestIDHeader = "X-Request-Id"

func RequestID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		requestID := c.Get(RequestIDHeader)

		if requestID == "" {
			requestID = uuid.NewString()
		}

		c.Locals(RequestIDKey, requestID)
		c.Set(RequestIDHeader, requestID)

		return c.Next()
	}
}

func GetRequestID(c *fiber.Ctx) string {
	requestID, _ := c.Locals(RequestIDKey).(string)
	return requestID
}
