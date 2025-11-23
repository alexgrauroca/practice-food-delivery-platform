//go:build unit

package refresh_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/alexgrauroca/practice-food-delivery-platform/pkg/clock"
	"github.com/alexgrauroca/practice-food-delivery-platform/pkg/log"
	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/refresh"
	refreshmocks "github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/refresh/mocks"
)

var errRepo = errors.New("repository error")

type refreshServiceTestCase[I, W any] struct {
	name       string
	input      I
	mocksSetup func(repo *refreshmocks.MockRepository)
	want       W
	wantErr    error
}

func TestService_Generate(t *testing.T) {
	logger, _ := log.NewTest()

	tests := []refreshServiceTestCase[refresh.GenerateTokenInput, refresh.GenerateTokenOutput]{
		{
			name: "when there is an error storing the refresh token, then it propagates the error",
			input: refresh.GenerateTokenInput{
				UserID:   "fake-user-id",
				Role:     "fake-role",
				TenantID: "fake-tenant-id",
			},
			mocksSetup: func(repo *refreshmocks.MockRepository) {
				repo.EXPECT().Create(gomock.Any(), gomock.Any()).
					Return(refresh.Token{}, errRepo)
			},
			want:    refresh.GenerateTokenOutput{},
			wantErr: errRepo,
		},
		{
			name: "when the refresh token is generated and stored, then it returns the token",
			input: refresh.GenerateTokenInput{
				UserID:   "fake-user-id",
				Role:     "fake-role",
				TenantID: "fake-tenant-id",
			},
			mocksSetup: func(repo *refreshmocks.MockRepository) {
				repo.EXPECT().Create(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, params refresh.CreateTokenParams) (refresh.Token, error) {
						require.Equal(t, "fake-user-id", params.UserID)
						require.Equal(t, "fake-role", params.Role)
						require.Equal(t, "fake-tenant-id", params.TenantID)
						require.NotEmpty(t, params.Token)
						require.NotEmpty(t, params.ExpiresAt)

						// Returning a fake-token to simplify the assertion
						return refresh.Token{Token: "fake-token"}, nil
					})
			},
			want:    refresh.GenerateTokenOutput{Token: "fake-token"},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := refreshmocks.NewMockRepository(ctrl)
			if tt.mocksSetup != nil {
				tt.mocksSetup(repo)
			}

			service := refresh.NewService(logger, repo, clock.RealClock{})
			got, err := service.Generate(context.Background(), tt.input)

			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestService_FindByActiveToken(t *testing.T) {
	var (
		now       = time.Date(2025, 1, 7, 0, 0, 0, 0, time.UTC)
		yesterday = time.Date(2025, 1, 6, 0, 0, 0, 0, time.UTC)
		tomorrow  = time.Date(2025, 1, 8, 0, 0, 0, 0, time.UTC)
	)
	logger, _ := log.NewTest()

	tests := []refreshServiceTestCase[refresh.FindActiveTokenInput, refresh.FindActiveTokenOutput]{
		{
			name: "when unable to find the token active, then it returns a refresh token not found error",
			input: refresh.FindActiveTokenInput{
				Token: "inactive-token",
			},
			mocksSetup: func(repo *refreshmocks.MockRepository) {
				repo.EXPECT().FindActiveToken(gomock.Any(), gomock.Any()).
					Return(refresh.Token{}, refresh.ErrRefreshTokenNotFound)
			},
			want:    refresh.FindActiveTokenOutput{},
			wantErr: refresh.ErrRefreshTokenNotFound,
		},
		{
			name: "when there is an unexpected error when finding the token, then it propagates the error",
			input: refresh.FindActiveTokenInput{
				Token: "active-token",
			},
			mocksSetup: func(repo *refreshmocks.MockRepository) {
				repo.EXPECT().FindActiveToken(gomock.Any(), gomock.Any()).
					Return(refresh.Token{}, errRepo)
			},
			want:    refresh.FindActiveTokenOutput{},
			wantErr: errRepo,
		},
		{
			name: "when the active refresh token is found, then it returns the token",
			input: refresh.FindActiveTokenInput{
				Token: "active-token",
			},
			mocksSetup: func(repo *refreshmocks.MockRepository) {
				repo.EXPECT().FindActiveToken(gomock.Any(), "active-token").
					Return(refresh.Token{
						ID:       "fake-id",
						UserID:   "fake-user-id",
						Role:     "fake-role",
						TenantID: "fake-tenant-id",
						Token:    "fake-token",
						Status:   refresh.TokenStatusActive,
						DeviceInfo: refresh.DeviceInfo{
							DeviceID:    "fake-device-id",
							UserAgent:   "fake-user-agent",
							IP:          "fake-ip",
							FirstUsedAt: yesterday,
							LastUsedAt:  now,
						},
						ExpiresAt: tomorrow,
						CreatedAt: yesterday,
						UpdatedAt: now,
					}, nil)
			},
			want: refresh.FindActiveTokenOutput{
				ID:       "fake-id",
				Token:    "fake-token",
				UserID:   "fake-user-id",
				Role:     "fake-role",
				TenantID: "fake-tenant-id",
				Device: refresh.DeviceInfo{
					DeviceID:    "fake-device-id",
					UserAgent:   "fake-user-agent",
					IP:          "fake-ip",
					FirstUsedAt: yesterday,
					LastUsedAt:  now,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := refreshmocks.NewMockRepository(ctrl)
			if tt.mocksSetup != nil {
				tt.mocksSetup(repo)
			}

			service := refresh.NewService(logger, repo, clock.RealClock{})
			got, err := service.FindActiveToken(context.Background(), tt.input)

			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestService_Expire(t *testing.T) {
	var (
		now       = time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
		yesterday = time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)
	)
	logger, _ := log.NewTest()

	tests := []refreshServiceTestCase[refresh.ExpireInput, refresh.ExpireOutput]{
		{
			name: "when unable to expire the token, then it propagates the error",
			input: refresh.ExpireInput{
				Token: "fake-token",
			},
			mocksSetup: func(repo *refreshmocks.MockRepository) {
				repo.EXPECT().Expire(gomock.Any(), gomock.Any()).Return(refresh.Token{}, errRepo)
			},
			want:    refresh.ExpireOutput{},
			wantErr: errRepo,
		},
		{
			name: "when the token is expired, then it returns the token",
			input: refresh.ExpireInput{
				Token: "fake-token",
			},
			mocksSetup: func(repo *refreshmocks.MockRepository) {
				repo.EXPECT().Expire(gomock.Any(), refresh.ExpireParams{
					Token:     "fake-token",
					ExpiresAt: now.Add(5 * time.Second),
				}).Return(refresh.Token{
					ID:        "fake-token-id",
					Role:      "fake-role",
					TenantID:  "fake-tenant-id",
					Token:     "fake-token",
					Status:    refresh.TokenStatusActive,
					ExpiresAt: now.Add(5 * time.Second),
					CreatedAt: yesterday,
					UpdatedAt: now,
				}, nil)
			},
			want: refresh.ExpireOutput{
				ID:        "fake-token-id",
				Token:     "fake-token",
				ExpiresAt: now.Add(5 * time.Second),
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := refreshmocks.NewMockRepository(ctrl)
			if tt.mocksSetup != nil {
				tt.mocksSetup(repo)
			}

			service := refresh.NewService(logger, repo, clock.FixedClock{FixedTime: now})
			got, err := service.Expire(context.Background(), tt.input)

			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}
