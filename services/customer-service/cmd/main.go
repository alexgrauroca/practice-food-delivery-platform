// Package main is the entry point for the customer service.
// It initializes and coordinates core components.

package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/alexgrauroca/practice-food-delivery-platform/pkg/clock"
	customhttp "github.com/alexgrauroca/practice-food-delivery-platform/pkg/http"
	"github.com/alexgrauroca/practice-food-delivery-platform/pkg/infraestructure/mongodb"

	"github.com/alexgrauroca/practice-food-delivery-platform/pkg/clients/authentication"
	customlog "github.com/alexgrauroca/practice-food-delivery-platform/pkg/log"
	"github.com/alexgrauroca/practice-food-delivery-platform/services/customer-service/internal/customers"
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
	router.Use(customhttp.RequestInfoMiddleware())

	// Initialize MongoDB connection
	client, err := mongodb.NewClient(ctx, logger)
	if err != nil {
		logger.Fatal("Failed to initialize MongoDB client", err)
		return
	}
	defer func(client *mongo.Client, ctx context.Context) {
		_ = client.Disconnect(ctx)
	}(client, ctx)

	db := client.Database("customer_service")

	// Initialize features
	authcli, authMiddleware, authctx := initAuthenticationFeature(logger)
	initCustomersFeature(logger, db, router, authcli, authMiddleware, authctx)

	logger.Info("Starting http server")
	// Start the server
	if err := router.Run(":8080"); err != nil {
		logger.Fatal("Failed to start server", err)
	}
}

func initAuthenticationFeature(logger customlog.Logger) (
	authentication.Client,
	authentication.Middleware,
	authentication.ContextReader,
) {
	authcli := authentication.NewClient(logger, authentication.Config{Debug: false})
	//TODO configure secret by env vars
	authService := authentication.NewService(logger, authcli, []byte("a-string-secret-at-least-256-bits-long"))
	authMiddleware := authentication.NewMiddleware(logger, authService)
	authctx := authentication.NewContextReader()

	return authcli, authMiddleware, authctx
}

func initCustomersFeature(
	logger customlog.Logger,
	db *mongo.Database,
	router *gin.Engine,
	authcli authentication.Client,
	authMiddleware authentication.Middleware,
	authctx authentication.ContextReader,
) {
	// Initialize the customer's repository
	repo := customers.NewRepository(logger, db, clock.RealClock{})

	// Initialize the customer's service
	service := customers.NewService(logger, repo, authcli, authctx)

	// Initialize the customer's handler and register routes
	handler := customers.NewHandler(logger, service, authMiddleware)
	handler.RegisterRoutes(router)
}
