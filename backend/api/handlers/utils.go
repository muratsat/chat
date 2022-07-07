package handlers

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// generate password hash
func generatePasswordHash(password string) (string, error) {
	// generate salt
	salt, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	// generate hash
	hash := string(salt)
	return hash, nil
}

// validate authentication request
func (auth *authRequest) validate() error {
	username := auth.Username

	// check if username is empty
	if username == "" {
		return fmt.Errorf("Username is empty")
	}

	// check if username is too long
	if len(username) > 25 {
		return fmt.Errorf("Username is too long")
	}

	// check if username has invalid characters
	for _, char := range username {
		if !(char >= 'a' && char <= 'z' || char >= 'A' && char <= 'Z' || char >= '0' && char <= '9' || char == '_') {
			return fmt.Errorf("Username contains invalid characters")
		}
	}

	// check if password is empty
	if auth.Password == "" {
		return fmt.Errorf("Password is empty")
	}

	return nil
}
