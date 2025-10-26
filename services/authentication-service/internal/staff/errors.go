// Package staff contains the error types used by the staff service.
package staff

import "errors"

var (
	// ErrStaffAlreadyExists indicates that a staff with the same identifying details already exists in the system.
	ErrStaffAlreadyExists = errors.New("staff already exists")
	// ErrStaffNotFound indicates that a staff with the specified details could not be found in the system.
	ErrStaffNotFound = errors.New("staff not found")
)
