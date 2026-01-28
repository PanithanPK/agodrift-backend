package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	fiberjwt "github.com/gofiber/jwt/v2"
	"github.com/golang-jwt/jwt/v4"

	"agodrift/internal/service"
)

// JWTConfig returns a Fiber middleware that validates JWT and checks blacklist
func JWTConfig(secret string) fiber.Handler {
	// ensure singleton auth using provided secret
	service.InitDefaultAuth(secret)
	cfg := fiberjwt.New(fiberjwt.Config{
		SigningKey: []byte(secret),
		ContextKey: "user",
		SuccessHandler: func(c *fiber.Ctx) error {
			// check blacklist
			if u := c.Locals("user"); u != nil {
				if tok, ok := u.(*jwt.Token); ok {
					if claims, ok := tok.Claims.(jwt.MapClaims); ok {
						jti, _ := claims["jti"].(string)
						if service.GetAuth().IsBlacklisted(jti) {
							return c.Status(fiber.StatusUnauthorized).SendString("token revoked")
						}
					}
				}
			}
			return c.Next()
		},
	})
	return func(c *fiber.Ctx) error {
		// Allow Authorization: Bearer <token>
		// fiberjwt will parse and validate token for us
		if err := cfg(c); err != nil {
			// fiberjwt returns an error which should already be a response
			// Map some errors to 401
			// If it's an unauthorized error, forward a 401
			if strings.Contains(err.Error(), "Missing or malformed JWT") || strings.Contains(err.Error(), "token is expired") {
				return c.Status(fiber.StatusUnauthorized).SendString(err.Error())
			}
			return c.Status(fiber.StatusUnauthorized).SendString(err.Error())
		}
		return nil
	}
}

// RequireRole returns middleware that ensures the token has the given role
func RequireRole(role string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user")
		if user == nil {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		if tok, ok := user.(*jwt.Token); ok {
			if claims, ok := tok.Claims.(jwt.MapClaims); ok {
				r, _ := claims["role"].(string)
				if r != role {
					return c.Status(fiber.StatusForbidden).SendString("forbidden")
				}
			}
		}
		return c.Next()
	}
}
