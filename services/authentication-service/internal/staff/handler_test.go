//go:build unit

package staff_test

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
	"github.com/alexgrauroca/practice-food-delivery-platform/pkg/log"
	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/authcore"
	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/customers"
	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/staff"
	staffmocks "github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/staff/mocks"
)

type staffHandlerTestCase struct {
	name        string
	token       string
	pathParams  map[string]string
	queryParams map[string]string
	jsonPayload string
	mocksSetup  func(service *staffmocks.MockService, authService *authmocks.MockService)
	wantJSON    string
	wantStatus  int
}

var errUnexpected = errors.New("unexpected error")

func TestHandler_RegisterStaff(t *testing.T) {
	logger := customhttp.SetupTestEnv()
	now := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []staffHandlerTestCase{
		{
			name:        "when invalid payload is provided, then it should return a 400 with invalid request error",
			jsonPayload: `{"password": 1.2, "email": true}`,
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
					"staff_id is required",
					"email is required",
					"password is required"
				]
			}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "when invalid email is provided, then it should return a 400 with the email validation error",
			jsonPayload: `{
				"staff_id": "fake-staff-id",
				"email": "invalid-email",
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
				"staff_id": "fake-staff-id",
				"email": "test@example.com",
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
			name: "when the staff already exists, then it should return a 409 with the staff already exists error",
			jsonPayload: `{
				"staff_id": "fake-staff-id",
				"email": "test@example.com",
				"password": "ValidPassword123"
			}`,
			mocksSetup: func(service *staffmocks.MockService, _ *authmocks.MockService) {
				service.EXPECT().RegisterStaff(gomock.Any(), gomock.Any()).
					Return(staff.RegisterStaffOutput{}, staff.ErrStaffAlreadyExists)
			},
			wantJSON: `{
				"code": "STAFF_ALREADY_EXISTS",
				"message": "staff already exists",
				"details": []
			}`,
			wantStatus: http.StatusConflict,
		},
		{
			name: "when unexpected error when registering the staff, " +
				"then it should return a 500 with the internal error",
			jsonPayload: `{
				"staff_id": "fake-staff-id",
				"email": "test@example.com",
				"password": "ValidPassword123"
			}`,
			mocksSetup: func(service *staffmocks.MockService, _ *authmocks.MockService) {
				service.EXPECT().RegisterStaff(gomock.Any(), gomock.Any()).
					Return(staff.RegisterStaffOutput{}, errUnexpected)
			},
			wantJSON: `{
				"code": "INTERNAL_ERROR",
				"message": "an unexpected error occurred",
				"details": []
			}`,
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "when the staff is successfully registered, then it should return a 201 with the staff details",
			jsonPayload: `{
				"staff_id": "fake-staff-id",
				"email": "test@example.com",
				"password": "ValidPassword123"
			}`,
			mocksSetup: func(service *staffmocks.MockService, _ *authmocks.MockService) {
				service.EXPECT().RegisterStaff(gomock.Any(), staff.RegisterStaffInput{
					StaffID:  "fake-staff-id",
					Email:    "test@example.com",
					Password: "ValidPassword123",
				}).Return(staff.RegisterStaffOutput{
					ID:        "fake-id",
					Email:     "test@example.com",
					CreatedAt: now,
				}, nil)
			},
			wantJSON: `{
				"id":"fake-id",
				"email":"test@example.com",
				"created_at":"2025-01-01T00:00:00Z"
			}`,
			wantStatus: http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runStaffHandlerTestCase(t, logger, http.MethodPost, "/v1.0/auth/staff", tt, "")
		})
	}
}

func TestHandler_LoginStaff(t *testing.T) {
	logger := customhttp.SetupTestEnv()

	tests := []staffHandlerTestCase{
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
			name: "when there is not an active staff with the same email and password, " +
				"then it should return a 401 with invalid credentials error",
			jsonPayload: `{"email": "test@example.com", "password": "ValidPassword123"}`,
			mocksSetup: func(service *staffmocks.MockService, _ *authmocks.MockService) {
				service.EXPECT().LoginStaff(gomock.Any(), gomock.Any()).
					Return(staff.LoginStaffOutput{}, authcore.ErrInvalidCredentials)
			},
			wantJSON: `{
				"code": "INVALID_CREDENTIALS",
				"message": "invalid credentials",
				"details": []
			}`,
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "when unexpected error when login the staff, " +
				"then it should return a 500 with the internal error",
			jsonPayload: `{"email": "test@example.com", "password": "ValidPassword123"}`,
			mocksSetup: func(service *staffmocks.MockService, _ *authmocks.MockService) {
				service.EXPECT().LoginStaff(gomock.Any(), gomock.Any()).
					Return(staff.LoginStaffOutput{}, errUnexpected)
			},
			wantJSON: `{
				"code": "INTERNAL_ERROR",
				"message": "an unexpected error occurred",
				"details": []
			}`,
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "when an active staff has the same email and password, " +
				"then it should return a 200 with the token",
			jsonPayload: `{"email": "test@example.com", "password": "ValidPassword123"}`,
			mocksSetup: func(service *staffmocks.MockService, _ *authmocks.MockService) {
				service.EXPECT().LoginStaff(gomock.Any(), staff.LoginStaffInput{
					Email:    "test@example.com",
					Password: "ValidPassword123",
				}).Return(staff.LoginStaffOutput{
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
			runStaffHandlerTestCase(t, logger, http.MethodPost, "/v1.0/staff/login", tt, "")
		})
	}
}

// runStaffHandlerTestCase executes a test case for the staff handler, which is common for all tests.
func runStaffHandlerTestCase(
	t *testing.T,
	logger log.Logger,
	httpMethod string,
	route string,
	tt staffHandlerTestCase,
	token string,
) {
	// Create a new mock service
	service := staffmocks.NewMockService(gomock.NewController(t))
	authService := authmocks.NewMockService(gomock.NewController(t))
	if tt.mocksSetup != nil {
		tt.mocksSetup(service, authService)
	}

	// Initialize the authentication middleware
	authMiddleware := auth.NewMiddleware(logger, authService)

	// Initialize the handler
	h := staff.NewHandler(logger, service, authMiddleware)

	// Make HTTP request
	w := customhttp.ServeTestHTTPRequest(t, h, httpMethod, route, token, tt.queryParams, tt.jsonPayload)

	assert.Equal(t, tt.wantStatus, w.Code)
	assert.JSONEq(t, tt.wantJSON, w.Body.String())
}
