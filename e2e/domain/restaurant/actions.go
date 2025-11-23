// Package restaurant provides functionality for managing restaurant operations in the e2e test suite.
package restaurant

import (
	"github.com/alexgrauroca/practice-food-delivery-platform/e2e/domain/staff"
	"github.com/alexgrauroca/practice-food-delivery-platform/e2e/pkg/api"
)

// Register sends a registration request for a TestCustomer and returns a RegisterResponse or an error.
func (c *TestRestaurant) Register(owner *staff.TestStaff) (*RegisterResponse, error) {
	req := RegisterRequest{
		Restaurant: RegisterRestaurantRequest{
			Name:       c.Name,
			LegalName:  c.LegalName,
			TaxID:      c.TaxID,
			TimezoneID: c.TimezoneID,
			Contact: ContactData{
				PhonePrefix: c.Contact.PhonePrefix,
				PhoneNumber: c.Contact.PhoneNumber,
				Email:       c.Contact.Email,
				Address:     c.Contact.Address,
				City:        c.Contact.City,
				PostalCode:  c.Contact.PostalCode,
				CountryCode: c.Contact.CountryCode,
			},
		},
		StaffOwner: staff.RegisterOwnerRequest{
			Email:       owner.Email,
			Password:    owner.Password,
			Name:        owner.Name,
			Address:     owner.Address,
			City:        owner.City,
			PostalCode:  owner.PostalCode,
			CountryCode: owner.CountryCode,
		},
	}
	res, err := api.DoPost[RegisterRequest, RegisterResponse](RegisterEndpoint, req, nil)
	if err == nil {
		if res == nil {
			err = ErrUnexpectedResponse
		} else {
			// Update restaurant data
			c.ID = res.Restaurant.ID

			// Update staff owner data
			owner.ID = res.StaffOwner.ID
			owner.RestaurantID = c.ID
		}
	}

	return res, err
}
