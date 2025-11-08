package restaurants

import "context"

// Restaurant represents a restaurant.
type Restaurant struct {
	ID         string  `bson:"_id,omitempty"`
	VatCode    string  `bson:"vat_code"`
	Name       string  `bson:"name"`
	LegalName  string  `bson:"legal_name"`
	TaxID      string  `bson:"tax_id"`
	TimezoneID string  `bson:"timezone_id"`
	Contact    Contact `bson:"contact"`
}

// Contact represents the contact details of a restaurant.
type Contact struct {
	PhonePrefix string `bson:"phone_prefix"`
	PhoneNumber string `bson:"phone_number"`
	Email       string `bson:"email"`
	Address     string `bson:"address"`
	City        string `bson:"city"`
	PostalCode  string `bson:"postal_code"`
	CountryCode string `bson:"country_code"`
}

// Repository represents the interface for operations related to restaurant management.
//
//go:generate mockgen -destination=./mocks/repository_mock.go -package=restaurants_mocks github.com/alexgrauroca/practice-food-delivery-platform/services/restaurant-service/internal/restaurants Repository
type Repository interface {
	CreateRestaurant(ctx context.Context, params CreateRestaurantParams) (Restaurant, error)
	PurgeRestaurant(ctx context.Context, vatCode string) error
}

// CreateRestaurantParams represents the parameters for creating a restaurant.
type CreateRestaurantParams struct {
	VatCode    string
	Name       string
	LegalName  string
	TaxID      string
	TimezoneID string
	Contact    CreateContactParams
}

// CreateContactParams represents the parameters for creating a restaurant's contact information.'
type CreateContactParams struct {
	PhonePrefix string
	PhoneNumber string
	Email       string
	Address     string
	City        string
	PostalCode  string
	CountryCode string
}
