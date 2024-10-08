package handlers

import (
	"InterestingChats/backend/microservice/redis/internal/consts"
	"InterestingChats/backend/microservice/redis/internal/models"
	"InterestingChats/backend/microservice/redis/internal/utils"
	"fmt"
	"net/http"
	"strings"
)

// GetSession retrieves the session for a user based on the provided userID from the request URL.
func (h *handler) GetSession(w http.ResponseWriter, r *http.Request) {
	user_id := r.URL.Query().Get("userID")
	if strings.ReplaceAll(user_id, " ", "") == "" {
		HandleError(w, http.StatusBadGateway, []string{consts.ErrMissingParametr}, fmt.Errorf(consts.InternalMissingParametr))
		return
	}

	session, err := h.rdb.GetSession(user_id)
	if err != nil {
		HandleError(w, http.StatusBadGateway, []string{consts.ErrInternalServer}, err)
		return
	}

	SendRespond(w, http.StatusOK, &models.Response{
		Errors: nil,
		Data:   session,
	})
}

// GetFriends retrieves the list of friends for a user based on the provided userID from the request URL.
func (h *handler) GetFriends(w http.ResponseWriter, r *http.Request) {
	var userID int
	if clientErr, err := utils.ConvertValue(r.URL.Query().Get("userID"), &userID); err != nil {
		HandleError(w, http.StatusInternalServerError, []string{clientErr}, err)
		return
	}

	friendList, err := h.rdb.GetFriends(userID)
	if err != nil {
		HandleError(w, http.StatusInternalServerError, []string{"failed get frineds list"}, err)
		return
	}

	SendRespond(w, http.StatusOK, &models.Response{
		Errors: nil,
		Data:   friendList,
	})
}

// GetSubscribers retrieves the list of subscribers for a user based on the provided userID from the request URL.
func (h *handler) GetSubscribers(w http.ResponseWriter, r *http.Request) {
	var userID int
	if clientErr, err := utils.ConvertValue(r.URL.Query().Get("userID"), &userID); err != nil {
		HandleError(w, http.StatusInternalServerError, []string{clientErr}, err)
		return
	}

	subsList, err := h.rdb.GetSubscribers(userID)
	if err != nil {
		HandleError(w, http.StatusInternalServerError, []string{"failed get frineds list"}, err)
		return
	}

	SendRespond(w, http.StatusOK, &models.Response{
		Errors: nil,
		Data:   subsList,
	})
}
