package staff

import (
	"context"
	"time"

	"github.com/alexgrauroca/practice-food-delivery-platform/pkg/clients/authentication"
	"github.com/alexgrauroca/practice-food-delivery-platform/pkg/log"
)

// Service represents the interface defining business operations related to staff management.
//
//go:generate mockgen -destination=./mocks/service_mock.go -package=staff_mocks github.com/alexgrauroca/practice-food-delivery-platform/services/restaurant-service/internal/staff Service
type Service interface {
	RegisterStaffOwner(ctx context.Context, input RegisterStaffOwnerInput) (RegisterStaffOwnerOutput, error)
}

type service struct {
	logger  log.Logger
	repo    Repository
	authcli authentication.Client
}

func NewService(logger log.Logger, repo Repository, authcli authentication.Client) Service {
	return &service{
		logger:  logger,
		repo:    repo,
		authcli: authcli,
	}
}

// RegisterStaffOwnerInput represents the input data required for registering a new staff member.
type RegisterStaffOwnerInput struct {
	Email        string
	Password     string
	RestaurantID string
	Name         string
	Address      string
	City         string
	PostalCode   string
	CountryCode  string
}

// RegisterStaffOwnerOutput represents the output data returned after successfully registering a new staff member.
type RegisterStaffOwnerOutput struct {
	ID           string
	Email        string
	RestaurantID string
	Owner        bool
	Name         string
	Address      string
	City         string
	PostalCode   string
	CountryCode  string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (s service) RegisterStaffOwner(ctx context.Context, input RegisterStaffOwnerInput) (RegisterStaffOwnerOutput, error) {
	//TODO implement me
	panic("implement me")
}
