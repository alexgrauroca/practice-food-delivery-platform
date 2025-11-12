//go:build integration

package staff_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/alexgrauroca/practice-food-delivery-platform/pkg/clock"
	"github.com/alexgrauroca/practice-food-delivery-platform/pkg/infraestructure/mongodb"
	"github.com/alexgrauroca/practice-food-delivery-platform/pkg/log"
	"github.com/alexgrauroca/practice-food-delivery-platform/services/restaurant-service/internal/staff"
)

const dbPrefix = "staff_test"

type repoTestCase[P, W any] struct {
	name            string
	insertDocuments func(t *testing.T, coll *mongo.Collection)
	params          P
	want            W
	wantErr         error
}

func TestRepository_CreateStaff(t *testing.T) {
	now := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	logger, _ := log.NewTest()

	tests := []repoTestCase[staff.CreateStaffParams, staff.Staff]{
		{
			name: "when there is an active staff with the same email and restaurant, " +
				"then it returns a staff already exists error",
			params: staff.CreateStaffParams{
				Email:        "test@example.com",
				RestaurantID: "fake-restaurant-id",
				Owner:        false,
			},
			insertDocuments: func(t *testing.T, coll *mongo.Collection) {
				mongodb.InsertTestDocument(t, coll, staff.Staff{
					Email:        "test@example.com",
					Active:       true,
					RestaurantID: "fake-restaurant-id",
					Owner:        false,
				})
			},
			want:    staff.Staff{},
			wantErr: staff.ErrStaffAlreadyExists,
		},
		{
			name: "when there is an active staff owner with the same restaurant, " +
				"then it returns a staff already exists error",
			params: staff.CreateStaffParams{
				Email:        "owner-test@example.com",
				RestaurantID: "fake-restaurant-id",
				Owner:        true,
			},
			insertDocuments: func(t *testing.T, coll *mongo.Collection) {
				mongodb.InsertTestDocument(t, coll, staff.Staff{
					Email:        "another-owner-test@example.com",
					Active:       true,
					RestaurantID: "fake-restaurant-id",
					Owner:        true,
				})
			},
			want:    staff.Staff{},
			wantErr: staff.ErrStaffAlreadyExists,
		},
		{
			name: "when there is an active staff with the same email but different restaurant, " +
				"then it creates the staff",
			params: staff.CreateStaffParams{
				Email:        "test@example.com",
				RestaurantID: "fake-restaurant-id",
				Owner:        false,
				Name:         "Test Staff",
				Address:      "123 Main St",
				City:         "London",
				PostalCode:   "SW1A 1AA",
				CountryCode:  "GB",
			},
			insertDocuments: func(t *testing.T, coll *mongo.Collection) {
				mongodb.InsertTestDocument(t, coll, staff.Staff{
					Email:        "test@example.com",
					Active:       true,
					RestaurantID: "another-fake-restaurant-id",
					Owner:        false,
				})

				mongodb.InsertTestDocument(t, coll, staff.Staff{
					Email:        "owner-test@example.com",
					Active:       true,
					RestaurantID: "fake-restaurant-id",
					Owner:        true,
				})

				mongodb.InsertTestDocument(t, coll, staff.Staff{
					Email:        "test@example.com",
					Active:       false,
					RestaurantID: "fake-restaurant-id",
					Owner:        false,
				})
			},
			want: staff.Staff{
				Email:        "test@example.com",
				RestaurantID: "fake-restaurant-id",
				Owner:        false,
				Active:       true,
				Name:         "Test Staff",
				Address:      "123 Main St",
				City:         "London",
				PostalCode:   "SW1A 1AA",
				CountryCode:  "GB",
				CreatedAt:    now,
				UpdatedAt:    now,
			},
		},
		{
			name: "when there is not an active staff owner with the same restaurant, then it creates the staff owner",
			params: staff.CreateStaffParams{
				Email:        "owner-test@example.com",
				RestaurantID: "fake-restaurant-id",
				Owner:        true,
				Name:         "Test Staff",
				Address:      "123 Main St",
				City:         "London",
				PostalCode:   "SW1A 1AA",
				CountryCode:  "GB",
			},
			insertDocuments: func(t *testing.T, coll *mongo.Collection) {
				mongodb.InsertTestDocument(t, coll, staff.Staff{
					Email:        "test@example.com",
					Active:       true,
					RestaurantID: "fake-restaurant-id",
					Owner:        false,
				})

				mongodb.InsertTestDocument(t, coll, staff.Staff{
					Email:        "owner-test@example.com",
					Active:       true,
					RestaurantID: "another-fake-restaurant-id",
					Owner:        true,
				})

				mongodb.InsertTestDocument(t, coll, staff.Staff{
					Email:        "owner-test@example.com",
					Active:       false,
					RestaurantID: "fake-restaurant-id",
					Owner:        true,
				})
			},
			want: staff.Staff{
				Email:        "owner-test@example.com",
				RestaurantID: "fake-restaurant-id",
				Owner:        true,
				Active:       true,
				Name:         "Test Staff",
				Address:      "123 Main St",
				City:         "London",
				PostalCode:   "SW1A 1AA",
				CountryCode:  "GB",
				CreatedAt:    now,
				UpdatedAt:    now,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tdb := mongodb.NewTestDB(t, dbPrefix)
			defer tdb.Close(t)

			coll := setupTestStaffCollection(t, tdb.DB)
			if tt.insertDocuments != nil {
				tt.insertDocuments(t, coll)
			}

			repo := staff.NewRepository(logger, tdb.DB, clock.FixedClock{FixedTime: now})
			got, err := repo.CreateStaff(context.Background(), tt.params)

			// Error assertion
			assert.ErrorIs(t, err, tt.wantErr)

			// Validating the got only if there is no error expected
			if tt.wantErr == nil {
				// As the StaffID is generated by MongoDB, we just check that it is not empty
				assert.NotEmpty(t, got.ID, "StaffID should not be empty")

				// Doing this as in that way, I can do a direct equal assertion between the want and got
				tt.want.ID = got.ID
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestRepository_CreateStaff_UnexpectedFailure(t *testing.T) {
	now := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	logger, _ := log.NewTest()

	tdb := mongodb.NewTestDB(t, "customers_test_authentication_service")
	repo := staff.NewRepository(logger, tdb.DB, clock.FixedClock{FixedTime: now})

	// Simulating an unexpected failure by closing the opened connection
	tdb.Close(t)

	_, err := repo.CreateStaff(context.Background(), staff.CreateStaffParams{})
	assert.Error(t, err, "Expected an error due to unexpected failure")
	assert.NotErrorIs(t, err, staff.ErrStaffAlreadyExists)
}

func TestRepository_PurgeStaff(t *testing.T) {}

func TestRepository_PurgeStaff_UnexpectedFailure(t *testing.T) {}

func setupTestStaffCollection(t *testing.T, db *mongo.Database) *mongo.Collection {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	coll := db.Collection(staff.CollectionName)

	// Email must be unique per restaurant with active staff users
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: staff.FieldRestaurantID, Value: 1},
			{Key: staff.FieldEmail, Value: 1},
		},
		Options: options.Index().
			SetUnique(true).
			SetPartialFilterExpression(bson.D{{Key: staff.FieldActive, Value: true}}),
	}
	if _, err := coll.Indexes().CreateOne(ctx, indexModel); err != nil {
		t.Fatalf("Failed to create unique index: %v", err)
	}

	// Staff owner must be unique per restaurant
	indexModel = mongo.IndexModel{
		Keys: bson.D{
			{Key: staff.FieldRestaurantID, Value: 1},
		},
		Options: options.Index().
			SetUnique(true).
			SetPartialFilterExpression(bson.D{
				{Key: staff.FieldActive, Value: true},
				{Key: staff.FieldOwner, Value: true},
			}),
	}
	if _, err := coll.Indexes().CreateOne(ctx, indexModel); err != nil {
		t.Fatalf("Failed to create unique index: %v", err)
	}

	return coll
}
