// Package staff contains the error types used by the staff service.
package staff

import "errors"

var (
	// ErrStaffAlreadyExists indicates that a staff with the same identifying details already exists in the system.
	ErrStaffAlreadyExists = errors.New("staff already exists")
	// ErrInvalidCredentials indicates that the provided credentials are invalid during authentication processes.
	ErrInvalidCredentials = errors.New("invalid credentials")
)
