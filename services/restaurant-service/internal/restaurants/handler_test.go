package restaurants_test

import (
	"errors"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/alexgrauroca/practice-food-delivery-platform/services/restaurant-service/internal/restaurants/testbuilder"

	customhttp "github.com/alexgrauroca/practice-food-delivery-platform/pkg/http"
	"github.com/alexgrauroca/practice-food-delivery-platform/pkg/log"
	"github.com/alexgrauroca/practice-food-delivery-platform/services/restaurant-service/internal/restaurants"
	restaurantsmocks "github.com/alexgrauroca/practice-food-delivery-platform/services/restaurant-service/internal/restaurants/mocks"
)

var errUnexpected = errors.New("unexpected error")

type handlerTestCase struct {
	name        string
	token       string
	pathParams  map[string]string
	queryParams map[string]string
	jsonPayload string
	mocksSetup  func(service *restaurantsmocks.MockService)
	wantJSON    string
	wantStatus  int
}

func TestHandler_RegisterRestaurant(t *testing.T) {
	logger := customhttp.SetupTestEnv()
	now := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []handlerTestCase{
		{
			name:        "when invalid payload is provided, then it should return a 400 with invalid request error",
			jsonPayload: `{"restaurant": 1.2, "staff_owner": true}`,
			wantJSON:    customhttp.NewInvalidRequestRespBuilder().Build(),
			wantStatus:  http.StatusBadRequest,
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
			jsonPayload: testbuilder.NewValidRegisterRestaurantPayload().
				WithContactEmail("invalid-email-1").
				WithOwnerEmail("invalid-email-2").
				Build(),
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
			jsonPayload: testbuilder.NewValidRegisterRestaurantPayload().
				WithRestaurantTimezone("Invalid/Timezone").
				Build(),
			wantJSON: `{
				"code": "VALIDATION_ERROR",
				"message": "validation failed",
				"details": [ "restaurant.timezone_id is invalid" ]
			}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "when invalid phone data is provided, then it should return a 400 with the phone validation error",
			jsonPayload: testbuilder.NewValidRegisterRestaurantPayload().
				WithContactPhone("a", "bbbbbbb").
				Build(),
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
			jsonPayload: testbuilder.NewValidRegisterRestaurantPayload().
				WithPostalCode("1").
				WithCountryCode("U").
				WithPassword("short").
				Build(),
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
		{
			name: "when fields length are longer than maximum required, " +
				"then it should return a 400 with the long length validation errors",
			jsonPayload: testbuilder.NewValidRegisterRestaurantPayload().
				WithVatCode(strings.Repeat("a", 41)).
				WithRestaurantName(strings.Repeat("a", 101)).
				WithLegalName(strings.Repeat("a", 101)).
				WithTaxID(strings.Repeat("a", 41)).
				WithOwnerName(strings.Repeat("a", 101)).
				WithAddress(strings.Repeat("a", 101)).
				WithCity(strings.Repeat("a", 101)).
				WithPostalCode(strings.Repeat("1", 33)).
				WithCountryCode("USA").
				Build(),
			wantJSON: `{
				"code": "VALIDATION_ERROR",
				"message": "validation failed",
				"details": [
					"restaurant.vat_code must not exceed 40 characters long",
					"restaurant.name must not exceed 100 characters long",
					"restaurant.legal_name must not exceed 100 characters long",
					"restaurant.tax_id must not exceed 40 characters long",
					"restaurant.contact.address must not exceed 100 characters long",
					"restaurant.contact.city must not exceed 100 characters long",
					"restaurant.contact.postal_code must not exceed 32 characters long",
					"restaurant.contact.country_code must not exceed 2 characters long",
					"staff_owner.name must not exceed 100 characters long",
					"staff_owner.address must not exceed 100 characters long",
					"staff_owner.city must not exceed 100 characters long",
					"staff_owner.postal_code must not exceed 32 characters long",
					"staff_owner.country_code must not exceed 2 characters long"
				]
			}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "when the restaurant already exists, " +
				"then it should return a 409 with the restaurant already exists error",
			jsonPayload: testbuilder.NewValidRegisterRestaurantPayload().Build(),
			mocksSetup: func(service *restaurantsmocks.MockService) {
				service.EXPECT().RegisterRestaurant(gomock.Any(), gomock.Any()).
					Return(restaurants.RegisterRestaurantOutput{}, restaurants.ErrRestaurantAlreadyExists)
			},
			wantJSON: `{
				"code": "RESTAURANT_ALREADY_EXISTS",
				"message": "restaurant already exists",
				"details": []
			}`,
			wantStatus: http.StatusConflict,
		},
		{
			name: "when unexpected error when registering the restaurant, " +
				"then it should return a 500 with the internal error",
			jsonPayload: testbuilder.NewValidRegisterRestaurantPayload().Build(),
			mocksSetup: func(service *restaurantsmocks.MockService) {
				service.EXPECT().RegisterRestaurant(gomock.Any(), gomock.Any()).
					Return(restaurants.RegisterRestaurantOutput{}, errUnexpected)
			},
			wantJSON:   customhttp.NewInternalErrorRespBuilder().Build(),
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "when the restaurant is registered successfully, " +
				"then it should return a 201 with the restaurant and staff owner details",
			jsonPayload: testbuilder.NewValidRegisterRestaurantPayload().Build(),
			mocksSetup: func(service *restaurantsmocks.MockService) {
				service.EXPECT().RegisterRestaurant(gomock.Any(), restaurants.RegisterRestaurantInput{
					Restaurant: restaurants.RestaurantInput{
						VatCode:    "GB123456789",
						Name:       "Acme Pizza",
						LegalName:  "Acme Pizza LLC",
						TaxID:      "99-1234567",
						TimezoneID: "America/New_York",
						Contact: restaurants.ContactInput{
							PhonePrefix: "+1",
							PhoneNumber: "1234567890",
							Email:       "restaurant@example.com",
							Address:     "123 Main St",
							City:        "New York",
							PostalCode:  "10001",
							CountryCode: "US",
						},
					},
					StaffOwner: restaurants.StaffOwnerInput{
						Email:       "user@example.com",
						Password:    "strongpassword123",
						Name:        "John Doe",
						Address:     "123 Main St",
						City:        "New York",
						PostalCode:  "10001",
						CountryCode: "US",
					},
				}).Return(restaurants.RegisterRestaurantOutput{
					Restaurant: restaurants.RestaurantOutput{
						ID:         "fake-restaurant-id",
						VatCode:    "GB123456789",
						Name:       "Acme Pizza",
						LegalName:  "Acme Pizza LLC",
						TaxID:      "99-1234567",
						TimezoneID: "America/New_York",
						Contact: restaurants.ContactOutput{
							PhonePrefix: "+1",
							PhoneNumber: "1234567890",
							Email:       "restaurant@example.com",
							Address:     "123 Main St",
							City:        "New York",
							PostalCode:  "10001",
							CountryCode: "US",
						},
						CreatedAt: now,
						UpdatedAt: now,
					},
					StaffOwner: restaurants.StaffOwnerOutput{
						ID:          "fake-owner-id",
						Owner:       true,
						Email:       "user@example.com",
						Name:        "John Doe",
						Address:     "123 Main St",
						City:        "New York",
						PostalCode:  "10001",
						CountryCode: "US",
						CreatedAt:   now,
						UpdatedAt:   now,
					},
				}, nil)
			},
			wantJSON:   testbuilder.NewRegisterRestaurantSuccessResponse().Build(),
			wantStatus: http.StatusCreated,
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
	service := restaurantsmocks.NewMockService(gomock.NewController(t))
	if tt.mocksSetup != nil {
		tt.mocksSetup(service)
	}

	h := restaurants.NewHandler(logger, service)
	w := customhttp.ServeTestHTTPRequest(t, h, httpMethod, route, token, tt.queryParams, tt.jsonPayload)

	assert.Equal(t, tt.wantStatus, w.Code)
	assert.JSONEq(t, tt.wantJSON, w.Body.String())
}
