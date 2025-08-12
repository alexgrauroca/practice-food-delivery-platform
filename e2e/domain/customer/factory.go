// Package customer provides test utilities and data structures for customer-related end-to-end tests
package customer

import (
	"time"
)

// New creates and returns a new TestCustomer with predefined and dynamically generated fields.
func New() TestCustomer {
	return TestCustomer{
		Email:       generateEmail(),
		Password:    "strongpassword123",
		Name:        generateName(),
		Address:     "123 Main St.",
		City:        "Anytown",
		PostalCode:  "12345",
		CountryCode: "US",
	}
}

func generateEmail() string {
	return "e2e_test_user_" + time.Now().Format("150405") + "@example.com"
}

func generateName() string {
	return "E2E Test User" + time.Now().Format("150405")
}
