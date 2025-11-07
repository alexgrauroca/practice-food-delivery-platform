package restaurants_test

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	customhttp "github.com/alexgrauroca/practice-food-delivery-platform/pkg/http"
	"github.com/alexgrauroca/practice-food-delivery-platform/pkg/log"
	"github.com/alexgrauroca/practice-food-delivery-platform/services/restaurant-service/internal/restaurants"
)

var errUnexpected = errors.New("unexpected error")

type handlerTestCase struct {
	name        string
	token       string
	pathParams  map[string]string
	queryParams map[string]string
	jsonPayload string
	mocksSetup  func()
	wantJSON    string
	wantStatus  int
}

func TestHandler_RegisterRestaurant(t *testing.T) {
	logger := customhttp.SetupTestEnv()
	_ = time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []handlerTestCase{
		{
			name:        "when invalid payload is provided, then it should return a 400 with invalid request error",
			jsonPayload: `{"restaurant": 1.2, "staff_owner": true}`,
			wantJSON: `{
				"code": "INVALID_REQUEST",
				"message": "invalid request",
				"details": []
			}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:        "when empty payload is provided, then it should return a 400 with the validation error",
			jsonPayload: `{}`,
			wantJSON: `{
				"code": "VALIDATION_ERROR",
				"message": "validation failed",
				"details": [
					"restaurant.vat_code is required",
					"restaurant.name is required",
					"restaurant.legal_name is required",
					"restaurant.timezone_id is required",
					"restaurant.contact.phone_prefix is required",
					"restaurant.contact.phone_number is required",
					"restaurant.contact.email is required",
					"restaurant.contact.address is required",
					"restaurant.contact.city is required",
					"restaurant.contact.postal_code is required",
					"restaurant.contact.country_code is required",
					"staff_owner.email is required",
					"staff_owner.password is required",
					"staff_owner.name is required",
					"staff_owner.address is required",
					"staff_owner.city is required",
					"staff_owner.postal_code is required",
					"staff_owner.country_code is required"
				]
			}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "when invalid email is provided, then it should return a 400 with the email validation error",
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
						"email": "invalid-email",
						"address": "123 Main St",
						"city": "New York",
						"postal_code": "10001",
						"country_code": "US"
					}
				},
				"staff_owner": {
					"email": "invalid-email-2",
					"password": "strongpassword123",
					"name": "John Doe",
					"address": "123 Main St",
					"city": "New York",
					"postal_code": "10001",
					"country_code": "US"
				}
			}`,
			wantJSON: `{
				"code": "VALIDATION_ERROR",
				"message": "validation failed",
				"details": [
					"restaurant.contact.email must be a valid email address",
					"staff_owner.email must be a valid email address"
				]
			}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "when invalid timezone is provided, then it should return a 400 with the timezone validation error",
			jsonPayload: `{
				"restaurant": {
					"vat_code": "GB123456789",
					"name": "Acme Pizza",
					"legal_name": "Acme Pizza LLC",
					"tax_id": "99-1234567",
					"timezone_id": "Invalid/Timezone",
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
			wantJSON: `{
				"code": "VALIDATION_ERROR",
				"message": "validation failed",
				"details": [ "restaurant.timezone_id is invalid" ]
			}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "when invalid phone data is provided, then it should return a 400 with the phone validation error",
			jsonPayload: `{
				"restaurant": {
					"vat_code": "GB123456789",
					"name": "Acme Pizza",
					"legal_name": "Acme Pizza LLC",
					"tax_id": "99-1234567",
					"timezone_id": "America/New_York",
					"contact": {
						"phone_prefix": "a",
						"phone_number": "bbbbbbb",
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
			wantJSON: `{
				"code": "VALIDATION_ERROR",
				"message": "validation failed",
				"details": [
					"restaurant.contact.phone_prefix is invalid",
					"restaurant.contact.phone_number is invalid"
				]
			}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "when fields length are shorten than minimum required, " +
				"then it should return a 400 with the short length validation errors",
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
						"postal_code": "1",
						"country_code": "U"
					}
				},
				"staff_owner": {
					"email": "user@example.com",
					"password": "short",
					"name": "John Doe",
					"address": "123 Main St",
					"city": "New York",
					"postal_code": "1",
					"country_code": "U"
				}
			}`,
			wantJSON: `{
				"code": "VALIDATION_ERROR",
				"message": "validation failed",
				"details": [
					"restaurant.contact.postal_code must be at least 5 characters long",
					"restaurant.contact.country_code must be at least 2 characters long",
					"staff_owner.password must be a valid password with at least 8 characters long",
					"staff_owner.postal_code must be at least 5 characters long",
					"staff_owner.country_code must be at least 2 characters long"
				]
			}`,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runHandlerTestCase(t, logger, http.MethodPost, "/v1.0/restaurants", tt, "")
		})
	}
}

func runHandlerTestCase(
	t *testing.T,
	logger log.Logger,
	httpMethod string,
	route string,
	tt handlerTestCase,
	token string,
) {
	if tt.mocksSetup != nil {
		tt.mocksSetup()
	}

	// Initialize the handler
	h := restaurants.NewHandler(logger)

	// Make HTTP request
	w := customhttp.ServeTestHTTPRequest(t, h, httpMethod, route, token, tt.queryParams, tt.jsonPayload)

	assert.Equal(t, tt.wantStatus, w.Code)
	assert.JSONEq(t, tt.wantJSON, w.Body.String())
}
