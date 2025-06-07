package customers_test

import "testing"

func TestRepository_CreateCustomer(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "when exists an active customer with the same email, it should return a customer already exists error",
		},
		{
			name: "when there is an error creating the customer, it should propagate the error",
		},
		{
			name: "when the customer is created successfully, it should return the created customer",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			
		})
	}
}
