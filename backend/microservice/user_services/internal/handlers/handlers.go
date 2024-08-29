package handlers

import (
	"InterestingChats/backend/user_services/internal/consts"
	"InterestingChats/backend/user_services/internal/logger"
	"InterestingChats/backend/user_services/internal/models"
	"InterestingChats/backend/user_services/internal/services"
	"InterestingChats/backend/user_services/internal/utils"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"net/http"
)

type UserService struct {
	log logger.Logger
}

func NewService(log logger.Logger) *UserService {
	return &UserService{
		log: log,
	}
}

func (h *UserService) Registrations(w http.ResponseWriter, r *http.Request) {
	user, err := utils.ValideUserData(r, consts.VALDIDATION_RegistrationType)
	if err != nil {
		h.HandleError(w, http.StatusBadRequest, []string{err.Error()}, err)
		return
	}

	user.Password, err = utils.HashPassword(user.Password)
	if err != nil || user.Password == "" {
		h.HandleError(w, http.StatusBadRequest, []string{consts.ErrInternalServer}, err)
		return
	}

	tokens, userInfo, clientErr, err := services.UserRequest(consts.POST_Method, consts.DB_Registration, user, http.StatusCreated)
	if err != nil {
		h.HandleError(w, http.StatusBadRequest, clientErr, err)
		return
	}

	h.log.Infof("successful registred user: %+v", user.Email)
	h.SendRespond(w, http.StatusCreated, &models.Response{
		Data: models.AuthResponse{
			Tokens: tokens.Tokens,
			User:   *userInfo,
		},
		Errors: nil,
	})
}

func (h *UserService) Login(w http.ResponseWriter, r *http.Request) {
	user, err := utils.ValideUserData(r, consts.VALDIDATION_LoginType)
	if err != nil {
		h.HandleError(w, http.StatusBadRequest, []string{consts.ErrUnexpectedValueFormat}, err)
		return
	}

	tokens, userInfo, clientErr, err := services.UserRequest(consts.POST_Method, consts.DB_Login, user, http.StatusOK)
	if err != nil {
		h.HandleError(w, http.StatusBadRequest, clientErr, err)
		return
	}

	h.log.Infof("successful login user: %+v", user.Email)
	h.SendRespond(w, http.StatusOK, &models.Response{
		Errors: nil,
		Data: models.AuthResponse{
			Tokens: tokens.Tokens,
			User:   *userInfo,
		},
	})
}

func (h *UserService) GetMyProfile(w http.ResponseWriter, r *http.Request) {
	response, statusCode, clientErr, err := services.GetMyProfile(r, h.log)
	if err != nil {
		h.HandleError(w, statusCode, []string{clientErr}, err)
		return
	}

	h.SendRespond(w, http.StatusOK, response)
}

func (h *UserService) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	response, statusCode, clientErr, err := services.GetUserInfo(r, h.log)
	if err != nil {
		h.HandleError(w, statusCode, []string{clientErr}, err)
		return
	}

	h.SendRespond(w, http.StatusOK, response)
}

func (h *UserService) GetSearchUserResult(w http.ResponseWriter, r *http.Request) {
	response, statusCode, clientErr, err := services.GetSearchUserResult(r)
	if err != nil {
		h.HandleError(w, statusCode, []string{clientErr}, err)
		return
	}

	h.SendRespond(w, http.StatusOK, response)
}

// TODO: add user request verification √
func (h *UserService) ChangeUserData(w http.ResponseWriter, r *http.Request) {
	// TODO: change '_' to response. Use response.Errors in errorhandler √

	response, statusCode, clientErr, err := services.ChangeUserData(r)
	if err != nil {
		log.Printf("here: %v", err)

		h.HandleError(w, statusCode, []string{clientErr}, err)
		return
	}

	h.SendRespond(w, http.StatusOK, &models.Response{
		Errors: nil,
		Data:   response,
	})
}

func (h *UserService) GetTokens(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		h.HandleError(w, http.StatusBadRequest, []string{consts.ErrMissingParametr}, fmt.Errorf(consts.InternalMissingParametr))
		return
	}

	redisResponse, statusCode, clientErr, err := utils.ProxyRequest(consts.GET_Method, fmt.Sprintf(consts.Redis_GetToken, email), nil, http.StatusOK)
	if err != nil {
		h.HandleError(w, statusCode, []string{clientErr}, err)
		return
	}

	redisResponseBytes, ok := redisResponse.Data.([]byte)
	if !ok {
		h.HandleError(w, http.StatusBadGateway, []string{consts.ErrInvalidValueFormat}, fmt.Errorf("unexpected data format"))
		return
	}

	var tokens models.Tokens
	if err := json.Unmarshal(redisResponseBytes, &tokens); err != nil {
		h.HandleError(w, http.StatusBadGateway, []string{consts.ErrInvalidValueFormat}, err)
		return
	}

	h.SendRespond(w, http.StatusOK, &models.Response{
		Errors: nil,
		Data:   tokens,
	})
}

func (h *UserService) CheckTokens(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		h.HandleError(w, http.StatusBadRequest, []string{consts.ErrMissingParametr}, fmt.Errorf(consts.InternalMissingParametr))
		return
	}

	id, statusCode, err := utils.ValidateJWT(token)
	if err != "" {
		h.HandleError(w, statusCode, []string{err}, fmt.Errorf(err))
		return
	}

	h.SendRespond(w, http.StatusOK, &models.Response{
		Errors: nil,
		Data:   id,
	})
}

func (h *UserService) RefreshAccessToken(w http.ResponseWriter, r *http.Request) {
	refreshToken := r.URL.Query().Get("refreshToken")
	if strings.ReplaceAll(refreshToken, " ", "") == "" {
		h.HandleError(w, http.StatusUnauthorized, []string{consts.ErrUserUnathorized}, fmt.Errorf(consts.InternalTokenInvalid))
		return
	}

	newAccessToken, clientErr, err := utils.RefreshToken(refreshToken)
	if err != nil {
		h.HandleError(w, http.StatusUnauthorized, []string{clientErr}, err)
		return
	}

	h.SendRespond(w, http.StatusOK, &models.Response{
		Errors: nil,
		Data:   newAccessToken,
	})
}
