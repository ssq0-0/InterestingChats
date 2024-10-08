package utils

import (
	"InterestingChats/backend/user_services/internal/consts"
	"InterestingChats/backend/user_services/internal/models"
	"encoding/json"
	"fmt"
	"net/http"
)

// ValideUserData checks the user data in the request against the validation type.
// Returns the user model or an error if something goes wrong
func ValideUserData(r *http.Request, validationType int) (models.User, error) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		return models.User{}, err
	}

	switch validationType {
	case consts.VALDIDATION_RegistrationType:
		if user.Email == "" || user.Username == "" || user.Password == "" {
			return models.User{}, fmt.Errorf(consts.InternalMissingParametr)
		}
	case consts.VALDIDATION_LoginType:
		if user.Email == "" || user.Password == "" {
			return models.User{}, fmt.Errorf(consts.InternalMissingParametr)
		}
	default:
		return models.User{}, fmt.Errorf(consts.InternalMissingParametr)
	}

	return user, nil
}
