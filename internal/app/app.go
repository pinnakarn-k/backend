package app

import (
	"backend/internal/health"
	"backend/internal/middleware"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/redis/go-redis/v9"
)

func New(log *slog.Logger, redisClient *redis.Client) *fiber.App {
	fiberApp := fiber.New()

	fiberApp.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173",
		AllowMethods: "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization, X-Request-Id",
	}))

	fiberApp.Use(middleware.RequestContext())
	fiberApp.Use(middleware.RequestID())
	fiberApp.Use(middleware.Logger(log))
	fiberApp.Use(middleware.Recover(log))

	healthRepo := health.NewRepository()
	healthService := health.NewService(healthRepo)
	healthHandler := health.NewHandler(healthService)

	health.RegisterRoutes(fiberApp, healthHandler)

	return fiberApp
}
