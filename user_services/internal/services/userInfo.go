package services

import (
	"InterestingChats/backend/user_services/internal/consts"
	"InterestingChats/backend/user_services/internal/logger"
	"InterestingChats/backend/user_services/internal/models"
	"InterestingChats/backend/user_services/internal/utils"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func GetMyProfile(r *http.Request, log logger.Logger) (*models.Response, int, string, error) {
	userID, statusCode, clientErr := utils.ValidateJWT(strings.ReplaceAll(r.Header.Get("Authorization"), "Bearer ", ""))
	if clientErr != "" {
		return nil, statusCode, consts.ErrInvalidToken, fmt.Errorf(consts.InternalTokenInvalid)
	}
	if userID == 0 {
		return nil, http.StatusUnauthorized, consts.ErrUserUnathorized, fmt.Errorf(consts.InternalTokenInvalid)
	}

	response, statusCode, clientErr, err := utils.ProxyRequest(consts.GET_Method, fmt.Sprintf(consts.DB_GetUserProfileInfo, userID), nil, http.StatusOK)
	if err != nil {
		return nil, statusCode, clientErr, err
	}

	return response, http.StatusOK, "", nil
}

func GetUserInfo(r *http.Request, log logger.Logger) (*models.Response, int, string, error) {
	requestUserID, statusCode, clientErr := utils.ValidateJWT(strings.ReplaceAll(r.Header.Get("Authorization"), "Bearer ", ""))
	if clientErr != "" {
		return nil, statusCode, consts.ErrInvalidToken, fmt.Errorf(consts.InternalTokenInvalid)
	}
	if requestUserID == 0 {
		return nil, http.StatusUnauthorized, consts.ErrUserUnathorized, fmt.Errorf(consts.InternalTokenInvalid)
	}

	userID, err := strconv.Atoi(strings.ReplaceAll(r.URL.Query().Get("userID"), " ", ""))
	if err != nil {
		return nil, http.StatusBadRequest, consts.ErrInternalServer, fmt.Errorf(consts.InternalTokenInvalid)
	}

	log.Infof("req id %d and user id %d", requestUserID, userID)
	response, statusCode, clientErr, err := utils.ProxyRequest(consts.GET_Method, fmt.Sprintf(consts.DB_GetUserProfileInfo, userID), nil, http.StatusOK)
	if err != nil {
		log.Errorf("err: %v", err)
		return nil, statusCode, clientErr, err
	}

	log.Infof("resp %+v", response)
	return response, http.StatusOK, "", nil
}

func GetSearchUserResult(r *http.Request) (*models.Response, int, string, error) {
	requestUserID, statusCode, clientErr := utils.ValidateJWT(strings.ReplaceAll(r.Header.Get("Authorization"), "Bearer ", ""))
	if clientErr != "" {
		return nil, statusCode, consts.ErrInvalidToken, fmt.Errorf(consts.InternalTokenInvalid)
	}
	if requestUserID == 0 {
		return nil, http.StatusUnauthorized, consts.ErrUserUnathorized, fmt.Errorf(consts.InternalTokenInvalid)
	}

	symbols := r.URL.Query().Get("symbols")
	if strings.ReplaceAll(symbols, " ", "") == "" {
		return nil, http.StatusBadRequest, consts.ErrMissingParametr, fmt.Errorf(consts.InternalMissingParametr)
	}

	repsonse, statusCode, clieclientErr, err := utils.ProxyRequest(consts.GET_Method, fmt.Sprintf(consts.DB_SearchUsersBySymbols, symbols), nil, http.StatusOK)
	if err != nil {
		log.Printf("err here: %v", err)
		return nil, statusCode, clieclientErr, err
	}

	return repsonse, http.StatusOK, "", nil

}
