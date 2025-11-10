// Package restaurants contains the business logic for the restaurant service.
package restaurants

import "errors"

var (
	// ErrRestaurantAlreadyExists indicates an error where a restaurant with the same identifier already exists in
	//the system.
	ErrRestaurantAlreadyExists = errors.New("restaurant already exists")
	// ErrRestaurantNotFound indicates an error where a restaurant with the specified identifier does not exist in
	//the system.
	ErrRestaurantNotFound = errors.New("restaurant not found")
)
