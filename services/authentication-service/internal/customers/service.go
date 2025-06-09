package customers

import (
	"context"
	"time"

	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/jwt"
	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/logctx"
	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/password"
	"go.uber.org/zap"
)

const (
	DefaultTokenExpiration = 3600 // 1 hour in seconds
	DefaultTokenRole       = "customer"
)

//go:generate mockgen -destination=./mocks/service_mock.go -package=mocks github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/customers Service
type Service interface {
	RegisterCustomer(ctx context.Context, input RegisterCustomerInput) (RegisterCustomerOutput, error)
	LoginCustomer(ctx context.Context, input LoginCustomerInput) (LoginCustomerOutput, error)
}

type RegisterCustomerInput struct {
	Email    string
	Password string
	Name     string
}

type RegisterCustomerOutput struct {
	ID        string
	Email     string
	Name      string
	CreatedAt time.Time
}

type LoginCustomerInput struct {
	Email    string
	Password string
}

type LoginCustomerOutput struct {
	Token     string
	ExpiresIn int // Number of seconds until the token expires
	TokenType string
}

type service struct {
	logger *zap.Logger
	repo   Repository
}

func NewService(logger *zap.Logger, repo Repository) Service {
	return &service{
		logger: logger,
		repo:   repo,
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
		if err == ErrCustomerNotFound {
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

	token, err := jwt.GenerateToken(customer.ID, jwt.Config{
		Expiration: DefaultTokenExpiration,
		Role:       DefaultTokenRole,
	})
	if err != nil {
		logctx.LoggerWithRequestInfo(ctx, s.logger).Error("failed to generate JWT", zap.Error(err))
		return LoginCustomerOutput{}, err
	}

	output := LoginCustomerOutput{
		TokenType: jwt.DefaultTokenType,
		Token:     token,
		ExpiresIn: DefaultTokenExpiration,
	}
	return output, nil
}
