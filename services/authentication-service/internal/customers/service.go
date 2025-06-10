package customers

import (
	"context"
	"errors"
	"time"

	"go.uber.org/zap"

	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/jwt"
	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/logctx"
	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/password"
	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/refresh"
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
}

// RegisterCustomerInput defines the input structure required for registering a new customer.
type RegisterCustomerInput struct {
	Email    string
	Password string
	Name     string
}

// RegisterCustomerOutput represents the output data returned after successfully registering a new customer.
type RegisterCustomerOutput struct {
	ID        string
	Email     string
	Name      string
	CreatedAt time.Time
}

// LoginCustomerInput represents the input required for the customer login process.
type LoginCustomerInput struct {
	Email    string
	Password string
}

// LoginCustomerOutput represents the output returned upon successful login of a customer.
type LoginCustomerOutput struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int // Number of seconds until the token expires
	TokenType    string
}

type service struct {
	logger         *zap.Logger
	repo           Repository
	refreshService refresh.Service
}

// NewService creates a new instance of Service with the provided logger and repository dependencies.
func NewService(logger *zap.Logger, repo Repository, refreshService refresh.Service) Service {
	return &service{
		logger:         logger,
		repo:           repo,
		refreshService: refreshService,
	}
}

func (s *service) RegisterCustomer(ctx context.Context, input RegisterCustomerInput) (RegisterCustomerOutput, error) {
	logctx.LoggerWithRequestInfo(ctx, s.logger).
		Info("registering customer", zap.String("email", input.Email), zap.String("name", input.Name))
	hashedPassword, err := password.Hash(input.Password)
	if err != nil {
		logctx.LoggerWithRequestInfo(ctx, s.logger).Error("failed to hash password", zap.Error(err))
		return RegisterCustomerOutput{}, err
	}

	params := CreateCustomerParams{
		Email:    input.Email,
		Password: hashedPassword,
		Name:     input.Name,
	}

	customer, err := s.repo.CreateCustomer(ctx, params)
	if err != nil {
		logctx.LoggerWithRequestInfo(ctx, s.logger).Error("failed to create customer", zap.Error(err))
		return RegisterCustomerOutput{}, err
	}

	output := RegisterCustomerOutput{
		ID:        customer.ID,
		Email:     customer.Email,
		Name:      customer.Name,
		CreatedAt: customer.CreatedAt,
	}
	logctx.LoggerWithRequestInfo(ctx, s.logger).
		Info("customer registered successfully", zap.String("customerID", customer.ID))
	return output, nil
}

func (s *service) LoginCustomer(ctx context.Context, input LoginCustomerInput) (LoginCustomerOutput, error) {
	logctx.LoggerWithRequestInfo(ctx, s.logger).Info("logging in", zap.String("email", input.Email))
	customer, err := s.repo.FindByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, ErrCustomerNotFound) {
			logctx.LoggerWithRequestInfo(ctx, s.logger).Warn("customer not found", zap.String("email", input.Email))
			return LoginCustomerOutput{}, ErrInvalidCredentials
		}
		logctx.LoggerWithRequestInfo(ctx, s.logger).Error("failed to find customer by email", zap.Error(err))
		return LoginCustomerOutput{}, err
	}

	// Check if the stored password matches the provided password
	if !password.Verify(customer.Password, input.Password) {
		logctx.LoggerWithRequestInfo(ctx, s.logger).Warn("invalid credentials", zap.Error(err))
		return LoginCustomerOutput{}, ErrInvalidCredentials
	}

	accessToken, err := jwt.GenerateToken(customer.ID, jwt.Config{
		Expiration: DefaultTokenExpiration,
		Role:       DefaultTokenRole,
	})
	if err != nil {
		logctx.LoggerWithRequestInfo(ctx, s.logger).Error("failed to generate JWT", zap.Error(err))
		return LoginCustomerOutput{}, err
	}

	refreshToken, err := s.refreshService.Generate(ctx, refresh.GenerateTokenInput{
		UserID: customer.ID,
		Role:   DefaultTokenRole,
	})
	if err != nil {
		logctx.LoggerWithRequestInfo(ctx, s.logger).Error("failed to generate refresh token", zap.Error(err))
		return LoginCustomerOutput{}, err
	}

	output := LoginCustomerOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken.RefreshToken,
		TokenType:    jwt.DefaultTokenType,
		ExpiresIn:    DefaultTokenExpiration,
	}
	return output, nil
}
