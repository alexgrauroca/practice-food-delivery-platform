package jwt

import (
	"crypto/rand"
	"encoding/base64"
)

// DefaultRefreshTokenLength defines the default length, in bytes, of a generated refresh token for authentication purposes.
const DefaultRefreshTokenLength = 32

// GenerateRefreshToken generates a cryptographically secure random refresh token.
func GenerateRefreshToken() (string, error) {
	// Creating a cryptographically secure random refresh token by:
	// 1. Allocating a byte slice of defined length (32 bytes)
	// 2. Filling it with random bytes using crypto/rand
	// 3. Encoding the random bytes to base64 string
	b := make([]byte, DefaultRefreshTokenLength)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
