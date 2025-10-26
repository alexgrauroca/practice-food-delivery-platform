//go:build unit

package staff_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	authmocks "github.com/alexgrauroca/practice-food-delivery-platform/pkg/auth/mocks"
	"github.com/alexgrauroca/practice-food-delivery-platform/pkg/log"
	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/password"
	refreshmocks "github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/refresh/mocks"
	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/staff"
	staffmocks "github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/staff/mocks"
)

var (
	errRepo  = errors.New("repository error")
	errToken = errors.New("token error")
)

type staffServiceTestCase[I, W any] struct {
	name       string
	input      I
	want       W
	mocksSetup func(
		repo *staffmocks.MockRepository,
		refreshService *refreshmocks.MockService,
		authService *authmocks.MockService,
		authctx *authmocks.MockContextReader,
	)
	wantErr error
}

func TestService_RegisterStaff(t *testing.T) {
	now := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	logger, _ := log.NewTest()

	tests := []staffServiceTestCase[staff.RegisterStaffInput, staff.RegisterStaffOutput]{
		{
			name: "when there is an active staff with the same email, " +
				"then it should return a staff already exists error",
			input: staff.RegisterStaffInput{
				StaffID:  "fake-staff-id",
				Email:    "test@example.com",
				Password: "ValidPassword123",
			},
			mocksSetup: func(
				repo *staffmocks.MockRepository,
				_ *refreshmocks.MockService,
				_ *authmocks.MockService,
				_ *authmocks.MockContextReader,
			) {
				repo.EXPECT().CreateStaff(gomock.Any(), gomock.Any()).
					Return(staff.Staff{}, staff.ErrStaffAlreadyExists)
			},
			want:    staff.RegisterStaffOutput{},
			wantErr: staff.ErrStaffAlreadyExists,
		},
		{
			name: "when there is an unexpected error when creating the staff, then it should propagate the error",
			input: staff.RegisterStaffInput{
				StaffID:  "fake-staff-id",
				Email:    "test@example.com",
				Password: "ValidPassword123",
			},
			mocksSetup: func(
				repo *staffmocks.MockRepository,
				_ *refreshmocks.MockService,
				_ *authmocks.MockService,
				_ *authmocks.MockContextReader,
			) {
				repo.EXPECT().CreateStaff(gomock.Any(), gomock.Any()).
					Return(staff.Staff{}, errRepo)
			},
			want:    staff.RegisterStaffOutput{},
			wantErr: errRepo,
		},
		{
			name: "when the staff can be created, then it should return the created staff",
			input: staff.RegisterStaffInput{
				StaffID:  "fake-staff-id",
				Email:    "test@example.com",
				Password: "ValidPassword123",
			},
			mocksSetup: func(
				repo *staffmocks.MockRepository,
				_ *refreshmocks.MockService,
				_ *authmocks.MockService,
				_ *authmocks.MockContextReader,
			) {
				repo.EXPECT().CreateStaff(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, params staff.CreateStaffParams) (staff.Staff, error) {
						// Assert that the password is hashed
						ok := password.Verify(params.Password, "ValidPassword123")
						require.True(t, ok, "Password should be hashed and match the input password")

						return staff.Staff{
							ID:        "fake-id",
							StaffID:   params.StaffID,
							Email:     params.Email,
							Password:  params.Password,
							CreatedAt: now,
							UpdatedAt: now,
							Active:    true,
						}, nil
					})
			},
			want: staff.RegisterStaffOutput{
				ID:        "fake-id",
				Email:     "test@example.com",
				CreatedAt: now,
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := staffmocks.NewMockRepository(ctrl)
			refreshService := refreshmocks.NewMockService(ctrl)
			authService := authmocks.NewMockService(ctrl)
			authctx := authmocks.NewMockContextReader(ctrl)

			if tt.mocksSetup != nil {
				tt.mocksSetup(repo, refreshService, authService, authctx)
			}

			service := staff.NewService(logger, repo, refreshService, authService, authctx)
			got, err := service.RegisterStaff(context.Background(), tt.input)

			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}
