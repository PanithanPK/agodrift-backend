package api

import (
	"agodrift/internal/api/handlers"
	"agodrift/internal/config"
	"agodrift/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

// NewApp builds and returns the Fiber app used by the server.
func NewApp() *fiber.App {
	app := fiber.New()
	secret := config.Get("JWT_SECRET", "changeme")

	// health check
	app.Get("/api/v1/health", handlers.Health)

	// auth routes
	app.Post("/api/v1/auth/login", handlers.Login)
	// protected routes
	app.Post("/api/v1/auth/logout", middleware.JWTConfig(secret), handlers.Logout)
	app.Get("/api/v1/auth/me", middleware.JWTConfig(secret), handlers.Me)

	// room routes
	app.Get("/api/v1/listrooms", handlers.ListRoomsHandler)
	app.Get("/api/v1/listrooms/:id", handlers.RoomByIDHandler)

	// require admin role to create room
	app.Post("/api/v1/AddRoom", middleware.JWTConfig(secret), middleware.RequireRole("admin"), handlers.AddRoomHandler)

	return app
}
