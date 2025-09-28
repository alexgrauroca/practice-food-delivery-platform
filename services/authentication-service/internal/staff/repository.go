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
	// CollectionName is the name of the staff collection in the database
	CollectionName = "staff"

	// FieldEmail represents the field name used to store or query email addresses in the database.
	FieldEmail = "email"
	// FieldActive represents the field name used to indicate the active status of a staff in the database.
	FieldActive = "active"
)

type Staff struct {
	ID        string    `bson:"_id,omitempty"`
	StaffID   string    `bson:"staff_id"`
	Email     string    `bson:"email"`
	Active    bool      `bson:"active"`
	Name      string    `bson:"name"`
	Password  string    `bson:"password,omitempty"`
	CreatedAt time.Time `bson:"created_at,omitempty"`
	UpdatedAt time.Time `bson:"updated_at,omitempty"`
}

// Repository defines the interface for the staff repository.
//
//go:generate mockgen -destination=./mocks/repository_mock.go -package=staff_mocks github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/staff Repository
type Repository interface {
	CreateStaff(ctx context.Context, params CreateStaffParams) (Staff, error)
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

// CreateStaffParams represents the parameters required to create a new staff user.
type CreateStaffParams struct {
	StaffID  string `json:"staff_id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (r *repository) CreateStaff(ctx context.Context, params CreateStaffParams) (Staff, error) {
	logger := r.logger.WithContext(ctx)

	now := r.clock.Now()
	c := Staff{
		StaffID: params.StaffID,
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
			logger.Warn("Staff already exists", log.Field{Key: "email", Value: params.Email})
			return Staff{}, ErrStaffAlreadyExists
		}
		logger.Error("Failed to insert staff", err)
		return Staff{}, err
	}
	c.ID = res.InsertedID.(primitive.ObjectID).Hex()
	logger.Info("Staff created successfully", log.Field{Key: "staff_id", Value: c.ID})
	return c, nil
}
