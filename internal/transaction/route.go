package transaction

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(router fiber.Router, handler *Handler) {
	router.Get("/transaction", handler.Search)
	router.Post("/transactions/download", handler.Download)
	router.Post("/transactions/sendemail", handler.SendEmail)
}
