package services

import (
	"InterestingChats/backend/user_services/internal/consts"
	"InterestingChats/backend/user_services/internal/models"
	"InterestingChats/backend/user_services/internal/utils"
	"fmt"
	"net/http"
	"strings"
)

func (us *userService) UserRequest(method, url string, user models.User, expectedStatusCode int) (*models.UserTokens, *models.User, string, error) {
	response, _, clientErr, err := utils.ProxyRequest(method, url, user, expectedStatusCode)
	if err != nil {
		return nil, nil, clientErr, err
	}

	if len(response.Errors) > 0 {
		resErrors := strings.Join(response.Errors, "; ")
		return nil, nil, resErrors, fmt.Errorf(resErrors)
	}

	var userResponseInfo models.User
	if err := utils.MapResponseDataToUser(response.Data, &userResponseInfo); err != nil {
		return nil, nil, consts.ErrInvalidValueFormat, err
	}

	response, _, _, err = utils.ProxyRequest(consts.POST_Method, consts.AS_GenerateTokens, userResponseInfo, http.StatusOK)
	if err != nil {
		return nil, nil, consts.ErrInternalServer, err
	}

	tokensData, ok := response.Data.(map[string]interface{})
	if !ok {
		return nil, nil, "Invalid token response format", fmt.Errorf("expected map[string]interface{}")
	}

	tokensMap, ok := tokensData["tokens"].(map[string]interface{})
	if !ok {
		return nil, nil, "Invalid tokens format", fmt.Errorf("expected tokens as map[string]interface{}")
	}

	tokens := models.Tokens{
		AccessToken:  tokensMap["accessToken"].(string),
		RefreshToken: tokensMap["refreshToken"].(string),
	}

	userTokens := &models.UserTokens{Tokens: tokens}
	return userTokens, &userResponseInfo, "", nil
}
