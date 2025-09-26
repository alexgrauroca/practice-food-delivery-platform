package staff

import (
	"context"
	"time"
)

// Service defines the interface for the staff service.
//
//go:generate mockgen -destination=./mocks/service_mock.go -package=staff_mocks github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/staff Service
type Service interface {
	RegisterStaff(ctx context.Context, input RegisterStaffInput) (RegisterStaffOutput, error)
}

// RegisterStaffInput defines the input structure required for registering a new staff.
type RegisterStaffInput struct {
	StaffID  string
	Email    string
	Password string
	Name     string
}

// RegisterStaffOutput represents the output data returned after successfully registering a new staff.
type RegisterStaffOutput struct {
	ID        string
	Email     string
	Name      string
	CreatedAt time.Time
}
