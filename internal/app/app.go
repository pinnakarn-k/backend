package app

import (
	"backend/internal/health"

	"github.com/gofiber/fiber/v2"
)

func New() *fiber.App {
	fiberApp := fiber.New()

	healthRepo := health.NewRepository()
	healthService := health.NewService(healthRepo)
	healthHandler := health.NewHandler(healthService)

	health.RegisterRoutes(fiberApp, healthHandler)

	return fiberApp
}
