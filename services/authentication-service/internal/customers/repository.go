package customers

import (
	"context"
	"time"

	"go.uber.org/zap"
)

//go:generate mockgen -destination=./mocks/repository_mock.go -package=mocks github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/customers Repository
type Repository interface {
	CreateCustomer(ctx context.Context, params CreateCustomerParams) (Customer, error)
}

type CreateCustomerParams struct {
	Email    string
	Password string
	Name     string
}

type Customer struct {
	ID        string
	Email     string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	Active    bool
}

type repository struct {
	logger *zap.Logger
}

func NewRepository(logger *zap.Logger) Repository {
	return &repository{
		logger: logger,
	}
}

func (r *repository) CreateCustomer(ctx context.Context, params CreateCustomerParams) (Customer, error) {
	//TODO implement me
	panic("implement me")
}
