//go:build e2e || authentication || restaurant

//nolint:revive
package authentication_test

import (
	"errors"
	"time"

	g "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/alexgrauroca/practice-food-delivery-platform/e2e/domain/restaurant"
	"github.com/alexgrauroca/practice-food-delivery-platform/e2e/domain/staff"
	"github.com/alexgrauroca/practice-food-delivery-platform/e2e/pkg/api"
)

var _ = g.Describe("Restaurant Authentication Workflow", func() {
	g.It("successfully registers a restaurant, logs in and refresh its staff owner ", func() {
		rest := restaurant.New()
		owner := staff.NewOwner()

		// Register the customer
		registerResponse, err := rest.Register(&owner)
		Expect(err).NotTo(HaveOccurred())

		restaurantResponse := registerResponse.Restaurant
		ownerResponse := registerResponse.StaffOwner

		// Check restaurant response
		Expect(restaurantResponse.ID).To(MatchRegexp(restaurant.IDRegexPattern))
		Expect(restaurantResponse.Name).To(Equal(rest.Name))
		Expect(restaurantResponse.LegalName).To(Equal(rest.LegalName))
		Expect(restaurantResponse.VatCode).To(Equal(rest.VatCode))
		Expect(restaurantResponse.TimezoneID).To(Equal(rest.TimezoneID))
		Expect(restaurantResponse.TaxID).To(Equal(rest.TaxID))
		Expect(restaurantResponse.Contact.PhonePrefix).To(Equal(rest.Contact.PhonePrefix))
		Expect(restaurantResponse.Contact.PhoneNumber).To(Equal(rest.Contact.PhoneNumber))
		Expect(restaurantResponse.Contact.Email).To(Equal(rest.Contact.Email))
		Expect(restaurantResponse.Contact.Address).To(Equal(rest.Contact.Address))
		Expect(restaurantResponse.Contact.City).To(Equal(rest.Contact.City))
		Expect(restaurantResponse.Contact.PostalCode).To(Equal(rest.Contact.PostalCode))
		Expect(restaurantResponse.Contact.CountryCode).To(Equal(rest.Contact.CountryCode))
		Expect(restaurantResponse.CreatedAt).NotTo(BeZero())
		Expect(restaurantResponse.UpdatedAt).NotTo(BeZero())

		// Check staff owner response
		Expect(ownerResponse.ID).To(MatchRegexp(staff.IDRegexPattern))
		Expect(ownerResponse.Email).To(Equal(owner.Email))
		Expect(ownerResponse.Name).To(Equal(owner.Name))
		Expect(ownerResponse.RestaurantID).To(Equal(rest.ID))
		Expect(ownerResponse.Owner).To(BeTrue())
		Expect(ownerResponse.CreatedAt).NotTo(BeZero())
		Expect(ownerResponse.UpdatedAt).NotTo(BeZero())
		Expect(owner.RestaurantID).To(Equal(rest.ID))

		// Log in the registered customer
		loginResponse, err := owner.Login()
		Expect(err).NotTo(HaveOccurred())
		Expect(loginResponse.AccessToken).NotTo(BeEmpty())
		Expect(loginResponse.RefreshToken).NotTo(BeEmpty())
		Expect(loginResponse.ExpiresIn).To(BeNumerically(">", 0))
		Expect(loginResponse.TokenType).To(Equal("Bearer"))
		claims, err := loginResponse.Token.GetClaims()
		Expect(err).NotTo(HaveOccurred())
		Expect(claims.GetSubject()).To(Equal(owner.ID))
		Expect(claims.GetTenant()).To(Equal(rest.ID))

		// Refresh the owner auth tokens
		// Wait 1 second to ensure a difference between ExpiresIn values from each new access_token
		time.Sleep(time.Second * 1)

		refreshResponse1, err := owner.Refresh()
		Expect(err).NotTo(HaveOccurred())
		Expect(refreshResponse1.AccessToken).NotTo(Equal(owner.Auth.AccessToken))
		Expect(refreshResponse1.RefreshToken).NotTo(Equal(owner.Auth.RefreshToken))
		Expect(refreshResponse1.TokenType).To(Equal(owner.Auth.TokenType))
		Expect(refreshResponse1.ExpiresIn).To(BeNumerically("==", owner.Auth.ExpiresIn))

		// If I wait 1 second and try to refresh again with the same refresh token, the auth should be generated anyway
		time.Sleep(time.Second * 1)

		refreshResponse2, err := owner.Refresh()
		Expect(err).NotTo(HaveOccurred())
		Expect(refreshResponse2.AccessToken).NotTo(Equal(refreshResponse1.AccessToken))
		Expect(refreshResponse2.RefreshToken).NotTo(Equal(refreshResponse1.RefreshToken))
		Expect(refreshResponse2.TokenType).To(Equal(owner.Auth.TokenType))
		Expect(refreshResponse2.ExpiresIn).To(BeNumerically("==", owner.Auth.ExpiresIn))

		// If I wait 5 seconds and try to refresh again with the same refresh token, the action should fail
		time.Sleep(time.Second * 5)

		_, err = owner.Refresh()
		Expect(err).To(HaveOccurred())

		var apiErr *api.ErrorResponse
		ok := errors.As(err, &apiErr)
		Expect(ok).To(BeTrue(), "Expected ErrorResponse type")
		Expect(apiErr.Code).To(Equal("INVALID_REFRESH_TOKEN"))
	})
})
