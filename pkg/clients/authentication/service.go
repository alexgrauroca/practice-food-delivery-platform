package authentication

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/alexgrauroca/practice-food-delivery-platform/pkg/log"
)

// Claims represent the authentication claims
type Claims struct {
	jwt.RegisteredClaims
	Role string `json:"role"`
}

// Service defines the interface for authentication operations
//
//go:generate mockgen -destination=./mocks/service_mock.go -package=authentication_mocks github.com/alexgrauroca/practice-food-delivery-platform/pkg/clients/authentication Service
type Service interface {
	RegisterCustomer(ctx context.Context, input RegisterCustomerInput) (RegisterCustomerOutput, error)
	ValidateAccessToken(ctx context.Context, input ValidateAccessTokenInput) (ValidateAccessTokenOutput, error)
	UpdateCustomer(ctx context.Context, input UpdateCustomerInput) (UpdateCustomerOutput, error)
}

type service struct {
	logger log.Logger
	cli    Client
	secret []byte
}

// NewService creates a new authentication service instance
func NewService(logger log.Logger, cli Client, secret []byte) Service {
	return &service{
		logger: logger,
		cli:    cli,
		secret: secret,
	}
}

// RegisterCustomerInput represents the input data required to register a new customer
type RegisterCustomerInput struct {
	CustomerID string
	Email      string
	Password   string
	Name       string
}

// RegisterCustomerOutput represents the output data returned after successfully registering a customer
type RegisterCustomerOutput struct {
	ID        string
	Email     string
	Name      string
	CreatedAt time.Time
}

func (s service) RegisterCustomer(ctx context.Context, input RegisterCustomerInput) (RegisterCustomerOutput, error) {
	s.logger.Info("registering customer in authentication service")
	resp, err := s.cli.RegisterCustomer(ctx, RegisterCustomerRequest(input))
	if err != nil {
		return RegisterCustomerOutput{}, err
	}
	return RegisterCustomerOutput(resp), nil
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

// UpdateCustomerInput represents the input data required to update a customer
type UpdateCustomerInput struct {
	CustomerID string
	Name       string
}

// UpdateCustomerOutput represents the output data returned after successfully updating a customer
type UpdateCustomerOutput struct {
	ID        string
	Email     string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (s service) UpdateCustomer(ctx context.Context, input UpdateCustomerInput) (UpdateCustomerOutput, error) {
	//TODO implement me
	panic("implement me")
}
