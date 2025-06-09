package customers

import (
	"context"
	"errors"
	"time"

	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/clock"
	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/logctx"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

const (
	CollectionName = "customers"

	FieldEmail  = "email"
	FieldActive = "active"
)

//go:generate mockgen -destination=./mocks/repository_mock.go -package=mocks github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/customers Repository
type Repository interface {
	CreateCustomer(ctx context.Context, params CreateCustomerParams) (Customer, error)
	FindByEmail(ctx context.Context, email string) (Customer, error)
}

type CreateCustomerParams struct {
	Email    string
	Password string
	Name     string
}

type Customer struct {
	ID        string    `bson:"_id,omitempty"`
	Email     string    `bson:"email"`
	Name      string    `bson:"name"`
	Password  string    `bson:"password,omitempty"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
	Active    bool      `bson:"active"`
}

type repository struct {
	logger     *zap.Logger
	collection *mongo.Collection
	clock      clock.Clock
}

func NewRepository(logger *zap.Logger, db *mongo.Database, clk clock.Clock) Repository {
	return &repository{
		logger:     logger,
		collection: db.Collection(CollectionName),
		clock:      clk,
	}
}

func (r *repository) CreateCustomer(ctx context.Context, params CreateCustomerParams) (Customer, error) {
	now := r.clock.Now()
	c := Customer{
		Email:     params.Email,
		Name:      params.Name,
		Password:  params.Password,
		CreatedAt: now,
		UpdatedAt: now,
		Active:    true,
	}
	res, err := r.collection.InsertOne(ctx, c)
	if err != nil {
		if isDuplicateKeyError(err) {
			logctx.LoggerWithRequestInfo(ctx, r.logger).
				Warn("Customer already exists", zap.String("email", params.Email))
			return Customer{}, ErrCustomerAlreadyExists
		}
		logctx.LoggerWithRequestInfo(ctx, r.logger).Error("Failed to insert customer", zap.Error(err))
		return Customer{}, err
	}
	c.ID = res.InsertedID.(primitive.ObjectID).Hex()
	logctx.LoggerWithRequestInfo(ctx, r.logger).
		Info("Customer created successfully", zap.String("customer_id", c.ID))
	return c, nil
}

func (r *repository) FindByEmail(ctx context.Context, email string) (Customer, error) {
	var customer Customer
	filter := bson.M{
		FieldEmail:  email,
		FieldActive: true,
	}

	if err := r.collection.FindOne(ctx, filter).Decode(&customer); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			logctx.LoggerWithRequestInfo(ctx, r.logger).
				Warn("Customer not found", zap.String("email", email))
			return Customer{}, ErrCustomerNotFound
		}
		logctx.LoggerWithRequestInfo(ctx, r.logger).Error("Failed to find customer", zap.Error(err))
		return Customer{}, err
	}
	return customer, nil
}

// isDuplicateKeyError checks if the error is a duplicate key error (MongoDB error code 11000).
func isDuplicateKeyError(err error) bool {
	var we mongo.WriteException
	if errors.As(err, &we) {
		for _, e := range we.WriteErrors {
			if e.Code == 11000 {
				return true
			}
		}
	}
	return false
}
