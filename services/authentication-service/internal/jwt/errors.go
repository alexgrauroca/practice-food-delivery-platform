package jwt

import "errors"

// ErrInvalidToken indicates that a provided token is invalid or cannot be authenticated.
var ErrInvalidToken = errors.New("invalid token")
