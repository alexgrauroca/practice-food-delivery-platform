package customers

import (
	"context"
	"errors"
	"time"

	"github.com/alexgrauroca/practice-food-delivery-platform/services/customer-service/internal/authentication"
	"github.com/alexgrauroca/practice-food-delivery-platform/services/customer-service/internal/log"
)

// Service defines the interface for customer management service.
//
//go:generate mockgen -destination=./mocks/service_mock.go -package=customers_mocks github.com/alexgrauroca/practice-food-delivery-platform/services/customer-service/internal/customers Service
type Service interface {
	RegisterCustomer(ctx context.Context, input RegisterCustomerInput) (RegisterCustomerOutput, error)
	GetCustomer(ctx context.Context, input GetCustomerInput) (GetCustomerOutput, error)
}

type service struct {
	logger  log.Logger
	repo    Repository
	authcli authentication.Client
	authctx authentication.ContextReader
}

// NewService creates a new instance of Service with the provided logger and repository dependencies.
func NewService(logger log.Logger, repo Repository, authcli authentication.Client, authctx authentication.ContextReader) Service {
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

	authReq := authentication.RegisterCustomerRequest{
		CustomerID: customer.ID,
		Email:      input.Email,
		Password:   input.Password,
		Name:       input.Name,
	}
	if _, err := s.authcli.RegisterCustomer(ctx, authReq); err != nil {
		logger.Error("failed to register customer at auth service", err)

		// Roll back the created customer in case of error when registering the customer at auth service.
		// Customer not found error is ignored.
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
	authCustomerID, ok := s.authctx.GetSubject(ctx)
	if !ok {
		s.logger.Warn("authentication context not found")
		return GetCustomerOutput{}, authentication.ErrInvalidToken
	}
	if authCustomerID != input.CustomerID {
		s.logger.Warn(
			"customer ID mismatch with the token",
			log.Field{Key: "customerID", Value: input.CustomerID},
			log.Field{Key: "authCustomerID", Value: authCustomerID},
		)
		return GetCustomerOutput{}, ErrCustomerIDMismatch
	}

	customer, err := s.repo.GetCustomer(ctx, input.CustomerID)
	if err != nil {
		if errors.Is(err, ErrCustomerNotFound) {
			s.logger.Warn("customer not found", log.Field{Key: "customerID", Value: input.CustomerID})
			return GetCustomerOutput{}, ErrCustomerNotFound
		}

		s.logger.Error("failed to get customer", err)
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
