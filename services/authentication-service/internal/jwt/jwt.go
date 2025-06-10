// Package jwt provides functionality for generating JWT tokens with custom claims and expiration.
package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// DefaultTokenType represents the default type of token used for authorization, typically set to "Bearer".
const DefaultTokenType = "Bearer"

var jwtSecret = []byte("your-very-secret-key")

// Config represents configuration settings for token generation.
// Expiration defines the token's lifetime in seconds.
// Role specifies the role assigned to the generated token.
type Config struct {
	Expiration int // AccessToken expiration duration in seconds
	Role       string
}

// GenerateToken generates a JWT token with claims based on the provided ID and configuration settings.
// It includes the ID as the subject, assigns the specified role, and sets the expiration time.
// Returns the signed JWT token string or an error in case of token creation failure.
func GenerateToken(id string, cfg Config) (string, error) {
	claims := jwt.MapClaims{
		"sub":  id,
		"role": cfg.Role,
		"exp":  time.Now().Add(time.Duration(cfg.Expiration) * time.Second).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
