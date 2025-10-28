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

	// CodeInvalidRefreshToken represents the error code for an invalid or expired refresh token used in authentication processes.
	CodeInvalidRefreshToken = "INVALID_REFRESH_TOKEN"
	// MsgInvalidRefreshToken represents an error message indicating an invalid or expired refresh token.
	MsgInvalidRefreshToken = "invalid or expired refresh token"

	// CodeTokenMismatch represents an error code indicating a mismatch between the provided token and the expected value.
	CodeTokenMismatch = "TOKEN_MISMATCH"
	// MsgTokenMismatch represents the error message for a token mismatch scenario.
	MsgTokenMismatch = "token mismatch"
)
