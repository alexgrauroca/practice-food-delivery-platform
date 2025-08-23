package authentication

import "errors"

var (
	// ErrInvalidToken represents an error indicating that the authentication token provided is invalid or malformed
	ErrInvalidToken = errors.New("invalid token")
	// ErrExpiredToken represents an error indicating that the authentication token has expired and is no longer valid
	ErrExpiredToken = errors.New("expired token")
)