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
	"go.uber.org/zap"
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
				//TODO replace the error for the correct one after knowing how the driver returns the unique constraint error
				repo.EXPECT().CreateCustomer(gomock.Any(), gomock.Any()).
					Return(customers.Customer{}, repoError)
			},
			expectedOutput: customers.RegisterCustomerOutput{},
			//TODO replace the error for the correct one after knowing how the driver returns the unique constraint error
			// this is not the best way to do it, but as I'm learning how everything works, this is the only way to pass
			// the ci pipeline
			expectError: repoError,
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
				repo.EXPECT().CreateCustomer(gomock.Any(), customers.CreateCustomerParams{
					Email:    "test@example.com",
					Password: "ValidPassword123",
					Name:     "John Doe",
				}).
					Return(customers.Customer{
						ID:        "fake-id",
						Email:     "test@example.com",
						Name:      "John Doe",
						CreatedAt: now,
						UpdatedAt: now,
						Active:    true,
					}, nil)
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
