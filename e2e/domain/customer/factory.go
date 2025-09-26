// Package customer provides test utilities and data structures for customer-related end-to-end tests
package customer

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/brianvoe/gofakeit/v7"
)

var counter atomic.Int64
var uniqueKey string

func init() {
	_ = gofakeit.Seed(time.Now().UnixNano())
}

// New creates and returns a new TestCustomer with predefined and dynamically generated fields.
func New() TestCustomer {
	uniqueKey = fmt.Sprintf("%d_%d", time.Now().Nanosecond(), counter.Add(1))

	return TestCustomer{
		Email:       generateEmail(),
		Password:    "strongpassword123",
		Name:        generateName(),
		Address:     gofakeit.Street(),
		City:        gofakeit.City(),
		PostalCode:  gofakeit.Zip(),
		CountryCode: gofakeit.CountryAbr(),
	}
}

func generateEmail() string {
	return "e2e_test_user_" + uniqueKey + "@example.com"
}

func generateName() string {
	return "E2E Test User " + uniqueKey
}
