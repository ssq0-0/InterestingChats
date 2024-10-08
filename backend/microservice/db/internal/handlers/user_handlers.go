package handlers

import (
	"InterestingChats/backend/microservice/db/internal/consts"

	"InterestingChats/backend/microservice/db/internal/models"
	"InterestingChats/backend/microservice/db/internal/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Registrations to handle user registration requests.
func (h *handler) Registrations(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		ErrorHandler(w, http.StatusBadRequest, h.log, []string{consts.ErrInvalidValueFormat}, err)
		return
	}

	if userInfo, clientErr, err := h.DBService.CreateNewUser(user); err != nil {
		ErrorHandler(w, http.StatusBadRequest, h.log, []string{clientErr}, err)
		return

	} else {
		SendRespond(w, http.StatusCreated, &models.Response{
			Errors: nil,
			Data:   userInfo,
		})
	}
}

// Login to handle user login requests.
func (h *handler) Login(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		ErrorHandler(w, http.StatusBadRequest, h.log, []string{consts.ErrInvalidValueFormat}, err)
		return
	}

	userPassword, clientErr, err := h.DBService.LoginData(&user)
	if err != nil {
		ErrorHandler(w, http.StatusBadRequest, h.log, []string{clientErr}, err)
		return
	}

	if err := utils.CompareHashAndPassword(userPassword, user.Password); err != nil {
		ErrorHandler(w, http.StatusBadRequest, h.log, []string{consts.ErrIncorrectEmailOrPassword}, err)
		return
	}
	user.Password = ""

	SendRespond(w, http.StatusOK, &models.Response{
		Errors: nil,
		Data:   user,
	})
}

// CheckUser to verify if a user exists.
func (h *handler) CheckUser(w http.ResponseWriter, r *http.Request) {
	userIDstr := r.URL.Query().Get("userID")
	if userIDstr == "" {
		ErrorHandler(w, http.StatusBadRequest, h.log, []string{consts.ErrMissingUserID}, fmt.Errorf(consts.InternalErrMissingURLVal))
		return
	}

	var userID int
	if clientErr, err := utils.ConvertValue(userIDstr, &userID); err != nil {
		ErrorHandler(w, http.StatusBadRequest, h.log, []string{clientErr}, err)
		return
	}

	accept, clientErr, err := h.DBService.CheckUser(userID)
	if err != nil {
		ErrorHandler(w, http.StatusBadRequest, h.log, []string{clientErr}, err)
		return
	}
	if !accept {
		ErrorHandler(w, http.StatusBadRequest, h.log, []string{consts.ErrUserNotFound}, fmt.Errorf(consts.InternalErrUserNotFoud))
		return
	}

	SendRespond(w, http.StatusOK, &models.Response{
		Errors: nil,
		Data:   userID,
	})
}

// GetUserProfileInfo to retrieve detailed user profile information.
func (h *handler) GetUserProfileInfo(w http.ResponseWriter, r *http.Request) {
	var userID int
	if clientErr, err := utils.ConvertValue(r.URL.Query().Get("userID"), &userID); err != nil {
		ErrorHandler(w, http.StatusBadRequest, h.log, []string{clientErr}, err)
		return
	}

	_, clientErr, err := h.DBService.CheckUser(userID)
	if err != nil {
		ErrorHandler(w, http.StatusBadRequest, h.log, []string{clientErr}, err)
		return
	}

	userInfo, clientErr, err := h.DBService.GetUserInfo(userID)
	if err != nil {
		ErrorHandler(w, http.StatusBadRequest, h.log, []string{clientErr}, err)
		return
	}

	SendRespond(w, http.StatusOK, &models.Response{
		Errors: nil,
		Data:   userInfo,
	})
}

// SearchUsers to search for users based on input parameters.
func (h *handler) SearchUsers(w http.ResponseWriter, r *http.Request) {
	userList, clientErr, err := h.DBService.SearchUsers(r.URL.Query().Get("symbols"))
	if err != nil {
		ErrorHandler(w, http.StatusBadRequest, h.log, []string{clientErr}, err)
		return
	}

	SendRespond(w, http.StatusOK, &models.Response{
		Errors: nil,
		Data:   userList,
	})
}

// ChangeUserData to modify user data in the system.
func (h *handler) ChangeUserData(w http.ResponseWriter, r *http.Request) {
	var data *models.ChangeUserData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		ErrorHandler(w, http.StatusBadRequest, h.log, []string{consts.ErrInternalServerError}, err)
		return
	}
	log.Printf("change data: %+v", data)
	switch data.Type {
	case "username":
		if clientErr, err := h.DBService.ChangeUsername(data.Data, data.UserID); err != nil {
			ErrorHandler(w, http.StatusBadRequest, h.log, []string{clientErr}, err)
			return
		}
	case "email":
		if clientErr, err := h.DBService.ChangeEmail(data.Data, data.UserID); err != nil {
			ErrorHandler(w, http.StatusBadRequest, h.log, []string{clientErr}, err)
			return
		}
	default:
		ErrorHandler(w, http.StatusBadRequest, h.log, []string{consts.ErrInvalidRequestFormat}, fmt.Errorf(consts.InternalErrChangedUserInfo))
		return
	}

	SendRespond(w, http.StatusOK, &models.Response{
		Errors: nil,
		Data:   "successful changed",
	})
}

// UploadPhoto to handle photo uploads by users.
func (h *handler) UploadPhoto(w http.ResponseWriter, r *http.Request) {
	var userFiles *models.UserFile
	if err := json.NewDecoder(r.Body).Decode(&userFiles); err != nil {
		ErrorHandler(w, http.StatusBadRequest, h.log, []string{consts.ErrInvalidValueFormat}, err)
		return
	}
	log.Printf("request files: %v", userFiles)
	if statusCode, clientErr, err := h.DBService.UploadAvatar(userFiles); err != nil {
		ErrorHandler(w, statusCode, h.log, []string{clientErr}, err)
		return
	}

	SendRespond(w, http.StatusOK, nil)
}

// RequestToFriendShip to send a friendship request to another user.
func (h *handler) RequestToFriendShip(w http.ResponseWriter, r *http.Request) {
	user, statusCode, clientErr, err := h.DBService.FriendshipOperation(r, consts.FRIENDSHIP_RequestType)
	if err != nil {
		ErrorHandler(w, statusCode, h.log, []string{clientErr}, err)
		return
	}

	SendRespond(w, http.StatusOK, &models.Response{
		Errors: nil,
		Data:   user,
	})
}

// AcceptFriendShip to accept a friendship request from another user.
func (h *handler) AcceptFriendShip(w http.ResponseWriter, r *http.Request) {
	user, statusCode, clientErr, err := h.DBService.FriendshipOperation(r, consts.FRIENDSHIP_AcceptType)
	if err != nil {
		ErrorHandler(w, statusCode, h.log, []string{clientErr}, err)
		return
	}

	SendRespond(w, http.StatusOK, &models.Response{
		Errors: nil,
		Data:   user,
	})
}

// DeleteFriend to remove a user from the friend list.
func (h *handler) DeleteFriend(w http.ResponseWriter, r *http.Request) {
	user, statusCode, clientErr, err := h.DBService.FriendshipOperation(r, consts.FRIENDSHIP_DeleteFriend)
	if err != nil {
		ErrorHandler(w, statusCode, h.log, []string{clientErr}, err)
		return
	}

	SendRespond(w, http.StatusOK, &models.Response{
		Errors: nil,
		Data:   user,
	})
}

// DeleteFriendRequest to cancel a friendship request.
func (h *handler) DeleteFriendRequest(w http.ResponseWriter, r *http.Request) {
	user, statusCode, clientErr, err := h.DBService.FriendshipOperation(r, consts.FRIENDSHIP_DeleteFriendRequest)
	if err != nil {
		ErrorHandler(w, statusCode, h.log, []string{clientErr}, err)
		return
	}

	SendRespond(w, http.StatusOK, &models.Response{
		Errors: nil,
		Data:   user,
	})
}

// GetFriendList to retrieve the list of a user's friends.
func (h *handler) GetFriendList(w http.ResponseWriter, r *http.Request) {
	friendList, statusCode, clientErr, err := h.DBService.FriendshipOperation(r, consts.FRIENDSHIP_GetUsersFriends)
	if err != nil {
		ErrorHandler(w, statusCode, h.log, []string{clientErr}, err)
		return
	}

	SendRespond(w, http.StatusOK, &models.Response{
		Errors: nil,
		Data:   friendList,
	})
}

// GetSubList to retrieve the list of a user's subscribers.
func (h *handler) GetSubList(w http.ResponseWriter, r *http.Request) {
	friendList, statusCode, clientErr, err := h.DBService.FriendshipOperation(r, consts.FRIENDSHIP_GetSubs)
	if err != nil {
		ErrorHandler(w, statusCode, h.log, []string{clientErr}, err)
		return
	}

	SendRespond(w, http.StatusOK, &models.Response{
		Errors: nil,
		Data:   friendList,
	})
}
