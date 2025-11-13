package staff

import "errors"

var (
	// ErrStaffAlreadyExists indicates an error where a staff with the same identifier already exists in the system.
	ErrStaffAlreadyExists = errors.New("staff already exists")
	// ErrStaffNotFound indicates an error where a staff with the specified identifier does not exist in the system.
	ErrStaffNotFound = errors.New("staff not found")
)
