// Package restaurant provides functionality for managing restaurant operations in the e2e test suite.
package restaurant

import (
	"github.com/alexgrauroca/practice-food-delivery-platform/e2e/domain/staff"
	"github.com/alexgrauroca/practice-food-delivery-platform/e2e/pkg/api"
)

// Register sends a registration request for a TestCustomer and returns a RegisterResponse or an error.
func (restaurant *TestRestaurant) Register(owner *staff.TestStaff) (*RegisterResponse, error) {
	req := RegisterRequest{
		Restaurant: RegisterRestaurantRequest{
			VatCode:    restaurant.VatCode,
			Name:       restaurant.Name,
			LegalName:  restaurant.LegalName,
			TaxID:      restaurant.TaxID,
			TimezoneID: restaurant.TimezoneID,
			Contact: ContactData{
				PhonePrefix: restaurant.Contact.PhonePrefix,
				PhoneNumber: restaurant.Contact.PhoneNumber,
				Email:       restaurant.Contact.Email,
				Address:     restaurant.Contact.Address,
				City:        restaurant.Contact.City,
				PostalCode:  restaurant.Contact.PostalCode,
				CountryCode: restaurant.Contact.CountryCode,
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
			restaurant.ID = res.Restaurant.ID

			// Update staff owner data
			owner.ID = res.StaffOwner.ID
			owner.RestaurantID = restaurant.ID
			owner.Owner = res.StaffOwner.Owner
		}
	}

	return res, err
}
