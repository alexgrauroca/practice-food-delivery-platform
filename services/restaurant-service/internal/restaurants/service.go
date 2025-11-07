package restaurants

import (
	"context"
)

// Service represents the interface for operations related to restaurant management.
//
//go:generate mockgen -destination=./mocks/service_mock.go -package=restaurants_mocks github.com/alexgrauroca/practice-food-delivery-platform/services/restaurant-service/internal/restaurants Service
type Service interface {
	RegisterRestaurant(ctx context.Context, input RegisterRestaurantInput) (RegisterRestaurantOutput, error)
}
