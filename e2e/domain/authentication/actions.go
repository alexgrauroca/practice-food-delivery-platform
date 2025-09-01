package authentication

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

// GetSubject extracts the subject claim from the access token.
func (t *Token) GetSubject() (string, error) {
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
		return "", fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return "", errors.New("invalid token")
	}

	return claims.Subject, nil

}
