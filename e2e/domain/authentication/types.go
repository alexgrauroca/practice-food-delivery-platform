// Package authentication handles identity and access management for all users (Customers, Staff,
// Couriers) across the platform
package authentication

// Token represents an authentication token containing access and refresh tokens, expiration time, and token type.
type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}
