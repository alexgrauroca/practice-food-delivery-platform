// Package customer provides test utilities and data structures for customer-related end-to-end tests
package customer

import (
	"fmt"
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
	return fmt.Sprintf("e2e_test_user_%d@example.com", time.Now().UnixNano())
}

func generateName() string {
	return fmt.Sprintf("E2E Test User%d", time.Now().UnixNano())
}
