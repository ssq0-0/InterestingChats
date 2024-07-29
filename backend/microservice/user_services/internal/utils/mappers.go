package utils

import (
	"InterestingChats/backend/user_services/internal/models"
	"fmt"
	"log"
)

func MapResponseDataToUser(data interface{}, user *models.User) error {
	responseData, ok := data.(map[string]interface{})
	if !ok {
		log.Printf("Invalid response data format")
		return fmt.Errorf("invalid response data format")
	}

	if id, ok := responseData["id"].(float64); ok {
		user.ID = int(id)
	} else {
		return fmt.Errorf("missing or invalid 'id' in response data")
	}

	if username, ok := responseData["username"].(string); ok {
		user.Username = username
	} else {
		return fmt.Errorf("missing or invalid 'username' in response data")
	}

	if email, ok := responseData["email"].(string); ok {
		user.Email = email
	} else {
		return fmt.Errorf("missing or invalid 'email' in response data")
	}

	return nil
}
