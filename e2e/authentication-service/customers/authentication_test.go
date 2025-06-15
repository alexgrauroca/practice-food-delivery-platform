//go:build e2e || authentication || customers

//nolint:revive
package customers

import (
	"time"

	g "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/alexgrauroca/practice-food-delivery-platform/e2e/authentication-service/customers/customer"
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

		refreshResponse, err := c.Refresh()
		Expect(err).NotTo(HaveOccurred())
		Expect(refreshResponse.AccessToken).NotTo(Equal(c.Auth.AccessToken))
		Expect(refreshResponse.RefreshToken).NotTo(Equal(c.Auth.RefreshToken))
		Expect(refreshResponse.TokenType).To(Equal(c.Auth.TokenType))
		Expect(refreshResponse.ExpiresIn).To(BeNumerically("==", c.Auth.ExpiresIn))
	})
})
