//go:build unit

package customers_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	customersmocks "github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/customers/mocks"

	"github.com/alexgrauroca/practice-food-delivery-platform/pkg/log"
	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/customers"
	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/jwt"
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
					"customer_id is required",
					"email is required",
					"password is required",
					"name is required"
				]
			}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "when invalid email is provided, then it should return a 400 with the email validation error",
			jsonPayload: `{
				"customer_id": "fake-customer-id",
				"email": "invalid-email",
				"name": "John Doe", 
				"password": "ValidPassword123"
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
			name: "when invalid password is provided, then it should return a 400 with the pwd validation error",
			jsonPayload: `{
				"customer_id": "fake-customer-id",
				"email": "test@example.com",
				"name": "John Doe", 
				"password": "short"
			}`,
			wantJSON: `{
				"code": "VALIDATION_ERROR",
				"message": "validation failed",
				"details": [
					"password must be a valid password with at least 8 characters long"
				]
			}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "when no name is provided, then it should return a 400 with the name validation error",
			jsonPayload: `{
				"customer_id": "fake-customer-id",
				"email": "test@example.com",
				"password": "ValidPassword123"
			}`,
			wantJSON: `{
				"code":"VALIDATION_ERROR",
				"message":"validation failed",
				"details":[
					"name is required"
				]
			}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "when no email is provided, then it should return a 400 with the email validation error",
			jsonPayload: `{
				"customer_id": "fake-customer-id",
				"name": "John Doe", 
				"password": "ValidPassword123"
			}`,
			wantJSON: `{
				"code":"VALIDATION_ERROR",
				"message":"validation failed",
				"details":[
					"email is required"
				]
			}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "when no password is provided, then it should return a 400 with the email validation error",
			jsonPayload: `{
				"customer_id": "fake-customer-id",
				"email": "test@example.com",
				"name": "John Doe"
			}`,
			wantJSON: `{
				"code":"VALIDATION_ERROR",
				"message":"validation failed",
				"details":[
					"password is required"
				]
			}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "when the customer already exists, then it should return a 409 with the customer already exists error",
			jsonPayload: `{
				"customer_id": "fake-customer-id",
				"email": "test@example.com",
				"name": "John Doe", 
				"password": "ValidPassword123"
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
			name: "when unexpected error when registering the customer, then it should return a 500 with the internal error",
			jsonPayload: `{
				"customer_id": "fake-customer-id",
				"email": "test@example.com",
				"name": "John Doe", 
				"password": "ValidPassword123"
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
			name: "when the customer is successfully registered, then it should return a 201 with the customer details",
			jsonPayload: `{
				"customer_id": "fake-customer-id",
				"email": "test@example.com",
				"name": "John Doe", 
				"password": "ValidPassword123"
			}`,
			mocksSetup: func(service *customersmocks.MockService) {
				service.EXPECT().RegisterCustomer(gomock.Any(), customers.RegisterCustomerInput{
					CustomerID: "fake-customer-id",
					Email:      "test@example.com",
					Password:   "ValidPassword123",
					Name:       "John Doe",
				}).Return(customers.RegisterCustomerOutput{
					ID:        "fake-id",
					Email:     "test@example.com",
					Name:      "John Doe",
					CreatedAt: now,
				}, nil)
			},
			wantJSON: `{
				"created_at":"2025-01-01T00:00:00Z",
				"email":"test@example.com",
				"id":"fake-id",
				"name":"John Doe"
			}`,
			wantStatus: http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runCustomerHandlerTestCase(t, logger, "/v1.0/customers/register", tt)
		})
	}
}

func TestHandler_LoginCustomer(t *testing.T) {
	logger := setupTestEnv()

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
					"password is required"
				]
			}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:        "when invalid email is provided, then it should return a 400 with the email validation error",
			jsonPayload: `{"email": "invalid-email", "password": "ValidPassword123"}`,
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
			name:        "when invalid password is provided, then it should return a 400 with the pwd validation error",
			jsonPayload: `{"email":"test@example.com", "password": "short"}`,
			wantJSON: `{
				"code": "VALIDATION_ERROR",
				"message": "validation failed",
				"details": [
					"password must be a valid password with at least 8 characters long"
				]
			}`,
			wantStatus: http.StatusBadRequest,
		},

		{
			name:        "when there is not an active customer with the same email and password, then it should return a 401 with invalid credentials error",
			jsonPayload: `{"email": "test@example.com", "password": "ValidPassword123"}`,
			mocksSetup: func(service *customersmocks.MockService) {
				service.EXPECT().LoginCustomer(gomock.Any(), gomock.Any()).
					Return(customers.LoginCustomerOutput{}, customers.ErrInvalidCredentials)
			},
			wantJSON: `{
				"code": "INVALID_CREDENTIALS",
				"message": "invalid credentials",
				"details": []
			}`,
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:        "when unexpected error when login the customer, then it should return a 500 with the internal error",
			jsonPayload: `{"email": "test@example.com", "password": "ValidPassword123"}`,
			mocksSetup: func(service *customersmocks.MockService) {
				service.EXPECT().LoginCustomer(gomock.Any(), gomock.Any()).
					Return(customers.LoginCustomerOutput{}, errUnexpected)
			},
			wantJSON: `{
				"code": "INTERNAL_ERROR",
				"message": "an unexpected error occurred",
				"details": []
			}`,
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:        "when an active customer has the same email and password, then it should return a 200 with the token",
			jsonPayload: `{"email": "test@example.com", "password": "ValidPassword123"}`,
			mocksSetup: func(service *customersmocks.MockService) {
				service.EXPECT().LoginCustomer(gomock.Any(), customers.LoginCustomerInput{
					Email:    "test@example.com",
					Password: "ValidPassword123",
				}).Return(customers.LoginCustomerOutput{
					customers.TokenPair{
						AccessToken:  "fake-token",
						RefreshToken: "fake-refresh-token",
						ExpiresIn:    customers.DefaultTokenExpiration,
						TokenType:    jwt.DefaultTokenType,
					},
				}, nil)
			},
			wantJSON: `{
			  "access_token": "fake-token",
			  "refresh_token": "fake-refresh-token",
			  "expires_in": 3600,
			  "token_type": "Bearer"
			}`,
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runCustomerHandlerTestCase(t, logger, "/v1.0/customers/login", tt)
		})
	}
}

func TestHandler_RefreshCustomer(t *testing.T) {
	logger := setupTestEnv()

	tests := []customerHandlerTestCase{
		{
			name:        "when invalid payload is provided, then it should return a 400 with invalid request error",
			jsonPayload: `{"access_token": 1.2, "refresh_token": true}`,
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
					"refresh_token is required",
					"access_token is required"
				]
			}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "when invalid refresh token provided, " +
				"then it should return a 401 with the invalid refresh token error",
			jsonPayload: `{"access_token": "valid-access-token", "refresh_token": "invalid-refresh-token"}`,
			mocksSetup: func(service *customersmocks.MockService) {
				service.EXPECT().RefreshCustomer(gomock.Any(), gomock.Any()).
					Return(customers.RefreshCustomerOutput{}, customers.ErrInvalidRefreshToken)
			},
			wantJSON: `{
				"code": "INVALID_REFRESH_TOKEN",
				"message": "invalid or expired refresh token",
				"details": []
			}`,
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "when there is a token mismatch between the access token and the refresh token, " +
				"then it should return a 403 with the token mismatch error",
			jsonPayload: `{"access_token": "invalid-access-token", "refresh_token": "valid-refresh-token"}`,
			mocksSetup: func(service *customersmocks.MockService) {
				service.EXPECT().RefreshCustomer(gomock.Any(), gomock.Any()).
					Return(customers.RefreshCustomerOutput{}, customers.ErrTokenMismatch)
			},
			wantJSON: `{
				"code": "TOKEN_MISMATCH",
				"message": "token mismatch",
				"details": []
			}`,
			wantStatus: http.StatusForbidden,
		},
		{
			name: "when unexpected error when refreshing the customer token, " +
				"then it should return a 500 with the internal error",
			jsonPayload: `{"access_token": "valid-access-token", "refresh_token": "valid-refresh-token"}`,
			mocksSetup: func(service *customersmocks.MockService) {
				service.EXPECT().RefreshCustomer(gomock.Any(), gomock.Any()).
					Return(customers.RefreshCustomerOutput{}, errUnexpected)
			},
			wantJSON: `{
				"code": "INTERNAL_ERROR",
				"message": "an unexpected error occurred",
				"details": []
			}`,
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:        "when the customer token is refreshed, then it should return a 200 with the new token",
			jsonPayload: `{"access_token": "valid-access-token", "refresh_token": "valid-refresh-token"}`,
			mocksSetup: func(service *customersmocks.MockService) {
				service.EXPECT().RefreshCustomer(gomock.Any(), customers.RefreshCustomerInput{
					AccessToken:  "valid-access-token",
					RefreshToken: "valid-refresh-token",
				}).Return(customers.RefreshCustomerOutput{
					TokenPair: customers.TokenPair{
						AccessToken:  "fake-token",
						RefreshToken: "fake-refresh-token",
						ExpiresIn:    customers.DefaultTokenExpiration,
						TokenType:    jwt.DefaultTokenType,
					},
				}, nil)
			},
			wantJSON: `{
			  "access_token": "fake-token",
			  "refresh_token": "fake-refresh-token",
			  "expires_in": 3600,
			  "token_type": "Bearer"
			}`,
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runCustomerHandlerTestCase(t, logger, "/v1.0/customers/refresh", tt)
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
