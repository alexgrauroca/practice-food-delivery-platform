//go:build unit

package customers_test

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/alexgrauroca/practice-food-delivery-platform/pkg/auth"
	authmocks "github.com/alexgrauroca/practice-food-delivery-platform/pkg/auth/mocks"
	customhttp "github.com/alexgrauroca/practice-food-delivery-platform/pkg/http"
	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/authcore"
	customersmocks "github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/customers/mocks"

	"github.com/alexgrauroca/practice-food-delivery-platform/pkg/log"
	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/customers"
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
	logger := customhttp.SetupTestEnv()
	now := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []customerHandlerTestCase{
		{
			name:        "when invalid payload is provided, then it should return a 400 with invalid request error",
			jsonPayload: `{"password": 1.2, "email": true}`,
			wantJSON:    customhttp.NewInvalidRequestRespBuilder().Build(),
			wantStatus:  http.StatusBadRequest,
		},
		{
			name:        "when empty payload is provided, then it should return a 400 with the validation error",
			jsonPayload: `{}`,
			wantJSON: customhttp.NewValidationErrorRespBuilder().
				WithDetails(
					"customer_id is required",
					"email is required",
					"password is required",
				).Build(),
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "when invalid email is provided, then it should return a 400 with the email validation error",
			jsonPayload: `{
				"customer_id": "fake-customer-id",
				"email": "invalid-email",
				"password": "ValidPassword123"
			}`,
			wantJSON: customhttp.NewValidationErrorRespBuilder().
				WithDetails("email must be a valid email address").
				Build(),
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "when invalid password is provided, then it should return a 400 with the pwd validation error",
			jsonPayload: `{
				"customer_id": "fake-customer-id",
				"email": "test@example.com",
				"password": "short"
			}`,
			wantJSON: customhttp.NewValidationErrorRespBuilder().
				WithDetails("password must be a valid password with at least 8 characters long").
				Build(),
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "when the customer already exists, then it should return a 409 with the customer already exists error",
			jsonPayload: `{
				"customer_id": "fake-customer-id",
				"email": "test@example.com",
				"password": "ValidPassword123"
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
			name: "when unexpected error when registering the customer, then it should return a 500 with the internal error",
			jsonPayload: `{
				"customer_id": "fake-customer-id",
				"email": "test@example.com",
				"password": "ValidPassword123"
			}`,
			mocksSetup: func(service *customersmocks.MockService, _ *authmocks.MockService) {
				service.EXPECT().RegisterCustomer(gomock.Any(), gomock.Any()).
					Return(customers.RegisterCustomerOutput{}, errUnexpected)
			},
			wantJSON:   customhttp.NewInternalErrorRespBuilder().Build(),
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "when the customer is successfully registered, then it should return a 201 with the customer details",
			jsonPayload: `{
				"customer_id": "fake-customer-id",
				"email": "test@example.com",
				"password": "ValidPassword123"
			}`,
			mocksSetup: func(service *customersmocks.MockService, _ *authmocks.MockService) {
				service.EXPECT().RegisterCustomer(gomock.Any(), customers.RegisterCustomerInput{
					CustomerID: "fake-customer-id",
					Email:      "test@example.com",
					Password:   "ValidPassword123",
				}).Return(customers.RegisterCustomerOutput{
					ID:        "fake-id",
					Email:     "test@example.com",
					CreatedAt: now,
				}, nil)
			},
			wantJSON: `{
				"created_at":"2025-01-01T00:00:00Z",
				"email":"test@example.com",
				"id":"fake-id"
			}`,
			wantStatus: http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runCustomerHandlerTestCase(t, logger, http.MethodPost, "/v1.0/auth/customers", tt, "")
		})
	}
}

func TestHandler_LoginCustomer(t *testing.T) {
	logger := customhttp.SetupTestEnv()

	tests := []customerHandlerTestCase{
		{
			name:        "when invalid payload is provided, then it should return a 400 with invalid request error",
			jsonPayload: `{"name": 1.2, "email": true}`,
			wantJSON:    customhttp.NewInvalidRequestRespBuilder().Build(),
			wantStatus:  http.StatusBadRequest,
		},
		{
			name:        "when empty payload is provided, then it should return a 400 with the validation error",
			jsonPayload: `{}`,
			wantJSON: customhttp.NewValidationErrorRespBuilder().
				WithDetails(
					"email is required",
					"password is required",
				).Build(),
			wantStatus: http.StatusBadRequest,
		},
		{
			name:        "when invalid email is provided, then it should return a 400 with the email validation error",
			jsonPayload: `{"email": "invalid-email", "password": "ValidPassword123"}`,
			wantJSON: customhttp.NewValidationErrorRespBuilder().
				WithDetails("email must be a valid email address").
				Build(),
			wantStatus: http.StatusBadRequest,
		},
		{
			name:        "when invalid password is provided, then it should return a 400 with the pwd validation error",
			jsonPayload: `{"email":"test@example.com", "password": "short"}`,
			wantJSON: customhttp.NewValidationErrorRespBuilder().
				WithDetails("password must be a valid password with at least 8 characters long").
				Build(),
			wantStatus: http.StatusBadRequest,
		},

		{
			name:        "when there is not an active customer with the same email and password, then it should return a 401 with invalid credentials error",
			jsonPayload: `{"email": "test@example.com", "password": "ValidPassword123"}`,
			mocksSetup: func(service *customersmocks.MockService, _ *authmocks.MockService) {
				service.EXPECT().LoginCustomer(gomock.Any(), gomock.Any()).
					Return(customers.LoginCustomerOutput{}, authcore.ErrInvalidCredentials)
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
			mocksSetup: func(service *customersmocks.MockService, _ *authmocks.MockService) {
				service.EXPECT().LoginCustomer(gomock.Any(), gomock.Any()).
					Return(customers.LoginCustomerOutput{}, errUnexpected)
			},
			wantJSON:   customhttp.NewInternalErrorRespBuilder().Build(),
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:        "when an active customer has the same email and password, then it should return a 200 with the token",
			jsonPayload: `{"email": "test@example.com", "password": "ValidPassword123"}`,
			mocksSetup: func(service *customersmocks.MockService, _ *authmocks.MockService) {
				service.EXPECT().LoginCustomer(gomock.Any(), customers.LoginCustomerInput{
					Email:    "test@example.com",
					Password: "ValidPassword123",
				}).Return(customers.LoginCustomerOutput{
					TokenPair: authcore.TokenPair{
						AccessToken:  "fake-token",
						RefreshToken: "fake-refresh-token",
						ExpiresIn:    customers.DefaultTokenExpiration,
						TokenType:    auth.DefaultTokenType,
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
			runCustomerHandlerTestCase(t, logger, http.MethodPost, "/v1.0/customers/login", tt, "")
		})
	}
}

func TestHandler_RefreshCustomer(t *testing.T) {
	logger := customhttp.SetupTestEnv()

	tests := []customerHandlerTestCase{
		{
			name:        "when invalid payload is provided, then it should return a 400 with invalid request error",
			jsonPayload: `{"access_token": 1.2, "refresh_token": true}`,
			wantJSON:    customhttp.NewInvalidRequestRespBuilder().Build(),
			wantStatus:  http.StatusBadRequest,
		},
		{
			name:        "when empty payload is provided, then it should return a 400 with the validation error",
			jsonPayload: `{}`,
			wantJSON: customhttp.NewValidationErrorRespBuilder().
				WithDetails(
					"refresh_token is required",
					"access_token is required",
				).Build(),
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "when invalid refresh token provided, " +
				"then it should return a 401 with the invalid refresh token error",
			jsonPayload: `{"access_token": "valid-access-token", "refresh_token": "invalid-refresh-token"}`,
			mocksSetup: func(service *customersmocks.MockService, _ *authmocks.MockService) {
				service.EXPECT().RefreshCustomer(gomock.Any(), gomock.Any()).
					Return(customers.RefreshCustomerOutput{}, authcore.ErrInvalidRefreshToken)
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
			mocksSetup: func(service *customersmocks.MockService, _ *authmocks.MockService) {
				service.EXPECT().RefreshCustomer(gomock.Any(), gomock.Any()).
					Return(customers.RefreshCustomerOutput{}, authcore.ErrTokenMismatch)
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
			mocksSetup: func(service *customersmocks.MockService, _ *authmocks.MockService) {
				service.EXPECT().RefreshCustomer(gomock.Any(), gomock.Any()).
					Return(customers.RefreshCustomerOutput{}, errUnexpected)
			},
			wantJSON:   customhttp.NewInternalErrorRespBuilder().Build(),
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:        "when the customer token is refreshed, then it should return a 200 with the new token",
			jsonPayload: `{"access_token": "valid-access-token", "refresh_token": "valid-refresh-token"}`,
			mocksSetup: func(service *customersmocks.MockService, _ *authmocks.MockService) {
				service.EXPECT().RefreshCustomer(gomock.Any(), customers.RefreshCustomerInput{
					AccessToken:  "valid-access-token",
					RefreshToken: "valid-refresh-token",
				}).Return(customers.RefreshCustomerOutput{
					TokenPair: authcore.TokenPair{
						AccessToken:  "fake-token",
						RefreshToken: "fake-refresh-token",
						ExpiresIn:    customers.DefaultTokenExpiration,
						TokenType:    auth.DefaultTokenType,
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
			runCustomerHandlerTestCase(t, logger, http.MethodPost, "/v1.0/customers/refresh", tt, "")
		})
	}
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
	authService := authmocks.NewMockService(gomock.NewController(t))
	if tt.mocksSetup != nil {
		tt.mocksSetup(service, authService)
	}

	// Initialize the authentication middleware
	authMiddleware := auth.NewMiddleware(logger, authService)

	// Initialize the handler
	h := customers.NewHandler(logger, service, authMiddleware)

	// Make HTTP request
	w := customhttp.ServeTestHTTPRequest(t, h, httpMethod, route, token, tt.queryParams, tt.jsonPayload)

	assert.Equal(t, tt.wantStatus, w.Code)
	assert.JSONEq(t, tt.wantJSON, w.Body.String())
}
