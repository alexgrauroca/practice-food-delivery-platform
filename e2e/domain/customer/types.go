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

// customerData captures the common customer fields returned by multiple endpoints.
type customerData struct {
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

// GetCustomerRequest represents the request to retrieve customer details, containing the customer ID as a path parameter.
type GetCustomerRequest struct {
	ID string `path:"customerID"`
}

// GetCustomerResponse represents the response data structure containing a customer's full profile information.
type GetCustomerResponse = customerData

// UpdateCustomerParams holds the parameters for updating a customer's profile information.
type UpdateCustomerParams struct {
	Name        string
	Address     string
	City        string
	PostalCode  string
	CountryCode string
}

// UpdateCustomerRequest represents the request payload for updating customer profile information.
type UpdateCustomerRequest struct {
	ID          string `path:"customerID"`
	Name        string `json:"name"`
	Address     string `json:"address"`
	City        string `json:"city"`
	PostalCode  string `json:"postal_code"`
	CountryCode string `json:"country_code"`
}

// UpdateCustomerResponse represents the response data structure containing a customer's full profile information.
type UpdateCustomerResponse = customerData
