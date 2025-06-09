package password

import "golang.org/x/crypto/bcrypt"

const hashCost = 12 // Cost factor for bcrypt hashing

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
