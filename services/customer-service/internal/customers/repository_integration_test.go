//go:build integration

package customers_test

import (
	"context"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/alexgrauroca/practice-food-delivery-platform/services/customer-service/internal/customers"
)

type customersRepositoryTestCase[P, W any] struct {
	name            string
	insertDocuments func(t *testing.T, coll *mongo.Collection)
	params          P
	want            W
	wantErr         error
}

func setupTestCustomersCollection(t *testing.T, db *mongo.Database) *mongo.Collection {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	coll := db.Collection(customers.CollectionName)

	// Create unique index on email
	indexModel := mongo.IndexModel{
		Keys: bson.D{{Key: customers.FieldEmail, Value: 1}},
		Options: options.Index().
			SetUnique(true).
			SetPartialFilterExpression(bson.D{{Key: customers.FieldActive, Value: true}}),
	}
	if _, err := coll.Indexes().CreateOne(ctx, indexModel); err != nil {
		t.Fatalf("Failed to create unique index: %v", err)
	}

	return coll
}
