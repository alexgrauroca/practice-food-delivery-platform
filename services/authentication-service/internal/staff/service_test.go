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

	"github.com/alexgrauroca/practice-food-delivery-platform/pkg/auth"
	authmocks "github.com/alexgrauroca/practice-food-delivery-platform/pkg/auth/mocks"
	"github.com/alexgrauroca/practice-food-delivery-platform/pkg/log"
	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/authcore"
	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/password"
	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/refresh"
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
				_ *refreshmocks.MockService,
				_ *authmocks.MockService,
				_ *authmocks.MockContextReader,
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
				_ *refreshmocks.MockService,
				_ *authmocks.MockService,
				_ *authmocks.MockContextReader,
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
				_ *refreshmocks.MockService,
				_ *authmocks.MockService,
				_ *authmocks.MockContextReader,
			) {
				repo.EXPECT().FindByEmail(gomock.Any(), gomock.Any()).
					Return(staff.Staff{}, errRepo)
			},
			want:    staff.LoginStaffOutput{},
			wantErr: errRepo,
		},
		{
			name: "when there is an error generating the jwt, then it should propagate the error",
			input: staff.LoginStaffInput{
				Email:    "test@example.com",
				Password: "ValidPassword123",
			},
			mocksSetup: func(
				repo *staffmocks.MockRepository,
				_ *refreshmocks.MockService,
				authService *authmocks.MockService,
				_ *authmocks.MockContextReader,
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

				authService.EXPECT().GenerateToken(gomock.Any(), gomock.Any()).
					Return(auth.GenerateTokenOutput{}, errToken)
			},
			want:    staff.LoginStaffOutput{},
			wantErr: errToken,
		},
		{
			name: "when there is an error generating the refresh token, then it should propagate the error",
			input: staff.LoginStaffInput{
				Email:    "test@example.com",
				Password: "ValidPassword123",
			},
			mocksSetup: func(
				repo *staffmocks.MockRepository,
				refreshService *refreshmocks.MockService,
				authService *authmocks.MockService,
				_ *authmocks.MockContextReader,
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

				authService.EXPECT().GenerateToken(gomock.Any(), gomock.Any()).
					Return(auth.GenerateTokenOutput{AccessToken: "fake-token"}, nil)

				refreshService.EXPECT().Generate(gomock.Any(), gomock.Any()).
					Return(refresh.GenerateTokenOutput{}, errToken)
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
				refreshService *refreshmocks.MockService,
				authService *authmocks.MockService,
				_ *authmocks.MockContextReader,
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

				authService.EXPECT().GenerateToken(gomock.Any(), gomock.Any()).
					Return(auth.GenerateTokenOutput{AccessToken: "fake-token"}, nil)

				refreshService.EXPECT().Generate(gomock.Any(), refresh.GenerateTokenInput{
					UserID: "fake-id",
					Role:   "staff",
				}).Return(refresh.GenerateTokenOutput{Token: "fake-refresh-token"}, nil)
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

func serviceSetup(t *testing.T, logger log.Logger, mocksSetup func(repo *staffmocks.MockRepository,
	refreshService *refreshmocks.MockService,
	authService *authmocks.MockService,
	authctx *authmocks.MockContextReader,
)) (staff.Service, func()) {
	ctrl := gomock.NewController(t)

	repo := staffmocks.NewMockRepository(ctrl)
	refreshService := refreshmocks.NewMockService(ctrl)
	authService := authmocks.NewMockService(ctrl)
	authctx := authmocks.NewMockContextReader(ctrl)

	if mocksSetup != nil {
		mocksSetup(repo, refreshService, authService, authctx)
	}

	service := staff.NewService(logger, repo, refreshService, authService, authctx)
	return service, func() {
		ctrl.Finish()
	}
}
