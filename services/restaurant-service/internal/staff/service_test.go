//go:build unit

package staff_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/alexgrauroca/practice-food-delivery-platform/pkg/clients/authentication"
	authclimocks "github.com/alexgrauroca/practice-food-delivery-platform/pkg/clients/authentication/mocks"
	"github.com/alexgrauroca/practice-food-delivery-platform/pkg/log"
	"github.com/alexgrauroca/practice-food-delivery-platform/services/restaurant-service/internal/staff"
	staffmocks "github.com/alexgrauroca/practice-food-delivery-platform/services/restaurant-service/internal/staff/mocks"
)

var (
	errRepo        = errors.New("repository error")
	errAuthService = errors.New("auth service error")
)

type serviceCase[I, W any] struct {
	name       string
	input      I
	mocksSetup func(repo *staffmocks.MockRepository, authcli *authclimocks.MockClient)
	want       W
	wantErr    error
}

func TestService_RegisterStaffOwner(t *testing.T) {
	now := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	logger, _ := log.NewTest()

	tests := []serviceCase[staff.RegisterStaffOwnerInput, staff.RegisterStaffOwnerOutput]{
		{
			name:  "when the staff owner cannot be created, then it should propagate the error",
			input: staff.RegisterStaffOwnerInput{Email: "test@example.com"},
			mocksSetup: func(repo *staffmocks.MockRepository, _ *authclimocks.MockClient) {
				repo.EXPECT().CreateStaff(gomock.Any(), gomock.Any()).
					Return(staff.Staff{}, errRepo)
			},
			want:    staff.RegisterStaffOwnerOutput{},
			wantErr: errRepo,
		},
		{
			name: "when the auth service returns an unexpected error and the staff cannot be purged, " +
				"then it returns the auth service error",
			input: staff.RegisterStaffOwnerInput{Email: "test@example.com"},
			mocksSetup: func(repo *staffmocks.MockRepository, authcli *authclimocks.MockClient) {
				repo.EXPECT().CreateStaff(gomock.Any(), gomock.Any()).
					Return(staff.Staff{ID: "fake-staff-id"}, nil)

				authcli.EXPECT().RegisterStaff(gomock.Any(), gomock.Any()).
					Return(authentication.RegisterStaffResponse{}, errAuthService)

				repo.EXPECT().PurgeStaff(gomock.Any(), gomock.Any()).Return(errRepo)
			},
			want:    staff.RegisterStaffOwnerOutput{},
			wantErr: errAuthService,
		},
		{
			name: "when the auth service returns an unexpected error and the staff is purged successfully, " +
				"then it returns the auth service error",
			input: staff.RegisterStaffOwnerInput{Email: "test@example.com"},
			mocksSetup: func(repo *staffmocks.MockRepository, authcli *authclimocks.MockClient) {
				repo.EXPECT().CreateStaff(gomock.Any(), gomock.Any()).
					Return(staff.Staff{ID: "fake-staff-id"}, nil)

				authcli.EXPECT().RegisterStaff(gomock.Any(), gomock.Any()).
					Return(authentication.RegisterStaffResponse{}, errAuthService)

				repo.EXPECT().PurgeStaff(gomock.Any(), "test@example.com").Return(nil)
			},
			want:    staff.RegisterStaffOwnerOutput{},
			wantErr: errAuthService,
		},
		{
			name: "when the staff owner is registered successfully, then it should return the created staff owner",
			input: staff.RegisterStaffOwnerInput{
				Email:        "test@example.com",
				Password:     "ValidPassword123",
				RestaurantID: "valid-restaurant-id",
				Name:         "Test Staff Owner",
				Address:      "123 Main St",
				City:         "London",
				PostalCode:   "SW1A 1AA",
				CountryCode:  "GB",
			},
			mocksSetup: func(repo *staffmocks.MockRepository, authcli *authclimocks.MockClient) {
				repo.EXPECT().CreateStaff(gomock.Any(), staff.CreateStaffParams{
					Email:        "test@example.com",
					RestaurantID: "valid-restaurant-id",
					Owner:        true,
					Name:         "Test Staff Owner",
					Address:      "123 Main St",
					City:         "London",
					PostalCode:   "SW1A 1AA",
					CountryCode:  "GB",
				}).Return(staff.Staff{
					ID:           "fake-staff-id",
					Email:        "test@example.com",
					RestaurantID: "valid-restaurant-id",
					Owner:        true,
					Name:         "Test Staff Owner",
					Address:      "123 Main St",
					City:         "London",
					PostalCode:   "SW1A 1AA",
					CountryCode:  "GB",
					CreatedAt:    now,
					UpdatedAt:    now,
				}, nil)

				authcli.EXPECT().RegisterStaff(gomock.Any(), authentication.RegisterStaffRequest{
					StaffID:      "fake-staff-id",
					Email:        "test@example.com",
					RestaurantID: "valid-restaurant-id",
					Password:     "ValidPassword123",
				}).Return(authentication.RegisterStaffResponse{
					ID:           "fake-auth-staff-id",
					StaffID:      "fake-staff-id",
					Email:        "test@example.com",
					RestaurantID: "valid-restaurant-id",
					CreatedAt:    now,
					UpdatedAt:    now,
				}, nil)
			},
			want: staff.RegisterStaffOwnerOutput{
				ID:           "fake-staff-id",
				Email:        "test@example.com",
				RestaurantID: "valid-restaurant-id",
				Owner:        true,
				Name:         "Test Staff Owner",
				Address:      "123 Main St",
				City:         "London",
				PostalCode:   "SW1A 1AA",
				CountryCode:  "GB",
				CreatedAt:    now,
				UpdatedAt:    now,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctr := gomock.NewController(t)
			defer ctr.Finish()

			repo := staffmocks.NewMockRepository(ctr)
			authcli := authclimocks.NewMockClient(ctr)
			if tt.mocksSetup != nil {
				tt.mocksSetup(repo, authcli)
			}

			service := staff.NewService(logger, repo, authcli)
			got, err := service.RegisterStaffOwner(context.Background(), tt.input)

			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}
