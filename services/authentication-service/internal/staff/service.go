package staff

import (
	"context"
	"time"

	"github.com/alexgrauroca/practice-food-delivery-platform/pkg/auth"
	"github.com/alexgrauroca/practice-food-delivery-platform/pkg/log"
	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/password"
	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/refresh"
)

// Service defines the interface for the staff service.
//
//go:generate mockgen -destination=./mocks/service_mock.go -package=staff_mocks github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/staff Service
type Service interface {
	RegisterStaff(ctx context.Context, input RegisterStaffInput) (RegisterStaffOutput, error)
}

type service struct {
	logger         log.Logger
	repo           Repository
	refreshService refresh.Service
	authService    auth.Service
	authctx        auth.ContextReader
}

// NewService creates a new instance of Service with the provided dependencies.
func NewService(
	logger log.Logger,
	repo Repository,
	refreshService refresh.Service,
	authService auth.Service,
	authctx auth.ContextReader,
) Service {
	return &service{
		logger:         logger,
		repo:           repo,
		refreshService: refreshService,
		authService:    authService,
		authctx:        authctx,
	}
}

// RegisterStaffInput defines the input structure required for registering a new staff.
type RegisterStaffInput struct {
	StaffID  string
	Email    string
	Password string
}

// RegisterStaffOutput represents the output data returned after successfully registering a new staff.
type RegisterStaffOutput struct {
	ID        string
	Email     string
	CreatedAt time.Time
}

func (s *service) RegisterStaff(ctx context.Context, input RegisterStaffInput) (RegisterStaffOutput, error) {
	logger := s.logger.WithContext(ctx)

	logger.Info("registering staff", log.Field{Key: "email", Value: input.Email})
	hashedPassword, err := password.Hash(input.Password)
	if err != nil {
		logger.Error("failed to hash password", err)
		return RegisterStaffOutput{}, err
	}

	params := CreateStaffParams{
		StaffID:  input.StaffID,
		Email:    input.Email,
		Password: hashedPassword,
	}

	staff, err := s.repo.CreateStaff(ctx, params)
	if err != nil {
		logger.Error("failed to create staff", err)
		return RegisterStaffOutput{}, err
	}

	output := RegisterStaffOutput{
		ID:        staff.ID,
		Email:     staff.Email,
		CreatedAt: staff.CreatedAt,
	}
	logger.Info("staff registered successfully", log.Field{Key: "staffID", Value: staff.ID})
	return output, nil
}
