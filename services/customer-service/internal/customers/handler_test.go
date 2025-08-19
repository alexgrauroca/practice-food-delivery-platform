//go:build unit

package customers_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
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
	queryParams map[string]string
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
			runCustomerHandlerTestCase(t, logger, http.MethodPost, "/v1.0/customers", tt, "")
		})
	}
}

func TestHandler_GetCustomers(t *testing.T) {
	logger := setupTestEnv()
	now := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []customerHandlerTestCase{
		{
			name: "when invalid params are provided, then it should return a 400 with invalid request error",
			queryParams: map[string]string{
				"page": "invalid-page",
				"size": "invalid-size",
				"sort": "1",
			},
			wantJSON: `{
				"code": "INVALID_REQUEST",
				"message": "invalid request",
				"details": []
			}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "when invalid sort is provided, then it should return a 400 with the sort validation error",
			queryParams: map[string]string{
				"sort": "invalid",
			},
			wantJSON: `{
				"code": "VALIDATION_ERROR",
				"message": "validation failed",
				"details": [
					"sort must be one of name, email, created_at"
				]
			}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "when values are lower than the minimum, " +
				"then it should return a 400 with the min value validation errors",
			queryParams: map[string]string{
				"page":      "0",
				"page-size": "0",
			},
			wantJSON: `{
				"code": "VALIDATION_ERROR",
				"message": "validation failed",
				"details": [
					"page must be greater than 1",
					"page-size must be greater than 1"
				]
			}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "when values are higher than the maximum, " +
				"then it should return a 400 with the max value validation errors",
			queryParams: map[string]string{
				"page-size": "101",
			},
			wantJSON: `{
				"code": "VALIDATION_ERROR",
				"message": "validation failed",
				"details": [
					"page-size must be lower than 100"
				]
			}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "when unexpected error when searching customers, " +
				"then it should return a 500 with the internal error",
			queryParams: map[string]string{},
			mocksSetup: func(service *customersmocks.MockService) {
				// TODO: set the right function
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
			name:        "when there are no customers, then it should return a 200 with an empty list",
			queryParams: map[string]string{},
			mocksSetup: func(service *customersmocks.MockService) {
				// TODO: set the right function
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
				"items": [],
				"pagination": {
					"total_items": 0,
					"total_pages": 0,
					"current_page": 1,
					"page_size": 10
				}
			}`,
			wantStatus: http.StatusOK,
		},
		{
			name: "when there are customers, then it should return a 200 with the list of customers",
			queryParams: map[string]string{
				"page":      "1",
				"page-size": "2",
				"sort":      "name,-email",
			},
			mocksSetup: func(service *customersmocks.MockService) {
				// TODO: set the right function
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
				"items": [
					{
						"id": "fake-id",
						"name": "John Doe",
						"email": "test@example.com",
						"address": "a valid address",
						"city": "a valid city",
						"postal_code": "12345",
						"country_code": "US",
						"created_at": "2025-01-01T00:00:00Z",
						"updated_at": "2025-01-01T00:00:00Z"
					},
					{
						"id": "fake-id-2",
						"name": "John Doe 2",
						"email": "test2@example.com",
						"address": "a valid address 2",
						"city": "a valid city 2",
						"postal_code": "67890",
						"country_code": "ES",
						"created_at": "2025-01-01T00:00:00Z",
						"updated_at": "2025-01-01T00:00:00Z"
					}
				],
				"pagination": {
					"total_items": 3,
					"total_pages": 2,
					"current_page": 1,
					"page_size": 2
				}
			}`,
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runCustomerHandlerTestCase(t, logger, http.MethodGet, "/v1.0/customers", tt, "")
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
	httpMethod string,
	route string,
	tt customerHandlerTestCase,
	token string,
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

	// Create a new HTTP request with the test case's params
	var req *http.Request

	switch httpMethod {
	case http.MethodGet:
		baseURL, err := url.Parse(route)
		if err != nil {
			t.Fatalf("failed to parse route: %v", err)
		}

		// Add query parameters
		if len(tt.queryParams) > 0 {
			q := baseURL.Query()
			for k, v := range tt.queryParams {
				q.Add(k, v)
			}
			baseURL.RawQuery = q.Encode()
		}

		req = httptest.NewRequest(http.MethodGet, baseURL.String(), nil)

	case http.MethodPost:
		req = httptest.NewRequest(httpMethod, route, strings.NewReader(tt.jsonPayload))
		req.Header.Set("Content-Type", "application/json")

	default:
		t.Fatalf("unsupported HTTP method: %s", httpMethod)
	}

	if token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}

	w := httptest.NewRecorder()

	// Make the request to the handler
	router.ServeHTTP(w, req)

	assert.Equal(t, tt.wantStatus, w.Code)
	assert.JSONEq(t, tt.wantJSON, w.Body.String())
}
