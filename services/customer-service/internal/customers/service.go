package customers

import (
	"context"
	"errors"
	"time"

	"github.com/alexgrauroca/practice-food-delivery-platform/pkg/auth"
	"github.com/alexgrauroca/practice-food-delivery-platform/pkg/clients/authentication"
	"github.com/alexgrauroca/practice-food-delivery-platform/pkg/log"
)

// Service defines the interface for customer management service.
//
//go:generate mockgen -destination=./mocks/service_mock.go -package=customers_mocks github.com/alexgrauroca/practice-food-delivery-platform/services/customer-service/internal/customers Service
type Service interface {
	RegisterCustomer(ctx context.Context, input RegisterCustomerInput) (RegisterCustomerOutput, error)
	GetCustomer(ctx context.Context, input GetCustomerInput) (GetCustomerOutput, error)
	UpdateCustomer(ctx context.Context, input UpdateCustomerInput) (UpdateCustomerOutput, error)
}

type service struct {
	logger  log.Logger
	repo    Repository
	authcli authentication.Client
	authctx auth.ContextReader
}

// NewService creates a new instance of Service with the provided logger and repository dependencies.
func NewService(
	logger log.Logger,
	repo Repository,
	authcli authentication.Client,
	authctx auth.ContextReader,
) Service {
	return &service{
		logger:  logger,
		repo:    repo,
		authcli: authcli,
		authctx: authctx,
	}
}

// RegisterCustomerInput defines the input structure required for registering a new customer.
type RegisterCustomerInput struct {
	Email       string
	Password    string
	Name        string
	Address     string
	City        string
	PostalCode  string
	CountryCode string
}

// RegisterCustomerOutput represents the output data returned after successfully registering a new customer.
type RegisterCustomerOutput struct {
	ID          string
	Email       string
	Name        string
	Address     string
	City        string
	PostalCode  string
	CountryCode string
	CreatedAt   time.Time
}

func (s *service) RegisterCustomer(ctx context.Context, input RegisterCustomerInput) (RegisterCustomerOutput, error) {
	logger := s.logger.WithContext(ctx)
	logger.Info("registering customer",
		log.Field{Key: "email", Value: input.Email}, log.Field{Key: "name", Value: input.Name})

	params := CreateCustomerParams{
		Email:       input.Email,
		Name:        input.Name,
		Address:     input.Address,
		City:        input.City,
		PostalCode:  input.PostalCode,
		CountryCode: input.CountryCode,
	}

	customer, err := s.repo.CreateCustomer(ctx, params)
	if err != nil {
		logger.Error("failed to create customer", err)
		return RegisterCustomerOutput{}, err
	}

	req := authentication.RegisterCustomerRequest{
		CustomerID: customer.ID,
		Email:      input.Email,
		Password:   input.Password,
	}
	if _, err := s.authcli.RegisterCustomer(ctx, req); err != nil {
		logger.Error("failed to register customer at auth service", err)

		// Roll back the created customer in case of error when registering the customer at auth service.
		if err := s.repo.PurgeCustomer(ctx, input.Email); err != nil && !errors.Is(err, ErrCustomerNotFound) {
			logger.Error("failed to purge customer", err)
			return RegisterCustomerOutput{}, err
		}
		return RegisterCustomerOutput{}, err
	}

	output := RegisterCustomerOutput{
		ID:          customer.ID,
		Email:       customer.Email,
		Name:        customer.Name,
		Address:     customer.Address,
		City:        customer.City,
		PostalCode:  customer.PostalCode,
		CountryCode: customer.CountryCode,
		CreatedAt:   customer.CreatedAt,
	}
	logger.Info("customer registered successfully", log.Field{Key: "customerID", Value: customer.ID})
	return output, nil
}

// GetCustomerInput represents the input parameters required for retrieving a customer.
type GetCustomerInput struct {
	CustomerID string
}

// GetCustomerOutput represents the output data containing customer details returned from GetCustomer operation.
type GetCustomerOutput struct {
	ID          string
	Email       string
	Name        string
	Address     string
	City        string
	PostalCode  string
	CountryCode string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (s *service) GetCustomer(ctx context.Context, input GetCustomerInput) (GetCustomerOutput, error) {
	err := s.authctx.RequireSubjectMatch(ctx, input.CustomerID)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidToken) {
			return GetCustomerOutput{}, err
		}
		return GetCustomerOutput{}, ErrCustomerIDMismatch
	}

	customer, err := s.loadCustomerByID(ctx, input.CustomerID)
	if err != nil {
		return GetCustomerOutput{}, err
	}

	return GetCustomerOutput{
		ID:          customer.ID,
		Email:       customer.Email,
		Name:        customer.Name,
		Address:     customer.Address,
		City:        customer.City,
		PostalCode:  customer.PostalCode,
		CountryCode: customer.CountryCode,
		CreatedAt:   customer.CreatedAt,
		UpdatedAt:   customer.UpdatedAt,
	}, nil
}

func (s *service) loadCustomerByID(ctx context.Context, customerID string) (Customer, error) {
	customer, err := s.repo.GetCustomer(ctx, customerID)
	if err != nil {
		if errors.Is(err, ErrCustomerNotFound) {
			s.logger.Warn("customer not found", log.Field{Key: "customerID"})
			return Customer{}, ErrCustomerNotFound
		}

		s.logger.Error("failed to get customer", err)
		return Customer{}, err
	}
	return customer, nil
}

// UpdateCustomerInput represents the input parameters required for updating a customer's details.
type UpdateCustomerInput struct {
	CustomerID  string
	Name        string
	Address     string
	City        string
	PostalCode  string
	CountryCode string
}

// UpdateCustomerOutput represents the output data containing updated customer details returned from UpdateCustomer
// operation.
type UpdateCustomerOutput struct {
	ID          string
	Email       string
	Name        string
	Address     string
	City        string
	PostalCode  string
	CountryCode string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (s *service) UpdateCustomer(ctx context.Context, input UpdateCustomerInput) (UpdateCustomerOutput, error) {
	err := s.authctx.RequireSubjectMatch(ctx, input.CustomerID)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidToken) {
			return UpdateCustomerOutput{}, err
		}
		return UpdateCustomerOutput{}, ErrCustomerIDMismatch
	}

	customer, err := s.repo.UpdateCustomer(ctx, UpdateCustomerParams(input))
	if err != nil {
		if errors.Is(err, ErrCustomerNotFound) {
			s.logger.Warn("customer not found", log.Field{Key: "customerID", Value: input.CustomerID})
			return UpdateCustomerOutput{}, ErrCustomerNotFound
		}

		s.logger.Error("failed to update customer", err)
		return UpdateCustomerOutput{}, err
	}

	return UpdateCustomerOutput{
		ID:          customer.ID,
		Email:       customer.Email,
		Name:        customer.Name,
		Address:     customer.Address,
		City:        customer.City,
		PostalCode:  customer.PostalCode,
		CountryCode: customer.CountryCode,
		CreatedAt:   customer.CreatedAt,
		UpdatedAt:   customer.UpdatedAt,
	}, nil
}
