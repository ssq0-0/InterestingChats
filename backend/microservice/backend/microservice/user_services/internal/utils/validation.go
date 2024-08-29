package utils

import (
	"InterestingChats/backend/user_services/internal/models"
	"encoding/json"
	"fmt"
	"net/http"
)

func ValidRegistrationData(r *http.Request, validationType int) (models.User, error) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		return models.User{}, err
	}

	switch validationType {
	case 0:
		if user.Email == "" || user.Username == "" || user.Password == "" {
			return models.User{}, fmt.Errorf("missing usename, password or email")
		}
	case 1:
		if user.Email == "" || user.Password == "" {
			return models.User{}, fmt.Errorf("missing usename, password or email")
		}
	default:
		return models.User{}, fmt.Errorf("invalid vadation type")
	}

	return user, nil
}
