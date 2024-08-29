package handlers

import (
	"InterestingChats/backend/microservice/db/internal/consts"
	"InterestingChats/backend/microservice/db/internal/db"
	"InterestingChats/backend/microservice/db/internal/handlers"
	"InterestingChats/backend/microservice/db/internal/logger"
	"InterestingChats/backend/microservice/db/internal/models"
	"InterestingChats/backend/microservice/db/internal/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

type ChatHandler struct {
	ChatService *db.ChatService
	log         logger.Logger
}

func NewChatHandler(chatService *db.ChatService, log logger.Logger) *ChatHandler {
	if chatService == nil {
		log.Fatal("ChatService is nil")
	}
	return &ChatHandler{
		ChatService: chatService,
		log:         log,
	}
}

func (ch *ChatHandler) GetChat(w http.ResponseWriter, r *http.Request) {
	chatIDstr, userIDstr := r.URL.Query().Get("chatID"), r.URL.Query().Get("userID")
	if chatIDstr == "" || userIDstr == "" {
		handlers.ErrorHandler(w, http.StatusBadRequest, ch.log, []string{consts.ErrMissingChatName, consts.ErrMissingUserID}, nil)
		return
	}

	var userID int
	if clientErr, err := utils.ConvertValue(userIDstr, &userID); err != nil {
		handlers.ErrorHandler(w, http.StatusBadRequest, ch.log, []string{clientErr}, err)
		return
	}

	var chatID int
	if clientErr, err := utils.ConvertValue(chatIDstr, &chatID); err != nil {
		handlers.ErrorHandler(w, http.StatusBadRequest, ch.log, []string{clientErr}, err)
		return
	}

	chat, clientErr, err := ch.ChatService.GetChatInfo(chatID, userID)
	if err != nil {
		handlers.ErrorHandler(w, http.StatusBadRequest, ch.log, []string{clientErr}, err)
		return
	}

	handlers.SendRespond(w, http.StatusOK, chat)
}

func (ch *ChatHandler) GetAllChats(w http.ResponseWriter, r *http.Request) {
	chatList, clientErr, err := ch.ChatService.GetAllChats()
	if err != nil {
		handlers.ErrorHandler(w, http.StatusBadRequest, ch.log, []string{clientErr}, err)
		return
	}

	handlers.SendRespond(w, http.StatusOK, &models.Response{
		Errors: nil,
		Data:   chatList,
	})
}

func (ch *ChatHandler) GetChatBySymbols(w http.ResponseWriter, r *http.Request) {
	chatList, clientErr, err := ch.ChatService.GetChat(r.URL.Query().Get("chatName"))
	if err != nil {
		handlers.ErrorHandler(w, http.StatusBadRequest, ch.log, []string{clientErr}, err)
		return
	}

	handlers.SendRespond(w, http.StatusOK, &models.Response{
		Errors: nil,
		Data:   chatList,
	})
}

func (ch *ChatHandler) GetUserChats(w http.ResponseWriter, r *http.Request) {
	var userID int
	if clientErr, err := utils.ConvertValue(r.URL.Query().Get("userID"), &userID); err != nil {
		handlers.ErrorHandler(w, http.StatusBadRequest, ch.log, []string{clientErr}, err)
		return
	}

	chatList, clientErr, err := ch.ChatService.GetUserChats(userID)
	if err != nil {
		handlers.ErrorHandler(w, http.StatusBadRequest, ch.log, []string{clientErr}, err)
		return
	}

	handlers.SendRespond(w, http.StatusOK, &models.Response{
		Errors: nil,
		Data:   chatList,
	})
}

func (ch *ChatHandler) CreateChat(w http.ResponseWriter, r *http.Request) {
	var chat models.CreateChat
	if err := json.NewDecoder(r.Body).Decode(&chat); err != nil {
		handlers.ErrorHandler(w, http.StatusBadRequest, ch.log, []string{consts.ErrInvalidValueFormat}, fmt.Errorf(consts.InternalErrFailedRequest, err))
		return
	}

	if clientErr, err := ch.ChatService.CreateChat(&chat); err != nil {
		handlers.ErrorHandler(w, http.StatusBadRequest, ch.log, []string{clientErr}, err)
		return
	}

	handlers.SendRespond(w, http.StatusCreated, chat)
}

func (ch *ChatHandler) DeleteChat(w http.ResponseWriter, r *http.Request) {
	chatIDStr := r.URL.Query().Get("chatID")
	if chatIDStr == "" {
		handlers.ErrorHandler(w, http.StatusBadRequest, ch.log, []string{consts.ErrMissingChatID}, fmt.Errorf(consts.InternalErrMissingURLVal))
		return
	}

	var chatID int
	if clientErr, err := utils.ConvertValue(chatIDStr, &chatID); err != nil {
		handlers.ErrorHandler(w, http.StatusBadRequest, ch.log, []string{clientErr}, err)
		return
	}

	if clientErr, err := ch.ChatService.DeleteChat(chatID); err != nil {
		handlers.ErrorHandler(w, http.StatusBadRequest, ch.log, []string{clientErr}, err)
		return
	}

	handlers.SendRespond(w, http.StatusNoContent, &models.Response{
		Data:   nil,
		Errors: nil,
	})
}

func (ch *ChatHandler) DeleteMember(w http.ResponseWriter, r *http.Request) {
	var deleteMember models.MemberRequest
	if err := json.NewDecoder(r.Body).Decode(&deleteMember); err != nil {
		handlers.ErrorHandler(w, http.StatusBadRequest, ch.log, []string{consts.ErrInvalidValueFormat}, fmt.Errorf(consts.InternalErrFailedRequest, err))
		return
	}

	if clientErr, err := ch.ChatService.DeleteMember(deleteMember.ChatID, deleteMember.UserID); err != nil {
		handlers.ErrorHandler(w, http.StatusBadRequest, ch.log, []string{clientErr}, err)
		return
	}

	handlers.SendRespond(w, http.StatusNoContent, &models.Response{
		Data:   nil,
		Errors: nil,
	})
}

func (ch *ChatHandler) AddMembers(w http.ResponseWriter, r *http.Request) {
	var addMember models.MemberRequest
	if err := json.NewDecoder(r.Body).Decode(&addMember); err != nil {
		handlers.ErrorHandler(w, http.StatusBadRequest, ch.log, []string{consts.ErrInvalidValueFormat}, fmt.Errorf(consts.InternalErrFailedRequest, err))
		return
	}

	if clientErr, err := ch.ChatService.AddMember(addMember.ChatID, addMember.UserID); err != nil || clientErr != "" {
		handlers.ErrorHandler(w, http.StatusBadRequest, ch.log, []string{clientErr}, err)
		return
	}

	handlers.SendRespond(w, http.StatusAccepted, &models.Response{
		Data:   "succesfull added member",
		Errors: nil,
	})
}

func (ch *ChatHandler) GetAuthor(w http.ResponseWriter, r *http.Request) {
	userIDstr, chatIDstr := r.URL.Query().Get("userID"), r.URL.Query().Get("chatID")
	if userIDstr == "" || chatIDstr == "" {
		handlers.ErrorHandler(w, http.StatusBadRequest, ch.log, []string{consts.ErrMissingChatName, consts.ErrMissingUserID}, fmt.Errorf(consts.InternalErrMissingURLVal))
		return
	}

	var chatID, userID int
	if clientErr, err := utils.ConvertValue(chatIDstr, &chatID); err != nil {
		handlers.ErrorHandler(w, http.StatusBadRequest, ch.log, []string{clientErr}, err)
		return
	}
	if clientErr, err := utils.ConvertValue(userIDstr, &userID); err != nil {
		handlers.ErrorHandler(w, http.StatusBadRequest, ch.log, []string{clientErr}, err)
		return
	}

	if clientErr, err := ch.ChatService.CheckAuthor(userID, chatID); err != nil {
		handlers.ErrorHandler(w, http.StatusBadRequest, ch.log, []string{clientErr}, err)
		return
	}

	handlers.SendRespond(w, http.StatusOK, &models.Response{
		Data:   userID,
		Errors: nil,
	})
}

func (ch *ChatHandler) SaveMessage(w http.ResponseWriter, r *http.Request) {
	var chatID int
	if userErr, err := utils.ConvertValue(r.URL.Query().Get("chatID"), &chatID); err != nil {
		handlers.ErrorHandler(w, http.StatusBadRequest, ch.log, []string{userErr}, err)
		return
	}
	ch.log.Infof("chat id here: %d", chatID)
	var message models.Message
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		handlers.ErrorHandler(w, http.StatusBadRequest, ch.log, []string{consts.ErrInvalidValueFormat}, fmt.Errorf(consts.InternalErrFailedRequest, err))
		return
	}

	if clientErr, err := ch.ChatService.SaveMessage(message, chatID); err != nil {
		handlers.ErrorHandler(w, http.StatusBadRequest, ch.log, []string{clientErr}, err)
		return
	}

	handlers.SendRespond(w, http.StatusOK, &models.Response{})
}
