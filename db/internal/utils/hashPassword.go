package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func CompareHashAndPassword(hash, userPassword string) error {
	// log.Printf("hash: %s", hash)
	// log.Printf("userpass: %s", userPassword)
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(userPassword)); err != nil {
		return err
	}
	// log.Printf("hash: %s", hash)
	// log.Printf("userpass: %s", userPassword)
	return nil
}
