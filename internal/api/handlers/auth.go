package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"

	"agodrift/internal/config"
	"agodrift/internal/model"
	"agodrift/internal/service"
)

// use the shared AuthService singleton
var authService = service.GetAuth()

// LoginRequest is the body for login
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(c *fiber.Ctx) error {
	// ensure default auth is initialized with configured secret
	service.InitDefaultAuth(config.Get("JWT_SECRET", "changeme"))
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("invalid body")
	}
	u, ok := authService.Authenticate(req.Username, req.Password)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).SendString("invalid credentials")
	}
	// create token (30m)
	token, err := authService.CreateToken(u, 30*time.Minute)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("failed to create token")
	}
	return c.JSON(fiber.Map{"token": token})
}

func Logout(c *fiber.Ctx) error {
	// extract token jti and exp from token in context set by jwt middleware
	user := c.Locals("user")
	if user == nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	if tok, ok := user.(*jwt.Token); ok {
		if claims, ok := tok.Claims.(jwt.MapClaims); ok {
			jti, _ := claims["jti"].(string)
			expFloat, _ := claims["exp"].(float64)
			exp := int64(expFloat)
			authService.BlacklistToken(jti, exp)
			return c.SendStatus(fiber.StatusOK)
		}
	}
	return c.SendStatus(fiber.StatusBadRequest)
}

func Me(c *fiber.Ctx) error {
	user := c.Locals("user")
	if user == nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	if tok, ok := user.(*jwt.Token); ok {
		if claims, ok := tok.Claims.(jwt.MapClaims); ok {
			username, _ := claims["sub"].(string)
			role, _ := claims["role"].(string)
			return c.JSON(model.User{Username: username, Role: role})
		}
	}
	return c.SendStatus(fiber.StatusBadRequest)
}
