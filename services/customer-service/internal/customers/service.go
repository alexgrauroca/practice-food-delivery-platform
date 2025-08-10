package customers

import (
	"context"
	"time"

	"github.com/alexgrauroca/practice-food-delivery-platform/services/customer-service/internal/authentication"
	"github.com/alexgrauroca/practice-food-delivery-platform/services/customer-service/internal/log"
)

// Service defines the interface for customer management service.
//
//go:generate mockgen -destination=./mocks/service_mock.go -package=customers_mocks github.com/alexgrauroca/practice-food-delivery-platform/services/customer-service/internal/customers Service
type Service interface {
	RegisterCustomer(ctx context.Context, input RegisterCustomerInput) (RegisterCustomerOutput, error)
}

type service struct {
	logger  log.Logger
	repo    Repository
	authcli authentication.Client
}

// NewService creates a new instance of Service with the provided logger and repository dependencies.
func NewService(logger log.Logger, repo Repository, authcli authentication.Client) Service {
	return &service{
		logger:  logger,
		repo:    repo,
		authcli: authcli,
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
		Email:    input.Email,
		Password: input.Password,
		Name:     input.Name,
	}
	if _, err := s.authcli.RegisterCustomer(ctx, authReq); err != nil {
		logger.Error("failed to register customer at auth service", err)
		// TODO: rollback customer creation
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
