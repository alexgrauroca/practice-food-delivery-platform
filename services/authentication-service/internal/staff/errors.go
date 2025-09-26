package staff

import "errors"

var (
	// ErrStaffAlreadyExists indicates that a staff with the same identifying details already exists in the system.
	ErrStaffAlreadyExists = errors.New("staff already exists")
)
