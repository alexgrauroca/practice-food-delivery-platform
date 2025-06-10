//go:build !integration

package customers_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/customers"
	customersmocks "github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/customers/mocks"
	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/jwt"
)

type customerHandlerTestCase struct {
	name                 string
	jsonPayload          string
	mocksSetup           func(service *customersmocks.MockService)
	expectedJSONResponse string
	expectedStatusCode   int
}

func TestHandler_RegisterCustomer(t *testing.T) {
	logger := setupTestEnv()
	now := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []customerHandlerTestCase{
		{
			name:        "when invalid payload is provided, then it should return a 400 with invalid request error",
			jsonPayload: `{"name": 1.2, "email": true}`,
			expectedJSONResponse: `{
				"code": "INVALID_REQUEST",
				"message": "invalid request",
				"details": []
			}`,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:        "when empty payload is provided, then it should return a 400 with the validation error",
			jsonPayload: `{}`,
			expectedJSONResponse: `{
				"code": "VALIDATION_ERROR",
				"message": "validation failed",
				"details": [
					"email is required",
					"password is required",
					"name is required"
				]
			}`,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:        "when invalid email is provided, then it should return a 400 with the email validation error",
			jsonPayload: `{"email": "invalid-email", "name": "John Doe", "password": "ValidPassword123"}`,
			expectedJSONResponse: `{
				"code": "VALIDATION_ERROR",
				"message": "validation failed",
				"details": [
					"email must be a valid email address"
				]
			}`,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:        "when invalid password is provided, then it should return a 400 with the pwd validation error",
			jsonPayload: `{"email":"test@example.com", "name": "John Doe", "password": "short"}`,
			expectedJSONResponse: `{
				"code": "VALIDATION_ERROR",
				"message": "validation failed",
				"details": [
					"password must be a valid password with at least 8 characters long"
				]
			}`,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:        "when no name is provided, then it should return a 400 with the name validation error",
			jsonPayload: `{"email":"test@example.com", "password": "ValidPassword123"}`,
			expectedJSONResponse: `{
				"code":"VALIDATION_ERROR",
				"message":"validation failed",
				"details":[
					"name is required"
				]
			}`,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:        "when no email is provided, then it should return a 400 with the email validation error",
			jsonPayload: `{"name": "John Doe", "password": "ValidPassword123"}`,
			expectedJSONResponse: `{
				"code":"VALIDATION_ERROR",
				"message":"validation failed",
				"details":[
					"email is required"
				]
			}`,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:        "when no password is provided, then it should return a 400 with the email validation error",
			jsonPayload: `{"email": "test@example.com", "name": "John Doe"}`,
			expectedJSONResponse: `{
				"code":"VALIDATION_ERROR",
				"message":"validation failed",
				"details":[
					"password is required"
				]
			}`,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:        "when the customer already exists, then it should return a 409 with the customer already exists error",
			jsonPayload: `{"email": "test@example.com", "name": "John Doe", "password": "ValidPassword123"}`,
			mocksSetup: func(service *customersmocks.MockService) {
				service.EXPECT().RegisterCustomer(gomock.Any(), gomock.Any()).
					Return(customers.RegisterCustomerOutput{}, customers.ErrCustomerAlreadyExists)
			},
			expectedJSONResponse: `{
				"code": "CUSTOMER_ALREADY_EXISTS",
				"message": "customer already exists",
				"details": []
			}`,
			expectedStatusCode: http.StatusConflict,
		},
		{
			name:        "when unexpected error when registering the customer, then it should return a 500 with the internal error",
			jsonPayload: `{"email": "test@example.com", "name": "John Doe", "password": "ValidPassword123"}`,
			mocksSetup: func(service *customersmocks.MockService) {
				service.EXPECT().RegisterCustomer(gomock.Any(), gomock.Any()).
					Return(customers.RegisterCustomerOutput{}, errors.New("unexpected error"))
			},
			expectedJSONResponse: `{
				"code": "INTERNAL_ERROR",
				"message": "failed to register the customer",
				"details": []
			}`,
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name:        "when the customer is successfully registered, then it should return a 201 with the customer details",
			jsonPayload: `{"email": "test@example.com", "name": "John Doe", "password": "ValidPassword123"}`,
			mocksSetup: func(service *customersmocks.MockService) {
				service.EXPECT().RegisterCustomer(gomock.Any(), customers.RegisterCustomerInput{
					Email:    "test@example.com",
					Password: "ValidPassword123",
					Name:     "John Doe",
				}).Return(customers.RegisterCustomerOutput{
					ID:        "fake-id",
					Email:     "test@example.com",
					Name:      "John Doe",
					CreatedAt: now,
				}, nil)
			},
			expectedJSONResponse: `{
				"created_at":"2025-01-01T00:00:00Z",
				"email":"test@example.com",
				"id":"fake-id",
				"name":"John Doe"
			}`,
			expectedStatusCode: http.StatusCreated,
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
			expectedJSONResponse: `{
				"code": "INVALID_REQUEST",
				"message": "invalid request",
				"details": []
			}`,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:        "when empty payload is provided, then it should return a 400 with the validation error",
			jsonPayload: `{}`,
			expectedJSONResponse: `{
				"code": "VALIDATION_ERROR",
				"message": "validation failed",
				"details": [
					"email is required",
					"password is required"
				]
			}`,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:        "when invalid email is provided, then it should return a 400 with the email validation error",
			jsonPayload: `{"email": "invalid-email", "password": "ValidPassword123"}`,
			expectedJSONResponse: `{
				"code": "VALIDATION_ERROR",
				"message": "validation failed",
				"details": [
					"email must be a valid email address"
				]
			}`,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:        "when invalid password is provided, then it should return a 400 with the pwd validation error",
			jsonPayload: `{"email":"test@example.com", "password": "short"}`,
			expectedJSONResponse: `{
				"code": "VALIDATION_ERROR",
				"message": "validation failed",
				"details": [
					"password must be a valid password with at least 8 characters long"
				]
			}`,
			expectedStatusCode: http.StatusBadRequest,
		},

		{
			name:        "when there is not an active customer with the same email and password, then it should return a 401 with invalid credentials error",
			jsonPayload: `{"email": "test@example.com", "password": "ValidPassword123"}`,
			mocksSetup: func(service *customersmocks.MockService) {
				service.EXPECT().LoginCustomer(gomock.Any(), gomock.Any()).
					Return(customers.LoginCustomerOutput{}, customers.ErrInvalidCredentials)
			},
			expectedJSONResponse: `{
				"code": "INVALID_CREDENTIALS",
				"message": "invalid credentials",
				"details": []
			}`,
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name:        "when unexpected error when login the customer, then it should return a 500 with the internal error",
			jsonPayload: `{"email": "test@example.com", "password": "ValidPassword123"}`,
			mocksSetup: func(service *customersmocks.MockService) {
				service.EXPECT().LoginCustomer(gomock.Any(), gomock.Any()).
					Return(customers.LoginCustomerOutput{}, errors.New("unexpected error"))
			},
			expectedJSONResponse: `{
				"code": "INTERNAL_ERROR",
				"message": "failed to login the customer",
				"details": []
			}`,
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name:        "when an active customer has the same email and password, then it should return a 200 with the token",
			jsonPayload: `{"email": "test@example.com", "password": "ValidPassword123"}`,
			mocksSetup: func(service *customersmocks.MockService) {
				service.EXPECT().LoginCustomer(gomock.Any(), customers.LoginCustomerInput{
					Email:    "test@example.com",
					Password: "ValidPassword123",
				}).Return(customers.LoginCustomerOutput{
					AccessToken:  "fake-token",
					RefreshToken: "fake-refresh-token",
					ExpiresIn:    customers.DefaultTokenExpiration,
					TokenType:    jwt.DefaultTokenType,
				}, nil)
			},
			expectedJSONResponse: `{
			  "access_token": "fake-token",
			  "refresh_token": "fake-refresh-token",
			  "expires_in": 3600,
			  "token_type": "Bearer"
			}`,
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runCustomerHandlerTestCase(t, logger, "/v1.0/customers/login", tt)
		})
	}
}

// setupTestEnv initializes the test environment with default values common to all tests.
func setupTestEnv() *zap.Logger {
	// Setting up the default values
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()
	return logger
}

// runCustomerHandlerTestCase executes a test case for the customer handler, which is common for all tests.
func runCustomerHandlerTestCase(
	t *testing.T,
	logger *zap.Logger,
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

	assert.Equal(t, tt.expectedStatusCode, w.Code)
	assert.JSONEq(t, tt.expectedJSONResponse, w.Body.String())
}
