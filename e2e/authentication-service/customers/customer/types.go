package customer

// TestCustomer represents a test customer with email, password, and name information.
type TestCustomer struct {
	Email    string
	Password string
	Name     string
}

// RegisterResponse represents the response data structure for a successful user registration.
type RegisterResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
}

// LoginResponse represents the response payload returned upon a successful login containing token information.
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

const (
	// IDRegexPattern defines a regular expression pattern to validate a 24-character hexadecimal string typically used as an ID.
	IDRegexPattern = `^[a-fA-F0-9]{24}$`
)
