package customers

import "errors"

var (
	ErrCustomerAlreadyExists = errors.New("customer already exists")
	ErrInvalidCredentials    = errors.New("invalid credentials")
	ErrCustomerNotFound      = errors.New("customer not found")
)
