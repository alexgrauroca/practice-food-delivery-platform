//go:build unit

package customers_test

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/alexgrauroca/practice-food-delivery-platform/pkg/auth"
	authmocks "github.com/alexgrauroca/practice-food-delivery-platform/pkg/auth/mocks"
	customhttp "github.com/alexgrauroca/practice-food-delivery-platform/pkg/http"
	customersmocks "github.com/alexgrauroca/practice-food-delivery-platform/services/customer-service/internal/customers/mocks"

	"github.com/alexgrauroca/practice-food-delivery-platform/pkg/log"
	"github.com/alexgrauroca/practice-food-delivery-platform/services/customer-service/internal/customers"
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
					"name is required",
					"address is required",
					"city is required",
					"postal_code is required",
					"country_code is required",
				).Build(),
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
			wantJSON: customhttp.NewValidationErrorRespBuilder().
				WithDetails("email must be a valid email address").
				Build(),
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
			wantJSON: customhttp.NewValidationErrorRespBuilder().
				WithDetails(
					"password must be a valid password with at least 8 characters long",
					"postal_code must be at least 5 characters long",
					"country_code must be at least 2 characters long",
				).Build(),
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
			wantJSON: customhttp.NewValidationErrorRespBuilder().
				WithDetails(
					"name must not exceed 100 characters long",
					"address must not exceed 100 characters long",
					"city must not exceed 100 characters long",
					"postal_code must not exceed 32 characters long",
					"country_code must not exceed 2 characters long",
				).Build(),
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
			wantJSON:   customhttp.NewInternalErrorRespBuilder().Build(),
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
	logger := customhttp.SetupTestEnv()
	now := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []customerHandlerTestCase{
		{
			name:       "when any token is provided, then it should return a 401 with the unauthorized error",
			token:      "",
			pathParams: map[string]string{"customerID": "fakeID"},
			wantJSON:   auth.NewUnauthorizedRespBuilder().Build(),
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "when invalid token is provided, then it should return a 401 with the unauthorized error",
			token:      "invalid-token",
			pathParams: map[string]string{"customerID": "fakeID"},
			mocksSetup: func(_ *customersmocks.MockService, authService *authmocks.MockService) {
				authService.EXPECT().GetClaims(gomock.Any(), gomock.Any()).
					Return(auth.GetClaimsOutput{}, auth.ErrInvalidToken)
			},
			wantJSON:   auth.NewUnauthorizedRespBuilder().Build(),
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "when authenticated user is not a customer, then it should return a 403 with the forbidden error",
			token:      "none-customer-token",
			pathParams: map[string]string{"customerID": "fakeID"},
			mocksSetup: func(_ *customersmocks.MockService, authService *authmocks.MockService) {
				authService.EXPECT().GetClaims(gomock.Any(), gomock.Any()).
					Return(auth.GetClaimsOutput{
						Claims: &auth.Claims{
							Role: "none-customer-role",
						},
					}, nil)
			},
			wantJSON:   auth.NewForbiddenRespBuilder().Build(),
			wantStatus: http.StatusForbidden,
		},
		{
			name: "when authenticated customer is not the same as the one requested, " +
				"then it should return a 403 with the forbidden error",
			token:      "none-customer-token",
			pathParams: map[string]string{"customerID": "fakeID"},
			mocksSetup: func(service *customersmocks.MockService, authService *authmocks.MockService) {
				authService.EXPECT().GetClaims(gomock.Any(), gomock.Any()).
					Return(auth.GetClaimsOutput{
						Claims: &auth.Claims{
							Role: string(auth.RoleCustomer),
						},
					}, nil)

				service.EXPECT().GetCustomer(gomock.Any(), gomock.Any()).
					Return(customers.GetCustomerOutput{}, customers.ErrCustomerIDMismatch)
			},
			wantJSON:   auth.NewForbiddenRespBuilder().Build(),
			wantStatus: http.StatusForbidden,
		},
		{
			name:       "when the customer is not found, then it should return a 404 with the not found error",
			token:      "valid-token",
			pathParams: map[string]string{"customerID": "unexistingID"},
			mocksSetup: func(service *customersmocks.MockService, authService *authmocks.MockService) {
				authService.EXPECT().GetClaims(gomock.Any(), gomock.Any()).
					Return(auth.GetClaimsOutput{
						Claims: &auth.Claims{
							Role: string(auth.RoleCustomer),
						},
					}, nil)

				service.EXPECT().GetCustomer(gomock.Any(), gomock.Any()).
					Return(customers.GetCustomerOutput{}, customers.ErrCustomerNotFound)
			},
			wantJSON:   customhttp.NewNotFoundRespBuilder().Build(),
			wantStatus: http.StatusNotFound,
		},
		{
			name: "when unexpected error when getting the customer, " +
				"then it should return a 500 with the internal error",
			token:      "valid-token",
			pathParams: map[string]string{"customerID": "fakeID"},
			mocksSetup: func(service *customersmocks.MockService, authService *authmocks.MockService) {
				authService.EXPECT().GetClaims(gomock.Any(), gomock.Any()).
					Return(auth.GetClaimsOutput{
						Claims: &auth.Claims{
							Role: string(auth.RoleCustomer),
						},
					}, nil)

				service.EXPECT().GetCustomer(gomock.Any(), gomock.Any()).
					Return(customers.GetCustomerOutput{}, errUnexpected)
			},
			wantJSON:   customhttp.NewInternalErrorRespBuilder().Build(),
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "when a valid customerID is provided, then it should return a 200 with the customer details",
			token:      "valid-token",
			pathParams: map[string]string{"customerID": "fakeID"},
			mocksSetup: func(service *customersmocks.MockService, authService *authmocks.MockService) {
				authService.EXPECT().GetClaims(gomock.Any(), auth.GetClaimsInput{
					AccessToken: "valid-token",
				}).Return(auth.GetClaimsOutput{
					Claims: &auth.Claims{
						Role: string(auth.RoleCustomer),
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

func TestHandler_UpdateCustomer(t *testing.T) {
	logger := customhttp.SetupTestEnv()
	now := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	yesterday := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)

	tests := []customerHandlerTestCase{
		{
			name:        "when any token is provided, then it should return a 401 with the unauthorized error",
			token:       "",
			pathParams:  map[string]string{"customerID": "fakeID"},
			jsonPayload: `{}`,
			wantJSON:    auth.NewUnauthorizedRespBuilder().Build(),
			wantStatus:  http.StatusUnauthorized,
		},
		{
			name:        "when invalid token is provided, then it should return a 401 with the unauthorized error",
			token:       "invalid-token",
			pathParams:  map[string]string{"customerID": "fakeID"},
			jsonPayload: `{}`,
			mocksSetup: func(_ *customersmocks.MockService, authService *authmocks.MockService) {
				authService.EXPECT().GetClaims(gomock.Any(), gomock.Any()).
					Return(auth.GetClaimsOutput{}, auth.ErrInvalidToken)
			},
			wantJSON:   auth.NewUnauthorizedRespBuilder().Build(),
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:        "when authenticated user is not a customer, then it should return a 403 with the forbidden error",
			token:       "none-customer-token",
			pathParams:  map[string]string{"customerID": "fakeID"},
			jsonPayload: `{}`,
			mocksSetup: func(_ *customersmocks.MockService, authService *authmocks.MockService) {
				authService.EXPECT().GetClaims(gomock.Any(), gomock.Any()).
					Return(auth.GetClaimsOutput{
						Claims: &auth.Claims{
							Role: "none-customer-role",
						},
					}, nil)
			},
			wantJSON:   auth.NewForbiddenRespBuilder().Build(),
			wantStatus: http.StatusForbidden,
		},
		{
			name:        "when invalid payload is provided, then it should return a 400 with invalid request error",
			token:       "valid-token",
			pathParams:  map[string]string{"customerID": "fakeID"},
			jsonPayload: `{"name": 1.2, "address": true}`,
			mocksSetup: func(_ *customersmocks.MockService, authService *authmocks.MockService) {
				authService.EXPECT().GetClaims(gomock.Any(), gomock.Any()).
					Return(auth.GetClaimsOutput{
						Claims: &auth.Claims{
							Role: string(auth.RoleCustomer),
						},
					}, nil)
			},
			wantJSON:   customhttp.NewInvalidRequestRespBuilder().Build(),
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "when required fields are not provided payload, " +
				"then it should return a 400 with the required validation errors",
			token:       "valid-token",
			pathParams:  map[string]string{"customerID": "fakeID"},
			jsonPayload: `{}`,
			mocksSetup: func(_ *customersmocks.MockService, authService *authmocks.MockService) {
				authService.EXPECT().GetClaims(gomock.Any(), gomock.Any()).
					Return(auth.GetClaimsOutput{
						Claims: &auth.Claims{
							Role: string(auth.RoleCustomer),
						},
					}, nil)
			},
			wantJSON: customhttp.NewValidationErrorRespBuilder().
				WithDetails(
					"name is required",
					"address is required",
					"city is required",
					"postal_code is required",
					"country_code is required",
				).Build(),
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "when fields length are shorten than minimum required, " +
				"then it should return a 400 with the short length validation errors",
			token:      "valid-token",
			pathParams: map[string]string{"customerID": "fakeID"},
			jsonPayload: `{
				"name": "New John Doe",
				"address": "New 123 Main St",
				"city": "Los Angeles",
				"postal_code": "1",
				"country_code": "U"
			}`,
			mocksSetup: func(_ *customersmocks.MockService, authService *authmocks.MockService) {
				authService.EXPECT().GetClaims(gomock.Any(), gomock.Any()).
					Return(auth.GetClaimsOutput{
						Claims: &auth.Claims{
							Role: string(auth.RoleCustomer),
						},
					}, nil)
			},
			wantJSON: customhttp.NewValidationErrorRespBuilder().
				WithDetails(
					"postal_code must be at least 5 characters long",
					"country_code must be at least 2 characters long",
				).Build(),
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "when fields length are longer than maximum required, " +
				"then it should return a 400 with the long length validation errors",
			token:      "valid-token",
			pathParams: map[string]string{"customerID": "fakeID"},
			jsonPayload: fmt.Sprintf(`{
					"name": "%s",
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
			mocksSetup: func(_ *customersmocks.MockService, authService *authmocks.MockService) {
				authService.EXPECT().GetClaims(gomock.Any(), gomock.Any()).
					Return(auth.GetClaimsOutput{
						Claims: &auth.Claims{
							Role: string(auth.RoleCustomer),
						},
					}, nil)
			},
			wantJSON: customhttp.NewValidationErrorRespBuilder().
				WithDetails(
					"name must not exceed 100 characters long",
					"address must not exceed 100 characters long",
					"city must not exceed 100 characters long",
					"postal_code must not exceed 32 characters long",
					"country_code must not exceed 2 characters long",
				).Build(),
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "when authenticated customer is not the same as the one requested, " +
				"then it should return a 403 with the forbidden error",
			token:      "none-customer-token",
			pathParams: map[string]string{"customerID": "fakeID"},
			jsonPayload: `{
				"name": "New John Doe",
				"address": "New 123 Main St",
				"city": "Los Angeles",
				"postal_code": "09001",
				"country_code": "SP"
			}`,
			mocksSetup: func(service *customersmocks.MockService, authService *authmocks.MockService) {
				authService.EXPECT().GetClaims(gomock.Any(), gomock.Any()).
					Return(auth.GetClaimsOutput{
						Claims: &auth.Claims{
							Role: string(auth.RoleCustomer),
						},
					}, nil)

				service.EXPECT().UpdateCustomer(gomock.Any(), gomock.Any()).
					Return(customers.UpdateCustomerOutput{}, customers.ErrCustomerIDMismatch)
			},
			wantJSON:   auth.NewForbiddenRespBuilder().Build(),
			wantStatus: http.StatusForbidden,
		},
		{
			name:  "when the customer is not found, then it should return a 404 with the not found error",
			token: "valid-token",
			pathParams: map[string]string{
				"customerID": "unexistingID",
			},
			jsonPayload: `{
				"name": "New John Doe",
				"address": "New 123 Main St",
				"city": "Los Angeles",
				"postal_code": "09001",
				"country_code": "SP"
			}`,
			mocksSetup: func(service *customersmocks.MockService, authService *authmocks.MockService) {
				authService.EXPECT().GetClaims(gomock.Any(), gomock.Any()).
					Return(auth.GetClaimsOutput{
						Claims: &auth.Claims{
							Role: string(auth.RoleCustomer),
						},
					}, nil)

				service.EXPECT().UpdateCustomer(gomock.Any(), gomock.Any()).
					Return(customers.UpdateCustomerOutput{}, customers.ErrCustomerNotFound)
			},
			wantJSON:   customhttp.NewNotFoundRespBuilder().Build(),
			wantStatus: http.StatusNotFound,
		},
		{
			name: "when unexpected error when updating the customer, " +
				"then it should return a 500 with the internal error",
			token:      "valid-token",
			pathParams: map[string]string{"customerID": "fakeID"},
			jsonPayload: `{
				"name": "New John Doe",
				"address": "New 123 Main St",
				"city": "Los Angeles",
				"postal_code": "09001",
				"country_code": "SP"
			}`,
			mocksSetup: func(service *customersmocks.MockService, authService *authmocks.MockService) {
				authService.EXPECT().GetClaims(gomock.Any(), gomock.Any()).
					Return(auth.GetClaimsOutput{
						Claims: &auth.Claims{
							Role: string(auth.RoleCustomer),
						},
					}, nil)

				service.EXPECT().UpdateCustomer(gomock.Any(), gomock.Any()).
					Return(customers.UpdateCustomerOutput{}, errUnexpected)
			},
			wantJSON:   customhttp.NewInternalErrorRespBuilder().Build(),
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "when the customer can be updated, then it should return a 200 with the customer details updated",
			token:      "valid-token",
			pathParams: map[string]string{"customerID": "fakeID"},
			jsonPayload: `{
				"name": "New John Doe",
				"address": "New 123 Main St",
				"city": "Los Angeles",
				"postal_code": "09001",
				"country_code": "SP"
			}`,
			mocksSetup: func(service *customersmocks.MockService, authService *authmocks.MockService) {
				authService.EXPECT().GetClaims(gomock.Any(), auth.GetClaimsInput{
					AccessToken: "valid-token",
				}).Return(auth.GetClaimsOutput{
					Claims: &auth.Claims{
						Role: string(auth.RoleCustomer),
					},
				}, nil)

				service.EXPECT().UpdateCustomer(gomock.Any(), customers.UpdateCustomerInput{
					CustomerID:  "fakeID",
					Name:        "New John Doe",
					Address:     "New 123 Main St",
					City:        "Los Angeles",
					PostalCode:  "09001",
					CountryCode: "SP",
				}).Return(customers.UpdateCustomerOutput{
					ID:          "fakeID",
					Name:        "New John Doe",
					Email:       "test@example.com",
					Address:     "New 123 Main St",
					City:        "Los Angeles",
					PostalCode:  "09001",
					CountryCode: "SP",
					CreatedAt:   yesterday,
					UpdatedAt:   now,
				}, nil)
			},
			wantJSON: `{
				"id": "fakeID",
				"name": "New John Doe",
			    "email": "test@example.com",
			    "address": "New 123 Main St",
			    "city": "Los Angeles",
			    "postal_code": "09001",
			    "country_code": "SP",
				"created_at": "2024-12-31T00:00:00Z",
				"updated_at": "2025-01-01T00:00:00Z"
			}`,
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updateCustomerPath := fmt.Sprintf("/v1.0/customers/%s", tt.pathParams["customerID"])
			runCustomerHandlerTestCase(t, logger, http.MethodPut, updateCustomerPath, tt, tt.token)
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
