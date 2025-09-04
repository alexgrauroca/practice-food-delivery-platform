// Package main is the entry point for the authentication service.
// It initializes and coordinates core components.

package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/alexgrauroca/practice-food-delivery-platform/pkg/clock"
	"github.com/alexgrauroca/practice-food-delivery-platform/pkg/infraestructure/mongodb"
	customlog "github.com/alexgrauroca/practice-food-delivery-platform/pkg/log"
	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/customers"
	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/jwt"
	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/middleware"
	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/refresh"
)

func main() {
	ctx := context.Background()

	// Initialize the logger
	logger, err := customlog.NewProduction()
	if err != nil {
		log.Fatalf("can't initialize logger: %v", err)
		return
	}
	defer func(logger customlog.Logger) {
		if err := logger.Sync(); err != nil {
			log.Printf("failed to sync logger: %v", err)
		}
	}(logger)

	// Initialize the Gin router
	router := gin.Default()
	router.Use(middleware.RequestInfoMiddleware())

	// Initialize MongoDB connection
	client, err := mongodb.NewClient(ctx, logger)
	if err != nil {
		logger.Fatal("Failed to initialize MongoDB client", err)
		return
	}
	defer func(client *mongo.Client, ctx context.Context) {
		_ = client.Disconnect(ctx)
	}(client, ctx)

	db := client.Database("authentication_service")

	// Initialize features
	refreshService := initRefreshFeature(logger, db)
	jwtService := initJWTFeature(logger)
	initCustomersFeature(logger, db, router, refreshService, jwtService)

	logger.Info("Starting http server")
	// Start the server
	if err := router.Run(":8080"); err != nil {
		logger.Fatal("Failed to start server", err)
	}
}

func initRefreshFeature(logger customlog.Logger, db *mongo.Database) refresh.Service {
	// Initialize the refresh repository
	repo := refresh.NewRepository(logger, db, clock.RealClock{})

	// Initialize the refresh service
	return refresh.NewService(logger, repo, clock.RealClock{})
}

func initJWTFeature(logger customlog.Logger) jwt.Service {
	/// Initialize the jwt service
	//TODO configure secret by env vars
	return jwt.NewService(logger, []byte("a-string-secret-at-least-256-bits-long"))
}

func initCustomersFeature(logger customlog.Logger, db *mongo.Database, router *gin.Engine,
	refreshService refresh.Service,
	jwtService jwt.Service) {
	// Initialize the customer's repository
	repo := customers.NewRepository(logger, db, clock.RealClock{})

	// Initialize the customer's service
	service := customers.NewService(logger, repo, refreshService, jwtService)

	// Initialize the customer's handler and register routes
	handler := customers.NewHandler(logger, service)
	handler.RegisterRoutes(router)
}
