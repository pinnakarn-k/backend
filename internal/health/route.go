package health

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(app *fiber.App, handler *Handler) {
	app.Get("/healthz", handler.Check)
	app.Get("/readyz", handler.Check)
}
