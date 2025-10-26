// Package main is the entry point for the authentication service.
// It initializes and coordinates core components.

package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/alexgrauroca/practice-food-delivery-platform/pkg/auth"
	"github.com/alexgrauroca/practice-food-delivery-platform/pkg/clock"
	"github.com/alexgrauroca/practice-food-delivery-platform/pkg/infraestructure/mongodb"
	customlog "github.com/alexgrauroca/practice-food-delivery-platform/pkg/log"
	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/authcore"
	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/staff"

	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/customers"
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
	authService, authMiddleware, authctx := initAuthFeature(logger)
	authCoreService := initAuthCoreFeature(logger, authService, refreshService)
	initCustomersFeature(logger, db, router, refreshService, authService, authMiddleware, authctx)
	initStaffFeature(logger, db, router, authCoreService, authMiddleware, authctx)

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

func initAuthFeature(logger customlog.Logger) (
	auth.Service,
	auth.Middleware,
	auth.ContextReader,
) {
	/// Initialize the jwt service
	//TODO configure secret by env vars
	authService := auth.NewService(logger, []byte("a-string-secret-at-least-256-bits-long"), clock.RealClock{})
	authMiddleware := auth.NewMiddleware(logger, authService)
	authContextReader := auth.NewContextReader(logger)

	return authService, authMiddleware, authContextReader
}

func initAuthCoreFeature(
	logger customlog.Logger,
	authService auth.Service,
	refreshService refresh.Service,
) authcore.Service {
	return authcore.NewService(logger, authService, refreshService)
}

func initCustomersFeature(
	logger customlog.Logger,
	db *mongo.Database,
	router *gin.Engine,
	refreshService refresh.Service,
	authService auth.Service,
	authMiddleware auth.Middleware,
	authctx auth.ContextReader,
) {
	// Initialize the customer's repository
	repo := customers.NewRepository(logger, db, clock.RealClock{})

	// Initialize the customer's service
	service := customers.NewService(logger, repo, refreshService, authService, authctx)

	// Initialize the customer's handler and register routes
	handler := customers.NewHandler(logger, service, authMiddleware)
	handler.RegisterRoutes(router)
}

func initStaffFeature(
	logger customlog.Logger,
	db *mongo.Database,
	router *gin.Engine,
	authCoreService authcore.Service,
	authMiddleware auth.Middleware,
	authctx auth.ContextReader,
) {
	repo := staff.NewRepository(logger, db, clock.RealClock{})
	service := staff.NewService(logger, repo, authCoreService, authctx)
	handler := staff.NewHandler(logger, service, authMiddleware)
	handler.RegisterRoutes(router)
}
