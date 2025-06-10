package customers

import (
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"

	"github.com/alexgrauroca/practice-food-delivery-platform/e2e/authentication-service/customers/customer"
)

var _ = ginkgo.Describe("Customer Registration", func() {
	ginkgo.It("registers a new customer successfully", func() {
		// Register the customer
		c := customer.New()
		result, err := c.Register()

		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		gomega.Expect(result.Email).To(gomega.Equal(c.Email))
		gomega.Expect(result.Name).To(gomega.Equal(c.Name))
		gomega.Expect(result.ID).To(gomega.MatchRegexp(customer.IDRegexPattern))
		gomega.Expect(result.CreatedAt).NotTo(gomega.BeEmpty())
	})
})
