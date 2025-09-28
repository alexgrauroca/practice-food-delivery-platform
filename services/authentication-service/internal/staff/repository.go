package staff

import (
	"context"
	"time"
)

const (
	// CollectionName is the name of the staff collection in the database
	CollectionName = "staff"
)

type Staff struct {
	ID        string    `bson:"_id,omitempty"`
	StaffID   string    `bson:"staff_id"`
	Email     string    `bson:"email"`
	Active    bool      `bson:"active"`
	Name      string    `bson:"name"`
	Password  string    `bson:"password,omitempty"`
	CreatedAt time.Time `bson:"created_at,omitempty"`
	UpdatedAt time.Time `bson:"updated_at,omitempty"`
}

// Repository defines the interface for the staff repository.
//
//go:generate mockgen -destination=./mocks/repository_mock.go -package=staff_mocks github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/staff Repository
type Repository interface {
	CreateStaff(ctx context.Context, params CreateStaffParams) (Staff, error)
}

// CreateStaffParams represents the parameters required to create a new staff user.
type CreateStaffParams struct {
	StaffID  string `json:"staff_id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}
