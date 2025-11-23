//go:build unit

package authcore_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/alexgrauroca/practice-food-delivery-platform/pkg/auth"
	authmocks "github.com/alexgrauroca/practice-food-delivery-platform/pkg/auth/mocks"
	"github.com/alexgrauroca/practice-food-delivery-platform/pkg/log"
	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/authcore"
	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/refresh"
	refreshmocks "github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/refresh/mocks"
)

type authCoreTestsCase[I, W any] struct {
	name       string
	input      I
	mocksSetup func(
		authService *authmocks.MockService,
		refreshService *refreshmocks.MockService,
	)
	want    W
	wantErr error
}

var errUnexpected = errors.New("unexpected error")

func TestService_GenerateTokenPair(t *testing.T) {
	logger, _ := log.NewTest()

	tests := []authCoreTestsCase[authcore.GenerateTokenPairInput, authcore.TokenPair]{
		{
			name: "when there is an error generating the token, then it propagates the error",
			input: authcore.GenerateTokenPairInput{
				UserID:     "fake-id",
				Expiration: 3600,
				Role:       "fake-role",
				TenantID:   "fake-tenant-id",
			},
			mocksSetup: func(authService *authmocks.MockService, _ *refreshmocks.MockService) {
				authService.EXPECT().GenerateToken(gomock.Any(), gomock.Any()).
					Return(auth.GenerateTokenOutput{}, errUnexpected)
			},
			want:    authcore.TokenPair{},
			wantErr: errUnexpected,
		},
		{
			name: "when there is an error generating the refresh token, then it propagates the error",
			input: authcore.GenerateTokenPairInput{
				UserID:     "fake-id",
				Expiration: 3600,
				Role:       "fake-role",
				TenantID:   "fake-tenant-id",
			},
			mocksSetup: func(authService *authmocks.MockService, refreshService *refreshmocks.MockService) {
				authService.EXPECT().GenerateToken(gomock.Any(), gomock.Any()).
					Return(auth.GenerateTokenOutput{
						AccessToken: "fake-access-token",
					}, nil)

				refreshService.EXPECT().Generate(gomock.Any(), gomock.Any()).
					Return(refresh.GenerateTokenOutput{}, errUnexpected)
			},
			want:    authcore.TokenPair{},
			wantErr: errUnexpected,
		},
		{
			name: "when the token is generated correctly, then it returns the token pair",
			input: authcore.GenerateTokenPairInput{
				UserID:     "fake-id",
				Expiration: 3600,
				Role:       "fake-role",
				TenantID:   "fake-tenant-id",
			},
			mocksSetup: func(authService *authmocks.MockService, refreshService *refreshmocks.MockService) {
				authService.EXPECT().GenerateToken(gomock.Any(), auth.GenerateTokenInput{
					ID:         "fake-id",
					Expiration: 3600,
					Role:       "fake-role",
					TenantID:   "fake-tenant-id",
				}).Return(auth.GenerateTokenOutput{
					AccessToken: "fake-access-token",
				}, nil)

				refreshService.EXPECT().Generate(gomock.Any(), refresh.GenerateTokenInput{
					UserID:   "fake-id",
					Role:     "fake-role",
					TenantID: "fake-tenant-id",
				}).Return(refresh.GenerateTokenOutput{
					Token: "fake-refresh-token",
				}, nil)
			},
			want: authcore.TokenPair{
				AccessToken:  "fake-access-token",
				RefreshToken: "fake-refresh-token",
				ExpiresIn:    3600,
				TokenType:    "Bearer",
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, cleanup := serviceSetup(t, logger, tt.mocksSetup)
			defer cleanup()

			got, err := service.GenerateTokenPair(context.Background(), tt.input)

			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestService_RefreshToken(t *testing.T) {
	logger, _ := log.NewTest()

	tests := []authCoreTestsCase[authcore.RefreshTokenInput, authcore.TokenPair]{
		{
			name: "when the refresh token is not found, then it returns an invalid refresh token error",
			input: authcore.RefreshTokenInput{
				RefreshToken: "InvalidRefreshToken",
			},
			mocksSetup: func(_ *authmocks.MockService, refreshService *refreshmocks.MockService) {
				refreshService.EXPECT().FindActiveToken(gomock.Any(), gomock.Any()).
					Return(refresh.FindActiveTokenOutput{}, refresh.ErrRefreshTokenNotFound)
			},
			want:    authcore.TokenPair{},
			wantErr: authcore.ErrInvalidRefreshToken,
		},
		{
			name: "when there is an unexpected error when finding the refresh token, " +
				"then it should propagate the error",
			input: authcore.RefreshTokenInput{
				RefreshToken: "ValidRefreshToken",
			},
			mocksSetup: func(_ *authmocks.MockService, refreshService *refreshmocks.MockService) {
				refreshService.EXPECT().FindActiveToken(gomock.Any(), gomock.Any()).
					Return(refresh.FindActiveTokenOutput{}, errUnexpected)
			},
			want:    authcore.TokenPair{},
			wantErr: errUnexpected,
		},
		{
			name: "when the access token is invalid, then it should return a token mismatch error",
			input: authcore.RefreshTokenInput{
				AccessToken:  "InvalidAccessToken",
				RefreshToken: "ValidRefreshToken",
			},
			mocksSetup: func(authService *authmocks.MockService, refreshService *refreshmocks.MockService) {
				refreshService.EXPECT().FindActiveToken(gomock.Any(), gomock.Any()).
					Return(refresh.FindActiveTokenOutput{
						ID:       "fake-id",
						Token:    "ValidRefreshToken",
						UserID:   "fake-user-id",
						Role:     "fake-role",
						TenantID: "fake-tenant-id",
						Device:   refresh.DeviceInfo{}, // device info is irrelevant here
					}, nil)

				authService.EXPECT().GetClaims(gomock.Any(), gomock.Any()).
					Return(auth.GetClaimsOutput{}, auth.ErrInvalidToken)
			},
			want:    authcore.TokenPair{},
			wantErr: authcore.ErrTokenMismatch,
		},
		{
			name: "when there is an unexpected error when getting the claims from the access token, " +
				"then it should propagate the error",
			input: authcore.RefreshTokenInput{
				AccessToken:  "ValidAccessToken",
				RefreshToken: "ValidRefreshToken",
			},
			mocksSetup: func(authService *authmocks.MockService, refreshService *refreshmocks.MockService) {
				refreshService.EXPECT().FindActiveToken(gomock.Any(), gomock.Any()).
					Return(refresh.FindActiveTokenOutput{
						ID:       "fake-id",
						Token:    "ValidRefreshToken",
						UserID:   "fake-user-id",
						Role:     "fake-role",
						TenantID: "fake-tenant-id",
						Device:   refresh.DeviceInfo{}, // device info is irrelevant here
					}, nil)

				authService.EXPECT().GetClaims(gomock.Any(), gomock.Any()).
					Return(auth.GetClaimsOutput{}, errUnexpected)
			},
			want:    authcore.TokenPair{},
			wantErr: errUnexpected,
		},
		{
			name: "when the user ID from the claims is different than the one from the refresh token, " +
				"then it should return a token mismatch error",
			input: authcore.RefreshTokenInput{
				AccessToken:  "ValidAccessToken",
				RefreshToken: "ValidRefreshToken",
			},
			mocksSetup: func(authService *authmocks.MockService, refreshService *refreshmocks.MockService) {
				refreshService.EXPECT().FindActiveToken(gomock.Any(), gomock.Any()).
					Return(refresh.FindActiveTokenOutput{
						ID:       "fake-id",
						Token:    "ValidRefreshToken",
						UserID:   "fake-user-id",
						Role:     "fake-role",
						TenantID: "fake-tenant-id",
						Device:   refresh.DeviceInfo{}, // device info is irrelevant here
					}, nil)

				authService.EXPECT().GetClaims(gomock.Any(), gomock.Any()).
					Return(auth.GetClaimsOutput{
						Claims: &auth.Claims{
							RegisteredClaims: jwt.RegisteredClaims{
								Subject: "invalid-user-id",
							},
							Role:   "fake-role",
							Tenant: "fake-tenant-id",
						},
					}, nil)
			},
			want:    authcore.TokenPair{},
			wantErr: authcore.ErrTokenMismatch,
		},
		{
			name: "when the role from the claims is different than the one from the refresh token, " +
				"then it should return a token mismatch error",
			input: authcore.RefreshTokenInput{
				AccessToken:  "ValidAccessToken",
				RefreshToken: "ValidRefreshToken",
			},
			mocksSetup: func(authService *authmocks.MockService, refreshService *refreshmocks.MockService) {
				refreshService.EXPECT().FindActiveToken(gomock.Any(), gomock.Any()).
					Return(refresh.FindActiveTokenOutput{
						ID:       "fake-id",
						Token:    "ValidRefreshToken",
						UserID:   "fake-user-id",
						Role:     "fake-role",
						TenantID: "fake-tenant-id",
						Device:   refresh.DeviceInfo{}, // device info is irrelevant here
					}, nil)

				authService.EXPECT().GetClaims(gomock.Any(), gomock.Any()).
					Return(auth.GetClaimsOutput{
						Claims: &auth.Claims{
							RegisteredClaims: jwt.RegisteredClaims{
								Subject: "fake-user-id",
							},
							Role:   "invalid-role",
							Tenant: "fake-tenant-id",
						},
					}, nil)
			},
			want:    authcore.TokenPair{},
			wantErr: authcore.ErrTokenMismatch,
		},
		{
			name: "when the tenant from the claims is different than the one from the refresh token, " +
				"then it should return a token mismatch error",
			input: authcore.RefreshTokenInput{
				AccessToken:  "ValidAccessToken",
				RefreshToken: "ValidRefreshToken",
			},
			mocksSetup: func(authService *authmocks.MockService, refreshService *refreshmocks.MockService) {
				refreshService.EXPECT().FindActiveToken(gomock.Any(), gomock.Any()).
					Return(refresh.FindActiveTokenOutput{
						ID:       "fake-id",
						Token:    "ValidRefreshToken",
						UserID:   "fake-user-id",
						Role:     "fake-role",
						TenantID: "fake-tenant-id",
						Device:   refresh.DeviceInfo{}, // device info is irrelevant here
					}, nil)

				authService.EXPECT().GetClaims(gomock.Any(), gomock.Any()).
					Return(auth.GetClaimsOutput{
						Claims: &auth.Claims{
							RegisteredClaims: jwt.RegisteredClaims{
								Subject: "fake-user-id",
							},
							Role:   "invalid-role",
							Tenant: "",
						},
					}, nil)
			},
			want:    authcore.TokenPair{},
			wantErr: authcore.ErrTokenMismatch,
		},
		{
			name: "when there is an error generating the new token, then it should propagate the error",
			input: authcore.RefreshTokenInput{
				AccessToken:  "ValidAccessToken",
				RefreshToken: "ValidRefreshToken",
			},
			mocksSetup: func(authService *authmocks.MockService, refreshService *refreshmocks.MockService) {
				refreshService.EXPECT().FindActiveToken(gomock.Any(), gomock.Any()).
					Return(refresh.FindActiveTokenOutput{
						ID:       "fake-id",
						Token:    "ValidRefreshToken",
						UserID:   "fake-user-id",
						Role:     "fake-role",
						TenantID: "fake-tenant-id",
						Device:   refresh.DeviceInfo{}, // device info is irrelevant here
					}, nil)

				authService.EXPECT().GetClaims(gomock.Any(), gomock.Any()).
					Return(auth.GetClaimsOutput{
						Claims: &auth.Claims{
							RegisteredClaims: jwt.RegisteredClaims{
								Subject: "fake-user-id",
							},
							Role:   "fake-role",
							Tenant: "fake-tenant-id",
						},
					}, nil)

				authService.EXPECT().GenerateToken(gomock.Any(), gomock.Any()).
					Return(auth.GenerateTokenOutput{}, errUnexpected)
			},
			want:    authcore.TokenPair{},
			wantErr: errUnexpected,
		},
		{
			name: "when there is an error generating the new refresh token, then it should propagate the error",
			input: authcore.RefreshTokenInput{
				AccessToken:  "ValidAccessToken",
				RefreshToken: "ValidRefreshToken",
			},
			mocksSetup: func(authService *authmocks.MockService, refreshService *refreshmocks.MockService) {
				refreshService.EXPECT().FindActiveToken(gomock.Any(), gomock.Any()).
					Return(refresh.FindActiveTokenOutput{
						ID:       "fake-id",
						Token:    "ValidRefreshToken",
						UserID:   "fake-user-id",
						Role:     "fake-role",
						TenantID: "fake-tenant-id",
						Device:   refresh.DeviceInfo{}, // device info is irrelevant here
					}, nil)

				authService.EXPECT().GetClaims(gomock.Any(), gomock.Any()).
					Return(auth.GetClaimsOutput{
						Claims: &auth.Claims{
							RegisteredClaims: jwt.RegisteredClaims{
								Subject: "fake-user-id",
							},
							Role:   "fake-role",
							Tenant: "fake-tenant-id",
						},
					}, nil)

				authService.EXPECT().GenerateToken(gomock.Any(), gomock.Any()).
					Return(auth.GenerateTokenOutput{
						AccessToken: "fake-access-token",
					}, nil)

				refreshService.EXPECT().Generate(gomock.Any(), gomock.Any()).
					Return(refresh.GenerateTokenOutput{}, errUnexpected)
			},
			want:    authcore.TokenPair{},
			wantErr: errUnexpected,
		},
		{
			name: "when there is an unexpected error expiring the refresh token, then it should propagate the error",
			input: authcore.RefreshTokenInput{
				AccessToken:  "ValidAccessToken",
				RefreshToken: "ValidRefreshToken",
			},
			mocksSetup: func(authService *authmocks.MockService, refreshService *refreshmocks.MockService) {
				refreshService.EXPECT().FindActiveToken(gomock.Any(), gomock.Any()).
					Return(refresh.FindActiveTokenOutput{
						ID:       "fake-id",
						Token:    "ValidRefreshToken",
						UserID:   "fake-user-id",
						Role:     "fake-role",
						TenantID: "fake-tenant-id",
						Device:   refresh.DeviceInfo{}, // device info is irrelevant here
					}, nil)

				authService.EXPECT().GetClaims(gomock.Any(), gomock.Any()).
					Return(auth.GetClaimsOutput{
						Claims: &auth.Claims{
							RegisteredClaims: jwt.RegisteredClaims{
								Subject: "fake-user-id",
							},
							Role:   "fake-role",
							Tenant: "fake-tenant-id",
						},
					}, nil)

				authService.EXPECT().GenerateToken(gomock.Any(), gomock.Any()).
					Return(auth.GenerateTokenOutput{
						AccessToken: "fake-access-token",
					}, nil)

				refreshService.EXPECT().Generate(gomock.Any(), gomock.Any()).
					Return(refresh.GenerateTokenOutput{
						Token: "fake-refresh-token",
					}, nil)

				refreshService.EXPECT().Expire(gomock.Any(), gomock.Any()).
					Return(refresh.ExpireOutput{}, errUnexpected)
			},
			want:    authcore.TokenPair{},
			wantErr: errUnexpected,
		},
		{
			name: "when the refresh token is not found, then it should return the new token pair",
			input: authcore.RefreshTokenInput{
				AccessToken:  "ValidAccessToken",
				RefreshToken: "ValidRefreshToken",
				Expiration:   3600,
				Role:         "ValidRole",
			},
			mocksSetup: func(authService *authmocks.MockService, refreshService *refreshmocks.MockService) {
				refreshService.EXPECT().FindActiveToken(gomock.Any(), gomock.Any()).
					Return(refresh.FindActiveTokenOutput{
						ID:       "fake-id",
						Token:    "ValidRefreshToken",
						UserID:   "fake-user-id",
						Role:     "fake-role",
						TenantID: "fake-tenant-id",
						Device:   refresh.DeviceInfo{}, // device info is irrelevant here
					}, nil)

				authService.EXPECT().GetClaims(gomock.Any(), gomock.Any()).
					Return(auth.GetClaimsOutput{
						Claims: &auth.Claims{
							RegisteredClaims: jwt.RegisteredClaims{
								Subject: "fake-user-id",
							},
							Role:   "fake-role",
							Tenant: "fake-tenant-id",
						},
					}, nil)

				authService.EXPECT().GenerateToken(gomock.Any(), gomock.Any()).
					Return(auth.GenerateTokenOutput{
						AccessToken: "fake-access-token",
					}, nil)

				refreshService.EXPECT().Generate(gomock.Any(), gomock.Any()).
					Return(refresh.GenerateTokenOutput{
						Token: "fake-refresh-token",
					}, nil)

				refreshService.EXPECT().Expire(gomock.Any(), gomock.Any()).
					Return(refresh.ExpireOutput{}, refresh.ErrRefreshTokenNotFound)
			},
			want: authcore.TokenPair{
				AccessToken:  "fake-access-token",
				RefreshToken: "fake-refresh-token",
				ExpiresIn:    3600,
				TokenType:    "Bearer",
			},
			wantErr: nil,
		},
		{
			name: "when the token is refreshed correctly, then it should return the new token pair",
			input: authcore.RefreshTokenInput{
				AccessToken:  "ValidAccessToken",
				RefreshToken: "ValidRefreshToken",
				Expiration:   3600,
				Role:         "ValidRole",
			},
			mocksSetup: func(authService *authmocks.MockService, refreshService *refreshmocks.MockService) {
				refreshService.EXPECT().FindActiveToken(gomock.Any(), refresh.FindActiveTokenInput{
					Token: "ValidRefreshToken",
				}).Return(refresh.FindActiveTokenOutput{
					ID:       "fake-id",
					Token:    "ValidRefreshToken",
					UserID:   "fake-user-id",
					Role:     "fake-role",
					TenantID: "fake-tenant-id",
					Device:   refresh.DeviceInfo{}, // device info is irrelevant here
				}, nil)

				authService.EXPECT().GetClaims(gomock.Any(), auth.GetClaimsInput{
					AccessToken: "ValidAccessToken",
				}).Return(auth.GetClaimsOutput{
					Claims: &auth.Claims{
						RegisteredClaims: jwt.RegisteredClaims{
							Subject: "fake-user-id",
						},
						Role:   "fake-role",
						Tenant: "fake-tenant-id",
					},
				}, nil)

				authService.EXPECT().GenerateToken(gomock.Any(), auth.GenerateTokenInput{
					ID:         "fake-user-id",
					Expiration: 3600,
					Role:       "ValidRole",
					TenantID:   "fake-tenant-id",
				}).Return(auth.GenerateTokenOutput{
					AccessToken: "fake-access-token",
				}, nil)

				refreshService.EXPECT().Generate(gomock.Any(), refresh.GenerateTokenInput{
					UserID:   "fake-user-id",
					Role:     "ValidRole",
					TenantID: "fake-tenant-id",
				}).Return(refresh.GenerateTokenOutput{
					Token: "fake-refresh-token",
				}, nil)

				refreshService.EXPECT().Expire(gomock.Any(), refresh.ExpireInput{
					Token: "ValidRefreshToken",
				}).Return(refresh.ExpireOutput{}, nil)
			},
			want: authcore.TokenPair{
				AccessToken:  "fake-access-token",
				RefreshToken: "fake-refresh-token",
				ExpiresIn:    3600,
				TokenType:    "Bearer",
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, cleanup := serviceSetup(t, logger, tt.mocksSetup)
			defer cleanup()

			got, err := service.RefreshToken(context.Background(), tt.input)

			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func serviceSetup(
	t *testing.T, logger log.Logger, mocksSetup func(
		authService *authmocks.MockService,
		refreshService *refreshmocks.MockService,
	),
) (authcore.Service, func()) {
	ctrl := gomock.NewController(t)

	authService := authmocks.NewMockService(ctrl)
	refreshService := refreshmocks.NewMockService(ctrl)

	if mocksSetup != nil {
		mocksSetup(authService, refreshService)
	}

	service := authcore.NewService(logger, authService, refreshService)
	return service, func() {
		ctrl.Finish()
	}
}
