//go:build integration

package customers_test

import (
	"testing"

	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

func setupTestDB(t *testing.T) (mongo.Database, func()) {
	logger := zap.NewNop()
	mongoCfg := config.LoadMongoConfig(logger)
	if err != nil {
		panic("Failed to load MongoDB configuration: " + err.Error())
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
		panic("Failed to connect to MongoDB: " + err.Error())
	}
	db := client.Database("customers_test_authentication_service")
	cleanup := func() {
		db.Drop(ctx)
		client.Disconnect(ctx)
		cancel()
	}
	return db, true
}

func TestRepository_CreateCustomer(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "when exists an active customer with the same email, it should return a customer already exists error",
		},
		{
			name: "when there is an error creating the customer, it should propagate the error",
		},
		{
			name: "when the customer is created successfully, it should return the created customer",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

		})
	}
}
