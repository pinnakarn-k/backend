package middleware

import (
	"backend/internal/requestcontext"

	"github.com/gofiber/fiber/v2"
)

func RequestContext() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := requestcontext.RequestContext{
			UserID:    c.Cookies("user_id"),
			Username:  c.Cookies("username"),
			SessionID: c.Cookies("session_id"),
			Values: map[string]string{
				"account_no": c.Cookies("account_no"),
				"branch":     c.Cookies("branch"),
			},
		}

		requestcontext.Set(c, ctx)

		return c.Next()
	}
}
