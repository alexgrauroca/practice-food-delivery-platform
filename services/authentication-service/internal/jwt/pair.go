package jwt

// TokenPair represents a pair of tokens consisting of an access token and a refresh token.
type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

// GenerateTokenPair generates an access token and a refresh token as a token pair for authentication purposes.
func GenerateTokenPair(id string, cfg Config) (TokenPair, error) {
	accessToken, err := GenerateToken(id, cfg)
	if err != nil {
		return TokenPair{}, err
	}
	refreshToken, err := GenerateRefreshToken()
	if err != nil {
		return TokenPair{}, err
	}
	return TokenPair{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}
