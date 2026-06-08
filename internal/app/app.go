package app

import (
	"backend/internal/health"
	"backend/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func New() *fiber.App {
	fiberApp := fiber.New()

	fiberApp.Use(middleware.RequestID())

	healthRepo := health.NewRepository()
	healthService := health.NewService(healthRepo)
	healthHandler := health.NewHandler(healthService)

	health.RegisterRoutes(fiberApp, healthHandler)

	return fiberApp
}
