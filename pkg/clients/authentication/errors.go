package authentication

import (
	"errors"
)

// Internal errors
var (
	// ErrInvalidToken represents an error when the token's signature, claims, or format is invalid
	ErrInvalidToken = errors.New("invalid token")
	// ErrTokenExpired represents an error when the token has exceeded its expiration time
	ErrTokenExpired = errors.New("token expired")
	// ErrAuthHeaderMissing represents an error when the Authorization header is not present in the request
	ErrAuthHeaderMissing = errors.New("authorization header is missing")
	// ErrInvalidAuthHeader represents an error when the Authorization header format does not match "Bearer <token>"
	ErrInvalidAuthHeader = errors.New("invalid authorization header format")
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
