package services

import (
	"InterestingChats/backend/user_services/internal/consts"
	"InterestingChats/backend/user_services/internal/models"
	"InterestingChats/backend/user_services/internal/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

func (us *userService) ChangeUserData(r *http.Request) (*models.Response, int, string, error) {
	var userID int
	if _, err := utils.ConvertValue(r.Header.Get("X-User-ID"), &userID); err != nil {
		return nil, http.StatusUnauthorized, consts.ErrUserUnathorized, fmt.Errorf(consts.InternalTokenInvalid)
	}

	var data *models.ChangeUserData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return nil, http.StatusBadRequest, "", err
	}

	data.UserID = userID
	response, statusCode, clientErr, err := utils.ProxyRequest(consts.POST_Method, consts.DB_ChangeUserData, data, http.StatusOK)
	if err != nil {
		return nil, statusCode, clientErr, err
	}

	if err := us.producer.Writer(data, consts.KAFKA_session_UPDATE); err != nil {
		return nil, 400, clientErr, err
	}

	return response, http.StatusOK, "", nil
}

func (us *userService) AddUserAvatar(r *http.Request) (*models.Response, int, string, error) {
	var userID int
	if _, err := utils.ConvertValue(r.Header.Get("X-User-ID"), &userID); err != nil {
		return nil, http.StatusUnauthorized, consts.ErrUserUnathorized, fmt.Errorf(consts.InternalTokenInvalid)
	}

	response, statusCode, clientErr, err := utils.FileProxyRequest(r, consts.POST_Method, consts.FS_UploadFile, http.StatusOK)
	if err != nil {
		return nil, statusCode, clientErr, err
	}

	_, statusCode, clientErr, err = utils.ProxyRequest(consts.POST_Method, consts.DB_UploadPhoto, models.UserFile{UserID: userID, URL: *response}, http.StatusOK)
	if err != nil {
		return nil, statusCode, clientErr, err
	}

	if err := us.producer.Writer(&models.ChangeUserData{UserID: userID, Type: "avatar", Data: response.TemporaryLink}, consts.KAFKA_session_UPDATE); err != nil {
		return nil, 400, clientErr, err
	}

	return &models.Response{Errors: nil, Data: response.StaticLink}, http.StatusOK, "", nil
}
