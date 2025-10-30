package staff

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
	DefaultTokenRole = "staff"
)

// Service defines the interface for the staff service.
//
//go:generate mockgen -destination=./mocks/service_mock.go -package=staff_mocks github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/staff Service
type Service interface {
	RegisterStaff(ctx context.Context, input RegisterStaffInput) (RegisterStaffOutput, error)
	LoginStaff(ctx context.Context, input LoginStaffInput) (LoginStaffOutput, error)
	RefreshStaff(ctx context.Context, input RefreshStaffInput) (RefreshStaffOutput, error)
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

// LoginStaffInput represents the input required for the staff user login process.
type LoginStaffInput struct {
	Email    string
	Password string
}

// LoginStaffOutput represents the output returned upon successful login of a staff user.
type LoginStaffOutput struct {
	authcore.TokenPair
}

func (s *service) LoginStaff(ctx context.Context, input LoginStaffInput) (LoginStaffOutput, error) {
	logger := s.logger.WithContext(ctx)

	logger.Info("logging in", log.Field{Key: "email", Value: input.Email})
	customer, err := s.repo.FindByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, ErrStaffNotFound) {
			logger.Warn("customer not found", log.Field{Key: "email", Value: input.Email})
			return LoginStaffOutput{}, authcore.ErrInvalidCredentials
		}
		logger.Error("failed to find customer by email", err)
		return LoginStaffOutput{}, err
	}

	// Check if the stored password matches the provided password
	if !password.Verify(customer.Password, input.Password) {
		logger.Warn("invalid credentials")
		return LoginStaffOutput{}, authcore.ErrInvalidCredentials
	}

	tokenPair, err := s.authCoreService.GenerateTokenPair(ctx, authcore.GenerateTokenPairInput{
		UserID:     customer.StaffID,
		Expiration: DefaultTokenExpiration,
		Role:       DefaultTokenRole,
	})
	if err != nil {
		logger.Error("failed to generate token pair", err)
		return LoginStaffOutput{}, err
	}

	return LoginStaffOutput{TokenPair: tokenPair}, nil
}

// RefreshStaffInput represents the input required to refresh a staff's authentication tokens.
type RefreshStaffInput struct {
	RefreshToken string
	AccessToken  string
}

// RefreshStaffOutput wraps the response of a successful staff token refresh operation.
type RefreshStaffOutput struct {
	authcore.TokenPair
}

func (s *service) RefreshStaff(ctx context.Context, input RefreshStaffInput) (RefreshStaffOutput, error) {
	//TODO implement me
	panic("implement me")
}
