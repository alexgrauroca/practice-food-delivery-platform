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

// NewService creates a new instance of the Service interface.
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
	logger := s.logger.WithContext(ctx)
	logger.Info(
		"Registering staff owner",
		log.Field{Key: "email", Value: input.Email},
		log.Field{Key: "restaurant_id", Value: input.RestaurantID},
	)

	params := CreateStaffParams{
		Email:        input.Email,
		RestaurantID: input.RestaurantID,
		Owner:        true,
		Name:         input.Name,
		Address:      input.Address,
		City:         input.City,
		PostalCode:   input.PostalCode,
		CountryCode:  input.CountryCode,
	}
	staff, err := s.repo.CreateStaff(ctx, params)
	if err != nil {
		logger.Error("failed to create staff owner", err)
		return RegisterStaffOwnerOutput{}, err
	}

	req := authentication.RegisterStaffRequest{
		StaffID:      staff.ID,
		Email:        input.Email,
		Password:     input.Password,
		RestaurantID: input.RestaurantID,
	}
	if _, err := s.authcli.RegisterStaff(ctx, req); err != nil {
		logger.Error("failed to register staff owner at auth service", err)

		// Roll back the created staff in case of error when registering the staff at auth service.
		if err := s.repo.PurgeStaff(ctx, input.Email); err != nil {
			logger.Error("failed to purge staff owner", err)
		}
		return RegisterStaffOwnerOutput{}, err
	}

	logger.Info("staff owner registered successfully", log.Field{Key: "id", Value: staff.ID})
	return RegisterStaffOwnerOutput{
		ID:           staff.ID,
		Email:        staff.Email,
		RestaurantID: staff.RestaurantID,
		Owner:        staff.Owner,
		Name:         staff.Name,
		Address:      staff.Address,
		City:         staff.City,
		PostalCode:   staff.PostalCode,
		CountryCode:  staff.CountryCode,
		CreatedAt:    staff.CreatedAt,
		UpdatedAt:    staff.UpdatedAt,
	}, nil
}
