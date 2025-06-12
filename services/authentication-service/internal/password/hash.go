// Package password provides password-related functionality for the authentication service.
// It implements industry-standard cryptographic methods for handling sensitive password data.
package password

import "golang.org/x/crypto/bcrypt"

const hashCost = 12 // Cost factor for bcrypt hashing

// Hash generates a bcrypt hashed version of the input password. Returns the hashed string or an error if hashing fails.
func Hash(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), hashCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

// Verify checks if the provided password matches the hashed password.
func Verify(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
