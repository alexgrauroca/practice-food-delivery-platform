package customers_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/customers"
	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/customers/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func TestService_RegisterCustomer(t *testing.T) {
	logger := zap.NewNop()
	now := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	repoError := errors.New("repository error")

	tests := []struct {
		name           string
		input          customers.RegisterCustomerInput
		mocksSetup     func(repo *mocks.MockRepository)
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
			mocksSetup: func(repo *mocks.MockRepository) {
				repo.EXPECT().CreateCustomer(gomock.Any(), gomock.Any()).
					Return(customers.Customer{}, customers.ErrCustomerAlreadyExists)
			},
			expectedOutput: customers.RegisterCustomerOutput{},
			expectError:    customers.ErrCustomerAlreadyExists,
		},
		{
			name: "when there is an error when creating the customer, then it should propagate the error",
			input: customers.RegisterCustomerInput{
				Email:    "test@example.com",
				Password: "ValidPassword123",
				Name:     "John Doe",
			},
			mocksSetup: func(repo *mocks.MockRepository) {
				repo.EXPECT().CreateCustomer(gomock.Any(), gomock.Any()).
					Return(customers.Customer{}, repoError)
			},
			expectedOutput: customers.RegisterCustomerOutput{},
			expectError:    repoError,
		},
		{
			name: "when the customer can be created, then it should return the created customer",
			input: customers.RegisterCustomerInput{
				Email:    "test@example.com",
				Password: "ValidPassword123",
				Name:     "John Doe",
			},
			mocksSetup: func(repo *mocks.MockRepository) {
				repo.EXPECT().CreateCustomer(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, params customers.CreateCustomerParams) (customers.Customer, error) {
						// Assert that the password is hashed
						err := bcrypt.CompareHashAndPassword([]byte(params.Password), []byte("ValidPassword123"))
						require.NoError(t, err, "Password should be hashed and match the input password")

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
			repo := mocks.NewMockRepository(gomock.NewController(t))
			if tt.mocksSetup != nil {
				tt.mocksSetup(repo)
			}

			service := customers.NewService(logger, repo)
			output, err := service.RegisterCustomer(context.Background(), tt.input)

			assert.Equal(t, tt.expectedOutput, output)
			assert.Equal(t, tt.expectError, err)
		})
	}
}
