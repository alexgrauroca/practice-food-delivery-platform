package customers

import (
	"context"
	"errors"
	"time"

	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/jwt"
	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/log"
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
	RefreshCustomer(ctx context.Context, input RefreshCustomerInput) (RefreshCustomerOutput, error)
}

type service struct {
	logger         log.Logger
	repo           Repository
	refreshService refresh.Service
	jwtService     jwt.Service
}

// NewService creates a new instance of Service with the provided logger and repository dependencies.
func NewService(logger log.Logger, repo Repository, refreshService refresh.Service, jwtService jwt.Service) Service {
	return &service{
		logger:         logger,
		repo:           repo,
		refreshService: refreshService,
		jwtService:     jwtService,
	}
}

// RegisterCustomerInput defines the input structure required for registering a new customer.
type RegisterCustomerInput struct {
	CustomerID string
	Email      string
	Password   string
	Name       string
}

// RegisterCustomerOutput represents the output data returned after successfully registering a new customer.
type RegisterCustomerOutput struct {
	ID        string
	Email     string
	Name      string
	CreatedAt time.Time
}

func (s *service) RegisterCustomer(ctx context.Context, input RegisterCustomerInput) (RegisterCustomerOutput, error) {
	logger := s.logger.WithContext(ctx)

	logger.Info("registering customer",
		log.Field{Key: "email", Value: input.Email}, log.Field{Key: "name", Value: input.Name})
	hashedPassword, err := password.Hash(input.Password)
	if err != nil {
		logger.Error("failed to hash password", err)
		return RegisterCustomerOutput{}, err
	}

	params := CreateCustomerParams{
		CustomerID: input.CustomerID,
		Email:      input.Email,
		Password:   hashedPassword,
		Name:       input.Name,
	}

	customer, err := s.repo.CreateCustomer(ctx, params)
	if err != nil {
		logger.Error("failed to create customer", err)
		return RegisterCustomerOutput{}, err
	}

	output := RegisterCustomerOutput{
		ID:        customer.ID,
		Email:     customer.Email,
		Name:      customer.Name,
		CreatedAt: customer.CreatedAt,
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
	TokenPair
}

func (s *service) LoginCustomer(ctx context.Context, input LoginCustomerInput) (LoginCustomerOutput, error) {
	logger := s.logger.WithContext(ctx)

	logger.Info("logging in", log.Field{Key: "email", Value: input.Email})
	customer, err := s.repo.FindByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, ErrCustomerNotFound) {
			logger.Warn("customer not found", log.Field{Key: "email", Value: input.Email})
			return LoginCustomerOutput{}, ErrInvalidCredentials
		}
		logger.Error("failed to find customer by email", err)
		return LoginCustomerOutput{}, err
	}

	// Check if the stored password matches the provided password
	if !password.Verify(customer.Password, input.Password) {
		logger.Warn("invalid credentials")
		return LoginCustomerOutput{}, ErrInvalidCredentials
	}

	tokenPair, err := s.generateTokenPair(ctx, customer)
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
	TokenPair
}

func (s *service) RefreshCustomer(ctx context.Context, input RefreshCustomerInput) (RefreshCustomerOutput, error) {
	logger := s.logger.WithContext(ctx)

	logger.Info("refreshing customer token")
	refreshToken, err := s.refreshService.FindActiveToken(ctx, refresh.FindActiveTokenInput{
		Token: input.RefreshToken,
	})
	if err != nil {
		if errors.Is(err, refresh.ErrRefreshTokenNotFound) {
			logger.Warn("refresh token not found")
			return RefreshCustomerOutput{}, ErrInvalidRefreshToken
		}
		logger.Error("failed to find active refresh token", err)
		return RefreshCustomerOutput{}, err
	}

	claims, err := s.jwtService.GetClaims(input.AccessToken)
	if err != nil {
		if errors.Is(err, jwt.ErrInvalidToken) {
			logger.Warn("access token is invalid")
			return RefreshCustomerOutput{}, ErrTokenMismatch
		}
		logger.Error("failed to get claims from access token", err)
		return RefreshCustomerOutput{}, err
	}
	if claims.Subject != refreshToken.UserID || claims.Role != refreshToken.Role {
		logger.Warn("token mismatch")
		return RefreshCustomerOutput{}, ErrTokenMismatch
	}

	tokenPair, err := s.generateTokenPair(ctx, Customer{ID: refreshToken.UserID})
	if err != nil {
		logger.Error("failed to generate token pair", err)
		return RefreshCustomerOutput{}, err
	}

	_, err = s.refreshService.Expire(ctx, refresh.ExpireInput{Token: input.RefreshToken})
	// ErrRefreshTokenNotFound is silent because it does not affect the result of the workflow
	if err != nil && !errors.Is(err, refresh.ErrRefreshTokenNotFound) {
		logger.Error("failed to expire refresh token", err)
		return RefreshCustomerOutput{}, err
	}

	return RefreshCustomerOutput{TokenPair: tokenPair}, nil
}

// TokenPair represents a pair of tokens typically used for authentication and session management.
type TokenPair struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int // Number of seconds until the token expires
	TokenType    string
}

func (s *service) generateTokenPair(ctx context.Context, customer Customer) (TokenPair, error) {
	logger := s.logger.WithContext(ctx)

	accessToken, err := s.jwtService.GenerateToken(customer.CustomerID, jwt.Config{
		Expiration: DefaultTokenExpiration,
		Role:       DefaultTokenRole,
	})
	if err != nil {
		logger.Error("failed to generate JWT", err)
		return TokenPair{}, err
	}

	refreshToken, err := s.refreshService.Generate(ctx, refresh.GenerateTokenInput{
		UserID: customer.ID,
		Role:   DefaultTokenRole,
	})
	if err != nil {
		logger.Error("failed to generate refresh token", err)
		return TokenPair{}, err
	}

	return TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken.Token,
		TokenType:    jwt.DefaultTokenType,
		ExpiresIn:    DefaultTokenExpiration,
	}, nil
}
