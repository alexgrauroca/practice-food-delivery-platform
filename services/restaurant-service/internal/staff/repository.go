package staff

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/alexgrauroca/practice-food-delivery-platform/pkg/clock"
	"github.com/alexgrauroca/practice-food-delivery-platform/pkg/infraestructure/mongodb"
	"github.com/alexgrauroca/practice-food-delivery-platform/pkg/log"
)

const (
	// CollectionName defines the name of the MongoDB collection used for storing staff documents.
	CollectionName = "staff"

	// FieldEmail represents the field name used to store or query email addresses in the database.
	FieldEmail = "email"
	// FieldActive represents the field name used to indicate the active status of a staff user in the database.
	FieldActive = "active"
	// FieldRestaurantID represents the field name used to store the unique RestaurantID of a staff user in the
	//database.
	FieldRestaurantID = "restaurant_id"
	// FieldOwner represents the field name used to indicate whether a staff user is a restaurant owner or not.
	FieldOwner = "owner"
)

// Staff represents the structure of a restaurant staff user.
type Staff struct {
	ID           string    `bson:"_id,omitempty"`
	Email        string    `bson:"email"`
	RestaurantID string    `bson:"restaurant_id"`
	Owner        bool      `bson:"owner"`
	Name         string    `bson:"name"`
	Active       bool      `bson:"active"`
	Address      string    `bson:"address"`
	City         string    `bson:"city"`
	PostalCode   string    `bson:"postal_code"`
	CountryCode  string    `bson:"country_code"`
	CreatedAt    time.Time `bson:"created_at"`
	UpdatedAt    time.Time `bson:"updated_at"`
}

// Repository represents the interface for operations related to staff management.
//
//go:generate mockgen -destination=./mocks/repository_mock.go -package=staff_mocks github.com/alexgrauroca/practice-food-delivery-platform/services/restaurant-service/internal/staff Repository
type Repository interface {
	CreateStaff(ctx context.Context, params CreateStaffParams) (Staff, error)
	PurgeStaff(ctx context.Context, email string) error
}

type repository struct {
	logger     log.Logger
	collection *mongo.Collection
	clock      clock.Clock
}

// NewRepository creates a new instance of the Repository interface with MongoDB implementation.
func NewRepository(logger log.Logger, db *mongo.Database, clk clock.Clock) Repository {
	return &repository{
		logger:     logger,
		collection: db.Collection(CollectionName),
		clock:      clk,
	}
}

// CreateStaffParams represents the input data required for creating a new staff user.
type CreateStaffParams struct {
	Email        string
	RestaurantID string
	Owner        bool
	Name         string
	Address      string
	City         string
	PostalCode   string
	CountryCode  string
}

func (r repository) CreateStaff(ctx context.Context, params CreateStaffParams) (Staff, error) {
	logger := r.logger.WithContext(ctx)

	now := r.clock.Now()
	staff := Staff{
		Email:        params.Email,
		RestaurantID: params.RestaurantID,
		Owner:        params.Owner,
		Name:         params.Name,
		Active:       true,
		Address:      params.Address,
		City:         params.City,
		PostalCode:   params.PostalCode,
		CountryCode:  params.CountryCode,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	res, err := r.collection.InsertOne(ctx, staff)
	if err != nil {
		if mongodb.IsDuplicateKeyError(err) {
			logger.Warn(
				"staff already exists",
				log.Field{Key: "email", Value: params.Email},
				log.Field{Key: "restaurant_id", Value: params.RestaurantID},
				log.Field{Key: "owner", Value: params.Owner},
			)
			return Staff{}, ErrStaffAlreadyExists
		}
		logger.Error("failed to insert staff", err)
		return Staff{}, err
	}

	staff.ID = res.InsertedID.(primitive.ObjectID).Hex()
	logger.Info("staff created successfully", log.Field{Key: "staff_id", Value: staff.ID})
	return staff, nil
}

func (r repository) PurgeStaff(ctx context.Context, email string) error {
	//TODO implement me
	panic("implement me")
}
