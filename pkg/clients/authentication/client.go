// Package authentication provides functionality for customer authentication and registration
// through integration with the authentication service.
package authentication

import (
	"context"
	"time"

	"github.com/alexgrauroca/practice-food-delivery-platform/authclient"
	"github.com/alexgrauroca/practice-food-delivery-platform/pkg/log"
)

// Client defines the interface for interacting with the authentication service.
// It provides methods for customer registration and authentication operations.
//
//go:generate mockgen -destination=./mocks/authclient_mock.go -package=authentication_mocks github.com/alexgrauroca/practice-food-delivery-platform/pkg/clients/authentication Client
type Client interface {
	RegisterCustomer(ctx context.Context, req RegisterCustomerRequest) (RegisterCustomerResponse, error)
	RegisterStaff(ctx context.Context, req RegisterStaffRequest) (RegisterStaffResponse, error)
}

// Config holds the configuration options for the authentication client.
type Config struct {
	Debug bool
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

// RegisterCustomerRequest represents the data required to register a new customer
// in the authentication service.
type RegisterCustomerRequest struct {
	CustomerID string
	Email      string
	Password   string
}

// RegisterCustomerResponse contains the data returned after successfully
// registering a customer in the authentication service.
type RegisterCustomerResponse struct {
	ID        string
	Email     string
	CreatedAt time.Time
}

func (c *client) RegisterCustomer(ctx context.Context, req RegisterCustomerRequest) (RegisterCustomerResponse, error) {
	c.logger.Info("Registering customer", log.Field{Key: "customerID", Value: req.CustomerID})

	authreq := authclient.NewRegisterCustomerRequest(req.CustomerID, req.Email, req.Password)
	resp, r, err := c.apicli.CustomersAPI.RegisterCustomer(ctx).RegisterCustomerRequest(*authreq).Execute()
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
		log.Field{Key: "customerID", Value: resp.GetId()},
	)
	return RegisterCustomerResponse{
		ID:        resp.GetId(),
		Email:     resp.GetEmail(),
		CreatedAt: resp.GetCreatedAt(),
	}, nil
}

// RegisterStaffRequest represents the data required to register a new staff user in the authentication service.
type RegisterStaffRequest struct {
	StaffID      string
	Email        string
	RestaurantID string
	Password     string
}

// RegisterStaffResponse contains the data returned after successfully registering a staff user in the authentication
// service.
type RegisterStaffResponse struct {
	ID           string
	StaffID      string
	Email        string
	RestaurantID string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (c *client) RegisterStaff(ctx context.Context, req RegisterStaffRequest) (RegisterStaffResponse, error) {
	c.logger.Info("Registering staff", log.Field{Key: "staffID", Value: req.StaffID})
	authreq := authclient.NewRegisterStaffRequest(req.StaffID, req.Email, req.RestaurantID, req.Password)
	resp, r, err := c.apicli.StaffAPI.RegisterStaff(ctx).RegisterStaffRequest(*authreq).Execute()
	if err != nil {
		c.logger.Warn(
			"Failed to register staff",
			log.Field{Key: "error", Value: err.Error()},
			log.Field{Key: "response", Value: r},
		)
		return RegisterStaffResponse{}, err
	}
	c.logger.Info(
		"Staff registered successfully at authentication service",
		log.Field{Key: "staffID", Value: resp.GetId()},
	)
	return RegisterStaffResponse{
		ID:           resp.GetId(),
		Email:        resp.GetEmail(),
		RestaurantID: resp.GetRestaurantId(),
		CreatedAt:    resp.GetCreatedAt(),
		UpdatedAt:    resp.GetUpdatedAt(),
	}, nil
}
