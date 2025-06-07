package customers

import (
	"context"
	"time"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

//go:generate mockgen -destination=./mocks/service_mock.go -package=mocks github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/customers Service
type Service interface {
	RegisterCustomer(ctx context.Context, input RegisterCustomerInput) (RegisterCustomerOutput, error)
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
	s.logger.Info("registering customer", zap.String("email", input.Email), zap.String("name", input.Name))
	hashedPassword, err := HashPassword(input.Password)
	if err != nil {
		s.logger.Error("failed to hash password", zap.Error(err))
		return RegisterCustomerOutput{}, err
	}

	params := CreateCustomerParams{
		Email:    input.Email,
		Password: hashedPassword,
		Name:     input.Name,
	}

	customer, err := s.repo.CreateCustomer(ctx, params)
	if err != nil {
		s.logger.Error("failed to create customer", zap.Error(err))
		return RegisterCustomerOutput{}, err
	}

	output := RegisterCustomerOutput{
		ID:        customer.ID,
		Email:     customer.Email,
		Name:      customer.Name,
		CreatedAt: customer.CreatedAt,
	}
	s.logger.Info("customer registered successfully", zap.String("customerID", customer.ID))
	return output, nil
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
