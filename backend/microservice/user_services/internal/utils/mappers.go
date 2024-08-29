package utils

import (
	"InterestingChats/backend/user_services/internal/consts"
	"InterestingChats/backend/user_services/internal/models"
	"fmt"
)

func MapResponseDataToUser(data interface{}, user *models.User) error {
	responseData, ok := data.(map[string]interface{})
	if !ok {
		return fmt.Errorf(consts.InternalGhostType)
	}

	if id, ok := responseData["id"].(float64); ok {
		user.ID = int(id)
	} else {
		return fmt.Errorf(consts.InternalMissingParametr)
	}

	if username, ok := responseData["username"].(string); ok {
		user.Username = username
	} else {
		return fmt.Errorf(consts.InternalMissingParametr)
	}

	if email, ok := responseData["email"].(string); ok {
		user.Email = email
	} else {
		return fmt.Errorf(consts.InternalMissingParametr)
	}

	return nil
}
