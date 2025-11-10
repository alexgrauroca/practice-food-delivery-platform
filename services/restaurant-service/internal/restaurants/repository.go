package restaurants

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/alexgrauroca/practice-food-delivery-platform/pkg/clock"
	"github.com/alexgrauroca/practice-food-delivery-platform/pkg/infraestructure/mongodb"
	"github.com/alexgrauroca/practice-food-delivery-platform/pkg/log"
)

const (
	// CollectionName defines the name of the MongoDB collection used for storing restaurant documents.
	CollectionName = "restaurants"

	// FieldVatCode represents the field name used to store the VAT code of a restaurant in the database.
	FieldVatCode = "vat_code"
	// FieldActive represents the field name used to indicate the active status of a restaurant in the database.
	FieldActive = "active"
)

// Restaurant represents a restaurant.
type Restaurant struct {
	ID         string    `bson:"_id,omitempty"`
	VatCode    string    `bson:"vat_code"`
	Name       string    `bson:"name"`
	LegalName  string    `bson:"legal_name"`
	TaxID      string    `bson:"tax_id"`
	TimezoneID string    `bson:"timezone_id"`
	Contact    Contact   `bson:"contact"`
	Active     bool      `bson:"active"`
	CreatedAt  time.Time `bson:"created_at"`
	UpdatedAt  time.Time `bson:"updated_at"`
}

// Contact represents the contact details of a restaurant.
type Contact struct {
	PhonePrefix string `bson:"phone_prefix"`
	PhoneNumber string `bson:"phone_number"`
	Email       string `bson:"email"`
	Address     string `bson:"address"`
	City        string `bson:"city"`
	PostalCode  string `bson:"postal_code"`
	CountryCode string `bson:"country_code"`
}

// Repository represents the interface for operations related to restaurant management.
//
//go:generate mockgen -destination=./mocks/repository_mock.go -package=restaurants_mocks github.com/alexgrauroca/practice-food-delivery-platform/services/restaurant-service/internal/restaurants Repository
type Repository interface {
	CreateRestaurant(ctx context.Context, params CreateRestaurantParams) (Restaurant, error)
	PurgeRestaurant(ctx context.Context, vatCode string) error
}

type repository struct {
	logger     log.Logger
	collection *mongo.Collection
	clock      clock.Clock
}

// NewRepository initializes and returns a new Repository instance.
func NewRepository(logger log.Logger, db *mongo.Database, clk clock.Clock) Repository {
	return &repository{
		logger:     logger,
		collection: db.Collection(CollectionName),
		clock:      clk,
	}
}

// CreateRestaurantParams represents the parameters for creating a restaurant.
type CreateRestaurantParams struct {
	VatCode    string
	Name       string
	LegalName  string
	TaxID      string
	TimezoneID string
	Contact    CreateContactParams
}

// CreateContactParams represents the parameters for creating a restaurant's contact information.'
type CreateContactParams struct {
	PhonePrefix string
	PhoneNumber string
	Email       string
	Address     string
	City        string
	PostalCode  string
	CountryCode string
}

func (r repository) CreateRestaurant(ctx context.Context, params CreateRestaurantParams) (Restaurant, error) {
	logger := r.logger.WithContext(ctx)

	now := r.clock.Now()
	c := Restaurant{
		VatCode:    params.VatCode,
		Name:       params.Name,
		LegalName:  params.LegalName,
		TaxID:      params.TaxID,
		TimezoneID: params.TimezoneID,
		Contact: Contact{
			PhonePrefix: params.Contact.PhonePrefix,
			PhoneNumber: params.Contact.PhoneNumber,
			Email:       params.Contact.Email,
			Address:     params.Contact.Address,
			City:        params.Contact.City,
			PostalCode:  params.Contact.PostalCode,
			CountryCode: params.Contact.CountryCode,
		},
		Active:    true,
		CreatedAt: now,
		UpdatedAt: now,
	}
	res, err := r.collection.InsertOne(ctx, c)
	if err != nil {
		if mongodb.IsDuplicateKeyError(err) {
			logger.Warn("restaurant already exists", log.Field{Key: "vat_code", Value: params.VatCode})
			return Restaurant{}, ErrRestaurantAlreadyExists
		}
		logger.Error("Failed to insert restaurant", err)
		return Restaurant{}, err
	}
	c.ID = res.InsertedID.(primitive.ObjectID).Hex()
	logger.Info("Restaurant created successfully", log.Field{Key: "restaurant_id", Value: c.ID})
	return c, nil
}

func (r repository) PurgeRestaurant(ctx context.Context, vatCode string) error {
	logger := r.logger.WithContext(ctx)
	logger.Info("Purging restaurant", log.Field{Key: "vat_code", Value: vatCode})

	res, err := r.collection.DeleteOne(ctx, bson.M{
		FieldVatCode: vatCode,
		FieldActive:  true,
	})
	if err != nil {
		logger.Error("Failed to purge restaurant", err)
		return err
	}
	if res.DeletedCount == 0 {
		logger.Warn("Restaurant not found", log.Field{Key: "vat_code", Value: vatCode})
		return ErrRestaurantNotFound
	}

	logger.Info("Restaurant purged successfully", log.Field{Key: "vat_code", Value: vatCode})
	return nil
}
