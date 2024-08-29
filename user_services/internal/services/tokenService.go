package services

import (
	"InterestingChats/backend/user_services/internal/consts"
	"InterestingChats/backend/user_services/internal/models"
	"InterestingChats/backend/user_services/internal/utils"
	"fmt"
)

func GetORSetToken(user models.User, method, url string, expectedStatusCode int) (map[string]models.UserTokens, error) {
	accessToken, refreshToken, err := utils.GenerateJWT(user)
	if err != nil {
		return nil, err
	}

	userTokens := map[string]models.UserTokens{
		user.Email: {
			Tokens: models.Tokens{
				AccessToken:  accessToken,
				RefreshToken: refreshToken,
			},
		},
	}

	if _, _, _, err = utils.ProxyRequest(method, url, userTokens, expectedStatusCode); err != nil {
		return nil, fmt.Errorf(consts.InternalFailedSetToken, err)
	}

	return userTokens, nil
}
