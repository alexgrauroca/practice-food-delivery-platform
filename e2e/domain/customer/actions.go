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
	res, err := api.DoPost[RegisterRequest, RegisterResponse](RegisterEndpoint, req)
	if err == nil {
		if res == nil {
			err = ErrUnexpectedResponse
		} else {
			c.ID = res.ID
		}
	}

	return res, err
}

// Login authenticates a user by sending their email and password to the login API and returns the login response or an error.
func (c *TestCustomer) Login() (*authentication.LoginResponse, error) {
	req := authentication.LoginRequest{
		Email:    c.Email,
		Password: c.Password,
	}
	res, err := api.DoPost[authentication.LoginRequest, authentication.LoginResponse](LoginEndpoint, req)
	if err == nil {
		if res == nil {
			err = ErrUnexpectedResponse
		} else {
			c.SetAuth(res.Token)
		}
	}

	return res, err
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

// RegisterAndLogin performs both registration and login operations for a TestCustomer sequentially.
// Returns an error if either registration or login fails.
func (c *TestCustomer) RegisterAndLogin() error {
	if _, err := c.Register(); err != nil {
		return err
	}
	if _, err := c.Login(); err != nil {
		return err
	}
	return nil
}
