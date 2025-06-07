//go:build integration

package integration

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/config"
	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/customers"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

func setupTestDB(t *testing.T) (*mongo.Database, func()) {
	logger := zap.NewNop()
	mongoCfg, err := config.LoadMongoConfig(logger)
	if err != nil {
		t.Fatalf("Failed to load MongoDB configuration: %v", err)
	}

	clientOpts := options.Client().ApplyURI(mongoCfg.URI)
	if mongoCfg.User != "" && mongoCfg.Password != "" {
		clientOpts.SetAuth(options.Credential{
			Username: mongoCfg.User,
			Password: mongoCfg.Password,
		})
	}

	// Context with timeout for connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	// Setting up a unique database name for each test to avoid conflicts
	dbName := fmt.Sprintf("customers_test_authentication_service_%d_%d", time.Now().UnixNano(), rand.Intn(10000))
	db := client.Database(dbName)
	cleanup := func() {
		if err := db.Drop(ctx); err != nil {
			t.Fatalf("Failed to drop MongoDB collection: %v", err)
			return
		}
		if err := client.Disconnect(ctx); err != nil {
			t.Fatalf("Failed to disconnect MongoDB client: %v", err)
			return
		}
		cancel()
	}
	return db, cleanup
}

func setupTestCustomersCollection(t *testing.T, db *mongo.Database) *mongo.Collection {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	coll := db.Collection(customers.CustomersCollectionName)

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

func insertTestCustomer(t *testing.T, coll *mongo.Collection, email, name, password string, active bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	doc := bson.M{
		"email":      email,
		"name":       name,
		"password":   password,
		"created_at": time.Now(),
		"updated_at": time.Now(),
		"active":     active,
	}
	if _, err := coll.InsertOne(ctx, doc); err != nil {
		t.Fatalf("Failed to insert test customer: %v", err)
	}
}
