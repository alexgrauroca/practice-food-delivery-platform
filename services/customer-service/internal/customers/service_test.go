package customers_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/alexgrauroca/practice-food-delivery-platform/services/customer-service/internal/customers"
	customersmocks "github.com/alexgrauroca/practice-food-delivery-platform/services/customer-service/internal/customers/mocks"
	"github.com/alexgrauroca/practice-food-delivery-platform/services/customer-service/internal/log"
)

var errRepo = errors.New("repository error")

type customersServiceTestCase[I, W any] struct {
	name       string
	input      I
	mocksSetup func(repo *customersmocks.MockRepository)
	want       W
	wantErr    error
}

func TestService_RegisterCustomer(t *testing.T) {
	now := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	logger, _ := log.NewTest()

	tests := []customersServiceTestCase[customers.RegisterCustomerInput, customers.RegisterCustomerOutput]{
		{
			name: "when there is an active customer with the same email, " +
				"then it should return a customer already exists error",
			input: customers.RegisterCustomerInput{
				Email:       "test@example.com",
				Password:    "ValidPassword123",
				Name:        "John Doe",
				Address:     "a valid address",
				City:        "a valid city",
				PostalCode:  "12345",
				CountryCode: "US",
			},
			mocksSetup: func(repo *customersmocks.MockRepository) {
				repo.EXPECT().CreateCustomer(gomock.Any(), gomock.Any()).
					Return(customers.Customer{}, customers.ErrCustomerAlreadyExists)
			},
			want:    customers.RegisterCustomerOutput{},
			wantErr: customers.ErrCustomerAlreadyExists,
		},
		{
			name: "when there is an unexpected error when creating the customer, then it should propagate the error",
			input: customers.RegisterCustomerInput{
				Email:       "test@example.com",
				Password:    "ValidPassword123",
				Name:        "John Doe",
				Address:     "a valid address",
				City:        "a valid city",
				PostalCode:  "12345",
				CountryCode: "US",
			},
			mocksSetup: func(repo *customersmocks.MockRepository) {
				repo.EXPECT().CreateCustomer(gomock.Any(), gomock.Any()).
					Return(customers.Customer{}, errRepo)
			},
			want:    customers.RegisterCustomerOutput{},
			wantErr: errRepo,
		},
		{
			name: "when the customer can be created, then it should return the created customer",
			input: customers.RegisterCustomerInput{
				Email:       "test@example.com",
				Password:    "ValidPassword123",
				Name:        "John Doe",
				Address:     "a valid address",
				City:        "a valid city",
				PostalCode:  "12345",
				CountryCode: "US",
			},
			mocksSetup: func(repo *customersmocks.MockRepository) {
				repo.EXPECT().CreateCustomer(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, params customers.CreateCustomerParams) (customers.Customer, error) {
						return customers.Customer{
							ID:          "fake-id",
							Email:       params.Email,
							Name:        params.Name,
							Address:     params.Address,
							City:        params.City,
							PostalCode:  params.PostalCode,
							CountryCode: params.CountryCode,
							Active:      true,
							CreatedAt:   now,
							UpdatedAt:   now,
						}, nil
					})
			},
			want: customers.RegisterCustomerOutput{
				ID:          "fake-id",
				Email:       "test@example.com",
				Name:        "John Doe",
				Address:     "a valid address",
				City:        "a valid city",
				PostalCode:  "12345",
				CountryCode: "US",
				CreatedAt:   now,
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := customersmocks.NewMockRepository(ctrl)
			if tt.mocksSetup != nil {
				tt.mocksSetup(repo)
			}

			service := customers.NewService(logger, repo)
			got, err := service.RegisterCustomer(context.Background(), tt.input)

			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}
