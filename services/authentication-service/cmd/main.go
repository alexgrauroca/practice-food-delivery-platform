package main

import (
	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/customers"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			log.Fatalf("failed to sync logger: %v", err)
		}
	}(logger)

	// Initialize the Gin router
	router := gin.Default()

	// Initialize the customers handler and register routes
	handler := customers.NewHandler(logger)
	handler.RegisterRoutes(router)

	if err := router.Run(":8080"); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}
