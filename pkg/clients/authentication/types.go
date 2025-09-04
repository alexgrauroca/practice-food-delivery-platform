package authentication

// Role represents a user role within the authentication system
type Role string

const (
	// RoleCustomer represents the role assigned to authenticated customers
	RoleCustomer Role = "customer"
)
