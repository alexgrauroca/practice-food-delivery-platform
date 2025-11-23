package auth

import "github.com/golang-jwt/jwt/v5"

// Role represents a user role within the authentication system
type Role string

const (
	// RoleCustomer represents the role assigned to authenticated customers
	RoleCustomer Role = "customer"
)

// Claims represent the authentication claims
type Claims struct {
	jwt.RegisteredClaims
	Role   string `json:"role"`
	Tenant string `json:"tenant"`
}
