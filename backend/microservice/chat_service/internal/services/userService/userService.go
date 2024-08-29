package userservice

import (
	"chat_service/internal/consts"
	"chat_service/internal/models"
	"chat_service/internal/utils"
	"fmt"
	"log"
	"net/http"
)

func UserVerification(r *http.Request) (int, error) {
	authToken := r.Header.Get("Authorization")
	userID, err := utils.CheckToken(authToken)
	if err != nil {
		log.Printf("access rejected: %v", err)
		return 0, fmt.Errorf("access rejected: %v", err)
	}
	return userID, nil
}

func UserExists(userID int) error {
	if _, _, err := utils.ProxyRequest(consts.GET_Method, fmt.Sprintf(consts.DB_CheckUser, userID), nil, http.StatusOK); err != nil {
		log.Printf("failed request: %v", err)
		return err
	}
	return nil
}

func ManageMember(method, url string, member models.MemberRequest, statusCode int) error {
	_, _, err := utils.ProxyRequest(method, url, member, statusCode)
	if err != nil {
		log.Printf("failed request: %v", err)
		return err
	}
	return nil
}
