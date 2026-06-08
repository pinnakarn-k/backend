package main

import (
	"os"
	"os/signal"
	"syscall"

	"backend/internal/app"
	"backend/internal/config"
	"backend/internal/logger"
)

func main() {
	cfg := config.Load()

	logger := logger.New()

	fiberApp := app.New(logger)

	go func() {
		if err := fiberApp.Listen(cfg.ListenAddress()); err != nil {
			logger.Error(
				"failed to start server",
				"error", err,
			)

			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(
		quit,
		os.Interrupt,
		syscall.SIGTERM,
	)

	<-quit

	logger.Info("shutting down server")

	if err := fiberApp.Shutdown(); err != nil {
		logger.Error(
			"failed to shutdown server",
			"error", err,
		)

		os.Exit(1)
	}

	logger.Info("server stopped")
}
