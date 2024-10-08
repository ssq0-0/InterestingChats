package services

import (
	"InterestingChats/backend/user_services/internal/consts"
	"InterestingChats/backend/user_services/internal/logger"
	"InterestingChats/backend/user_services/internal/models"
	"InterestingChats/backend/user_services/internal/utils"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func (us *userService) GetMyProfile(r *http.Request, log logger.Logger) (*models.Response, int, string, error) {
	var userID int
	if _, err := utils.ConvertValue(r.Header.Get("X-User-ID"), &userID); err != nil {
		return nil, http.StatusUnauthorized, consts.ErrUserUnathorized, fmt.Errorf(consts.InternalTokenInvalid)
	}

	response, _, _, err := utils.ProxyRequest(consts.GET_Method, fmt.Sprintf(consts.CACHE_GetUserProfileInfo, userID), nil, http.StatusOK)
	if err != nil {
		response, statusCode, clientErr, err := utils.ProxyRequest(consts.GET_Method, fmt.Sprintf(consts.DB_GetUserProfileInfo, userID), nil, http.StatusOK)
		if err != nil {
			return nil, statusCode, clientErr, err
		}
		if err := us.producer.Writer(response.Data, consts.KAFKA_Session); err != nil {
			return nil, 400, "", err
		}
		return response, http.StatusOK, "", nil
	}

	return response, http.StatusOK, "", nil
}

func (us *userService) GetUserInfo(r *http.Request, log logger.Logger) (*models.Response, int, string, error) {
	var reqUserID int
	if _, err := utils.ConvertValue(r.URL.Query().Get("userID"), &reqUserID); err != nil {
		return nil, http.StatusUnauthorized, consts.ErrUserUnathorized, fmt.Errorf(consts.InternalTokenInvalid)
	}

	response, _, _, err := utils.ProxyRequest(consts.GET_Method, fmt.Sprintf(consts.CACHE_GetUserProfileInfo, reqUserID), nil, http.StatusOK)
	if err != nil {
		response, statusCode, clientErr, err := utils.ProxyRequest(consts.GET_Method, fmt.Sprintf(consts.DB_GetUserProfileInfo, reqUserID), nil, http.StatusOK)
		if err != nil {
			return nil, statusCode, clientErr, err
		}
		if err := us.producer.Writer(response.Data, consts.KAFKA_Session); err != nil {
			return nil, 400, "", err
		}
		return response, http.StatusOK, "", nil
	}

	return response, http.StatusOK, "", nil
}

func (us *userService) GetSearchUserResult(r *http.Request) (*models.Response, int, string, error) {
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
