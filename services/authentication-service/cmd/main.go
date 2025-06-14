// Package main is the entry point for the authentication service.
// It initializes and coordinates core components.

package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"

	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/clock"
	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/config"
	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/customers"
	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/jwt"
	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/middleware"
	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/refresh"
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

	// Initialize features
	refreshService := initRefreshFeature(logger, db)
	jwtService := initJWTFeature(logger)
	initCustomersFeature(logger, db, router, refreshService, jwtService)

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

func initRefreshFeature(logger *zap.Logger, db *mongo.Database) refresh.Service {
	// Initialize the refresh repository
	repo := refresh.NewRepository(logger, db, clock.RealClock{})

	// Initialize the refresh service
	return refresh.NewService(logger, repo, clock.RealClock{})
}

func initJWTFeature(logger *zap.Logger) jwt.Service {
	/// Initialize the jwt service
	//TODO configure secret by env vars
	return jwt.NewService(logger, []byte("your-fancy-secret"))
}

func initCustomersFeature(logger *zap.Logger, db *mongo.Database, router *gin.Engine, refreshService refresh.Service,
	jwtService jwt.Service) {
	// Initialize the customer's repository
	repo := customers.NewRepository(logger, db, clock.RealClock{})

	// Initialize the customer's service
	service := customers.NewService(logger, repo, refreshService, jwtService)

	// Initialize the customer's handler and register routes
	handler := customers.NewHandler(logger, service)
	handler.RegisterRoutes(router)
}
