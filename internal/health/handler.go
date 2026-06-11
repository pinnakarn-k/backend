package health

import (
	"backend/internal/requestcontext"
	"backend/internal/response"
	"fmt"

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
	reqCtx := requestcontext.Get(c)
	fmt.Printf("%+v\n", reqCtx)

	return response.Success(c, fiber.Map{
		"status": "ok",
	})
}

func (h *Handler) Ready(c *fiber.Ctx) error {
	reqCtx := requestcontext.Get(c)
	fmt.Printf("%+v\n", reqCtx)

	if err := h.service.Ready(reqCtx); err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, fiber.Map{
		"status": "ready",
	})
}
