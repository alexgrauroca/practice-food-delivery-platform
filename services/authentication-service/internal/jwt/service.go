// Package jwt provides functionality for generating JWT tokens with custom claims and expiration.
package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/alexgrauroca/practice-food-delivery-platform/pkg/log"
)

// DefaultTokenType represents the default type of token used for authorization, typically set to "Bearer".
const DefaultTokenType = "Bearer"

// Service defines the interface for JWT token operations
//
//go:generate mockgen -destination=./mocks/service_mock.go -package=jwt_mocks github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/jwt Service
type Service interface {
	GenerateToken(id string, cfg Config) (string, error)
	GetClaims(token string) (Claims, error)
}

// Claims represents a JWT payload.
type Claims struct {
	jwt.RegisteredClaims
	Role string `json:"role"`
}

// Config represents configuration settings for token generation
type Config struct {
	Expiration int // AccessToken expiration duration in seconds
	Role       string
}

type service struct {
	logger log.Logger
	secret []byte // jwt secret key
}

// NewService creates a new JWT service instance
func NewService(logger log.Logger, secret []byte) Service {
	return &service{
		logger: logger,
		secret: secret,
	}
}

func (s *service) GenerateToken(id string, cfg Config) (string, error) {
	now := time.Now()
	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   id,
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(cfg.Expiration) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
		Role: cfg.Role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

func (s *service) GetClaims(token string) (Claims, error) {
	claims, err := s.validateToken(token)
	if err != nil {
		return Claims{}, err
	}

	return claims, nil
}

func (s *service) validateToken(tokenString string) (Claims, error) {
	// Parse and validate the token with explicit validation options
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}))
	if err != nil {
		return Claims{}, ErrInvalidToken
	}

	// Type assertion of the claims
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return Claims{}, ErrInvalidToken
	}

	return *claims, nil
}
