// Package staff provides functionality for managing staff operations in the e2e test suite.
package staff

import (
	"github.com/alexgrauroca/practice-food-delivery-platform/e2e/domain/authentication"
	"github.com/alexgrauroca/practice-food-delivery-platform/e2e/pkg/api"
)

// Login authenticates a staff by sending their email and password to the login API and returns the login response or an
// error.
func (staff *TestStaff) Login() (*authentication.LoginResponse, error) {
	req := authentication.StaffLoginRequest{
		LoginRequest: authentication.LoginRequest{
			Email:    staff.Email,
			Password: staff.Password,
		},
		RestaurantID: staff.RestaurantID,
	}
	res, err := api.DoPost[authentication.StaffLoginRequest, authentication.LoginResponse](LoginEndpoint, req, nil)
	if err == nil {
		if res == nil {
			err = ErrUnexpectedResponse
		} else {
			staff.SetAuth(res.Token)
		}
	}

	return res, err
}

// Refresh attempts to refresh the authentication token for the TestStaff using the provided access and refresh tokens.
func (staff *TestStaff) Refresh() (*authentication.RefreshResponse, error) {
	req := authentication.RefreshRequest{
		AccessToken:  staff.Auth.AccessToken,
		RefreshToken: staff.Auth.RefreshToken,
	}

	return api.DoPost[authentication.RefreshRequest, authentication.RefreshResponse](RefreshEndpoint, req, nil)
}

// SetAuth sets the authentication token for the TestStaff instance.
func (staff *TestStaff) SetAuth(auth authentication.Token) {
	staff.Auth = auth
}
