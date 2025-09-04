//go:build unit

package customers_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/alexgrauroca/practice-food-delivery-platform/pkg/clients/authentication"
	authmocks "github.com/alexgrauroca/practice-food-delivery-platform/pkg/clients/authentication/mocks"
	"github.com/alexgrauroca/practice-food-delivery-platform/pkg/log"
	"github.com/alexgrauroca/practice-food-delivery-platform/services/customer-service/internal/customers"
	customersmocks "github.com/alexgrauroca/practice-food-delivery-platform/services/customer-service/internal/customers/mocks"
)

var (
	errRepo        = errors.New("repository error")
	errAuthService = errors.New("authentication service error")
)

type customersServiceTestCase[I, W any] struct {
	name       string
	input      I
	mocksSetup func(
		repo *customersmocks.MockRepository,
		authservice *authmocks.MockService,
		authctx *authmocks.MockContextReader,
	)
	want    W
	wantErr error
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
				_ *authmocks.MockService,
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
				_ *authmocks.MockService,
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
				"then it should propagate the authservice error",
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
				authservice *authmocks.MockService,
				_ *authmocks.MockContextReader,
			) {
				// The returned customer is not relevant for this case
				repo.EXPECT().CreateCustomer(gomock.Any(), gomock.Any()).
					Return(customers.Customer{}, nil)

				authservice.EXPECT().RegisterCustomer(gomock.Any(), gomock.Any()).
					Return(authentication.RegisterCustomerOutput{}, errAuthService)

				repo.EXPECT().PurgeCustomer(gomock.Any(), gomock.Any()).Return(customers.ErrCustomerNotFound)
			},
			want:    customers.RegisterCustomerOutput{},
			wantErr: errAuthService,
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
				authservice *authmocks.MockService,
				_ *authmocks.MockContextReader,
			) {
				// The returned customer is not relevant for this case
				repo.EXPECT().CreateCustomer(gomock.Any(), gomock.Any()).
					Return(customers.Customer{}, nil)

				authservice.EXPECT().RegisterCustomer(gomock.Any(), gomock.Any()).
					Return(authentication.RegisterCustomerOutput{}, errAuthService)

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
				authservice *authmocks.MockService,
				_ *authmocks.MockContextReader,
			) {
				// The returned customer is not relevant for this case
				repo.EXPECT().CreateCustomer(gomock.Any(), gomock.Any()).
					Return(customers.Customer{}, nil)

				authservice.EXPECT().RegisterCustomer(gomock.Any(), gomock.Any()).
					Return(authentication.RegisterCustomerOutput{}, errAuthService)

				repo.EXPECT().PurgeCustomer(gomock.Any(), "test@example.com").Return(nil)
			},
			want:    customers.RegisterCustomerOutput{},
			wantErr: errAuthService,
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
				authservice *authmocks.MockService,
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

				authservice.EXPECT().RegisterCustomer(gomock.Any(), authentication.RegisterCustomerInput{
					CustomerID: "fake-id",
					Email:      "test@example.com",
					Password:   "ValidPassword123",
					Name:       "John Doe",
				}).Return(authentication.RegisterCustomerOutput{
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
			authservice := authmocks.NewMockService(ctrl)
			authctx := authmocks.NewMockContextReader(ctrl)

			if tt.mocksSetup != nil {
				tt.mocksSetup(repo, authservice, authctx)
			}

			service := customers.NewService(logger, repo, authservice, authctx)
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
				_ *authmocks.MockService,
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
				_ *authmocks.MockService,
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
				_ *authmocks.MockService,
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
				_ *authmocks.MockService,
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
				_ *authmocks.MockService,
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
			authservice := authmocks.NewMockService(ctrl)
			authctx := authmocks.NewMockContextReader(ctrl)

			if tt.mocksSetup != nil {
				tt.mocksSetup(repo, authservice, authctx)
			}

			service := customers.NewService(logger, repo, authservice, authctx)
			got, err := service.GetCustomer(context.Background(), tt.input)

			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestService_UpdateCustomer(t *testing.T) {
	now := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	yesterday := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)
	logger, _ := log.NewTest()

	tests := []customersServiceTestCase[customers.UpdateCustomerInput, customers.UpdateCustomerOutput]{
		{
			name: "when there is an unexpected error when getting the auth subject, " +
				"then it should return an invalid token error",
			input: customers.UpdateCustomerInput{},
			mocksSetup: func(
				_ *customersmocks.MockRepository,
				_ *authmocks.MockService,
				authctx *authmocks.MockContextReader,
			) {
				authctx.EXPECT().GetSubject(gomock.Any()).Return("", false)
			},
			want:    customers.UpdateCustomerOutput{},
			wantErr: authentication.ErrInvalidToken,
		},
		{
			name: "when the authenticated customer id is different than the requested customer id, " +
				"then it should return a customer id mismatch error",
			input: customers.UpdateCustomerInput{CustomerID: "fake-id"},
			mocksSetup: func(
				_ *customersmocks.MockRepository,
				_ *authmocks.MockService,
				authctx *authmocks.MockContextReader,
			) {
				authctx.EXPECT().GetSubject(gomock.Any()).Return("another-fake-id", true)
			},
			want:    customers.UpdateCustomerOutput{},
			wantErr: customers.ErrCustomerIDMismatch,
		},
		{
			name:  "when the customer is not found, then it should return a customer not found error",
			input: customers.UpdateCustomerInput{CustomerID: "fake-id"},
			mocksSetup: func(
				repo *customersmocks.MockRepository,
				_ *authmocks.MockService,
				authctx *authmocks.MockContextReader,
			) {
				authctx.EXPECT().GetSubject(gomock.Any()).Return("fake-id", true)

				repo.EXPECT().GetCustomer(gomock.Any(), gomock.Any()).
					Return(customers.Customer{}, customers.ErrCustomerNotFound)
			},
			want:    customers.UpdateCustomerOutput{},
			wantErr: customers.ErrCustomerNotFound,
		},
		{
			name:  "when there is an unexpected error when getting the customer, then it should propagate the error",
			input: customers.UpdateCustomerInput{CustomerID: "fake-id"},
			mocksSetup: func(
				repo *customersmocks.MockRepository,
				_ *authmocks.MockService,
				authctx *authmocks.MockContextReader,
			) {
				authctx.EXPECT().GetSubject(gomock.Any()).Return("fake-id", true)

				repo.EXPECT().GetCustomer(gomock.Any(), gomock.Any()).
					Return(customers.Customer{}, errUnexpected)
			},
			want:    customers.UpdateCustomerOutput{},
			wantErr: errUnexpected,
		},
		{
			name:  "when the customer is not found when updating, then it should return a customer not found error",
			input: customers.UpdateCustomerInput{CustomerID: "fake-id"},
			mocksSetup: func(
				repo *customersmocks.MockRepository,
				_ *authmocks.MockService,
				authctx *authmocks.MockContextReader,
			) {
				authctx.EXPECT().GetSubject(gomock.Any()).Return("fake-id", true)

				repo.EXPECT().GetCustomer(gomock.Any(), gomock.Any()).
					Return(customers.Customer{ID: "fake-id"}, nil)

				repo.EXPECT().UpdateCustomer(gomock.Any(), gomock.Any()).
					Return(customers.Customer{}, customers.ErrCustomerNotFound)
			},
			want:    customers.UpdateCustomerOutput{},
			wantErr: customers.ErrCustomerNotFound,
		},
		{
			name:  "when there is an unexpected error when updating the customer, then it should propagate the error",
			input: customers.UpdateCustomerInput{CustomerID: "fake-id"},
			mocksSetup: func(
				repo *customersmocks.MockRepository,
				_ *authmocks.MockService,
				authctx *authmocks.MockContextReader,
			) {
				authctx.EXPECT().GetSubject(gomock.Any()).Return("fake-id", true)

				repo.EXPECT().GetCustomer(gomock.Any(), gomock.Any()).
					Return(customers.Customer{ID: "fake-id"}, nil)

				repo.EXPECT().UpdateCustomer(gomock.Any(), gomock.Any()).
					Return(customers.Customer{}, errUnexpected)
			},
			want:    customers.UpdateCustomerOutput{},
			wantErr: errUnexpected,
		},
		{
			name: "when there is an unexpected error rolling back the customer data, " +
				"then it should propagate the authservice error",
			input: customers.UpdateCustomerInput{
				CustomerID:  "fake-id",
				Name:        "New John Doe",
				Address:     "New 123 Main St",
				City:        "Los Angeles",
				PostalCode:  "09001",
				CountryCode: "SP",
			},
			mocksSetup: func(
				repo *customersmocks.MockRepository,
				authservice *authmocks.MockService,
				authctx *authmocks.MockContextReader,
			) {
				authctx.EXPECT().GetSubject(gomock.Any()).Return("fake-id", true)

				repo.EXPECT().GetCustomer(gomock.Any(), gomock.Any()).
					Return(customers.Customer{ID: "fake-id"}, nil)

				repo.EXPECT().UpdateCustomer(gomock.Any(), gomock.Any()).
					Return(customers.Customer{ID: "fake-id"}, nil).Times(1)

				authservice.EXPECT().UpdateCustomer(gomock.Any(), gomock.Any()).
					Return(authentication.UpdateCustomerOutput{}, errAuthService)

				repo.EXPECT().UpdateCustomer(gomock.Any(), gomock.Any()).
					Return(customers.Customer{}, errUnexpected).Times(1)
			},
			want:    customers.UpdateCustomerOutput{},
			wantErr: errAuthService,
		},
		{
			// rollback happy path
			name: "when there is an unexpected error updating the customer data at auth service, " +
				"then it should propagate the error",
			input: customers.UpdateCustomerInput{
				CustomerID:  "fake-id",
				Name:        "New John Doe",
				Address:     "New 123 Main St",
				City:        "Los Angeles",
				PostalCode:  "09001",
				CountryCode: "SP",
			},
			mocksSetup: func(
				repo *customersmocks.MockRepository,
				authservice *authmocks.MockService,
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
						CreatedAt:   yesterday,
						UpdatedAt:   yesterday,
					}, nil)

				repo.EXPECT().UpdateCustomer(gomock.Any(), customers.UpdateCustomerParams{
					CustomerID:  "fake-id",
					Name:        "New John Doe",
					Address:     "New 123 Main St",
					City:        "Los Angeles",
					PostalCode:  "09001",
					CountryCode: "SP",
				}).Return(customers.Customer{
					ID:          "fake-id",
					Email:       "test@example.com",
					Name:        "New John Doe",
					Active:      true,
					Address:     "New 123 Main St",
					City:        "Los Angeles",
					PostalCode:  "09001",
					CountryCode: "SP",
					CreatedAt:   yesterday,
					UpdatedAt:   now,
				}, nil).Times(1)

				authservice.EXPECT().UpdateCustomer(gomock.Any(), gomock.Any()).
					Return(authentication.UpdateCustomerOutput{}, errAuthService)

				repo.EXPECT().UpdateCustomer(gomock.Any(), customers.UpdateCustomerParams{
					CustomerID:  "fake-id",
					Name:        "John Doe",
					Address:     "123 Main St",
					City:        "New York",
					PostalCode:  "10001",
					CountryCode: "US",
				}).Return(customers.Customer{
					ID:          "fake-id",
					Email:       "test@example.com",
					Name:        "John Doe",
					Active:      true,
					Address:     "123 Main St",
					City:        "New York",
					PostalCode:  "10001",
					CountryCode: "US",
					CreatedAt:   yesterday,
					UpdatedAt:   now.Add(time.Second),
				}, nil).Times(1)
			},
			want:    customers.UpdateCustomerOutput{},
			wantErr: errAuthService,
		},
		{
			name: "when the customer is updated, then it should return the updated customer data",
			input: customers.UpdateCustomerInput{
				CustomerID:  "fake-id",
				Name:        "New John Doe",
				Address:     "New 123 Main St",
				City:        "Los Angeles",
				PostalCode:  "09001",
				CountryCode: "SP",
			},
			mocksSetup: func(
				repo *customersmocks.MockRepository,
				authservice *authmocks.MockService,
				authctx *authmocks.MockContextReader,
			) {
				authctx.EXPECT().GetSubject(gomock.Any()).Return("fake-id", true)

				repo.EXPECT().GetCustomer(gomock.Any(), gomock.Any()).
					Return(customers.Customer{ID: "fake-id"}, nil)

				repo.EXPECT().UpdateCustomer(gomock.Any(), customers.UpdateCustomerParams{
					CustomerID:  "fake-id",
					Name:        "New John Doe",
					Address:     "New 123 Main St",
					City:        "Los Angeles",
					PostalCode:  "09001",
					CountryCode: "SP",
				}).Return(customers.Customer{
					ID:          "fake-id",
					Email:       "test@example.com",
					Name:        "New John Doe",
					Active:      true,
					Address:     "New 123 Main St",
					City:        "Los Angeles",
					PostalCode:  "09001",
					CountryCode: "SP",
					CreatedAt:   yesterday,
					UpdatedAt:   now,
				}, nil).Times(1)

				authservice.EXPECT().UpdateCustomer(gomock.Any(), authentication.UpdateCustomerInput{
					CustomerID: "fake-id",
					Name:       "New John Doe",
				}).Return(authentication.UpdateCustomerOutput{
					ID:        "auth-fake-id",
					Email:     "test@example.com",
					Name:      "New John Doe",
					CreatedAt: yesterday,
					UpdatedAt: now,
				}, nil)
			},
			want: customers.UpdateCustomerOutput{
				ID:          "fake-id",
				Email:       "test@example.com",
				Name:        "New John Doe",
				Address:     "New 123 Main St",
				City:        "Los Angeles",
				PostalCode:  "09001",
				CountryCode: "SP",
				CreatedAt:   yesterday,
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
			authservice := authmocks.NewMockService(ctrl)
			authctx := authmocks.NewMockContextReader(ctrl)

			if tt.mocksSetup != nil {
				tt.mocksSetup(repo, authservice, authctx)
			}

			service := customers.NewService(logger, repo, authservice, authctx)
			got, err := service.UpdateCustomer(context.Background(), tt.input)

			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}
