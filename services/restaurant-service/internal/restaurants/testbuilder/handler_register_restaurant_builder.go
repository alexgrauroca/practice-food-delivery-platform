// Package testbuilder provides a set of test builders for the restaurant service.
package testbuilder

import "strings"

// HandlerRegisterRestaurantBuilder represents a builder for the register restaurant handler.
type HandlerRegisterRestaurantBuilder struct {
	jsonPayload string
}

// NewValidRegisterRestaurantPayload creates a new instance with a valid REQUEST payload.
func NewValidRegisterRestaurantPayload() HandlerRegisterRestaurantBuilder {
	return HandlerRegisterRestaurantBuilder{
		jsonPayload: `{
			"restaurant": {
				"vat_code": "GB123456789",
				"name": "Acme Pizza",
				"legal_name": "Acme Pizza LLC",
				"tax_id": "99-1234567",
				"timezone_id": "America/New_York",
				"contact": {
					"phone_prefix": "+1",
					"phone_number": "1234567890",
					"email": "restaurant@example.com",
					"address": "123 Main St",
					"city": "New York",
					"postal_code": "10001",
					"country_code": "US"
				}
			},
			"staff_owner": {
				"email": "user@example.com",
				"password": "strongpassword123",
				"name": "John Doe",
				"address": "123 Main St",
				"city": "New York",
				"postal_code": "10001",
				"country_code": "US"
			}
		}`,
	}
}

// NewRegisterRestaurantSuccessResponse creates a builder for the SUCCESS RESPONSE JSON used in tests.
func NewRegisterRestaurantSuccessResponse() HandlerRegisterRestaurantBuilder {
	return HandlerRegisterRestaurantBuilder{
		jsonPayload: `{
			"restaurant": {
				"id": "fake-restaurant-id",
				"vat_code": "GB123456789",
				"name": "Acme Pizza",
				"legal_name": "Acme Pizza LLC",
				"tax_id": "99-1234567",
				"timezone_id": "America/New_York",
				"contact": {
					"phone_prefix": "+1",
					"phone_number": "1234567890",
					"email": "restaurant@example.com",
					"address": "123 Main St",
					"city": "New York",
					"postal_code": "10001",
					"country_code": "US"
				},
				"created_at": "2025-01-01T00:00:00Z",
				"updated_at": "2025-01-01T00:00:00Z"
			},
			"staff_owner": {
				"id": "fake-owner-id",
				"owner": true,
				"email": "user@example.com",
				"name": "John Doe",
				"address": "123 Main St",
				"city": "New York",
				"postal_code": "10001",
				"country_code": "US",
				"created_at": "2025-01-01T00:00:00Z",
				"updated_at": "2025-01-01T00:00:00Z"
			}
		}`,
	}
}

// Build returns the current JSON string.
func (b HandlerRegisterRestaurantBuilder) Build() string {
	return b.jsonPayload
}

// WithContactEmail sets the restaurant contact email.
func (b HandlerRegisterRestaurantBuilder) WithContactEmail(email string) HandlerRegisterRestaurantBuilder {
	nb := b
	nb.jsonPayload = strings.ReplaceAll(
		nb.jsonPayload,
		`"email": "restaurant@example.com"`,
		`"email": "`+email+`"`,
	)
	return nb
}

// WithOwnerEmail sets the staff owner email.
func (b HandlerRegisterRestaurantBuilder) WithOwnerEmail(email string) HandlerRegisterRestaurantBuilder {
	nb := b
	nb.jsonPayload = strings.ReplaceAll(
		nb.jsonPayload,
		`"email": "user@example.com"`,
		`"email": "`+email+`"`,
	)
	return nb
}

// WithRestaurantTimezone sets the restaurant timezone_id.
func (b HandlerRegisterRestaurantBuilder) WithRestaurantTimezone(tz string) HandlerRegisterRestaurantBuilder {
	nb := b
	nb.jsonPayload = strings.ReplaceAll(
		nb.jsonPayload,
		`"timezone_id": "America/New_York"`,
		`"timezone_id": "`+tz+`"`,
	)
	return nb
}

// WithContactPhone sets both phone prefix and number.
func (b HandlerRegisterRestaurantBuilder) WithContactPhone(prefix, number string) HandlerRegisterRestaurantBuilder {
	nb := b
	nb.jsonPayload = strings.ReplaceAll(
		nb.jsonPayload,
		`"phone_prefix": "+1"`,
		`"phone_prefix": "`+prefix+`"`,
	)
	nb.jsonPayload = strings.ReplaceAll(
		nb.jsonPayload,
		`"phone_number": "1234567890"`,
		`"phone_number": "`+number+`"`,
	)
	return nb
}

// WithPostalCode sets all postal_code occurrences (restaurant.contact and staff_owner).
func (b HandlerRegisterRestaurantBuilder) WithPostalCode(postalCode string) HandlerRegisterRestaurantBuilder {
	nb := b
	nb.jsonPayload = strings.ReplaceAll(
		nb.jsonPayload,
		`"postal_code": "10001"`,
		`"postal_code": "`+postalCode+`"`,
	)
	return nb
}

// WithCountryCode sets all country_code occurrences (restaurant.contact and staff_owner).
func (b HandlerRegisterRestaurantBuilder) WithCountryCode(countryCode string) HandlerRegisterRestaurantBuilder {
	nb := b
	nb.jsonPayload = strings.ReplaceAll(
		nb.jsonPayload,
		`"country_code": "US"`,
		`"country_code": "`+countryCode+`"`,
	)
	return nb
}

// WithPassword sets the staff owner password.
func (b HandlerRegisterRestaurantBuilder) WithPassword(password string) HandlerRegisterRestaurantBuilder {
	nb := b
	nb.jsonPayload = strings.ReplaceAll(
		nb.jsonPayload,
		`"password": "strongpassword123"`,
		`"password": "`+password+`"`,
	)
	return nb
}

// WithVatCode sets restaurant.vat_code.
func (b HandlerRegisterRestaurantBuilder) WithVatCode(vat string) HandlerRegisterRestaurantBuilder {
	nb := b
	nb.jsonPayload = strings.ReplaceAll(
		nb.jsonPayload,
		`"vat_code": "GB123456789"`,
		`"vat_code": "`+vat+`"`,
	)
	return nb
}

// WithRestaurantName sets restaurant.name.
func (b HandlerRegisterRestaurantBuilder) WithRestaurantName(name string) HandlerRegisterRestaurantBuilder {
	nb := b
	nb.jsonPayload = strings.ReplaceAll(
		nb.jsonPayload,
		`"name": "Acme Pizza"`,
		`"name": "`+name+`"`,
	)
	return nb
}

// WithOwnerName sets staff_owner.name.
func (b HandlerRegisterRestaurantBuilder) WithOwnerName(name string) HandlerRegisterRestaurantBuilder {
	nb := b
	nb.jsonPayload = strings.ReplaceAll(
		nb.jsonPayload,
		`"name": "John Doe"`,
		`"name": "`+name+`"`,
	)
	return nb
}

// WithLegalName sets restaurant.legal_name.
func (b HandlerRegisterRestaurantBuilder) WithLegalName(legalName string) HandlerRegisterRestaurantBuilder {
	nb := b
	nb.jsonPayload = strings.ReplaceAll(
		nb.jsonPayload,
		`"legal_name": "Acme Pizza LLC"`,
		`"legal_name": "`+legalName+`"`,
	)
	return nb
}

// WithTaxID sets restaurant.tax_id.
func (b HandlerRegisterRestaurantBuilder) WithTaxID(taxID string) HandlerRegisterRestaurantBuilder {
	nb := b
	nb.jsonPayload = strings.ReplaceAll(
		nb.jsonPayload,
		`"tax_id": "99-1234567"`,
		`"tax_id": "`+taxID+`"`,
	)
	return nb
}

// WithAddress sets all address fields (restaurant.contact and staff_owner).
func (b HandlerRegisterRestaurantBuilder) WithAddress(addr string) HandlerRegisterRestaurantBuilder {
	nb := b
	nb.jsonPayload = strings.ReplaceAll(
		nb.jsonPayload,
		`"address": "123 Main St"`,
		`"address": "`+addr+`"`,
	)
	return nb
}

// WithCity sets all city fields (restaurant.contact and staff_owner).
func (b HandlerRegisterRestaurantBuilder) WithCity(city string) HandlerRegisterRestaurantBuilder {
	nb := b
	nb.jsonPayload = strings.ReplaceAll(
		nb.jsonPayload,
		`"city": "New York"`,
		`"city": "`+city+`"`,
	)
	return nb
}


