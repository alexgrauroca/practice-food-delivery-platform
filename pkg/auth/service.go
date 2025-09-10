package auth

import (
	"context"
	"fmt"

	"github.com/golang-jwt/jwt/v5"

	"github.com/alexgrauroca/practice-food-delivery-platform/pkg/log"
)

// Service defines the interface for auth operations
//
//go:generate mockgen -destination=./mocks/service_mock.go -package=auth_mocks github.com/alexgrauroca/practice-food-delivery-platform/pkg/auth Service
type Service interface {
	ValidateAccessToken(ctx context.Context, input ValidateAccessTokenInput) (ValidateAccessTokenOutput, error)
}

type service struct {
	logger log.Logger
	secret []byte
}

// NewService creates a new auth service instance
func NewService(logger log.Logger, secret []byte) Service {
	return &service{
		logger: logger,
		secret: secret,
	}
}

// ValidateAccessTokenInput contains the access token to be validated
type ValidateAccessTokenInput struct {
	AccessToken string
}

// ValidateAccessTokenOutput contains the validated claims from the access token
type ValidateAccessTokenOutput struct {
	Claims *Claims
}

func (s service) ValidateAccessToken(
	_ context.Context,
	input ValidateAccessTokenInput,
) (ValidateAccessTokenOutput, error) {
	token, err := jwt.ParseWithClaims(input.AccessToken, &Claims{}, func(token *jwt.Token) (any, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}))

	if err != nil {
		return ValidateAccessTokenOutput{}, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return ValidateAccessTokenOutput{}, ErrInvalidToken
	}

	return ValidateAccessTokenOutput{Claims: claims}, nil
}
