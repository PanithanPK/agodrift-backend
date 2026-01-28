package middleware

import "github.com/gofiber/fiber/v2"

// DummyAuth is a placeholder auth middleware.
func DummyAuth(c *fiber.Ctx) error {
	// TODO: validate token / set user in context
	return c.Next()
}
