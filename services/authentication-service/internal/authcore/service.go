package authcore

// TokenPair represents a pair of tokens typically used for authentication and session management.
type TokenPair struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int // Number of seconds until the token expires
	TokenType    string
}
