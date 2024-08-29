package userservice

import (
	"chat_service/internal/consts"
	"chat_service/internal/models"
	"chat_service/internal/utils"
	"fmt"
	"net/http"
)

func TokenVerification(token string) (int, int, string, error) {
	userID, statusCode, clientErr, err := utils.CheckToken(token)
	if err != nil {
		return 0, statusCode, clientErr, err
	}

	if clientErr, statusCode, err := UserExists(userID); err != nil {
		return 0, statusCode, clientErr, err
	}

	return userID, http.StatusOK, "", nil
}

func UserExists(userID int) (string, int, error) {
	_, statusCode, clientErr, err := utils.ProxyRequest(consts.GET_Method, fmt.Sprintf(consts.DB_CheckUser, userID), nil, http.StatusOK)
	if err != nil {
		return clientErr, statusCode, err
	}
	return "", http.StatusOK, nil
}

func ManageMember(method, url string, member models.MemberRequest, statusCode int) (string, error) {
	_, _, clientErr, err := utils.ProxyRequest(method, url, member, statusCode)
	if err != nil {
		return clientErr, err
	}

	return "", nil
}
