package staff

import (
	"context"
	"time"
)

// Staff represents the structure of a restaurant staff user.
type Staff struct {
	ID           string    `bson:"_id,omitempty"`
	Email        string    `bson:"email"`
	RestaurantID string    `bson:"restaurant_id"`
	Owner        bool      `bson:"owner"`
	Name         string    `bson:"name"`
	Active       bool      `bson:"active"`
	Address      string    `bson:"address"`
	City         string    `bson:"city"`
	PostalCode   string    `bson:"postal_code"`
	CountryCode  string    `bson:"country_code"`
	CreatedAt    time.Time `bson:"created_at"`
	UpdatedAt    time.Time `bson:"updated_at"`
}

// Repository represents the interface for operations related to staff management.
//
//go:generate mockgen -destination=./mocks/repository_mock.go -package=staff_mocks github.com/alexgrauroca/practice-food-delivery-platform/services/restaurant-service/internal/staff Repository
type Repository interface {
	CreateStaff(ctx context.Context, params CreateStaffParams) (Staff, error)
	PurgeStaff(ctx context.Context, email string) error
}

// CreateStaffParams represents the input data required for creating a new staff user.
type CreateStaffParams struct {
	Email        string
	RestaurantID string
	Owner        bool
	Name         string
	Address      string
	City         string
	PostalCode   string
	CountryCode  string
}
