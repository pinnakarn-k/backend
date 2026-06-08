package response

import (
	"errors"

	"backend/internal/apperror"

	"github.com/gofiber/fiber/v2"
)

type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ValidationErrorBody struct {
	Code    string       `json:"code"`
	Message string       `json:"message"`
	Errors  []FieldError `json:"errors"`
}

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func Error(c *fiber.Ctx, err error) error {
	var appErr *apperror.Error

	if errors.As(err, &appErr) {
		if len(appErr.Fields) > 0 {
			return validationError(c, appErr)
		}

		return appError(c, appErr)
	}

	return appError(c, apperror.ErrInternalServer)
}

func appError(c *fiber.Ctx, err *apperror.Error) error {
	return c.Status(err.StatusCode).JSON(ErrorBody{
		Code:    err.Code,
		Message: err.Message,
	})
}

func validationError(c *fiber.Ctx, err *apperror.Error) error {
	return c.Status(err.StatusCode).JSON(ValidationErrorBody{
		Code:    err.Code,
		Message: err.Message,
		Errors:  toFieldErrors(err.Fields),
	})
}

func toFieldErrors(fields []apperror.FieldError) []FieldError {
	result := make([]FieldError, 0, len(fields))

	for _, field := range fields {
		result = append(result, FieldError{
			Field:   field.Field,
			Message: field.Message,
		})
	}

	return result
}
