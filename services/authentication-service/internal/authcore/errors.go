package authcore

import "errors"

var (
	// ErrInvalidCredentials indicates that the provided credentials are invalid during authentication processes.
	ErrInvalidCredentials = errors.New("invalid credentials")
)

const (
	// CodeInvalidCredentials represents the error code for failed authentication due to invalid login credentials.
	CodeInvalidCredentials = "INVALID_CREDENTIALS"
	// MsgInvalidCredentials represents the error message returned when login authentication fails due to invalid credentials.
	MsgInvalidCredentials = "invalid credentials"
)
