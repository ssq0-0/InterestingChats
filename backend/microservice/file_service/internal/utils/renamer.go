package utils

import (
	"crypto/rand"
	"encoding/base64"
	"file_service/internal/consts"
	"fmt"
	"time"
)

// RenameFile generates a new file name using a timestamp and a random string.
func RenameFile(filename string) (string, string, error) {
	str, err := generateRandomString(8)
	if err != nil {
		return "", consts.ErrFailedSaveImage, err
	}

	return fmt.Sprintf("%s_%s_%s", time.Now().Format("20060102150405"), str, filename), "", nil
}

// generateRandomString generates a random string of specified length.
func generateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	encoded := base64.URLEncoding.EncodeToString(bytes)
	if len(encoded) < length {
		return "", fmt.Errorf("encoded string length is less than required")
	}
	return encoded[:length], nil
}
