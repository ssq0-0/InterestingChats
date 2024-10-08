package userservice

import (
	"chat_service/internal/consts"
	"chat_service/internal/models"
	"chat_service/internal/utils"
	"fmt"
	"net/http"
)

// UserExists checks if a user exists in the database by their user ID.
// It returns any client error and the HTTP status code.
func UserExists(userID int) (string, int, error) {
	_, statusCode, clientErr, err := utils.ProxyRequest(consts.GET_Method, fmt.Sprintf(consts.DB_CheckUser, userID), nil, http.StatusOK)
	if err != nil {
		return clientErr, statusCode, err
	}
	return "", http.StatusOK, nil
}

// ManageMember handles the addition or removal of a member in a chat.
// It sends a request to the specified URL using the provided HTTP method and member data.
func ManageMember(method, url string, member models.MemberRequest, statusCode int) (string, error) {
	_, _, clientErr, err := utils.ProxyRequest(method, url, member, statusCode)
	if err != nil {
		return clientErr, err
	}

	return "", nil
}
