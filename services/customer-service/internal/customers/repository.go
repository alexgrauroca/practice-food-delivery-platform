package customers

import (
	"context"
	"time"
)

const (
	// CollectionName defines the name of the MongoDB collection used for storing customer documents.
	CollectionName = "customers"

	// FieldEmail represents the field name used to store or query email addresses in the database.
	FieldEmail = "email"
	// FieldActive represents the field name used to indicate the active status of a customer in the database.
	FieldActive = "active"
)

// Customer represents a user in the system with associated details such as email, name, and account activation status.
type Customer struct {
	ID          string    `bson:"_id,omitempty"`
	Email       string    `bson:"email"`
	Name        string    `bson:"name"`
	Active      bool      `bson:"active"`
	Address     string    `bson:"address"`
	City        string    `bson:"city"`
	PostalCode  string    `bson:"postal_code"`
	CountryCode string    `bson:"country_code"`
	CreatedAt   time.Time `bson:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at"`
}

// Repository defines the interface for customer repository operations.
// It includes methods to create a customer and find a customer by email.
//
//go:generate mockgen -destination=./mocks/repository_mock.go -package=customers_mocks github.com/alexgrauroca/practice-food-delivery-platform/services/customer-service/internal/customers Repository
type Repository interface {
	CreateCustomer(ctx context.Context, params CreateCustomerParams) (Customer, error)
	PurgeCustomer(ctx context.Context, id string) error
}

// CreateCustomerParams represents the parameters needed to create a new customer.
type CreateCustomerParams struct {
	Email       string
	Name        string
	Address     string
	City        string
	PostalCode  string
	CountryCode string
}
