package customers

import "errors"

var (
	ErrCustomerAlreadyExists = errors.New("customer already exists")
	ErrInvalidCredentials    = errors.New("invalid credentials")
)
