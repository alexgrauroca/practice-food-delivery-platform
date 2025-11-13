// Package staff provides functionality for managing restaurant staff members,
// including creation, retrieval, and management of staff accounts. It handles
// staff-related operations such as authentication, authorization, and maintains
// staff records in the system.
package staff

import "errors"

var (
	// ErrStaffAlreadyExists indicates an error where a staff with the same identifier already exists in the system.
	ErrStaffAlreadyExists = errors.New("staff already exists")
	// ErrStaffNotFound indicates an error where a staff with the specified identifier does not exist in the system.
	ErrStaffNotFound = errors.New("staff not found")
)
