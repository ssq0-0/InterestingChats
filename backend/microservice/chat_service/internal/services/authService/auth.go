package services

import (
	"chat_service/internal/consts"
	"chat_service/internal/models"
	"chat_service/internal/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func AuthenticateAndAuthorize(r *http.Request) ([]models.MemberRequest, error) {
	authToken := r.Header.Get("Authorization")
	userID, err := utils.CheckToken(authToken)
	if err != nil {
		return nil, fmt.Errorf("error check token: %v", err)
	}

	var membersList []models.MemberRequest
	if err := json.NewDecoder(r.Body).Decode(&membersList); err != nil {
		return nil, fmt.Errorf("error decoding r body to json: %v", err)
	}
	if len(membersList) == 0 {
		return nil, fmt.Errorf("no members to delete")

	}

	chatID := membersList[0].ChatID
	log.Printf("request going to: %s", fmt.Sprintf(consts.DB_GetAuthor, userID, chatID))
	_, authorStatusCode, err := utils.ProxyRequest(consts.GET_Method, fmt.Sprintf(consts.DB_GetAuthor, userID, chatID), nil, http.StatusOK)
	if err != nil {
		return nil, fmt.Errorf("access foribben: %v, status code: %d", err, authorStatusCode)
	}

	return membersList, nil
}
