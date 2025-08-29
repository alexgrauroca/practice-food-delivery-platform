// Package customers provide customer-related functionality and error definitions
// for the customer service. It defines custom errors for handling common
// customer-related scenarios.
package customers

import "errors"

var (
	// ErrCustomerAlreadyExists indicates that a customer with the same identifying details already exists in the system.
	ErrCustomerAlreadyExists = errors.New("customer already exists")
	// ErrCustomerNotFound indicates that a customer with the specified details could not be found in the system.
	ErrCustomerNotFound = errors.New("customer not found")
	// ErrCustomerIDMismatch indicates that the requested customer ID does not match the authenticated customer's identity,
	// which is typically used when validating access permissions for customer-specific operations.
	ErrCustomerIDMismatch = errors.New("customer ID does not match authenticated customer")
)
