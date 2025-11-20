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

	"github.com/alexgrauroca/practice-food-delivery-platform/pkg/log"
	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/authcore"
	authcoremocks "github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/authcore/mocks"
	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/password"
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
		authCoreService *authcoremocks.MockService,
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
				StaffID:      "fake-staff-id",
				Email:        "test@example.com",
				RestaurantID: "fake-restaurant-id",
				Password:     "ValidPassword123",
			},
			mocksSetup: func(
				repo *staffmocks.MockRepository,
				_ *authcoremocks.MockService,
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
				StaffID:      "fake-staff-id",
				Email:        "test@example.com",
				RestaurantID: "fake-restaurant-id",
				Password:     "ValidPassword123",
			},
			mocksSetup: func(
				repo *staffmocks.MockRepository,
				_ *authcoremocks.MockService,
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
				StaffID:      "fake-staff-id",
				Email:        "test@example.com",
				RestaurantID: "fake-restaurant-id",
				Password:     "ValidPassword123",
			},
			mocksSetup: func(
				repo *staffmocks.MockRepository,
				_ *authcoremocks.MockService,
			) {
				repo.EXPECT().CreateStaff(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, params staff.CreateStaffParams) (staff.Staff, error) {
						// Assert that the password is hashed
						ok := password.Verify(params.Password, "ValidPassword123")
						require.True(t, ok, "Password should be hashed and match the input password")

						return staff.Staff{
							ID:           "fake-id",
							StaffID:      params.StaffID,
							Email:        params.Email,
							Password:     params.Password,
							RestaurantID: params.RestaurantID,
							CreatedAt:    now,
							UpdatedAt:    now,
							Active:       true,
						}, nil
					})
			},
			want: staff.RegisterStaffOutput{
				ID:           "fake-id",
				Email:        "test@example.com",
				RestaurantID: "fake-restaurant-id",
				CreatedAt:    now,
				UpdatedAt:    now,
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, cleanup := serviceSetup(t, logger, tt.mocksSetup)
			defer cleanup()

			got, err := service.RegisterStaff(context.Background(), tt.input)

			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestService_LoginStaff(t *testing.T) {
	now := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	logger, _ := log.NewTest()

	tests := []staffServiceTestCase[staff.LoginStaffInput, staff.LoginStaffOutput]{
		{
			name: "when there is not an active staff with the same email, " +
				"then it should return an invalid credentials error",
			input: staff.LoginStaffInput{
				Email:    "test@example.com",
				Password: "ValidPassword123",
			},
			mocksSetup: func(
				repo *staffmocks.MockRepository,
				_ *authcoremocks.MockService,
			) {
				repo.EXPECT().FindByEmail(gomock.Any(), gomock.Any()).
					Return(staff.Staff{}, staff.ErrStaffNotFound)
			},
			want:    staff.LoginStaffOutput{},
			wantErr: authcore.ErrInvalidCredentials,
		},
		{
			name: "when there is not an active staff with the same password, " +
				"then it should return an invalid credentials error",
			input: staff.LoginStaffInput{
				Email:    "test@example.com",
				Password: "InvalidPassword123",
			},
			mocksSetup: func(
				repo *staffmocks.MockRepository,
				_ *authcoremocks.MockService,
			) {
				hashedPassword, err := password.Hash("ValidPassword123")
				require.NoError(t, err)

				repo.EXPECT().FindByEmail(gomock.Any(), gomock.Any()).
					Return(staff.Staff{
						ID:        "fake-id",
						Email:     "test@example.com",
						Password:  hashedPassword, // This should be a hashed password
						CreatedAt: now,
						UpdatedAt: now,
						Active:    true,
					}, nil)
			},
			want:    staff.LoginStaffOutput{},
			wantErr: authcore.ErrInvalidCredentials,
		},
		{
			name: "when there is an unexpected error when fetching the staff, then it should propagate the error",
			input: staff.LoginStaffInput{
				Email:    "test@example.com",
				Password: "ValidPassword123",
			},
			mocksSetup: func(
				repo *staffmocks.MockRepository,
				_ *authcoremocks.MockService,
			) {
				repo.EXPECT().FindByEmail(gomock.Any(), gomock.Any()).
					Return(staff.Staff{}, errRepo)
			},
			want:    staff.LoginStaffOutput{},
			wantErr: errRepo,
		},
		{
			name: "when there is an error generating the token pair, then it should propagate the error",
			input: staff.LoginStaffInput{
				Email:    "test@example.com",
				Password: "ValidPassword123",
			},
			mocksSetup: func(
				repo *staffmocks.MockRepository,
				authCoreService *authcoremocks.MockService,
			) {
				hashedPassword, err := password.Hash("ValidPassword123")
				require.NoError(t, err)

				repo.EXPECT().FindByEmail(gomock.Any(), "test@example.com").
					Return(staff.Staff{
						ID:        "fake-id",
						Email:     "test@example.com",
						Password:  hashedPassword, // This should be a hashed password
						CreatedAt: now,
						UpdatedAt: now,
						Active:    true,
					}, nil)

				authCoreService.EXPECT().GenerateTokenPair(gomock.Any(), gomock.Any()).
					Return(authcore.TokenPair{}, errToken)
			},
			want:    staff.LoginStaffOutput{},
			wantErr: errToken,
		},
		{
			name: "when there is an active staff with the same email and password, then it should return its token",
			input: staff.LoginStaffInput{
				Email:    "test@example.com",
				Password: "ValidPassword123",
			},
			mocksSetup: func(
				repo *staffmocks.MockRepository,
				authCoreService *authcoremocks.MockService,
			) {
				hashedPassword, err := password.Hash("ValidPassword123")
				require.NoError(t, err)

				repo.EXPECT().FindByEmail(gomock.Any(), "test@example.com").
					Return(staff.Staff{
						ID:        "fake-staff-id",
						StaffID:   "fake-id",
						Email:     "test@example.com",
						Password:  hashedPassword, // This should be a hashed password
						CreatedAt: now,
						UpdatedAt: now,
						Active:    true,
					}, nil)

				authCoreService.EXPECT().GenerateTokenPair(gomock.Any(), gomock.Any()).
					Return(authcore.TokenPair{
						AccessToken:  "fake-token",
						RefreshToken: "fake-refresh-token",
						ExpiresIn:    3600,
						TokenType:    "Bearer",
					}, nil)
			},
			want: staff.LoginStaffOutput{
				TokenPair: authcore.TokenPair{
					AccessToken:  "fake-token",
					ExpiresIn:    3600, // 1 hour
					TokenType:    "Bearer",
					RefreshToken: "fake-refresh-token",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, cleanup := serviceSetup(t, logger, tt.mocksSetup)
			defer cleanup()

			got, err := service.LoginStaff(context.Background(), tt.input)

			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestService_RefreshStaff(t *testing.T) {
	logger, _ := log.NewTest()

	tests := []staffServiceTestCase[staff.RefreshStaffInput, staff.RefreshStaffOutput]{
		{
			name: "when there is an error refreshing the token, then it should propagate the error",
			input: staff.RefreshStaffInput{
				RefreshToken: "InvalidRefreshToken",
				AccessToken:  "ValidAccessToken",
			},
			mocksSetup: func(
				_ *staffmocks.MockRepository,
				authCoreService *authcoremocks.MockService,
			) {
				authCoreService.EXPECT().RefreshToken(gomock.Any(), gomock.Any()).
					Return(authcore.TokenPair{}, errToken)
			},
			want:    staff.RefreshStaffOutput{},
			wantErr: errToken,
		},
		{
			name: "when the new access token is generated correctly, then it should return the new token",
			input: staff.RefreshStaffInput{
				RefreshToken: "ValidRefreshToken",
				AccessToken:  "ValidAccessToken",
			},
			mocksSetup: func(
				_ *staffmocks.MockRepository,
				authCoreService *authcoremocks.MockService,
			) {
				authCoreService.EXPECT().RefreshToken(gomock.Any(), gomock.Any()).
					Return(authcore.TokenPair{
						AccessToken:  "fake-token",
						RefreshToken: "fake-refresh-token",
						ExpiresIn:    3600,
						TokenType:    "Bearer",
					}, nil)
			},
			want: staff.RefreshStaffOutput{
				TokenPair: authcore.TokenPair{
					AccessToken:  "fake-token",
					RefreshToken: "fake-refresh-token",
					ExpiresIn:    3600,
					TokenType:    "Bearer",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, cleanup := serviceSetup(t, logger, tt.mocksSetup)
			defer cleanup()

			got, err := service.RefreshStaff(context.Background(), tt.input)

			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func serviceSetup(t *testing.T, logger log.Logger, mocksSetup func(
	repo *staffmocks.MockRepository,
	authCoreService *authcoremocks.MockService,
)) (staff.Service, func()) {
	ctrl := gomock.NewController(t)

	repo := staffmocks.NewMockRepository(ctrl)
	authCoreService := authcoremocks.NewMockService(ctrl)

	if mocksSetup != nil {
		mocksSetup(repo, authCoreService)
	}

	service := staff.NewService(logger, repo, authCoreService)
	return service, func() {
		ctrl.Finish()
	}
}
