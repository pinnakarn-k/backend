package apperror

type FieldError struct {
	Field   string
	Message string
}

type Error struct {
	StatusCode int
	Code       string
	Message    string
	Fields     []FieldError
}

func New(statusCode int, code string, message string) *Error {
	return &Error{
		StatusCode: statusCode,
		Code:       code,
		Message:    message,
	}
}

func Validation(fields []FieldError) *Error {
	return &Error{
		StatusCode: 400,
		Code:       "VALIDATION_ERROR",
		Message:    "validation error",
		Fields:     fields,
	}
}

func (e *Error) Error() string {
	return e.Message
}
