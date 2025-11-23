// Package staff provides functionality for managing staff operations in the e2e test suite.
package staff

import (
	"github.com/alexgrauroca/practice-food-delivery-platform/e2e/domain/authentication"
	"github.com/alexgrauroca/practice-food-delivery-platform/e2e/pkg/api"
)

// Login authenticates a staff by sending their email and password to the login API and returns the login response or an
// error.
func (c *TestStaff) Login() (*authentication.LoginResponse, error) {
	req := authentication.StaffLoginRequest{
		LoginRequest: authentication.LoginRequest{
			Email:    c.Email,
			Password: c.Password,
		},
		RestaurantID: c.RestaurantID,
	}
	res, err := api.DoPost[authentication.StaffLoginRequest, authentication.LoginResponse](LoginEndpoint, req, nil)
	if err == nil {
		if res == nil {
			err = ErrUnexpectedResponse
		} else {
			c.SetAuth(res.Token)
		}
	}

	return res, err
}

// Refresh attempts to refresh the authentication token for the TestStaff using the provided access and refresh tokens.
func (c *TestStaff) Refresh() (*authentication.RefreshResponse, error) {
	req := authentication.RefreshRequest{
		AccessToken:  c.Auth.AccessToken,
		RefreshToken: c.Auth.RefreshToken,
	}

	return api.DoPost[authentication.RefreshRequest, authentication.RefreshResponse](RefreshEndpoint, req, nil)
}

// SetAuth sets the authentication token for the TestStaff instance.
func (c *TestStaff) SetAuth(auth authentication.Token) {
	c.Auth = auth
}
