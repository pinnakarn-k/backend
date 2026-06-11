package main

import (
	"os"
	"os/signal"
	"syscall"

	"backend/internal/app"
	"backend/internal/config"
	"backend/internal/logger"
	redisinfra "backend/internal/redis"
)

func main() {
	cfg := config.Load()

	logger := logger.New(
		cfg.Service,
		cfg.Env,
	)

	redisClient, err := redisinfra.New(
		redisinfra.Config{
			Enabled:  cfg.RedisEnabled,
			Host:     cfg.RedisHost,
			Port:     cfg.RedisPort,
			Password: cfg.RedisPassword,
			DB:       cfg.RedisDB,
		},
	)
	if err != nil {
		logger.Error(
			"failed to connect redis",
			"error", err,
		)

		os.Exit(1)
	}

	defer func() {
		if err := redisClient.Close(); err != nil {
			logger.Error(
				"failed to close redis",
				"error", err,
			)
		}
	}()

	fiberApp := app.New(logger, redisClient)

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
