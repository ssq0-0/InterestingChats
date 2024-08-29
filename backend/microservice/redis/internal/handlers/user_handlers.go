package handlers

import (
	"InterestingChats/backend/microservice/redis/internal/consts"
	"InterestingChats/backend/microservice/redis/internal/logger"
	"InterestingChats/backend/microservice/redis/internal/models"
	"InterestingChats/backend/microservice/redis/internal/rdb"
	"encoding/json"
	"io"
	"net/http"
)

type UserHandler struct {
	rdb *rdb.RedisClient
	log logger.Logger
}

func NewUserHandler(rdb *rdb.RedisClient, log logger.Logger) *UserHandler {
	return &UserHandler{
		rdb: rdb,
		log: log,
	}
}

func (uh *UserHandler) GetUsersTokens(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		HandleError(w, http.StatusBadRequest, []string{consts.ErrMissingParametr}, consts.InternalMissingParametr)
		return
	}

	tokens, err := uh.rdb.GetTokens(email)
	if err != nil {
		HandleError(w, http.StatusBadRequest, []string{consts.ErrInternalServer}, err.Error())
		return
	}

	uh.log.Infof("successful sending of %s tokens", email)
	SendRespond(w, http.StatusOK, &models.Response{
		Errors: nil,
		Data:   tokens,
	})
}

func (uh *UserHandler) SetTokens(w http.ResponseWriter, r *http.Request) {
	var userTokens map[string]models.UserTokens
	body, err := io.ReadAll(r.Body)
	if err != nil {
		HandleError(w, http.StatusBadGateway, []string{consts.ErrUnexpectedValueFormat}, consts.InternalInvalidValueFormat)
		return
	}

	if err := json.Unmarshal(body, &userTokens); err != nil {
		HandleError(w, http.StatusBadGateway, []string{consts.ErrUnexpectedValueFormat}, consts.InternalInvalidValueFormat)
		return
	}

	if err := uh.rdb.SetToken(userTokens); err != nil {
		HandleError(w, http.StatusInternalServerError, []string{consts.ErrInternalServer}, err.Error())
		return
	}

	uh.log.Infof("successful set of tokens")
	SendRespond(w, http.StatusOK, &models.Response{
		Data:   "success",
		Errors: nil,
	})
}

func (uh *UserHandler) DeleteTokens(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		HandleError(w, http.StatusBadRequest, []string{consts.ErrInternalServer}, consts.InternalMissingParametr)
		return
	}

	if err := uh.rdb.DeleteTokens(email); err != nil {
		HandleError(w, http.StatusBadRequest, []string{consts.ErrInternalServer}, err.Error())
		return
	}

	SendRespond(w, http.StatusOK, &models.Response{
		Errors: nil,
		Data:   "successfull deleted tokens",
	})
}
