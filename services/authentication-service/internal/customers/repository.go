package customers

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

const customersCollection = "customers"

//go:generate mockgen -destination=./mocks/repository_mock.go -package=mocks github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/customers Repository
type Repository interface {
	CreateCustomer(ctx context.Context, params CreateCustomerParams) (Customer, error)
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
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
	Active    bool      `bson:"active"`
}

type repository struct {
	logger     *zap.Logger
	collection *mongo.Collection
}

func NewRepository(logger *zap.Logger, db *mongo.Database) Repository {
	return &repository{
		logger:     logger,
		collection: db.Collection(customersCollection),
	}
}

func (r *repository) CreateCustomer(ctx context.Context, params CreateCustomerParams) (Customer, error) {
	now := time.Now()

	c := Customer{
		Email:     params.Email,
		Name:      params.Name,
		CreatedAt: now,
		UpdatedAt: now,
		Active:    true,
	}
	res, err := r.collection.InsertOne(ctx, c)
	if err != nil {
		r.logger.Error("Failed to insert customer", zap.Error(err))
		return Customer{}, err
	}
	c.ID = res.InsertedID.(primitive.ObjectID).Hex()
	r.logger.Info("Customer created successfully", zap.String("customer_id", c.ID))
	return c, nil
}
