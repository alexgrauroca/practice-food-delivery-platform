package main

import (
	"context"
	"log"

	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/clock"
	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/config"
	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/customers"
	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/middleware"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

func main() {
	// Initialize the logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
		return
	}
	defer func(logger *zap.Logger) {
		if err := logger.Sync(); err != nil && err.Error() != "sync /dev/stderr: invalid argument" {
			log.Printf("failed to sync logger: %v", err)
		}
	}(logger)

	// Initialize the Gin router
	router := gin.Default()
	router.Use(middleware.RequestInfoMiddleware())

	// Initialize MongoDB connection
	db, done := initMongoDB(logger)
	if !done {
		return
	}

	// Initialize the customers feature
	initCustomersFeature(logger, db, router)

	logger.Info("Starting http server")
	// Start the server
	if err := router.Run(":8080"); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}

func initMongoDB(logger *zap.Logger) (*mongo.Database, bool) {
	mongoCfg, err := config.LoadMongoConfig(logger)
	if err != nil {
		logger.Fatal("Failed to load MongoDB configuration", zap.Error(err))
		return nil, true
	}

	clientOpts := options.Client().ApplyURI(mongoCfg.URI)
	if mongoCfg.User != "" && mongoCfg.Password != "" {
		clientOpts.SetAuth(options.Credential{
			Username: mongoCfg.User,
			Password: mongoCfg.Password,
		})
	}

	client, err := mongo.Connect(context.Background(), clientOpts)
	if err != nil {
		logger.Fatal("Failed to connect to MongoDB", zap.Error(err))
		return nil, false
	}
	db := client.Database("authentication_service")
	logger.Info("Connected to MongoDB")
	return db, true
}

func initCustomersFeature(logger *zap.Logger, db *mongo.Database, router *gin.Engine) {
	// Initialize the customers repository
	repo := customers.NewRepository(logger, db, clock.RealClock{})

	// Initialize the customers service
	service := customers.NewService(logger, repo)

	// Initialize the customers handler and register routes
	handler := customers.NewHandler(logger, service)
	handler.RegisterRoutes(router)
}
