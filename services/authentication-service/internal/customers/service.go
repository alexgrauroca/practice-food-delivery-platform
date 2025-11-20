package customers

import (
	"context"
	"errors"
	"time"

	"github.com/alexgrauroca/practice-food-delivery-platform/pkg/log"
	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/authcore"
	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/password"
)

const (
	// DefaultTokenExpiration defines the duration in seconds for which a JWT token remains valid
	// after being issued during customer authentication. The default value is 3600 seconds (1 hour).
	DefaultTokenExpiration = 3600
	// DefaultTokenRole represents the default role assigned to a generated JWT token for customers.
	DefaultTokenRole = "customer"
)

// Service defines the interface for customer authentication management service.
//
//go:generate mockgen -destination=./mocks/service_mock.go -package=customers_mocks github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/customers Service
type Service interface {
	RegisterCustomer(ctx context.Context, input RegisterCustomerInput) (RegisterCustomerOutput, error)
	LoginCustomer(ctx context.Context, input LoginCustomerInput) (LoginCustomerOutput, error)
	RefreshCustomer(ctx context.Context, input RefreshCustomerInput) (RefreshCustomerOutput, error)
}

type service struct {
	logger          log.Logger
	repo            Repository
	authCoreService authcore.Service
}

// NewService creates a new instance of Service with the provided dependencies.
func NewService(
	logger log.Logger,
	repo Repository,
	authCoreService authcore.Service,
) Service {
	return &service{
		logger:          logger,
		repo:            repo,
		authCoreService: authCoreService,
	}
}

// RegisterCustomerInput defines the input structure required for registering a new customer.
type RegisterCustomerInput struct {
	CustomerID string
	Email      string
	Password   string
}

// RegisterCustomerOutput represents the output data returned after successfully registering a new customer.
type RegisterCustomerOutput struct {
	ID        string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (s *service) RegisterCustomer(ctx context.Context, input RegisterCustomerInput) (RegisterCustomerOutput, error) {
	logger := s.logger.WithContext(ctx)

	logger.Info("registering customer", log.Field{Key: "email", Value: input.Email})
	hashedPassword, err := password.Hash(input.Password)
	if err != nil {
		logger.Error("failed to hash password", err)
		return RegisterCustomerOutput{}, err
	}

	params := CreateCustomerParams{
		CustomerID: input.CustomerID,
		Email:      input.Email,
		Password:   hashedPassword,
	}

	customer, err := s.repo.CreateCustomer(ctx, params)
	if err != nil {
		logger.Error("failed to create customer", err)
		return RegisterCustomerOutput{}, err
	}

	output := RegisterCustomerOutput{
		ID:        customer.ID,
		Email:     customer.Email,
		CreatedAt: customer.CreatedAt,
		UpdatedAt: customer.UpdatedAt,
	}
	logger.Info("customer registered successfully", log.Field{Key: "customerID", Value: customer.ID})
	return output, nil
}

// LoginCustomerInput represents the input required for the customer login process.
type LoginCustomerInput struct {
	Email    string
	Password string
}

// LoginCustomerOutput represents the output returned upon successful login of a customer.
type LoginCustomerOutput struct {
	authcore.TokenPair
}

func (s *service) LoginCustomer(ctx context.Context, input LoginCustomerInput) (LoginCustomerOutput, error) {
	logger := s.logger.WithContext(ctx)

	logger.Info("logging in", log.Field{Key: "email", Value: input.Email})
	customer, err := s.repo.FindByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, ErrCustomerNotFound) {
			logger.Warn("customer not found", log.Field{Key: "email", Value: input.Email})
			return LoginCustomerOutput{}, authcore.ErrInvalidCredentials
		}
		logger.Error("failed to find customer by email", err)
		return LoginCustomerOutput{}, err
	}

	// Check if the stored password matches the provided password
	if !password.Verify(customer.Password, input.Password) {
		logger.Warn("invalid credentials")
		return LoginCustomerOutput{}, authcore.ErrInvalidCredentials
	}

	tokenPair, err := s.authCoreService.GenerateTokenPair(
		ctx, authcore.GenerateTokenPairInput{
			UserID:     customer.CustomerID,
			Expiration: DefaultTokenExpiration,
			Role:       DefaultTokenRole,
		},
	)
	if err != nil {
		logger.Error("failed to generate token pair", err)
		return LoginCustomerOutput{}, err
	}

	return LoginCustomerOutput{TokenPair: tokenPair}, nil
}

// RefreshCustomerInput represents the input required to refresh a customer's authentication tokens.
type RefreshCustomerInput struct {
	RefreshToken string
	AccessToken  string
}

// RefreshCustomerOutput wraps the response of a successful customer token refresh operation.
type RefreshCustomerOutput struct {
	authcore.TokenPair
}

func (s *service) RefreshCustomer(ctx context.Context, input RefreshCustomerInput) (RefreshCustomerOutput, error) {
	logger := s.logger.WithContext(ctx)

	logger.Info("refreshing customer token")

	tokenPair, err := s.authCoreService.RefreshToken(
		ctx, authcore.RefreshTokenInput{
			RefreshToken: input.RefreshToken,
			AccessToken:  input.AccessToken,
			Expiration:   DefaultTokenExpiration,
			Role:         DefaultTokenRole,
		},
	)
	if err != nil {
		logger.Error("failed to refresh the customer token", err)
		return RefreshCustomerOutput{}, err
	}

	return RefreshCustomerOutput{TokenPair: tokenPair}, nil
}
