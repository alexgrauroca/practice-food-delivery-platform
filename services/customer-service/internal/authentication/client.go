// Package authentication provides functionality for customer authentication and registration
// through integration with the authentication service. It handles customer registration,
// login, and token management operations.
package authentication

import (
	"context"
	"time"

	"github.com/alexgrauroca/practice-food-delivery-platform/authclient"
)

// Client defines the interface for interacting with the authentication service.
// It provides methods for customer registration and authentication operations.
//
//go:generate mockgen -destination=./mocks/authclient_mock.go -package=authentication_mocks github.com/alexgrauroca/practice-food-delivery-platform/services/customer-service/internal/authentication Client
type Client interface {
	RegisterCustomer(ctx context.Context, req RegisterCustomerRequest) (RegisterCustomerResponse, error)
}

// Config holds the configuration options for the authentication client.
type Config struct {
	Debug bool
}

// RegisterCustomerRequest represents the data required to register a new customer
// in the authentication service.
type RegisterCustomerRequest struct {
	Email    string
	Password string
	Name     string
}

// RegisterCustomerResponse contains the data returned after successfully
// registering a customer in the authentication service.
type RegisterCustomerResponse struct {
	ID        string
	Email     string
	Name      string
	CreatedAt time.Time
}

type client struct {
	conf   *authclient.Configuration
	apicli *authclient.APIClient
}

func NewClient(config Config) Client {
	conf := authclient.NewConfiguration()
	conf.Debug = config.Debug

	apiclient := authclient.NewAPIClient(conf)
	return &client{
		conf:   conf,
		apicli: apiclient,
	}
}

func (c *client) RegisterCustomer(ctx context.Context, req RegisterCustomerRequest) (RegisterCustomerResponse, error) {
	authreq := authclient.RegisterCustomerRequest(req)
	resp, _, err := c.apicli.CustomersAPI.RegisterCustomer(ctx).RegisterCustomerRequest(authreq).Execute()
	if err != nil {
		return RegisterCustomerResponse{}, err
	}
	return RegisterCustomerResponse{
		ID:        resp.Id,
		Email:     resp.Email,
		Name:      resp.Name,
		CreatedAt: resp.CreatedAt,
	}, nil
}
