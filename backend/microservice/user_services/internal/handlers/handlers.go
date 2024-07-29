package handlers

import (
	"InterestingChats/backend/user_services/internal/consts"
	"InterestingChats/backend/user_services/internal/models"
	"InterestingChats/backend/user_services/internal/services"
	"InterestingChats/backend/user_services/internal/utils"
	"fmt"

	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Registrations(w http.ResponseWriter, r *http.Request) {
	user, err := utils.ValidRegistrationData(r, 0)
	if err != nil {
		h.HandleError(w, http.StatusBadRequest, []string{err.Error()}, fmt.Sprintf("failed to decode JSON data: %v", err))
		return
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		h.HandleError(w, http.StatusBadRequest, []string{"Failed hash password"}, fmt.Sprintf("failed to to hash password: %v", err))
		return
	}
	user.Password = string(hashPassword)

	data, errors, err := services.HandleUserRequest(consts.POST_Method, consts.DB_Registration, user, http.StatusCreated)
	if err != nil {
		h.HandleError(w, http.StatusBadRequest, errors, fmt.Sprintf("error method: %v", err))
		return
	}

	h.SendRespond(w, http.StatusCreated, &models.Response{
		Data:   data,
		Errors: nil,
	})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	user, err := utils.ValidRegistrationData(r, 1)
	if err != nil {
		h.HandleError(w, http.StatusBadRequest, []string{fmt.Sprintf("Failed to recived user data: %v", err)}, fmt.Sprintf("failed decode user data: %v", err))
		return
	}

	tokens, errors, err := services.HandleUserRequest(consts.POST_Method, consts.DB_Login, user, http.StatusOK)
	if err != nil {
		h.HandleError(w, http.StatusBadRequest, errors, fmt.Sprintf("failed recived tokens:: %v", err))
		return
	}

	h.SendRespond(w, http.StatusOK, &models.Response{
		Errors: nil,
		Data:   tokens,
	})
}

func (h *Handler) GetTokens(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		h.HandleError(w, http.StatusBadRequest, []string{"missing email in URL"}, "missing email in URL")
		return
	}

	redisResponse, statusCode, err := utils.ProxyRequest(consts.GET_Method, fmt.Sprintf(consts.Redis_GetToken, email), nil, http.StatusOK)
	if err != nil {
		h.HandleError(w, statusCode, []string{"Failed recived tokens"}, fmt.Sprintf("failed recived tokens from redis: %v", err))
		return
	}

	var tokens models.Tokens
	err = json.Unmarshal(redisResponse, &tokens)
	if err != nil {
		h.HandleError(w, statusCode, []string{"Failed decode tokens"}, fmt.Sprintf("failed decode tokens from redis response: %v", err))
		return
	}

	h.SendRespond(w, http.StatusOK, &models.Response{
		Errors: nil,
		Data:   tokens,
	})
}

func (h *Handler) CheckTokens(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		h.HandleError(w, http.StatusBadRequest, []string{"missin token in url"}, "failed get token from url")
		return
	}

	id, err := utils.ValidateJWT(token)
	if err != nil {
		h.HandleError(w, http.StatusBadRequest, []string{"incorrect token"}, fmt.Sprintf("incorrect token: %v", err))
		return
	}

	h.SendRespond(w, http.StatusOK, &models.Response{
		Errors: nil,
		Data:   id,
	})
}
