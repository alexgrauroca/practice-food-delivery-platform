// Package main is the entry point for the restaurant service.
// It initializes and coordinates core components.

package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/alexgrauroca/practice-food-delivery-platform/pkg/clients/authentication"
	"github.com/alexgrauroca/practice-food-delivery-platform/pkg/clock"
	customhttp "github.com/alexgrauroca/practice-food-delivery-platform/pkg/http"
	"github.com/alexgrauroca/practice-food-delivery-platform/pkg/infraestructure/mongodb"
	customlog "github.com/alexgrauroca/practice-food-delivery-platform/pkg/log"
	"github.com/alexgrauroca/practice-food-delivery-platform/services/restaurant-service/internal/restaurants"
	"github.com/alexgrauroca/practice-food-delivery-platform/services/restaurant-service/internal/staff"
)

const dbName = "restaurant_service"

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
	defer func(ctx context.Context, client *mongo.Client) {
		_ = client.Disconnect(ctx)
	}(ctx, client)

	db := client.Database(dbName)

	// Initialize features
	authcli := initAuthenticationFeature(logger)
	staffService := initStaffFeature(logger, db, authcli)
	initRestaurantsFeature(router, logger, db, staffService)

	logger.Info("Starting http server")
	// Start the server
	if err := router.Run(":8080"); err != nil {
		logger.Fatal("Failed to start server", err)
	}
}

func initAuthenticationFeature(logger customlog.Logger) authentication.Client {
	return authentication.NewClient(logger, authentication.Config{Debug: false})
}

func initStaffFeature(logger customlog.Logger, db *mongo.Database, authcli authentication.Client) staff.Service {
	repo := staff.NewRepository(logger, db, clock.RealClock{})
	return staff.NewService(logger, repo, authcli)
}

func initRestaurantsFeature(
	router *gin.Engine,
	logger customlog.Logger,
	db *mongo.Database,
	staffService staff.Service,
) {
	repo := restaurants.NewRepository(logger, db, clock.RealClock{})
	service := restaurants.NewService(logger, repo, staffService)
	handler := restaurants.NewHandler(logger, service)
	handler.RegisterRoutes(router)
}
