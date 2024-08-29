package services

import (
	"InterestingChats/backend/user_services/internal/consts"
	"InterestingChats/backend/user_services/internal/models"
	"InterestingChats/backend/user_services/internal/utils"
	"fmt"
	"net/http"
	"strings"
)

func UserRequest(method, url string, user models.User, expectedStatusCode int) (*models.UserTokens, *models.User, []string, error) {
	response, _, clientErr, err := utils.ProxyRequest(method, url, user, expectedStatusCode)
	if err != nil {
		return nil, nil, []string{clientErr}, err
	}

	if len(response.Errors) > 0 {
		resErrors := strings.Join(response.Errors, "; ")
		return nil, nil, response.Errors, fmt.Errorf(resErrors)
	}

	var userResponseInfo models.User
	if err := utils.MapResponseDataToUser(response.Data, &userResponseInfo); err != nil {
		return nil, nil, []string{consts.ErrInvalidValueFormat}, err
	}

	userTokens, err := GetORSetToken(userResponseInfo, consts.POST_Method, consts.Redis_SetToken, http.StatusOK)
	if err != nil {
		return nil, nil, []string{consts.ErrInternalServer}, err
	}

	tokens, ok := userTokens[userResponseInfo.Email]
	if !ok {
		return nil, nil, []string{consts.ErrUnexpectedRecivedEmail}, fmt.Errorf(consts.InternalUserEmailInToken)
	}

	return &tokens, &userResponseInfo, nil, nil
}
