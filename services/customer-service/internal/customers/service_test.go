//go:build unit

package customers_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/alexgrauroca/practice-food-delivery-platform/services/customer-service/internal/authentication"
	authmocks "github.com/alexgrauroca/practice-food-delivery-platform/services/customer-service/internal/authentication/mocks"
	"github.com/alexgrauroca/practice-food-delivery-platform/services/customer-service/internal/customers"
	customersmocks "github.com/alexgrauroca/practice-food-delivery-platform/services/customer-service/internal/customers/mocks"
	"github.com/alexgrauroca/practice-food-delivery-platform/services/customer-service/internal/log"
)

var (
	errRepo    = errors.New("repository error")
	errAuthCli = errors.New("authentication client error")
)

type customersServiceTestCase[I, W any] struct {
	name       string
	input      I
	mocksSetup func(repo *customersmocks.MockRepository, authcli *authmocks.MockClient, authctx *authmocks.MockContextReader)
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
			mocksSetup: func(
				repo *customersmocks.MockRepository,
				_ *authmocks.MockClient,
				_ *authmocks.MockContextReader,
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
				Email:       "test@example.com",
				Password:    "ValidPassword123",
				Name:        "John Doe",
				Address:     "a valid address",
				City:        "a valid city",
				PostalCode:  "12345",
				CountryCode: "US",
			},
			mocksSetup: func(
				repo *customersmocks.MockRepository,
				_ *authmocks.MockClient,
				_ *authmocks.MockContextReader,
			) {
				repo.EXPECT().CreateCustomer(gomock.Any(), gomock.Any()).
					Return(customers.Customer{}, errRepo)
			},
			want:    customers.RegisterCustomerOutput{},
			wantErr: errRepo,
		},
		{
			name: "when there is a customer not found error when purging the created customer, " +
				"then it should propagate the authcli error",
			input: customers.RegisterCustomerInput{
				Email:       "test@example.com",
				Password:    "ValidPassword123",
				Name:        "John Doe",
				Address:     "a valid address",
				City:        "a valid city",
				PostalCode:  "12345",
				CountryCode: "US",
			},
			mocksSetup: func(
				repo *customersmocks.MockRepository,
				authcli *authmocks.MockClient,
				_ *authmocks.MockContextReader,
			) {
				// The returned customer is not relevant for this case
				repo.EXPECT().CreateCustomer(gomock.Any(), gomock.Any()).
					Return(customers.Customer{}, nil)

				authcli.EXPECT().RegisterCustomer(gomock.Any(), gomock.Any()).
					Return(authentication.RegisterCustomerResponse{}, errAuthCli)

				repo.EXPECT().PurgeCustomer(gomock.Any(), gomock.Any()).Return(customers.ErrCustomerNotFound)
			},
			want:    customers.RegisterCustomerOutput{},
			wantErr: errAuthCli,
		},
		{
			name: "when there is an unexpected error when purging the created customer, " +
				"then it should propagate the error",
			input: customers.RegisterCustomerInput{
				Email:       "test@example.com",
				Password:    "ValidPassword123",
				Name:        "John Doe",
				Address:     "a valid address",
				City:        "a valid city",
				PostalCode:  "12345",
				CountryCode: "US",
			},
			mocksSetup: func(
				repo *customersmocks.MockRepository,
				authcli *authmocks.MockClient,
				_ *authmocks.MockContextReader,
			) {
				// The returned customer is not relevant for this case
				repo.EXPECT().CreateCustomer(gomock.Any(), gomock.Any()).
					Return(customers.Customer{}, nil)

				authcli.EXPECT().RegisterCustomer(gomock.Any(), gomock.Any()).
					Return(authentication.RegisterCustomerResponse{}, errAuthCli)

				repo.EXPECT().PurgeCustomer(gomock.Any(), gomock.Any()).Return(errRepo)
			},
			want:    customers.RegisterCustomerOutput{},
			wantErr: errRepo,
		},
		{
			name: "when there is an unexpected error when registering the customer at auth service, " +
				"then it should propagate the error",
			input: customers.RegisterCustomerInput{
				Email:       "test@example.com",
				Password:    "ValidPassword123",
				Name:        "John Doe",
				Address:     "a valid address",
				City:        "a valid city",
				PostalCode:  "12345",
				CountryCode: "US",
			},
			mocksSetup: func(
				repo *customersmocks.MockRepository,
				authcli *authmocks.MockClient,
				_ *authmocks.MockContextReader,
			) {
				// The returned customer is not relevant for this case
				repo.EXPECT().CreateCustomer(gomock.Any(), gomock.Any()).
					Return(customers.Customer{}, nil)

				authcli.EXPECT().RegisterCustomer(gomock.Any(), gomock.Any()).
					Return(authentication.RegisterCustomerResponse{}, errAuthCli)

				repo.EXPECT().PurgeCustomer(gomock.Any(), "test@example.com").Return(nil)
			},
			want:    customers.RegisterCustomerOutput{},
			wantErr: errAuthCli,
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
			mocksSetup: func(
				repo *customersmocks.MockRepository,
				authcli *authmocks.MockClient,
				_ *authmocks.MockContextReader,
			) {
				repo.EXPECT().CreateCustomer(gomock.Any(), customers.CreateCustomerParams{
					Email:       "test@example.com",
					Name:        "John Doe",
					Address:     "a valid address",
					City:        "a valid city",
					PostalCode:  "12345",
					CountryCode: "US",
				}).DoAndReturn(func(_ context.Context, params customers.CreateCustomerParams) (customers.Customer, error) {
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

				authcli.EXPECT().RegisterCustomer(gomock.Any(), authentication.RegisterCustomerRequest{
					CustomerID: "fake-id",
					Email:      "test@example.com",
					Password:   "ValidPassword123",
					Name:       "John Doe",
				}).Return(authentication.RegisterCustomerResponse{
					ID:        "auth-fake-id",
					Email:     "test@example.com",
					Name:      "John Doe",
					CreatedAt: now,
				}, nil)
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
			authcli := authmocks.NewMockClient(ctrl)
			authctx := authmocks.NewMockContextReader(ctrl)

			if tt.mocksSetup != nil {
				tt.mocksSetup(repo, authcli, authctx)
			}

			service := customers.NewService(logger, repo, authcli, authctx)
			got, err := service.RegisterCustomer(context.Background(), tt.input)

			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestService_GetCustomer(t *testing.T) {
	now := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	logger, _ := log.NewTest()

	tests := []customersServiceTestCase[customers.GetCustomerInput, customers.GetCustomerOutput]{
		{
			name: "when there is an unexpected error when getting the auth subject, " +
				"then it should return an invalid token error",
			input: customers.GetCustomerInput{},
			mocksSetup: func(
				_ *customersmocks.MockRepository,
				_ *authmocks.MockClient,
				authctx *authmocks.MockContextReader,
			) {
				authctx.EXPECT().GetSubject(gomock.Any()).Return("", false)
			},
			want:    customers.GetCustomerOutput{},
			wantErr: authentication.ErrInvalidToken,
		},
		{
			name: "when the authenticated customer id is different than the requested customer id, " +
				"then it should return a customer id mismatch error",
			input: customers.GetCustomerInput{CustomerID: "fake-id"},
			mocksSetup: func(
				_ *customersmocks.MockRepository,
				_ *authmocks.MockClient,
				authctx *authmocks.MockContextReader,
			) {
				authctx.EXPECT().GetSubject(gomock.Any()).Return("another-fake-id", true)
			},
			want:    customers.GetCustomerOutput{},
			wantErr: customers.ErrCustomerIDMismatch,
		},
		{
			name:  "when there is an unexpected error when getting the customer, then it should propagate the error",
			input: customers.GetCustomerInput{CustomerID: "fake-id"},
			mocksSetup: func(
				repo *customersmocks.MockRepository,
				_ *authmocks.MockClient,
				authctx *authmocks.MockContextReader,
			) {
				authctx.EXPECT().GetSubject(gomock.Any()).Return("fake-id", true)

				repo.EXPECT().GetCustomer(gomock.Any(), gomock.Any()).Return(customers.Customer{}, errRepo)
			},
			want:    customers.GetCustomerOutput{},
			wantErr: errRepo,
		},
		{
			name:  "when the customer is not found, then it should return a customer not found error",
			input: customers.GetCustomerInput{CustomerID: "fake-id"},
			mocksSetup: func(
				repo *customersmocks.MockRepository,
				_ *authmocks.MockClient,
				authctx *authmocks.MockContextReader,
			) {
				authctx.EXPECT().GetSubject(gomock.Any()).Return("fake-id", true)

				repo.EXPECT().GetCustomer(gomock.Any(), gomock.Any()).
					Return(customers.Customer{}, customers.ErrCustomerNotFound)
			},
			want:    customers.GetCustomerOutput{},
			wantErr: customers.ErrCustomerNotFound,
		},
		{
			name:  "when the customer is found, then it should return the customer data",
			input: customers.GetCustomerInput{CustomerID: "fake-id"},
			mocksSetup: func(
				repo *customersmocks.MockRepository,
				_ *authmocks.MockClient,
				authctx *authmocks.MockContextReader,
			) {
				authctx.EXPECT().GetSubject(gomock.Any()).Return("fake-id", true)

				repo.EXPECT().GetCustomer(gomock.Any(), "fake-id").
					Return(customers.Customer{
						ID:          "fake-id",
						Email:       "test@example.com",
						Name:        "John Doe",
						Active:      true,
						Address:     "123 Main St",
						City:        "New York",
						PostalCode:  "10001",
						CountryCode: "US",
						CreatedAt:   now,
						UpdatedAt:   now,
					}, nil)
			},
			want: customers.GetCustomerOutput{
				ID:          "fake-id",
				Email:       "test@example.com",
				Name:        "John Doe",
				Address:     "123 Main St",
				City:        "New York",
				PostalCode:  "10001",
				CountryCode: "US",
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := customersmocks.NewMockRepository(ctrl)
			authcli := authmocks.NewMockClient(ctrl)
			authctx := authmocks.NewMockContextReader(ctrl)

			if tt.mocksSetup != nil {
				tt.mocksSetup(repo, authcli, authctx)
			}

			service := customers.NewService(logger, repo, authcli, authctx)
			got, err := service.GetCustomer(context.Background(), tt.input)

			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}
