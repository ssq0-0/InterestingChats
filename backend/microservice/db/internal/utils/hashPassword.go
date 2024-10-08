package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// CompareHashAndPassword to compare the given password with the hashed password.
func CompareHashAndPassword(hash, userPassword string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(userPassword)); err != nil {
		return err
	}
	return nil
}
