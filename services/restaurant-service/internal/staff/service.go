package staff

// Service represents the interface defining business operations related to staff management.
//
//go:generate mockgen -destination=./mocks/service_mock.go -package=staff_mocks github.com/alexgrauroca/practice-food-delivery-platform/services/restaurant-service/internal/staff Service
type Service interface {
	RegisterStaff(input RegisterStaffInput) (RegisterStaffOutput, error)
}

// RegisterStaffInput represents the input data required for registering a new staff member.
type RegisterStaffInput struct{}

// RegisterStaffOutput represents the output data returned after successfully registering a new staff member.
type RegisterStaffOutput struct{}
