// Package refresh provides functionality for managing refresh tokens in the authentication system,
// including token refresh operations and related business logic.
package refresh

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/clock"
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
	FindActiveToken(ctx context.Context, input FindActiveTokenInput) (FindActiveTokenOutput, error)
	Expire(ctx context.Context, input ExpireInput) (ExpireOutput, error)
}

// GenerateTokenInput represents the input data required for generating a token.
type GenerateTokenInput struct {
	UserID string
	Role   string
}

// GenerateTokenOutput represents the output result of a token generation operation.
type GenerateTokenOutput struct {
	Token string
}

// FindActiveTokenInput represents the input required to locate an active refresh token.
type FindActiveTokenInput struct {
	Token string
}

// FindActiveTokenOutput represents the result of a query to locate an active refresh token associated with a user.
type FindActiveTokenOutput struct {
	ID     string
	Token  string
	UserID string
	Role   string
	Device DeviceInfo
}

// ExpireInput represents the input required to mark a token as expired.
type ExpireInput struct {
	Token string
}

// ExpireOutput represents the output structure of a token expiration operation.
type ExpireOutput struct {
	ID        string
	Token     string
	ExpiresAt time.Time
}

type service struct {
	logger *zap.Logger
	repo   Repository
	clock  clock.Clock
}

// NewService initializes and returns a new Service implementation.
func NewService(logger *zap.Logger, repo Repository, clk clock.Clock) Service {
	return &service{logger: logger, repo: repo, clock: clk}
}

func (s *service) Generate(ctx context.Context, input GenerateTokenInput) (GenerateTokenOutput, error) {
	token, err := generateToken()
	if err != nil {
		logctx.LoggerWithRequestInfo(ctx, s.logger).Error("failed to generate refresh token", zap.Error(err))
		return GenerateTokenOutput{}, err
	}

	device := getDeviceFromContext(ctx)
	params := CreateTokenParams{
		UserID:    input.UserID,
		Role:      input.Role,
		Token:     token,
		ExpiresAt: time.Now().Add(DefaultTokenExpiration),
		Device:    device,
	}

	refreshToken, err := s.repo.Create(ctx, params)
	if err != nil {
		logctx.LoggerWithRequestInfo(ctx, s.logger).Error("failed to store refresh token", zap.Error(err))
		return GenerateTokenOutput{}, err
	}
	return GenerateTokenOutput{Token: refreshToken.Token}, nil
}

func (s *service) FindActiveToken(ctx context.Context, input FindActiveTokenInput) (FindActiveTokenOutput, error) {
	token, err := s.repo.FindActiveToken(ctx, input.Token)
	if err != nil {
		return FindActiveTokenOutput{}, err
	}

	return FindActiveTokenOutput{
		ID:     token.ID,
		Token:  token.Token,
		UserID: token.UserID,
		Role:   token.Role,
		Device: token.DeviceInfo,
	}, nil
}

func (s *service) Expire(ctx context.Context, input ExpireInput) (ExpireOutput, error) {
	token, err := s.repo.Expire(ctx, ExpireParams{Token: input.Token})
	if err != nil {
		return ExpireOutput{}, err
	}
	return ExpireOutput{
		ID:        token.ID,
		Token:     token.Token,
		ExpiresAt: token.ExpiresAt,
	}, nil
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

func generateDeviceID(userAgent, ip string) string {
	data := fmt.Sprintf("%s|%s", userAgent, ip)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

func getDeviceFromContext(ctx context.Context) DeviceInfo {
	ip := logctx.RealIPFromContext(ctx)
	userAgent := logctx.UserAgentFromContext(ctx)
	deviceID := generateDeviceID(userAgent, ip)

	return DeviceInfo{
		DeviceID:    deviceID,
		UserAgent:   userAgent,
		IP:          ip,
		FirstUsedAt: time.Now(),
		LastUsedAt:  time.Now(),
	}
}
