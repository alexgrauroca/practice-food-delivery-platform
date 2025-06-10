// Package customer provides test utilities and data structures for customer-related end-to-end tests
package customer

import (
	"time"

	"github.com/alexgrauroca/practice-food-delivery-platform/e2e/authentication-service/customers/api"
)

// New creates and returns a new TestCustomer with predefined and dynamically generated fields.
func New() TestCustomer {
	return TestCustomer{
		Email:    generateEmail(),
		Password: "strongpassword123",
		Name:     generateName(),
	}
}

// Login authenticates a user by sending their email and password to the login API and returns the login response or an error.
func (c TestCustomer) Login() (*LoginResponse, error) {
	payload := map[string]any{
		"email":    c.Email,
		"password": c.Password,
	}

	return api.DoPost[LoginResponse](api.LoginEndpoint, payload)
}

// Register sends a registration request for a TestCustomer and returns a RegisterResponse or an error.
func (c TestCustomer) Register() (*RegisterResponse, error) {
	payload := map[string]any{
		"email":    c.Email,
		"password": c.Password,
		"name":     c.Name,
	}

	return api.DoPost[RegisterResponse](api.RegisterEndpoint, payload)
}

func generateEmail() string {
	return "e2e_test_user_" + time.Now().Format("150405") + "@example.com"
}

func generateName() string {
	return "E2E Test User" + time.Now().Format("150405")
}
