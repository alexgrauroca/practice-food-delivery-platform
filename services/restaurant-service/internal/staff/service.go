package staff

import (
	"context"
	"time"
)

// Service represents the interface defining business operations related to staff management.
//
//go:generate mockgen -destination=./mocks/service_mock.go -package=staff_mocks github.com/alexgrauroca/practice-food-delivery-platform/services/restaurant-service/internal/staff Service
type Service interface {
	RegisterStaffOwner(ctx context.Context, input RegisterStaffOwnerInput) (RegisterStaffOwnerOutput, error)
}

// RegisterStaffOwnerInput represents the input data required for registering a new staff member.
type RegisterStaffOwnerInput struct {
	Email        string
	Password     string
	RestaurantID string
	Name         string
	Address      string
	City         string
	PostalCode   string
	CountryCode  string
}

// RegisterStaffOwnerOutput represents the output data returned after successfully registering a new staff member.
type RegisterStaffOwnerOutput struct {
	ID           string
	Email        string
	RestaurantID string
	Owner        bool
	Name         string
	Address      string
	City         string
	PostalCode   string
	CountryCode  string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
