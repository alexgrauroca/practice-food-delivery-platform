package restaurant

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

// New creates and returns a new TestRestaurant with predefined and dynamically generated fields.
func New() TestRestaurant {
	uniqueKey = fmt.Sprintf("%d_%d", time.Now().Nanosecond(), counter.Add(1))

	return TestRestaurant{
		VatCode:    gofakeit.Zip(),
		Name:       generateName(),
		LegalName:  generateLegalName(),
		TaxID:      gofakeit.SSN(),
		TimezoneID: gofakeit.TimeZoneRegion(),
		Contact: TestContact{
			PhonePrefix: "+1",
			PhoneNumber: gofakeit.Phone(),
			Email:       generateEmail(),
			Address:     gofakeit.Street(),
			City:        gofakeit.City(),
			PostalCode:  gofakeit.Zip(),
			CountryCode: gofakeit.CountryAbr(),
		},
	}
}

func generateEmail() string {
	return "e2e_test_restaurant_" + uniqueKey + "@example.com"
}

func generateName() string {
	return "E2E Test Restaurant " + uniqueKey
}

func generateLegalName() string {
	return "Legal Restaurant " + uniqueKey
}
