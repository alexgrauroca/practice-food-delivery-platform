//go:build e2e || authentication || customer

//nolint:revive
package authentication_test

import (
	"errors"
	"time"

	g "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/alexgrauroca/practice-food-delivery-platform/e2e/domain/customer"
	"github.com/alexgrauroca/practice-food-delivery-platform/e2e/pkg/api"
)

var _ = g.Describe("Customer Authentication Workflow", func() {
	g.It("successfully registers, logs in and refresh a customer", func() {
		c := customer.New()

		// Register the customer
		registerResponse, err := c.Register()

		Expect(err).NotTo(HaveOccurred())
		Expect(registerResponse.Email).To(Equal(c.Email))
		Expect(registerResponse.Name).To(Equal(c.Name))
		Expect(registerResponse.ID).To(MatchRegexp(customer.IDRegexPattern))
		Expect(registerResponse.CreatedAt).NotTo(BeEmpty())

		// Log in the registered customer
		loginResponse, err := c.Login()
		Expect(err).NotTo(HaveOccurred())
		Expect(loginResponse.AccessToken).NotTo(BeEmpty())
		Expect(loginResponse.RefreshToken).NotTo(BeEmpty())
		Expect(loginResponse.ExpiresIn).To(BeNumerically(">", 0))
		Expect(loginResponse.TokenType).To(Equal("Bearer"))

		c.SetAuth(loginResponse.Token)

		// Refresh the customer auth tokens
		// Wait 1 second to ensure a difference between ExpiresIn values from each new access_token
		time.Sleep(time.Second * 1)

		refreshResponse1, err := c.Refresh()
		Expect(err).NotTo(HaveOccurred())
		Expect(refreshResponse1.AccessToken).NotTo(Equal(c.Auth.AccessToken))
		Expect(refreshResponse1.RefreshToken).NotTo(Equal(c.Auth.RefreshToken))
		Expect(refreshResponse1.TokenType).To(Equal(c.Auth.TokenType))
		Expect(refreshResponse1.ExpiresIn).To(BeNumerically("==", c.Auth.ExpiresIn))

		// If I wait 1 second and try to refresh again with the same refresh token, the auth should be generated anyway
		time.Sleep(time.Second * 1)

		refreshResponse2, err := c.Refresh()
		Expect(err).NotTo(HaveOccurred())
		Expect(refreshResponse2.AccessToken).NotTo(Equal(refreshResponse1.AccessToken))
		Expect(refreshResponse2.RefreshToken).NotTo(Equal(refreshResponse1.RefreshToken))
		Expect(refreshResponse2.TokenType).To(Equal(c.Auth.TokenType))
		Expect(refreshResponse2.ExpiresIn).To(BeNumerically("==", c.Auth.ExpiresIn))

		// If I wait 5 seconds and try to refresh again with the same refresh token, the action should fail
		time.Sleep(time.Second * 5)

		_, err = c.Refresh()
		Expect(err).To(HaveOccurred())

		var apiErr *api.ErrorResponse
		ok := errors.As(err, &apiErr)
		Expect(ok).To(BeTrue(), "Expected ErrorResponse type")
		Expect(apiErr.Code).To(Equal("INVALID_REFRESH_TOKEN"))
	})
})
