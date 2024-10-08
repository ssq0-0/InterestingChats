package utils

import (
	"InterestingChats/backend/user_services/internal/consts"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword accepts a password and returns its hash created with bcrypt.
// If an error occurs while creating the hash, it is returned along with an empty string.
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf(consts.InternalGenerateHash, err)
	}
	stringPassword := string(hash)

	return stringPassword, nil
}
