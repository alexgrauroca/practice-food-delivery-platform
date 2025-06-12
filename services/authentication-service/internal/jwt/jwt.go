// Package jwt provides functionality for generating JWT tokens with custom claims and expiration.
package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// DefaultTokenType represents the default type of token used for authorization, typically set to "Bearer".
const DefaultTokenType = "Bearer"

// Service defines the interface for JWT token operations
//
//go:generate mockgen -destination=./mocks/service_mock.go -package=jwt_mocks github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/jwt Service
type Service interface {
	GenerateToken(id string, cfg Config) (string, error)
}

// Claims represents a JWT payload.
type Claims struct {
	Subject   string    // "sub" claim
	Role      string    // "role" claim
	ExpiresAt time.Time // "exp" claim
}

// Config represents configuration settings for token generation
type Config struct {
	Expiration int // AccessToken expiration duration in seconds
	Role       string
}

type service struct {
	secret []byte // jwt secret key
}

// NewService creates a new JWT service instance
func NewService(secret []byte) Service {
	return &service{
		secret: secret,
	}
}

func (s *service) GenerateToken(id string, cfg Config) (string, error) {
	claims := jwt.MapClaims{
		"sub":  id,
		"role": cfg.Role,
		"exp":  time.Now().Add(time.Duration(cfg.Expiration) * time.Second).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}
