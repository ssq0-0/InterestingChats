package services

import (
	"InterestingChats/backend/user_services/internal/consts"
	"InterestingChats/backend/user_services/internal/models"
	"InterestingChats/backend/user_services/internal/utils"
	"fmt"
	"net/http"
)

func HandleUserRequest(method, url string, user models.User, expectedStatusCode int) (models.UserTokens, []string, error) {
	body, _, err := utils.ProxyRequest(method, url, user, expectedStatusCode)
	if err != nil {
		return models.UserTokens{}, []string{"Failed to complete the request"}, fmt.Errorf("failed to complete the request: %v", err)
	}

	var response models.Response
	if err := utils.ParseBody(body, &response); err != nil {
		return models.UserTokens{}, []string{"Failed to parse response"}, fmt.Errorf("failed to parse response: %v", err)
	}

	if len(response.Errors) > 0 {
		return models.UserTokens{}, response.Errors, fmt.Errorf("registration errors: %v", response.Errors)
	}

	var userResponseInfo models.User
	if err := utils.MapResponseDataToUser(response.Data, &userResponseInfo); err != nil {
		return models.UserTokens{}, []string{"Failed to parse user info from response"}, err
	}

	userTokens, errors, err := GetORSetToken(userResponseInfo, consts.POST_Method, consts.Redis_SetToken, http.StatusOK)
	if err != nil {
		return models.UserTokens{}, errors, fmt.Errorf("failed to store tokens in Redis: %v", err)
	}

	tokens, ok := userTokens[userResponseInfo.Email]
	if !ok {
		return models.UserTokens{}, []string{"Can't find email in user tokens"}, fmt.Errorf("can't find email in user tokens")
	}

	return tokens, nil, nil
}
