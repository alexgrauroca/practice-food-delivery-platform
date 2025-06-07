package customers_test

import (
	"testing"

	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/customers"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestService_RegisterCustomer(t *testing.T) {
	logger := zap.NewNop()

	tests := []struct {
		name           string
		input          customers.RegisterCustomerInput
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
			expectedOutput: customers.RegisterCustomerOutput{},
			expectError:    customers.ErrCustomerAlreadyExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := customers.NewService(logger)
			output, err := service.RegisterCustomer(tt.input)

			assert.Equal(t, tt.expectedOutput, output)
			assert.Equal(t, tt.expectError, err)
		})
	}
}
