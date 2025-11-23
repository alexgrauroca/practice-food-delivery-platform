package restaurant

import (
	"time"

	"github.com/alexgrauroca/practice-food-delivery-platform/e2e/domain/staff"
)

const (
	// IDRegexPattern defines a regular expression pattern to validate a 24-character hexadecimal string typically used as an ID.
	IDRegexPattern = `^[a-fA-F0-9]{24}$`
)

// TestRestaurant represents a test restaurant.
type TestRestaurant struct {
	ID         string
	VatCode    string
	Name       string
	LegalName  string
	TaxID      string
	TimezoneID string
	Contact    TestContact
}

// TestContact represents the contact details of a restaurant.
type TestContact struct {
	PhonePrefix string
	PhoneNumber string
	Email       string
	Address     string
	City        string
	PostalCode  string
	CountryCode string
}

// APIData captures the common restaurant fields returned by multiple endpoints.
type APIData struct {
	ID         string      `json:"id"`
	VatCode    string      `json:"vat_code"`
	Name       string      `json:"name"`
	LegalName  string      `json:"legal_name"`
	TaxID      string      `json:"tax_id"`
	TimezoneID string      `json:"timezone_id"`
	Contact    ContactData `json:"contact"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
}

// ContactData represents the contact details of a restaurant.
type ContactData struct {
	PhonePrefix string `json:"phone_prefix"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	Address     string `json:"address"`
	City        string `json:"city"`
	PostalCode  string `json:"postal_code"`
	CountryCode string `json:"country_code"`
}

// RegisterRequest represents the payload required to register a new restaurant.
type RegisterRequest struct {
	Restaurant RegisterRestaurantRequest  `json:"restaurant"`
	StaffOwner staff.RegisterOwnerRequest `json:"staff_owner"`
}

type RegisterRestaurantRequest struct {
	VatCode    string      `json:"vat_code"`
	Name       string      `json:"name"`
	LegalName  string      `json:"legal_name"`
	TaxID      string      `json:"tax_id"`
	TimezoneID string      `json:"timezone_id"`
	Contact    ContactData `json:"contact"`
}

// RegisterResponse represents the response data structure for a successful restaurant registration.
type RegisterResponse struct {
	Restaurant APIData       `json:"restaurant"`
	StaffOwner staff.APIData `json:"staff_owner"`
}

// GetCustomerRequest represents the request to retrieve customer details, containing the customer ID as a path parameter.
type GetCustomerRequest struct {
	ID string `path:"customerID"`
}
