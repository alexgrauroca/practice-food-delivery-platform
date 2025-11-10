//go:build unit

package staff_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/alexgrauroca/practice-food-delivery-platform/pkg/log"
	"github.com/alexgrauroca/practice-food-delivery-platform/services/restaurant-service/internal/staff"
)

var (
	errRepo        = errors.New("repository error")
	errAuthService = errors.New("auth service error")
)

type serviceCase[I, W any] struct {
	name       string
	input      I
	mocksSetup func(repo, auth)
	want       W
	wantErr    error
}

func TestService_RegisterStaffOwner(t *testing.T) {
	now := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	logger, _ := log.NewTest()

	tests := []serviceCase[staff.RegisterStaffOwnerInput, staff.RegisterStaffOwnerOutput]{
		{
			name:    "when the staff owner cannot be created, then it should propagate the error",
			input:   staff.RegisterStaffOwnerInput{Email: "test@example.com"},
			want:    staff.RegisterStaffOwnerOutput{},
			wantErr: errRepo,
		},
		{
			name: "when the auth service returns an unexpected error and the staff cannot be purged, " +
				"then it returns the auth service error",
			input:   staff.RegisterStaffOwnerInput{Email: "test@example.com"},
			want:    staff.RegisterStaffOwnerOutput{},
			wantErr: errAuthService,
		},
		{
			name: "when the auth service returns an unexpected error and the staff is purged successfully, " +
				"then it returns the auth service error",
			input:   staff.RegisterStaffOwnerInput{Email: "test@example.com"},
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
			want: staff.RegisterStaffOwnerOutput{
				ID:           "valid-staff-owner-id",
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

			if tt.mocksSetup != nil {
				tt.mocksSetup()
			}

			service := staff.NewService(logger)
			got, err := service.RegisterStaffOwner(context.Background(), tt.input)

			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}
