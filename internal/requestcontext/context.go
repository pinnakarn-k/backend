package requestcontext

import "github.com/gofiber/fiber/v2"

const localsKey = "requestContext"

// Example
// Key=cookie value=user_id=1001; username=somchai; session_id=abc123; account_no=123456; branch=001

type RequestContext struct {
	UserID    string
	Username  string
	SessionID string

	// เปิดช่องไว้สำหรับอนาคต
	Values map[string]string
}

func Set(c *fiber.Ctx, ctx RequestContext) {
	c.Locals(localsKey, ctx)
}

func Get(c *fiber.Ctx) RequestContext {
	v := c.Locals(localsKey)
	if v == nil {
		return RequestContext{
			Values: map[string]string{},
		}
	}

	ctx, ok := v.(RequestContext)
	if !ok {
		return RequestContext{
			Values: map[string]string{},
		}
	}

	return ctx
}
