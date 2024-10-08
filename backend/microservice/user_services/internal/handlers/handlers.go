package handlers

import (
	"InterestingChats/backend/user_services/internal/consts"
	"InterestingChats/backend/user_services/internal/logger"
	"InterestingChats/backend/user_services/internal/models"
	"InterestingChats/backend/user_services/internal/services"
	"InterestingChats/backend/user_services/internal/utils"
	"log"

	"net/http"
)

// Handler handles all user-related operations.
type Handler struct {
	US  services.UserService
	log logger.Logger
}

// NewService creates a new instance of Handler with a logger.
func NewService(log logger.Logger, us services.UserService) *Handler {
	return &Handler{
		US:  us,
		log: log,
	}
}

// Registrations handles user registration.
func (h *Handler) Registrations(w http.ResponseWriter, r *http.Request) {
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

	tokens, userInfo, clientErr, err := h.US.UserRequest(consts.POST_Method, consts.DB_Registration, user, http.StatusCreated)
	if err != nil {
		h.HandleError(w, http.StatusBadRequest, []string{clientErr}, err)
		return
	}

	h.SendRespond(w, http.StatusCreated, &models.Response{
		Data: models.AuthResponse{
			Tokens: tokens.Tokens,
			User:   *userInfo,
		},
		Errors: nil,
	})
}

// Login handles user login.
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	user, err := utils.ValideUserData(r, consts.VALDIDATION_LoginType)
	if err != nil {
		h.HandleError(w, http.StatusBadRequest, []string{consts.ErrUnexpectedValueFormat}, err)
		return
	}
	tokens, userInfo, clientErr, err := h.US.UserRequest(consts.POST_Method, consts.DB_Login, user, http.StatusOK)
	if err != nil {
		h.HandleError(w, http.StatusBadRequest, []string{clientErr}, err)
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

// GetMyProfile retrieves the logged-in user's profile.
func (h *Handler) GetMyProfile(w http.ResponseWriter, r *http.Request) {
	response, statusCode, clientErr, err := h.US.GetMyProfile(r, h.log)
	if err != nil {
		h.HandleError(w, statusCode, []string{clientErr}, err)
		return
	}

	h.SendRespond(w, http.StatusOK, response)
}

// GetUserProfile retrieves another user's profile.
func (h *Handler) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	response, statusCode, clientErr, err := h.US.GetUserInfo(r, h.log)
	if err != nil {
		h.HandleError(w, statusCode, []string{clientErr}, err)
		return
	}

	h.SendRespond(w, http.StatusOK, response)
}

// GetSearchUserResult handles searching for users.
func (h *Handler) GetSearchUserResult(w http.ResponseWriter, r *http.Request) {
	response, statusCode, clientErr, err := h.US.GetSearchUserResult(r)
	if err != nil {
		h.HandleError(w, statusCode, []string{clientErr}, err)
		return
	}

	h.SendRespond(w, http.StatusOK, response)
}

// ChangeUserData allows the user to update their information.
func (h *Handler) ChangeUserData(w http.ResponseWriter, r *http.Request) {
	response, statusCode, clientErr, err := h.US.ChangeUserData(r)
	if err != nil {
		h.HandleError(w, statusCode, []string{clientErr}, err)
		return
	}

	h.SendRespond(w, http.StatusOK, &models.Response{
		Errors: nil,
		Data:   response,
	})
}

// UploadAvatar allows the user to upload their avatar.
func (h *Handler) UploadAvatar(w http.ResponseWriter, r *http.Request) {
	log.Printf("получен запрос")
	response, statusCode, clientErr, err := h.US.AddUserAvatar(r)
	if err != nil {
		h.HandleError(w, statusCode, []string{clientErr}, err)
		return
	}

	h.SendRespond(w, http.StatusOK, response)
}

// Friend request operations
// RequestToFriendShip sends a friend request.
func (h *Handler) RequestToFriendShip(w http.ResponseWriter, r *http.Request) {
	response, statusCode, clientErr, err := h.US.FriendsOperations(r, consts.FRIENDSHIP_RequestType)
	if err != nil {
		h.HandleError(w, statusCode, []string{clientErr}, err)
		return
	}

	h.SendRespond(w, http.StatusOK, response)
}

// AcceptFriendShip accepts a friend request.
func (h *Handler) AcceptFriendShip(w http.ResponseWriter, r *http.Request) {
	response, statusCode, clientErr, err := h.US.FriendsOperations(r, consts.FRIENDSHIP_AcceptType)
	if err != nil {
		h.HandleError(w, statusCode, []string{clientErr}, err)
		return
	}

	h.SendRespond(w, http.StatusOK, response)
}

// DeleteFriend removes a user from the friend list.
func (h *Handler) DeleteFriend(w http.ResponseWriter, r *http.Request) {
	response, statusCode, clientErr, err := h.US.FriendsOperations(r, consts.FRIENDSHIP_DeleteFriendType)
	if err != nil {
		h.HandleError(w, statusCode, []string{clientErr}, err)
		return
	}

	h.SendRespond(w, http.StatusOK, response)
}

// DeleteFriendRequest cancels a friend request.
func (h *Handler) DeleteFriendRequest(w http.ResponseWriter, r *http.Request) {
	response, statusCode, clientErr, err := h.US.FriendsOperations(r, consts.FRIENDSHIP_DeleteFriendRequestType)
	if err != nil {
		h.HandleError(w, statusCode, []string{clientErr}, err)
		return
	}

	h.SendRespond(w, http.StatusOK, response)
}

// GetFriends retrieves the list of user's friends.
func (h *Handler) GetFriends(w http.ResponseWriter, r *http.Request) {
	response, statusCode, clientErr, err := h.US.FriendsOperations(r, consts.FRIENDSHIP_GetFriendsType)
	if err != nil {
		h.HandleError(w, statusCode, []string{clientErr}, err)
		return
	}

	h.SendRespond(w, http.StatusOK, response)
}

// GetSubscribers retrieves the list of user's subscribers.
func (h *Handler) GetSubscribers(w http.ResponseWriter, r *http.Request) {
	response, statusCode, clientErr, err := h.US.FriendsOperations(r, consts.FRIENDSHIP_GetSubsType)
	if err != nil {
		h.HandleError(w, statusCode, []string{clientErr}, err)
		return
	}

	h.SendRespond(w, http.StatusOK, response)
}
