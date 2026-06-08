package health

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(router fiber.Router, handler *Handler) {
	router.Get("/healthz", handler.Health)
	router.Get("/readyz", handler.Ready)
}
