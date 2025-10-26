package authcore

import (
	"context"

	"github.com/alexgrauroca/practice-food-delivery-platform/pkg/auth"
	"github.com/alexgrauroca/practice-food-delivery-platform/pkg/log"
	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/refresh"
)

// Service defines the interface for the core authentication management service.
//
//go:generate mockgen -destination=./mocks/service_mock.go -package=authcore_mocks github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/authcore Service
type Service interface {
	GenerateTokenPair(ctx context.Context, input GenerateTokenPairInput) (TokenPair, error)
}

type service struct {
	logger         log.Logger
	authService    auth.Service
	refreshService refresh.Service
}

// NewService creates a new instance of Service with the provided dependencies.
func NewService(
	logger log.Logger,
	authService auth.Service,
	refreshService refresh.Service,
) Service {
	return &service{
		logger:         logger,
		authService:    authService,
		refreshService: refreshService,
	}
}

// GenerateTokenPairInput defines the input structure required for generating a new token pair.
type GenerateTokenPairInput struct {
	UserID     string
	Expiration int
	Role       string
}

// TokenPair represents a pair of tokens typically used for authentication and session management.
type TokenPair struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int // Number of seconds until the token expires
	TokenType    string
}

func (s service) GenerateTokenPair(ctx context.Context, input GenerateTokenPairInput) (TokenPair, error) {
	logger := s.logger.WithContext(ctx)

	generateOutput, err := s.authService.GenerateToken(ctx, auth.GenerateTokenInput{
		ID:         input.UserID,
		Expiration: input.Expiration,
		Role:       input.Role,
	})
	if err != nil {
		logger.Error("failed to generate JWT", err)
		return TokenPair{}, err
	}

	refreshToken, err := s.refreshService.Generate(ctx, refresh.GenerateTokenInput{
		UserID: input.UserID,
		Role:   input.Role,
	})
	if err != nil {
		logger.Error("failed to generate refresh token", err)
		return TokenPair{}, err
	}

	return TokenPair{
		AccessToken:  generateOutput.AccessToken,
		RefreshToken: refreshToken.Token,
		TokenType:    auth.DefaultTokenType,
		ExpiresIn:    input.Expiration,
	}, nil
}
