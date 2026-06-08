package response

import (
	"errors"

	"backend/internal/apperror"

	"github.com/gofiber/fiber/v2"
)

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ErrorBody struct {
	Code    string       `json:"code"`
	Message string       `json:"message"`
	Errors  []FieldError `json:"errors,omitempty"`
}

func Error(c *fiber.Ctx, err error) error {
	var appErr *apperror.Error

	if errors.As(err, &appErr) {
		return c.Status(appErr.StatusCode).JSON(ErrorBody{
			Code:    appErr.Code,
			Message: appErr.Message,
			Errors:  toFieldErrors(appErr.Fields),
		})
	}

	return c.Status(apperror.ErrInternalServer.StatusCode).JSON(ErrorBody{
		Code:    apperror.ErrInternalServer.Code,
		Message: apperror.ErrInternalServer.Message,
	})
}

func toFieldErrors(fields []apperror.FieldError) []FieldError {
	if len(fields) == 0 {
		return nil
	}

	result := make([]FieldError, 0, len(fields))

	for _, field := range fields {
		result = append(result, FieldError{
			Field:   field.Field,
			Message: field.Message,
		})
	}

	return result
}
