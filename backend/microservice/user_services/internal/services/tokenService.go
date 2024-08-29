package services

import (
	"InterestingChats/backend/user_services/internal/models"
	"InterestingChats/backend/user_services/internal/utils"
	"fmt"
)

func GetORSetToken(user models.User, method, url string, expectedStatusCode int) (map[string]models.UserTokens, []string, error) {
	var errors []string
	accessToken, refreshToken, err := utils.GenerateJWT(user)
	if err != nil {
		errors = append(errors, err.Error())
		return nil, errors, fmt.Errorf("problems with generate JWT: %v", err)
	}

	userTokens := map[string]models.UserTokens{
		user.Email: {
			Tokens: models.Tokens{
				AccessToken:  accessToken,
				RefreshToken: refreshToken,
			},
		},
	}

	if _, _, err = utils.ProxyRequest(method, url, userTokens, expectedStatusCode); err != nil {
		errors = append(errors, err.Error())
		return nil, errors, fmt.Errorf("failed to store tokens in Redis: %v", err)
	}

	return userTokens, errors, nil
}
