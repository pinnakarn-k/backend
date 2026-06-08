package validator

import (
	"errors"
	"reflect"
	"strings"

	"backend/internal/apperror"

	govalidator "github.com/go-playground/validator/v10"
)

var validate = newValidator()

func newValidator() *govalidator.Validate {
	v := govalidator.New()

	v.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(
			field.Tag.Get("json"),
			",",
			2,
		)[0]

		if name == "-" {
			return ""
		}

		return name
	})

	return v
}

func Validate(input any) error {
	if err := validate.Struct(input); err != nil {
		var validationErrors govalidator.ValidationErrors

		if errors.As(err, &validationErrors) {
			fields := make([]apperror.FieldError, 0, len(validationErrors))

			for _, fieldErr := range validationErrors {
				fields = append(fields, apperror.FieldError{
					Field:   fieldErr.Field(),
					Message: messageOf(fieldErr),
				})
			}

			return apperror.Validation(fields)
		}

		return err
	}

	return nil
}

func messageOf(err govalidator.FieldError) string {
	switch err.Tag() {
	case "required":
		return err.Field() + " is required"

	case "email":
		return err.Field() + " must be a valid email"

	case "min":
		return err.Field() + " must be at least " + err.Param()

	case "max":
		return err.Field() + " must be at most " + err.Param()

	default:
		return err.Field() + " is invalid"
	}
}
