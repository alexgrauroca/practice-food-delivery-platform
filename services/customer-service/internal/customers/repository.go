package customers

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

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
	// FieldID represents the field name used to store the unique CustomerID of a customer in the database.
	FieldID = "_id"
	// FieldName represents the field name used to store the customer's name in the database.
	FieldName = "name"
	// FieldAddress represents the field name used to store the customer's address in the database.
	FieldAddress = "address"
	// FieldCity represents the field name used to store the customer's city in the database.
	FieldCity = "city"
	// FieldPostalCode represents the field name used to store the customer's postal code in the database.
	FieldPostalCode = "postal_code"
	// FieldCountryCode represents the field name used to store the customer's country code in the database.
	FieldCountryCode = "country_code"
	// FieldUpdatedAt represents the field name used to store the timestamp when the customer was last updated.
	FieldUpdatedAt = "updated_at"
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
	PurgeCustomer(ctx context.Context, email string) error
	GetCustomer(ctx context.Context, customerID string) (Customer, error)
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
	Email       string
	Name        string
	Address     string
	City        string
	PostalCode  string
	CountryCode string
}

// CreateCustomer creates a new customer record in the database.
// It returns the created customer with an assigned CustomerID or an error if the operation fails.
// If a customer with the same email already exists, it returns ErrCustomerAlreadyExists.
func (r *repository) CreateCustomer(ctx context.Context, params CreateCustomerParams) (Customer, error) {
	logger := r.logger.WithContext(ctx)

	now := r.clock.Now()
	c := Customer{
		Email:       params.Email,
		Name:        params.Name,
		Active:      true,
		Address:     params.Address,
		City:        params.City,
		PostalCode:  params.PostalCode,
		CountryCode: params.CountryCode,
		CreatedAt:   now,
		UpdatedAt:   now,
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

func (r *repository) PurgeCustomer(ctx context.Context, email string) error {
	logger := r.logger.WithContext(ctx)
	logger.Info("Purging customer", log.Field{Key: "email", Value: email})

	res, err := r.collection.DeleteOne(ctx, bson.M{
		FieldEmail:  email,
		FieldActive: true,
	})
	if err != nil {
		logger.Error("Failed to purge customer", err)
		return err
	}
	if res.DeletedCount == 0 {
		logger.Warn("Customer not found", log.Field{Key: "email", Value: email})
		return ErrCustomerNotFound
	}

	logger.Info("Customer purged successfully", log.Field{Key: "email", Value: email})
	return nil
}

func (r *repository) GetCustomer(ctx context.Context, customerID string) (Customer, error) {
	logger := r.logger.WithContext(ctx)
	logger.Info("Getting customer", log.Field{Key: "customer_id", Value: customerID})

	id, err := primitive.ObjectIDFromHex(customerID)
	if err != nil {
		logger.Warn("Invalid customer CustomerID format", log.Field{Key: "customer_id", Value: customerID})
		return Customer{}, ErrCustomerNotFound
	}

	var customer Customer
	err = r.collection.FindOne(ctx, bson.M{FieldID: id}).Decode(&customer)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			logger.Warn("Customer not found", log.Field{Key: "customerID", Value: customerID})
			return Customer{}, ErrCustomerNotFound
		}
		logger.Error("Failed to get customer", err)
		return Customer{}, err
	}

	return customer, nil
}

// UpdateCustomerParams represents the data required for updating an existing customer's information.
type UpdateCustomerParams struct {
	CustomerID  string
	Name        string
	Address     string
	City        string
	PostalCode  string
	CountryCode string
}

func (r *repository) UpdateCustomer(ctx context.Context, params UpdateCustomerParams) (Customer, error) {
	logger := r.logger.WithContext(ctx)
	logger.Info("Updating customer", log.Field{Key: "customer_id", Value: params.CustomerID})

	id, err := primitive.ObjectIDFromHex(params.CustomerID)
	if err != nil {
		logger.Warn("Invalid customer CustomerID format", log.Field{Key: "customer_id", Value: params.CustomerID})
		return Customer{}, ErrCustomerNotFound
	}

	var customer Customer
	update := bson.M{
		"$set": bson.M{
			FieldName:        params.Name,
			FieldAddress:     params.Address,
			FieldCity:        params.City,
			FieldPostalCode:  params.PostalCode,
			FieldCountryCode: params.CountryCode,
			FieldUpdatedAt:   r.clock.Now(),
		},
	}

	// Update the customer document in the database and return the updated document
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	err = r.collection.FindOneAndUpdate(ctx, bson.M{FieldID: id}, update, opts).Decode(&customer)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			logger.Warn("Customer not found", log.Field{Key: "customer_id", Value: params.CustomerID})
			return Customer{}, ErrCustomerNotFound
		}
		logger.Error("Failed to update customer", err)
		return Customer{}, err
	}

	return customer, nil
}
