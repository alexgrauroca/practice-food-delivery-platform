package authentication

import "errors"

// ErrAccessTokenRequired represents an error when a required access token is missing or not provided.
var ErrAccessTokenRequired = errors.New("access token required")
