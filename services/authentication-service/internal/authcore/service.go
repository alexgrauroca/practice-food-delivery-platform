package authcore

import (
	"context"
	"errors"

	"github.com/alexgrauroca/practice-food-delivery-platform/pkg/auth"
	"github.com/alexgrauroca/practice-food-delivery-platform/pkg/log"
	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/refresh"
)

// TokenPair represents a pair of tokens typically used for authentication and session management.
type TokenPair struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int // Number of seconds until the token expires
	TokenType    string
}

// Service defines the interface for the core authentication management service.
//
//go:generate mockgen -destination=./mocks/service_mock.go -package=authcore_mocks github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/authcore Service
type Service interface {
	GenerateTokenPair(ctx context.Context, input GenerateTokenPairInput) (TokenPair, error)
	RefreshToken(ctx context.Context, input RefreshTokenInput) (TokenPair, error)
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
	TenantID   string
}

func (s service) GenerateTokenPair(ctx context.Context, input GenerateTokenPairInput) (TokenPair, error) {
	logger := s.logger.WithContext(ctx)

	generateOutput, err := s.authService.GenerateToken(ctx, auth.GenerateTokenInput{
		ID:         input.UserID,
		Expiration: input.Expiration,
		Role:       input.Role,
		TenantID:   input.TenantID,
	})
	if err != nil {
		logger.Error("failed to generate JWT", err)
		return TokenPair{}, err
	}

	refreshToken, err := s.refreshService.Generate(ctx, refresh.GenerateTokenInput{
		UserID:   input.UserID,
		Role:     input.Role,
		TenantID: input.TenantID,
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

// RefreshTokenInput defines the input structure required for refreshing a token pair.
type RefreshTokenInput struct {
	AccessToken  string
	RefreshToken string
	Expiration   int
	Role         string
}

func (s service) RefreshToken(ctx context.Context, input RefreshTokenInput) (TokenPair, error) {
	logger := s.logger.WithContext(ctx)

	logger.Info("refreshing token")
	refreshToken, err := s.refreshService.FindActiveToken(ctx, refresh.FindActiveTokenInput{
		Token: input.RefreshToken,
	})
	if err != nil {
		if errors.Is(err, refresh.ErrRefreshTokenNotFound) {
			logger.Warn("refresh token not found")
			return TokenPair{}, ErrInvalidRefreshToken
		}
		logger.Error("failed to find active refresh token", err)
		return TokenPair{}, err
	}

	claimsOutput, err := s.authService.GetClaims(ctx, auth.GetClaimsInput{AccessToken: input.AccessToken})
	if err != nil {
		if errors.Is(err, auth.ErrInvalidToken) {
			logger.Warn("access token is invalid")
			return TokenPair{}, ErrTokenMismatch
		}
		logger.Error("failed to get claims from access token", err)
		return TokenPair{}, err
	}

	claims := claimsOutput.Claims
	if claims.Subject != refreshToken.UserID ||
		claims.Role != refreshToken.Role ||
		claims.Tenant != refreshToken.TenantID {
		
		logger.Warn("token mismatch")
		return TokenPair{}, ErrTokenMismatch
	}

	tokenPair, err := s.GenerateTokenPair(ctx, GenerateTokenPairInput{
		UserID:     refreshToken.UserID,
		Expiration: input.Expiration,
		Role:       input.Role,
		TenantID:   refreshToken.TenantID,
	})
	if err != nil {
		logger.Error("failed to generate token pair", err)
		return TokenPair{}, err
	}

	_, err = s.refreshService.Expire(ctx, refresh.ExpireInput{Token: input.RefreshToken})
	// ErrRefreshTokenNotFound is silent because it does not affect the result of the workflow
	if err != nil && !errors.Is(err, refresh.ErrRefreshTokenNotFound) {
		logger.Error("failed to expire refresh token", err)
		return TokenPair{}, err
	}

	return tokenPair, nil
}
