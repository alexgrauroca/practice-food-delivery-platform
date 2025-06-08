package customers_test

import (
	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/customers/mocks"
	"github.com/golang/mock/gomock"
	"go.uber.org/zap"

	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/customers"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHandler_RegisterCustomer(t *testing.T) {
	// Setting up the default values
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()

	now := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name                 string
		jsonPayload          string
		mocksSetup           func(service *mocks.MockService)
		expectedJsonResponse string
		expectedStatusCode   int
	}{
		{
			name:                 "when invalid payload is provided, then it should return a 400 with validation errors",
			jsonPayload:          `{`,
			expectedJsonResponse: `{"error": "Invalid request"}`,
			expectedStatusCode:   http.StatusBadRequest,
		},
		//TODO those error response were used to learn about how gin handles validation errors, but they should simpler
		{
			name:                 "when empty payload is provided, then it should return a 400 with the validation error",
			jsonPayload:          `{}`,
			expectedJsonResponse: `{"error":"Key: 'RegisterCustomerRequest.Email' Error:Field validation for 'Email' failed on the 'required' tag\nKey: 'RegisterCustomerRequest.Password' Error:Field validation for 'Password' failed on the 'required' tag\nKey: 'RegisterCustomerRequest.Name' Error:Field validation for 'Name' failed on the 'required' tag"}`,
			expectedStatusCode:   http.StatusBadRequest,
		},
		{
			name:                 "when invalid email is provided, then it should return a 400 with the email validation error",
			jsonPayload:          `{"email": "invalid-email", "name": "John Doe", "password": "ValidPassword123"}`,
			expectedJsonResponse: `{"error": "Key: 'RegisterCustomerRequest.Email' Error:Field validation for 'Email' failed on the 'email' tag"}`,
			expectedStatusCode:   http.StatusBadRequest,
		},
		{
			name:                 "when invalid password is provided, then it should return a 400 with the pwd validation error",
			jsonPayload:          `{"email":"test@example.com", "name": "John Doe", "password": "short"}`,
			expectedJsonResponse: `{"error":"Key: 'RegisterCustomerRequest.Password' Error:Field validation for 'Password' failed on the 'min' tag"}`,
			expectedStatusCode:   http.StatusBadRequest,
		},
		{
			name:                 "when no name is provided, then it should return a 400 with the name validation error",
			jsonPayload:          `{"email":"test@example.com", "password": "ValidPassword123"}`,
			expectedJsonResponse: `{"error":"Key: 'RegisterCustomerRequest.Name' Error:Field validation for 'Name' failed on the 'required' tag"}`,
			expectedStatusCode:   http.StatusBadRequest,
		},
		{
			name:                 "when no email is provided, then it should return a 400 with the email validation error",
			jsonPayload:          `{"name": "John Doe", "password": "ValidPassword123"}`,
			expectedJsonResponse: `{"error":"Key: 'RegisterCustomerRequest.Email' Error:Field validation for 'Email' failed on the 'required' tag"}`,
			expectedStatusCode:   http.StatusBadRequest,
		},
		{
			name:                 "when no password is provided, then it should return a 400 with the email validation error",
			jsonPayload:          `{"email": "test@example.com", "name": "John Doe"}`,
			expectedJsonResponse: `{"error":"Key: 'RegisterCustomerRequest.Password' Error:Field validation for 'Password' failed on the 'required' tag"}`,
			expectedStatusCode:   http.StatusBadRequest,
		},
		{
			name:        "when the customer already exists, then it should return a 409 with the customer already exists error",
			jsonPayload: `{"email": "test@example.com", "name": "John Doe", "password": "ValidPassword123"}`,
			mocksSetup: func(service *mocks.MockService) {
				service.EXPECT().RegisterCustomer(gomock.Any(), gomock.Any()).
					Return(customers.RegisterCustomerOutput{}, customers.ErrCustomerAlreadyExists)
			},
			expectedJsonResponse: `{"error": "Customer already exists"}`,
			expectedStatusCode:   http.StatusConflict,
		},
		{
			name:        "when the customer is successfully registered, then it should return a 201 with the customer details",
			jsonPayload: `{"email": "test@example.com", "name": "John Doe", "password": "ValidPassword123"}`,
			mocksSetup: func(service *mocks.MockService) {
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
			expectedJsonResponse: `{"created_at":"2025-01-01T00:00:00Z", "email":"test@example.com", "id":"fake-id", "name":"John Doe"}`,
			expectedStatusCode:   http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create new mock service
			service := mocks.NewMockService(gomock.NewController(t))
			if tt.mocksSetup != nil {
				tt.mocksSetup(service)
			}

			// Initialize the handler
			h := customers.NewHandler(logger, service)

			// Initialize the Gin router and register the routes
			router := gin.New()
			h.RegisterRoutes(router)

			// Create a new HTTP request with the test case's JSON payload
			req := httptest.NewRequest(http.MethodPost, "/v1.0/customers/register", strings.NewReader(tt.jsonPayload))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Make the request to the handler
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.JSONEq(t, tt.expectedJsonResponse, w.Body.String())
		})
	}
}
