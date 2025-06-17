package refresh

import "errors"

var (
	// ErrRefreshTokenNotFound indicates that the specified refresh token could not be found.
	ErrRefreshTokenNotFound = errors.New("refresh token not found")
	// ErrRefreshTokenAlreadyExists represents an error indicating a refresh token already exists in the system.
	ErrRefreshTokenAlreadyExists = errors.New("refresh token already exists")
)
