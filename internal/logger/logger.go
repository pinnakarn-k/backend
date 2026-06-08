package logger

import (
	"log/slog"
	"os"
)

func New(service string, env string) *slog.Logger {
	return slog.New(
		slog.NewJSONHandler(os.Stdout, nil),
	).With(
		"service", service,
		"env", env,
	)
}
