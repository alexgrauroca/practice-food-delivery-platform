package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	DefaultTokenType = "Bearer"
)

var jwtSecret = []byte("your-very-secret-key")

type Config struct {
	Expiration int // Token expiration duration in seconds
	Role       string
}

func GenerateToken(id string, cfg Config) (string, error) {
	claims := jwt.MapClaims{
		"sub":  id,
		"role": cfg.Role,
		"exp":  time.Now().Add(time.Duration(cfg.Expiration) * time.Second).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
