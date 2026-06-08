package health

import (
	"backend/internal/response"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Health(c *fiber.Ctx) error {
	return response.Success(c, fiber.Map{
		"status": "ok",
	})
}

func (h *Handler) Ready(c *fiber.Ctx) error {
	if err := h.service.Ready(); err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, fiber.Map{
		"status": "ready",
	})
}
