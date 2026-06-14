package response

import (
	"backend/internal/pagination"

	"github.com/gofiber/fiber/v2"
)

type SuccessBody struct {
	Data any   `json:"data"`
	Meta *Meta `json:"meta,omitempty"`
}

type Meta struct {
	Pagination *pagination.Pagination `json:"pagination,omitempty"`
}

type Pagination struct {
	Page       int `json:"page"`
	PerPage    int `json:"perPage"`
	Total      int `json:"total"`
	TotalPages int `json:"totalPages"`
}

func Success(c *fiber.Ctx, data any) error {
	return c.JSON(SuccessBody{
		Data: data,
	})
}

func SuccessWithPagination(
	c *fiber.Ctx,
	data any,
	p pagination.Pagination,
) error {
	return c.JSON(SuccessBody{
		Data: data,
		Meta: &Meta{
			Pagination: &p,
		},
	})
}
