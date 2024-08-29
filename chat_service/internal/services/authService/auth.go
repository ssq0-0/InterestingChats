package services

import (
	"chat_service/internal/consts"
	"chat_service/internal/models"
	"chat_service/internal/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

func AuthenticateAndAuthorize(r *http.Request) ([]models.MemberRequest, int, string, error) {
	authToken := r.Header.Get("Authorization")
	userID, statusCode, clientErr, err := utils.CheckToken(authToken)
	if err != nil {
		return nil, statusCode, clientErr, err
	}

	var membersList []models.MemberRequest
	if err := json.NewDecoder(r.Body).Decode(&membersList); err != nil {
		return nil, http.StatusInternalServerError, consts.ErrInternalServer, fmt.Errorf(consts.InternalInvalidValueFormat, err)
	}
	if len(membersList) == 0 {
		return nil, http.StatusInternalServerError, consts.ErrUsersNotFoundInList, fmt.Errorf(consts.InternalNoFoundUsers)

	}

	if _, statusCode, clientErr, err := utils.ProxyRequest(consts.GET_Method, fmt.Sprintf(consts.DB_GetAuthor, userID, membersList[0].ChatID), nil, http.StatusOK); err != nil {
		return nil, statusCode, clientErr, err
	}

	return membersList, http.StatusOK, "", nil
}
