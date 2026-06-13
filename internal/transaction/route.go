package transaction

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(router fiber.Router, handler *Handler) {
	router.Get("/transaction", handler.Search)
}
