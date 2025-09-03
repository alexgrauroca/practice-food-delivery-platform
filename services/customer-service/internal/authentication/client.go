// Package authentication provides functionality for customer authentication and registration
// through integration with the authentication service. It handles customer registration,
// login, and token management operations.
package authentication

import (
	"context"
	"time"

	"github.com/alexgrauroca/practice-food-delivery-platform/authclient"
	"github.com/alexgrauroca/practice-food-delivery-platform/services/customer-service/internal/log"
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
	CustomerID string
	Email      string
	Password   string
	Name       string
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
	logger log.Logger
	conf   *authclient.Configuration
	apicli *authclient.APIClient
}

// NewClient creates and initializes a new authentication client with the provided logger and configuration.
// It sets up the underlying API client with the specified debug mode and returns an interface
// implementation for interacting with the authentication service.
func NewClient(logger log.Logger, config Config) Client {
	conf := authclient.NewConfiguration()
	conf.Debug = config.Debug
	conf.Host = "authentication-service:8080"

	apiclient := authclient.NewAPIClient(conf)
	return &client{
		logger: logger,
		conf:   conf,
		apicli: apiclient,
	}
}

func (c *client) RegisterCustomer(ctx context.Context, req RegisterCustomerRequest) (RegisterCustomerResponse, error) {
	authreq := authclient.RegisterCustomerRequest{
		CustomerId: req.CustomerID,
		Email:      req.Email,
		Password:   req.Password,
		Name:       req.Name,
	}
	resp, r, err := c.apicli.CustomersAPI.RegisterCustomer(ctx).RegisterCustomerRequest(authreq).Execute()
	if err != nil {
		c.logger.Warn(
			"Failed to register customer",
			log.Field{Key: "error", Value: err.Error()},
			log.Field{Key: "response", Value: r},
		)
		return RegisterCustomerResponse{}, err
	}
	c.logger.Info(
		"Customer registered successfully at authentication service",
		log.Field{Key: "customerID", Value: resp.Id},
	)
	return RegisterCustomerResponse{
		ID:        resp.Id,
		Email:     resp.Email,
		Name:      resp.Name,
		CreatedAt: resp.CreatedAt,
	}, nil
}
