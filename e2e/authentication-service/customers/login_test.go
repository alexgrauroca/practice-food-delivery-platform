package customers

import (
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"

	"github.com/alexgrauroca/practice-food-delivery-platform/e2e/authentication-service/customers/customer"
)

var _ = ginkgo.Describe("Customer Login", func() {
	ginkgo.It("logins a customer successfully", func() {
		// Register the customer
		c := customer.New()
		_, err := c.Register()
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		// Log in the customer
		result, err := c.Login()
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		gomega.Expect(result.AccessToken).NotTo(gomega.BeEmpty())
		gomega.Expect(result.RefreshToken).NotTo(gomega.BeEmpty())
		gomega.Expect(result.ExpiresIn).To(gomega.BeNumerically(">", 0))
		gomega.Expect(result.TokenType).To(gomega.Equal("Bearer"))
	})
})
