package authentication

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

func (t *Token) GetClaims() (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(t.AccessToken, claims, func(token *jwt.Token) (any, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		//TODO configure secret by env vars
		return []byte("a-string-secret-at-least-256-bits-long"), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}))
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// GetTenant returns the tenant ID from the JWT claims.
func (c *Claims) GetTenant() (string, error) {
	// Returning an error following the same jwt package's convention. For example, GetSubject() (string, error)
	return c.Tenant, nil
}
