// Package authentication provides functionality for customer authentication and registration
// through integration with the authentication service. It handles customer registration,
// login, and token management operations.
package authentication

import (
	"context"
)

// Client defines the interface for interacting with the authentication service.
// It provides methods for customer registration and authentication operations.
//
//go:generate mockgen -destination=./mocks/authclient_mock.go -package=authentication_mocks github.com/alexgrauroca/practice-food-delivery-platform/services/customer-service/internal/authentication Client
type Client interface {
	RegisterCustomer(ctx context.Context, req RegisterCustomerRequest) (RegisterCustomerResponse, error)
}

// Config holds the configuration options for the authentication client.
type Config struct {
	Debug bool
}

// RegisterCustomerRequest represents the data required to register a new customer
// in the authentication service.
type RegisterCustomerRequest struct {
	Email    string
	Password string
	Name     string
}

// RegisterCustomerResponse contains the data returned after successfully
// registering a customer in the authentication service.
type RegisterCustomerResponse struct {
	ID    string
	Email string
}
