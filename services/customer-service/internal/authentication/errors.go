package authentication

import (
	"errors"
)

// Internal errors
var (
	// ErrInvalidToken represents an error indicating that the authentication token provided is invalid or malformed
	ErrInvalidToken = errors.New("invalid token")
	// ErrTokenExpired represents an error indicating that the authentication token has expired and is no longer valid
	ErrTokenExpired = errors.New("token expired")
)

// HTTP errors
const (
	// CodeUnauthorizedError represents the error code indicating that authentication is required for accessing a
	// resource
	CodeUnauthorizedError = "UNAUTHORIZED"
	// MessageUnauthorizedError represents the error message indicating that authentication is required for
	// accessing a resource
	MessageUnauthorizedError = "Authentication is required to access this resource"

	// CodeTokenExpiredError represents the error code indicating that the authentication token has expired
	CodeTokenExpiredError = "TOKEN_EXPIRED"
	// MessageTokenExpiredError represents the error message indicating that the authentication token has expired
	MessageTokenExpiredError = "Token has expired"

	// CodeForbiddenError represents the error code indicating that the user does not have sufficient permissions
	CodeForbiddenError = "FORBIDDEN"
	// MessageForbiddenError represents the error message indicating that the user does not have sufficient permissions
	MessageForbiddenError = "You do not have permission to access this resource"
)

// ErrorResponse represents a standardized structure for API error responses containing code, message, and optional details.
type ErrorResponse struct {
	Code    string   `json:"code"`
	Message string   `json:"message"`
	Details []string `json:"details"`
}

func newErrorResponse(code, message string) ErrorResponse {
	return ErrorResponse{
		Code:    code,
		Message: message,
		Details: make([]string, 0),
	}
}
