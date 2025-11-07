// Package restaurants contains the business logic for the restaurant service.
package restaurants

import "errors"

var (
	// ErrRestaurantAlreadyExists indicates an error where a restaurant with the same identifier already exists in
	//the system.
	ErrRestaurantAlreadyExists = errors.New("restaurant already exists")
)
