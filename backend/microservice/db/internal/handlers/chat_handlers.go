package handlers

import (
	"InterestingChats/backend/microservice/db/internal/consts"
	"InterestingChats/backend/microservice/db/internal/models"
	"InterestingChats/backend/microservice/db/internal/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// GetChat to retrieve detailed chat room information.
func (h *handler) GetChat(w http.ResponseWriter, r *http.Request) {
	chatIDstr, userIDstr := r.URL.Query().Get("chatID"), r.URL.Query().Get("userID")
	if chatIDstr == "" || userIDstr == "" {
		ErrorHandler(w, http.StatusBadRequest, h.log, []string{consts.ErrMissingChatName, consts.ErrMissingUserID}, nil)
		return
	}

	var userID int
	if clientErr, err := utils.ConvertValue(userIDstr, &userID); err != nil {
		ErrorHandler(w, http.StatusBadRequest, h.log, []string{clientErr}, err)
		return
	}

	var chatID int
	if clientErr, err := utils.ConvertValue(chatIDstr, &chatID); err != nil {
		ErrorHandler(w, http.StatusBadRequest, h.log, []string{clientErr}, err)
		return
	}

	chat, clientErr, err := h.DBService.GetChatInfo(chatID, userID)
	if err != nil {
		ErrorHandler(w, http.StatusBadRequest, h.log, []string{clientErr}, err)
		return
	}

	SendRespond(w, http.StatusOK, &models.Response{
		Errors: nil,
		Data:   chat,
	})
}

// GetAllChats to retrieve a list of all available chat rooms.
func (h *handler) GetAllChats(w http.ResponseWriter, r *http.Request) {
	chatList, clientErr, err := h.DBService.GetAllChats()
	if err != nil {
		ErrorHandler(w, http.StatusBadRequest, h.log, []string{clientErr}, err)
		return
	}

	SendRespond(w, http.StatusOK, &models.Response{
		Errors: nil,
		Data:   chatList,
	})
}

// GetChatBySymbols to search for chat rooms by specific characters or keywords.
func (h *handler) GetChatBySymbols(w http.ResponseWriter, r *http.Request) {
	chatList, clientErr, err := h.DBService.GetChat(r.URL.Query().Get("chatName"))
	if err != nil {
		ErrorHandler(w, http.StatusBadRequest, h.log, []string{clientErr}, err)
		return
	}

	SendRespond(w, http.StatusOK, &models.Response{
		Errors: nil,
		Data:   chatList,
	})
}

// GetUserChats to retrieve a list of chats a user is a member of.
func (h *handler) GetUserChats(w http.ResponseWriter, r *http.Request) {
	var userID int
	if clientErr, err := utils.ConvertValue(r.URL.Query().Get("userID"), &userID); err != nil {
		ErrorHandler(w, http.StatusBadRequest, h.log, []string{clientErr}, err)
		return
	}

	chatList, clientErr, err := h.DBService.GetUserChats(userID)
	if err != nil {
		ErrorHandler(w, http.StatusBadRequest, h.log, []string{clientErr}, err)
		return
	}

	SendRespond(w, http.StatusOK, &models.Response{
		Errors: nil,
		Data:   chatList,
	})
}

// CreateChat to create a new chat room.
func (h *handler) CreateChat(w http.ResponseWriter, r *http.Request) {
	var chat models.CreateChat
	if err := json.NewDecoder(r.Body).Decode(&chat); err != nil {
		ErrorHandler(w, http.StatusBadRequest, h.log, []string{consts.ErrInvalidValueFormat}, fmt.Errorf(consts.InternalErrFailedRequest, err))
		return
	}

	if clientErr, err := h.DBService.CreateChat(&chat); err != nil {
		ErrorHandler(w, http.StatusBadRequest, h.log, []string{clientErr}, err)
		return
	}

	SendRespond(w, http.StatusCreated, &models.Response{
		Errors: nil,
		Data:   chat,
	})
}

// DeleteChat to delete a chat room.
func (h *handler) DeleteChat(w http.ResponseWriter, r *http.Request) {
	chatIDStr := r.URL.Query().Get("chatID")
	if chatIDStr == "" {
		ErrorHandler(w, http.StatusBadRequest, h.log, []string{consts.ErrMissingChatID}, fmt.Errorf(consts.InternalErrMissingURLVal))
		return
	}

	var chatID int
	if clientErr, err := utils.ConvertValue(chatIDStr, &chatID); err != nil {
		ErrorHandler(w, http.StatusBadRequest, h.log, []string{clientErr}, err)
		return
	}

	if clientErr, err := h.DBService.DeleteChat(chatID); err != nil {
		ErrorHandler(w, http.StatusBadRequest, h.log, []string{clientErr}, err)
		return
	}

	SendRespond(w, http.StatusNoContent, &models.Response{
		Data:   nil,
		Errors: nil,
	})
}

// DeleteMember to remove a user from a chat room.
func (h *handler) DeleteMember(w http.ResponseWriter, r *http.Request) {
	var deleteMember models.MemberRequest
	if err := json.NewDecoder(r.Body).Decode(&deleteMember); err != nil {
		ErrorHandler(w, http.StatusBadRequest, h.log, []string{consts.ErrInvalidValueFormat}, fmt.Errorf(consts.InternalErrFailedRequest, err))
		return
	}

	if clientErr, err := h.DBService.DeleteMember(deleteMember.ChatID, deleteMember.UserID); err != nil {
		ErrorHandler(w, http.StatusBadRequest, h.log, []string{clientErr}, err)
		return
	}

	SendRespond(w, http.StatusNoContent, &models.Response{
		Data:   nil,
		Errors: nil,
	})
}

// AddMembers to add one or more members to a chat room.
func (h *handler) AddMembers(w http.ResponseWriter, r *http.Request) {
	var addMember models.MemberRequest
	if err := json.NewDecoder(r.Body).Decode(&addMember); err != nil {
		ErrorHandler(w, http.StatusBadRequest, h.log, []string{consts.ErrInvalidValueFormat}, fmt.Errorf(consts.InternalErrFailedRequest, err))
		return
	}

	if clientErr, err := h.DBService.AddMember(addMember.ChatID, addMember.UserID); err != nil || clientErr != "" {
		ErrorHandler(w, http.StatusBadRequest, h.log, []string{clientErr}, err)
		return
	}

	SendRespond(w, http.StatusAccepted, &models.Response{
		Data:   "succesfull added member",
		Errors: nil,
	})
}

// GetAuthor to retrieve information about the creator of a chat room.
func (h *handler) GetAuthor(w http.ResponseWriter, r *http.Request) {
	userIDstr, chatIDstr := r.URL.Query().Get("userID"), r.URL.Query().Get("chatID")
	if userIDstr == "" || chatIDstr == "" {
		ErrorHandler(w, http.StatusBadRequest, h.log, []string{consts.ErrMissingChatName, consts.ErrMissingUserID}, fmt.Errorf(consts.InternalErrMissingURLVal))
		return
	}

	var chatID, userID int
	if clientErr, err := utils.ConvertValue(chatIDstr, &chatID); err != nil {
		ErrorHandler(w, http.StatusBadRequest, h.log, []string{clientErr}, err)
		return
	}
	if clientErr, err := utils.ConvertValue(userIDstr, &userID); err != nil {
		ErrorHandler(w, http.StatusBadRequest, h.log, []string{clientErr}, err)
		return
	}

	if clientErr, err := h.DBService.CheckAuthor(userID, chatID); err != nil {
		ErrorHandler(w, http.StatusBadRequest, h.log, []string{clientErr}, err)
		return
	}

	SendRespond(w, http.StatusOK, &models.Response{
		Data:   userID,
		Errors: nil,
	})
}

// ChangeChatName to update the name of a chat room.
func (h *handler) ChangeChatName(w http.ResponseWriter, r *http.Request) {
	chatID, chatName := r.URL.Query().Get("chatID"), r.URL.Query().Get("chatName")
	if strings.ReplaceAll(chatID, " ", "") == "" || strings.ReplaceAll(chatName, " ", "") == "" {
		ErrorHandler(w, http.StatusBadRequest, h.log, []string{consts.ErrMissingChatName, consts.ErrMissingUserID}, fmt.Errorf(consts.InternalErrMissingURLVal))
		return
	}

	var chatIDint int
	if clientErr, err := utils.ConvertValue(chatID, &chatIDint); err != nil {
		ErrorHandler(w, http.StatusBadRequest, h.log, []string{clientErr}, err)
		return
	}

	if clientErr, err := h.DBService.ChangeChatName(chatIDint, chatName); err != nil {
		ErrorHandler(w, http.StatusBadRequest, h.log, []string{clientErr}, err)
		return
	}

	SendRespond(w, http.StatusOK, &models.Response{
		Data:   chatName,
		Errors: nil,
	})
}

// SetTag to assign tags to a chat or user.
func (h *handler) SetTag(w http.ResponseWriter, r *http.Request) {
	if statusCode, clientErr, err := h.DBService.SetTags(r); err != nil {
		ErrorHandler(w, statusCode, h.log, []string{clientErr}, err)
		return
	}

	SendRespond(w, http.StatusOK, &models.Response{
		Errors: nil,
		Data:   "successul set tags!",
	})
}

// GetTags to retrieve all tags assigned to a chat or user.
func (h *handler) GetTags(w http.ResponseWriter, r *http.Request) {
	tagList, statusCode, clientErr, err := h.DBService.GetTags(r)
	if err != nil {
		ErrorHandler(w, statusCode, h.log, []string{clientErr}, err)
		return
	}

	SendRespond(w, http.StatusOK, &models.Response{
		Errors: nil,
		Data:   tagList,
	})
}

// DeleteTags to remove tags from a chat or user.
func (h *handler) DeleteTags(w http.ResponseWriter, r *http.Request) {
	if statusCode, clientErr, err := h.DBService.DeleteTags(r); err != nil {
		ErrorHandler(w, statusCode, h.log, []string{clientErr}, err)
		return
	}

	SendRespond(w, http.StatusOK, &models.Response{
		Errors: nil,
		Data:   "successul delete tags!",
	})
}
