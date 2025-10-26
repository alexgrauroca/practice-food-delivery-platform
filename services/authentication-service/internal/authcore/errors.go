package authcore

import "errors"

var (
	// ErrInvalidCredentials indicates that the provided credentials are invalid during authentication processes.
	ErrInvalidCredentials = errors.New("invalid credentials")
	// ErrInvalidRefreshToken indicates that the provided refresh token is invalid or cannot be used for token renewal.
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	// ErrTokenMismatch indicates a mismatch between the provided access token and the refresh token.
	ErrTokenMismatch = errors.New("token mismatch")
)

const (
	// CodeInvalidCredentials represents the error code for failed authentication due to invalid login credentials.
	CodeInvalidCredentials = "INVALID_CREDENTIALS"
	// MsgInvalidCredentials represents the error message returned when login authentication fails due to invalid credentials.
	MsgInvalidCredentials = "invalid credentials"
)
