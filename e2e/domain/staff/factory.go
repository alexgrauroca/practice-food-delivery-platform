package staff

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

// NewOwner creates and returns a new TestStaff owner with predefined and dynamically generated fields.
func NewOwner() TestStaff {
	uniqueKey = fmt.Sprintf("%d_%d", time.Now().Nanosecond(), counter.Add(1))

	return TestStaff{
		Email:       generateEmail(),
		Owner:       true,
		Password:    "strongpassword123",
		Name:        generateName(),
		Address:     gofakeit.Street(),
		City:        gofakeit.City(),
		PostalCode:  gofakeit.Zip(),
		CountryCode: gofakeit.CountryAbr(),
	}
}

func generateEmail() string {
	return "e2e_test_staff_" + uniqueKey + "@example.com"
}

func generateName() string {
	return "E2E Test Staff " + uniqueKey
}
