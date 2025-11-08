//go:build unit

package restaurants_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/alexgrauroca/practice-food-delivery-platform/pkg/log"
	"github.com/alexgrauroca/practice-food-delivery-platform/services/restaurant-service/internal/restaurants"
	restaurantsmocks "github.com/alexgrauroca/practice-food-delivery-platform/services/restaurant-service/internal/restaurants/mocks"
	"github.com/alexgrauroca/practice-food-delivery-platform/services/restaurant-service/internal/staff"
	staffmocks "github.com/alexgrauroca/practice-food-delivery-platform/services/restaurant-service/internal/staff/mocks"
)

var (
	errRepo  = errors.New("repository error")
	errStaff = errors.New("staff error")
)

type serviceTestCase[I, W any] struct {
	name       string
	input      I
	mocksSetup func(repo *restaurantsmocks.MockRepository, staffServ *staffmocks.MockService)
	want       W
	wantErr    error
}

func TestService_RegisterRestaurant(t *testing.T) {
	now := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	logger, _ := log.NewTest()

	tests := []serviceTestCase[restaurants.RegisterRestaurantInput, restaurants.RegisterRestaurantOutput]{
		{
			name: "when there is an active restaurant with the same vat code, " +
				"then it should return a restaurant already exists error",
			input: restaurants.RegisterRestaurantInput{
				Restaurant: restaurants.RestaurantInput{VatCode: "duplicated-vat-code"},
			},
			mocksSetup: func(repo *restaurantsmocks.MockRepository, _ *staffmocks.MockService) {
				repo.EXPECT().CreateRestaurant(gomock.Any(), gomock.Any()).
					Return(restaurants.Restaurant{}, restaurants.ErrRestaurantAlreadyExists)
			},
			want:    restaurants.RegisterRestaurantOutput{},
			wantErr: restaurants.ErrRestaurantAlreadyExists,
		},
		{
			name: "when there is an unexpected error when registering the restaurant, " +
				"then it should propagate the error",
			input: restaurants.RegisterRestaurantInput{
				Restaurant: restaurants.RestaurantInput{VatCode: "valid-vat-code"},
			},
			mocksSetup: func(repo *restaurantsmocks.MockRepository, _ *staffmocks.MockService) {
				repo.EXPECT().CreateRestaurant(gomock.Any(), gomock.Any()).
					Return(restaurants.Restaurant{}, errRepo)
			},
			want:    restaurants.RegisterRestaurantOutput{},
			wantErr: errRepo,
		},
		{
			name: "when cannot register the staff owner and there is an unexpected error purging the restaurant, " +
				"then it should propagate the staff error",
			input: restaurants.RegisterRestaurantInput{
				Restaurant: restaurants.RestaurantInput{VatCode: "valid-vat-code"},
				StaffOwner: restaurants.StaffOwnerInput{Email: "user@example.com"},
			},
			mocksSetup: func(repo *restaurantsmocks.MockRepository, staffServ *staffmocks.MockService) {
				repo.EXPECT().CreateRestaurant(gomock.Any(), gomock.Any()).
					Return(restaurants.Restaurant{ID: "fake-restaurant-id"}, nil)

				staffServ.EXPECT().RegisterStaffOwner(gomock.Any(), gomock.Any()).
					Return(staff.RegisterStaffOwnerOutput{}, errStaff)

				repo.EXPECT().PurgeRestaurant(gomock.Any(), "valid-vat-code").Return(errRepo)
			},
			want:    restaurants.RegisterRestaurantOutput{},
			wantErr: errStaff,
		},
		{
			name: "when there is an unexpected error when registering the staff owner and the restaurant is purged" +
				" successfully, then it should propagate the staff error",
			input: restaurants.RegisterRestaurantInput{
				Restaurant: restaurants.RestaurantInput{VatCode: "valid-vat-code"},
				StaffOwner: restaurants.StaffOwnerInput{Email: "user@example.com"},
			},
			mocksSetup: func(repo *restaurantsmocks.MockRepository, staffServ *staffmocks.MockService) {
				repo.EXPECT().CreateRestaurant(gomock.Any(), gomock.Any()).
					Return(restaurants.Restaurant{ID: "fake-restaurant-id"}, nil)

				staffServ.EXPECT().RegisterStaffOwner(gomock.Any(), gomock.Any()).
					Return(staff.RegisterStaffOwnerOutput{}, errStaff)

				repo.EXPECT().PurgeRestaurant(gomock.Any(), "valid-vat-code").Return(nil)
			},
			want:    restaurants.RegisterRestaurantOutput{},
			wantErr: errStaff,
		},
		{
			name: "when the restaurant and staff owner are registered successfully, " +
				"then it should return the created restaurant and staff owner",
			input: restaurants.RegisterRestaurantInput{
				Restaurant: restaurants.RestaurantInput{
					VatCode:    "valid-vat-code",
					Name:       "valid-restaurant-name",
					LegalName:  "valid-legal-name",
					TaxID:      "valid-tax-id",
					TimezoneID: "Europe/London",
					Contact: restaurants.ContactInput{
						PhonePrefix: "+44",
						PhoneNumber: "1234567890",
						Email:       "contact@example.com",
						Address:     "123 Main St",
						City:        "London",
						PostalCode:  "SW1A 1AA",
						CountryCode: "GB",
					},
				},
				StaffOwner: restaurants.StaffOwnerInput{
					Email:       "user@example.com",
					Password:    "ValidPassword123",
					Name:        "John Doe",
					Address:     "123 Main St",
					City:        "New York",
					PostalCode:  "10001",
					CountryCode: "US",
				},
			},
			mocksSetup: func(repo *restaurantsmocks.MockRepository, staffServ *staffmocks.MockService) {
				repo.EXPECT().CreateRestaurant(gomock.Any(), restaurants.CreateRestaurantParams{
					VatCode:    "valid-vat-code",
					Name:       "valid-restaurant-name",
					LegalName:  "valid-legal-name",
					TaxID:      "valid-tax-id",
					TimezoneID: "Europe/London",
					Contact: restaurants.CreateContactParams{
						PhonePrefix: "+44",
						PhoneNumber: "1234567890",
						Email:       "contact@example.com",
						Address:     "123 Main St",
						City:        "London",
						PostalCode:  "SW1A 1AA",
						CountryCode: "GB",
					},
				}).Return(restaurants.Restaurant{
					ID:         "fake-restaurant-id",
					VatCode:    "valid-vat-code",
					Name:       "valid-restaurant-name",
					LegalName:  "valid-legal-name",
					TaxID:      "valid-tax-id",
					TimezoneID: "Europe/London",
					Contact: restaurants.Contact{
						PhonePrefix: "+44",
						PhoneNumber: "1234567890",
						Email:       "contact@example.com",
						Address:     "123 Main St",
						City:        "London",
						PostalCode:  "SW1A 1AA",
						CountryCode: "GB",
					},
					CreatedAt: now,
					UpdatedAt: now,
				}, nil)

				staffServ.EXPECT().RegisterStaffOwner(gomock.Any(), staff.RegisterStaffOwnerInput{
					Email:        "user@example.com",
					Password:     "ValidPassword123",
					RestaurantID: "fake-restaurant-id",
					Name:         "John Doe",
					Address:      "123 Main St",
					City:         "New York",
					PostalCode:   "10001",
					CountryCode:  "US",
				}).Return(staff.RegisterStaffOwnerOutput{
					ID:           "fake-staff-owner-id",
					Email:        "user@example.com",
					RestaurantID: "fake-restaurant-id",
					Owner:        true,
					Name:         "John Doe",
					Address:      "123 Main St",
					City:         "New York",
					PostalCode:   "10001",
					CountryCode:  "US",
					CreatedAt:    now,
					UpdatedAt:    now,
				}, nil)
			},
			want: restaurants.RegisterRestaurantOutput{
				Restaurant: restaurants.RestaurantOutput{
					ID:         "fake-restaurant-id",
					VatCode:    "valid-vat-code",
					Name:       "valid-restaurant-name",
					LegalName:  "valid-legal-name",
					TaxID:      "valid-tax-id",
					TimezoneID: "Europe/London",
					Contact: restaurants.ContactOutput{
						PhonePrefix: "+44",
						PhoneNumber: "1234567890",
						Email:       "contact@example.com",
						Address:     "123 Main St",
						City:        "London",
						PostalCode:  "SW1A 1AA",
						CountryCode: "GB",
					},
					CreatedAt: now,
					UpdatedAt: now,
				},
				StaffOwner: restaurants.StaffOwnerOutput{
					ID:           "fake-staff-owner-id",
					Email:        "user@example.com",
					RestaurantID: "fake-restaurant-id",
					Owner:        true,
					Name:         "John Doe",
					Address:      "123 Main St",
					City:         "New York",
					PostalCode:   "10001",
					CountryCode:  "US",
					CreatedAt:    now,
					UpdatedAt:    now,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctr := gomock.NewController(t)
			defer ctr.Finish()

			repo := restaurantsmocks.NewMockRepository(ctr)
			staffServ := staffmocks.NewMockService(ctr)
			if tt.mocksSetup != nil {
				tt.mocksSetup(repo, staffServ)
			}

			service := restaurants.NewService(logger, repo, staffServ)
			got, err := service.RegisterRestaurant(context.Background(), tt.input)

			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}
