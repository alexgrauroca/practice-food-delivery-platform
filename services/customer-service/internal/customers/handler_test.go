//go:build unit

package customers_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	customersmocks "github.com/alexgrauroca/practice-food-delivery-platform/services/customer-service/internal/customers/mocks"

	"github.com/alexgrauroca/practice-food-delivery-platform/services/customer-service/internal/customers"
	"github.com/alexgrauroca/practice-food-delivery-platform/services/customer-service/internal/log"
)

type customerHandlerTestCase struct {
	name        string
	jsonPayload string
	mocksSetup  func(service *customersmocks.MockService)
	wantJSON    string
	wantStatus  int
}

var errUnexpected = errors.New("unexpected error")

func TestHandler_RegisterCustomer(t *testing.T) {
	logger := setupTestEnv()
	now := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []customerHandlerTestCase{
		{
			name:        "when invalid payload is provided, then it should return a 400 with invalid request error",
			jsonPayload: `{"name": 1.2, "email": true}`,
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
					"email is required",
					"password is required",
					"name is required",
					"address is required",
					"city is required",
					"postal_code is required",
					"country_code is required"
				]
			}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "when invalid email is provided, then it should return a 400 with the email validation error",
			jsonPayload: `{
				"email": "invalid-email",
				"name": "John Doe",
				"password": "ValidPassword123",
				"address": "a valid address",
				"city": "a valid city",
				"postal_code": "12345",
				"country_code": "US"
			}`,
			wantJSON: `{
				"code": "VALIDATION_ERROR",
				"message": "validation failed",
				"details": [
					"email must be a valid email address"
				]
			}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "when fields length are shorten than minimum required, " +
				"then it should return a 400 with the short length validation errors",
			jsonPayload: `{
				"email": "test@example.com",
				"name": "John Doe",
				"password": "short",
				"address": "a valid address",
				"city": "a valid city",
				"postal_code": "1",
				"country_code": "U"
			}`,
			wantJSON: `{
				"code": "VALIDATION_ERROR",
				"message": "validation failed",
				"details": [
					"password must be a valid password with at least 8 characters long",
					"postal_code must be at least 5 characters long",
					"country_code must be at least 2 characters long"
				]
			}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "when fields length are longer than maximum required, " +
				"then it should return a 400 with the long length validation errors",
			jsonPayload: fmt.Sprintf(`{
					"email": "test@example.com",
					"name": "%s",
					"password": "ValidPassword123",
					"address": "%s",
					"city": "%s",
					"postal_code": "%s",
					"country_code": "USA"
				}`,
				strings.Repeat("a", 101),
				strings.Repeat("a", 101),
				strings.Repeat("a", 101),
				strings.Repeat("a", 33),
			),
			wantJSON: `{
				"code": "VALIDATION_ERROR",
				"message": "validation failed",
				"details": [
					"name must not exceed 100 characters long",
					"address must not exceed 100 characters long",
					"city must not exceed 100 characters long",
					"postal_code must not exceed 32 characters long",
					"country_code must not exceed 2 characters long"
				]
			}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "when the customer already exists, " +
				"then it should return a 409 with the customer already exists error",
			jsonPayload: `{
				"email": "test@example.com",
				"name": "John Doe",
				"password": "ValidPassword123",
				"address": "a valid address",
				"city": "a valid city",
				"postal_code": "12345",
				"country_code": "US"
			}`,
			mocksSetup: func(service *customersmocks.MockService) {
				service.EXPECT().RegisterCustomer(gomock.Any(), gomock.Any()).
					Return(customers.RegisterCustomerOutput{}, customers.ErrCustomerAlreadyExists)
			},
			wantJSON: `{
				"code": "CUSTOMER_ALREADY_EXISTS",
				"message": "customer already exists",
				"details": []
			}`,
			wantStatus: http.StatusConflict,
		},
		{
			name: "when unexpected error when registering the customer, " +
				"then it should return a 500 with the internal error",
			jsonPayload: `{
				"email": "test@example.com",
				"name": "John Doe",
				"password": "ValidPassword123",
				"address": "a valid address",
				"city": "a valid city",
				"postal_code": "12345",
				"country_code": "US"
			}`,
			mocksSetup: func(service *customersmocks.MockService) {
				service.EXPECT().RegisterCustomer(gomock.Any(), gomock.Any()).
					Return(customers.RegisterCustomerOutput{}, errUnexpected)
			},
			wantJSON: `{
				"code": "INTERNAL_ERROR",
				"message": "an unexpected error occurred",
				"details": []
			}`,
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "when the customer is successfully registered, " +
				"then it should return a 201 with the customer details",
			jsonPayload: `{
				"email": "test@example.com",
				"name": "John Doe",
				"password": "ValidPassword123",
				"address": "a valid address",
				"city": "a valid city",
				"postal_code": "12345",
				"country_code": "US"
			}`,
			mocksSetup: func(service *customersmocks.MockService) {
				service.EXPECT().RegisterCustomer(gomock.Any(), customers.RegisterCustomerInput{
					Email:       "test@example.com",
					Password:    "ValidPassword123",
					Name:        "John Doe",
					Address:     "a valid address",
					City:        "a valid city",
					PostalCode:  "12345",
					CountryCode: "US",
				}).Return(customers.RegisterCustomerOutput{
					ID:          "fake-id",
					Email:       "test@example.com",
					Name:        "John Doe",
					Address:     "a valid address",
					City:        "a valid city",
					PostalCode:  "12345",
					CountryCode: "US",
					CreatedAt:   now,
				}, nil)
			},
			wantJSON: `{
				"created_at":"2025-01-01T00:00:00Z",
				"email":"test@example.com",
				"id":"fake-id",
				"name":"John Doe",
				"address":"a valid address",
				"city":"a valid city",
				"postal_code":"12345",
				"country_code":"US"
			}`,
			wantStatus: http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runCustomerHandlerTestCase(t, logger, "/v1.0/customers", tt)
		})
	}
}

// setupTestEnv initializes the test environment with default values common to all tests.
func setupTestEnv() log.Logger {
	// Setting up the default values
	gin.SetMode(gin.TestMode)
	logger, _ := log.NewTest()
	return logger
}

// runCustomerHandlerTestCase executes a test case for the customer handler, which is common for all tests.
func runCustomerHandlerTestCase(
	t *testing.T,
	logger log.Logger,
	route string,
	tt customerHandlerTestCase,
) {
	// Create a new mock service
	service := customersmocks.NewMockService(gomock.NewController(t))
	if tt.mocksSetup != nil {
		tt.mocksSetup(service)
	}

	// Initialize the handler
	h := customers.NewHandler(logger, service)

	// Initialize the Gin router and register the routes
	router := gin.New()
	h.RegisterRoutes(router)

	// Create a new HTTP request with the test case's JSON payload
	req := httptest.NewRequest(http.MethodPost, route, strings.NewReader(tt.jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Make the request to the handler
	router.ServeHTTP(w, req)

	assert.Equal(t, tt.wantStatus, w.Code)
	assert.JSONEq(t, tt.wantJSON, w.Body.String())
}
