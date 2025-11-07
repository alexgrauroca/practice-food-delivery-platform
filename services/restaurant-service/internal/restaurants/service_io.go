package restaurants

import "time"

// RegisterRestaurantInput represents the input payload for registering a new restaurant.
type RegisterRestaurantInput struct {
	Restaurant RestaurantInput
	StaffOwner StaffOwnerInput
}

// RestaurantInput represents the input payload for retrieving a restaurant.
type RestaurantInput struct {
	VatCode    string
	Name       string
	LegalName  string
	TaxID      string
	TimezoneID string
	Contact    ContactInput
}

// ContactInput represents the input payload for retrieving a restaurant's contact information.
type ContactInput struct {
	PhonePrefix string
	PhoneNumber string
	Email       string
	Address     string
	City        string
	PostalCode  string
	CountryCode string
}

// StaffOwnerInput represents the input payload for retrieving a restaurant's staff owner.
type StaffOwnerInput struct {
	Email       string
	Password    string
	Name        string
	Address     string
	City        string
	PostalCode  string
	CountryCode string
}

// RegisterRestaurantOutput represents the output payload for registering a new restaurant.
type RegisterRestaurantOutput struct {
	Restaurant RestaurantOutput
	StaffOwner StaffOwnerOutput
}

// RestaurantOutput represents the output payload for retrieving a restaurant.
type RestaurantOutput struct {
	ID         string
	VatCode    string
	Name       string
	LegalName  string
	TaxID      string
	TimezoneID string
	Contact    ContactOutput
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// ContactOutput represents the output payload for retrieving a restaurant's contact information.
type ContactOutput struct {
	PhonePrefix string
	PhoneNumber string
	Email       string
	Address     string
	City        string
	PostalCode  string
	CountryCode string
}

// StaffOwnerOutput represents the output payload for retrieving a restaurant's staff owner.
type StaffOwnerOutput struct {
	ID          string
	Owner       bool
	Email       string
	Name        string
	Address     string
	City        string
	PostalCode  string
	CountryCode string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
