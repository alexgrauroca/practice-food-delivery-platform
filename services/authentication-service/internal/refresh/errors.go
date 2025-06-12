package refresh

import "errors"

// ErrRefreshTokenNotFound indicates that the specified refresh token could not be found.
var ErrRefreshTokenNotFound = errors.New("refresh token not found")
