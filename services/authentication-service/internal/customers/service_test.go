//go:build unit

package customers_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/authcore"
	authcoremocks "github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/authcore/mocks"
	customersmocks "github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/customers/mocks"

	"github.com/alexgrauroca/practice-food-delivery-platform/pkg/log"
	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/customers"
	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/password"
)

var (
	errRepo  = errors.New("repository error")
	errToken = errors.New("token error")
)

type customersServiceTestCase[I, W any] struct {
	name       string
	input      I
	want       W
	mocksSetup func(
		repo *customersmocks.MockRepository,
		authCoreService *authcoremocks.MockService,
	)
	wantErr error
}

func TestService_RegisterCustomer(t *testing.T) {
	now := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	logger, _ := log.NewTest()

	tests := []customersServiceTestCase[customers.RegisterCustomerInput, customers.RegisterCustomerOutput]{
		{
			name: "when there is an active customer with the same email, then it should return a customer already exists error",
			input: customers.RegisterCustomerInput{
				CustomerID: "fake-customer-id",
				Email:      "test@example.com",
				Password:   "ValidPassword123",
			},
			mocksSetup: func(
				repo *customersmocks.MockRepository,
				_ *authcoremocks.MockService,
			) {
				repo.EXPECT().CreateCustomer(gomock.Any(), gomock.Any()).
					Return(customers.Customer{}, customers.ErrCustomerAlreadyExists)
			},
			want:    customers.RegisterCustomerOutput{},
			wantErr: customers.ErrCustomerAlreadyExists,
		},
		{
			name: "when there is an unexpected error when creating the customer, then it should propagate the error",
			input: customers.RegisterCustomerInput{
				CustomerID: "fake-customer-id",
				Email:      "test@example.com",
				Password:   "ValidPassword123",
			},
			mocksSetup: func(
				repo *customersmocks.MockRepository,
				_ *authcoremocks.MockService,
			) {
				repo.EXPECT().CreateCustomer(gomock.Any(), gomock.Any()).
					Return(customers.Customer{}, errRepo)
			},
			want:    customers.RegisterCustomerOutput{},
			wantErr: errRepo,
		},
		{
			name: "when the customer can be created, then it should return the created customer",
			input: customers.RegisterCustomerInput{
				CustomerID: "fake-customer-id",
				Email:      "test@example.com",
				Password:   "ValidPassword123",
			},
			mocksSetup: func(
				repo *customersmocks.MockRepository,
				_ *authcoremocks.MockService,
			) {
				repo.EXPECT().CreateCustomer(gomock.Any(), gomock.Any()).
					DoAndReturn(
						func(_ context.Context, params customers.CreateCustomerParams) (customers.Customer, error) {
							// Assert that the password is hashed
							ok := password.Verify(params.Password, "ValidPassword123")
							require.True(t, ok, "Password should be hashed and match the input password")

							return customers.Customer{
								ID:         "fake-id",
								CustomerID: params.CustomerID,
								Email:      params.Email,
								Password:   params.Password,
								CreatedAt:  now,
								UpdatedAt:  now,
								Active:     true,
							}, nil
						},
					)
			},
			want: customers.RegisterCustomerOutput{
				ID:        "fake-id",
				Email:     "test@example.com",
				CreatedAt: now,
				UpdatedAt: now,
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				service, cleanup := serviceSetup(t, logger, tt.mocksSetup)
				defer cleanup()

				got, err := service.RegisterCustomer(context.Background(), tt.input)

				assert.ErrorIs(t, err, tt.wantErr)
				assert.Equal(t, tt.want, got)
			},
		)
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
			mocksSetup: func(
				repo *customersmocks.MockRepository,
				_ *authcoremocks.MockService,
			) {
				repo.EXPECT().FindByEmail(gomock.Any(), gomock.Any()).
					Return(customers.Customer{}, customers.ErrCustomerNotFound)
			},
			want:    customers.LoginCustomerOutput{},
			wantErr: authcore.ErrInvalidCredentials,
		},
		{
			name: "when there is not an active customer with the same password, " +
				"then it should return an invalid credentials error",
			input: customers.LoginCustomerInput{
				Email:    "test@example.com",
				Password: "InvalidPassword123",
			},
			mocksSetup: func(
				repo *customersmocks.MockRepository,
				_ *authcoremocks.MockService,
			) {
				hashedPassword, err := password.Hash("ValidPassword123")
				require.NoError(t, err)

				repo.EXPECT().FindByEmail(gomock.Any(), gomock.Any()).
					Return(
						customers.Customer{
							ID:        "fake-id",
							Email:     "test@example.com",
							Password:  hashedPassword, // This should be a hashed password
							CreatedAt: now,
							UpdatedAt: now,
							Active:    true,
						}, nil,
					)
			},
			want:    customers.LoginCustomerOutput{},
			wantErr: authcore.ErrInvalidCredentials,
		},
		{
			name: "when there is an unexpected error when fetching the customer, then it should propagate the error",
			input: customers.LoginCustomerInput{
				Email:    "test@example.com",
				Password: "ValidPassword123",
			},
			mocksSetup: func(
				repo *customersmocks.MockRepository,
				_ *authcoremocks.MockService,
			) {
				repo.EXPECT().FindByEmail(gomock.Any(), gomock.Any()).
					Return(customers.Customer{}, errRepo)
			},
			want:    customers.LoginCustomerOutput{},
			wantErr: errRepo,
		},
		{
			name: "when there is an error generating the token pair, then it should propagate the error",
			input: customers.LoginCustomerInput{
				Email:    "test@example.com",
				Password: "ValidPassword123",
			},
			mocksSetup: func(
				repo *customersmocks.MockRepository,
				authCoreService *authcoremocks.MockService,
			) {
				hashedPassword, err := password.Hash("ValidPassword123")
				require.NoError(t, err)

				repo.EXPECT().FindByEmail(gomock.Any(), "test@example.com").
					Return(
						customers.Customer{
							ID:        "fake-id",
							Email:     "test@example.com",
							Password:  hashedPassword, // This should be a hashed password
							CreatedAt: now,
							UpdatedAt: now,
							Active:    true,
						}, nil,
					)

				authCoreService.EXPECT().GenerateTokenPair(gomock.Any(), gomock.Any()).
					Return(authcore.TokenPair{}, errToken)
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
			mocksSetup: func(
				repo *customersmocks.MockRepository,
				authCoreService *authcoremocks.MockService,
			) {
				hashedPassword, err := password.Hash("ValidPassword123")
				require.NoError(t, err)

				repo.EXPECT().FindByEmail(gomock.Any(), "test@example.com").
					Return(
						customers.Customer{
							ID:         "fake-customer-id",
							CustomerID: "fake-id",
							Email:      "test@example.com",
							Password:   hashedPassword, // This should be a hashed password
							CreatedAt:  now,
							UpdatedAt:  now,
							Active:     true,
						}, nil,
					)

				authCoreService.EXPECT().GenerateTokenPair(gomock.Any(), gomock.Any()).
					Return(
						authcore.TokenPair{
							AccessToken:  "fake-token",
							RefreshToken: "fake-refresh-token",
							ExpiresIn:    3600,
							TokenType:    "Bearer",
						}, nil,
					)
			},
			want: customers.LoginCustomerOutput{
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
		t.Run(
			tt.name, func(t *testing.T) {
				service, cleanup := serviceSetup(t, logger, tt.mocksSetup)
				defer cleanup()

				got, err := service.LoginCustomer(context.Background(), tt.input)

				assert.ErrorIs(t, err, tt.wantErr)
				assert.Equal(t, tt.want, got)
			},
		)
	}
}

func TestService_RefreshCustomer(t *testing.T) {
	logger, _ := log.NewTest()

	tests := []customersServiceTestCase[customers.RefreshCustomerInput, customers.RefreshCustomerOutput]{
		{
			name: "when there is an error refreshing the token, then it should propagate the error",
			input: customers.RefreshCustomerInput{
				RefreshToken: "InvalidRefreshToken",
				AccessToken:  "ValidAccessToken",
			},
			mocksSetup: func(
				_ *customersmocks.MockRepository,
				authCoreService *authcoremocks.MockService,
			) {
				authCoreService.EXPECT().RefreshToken(gomock.Any(), gomock.Any()).
					Return(authcore.TokenPair{}, errToken)
			},
			want:    customers.RefreshCustomerOutput{},
			wantErr: errToken,
		},
		{
			name: "when the new access token is generated correctly, then it should return the new token",
			input: customers.RefreshCustomerInput{
				RefreshToken: "ValidRefreshToken",
				AccessToken:  "ValidAccessToken",
			},
			mocksSetup: func(
				_ *customersmocks.MockRepository,
				authCoreService *authcoremocks.MockService,
			) {
				authCoreService.EXPECT().RefreshToken(gomock.Any(), gomock.Any()).
					Return(
						authcore.TokenPair{
							AccessToken:  "fake-token",
							RefreshToken: "fake-refresh-token",
							ExpiresIn:    3600,
							TokenType:    "Bearer",
						}, nil,
					)
			},
			want: customers.RefreshCustomerOutput{
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
		t.Run(
			tt.name, func(t *testing.T) {
				service, cleanup := serviceSetup(t, logger, tt.mocksSetup)
				defer cleanup()

				got, err := service.RefreshCustomer(context.Background(), tt.input)

				assert.ErrorIs(t, err, tt.wantErr)
				assert.Equal(t, tt.want, got)
			},
		)
	}
}

func serviceSetup(
	t *testing.T, logger log.Logger, mocksSetup func(
		repo *customersmocks.MockRepository,
		authCoreService *authcoremocks.MockService,
	),
) (customers.Service, func()) {
	ctrl := gomock.NewController(t)

	repo := customersmocks.NewMockRepository(ctrl)
	authCoreService := authcoremocks.NewMockService(ctrl)

	if mocksSetup != nil {
		mocksSetup(repo, authCoreService)
	}

	service := customers.NewService(logger, repo, authCoreService)
	return service, func() {
		ctrl.Finish()
	}
}
