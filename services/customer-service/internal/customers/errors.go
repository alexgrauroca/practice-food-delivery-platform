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
)
