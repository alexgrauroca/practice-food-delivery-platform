package staff

import (
	"time"

	"github.com/alexgrauroca/practice-food-delivery-platform/e2e/domain/authentication"
)

const (
	// IDRegexPattern defines a regular expression pattern to validate a 24-character hexadecimal string typically used as an ID.
	IDRegexPattern = `^[a-fA-F0-9]{24}$`
)

// TestStaff represents test staff information.
type TestStaff struct {
	ID           string
	Email        string
	RestaurantID string
	Owner        bool
	Password     string
	Name         string
	Address      string
	City         string
	PostalCode   string
	CountryCode  string
	Auth         authentication.Token
}

// APIData captures the common staff fields returned by multiple endpoints.
type APIData struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	RestaurantID string    `json:"restaurant_id"`
	Owner        bool      `json:"owner"`
	Name         string    `json:"name"`
	Address      string    `json:"address"`
	City         string    `json:"city"`
	PostalCode   string    `json:"postal_code"`
	CountryCode  string    `json:"country_code"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// RegisterOwnerRequest represents the payload required to register a new staff owner.
type RegisterOwnerRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	Name        string `json:"name"`
	Address     string `json:"address"`
	City        string `json:"city"`
	PostalCode  string `json:"postal_code"`
	CountryCode string `json:"country_code"`
}
