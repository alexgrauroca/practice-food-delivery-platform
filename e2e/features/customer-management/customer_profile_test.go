//go:build e2e || customer

//nolint:revive
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

		getCustomerResponse, err := c.Get()
		Expect(err).NotTo(HaveOccurred())
		Expect(getCustomerResponse.ID).To(Equal(c.ID))
		Expect(getCustomerResponse.Name).To(Equal(c.Name))
		Expect(getCustomerResponse.Email).To(Equal(c.Email))
		Expect(getCustomerResponse.Address).To(Equal(c.Address))
		Expect(getCustomerResponse.City).To(Equal(c.City))
		Expect(getCustomerResponse.PostalCode).To(Equal(c.PostalCode))
		Expect(getCustomerResponse.CountryCode).To(Equal(c.CountryCode))
		Expect(getCustomerResponse.CreatedAt).NotTo(BeEmpty())
		Expect(getCustomerResponse.UpdatedAt).NotTo(BeEmpty())
	})
})
