// Package mongodb provides MongoDB client implementation and database operations for authentication service,
// including configuration, connection management, and database operations for customer data storage.
package mongodb

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"

	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/config"
)

// NewClient creates and returns a new MongoDB Client with the given context and logger. It applies authentication if required.
// Returns an error if loading configuration or connecting to MongoDB fails.
func NewClient(ctx context.Context, logger *zap.Logger) (*mongo.Client, error) {
	mongoCfg, err := config.LoadMongoConfig(logger)
	if err != nil {
		return nil, fmt.Errorf("failed to load MongoDB configuration: %w", err)
	}

	clientOpts := options.Client().ApplyURI(mongoCfg.URI)
	if mongoCfg.User != "" && mongoCfg.Password != "" {
		clientOpts.SetAuth(options.Credential{
			Username: mongoCfg.User,
			Password: mongoCfg.Password,
		})
	}

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	return client, nil
}

// IsDuplicateKeyError checks if the error is a duplicate key error (MongoDB error code 11000).
func IsDuplicateKeyError(err error) bool {
	var we mongo.WriteException
	if errors.As(err, &we) {
		for _, e := range we.WriteErrors {
			if e.Code == 11000 {
				return true
			}
		}
	}
	return false
}
