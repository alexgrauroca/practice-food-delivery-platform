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

	"github.com/alexgrauroca/practice-food-delivery-platform/services/customer-service/internal/authentication"
	authmocks "github.com/alexgrauroca/practice-food-delivery-platform/services/customer-service/internal/authentication/mocks"
	customersmocks "github.com/alexgrauroca/practice-food-delivery-platform/services/customer-service/internal/customers/mocks"

	"github.com/alexgrauroca/practice-food-delivery-platform/services/customer-service/internal/customers"
	"github.com/alexgrauroca/practice-food-delivery-platform/services/customer-service/internal/log"
)

type customerHandlerTestCase struct {
	name        string
	token       string
	pathParams  map[string]string
	queryParams map[string]string
	jsonPayload string
	mocksSetup  func(service *customersmocks.MockService, authService *authmocks.MockService)
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
			mocksSetup: func(service *customersmocks.MockService, _ *authmocks.MockService) {
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
			mocksSetup: func(service *customersmocks.MockService, _ *authmocks.MockService) {
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
			mocksSetup: func(service *customersmocks.MockService, _ *authmocks.MockService) {
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

func TestHandler_GetCustomer(t *testing.T) {
	logger := setupTestEnv()
	now := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []customerHandlerTestCase{
		{
			name:  "when any token is provided, then it should return a 401 with the unauthorized error",
			token: "",
			pathParams: map[string]string{
				"customerID": "fakeID",
			},
			wantJSON: `{
				"code": "UNAUTHORIZED",
				"message": "Authentication is required to access this resource",
				"details": []
			}`,
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:  "when invalid token is provided, then it should return a 401 with the unauthorized error",
			token: "invalid-token",
			pathParams: map[string]string{
				"customerID": "fakeID",
			},
			mocksSetup: func(_ *customersmocks.MockService, authService *authmocks.MockService) {
				authService.EXPECT().ValidateAccessToken(gomock.Any(), gomock.Any()).
					Return(authentication.ValidateAccessTokenOutput{}, authentication.ErrInvalidToken)
			},
			wantJSON: `{
				"code": "UNAUTHORIZED",
				"message": "Authentication is required to access this resource",
				"details": []
			}`,
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "when authenticated customer is not the same as the one requested, " +
				"then it should return a 403 with the forbidden error",
			token: "none-customer-token",
			pathParams: map[string]string{
				"customerID": "fakeID",
			},
			mocksSetup: func(_ *customersmocks.MockService, authService *authmocks.MockService) {
				authService.EXPECT().ValidateAccessToken(gomock.Any(), gomock.Any()).
					Return(authentication.ValidateAccessTokenOutput{
						Claims: &authentication.Claims{
							Role: "none-customer-role",
						},
					}, nil)
			},
			wantJSON: `{
				"code": "FORBIDDEN",
				"message": "You do not have permission to access this resource",
				"details": []
			}`,
			wantStatus: http.StatusForbidden,
		},
		{
			name:  "when the customer is not found, then it should return a 404 with the not found error",
			token: "valid-token",
			pathParams: map[string]string{
				"customerID": "unexistingID",
			},
			mocksSetup: func(service *customersmocks.MockService, authService *authmocks.MockService) {
				authService.EXPECT().ValidateAccessToken(gomock.Any(), gomock.Any()).
					Return(authentication.ValidateAccessTokenOutput{
						Claims: &authentication.Claims{
							Role: string(authentication.RoleCustomer),
						},
					}, nil)

				service.EXPECT().GetCustomer(gomock.Any(), gomock.Any()).
					Return(customers.GetCustomerOutput{}, customers.ErrCustomerNotFound)
			},
			wantJSON: `{
				"code": "NOT_FOUND",
				"message": "resource not found",
				"details": []
			}`,
			wantStatus: http.StatusNotFound,
		},
		{
			name:  "when a valid customerID is provided, then it should return a 200 with the customer details",
			token: "valid-token",
			pathParams: map[string]string{
				"customerID": "fakeID",
			},
			mocksSetup: func(service *customersmocks.MockService, authService *authmocks.MockService) {
				authService.EXPECT().ValidateAccessToken(gomock.Any(), authentication.ValidateAccessTokenInput{
					AccessToken: "valid-token",
				}).Return(authentication.ValidateAccessTokenOutput{
					Claims: &authentication.Claims{
						Role: string(authentication.RoleCustomer),
					},
				}, nil)

				service.EXPECT().GetCustomer(gomock.Any(), customers.GetCustomerInput{CustomerID: "fakeID"}).
					Return(customers.GetCustomerOutput{
						ID:          "fakeID",
						Name:        "John Doe",
						Email:       "test@example.com",
						Address:     "123 Main St",
						City:        "New York",
						PostalCode:  "10001",
						CountryCode: "US",
						CreatedAt:   now,
						UpdatedAt:   now,
					}, nil)
			},
			wantJSON: `{
				"id": "fakeID",
				"name": "John Doe",
			    "email": "test@example.com",
			    "address": "123 Main St",
			    "city": "New York",
			    "postal_code": "10001",
			    "country_code": "US",
				"created_at": "2025-01-01T00:00:00Z",
				"updated_at": "2025-01-01T00:00:00Z"
			}`,
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getCustomerPath := fmt.Sprintf("/v1.0/customers/%s", tt.pathParams["customerID"])
			runCustomerHandlerTestCase(t, logger, http.MethodGet, getCustomerPath, tt, tt.token)
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
	service := customersmocks.NewMockService(gomock.NewController(t))
	authService := authmocks.NewMockService(gomock.NewController(t))
	if tt.mocksSetup != nil {
		tt.mocksSetup(service, authService)
	}

	// Initialize the authentication middleware
	authMiddleware := authentication.NewMiddleware(logger, authService)

	// Initialize the handler
	h := customers.NewHandler(logger, service, authMiddleware)

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
