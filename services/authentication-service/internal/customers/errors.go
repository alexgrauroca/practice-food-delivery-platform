// Package customers provide customer-related functionality and error definitions
// for the authentication service. It defines custom errors for handling common
// customer-related scenarios.
package customers

import "errors"

var (
	// ErrCustomerAlreadyExists indicates that a customer with the same identifying details already exists in the system.
	ErrCustomerAlreadyExists = errors.New("customer already exists")
	// ErrInvalidCredentials indicates that the provided credentials are invalid during authentication processes.
	ErrInvalidCredentials = errors.New("invalid credentials")
	// ErrCustomerNotFound indicates that a customer with the specified details could not be found in the system.
	ErrCustomerNotFound = errors.New("customer not found")
	// ErrInvalidRefreshToken indicates that the provided refresh token is invalid or cannot be used for token renewal.
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	// ErrTokenMismatch indicates a mismatch between the provided access token and the refresh token.
	ErrTokenMismatch = errors.New("token mismatch")
)
