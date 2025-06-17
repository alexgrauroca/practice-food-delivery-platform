package mongodb

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
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
func NewTestDB(t *testing.T) *TestDB {
	t.Helper()

	ctx := context.Background()
	client, err := NewClient(ctx, zap.NewNop())
	if err != nil {
		t.Fatalf("Failed to create MongoDB client: %v", err)
	}

	// Setting up a unique database name for each test to avoid conflicts
	dbName := fmt.Sprintf("customers_test_authentication_service_%d_%d", time.Now().UnixNano(), rand.Intn(10000))
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

