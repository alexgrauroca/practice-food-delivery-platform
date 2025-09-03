package customer

import "github.com/alexgrauroca/practice-food-delivery-platform/e2e/domain/authentication"

const (
	// IDRegexPattern defines a regular expression pattern to validate a 24-character hexadecimal string typically used as an ID.
	IDRegexPattern = `^[a-fA-F0-9]{24}$`
)

// TestCustomer represents a test customer with email, password, and name information.
type TestCustomer struct {
	ID          string
	Email       string
	Password    string
	Name        string
	Address     string
	City        string
	PostalCode  string
	CountryCode string
	Auth        authentication.Token
}

// RegisterRequest represents the payload required to register a new customer, containing personal details and credentials.
type RegisterRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	Name        string `json:"name"`
	Address     string `json:"address"`
	City        string `json:"city"`
	PostalCode  string `json:"postal_code"`
	CountryCode string `json:"country_code"`
}

// RegisterResponse represents the response data structure for a successful user registration.
type RegisterResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
}

type GetCustomerRequest struct {
	ID string `path:"customerID"`
}

type GetCustomerResponse struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	Name        string `json:"name"`
	Address     string `json:"address"`
	City        string `json:"city"`
	PostalCode  string `json:"postal_code"`
	CountryCode string `json:"country_code"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
