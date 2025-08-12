package customer

import (
	"github.com/alexgrauroca/practice-food-delivery-platform/e2e/domain/authentication"
	"github.com/alexgrauroca/practice-food-delivery-platform/e2e/pkg/api"
)

// Register sends a registration request for a TestCustomer and returns a RegisterResponse or an error.
func (c *TestCustomer) Register() (*RegisterResponse, error) {
	req := RegisterRequest{
		Email:       c.Email,
		Password:    c.Password,
		Name:        c.Name,
		Address:     c.Address,
		City:        c.City,
		PostalCode:  c.PostalCode,
		CountryCode: c.CountryCode,
	}

	return api.DoPost[RegisterRequest, RegisterResponse](RegisterEndpoint, req)
}

// Login authenticates a user by sending their email and password to the login API and returns the login response or an error.
func (c *TestCustomer) Login() (*authentication.LoginResponse, error) {
	req := authentication.LoginRequest{
		Email:    c.Email,
		Password: c.Password,
	}

	return api.DoPost[authentication.LoginRequest, authentication.LoginResponse](LoginEndpoint, req)
}

// Refresh attempts to refresh the authentication token for the TestCustomer using the provided access and refresh tokens.
func (c *TestCustomer) Refresh() (*authentication.RefreshResponse, error) {
	req := authentication.RefreshRequest{
		AccessToken:  c.Auth.AccessToken,
		RefreshToken: c.Auth.RefreshToken,
	}

	return api.DoPost[authentication.RefreshRequest, authentication.RefreshResponse](RefreshEndpoint, req)
}

// SetAuth sets the authentication token for the TestCustomer instance.
func (c *TestCustomer) SetAuth(auth authentication.Token) {
	c.Auth = auth
}
