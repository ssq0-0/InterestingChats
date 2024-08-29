package services

import (
	"InterestingChats/backend/user_services/internal/consts"
	"InterestingChats/backend/user_services/internal/models"
	"InterestingChats/backend/user_services/internal/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func ChangeUserData(r *http.Request) (*models.Response, int, string, error) {
	token := strings.ReplaceAll(r.Header.Get("Authorization"), "Bearer ", "")
	if _, statusCode, err := utils.ValidateJWT(token); err != "" {
		return nil, statusCode, consts.ErrUserUnathorized, fmt.Errorf(consts.InternalTokenInvalid)
	}
	var data *models.ChangeUserData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return nil, http.StatusBadRequest, "", err
	}

	response, statusCode, clientErr, err := utils.ProxyRequest(consts.POST_Method, consts.DB_ChangeUserData, data, http.StatusOK)
	if err != nil {
		log.Printf("here1: %+v", err)
		return nil, statusCode, clientErr, err
	}

	return response, http.StatusOK, "", nil
}
