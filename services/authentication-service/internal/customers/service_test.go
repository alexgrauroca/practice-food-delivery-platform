//go:build !integration

package customers_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/customers"
	customersmocks "github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/customers/mocks"
	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/password"
	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/refresh"
	refreshmocks "github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/refresh/mocks"
)

var (
	errRepo  = errors.New("repository error")
	errToken = errors.New("token error")
	logger   = zap.NewNop()
)

func TestService_RegisterCustomer(t *testing.T) {
	now := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	tests := []struct {
		name           string
		input          customers.RegisterCustomerInput
		mocksSetup     func(repo *customersmocks.MockRepository)
		expectedOutput customers.RegisterCustomerOutput
		expectError    error
	}{
		{
			name: "when there is an active customer with the same email, then it should return a customer already exists error",
			input: customers.RegisterCustomerInput{
				Email:    "test@example.com",
				Password: "ValidPassword123",
				Name:     "John Doe",
			},
			mocksSetup: func(repo *customersmocks.MockRepository) {
				repo.EXPECT().CreateCustomer(gomock.Any(), gomock.Any()).
					Return(customers.Customer{}, customers.ErrCustomerAlreadyExists)
			},
			expectedOutput: customers.RegisterCustomerOutput{},
			expectError:    customers.ErrCustomerAlreadyExists,
		},
		{
			name: "when there is an unexpected error when creating the customer, then it should propagate the error",
			input: customers.RegisterCustomerInput{
				Email:    "test@example.com",
				Password: "ValidPassword123",
				Name:     "John Doe",
			},
			mocksSetup: func(repo *customersmocks.MockRepository) {
				repo.EXPECT().CreateCustomer(gomock.Any(), gomock.Any()).
					Return(customers.Customer{}, errRepo)
			},
			expectedOutput: customers.RegisterCustomerOutput{},
			expectError:    errRepo,
		},
		{
			name: "when the customer can be created, then it should return the created customer",
			input: customers.RegisterCustomerInput{
				Email:    "test@example.com",
				Password: "ValidPassword123",
				Name:     "John Doe",
			},
			mocksSetup: func(repo *customersmocks.MockRepository) {
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
			expectedOutput: customers.RegisterCustomerOutput{
				ID:        "fake-id",
				Email:     "test@example.com",
				Name:      "John Doe",
				CreatedAt: now,
			},
			expectError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := customersmocks.NewMockRepository(ctrl)
			refreshService := refreshmocks.NewMockService(ctrl)
			if tt.mocksSetup != nil {
				tt.mocksSetup(repo)
			}

			service := customers.NewService(logger, repo, refreshService)
			output, err := service.RegisterCustomer(context.Background(), tt.input)

			assert.ErrorIs(t, err, tt.expectError)
			assert.Equal(t, tt.expectedOutput, output)
		})
	}
}

func TestService_LoginCustomer(t *testing.T) {
	now := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name           string
		input          customers.LoginCustomerInput
		mocksSetup     func(repo *customersmocks.MockRepository, refreshService *refreshmocks.MockService)
		expectedOutput customers.LoginCustomerOutput
		expectError    error
	}{
		{
			name: "when there is not an active customer with the same email, then it should an invalid credentials error",
			input: customers.LoginCustomerInput{
				Email:    "test@example.com",
				Password: "ValidPassword123",
			},
			mocksSetup: func(repo *customersmocks.MockRepository, refreshService *refreshmocks.MockService) {
				repo.EXPECT().FindByEmail(gomock.Any(), gomock.Any()).
					Return(customers.Customer{}, customers.ErrCustomerNotFound)
			},
			expectedOutput: customers.LoginCustomerOutput{},
			expectError:    customers.ErrInvalidCredentials,
		},
		{
			name: "when there is not an active customer with the same password, then it should an invalid credentials error",
			input: customers.LoginCustomerInput{
				Email:    "test@example.com",
				Password: "InvalidPassword123",
			},
			mocksSetup: func(repo *customersmocks.MockRepository, refreshService *refreshmocks.MockService) {
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
			expectedOutput: customers.LoginCustomerOutput{},
			expectError:    customers.ErrInvalidCredentials,
		},
		{
			name: "when there is an unexpected error when fetching the customer, then it should propagate the error",
			input: customers.LoginCustomerInput{
				Email:    "test@example.com",
				Password: "ValidPassword123",
			},
			mocksSetup: func(repo *customersmocks.MockRepository, refreshService *refreshmocks.MockService) {
				repo.EXPECT().FindByEmail(gomock.Any(), gomock.Any()).
					Return(customers.Customer{}, errRepo)
			},
			expectedOutput: customers.LoginCustomerOutput{},
			expectError:    errRepo,
		},
		{
			name: "when there is an error generating the refresh token, then it should propagate the error",
			input: customers.LoginCustomerInput{
				Email:    "test@example.com",
				Password: "ValidPassword123",
			},
			mocksSetup: func(repo *customersmocks.MockRepository, refreshService *refreshmocks.MockService) {
				repo.EXPECT().FindByEmail(gomock.Any(), gomock.Any()).
					Return(customers.Customer{}, nil)

				refreshService.EXPECT().Generate(gomock.Any(), gomock.Any()).
					Return("", errToken)
			},
			expectedOutput: customers.LoginCustomerOutput{},
			expectError:    errToken,
		},
		{
			name: "when there is an active customer with the same email and password, then it should return its token",
			input: customers.LoginCustomerInput{
				Email:    "test@example.com",
				Password: "ValidPassword123",
			},
			mocksSetup: func(repo *customersmocks.MockRepository, refreshService *refreshmocks.MockService) {
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

				refreshService.EXPECT().Generate(gomock.Any(), refresh.GenerateTokenInput{
					UserID: "fake-id",
					Role:   "customer",
				}).Return("fake-refresh-token", nil)
			},
			expectedOutput: customers.LoginCustomerOutput{
				ExpiresIn:    3600, // 1 hour
				TokenType:    "Bearer",
				RefreshToken: "fake-refresh-token",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := customersmocks.NewMockRepository(ctrl)
			refreshService := refreshmocks.NewMockService(ctrl)
			if tt.mocksSetup != nil {
				tt.mocksSetup(repo, refreshService)
			}

			service := customers.NewService(logger, repo, refreshService)
			output, err := service.LoginCustomer(context.Background(), tt.input)

			assert.ErrorIs(t, err, tt.expectError)

			// We only assert the expectedOutput if there is any error
			if tt.expectError == nil {
				assert.Equal(t, tt.expectedOutput.TokenType, output.TokenType)
				assert.Equal(t, tt.expectedOutput.ExpiresIn, output.ExpiresIn)
				assert.Equal(t, tt.expectedOutput.RefreshToken, output.RefreshToken)

				// As tokens are generated depending on the moment of the time, we just need to check if the token
				// is not empty
				assert.NotEmpty(t, output.AccessToken)
			}
		})
	}
}
