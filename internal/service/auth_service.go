package service

import (
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"

	"agodrift/internal/config"
	"agodrift/internal/model"
	"agodrift/internal/repository"
)

// AuthService handles authentication and token generation.
type AuthService struct {
	users  repository.UserRepository
	secret []byte
	// simple in-memory blacklist of token jtis -> expiry
	blacklist map[string]int64
	mu        sync.Mutex
}

// DefaultAuth is a shared singleton used by handlers and middleware
var DefaultAuth *AuthService

// InitDefaultAuth initializes the default AuthService singleton
func InitDefaultAuth(secret string) {
	if DefaultAuth == nil {
		DefaultAuth = NewAuthService(secret)
	}
}

func NewAuthService(secret string) *AuthService {
	return &AuthService{
		users:     repository.NewMySQLUserRepo(config.GetDB()),
		secret:    []byte(secret),
		blacklist: make(map[string]int64),
	}
}

// GetAuth returns the singleton if initialized otherwise creates one with default secret
func GetAuth() *AuthService {
	if DefaultAuth == nil {
		InitDefaultAuth("changeme")
	}
	return DefaultAuth
}

// Authenticate checks username/password and returns user if valid.
func (s *AuthService) Authenticate(username, password string) (model.User, bool) {
	u, ok := s.users.GetByUsername(username)
	if !ok {
		return model.User{}, false
	}
	// plain password check for demo - replace with hashed compare in prod
	if u.Password != password {
		return model.User{}, false
	}
	return u, true
}

// CreateToken generates a JWT with role claim and jti.
func (s *AuthService) CreateToken(u model.User, ttl time.Duration) (string, error) {
	jti := uuid.NewString()
	claims := jwt.MapClaims{
		"sub":  u.Username,
		"role": u.Role,
		"jti":  jti,
		"exp":  time.Now().Add(ttl).Unix(),
		"iat":  time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

// BlacklistToken marks a token's jti as revoked until expiry
func (s *AuthService) BlacklistToken(jti string, exp int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.blacklist[jti] = exp
}

// IsBlacklisted checks whether a jti is revoked
func (s *AuthService) IsBlacklisted(jti string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	exp, ok := s.blacklist[jti]
	if !ok {
		return false
	}
	if time.Now().Unix() > exp {
		// expired - remove
		delete(s.blacklist, jti)
		return false
	}
	return true
}
