package customers

import (
	"errors"
	"go.uber.org/zap"
	"time"
)

var (
	ErrCustomerAlreadyExists = errors.New("customer already exists")
)

//go:generate mockgen -destination=./mocks/service_mock.go -package=mocks github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/customers Service
type Service interface {
	RegisterCustomer(input RegisterCustomerInput) (RegisterCustomerOutput, error)
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
}

func NewService(logger *zap.Logger) Service {
	return &service{
		logger: logger,
	}
}

func (s *service) RegisterCustomer(input RegisterCustomerInput) (RegisterCustomerOutput, error) {
	//TODO implement me
	panic("implement me")
}
