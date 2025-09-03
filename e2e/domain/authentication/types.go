// Package authentication handles identity and access management for all users (Customers, Staff,
// Couriers) across the platform
package authentication

import "github.com/golang-jwt/jwt/v5"

// Token represents an authentication token containing access and refresh tokens, expiration time, and token type.
type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

// Claims represent the custom claims structure for JWT tokens
type Claims struct {
	jwt.RegisteredClaims
	Role string `json:"role"`
}

// LoginRequest represents the payload needed for user authentication, containing email and password.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse contains the authentication token details returned after successful login.
type LoginResponse struct {
	Token
}

// RefreshRequest contains the current access and refresh tokens needed to get new tokens.
type RefreshRequest struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// RefreshResponse contains the new authentication token details returned after successful token refresh.
type RefreshResponse struct {
	Token
}
