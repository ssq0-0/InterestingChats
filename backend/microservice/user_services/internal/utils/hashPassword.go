package utils

import (
	"InterestingChats/backend/user_services/internal/consts"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf(consts.InternalGenerateHash, err)
	}
	stringPassword := string(hash)

	return stringPassword, nil
}
