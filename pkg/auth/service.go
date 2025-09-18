package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/alexgrauroca/practice-food-delivery-platform/pkg/clock"
	"github.com/alexgrauroca/practice-food-delivery-platform/pkg/log"
)

// DefaultTokenType represents the default type of token used for authorization, typically set to "Bearer".
const DefaultTokenType = "Bearer"

// Service defines the interface for auth operations
//
//go:generate mockgen -destination=./mocks/service_mock.go -package=auth_mocks github.com/alexgrauroca/practice-food-delivery-platform/pkg/auth Service
type Service interface {
	GenerateToken(ctx context.Context, input GenerateTokenInput) (GenerateTokenOutput, error)
	GetClaims(ctx context.Context, input GetClaimsInput) (GetClaimsOutput, error)
}

type service struct {
	logger log.Logger
	secret []byte
	clock  clock.Clock
}

// NewService creates a new auth service instance
func NewService(logger log.Logger, secret []byte, clock clock.Clock) Service {
	return &service{
		logger: logger,
		secret: secret,
		clock:  clock,
	}
}

// GenerateTokenInput contains the required information to generate a JWT token
type GenerateTokenInput struct {
	ID         string
	Expiration int // AccessToken expiration duration in seconds
	Role       string
}

// GenerateTokenOutput contains the generated access token
type GenerateTokenOutput struct {
	AccessToken string
}

func (s service) GenerateToken(_ context.Context, input GenerateTokenInput) (GenerateTokenOutput, error) {
	now := s.clock.Now()
	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   input.ID,
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(input.Expiration) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
		Role: input.Role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString(s.secret)
	if err != nil {
		return GenerateTokenOutput{}, err
	}
	return GenerateTokenOutput{AccessToken: accessToken}, nil
}

// GetClaimsInput contains the access token from which to extract claims
type GetClaimsInput struct {
	AccessToken string
}

// GetClaimsOutput contains the claims extracted from the access token
type GetClaimsOutput struct {
	Claims *Claims
}

func (s service) GetClaims(_ context.Context, input GetClaimsInput) (GetClaimsOutput, error) {
	// Parse and validate the token with explicit validation options
	token, err := jwt.ParseWithClaims(input.AccessToken, &Claims{}, func(token *jwt.Token) (any, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}))
	if err != nil {
		return GetClaimsOutput{}, ErrInvalidToken
	}

	// Type assertion of the claims
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return GetClaimsOutput{}, ErrInvalidToken
	}

	return GetClaimsOutput{Claims: claims}, nil
}
