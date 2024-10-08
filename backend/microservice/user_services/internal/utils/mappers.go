package utils

import (
	"InterestingChats/backend/user_services/internal/consts"
	"InterestingChats/backend/user_services/internal/models"
	"fmt"
)

// MapResponseDataToUser maps the response data to the user model.
// If the data does not match the expected format, an error is returned.
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

// IsEmptyResponseData checks if the data is empty.
// Returns true if the data is empty, otherwise false.
func IsEmptyResponseData(data interface{}) bool {
	if slice, ok := data.([]interface{}); ok {
		return len(slice) == 0
	}
	return true
}
