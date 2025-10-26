package authcore

// TokenPairResponse represents the structure for holding both access and refresh tokens along with metadata.
type TokenPairResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"` // the number of seconds until the token expires
	TokenType    string `json:"token_type"`
}
