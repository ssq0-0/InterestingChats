package utils

import (
	"InterestingChats/backend/user_services/internal/consts"
	"InterestingChats/backend/user_services/internal/models"
	"encoding/json"
	"fmt"
	"net/http"
)

func ValideUserData(r *http.Request, validationType int) (models.User, error) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		return models.User{}, err
	}

	switch validationType {
	case 0:
		if user.Email == "" || user.Username == "" || user.Password == "" {
			return models.User{}, fmt.Errorf(consts.InternalMissingParametr)
		}
	case 1:
		if user.Email == "" || user.Password == "" {
			return models.User{}, fmt.Errorf(consts.InternalMissingParametr)
		}
	default:
		return models.User{}, fmt.Errorf(consts.InternalMissingParametr)
	}

	return user, nil
}
