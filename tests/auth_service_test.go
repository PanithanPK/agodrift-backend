package tests

import (
	"testing"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"

	"agodrift/internal/service"
)

func TestAuthTokenLifecycle(t *testing.T) {
	service.InitDefaultAuth("testsecret")
	a := service.GetAuth()
	// Authenticate seeded user
	u, ok := a.Authenticate("admin", "adminpass")
	if !ok {
		t.Fatalf("expected admin to authenticate")
	}
	// Create token with short ttl
	tok, err := a.CreateToken(u, 1*time.Minute)
	if err != nil {
		t.Fatalf("failed to create token: %v", err)
	}
	if tok == "" {
		t.Fatalf("empty token")
	}
	// Parse unverified to read claims
	var claims jwt.MapClaims
	parser := new(jwt.Parser)
	parsed, _, err := parser.ParseUnverified(tok, &claims)
	if err != nil {
		t.Fatalf("failed parse token: %v", err)
	}
	if parsed == nil {
		t.Fatalf("parsed token is nil")
	}
	jti, ok := claims["jti"].(string)
	if !ok || jti == "" {
		t.Fatalf("jti missing")
	}
	expF, ok := claims["exp"].(float64)
	if !ok {
		t.Fatalf("exp missing or wrong type")
	}
	exp := int64(expF)
	a.BlacklistToken(jti, exp)
	if !a.IsBlacklisted(jti) {
		t.Fatalf("jti should be blacklisted")
	}
}
