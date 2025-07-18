//go:build unit

package customers_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/customers"
	customersmocks "github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/customers/mocks"
	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/jwt"
	jwtmocks "github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/jwt/mocks"
	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/log"
	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/password"
	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/refresh"
	refreshmocks "github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/refresh/mocks"
)

var (
	errRepo  = errors.New("repository error")
	errToken = errors.New("token error")
)

type customersServiceTestCase[I, W any] struct {
	name       string
	input      I
	mocksSetup func(repo *customersmocks.MockRepository, refreshService *refreshmocks.MockService,
		jwtService *jwtmocks.MockService)
	want    W
	wantErr error
}

func TestService_RegisterCustomer(t *testing.T) {
	now := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	logger, _ := log.NewTest()

	tests := []customersServiceTestCase[customers.RegisterCustomerInput, customers.RegisterCustomerOutput]{
		{
			name: "when there is an active customer with the same email, then it should return a customer already exists error",
			input: customers.RegisterCustomerInput{
				Email:    "test@example.com",
				Password: "ValidPassword123",
				Name:     "John Doe",
			},
			mocksSetup: func(repo *customersmocks.MockRepository, _ *refreshmocks.MockService,
				_ *jwtmocks.MockService) {
				repo.EXPECT().CreateCustomer(gomock.Any(), gomock.Any()).
					Return(customers.Customer{}, customers.ErrCustomerAlreadyExists)
			},
			want:    customers.RegisterCustomerOutput{},
			wantErr: customers.ErrCustomerAlreadyExists,
		},
		{
			name: "when there is an unexpected error when creating the customer, then it should propagate the error",
			input: customers.RegisterCustomerInput{
				Email:    "test@example.com",
				Password: "ValidPassword123",
				Name:     "John Doe",
			},
			mocksSetup: func(repo *customersmocks.MockRepository, _ *refreshmocks.MockService,
				_ *jwtmocks.MockService) {
				repo.EXPECT().CreateCustomer(gomock.Any(), gomock.Any()).
					Return(customers.Customer{}, errRepo)
			},
			want:    customers.RegisterCustomerOutput{},
			wantErr: errRepo,
		},
		{
			name: "when the customer can be created, then it should return the created customer",
			input: customers.RegisterCustomerInput{
				Email:    "test@example.com",
				Password: "ValidPassword123",
				Name:     "John Doe",
			},
			mocksSetup: func(repo *customersmocks.MockRepository, _ *refreshmocks.MockService,
				_ *jwtmocks.MockService) {
				repo.EXPECT().CreateCustomer(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, params customers.CreateCustomerParams) (customers.Customer, error) {
						// Assert that the password is hashed
						ok := password.Verify(params.Password, "ValidPassword123")
						require.True(t, ok, "Password should be hashed and match the input password")

						return customers.Customer{
							ID:        "fake-id",
							Email:     params.Email,
							Name:      params.Name,
							Password:  params.Password,
							CreatedAt: now,
							UpdatedAt: now,
							Active:    true,
						}, nil
					})
			},
			want: customers.RegisterCustomerOutput{
				ID:        "fake-id",
				Email:     "test@example.com",
				Name:      "John Doe",
				CreatedAt: now,
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := customersmocks.NewMockRepository(ctrl)
			refreshService := refreshmocks.NewMockService(ctrl)
			jwtService := jwtmocks.NewMockService(ctrl)
			if tt.mocksSetup != nil {
				tt.mocksSetup(repo, refreshService, jwtService)
			}

			service := customers.NewService(logger, repo, refreshService, jwtService)
			got, err := service.RegisterCustomer(context.Background(), tt.input)

			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestService_LoginCustomer(t *testing.T) {
	now := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	logger, _ := log.NewTest()

	tests := []customersServiceTestCase[customers.LoginCustomerInput, customers.LoginCustomerOutput]{
		{
			name: "when there is not an active customer with the same email, " +
				"then it should return an invalid credentials error",
			input: customers.LoginCustomerInput{
				Email:    "test@example.com",
				Password: "ValidPassword123",
			},
			mocksSetup: func(repo *customersmocks.MockRepository, _ *refreshmocks.MockService,
				_ *jwtmocks.MockService) {
				repo.EXPECT().FindByEmail(gomock.Any(), gomock.Any()).
					Return(customers.Customer{}, customers.ErrCustomerNotFound)
			},
			want:    customers.LoginCustomerOutput{},
			wantErr: customers.ErrInvalidCredentials,
		},
		{
			name: "when there is not an active customer with the same password, " +
				"then it should return an invalid credentials error",
			input: customers.LoginCustomerInput{
				Email:    "test@example.com",
				Password: "InvalidPassword123",
			},
			mocksSetup: func(repo *customersmocks.MockRepository, _ *refreshmocks.MockService,
				_ *jwtmocks.MockService) {
				hashedPassword, err := password.Hash("ValidPassword123")
				require.NoError(t, err)

				repo.EXPECT().FindByEmail(gomock.Any(), gomock.Any()).
					Return(customers.Customer{
						ID:        "fake-id",
						Email:     "test@example.com",
						Name:      "John Doe",
						Password:  hashedPassword, // This should be a hashed password
						CreatedAt: now,
						UpdatedAt: now,
						Active:    true,
					}, nil)
			},
			want:    customers.LoginCustomerOutput{},
			wantErr: customers.ErrInvalidCredentials,
		},
		{
			name: "when there is an unexpected error when fetching the customer, then it should propagate the error",
			input: customers.LoginCustomerInput{
				Email:    "test@example.com",
				Password: "ValidPassword123",
			},
			mocksSetup: func(repo *customersmocks.MockRepository, _ *refreshmocks.MockService,
				_ *jwtmocks.MockService) {
				repo.EXPECT().FindByEmail(gomock.Any(), gomock.Any()).
					Return(customers.Customer{}, errRepo)
			},
			want:    customers.LoginCustomerOutput{},
			wantErr: errRepo,
		},
		{
			name: "when there is an error generating the jwt, then it should propagate the error",
			input: customers.LoginCustomerInput{
				Email:    "test@example.com",
				Password: "ValidPassword123",
			},
			mocksSetup: func(repo *customersmocks.MockRepository, _ *refreshmocks.MockService,
				jwtService *jwtmocks.MockService) {
				hashedPassword, err := password.Hash("ValidPassword123")
				require.NoError(t, err)

				repo.EXPECT().FindByEmail(gomock.Any(), "test@example.com").
					Return(customers.Customer{
						ID:        "fake-id",
						Email:     "test@example.com",
						Name:      "John Doe",
						Password:  hashedPassword, // This should be a hashed password
						CreatedAt: now,
						UpdatedAt: now,
						Active:    true,
					}, nil)

				jwtService.EXPECT().GenerateToken(gomock.Any(), gomock.Any()).
					Return("", errToken)
			},
			want:    customers.LoginCustomerOutput{},
			wantErr: errToken,
		},
		{
			name: "when there is an error generating the refresh token, then it should propagate the error",
			input: customers.LoginCustomerInput{
				Email:    "test@example.com",
				Password: "ValidPassword123",
			},
			mocksSetup: func(repo *customersmocks.MockRepository, refreshService *refreshmocks.MockService,
				jwtService *jwtmocks.MockService) {
				hashedPassword, err := password.Hash("ValidPassword123")
				require.NoError(t, err)

				repo.EXPECT().FindByEmail(gomock.Any(), "test@example.com").
					Return(customers.Customer{
						ID:        "fake-id",
						Email:     "test@example.com",
						Name:      "John Doe",
						Password:  hashedPassword, // This should be a hashed password
						CreatedAt: now,
						UpdatedAt: now,
						Active:    true,
					}, nil)

				jwtService.EXPECT().GenerateToken(gomock.Any(), gomock.Any()).
					Return("fake-token", nil)

				refreshService.EXPECT().Generate(gomock.Any(), gomock.Any()).
					Return(refresh.GenerateTokenOutput{}, errToken)
			},
			want:    customers.LoginCustomerOutput{},
			wantErr: errToken,
		},
		{
			name: "when there is an active customer with the same email and password, then it should return its token",
			input: customers.LoginCustomerInput{
				Email:    "test@example.com",
				Password: "ValidPassword123",
			},
			mocksSetup: func(repo *customersmocks.MockRepository, refreshService *refreshmocks.MockService,
				jwtService *jwtmocks.MockService) {
				hashedPassword, err := password.Hash("ValidPassword123")
				require.NoError(t, err)

				repo.EXPECT().FindByEmail(gomock.Any(), "test@example.com").
					Return(customers.Customer{
						ID:        "fake-id",
						Email:     "test@example.com",
						Name:      "John Doe",
						Password:  hashedPassword, // This should be a hashed password
						CreatedAt: now,
						UpdatedAt: now,
						Active:    true,
					}, nil)

				jwtService.EXPECT().GenerateToken(gomock.Any(), gomock.Any()).
					Return("fake-token", nil)

				refreshService.EXPECT().Generate(gomock.Any(), refresh.GenerateTokenInput{
					UserID: "fake-id",
					Role:   "customer",
				}).Return(refresh.GenerateTokenOutput{Token: "fake-refresh-token"}, nil)
			},
			want: customers.LoginCustomerOutput{
				TokenPair: customers.TokenPair{
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
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := customersmocks.NewMockRepository(ctrl)
			refreshService := refreshmocks.NewMockService(ctrl)
			jwtService := jwtmocks.NewMockService(ctrl)
			if tt.mocksSetup != nil {
				tt.mocksSetup(repo, refreshService, jwtService)
			}

			service := customers.NewService(logger, repo, refreshService, jwtService)
			got, err := service.LoginCustomer(context.Background(), tt.input)

			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestService_RefreshCustomer(t *testing.T) {
	logger, _ := log.NewTest()

	tests := []customersServiceTestCase[customers.RefreshCustomerInput, customers.RefreshCustomerOutput]{
		{
			name: "when there is not an active refresh token, then it should return an invalid refresh token error",
			input: customers.RefreshCustomerInput{
				RefreshToken: "InvalidRefreshToken",
				AccessToken:  "ValidAccessToken",
			},
			mocksSetup: func(_ *customersmocks.MockRepository, refreshService *refreshmocks.MockService,
				_ *jwtmocks.MockService) {

				refreshService.EXPECT().FindActiveToken(gomock.Any(), gomock.Any()).
					Return(refresh.FindActiveTokenOutput{}, refresh.ErrRefreshTokenNotFound)
			},
			want:    customers.RefreshCustomerOutput{},
			wantErr: customers.ErrInvalidRefreshToken,
		},
		{
			name: "when there is not an unexpected error when finding the active refresh token, " +
				"then it should propagate the error",
			input: customers.RefreshCustomerInput{
				RefreshToken: "InvalidRefreshToken",
				AccessToken:  "ValidAccessToken",
			},
			mocksSetup: func(_ *customersmocks.MockRepository, refreshService *refreshmocks.MockService,
				_ *jwtmocks.MockService) {

				refreshService.EXPECT().FindActiveToken(gomock.Any(), gomock.Any()).
					Return(refresh.FindActiveTokenOutput{}, errToken)
			},
			want:    customers.RefreshCustomerOutput{},
			wantErr: errToken,
		},
		{
			name: "when the expired access token is invalid, then it should return a token mismatch error",
			input: customers.RefreshCustomerInput{
				RefreshToken: "ValidRefreshToken",
				AccessToken:  "InvalidAccessToken",
			},
			mocksSetup: func(_ *customersmocks.MockRepository, refreshService *refreshmocks.MockService,
				jwtService *jwtmocks.MockService) {

				refreshService.EXPECT().FindActiveToken(gomock.Any(), gomock.Any()).
					Return(refresh.FindActiveTokenOutput{}, nil)

				jwtService.EXPECT().GetClaims(gomock.Any()).Return(jwt.Claims{}, jwt.ErrInvalidToken)
			},
			want:    customers.RefreshCustomerOutput{},
			wantErr: customers.ErrTokenMismatch,
		},
		{
			name: "when there is an unexpected error when reading access token claims, " +
				"then it should propagate the error",
			input: customers.RefreshCustomerInput{
				RefreshToken: "ValidRefreshToken",
				AccessToken:  "InvalidAccessToken",
			},
			mocksSetup: func(_ *customersmocks.MockRepository, refreshService *refreshmocks.MockService,
				jwtService *jwtmocks.MockService) {

				refreshService.EXPECT().FindActiveToken(gomock.Any(), gomock.Any()).
					Return(refresh.FindActiveTokenOutput{}, nil)

				jwtService.EXPECT().GetClaims(gomock.Any()).Return(jwt.Claims{}, errToken)
			},
			want:    customers.RefreshCustomerOutput{},
			wantErr: errToken,
		},
		{
			name: "when the user of the access token is different than the refresh, " +
				"then it should return a token mismatch error",
			input: customers.RefreshCustomerInput{
				RefreshToken: "ValidRefreshToken",
				AccessToken:  "ValidAccessToken",
			},
			mocksSetup: func(_ *customersmocks.MockRepository, refreshService *refreshmocks.MockService,
				jwtService *jwtmocks.MockService) {

				refreshService.EXPECT().FindActiveToken(gomock.Any(), gomock.Any()).
					Return(refresh.FindActiveTokenOutput{UserID: "fake-user-id-1"}, nil)

				//TODO As subject comes from the internal JWT library,
				//I don't want to expose it here. As a quick mitigation I'm doing it in that way,
				//but in long term I would like to do it in a better way.
				claims := jwt.Claims{}
				claims.Subject = "fake-user-id-2"
				jwtService.EXPECT().GetClaims(gomock.Any()).Return(claims, nil)
			},
			want:    customers.RefreshCustomerOutput{},
			wantErr: customers.ErrTokenMismatch,
		},
		{
			name: "when the role of the access token is different than the refresh, " +
				"then it should return a token mismatch error",
			input: customers.RefreshCustomerInput{
				RefreshToken: "ValidRefreshToken",
				AccessToken:  "ValidAccessToken",
			},
			mocksSetup: func(_ *customersmocks.MockRepository, refreshService *refreshmocks.MockService,
				jwtService *jwtmocks.MockService) {

				refreshService.EXPECT().FindActiveToken(gomock.Any(), gomock.Any()).
					Return(refresh.FindActiveTokenOutput{
						UserID: "fake-user-id-1",
						Role:   "role-1",
					}, nil)

				claims := jwt.Claims{Role: "role-2"}
				claims.Subject = "fake-user-id-1"
				jwtService.EXPECT().GetClaims(gomock.Any()).Return(claims, nil)
			},
			want:    customers.RefreshCustomerOutput{},
			wantErr: customers.ErrTokenMismatch,
		},
		{
			name: "when there is an error generating the new access token, then it should propagate the error",
			input: customers.RefreshCustomerInput{
				RefreshToken: "ValidRefreshToken",
				AccessToken:  "ValidAccessToken",
			},
			mocksSetup: func(_ *customersmocks.MockRepository, refreshService *refreshmocks.MockService,
				jwtService *jwtmocks.MockService) {

				refreshService.EXPECT().FindActiveToken(gomock.Any(), gomock.Any()).
					Return(refresh.FindActiveTokenOutput{
						UserID: "fake-user-id",
						Role:   "fake-role",
					}, nil)

				claims := jwt.Claims{Role: "fake-role"}
				claims.Subject = "fake-user-id"
				jwtService.EXPECT().GetClaims(gomock.Any()).Return(claims, nil)

				jwtService.EXPECT().GenerateToken(gomock.Any(), gomock.Any()).
					Return("", errToken)
			},
			want:    customers.RefreshCustomerOutput{},
			wantErr: errToken,
		},
		{
			name: "when there is an error generating the new refresh token, then it should propagate the error",
			input: customers.RefreshCustomerInput{
				RefreshToken: "ValidRefreshToken",
				AccessToken:  "ValidAccessToken",
			},
			mocksSetup: func(_ *customersmocks.MockRepository, refreshService *refreshmocks.MockService,
				jwtService *jwtmocks.MockService) {

				refreshService.EXPECT().FindActiveToken(gomock.Any(), gomock.Any()).
					Return(refresh.FindActiveTokenOutput{
						UserID: "fake-user-id",
						Role:   "fake-role",
					}, nil)

				claims := jwt.Claims{Role: "fake-role"}
				claims.Subject = "fake-user-id"
				jwtService.EXPECT().GetClaims(gomock.Any()).Return(claims, nil)

				jwtService.EXPECT().GenerateToken(gomock.Any(), gomock.Any()).
					Return("fake-token", nil)

				refreshService.EXPECT().Generate(gomock.Any(), gomock.Any()).
					Return(refresh.GenerateTokenOutput{}, errToken)
			},
			want:    customers.RefreshCustomerOutput{},
			wantErr: errToken,
		},
		{
			name: "when there is an error expiring the old refresh token, then it should propagate the error",
			input: customers.RefreshCustomerInput{
				RefreshToken: "ValidRefreshToken",
				AccessToken:  "ValidAccessToken",
			},
			mocksSetup: func(_ *customersmocks.MockRepository, refreshService *refreshmocks.MockService,
				jwtService *jwtmocks.MockService) {

				refreshService.EXPECT().FindActiveToken(gomock.Any(), gomock.Any()).
					Return(refresh.FindActiveTokenOutput{
						UserID: "fake-user-id",
						Role:   "fake-role",
					}, nil)

				claims := jwt.Claims{Role: "fake-role"}
				claims.Subject = "fake-user-id"
				jwtService.EXPECT().GetClaims(gomock.Any()).Return(claims, nil)

				jwtService.EXPECT().GenerateToken(gomock.Any(), gomock.Any()).
					Return("fake-token", nil)

				refreshService.EXPECT().Generate(gomock.Any(), gomock.Any()).
					Return(refresh.GenerateTokenOutput{Token: "fake-refresh-token"}, nil)

				refreshService.EXPECT().Expire(gomock.Any(), gomock.Any()).Return(refresh.ExpireOutput{}, errToken)
			},
			want:    customers.RefreshCustomerOutput{},
			wantErr: errToken,
		},
		{
			name: "when refresh token expiration returns a not found error, then it should return the new token",
			input: customers.RefreshCustomerInput{
				RefreshToken: "ValidRefreshToken",
				AccessToken:  "ValidAccessToken",
			},
			mocksSetup: func(_ *customersmocks.MockRepository, refreshService *refreshmocks.MockService,
				jwtService *jwtmocks.MockService) {

				refreshService.EXPECT().FindActiveToken(gomock.Any(), refresh.FindActiveTokenInput{
					Token: "ValidRefreshToken",
				}).Return(refresh.FindActiveTokenOutput{
					UserID: "fake-user-id",
					Role:   "customer",
				}, nil)

				claims := jwt.Claims{Role: "customer"}
				claims.Subject = "fake-user-id"
				jwtService.EXPECT().GetClaims("ValidAccessToken").Return(claims, nil)

				jwtService.EXPECT().GenerateToken("fake-user-id", jwt.Config{
					Expiration: 3600,
					Role:       "customer",
				}).Return("fake-token", nil)

				refreshService.EXPECT().Generate(gomock.Any(), refresh.GenerateTokenInput{
					UserID: "fake-user-id",
					Role:   "customer",
				}).Return(refresh.GenerateTokenOutput{Token: "fake-refresh-token"}, nil)

				refreshService.EXPECT().Expire(gomock.Any(), refresh.ExpireInput{
					Token: "ValidRefreshToken",
				}).Return(refresh.ExpireOutput{}, refresh.ErrRefreshTokenNotFound)
			},
			want: customers.RefreshCustomerOutput{
				TokenPair: customers.TokenPair{
					AccessToken:  "fake-token",
					RefreshToken: "fake-refresh-token",
					ExpiresIn:    3600,
					TokenType:    "Bearer",
				},
			},
		},
		{
			name: "when the new access token is generated correctly, then it should return the new token",
			input: customers.RefreshCustomerInput{
				RefreshToken: "ValidRefreshToken",
				AccessToken:  "ValidAccessToken",
			},
			mocksSetup: func(_ *customersmocks.MockRepository, refreshService *refreshmocks.MockService,
				jwtService *jwtmocks.MockService) {

				refreshService.EXPECT().FindActiveToken(gomock.Any(), refresh.FindActiveTokenInput{
					Token: "ValidRefreshToken",
				}).Return(refresh.FindActiveTokenOutput{
					UserID: "fake-user-id",
					Role:   "customer",
				}, nil)

				claims := jwt.Claims{Role: "customer"}
				claims.Subject = "fake-user-id"
				jwtService.EXPECT().GetClaims("ValidAccessToken").Return(claims, nil)

				jwtService.EXPECT().GenerateToken("fake-user-id", jwt.Config{
					Expiration: 3600,
					Role:       "customer",
				}).Return("fake-token", nil)

				refreshService.EXPECT().Generate(gomock.Any(), refresh.GenerateTokenInput{
					UserID: "fake-user-id",
					Role:   "customer",
				}).Return(refresh.GenerateTokenOutput{Token: "fake-refresh-token"}, nil)

				refreshService.EXPECT().Expire(gomock.Any(), refresh.ExpireInput{
					Token: "ValidRefreshToken",
				}).Return(refresh.ExpireOutput{}, nil)
			},
			want: customers.RefreshCustomerOutput{
				TokenPair: customers.TokenPair{
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
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := customersmocks.NewMockRepository(ctrl)
			refreshService := refreshmocks.NewMockService(ctrl)
			jwtService := jwtmocks.NewMockService(ctrl)
			if tt.mocksSetup != nil {
				tt.mocksSetup(repo, refreshService, jwtService)
			}

			service := customers.NewService(logger, repo, refreshService, jwtService)
			got, err := service.RefreshCustomer(context.Background(), tt.input)

			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}
