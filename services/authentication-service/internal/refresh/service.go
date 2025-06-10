// Package refresh provides functionality for managing refresh tokens in the authentication system,
// including token refresh operations and related business logic.
package refresh

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"time"

	"go.uber.org/zap"

	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/logctx"
)

const (
	// DefaultRefreshTokenLength defines the default length, in bytes, of a generated refresh token for authentication purposes.
	DefaultRefreshTokenLength = 32

	// DefaultTokenExpiration specifies the default duration for which a token remains valid, set to one hour.
	DefaultTokenExpiration = 7 * 24 * time.Hour
)

// Service represents the core interface for refresh tokens.
//
//go:generate mockgen -destination=./mocks/service_mock.go -package=refresh_mocks github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/refresh Service
type Service interface {
	Generate(ctx context.Context, input GenerateTokenInput) (GenerateTokenOutput, error)
}

// GenerateTokenInput represents the input data required for generating a token.
type GenerateTokenInput struct {
	UserID string
	Role   string
}

// GenerateTokenOutput represents the output result of a token generation operation.
type GenerateTokenOutput struct {
	RefreshToken string
}

type service struct {
	logger *zap.Logger
	repo   Repository
}

// NewService initializes and returns a new Service implementation.
func NewService(logger *zap.Logger, repo Repository) Service {
	return &service{logger: logger, repo: repo}
}

func (s *service) Generate(ctx context.Context, input GenerateTokenInput) (GenerateTokenOutput, error) {
	token, err := generateToken()
	if err != nil {
		logctx.LoggerWithRequestInfo(ctx, s.logger).Error("failed to generate refresh token", zap.Error(err))
		return GenerateTokenOutput{}, err
	}

	params := CreateTokenParams{
		UserID:    input.UserID,
		Role:      input.Role,
		Token:     token,
		ExpiresAt: time.Now().Add(DefaultTokenExpiration),
	}

	if _, err := s.repo.Store(ctx, params); err != nil {
		logctx.LoggerWithRequestInfo(ctx, s.logger).Error("failed to store refresh token", zap.Error(err))
		return GenerateTokenOutput{}, err
	}
	return GenerateTokenOutput{RefreshToken: token}, nil
}

func generateToken() (string, error) {
	// Creating a cryptographically secure random refresh token by:
	// 1. Allocating a byte slice of defined length (32 bytes)
	// 2. Filling it with random bytes using crypto/rand
	// 3. Encoding the random bytes to base64 string
	b := make([]byte, DefaultRefreshTokenLength)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
