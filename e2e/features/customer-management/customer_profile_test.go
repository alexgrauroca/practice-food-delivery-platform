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
		Expect(getCustomerResponse.CreatedAt).NotTo(BeZero())
		Expect(getCustomerResponse.UpdatedAt).NotTo(BeZero())

		tmpCustomer := customer.New()
		updateCustomerResponse, err := c.Update(customer.UpdateCustomerParams{
			Name:        tmpCustomer.Name,
			Address:     tmpCustomer.Address,
			City:        tmpCustomer.City,
			PostalCode:  tmpCustomer.PostalCode,
			CountryCode: tmpCustomer.CountryCode,
		})

		Expect(err).NotTo(HaveOccurred())
		Expect(updateCustomerResponse.ID).To(Equal(c.ID))
		Expect(updateCustomerResponse.Name).To(Equal(tmpCustomer.Name))
		Expect(updateCustomerResponse.Address).To(Equal(tmpCustomer.Address))
		Expect(updateCustomerResponse.City).To(Equal(tmpCustomer.City))
		Expect(updateCustomerResponse.PostalCode).To(Equal(tmpCustomer.PostalCode))
		Expect(updateCustomerResponse.CountryCode).To(Equal(tmpCustomer.CountryCode))
		Expect(updateCustomerResponse.CreatedAt).To(Equal(getCustomerResponse.CreatedAt))
		Expect(updateCustomerResponse.UpdatedAt).NotTo(Equal(getCustomerResponse.UpdatedAt))
	})
})
