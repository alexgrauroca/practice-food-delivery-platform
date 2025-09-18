package customers

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/alexgrauroca/practice-food-delivery-platform/pkg/clock"
	"github.com/alexgrauroca/practice-food-delivery-platform/pkg/infraestructure/mongodb"
	"github.com/alexgrauroca/practice-food-delivery-platform/pkg/log"
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
	ID         string    `bson:"_id,omitempty"`
	CustomerID string    `bson:"customer_id"`
	Email      string    `bson:"email"`
	Name       string    `bson:"name"`
	Password   string    `bson:"password,omitempty"`
	CreatedAt  time.Time `bson:"created_at"`
	UpdatedAt  time.Time `bson:"updated_at"`
	Active     bool      `bson:"active"`
}

// Repository defines the interface for customer repository operations.
// It includes methods to create a customer and find a customer by email.
//
//go:generate mockgen -destination=./mocks/repository_mock.go -package=customers_mocks github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/customers Repository
type Repository interface {
	CreateCustomer(ctx context.Context, params CreateCustomerParams) (Customer, error)
	FindByEmail(ctx context.Context, email string) (Customer, error)
	UpdateCustomer(ctx context.Context, params UpdateCustomerParams) (Customer, error)
}

type repository struct {
	logger     log.Logger
	collection *mongo.Collection
	clock      clock.Clock
}

// NewRepository creates a new instance of the Repository interface with MongoDB implementation.
// It requires a logger for operational logging, a database connection, and a clock implementation for timestamp generation.
func NewRepository(logger log.Logger, db *mongo.Database, clk clock.Clock) Repository {
	return &repository{
		logger:     logger,
		collection: db.Collection(CollectionName),
		clock:      clk,
	}
}

// CreateCustomerParams represents the parameters needed to create a new customer.
type CreateCustomerParams struct {
	CustomerID string
	Email      string
	Password   string
	Name       string
}

// CreateCustomer creates a new customer record in the database.
// It returns the created customer with an assigned ID or an error if the operation fails.
// If a customer with the same email already exists, it returns ErrCustomerAlreadyExists.
func (r *repository) CreateCustomer(ctx context.Context, params CreateCustomerParams) (Customer, error) {
	logger := r.logger.WithContext(ctx)

	now := r.clock.Now()
	c := Customer{
		CustomerID: params.CustomerID,
		Email:      params.Email,
		Name:       params.Name,
		Password:   params.Password,
		CreatedAt:  now,
		UpdatedAt:  now,
		Active:     true,
	}
	res, err := r.collection.InsertOne(ctx, c)
	if err != nil {
		if mongodb.IsDuplicateKeyError(err) {
			logger.Warn("Customer already exists", log.Field{Key: "email", Value: params.Email})
			return Customer{}, ErrCustomerAlreadyExists
		}
		logger.Error("Failed to insert customer", err)
		return Customer{}, err
	}
	c.ID = res.InsertedID.(primitive.ObjectID).Hex()
	logger.Info("Customer created successfully", log.Field{Key: "customer_id", Value: c.ID})
	return c, nil
}

// FindByEmail searches for an active customer with the specified email address.
// It returns the customer if found or ErrCustomerNotFound if no matching active customer exists.
func (r *repository) FindByEmail(ctx context.Context, email string) (Customer, error) {
	logger := r.logger.WithContext(ctx)

	var customer Customer
	filter := bson.M{
		FieldEmail:  email,
		FieldActive: true,
	}

	if err := r.collection.FindOne(ctx, filter).Decode(&customer); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			logger.Warn("Customer not found", log.Field{Key: "email", Value: email})
			return Customer{}, ErrCustomerNotFound
		}
		logger.Error("Failed to find customer", err)
		return Customer{}, err
	}
	return customer, nil
}

type UpdateCustomerParams struct {
	CustomerID string
	Name       string
}

func (r *repository) UpdateCustomer(ctx context.Context, params UpdateCustomerParams) (Customer, error) {
	//TODO implement me
	panic("implement me")
}
