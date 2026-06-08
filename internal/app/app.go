package app

import (
	"backend/internal/health"
	"backend/internal/middleware"
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

func New(log *slog.Logger) *fiber.App {
	fiberApp := fiber.New()

	fiberApp.Use(middleware.RequestID())
	fiberApp.Use(middleware.Logger(log))
	fiberApp.Use(middleware.Recover(log))

	healthRepo := health.NewRepository()
	healthService := health.NewService(healthRepo)
	healthHandler := health.NewHandler(healthService)

	health.RegisterRoutes(fiberApp, healthHandler)

	return fiberApp
}
