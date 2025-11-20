package mongodb

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/alexgrauroca/practice-food-delivery-platform/pkg/log"
)

// TestDB represents a test database wrapper for managing a MongoDB instance in a testing environment.
// Client is the MongoDB client connected to the test database.
// DB holds a reference to the MongoDB database being used for testing.
// cleanup is a function used to clean up and release resources after testing.
type TestDB struct {
	Client  *mongo.Client
	DB      *mongo.Database
	cleanup func() error
}

// NewTestDB creates and returns a new TestDB instance for testing purposes with an isolated MongoDB database.
// It generates a unique database name to avoid conflicts and initializes cleanup logic for resource management.
// The function requires a *testing.T instance for logging and error handling during test execution.
func NewTestDB(t *testing.T, dbPrefix string) *TestDB {
	t.Helper()

	ctx := context.Background()
	logger, _ := log.NewTest()
	client, err := NewClient(ctx, logger)
	if err != nil {
		t.Fatalf("Failed to create MongoDB client: %v", err)
	}

	// Setting up a unique database name for each test to avoid conflicts
	dbName := fmt.Sprintf("%s_%d_%d", dbPrefix, time.Now().UnixNano(), rand.Intn(10000))
	db := client.Database(dbName)

	return &TestDB{
		Client: client,
		DB:     db,
		cleanup: func() error {
			ctx := context.Background()
			if err := db.Drop(ctx); err != nil {
				return fmt.Errorf("failed to drop test database: %w", err)
			}
			if err := client.Disconnect(ctx); err != nil {
				return fmt.Errorf("failed to disconnect test client: %w", err)
			}
			return nil
		},
	}
}

// Close releases resources associated with the TestDB by invoking the cleanup function. Logs errors if cleanup fails.
func (tdb *TestDB) Close(t *testing.T) {
	t.Helper()
	if err := tdb.cleanup(); err != nil {
		t.Errorf("Failed to cleanup test database: %v", err)
	}
}

// InsertTestDocument inserts a test document into a Mongo collection for testing purposes.
// This function marshals the provided document, converts it to bson.M, and inserts it into the specified collection.
// It uses a timeout context of 5 seconds and fails the test if any errors occur during the process.
func InsertTestDocument(t *testing.T, coll *mongo.Collection, doc any) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	data, err := bson.Marshal(doc)
	if err != nil {
		t.Fatalf("Failed to marshal test refresh token: %v", err)
	}

	var bdoc bson.M
	if err := bson.Unmarshal(data, &bdoc); err != nil {
		t.Fatalf("Failed to unmarshal test refresh token: %v", err)
	}

	// Check if there's an _id field with a string value and convert it to ObjectID
	if idStr, ok := bdoc["_id"].(string); ok && idStr != "" {
		id, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			t.Fatalf("Failed to convert string ID to ObjectID: %v", err)
		}
		bdoc["_id"] = id
	}

	if _, err := coll.InsertOne(ctx, bdoc); err != nil {
		t.Fatalf("Failed to insert test refresh token: %v", err)
	}
}
