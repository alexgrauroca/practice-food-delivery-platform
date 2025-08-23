package authentication

import (
	"context"
)

// Claims represent the authentication claims containing subject and role information
type Claims struct {
	Subject string
	Role    string
}

// Service defines the interface for authentication operations
//
//go:generate mockgen -destination=./mocks/service_mock.go -package=authentication_mocks github.com/alexgrauroca/practice-food-delivery-platform/services/customer-service/internal/authentication Service
type Service interface {
	RegisterCustomer(ctx context.Context, input RegisterCustomerInput) (RegisterCustomerOutput, error)
	ValidateAccessToken(ctx context.Context, input ValidateAccessTokenInput) (ValidateAccessTokenOutput, error)
}

// RegisterCustomerInput TODO: implement me
type RegisterCustomerInput struct{}

// RegisterCustomerOutput TODO: implement me
type RegisterCustomerOutput struct{}

// ValidateAccessTokenInput contains the access token to be validated
type ValidateAccessTokenInput struct {
	AccessToken string
}

// ValidateAccessTokenOutput contains the validated claims from the access token
type ValidateAccessTokenOutput struct {
	Claims *Claims
}
