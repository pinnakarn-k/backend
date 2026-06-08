package apperror

import "net/http"

var (
	ErrBadRequest = New(
		http.StatusBadRequest,
		"BAD_REQUEST",
		"bad request",
	)

	ErrUnauthorized = New(
		http.StatusUnauthorized,
		"UNAUTHORIZED",
		"unauthorized",
	)

	ErrForbidden = New(
		http.StatusForbidden,
		"FORBIDDEN",
		"forbidden",
	)

	ErrInternalServer = New(
		http.StatusInternalServerError,
		"INTERNAL_SERVER_ERROR",
		"internal server error",
	)

	ErrServiceUnavailable = New(
		http.StatusServiceUnavailable,
		"SERVICE_UNAVAILABLE",
		"service unavailable",
	)
)
