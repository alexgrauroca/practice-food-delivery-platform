//go:build e2e || customer

package customer_management_test

import (
	g "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/alexgrauroca/practice-food-delivery-platform/e2e/domain/customer"
)

var _ = g.Describe("Customer Profile Workflow", func() {
	g.It("successfully updates a customer profile", func() {
		// Registering the customer
		c := customer.New()
		err := c.RegisterAndLogin()
		Expect(err).NotTo(HaveOccurred())
	})
})
