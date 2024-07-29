package handlers

import (
	"InterestingChats/backend/microservice/db/internal/db"
	"InterestingChats/backend/microservice/db/internal/handlers"
	"InterestingChats/backend/microservice/db/internal/models"
	"InterestingChats/backend/microservice/db/internal/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type ChatHandler struct {
	ChatService *db.ChatService
}

func NewChatHandler(chatService *db.ChatService) *ChatHandler {
	if chatService == nil {
		log.Fatal("ChatService is nil")
	}
	return &ChatHandler{
		ChatService: chatService,
	}
}

func (ch *ChatHandler) GetChat(w http.ResponseWriter, r *http.Request) {
	chatName := r.URL.Query().Get("chatName")
	if chatName == "" {
		handlers.ErrorHandler(w, http.StatusBadRequest, []string{"Missing chat name in url"}, "Missing chat name in url")
		return
	}

	chat, err := ch.ChatService.GetChatInfo(chatName)
	if err != nil {
		handlers.ErrorHandler(w, http.StatusBadRequest, []string{"Failed to get chat info"}, fmt.Sprintf("Failed to get chat info: %v", err))
		return
	}
	handlers.SendRespond(w, http.StatusOK, chat)
}

func (ch *ChatHandler) CreateChat(w http.ResponseWriter, r *http.Request) {
	var chat models.Chat
	if err := json.NewDecoder(r.Body).Decode(&chat); err != nil {
		handlers.ErrorHandler(w, http.StatusBadRequest, []string{"Failed recived data"}, fmt.Sprintf("failed recived data: %v", err))
		return
	}

	chatID, err := ch.ChatService.CreateChat(&chat)
	if err != nil {
		handlers.ErrorHandler(w, http.StatusBadRequest, []string{"Failed create chat"}, fmt.Sprintf("failed create chat: %v", err))
		return
	}
	chat.ID = chatID

	handlers.SendRespond(w, http.StatusCreated, chat)
}

func (ch *ChatHandler) DeleteChat(w http.ResponseWriter, r *http.Request) {
	chatIDStr := r.URL.Query().Get("chatID")
	if chatIDStr == "" {
		handlers.ErrorHandler(w, http.StatusBadRequest, []string{"missing chat id in request"}, "missing chat id in request")
		return
	}

	var chatID int
	if err := utils.ConvertValue(chatIDStr, &chatID); err != nil {
		handlers.ErrorHandler(w, http.StatusBadRequest, []string{"failed convert value"}, fmt.Sprintf("failed convert value: %v", err))
		return
	}
	// ////
	err := ch.ChatService.DeleteChat(chatID)
	if err != nil {
		handlers.ErrorHandler(w, http.StatusBadRequest, []string{"failed delete chat"}, fmt.Sprintf("failed delete chat: %v", err))
		return
	}

	handlers.SendRespond(w, http.StatusNoContent, &models.Response{
		Data:   nil,
		Errors: nil,
	})
}

func (ch *ChatHandler) DeleteMember(w http.ResponseWriter, r *http.Request) {
	var deleteMember models.AddMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&deleteMember); err != nil {
		handlers.ErrorHandler(w, http.StatusBadRequest, []string{"Failed recived data"}, fmt.Sprintf("failed recived data: %v", err))
		return
	}

	err := ch.ChatService.DeleteMember(deleteMember.ChatID, deleteMember.UserID)
	if err != nil {
		handlers.ErrorHandler(w, http.StatusBadRequest, []string{"Failed delete member"}, fmt.Sprintf("Failed delete member: %v", err))
		return
	}

	handlers.SendRespond(w, http.StatusNoContent, &models.Response{
		Data:   "successful deleted user",
		Errors: nil,
	})
}

func (ch *ChatHandler) AddMembers(w http.ResponseWriter, r *http.Request) {
	var addMember models.AddMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&addMember); err != nil {
		handlers.ErrorHandler(w, http.StatusBadRequest, []string{"Failed recived data"}, fmt.Sprintf("failed recived data: %v", err))
		return
	}

	err := ch.ChatService.AddMember(addMember.ChatID, addMember.UserID)
	if err != nil {
		handlers.ErrorHandler(w, http.StatusBadRequest, []string{"failed operation"}, fmt.Sprintf("failed operation: %v", err))
		return
	}

	handlers.SendRespond(w, http.StatusAccepted, &models.Response{
		Data:   "succesfull added member",
		Errors: nil,
	})
}

func (ch *ChatHandler) GetAuthor(w http.ResponseWriter, r *http.Request) {
	if ch.ChatService == nil {
		log.Fatal("ChatService is nil")
	}

	userIDstr, chatIDstr := r.URL.Query().Get("userID"), r.URL.Query().Get("chatID")
	if userIDstr == "" || chatIDstr == "" {
		http.Error(w, "missing email or chat name", http.StatusBadRequest)
		log.Printf("missing email or chat name")
		return
	}

	var chatID, userID int
	if err := utils.ConvertValue(chatIDstr, &chatID); err != nil {
		handlers.ErrorHandler(w, http.StatusBadRequest, []string{"failed convert chat id value"}, fmt.Sprintf("failed convert value: %v", err))
		return
	}
	if err := utils.ConvertValue(userIDstr, &userID); err != nil {
		handlers.ErrorHandler(w, http.StatusBadRequest, []string{"failed convert user id value"}, fmt.Sprintf("failed convert value: %v", err))
		return
	}

	accept, err := ch.ChatService.CheckAuthor(userID, chatID)
	if err != nil {
		handlers.ErrorHandler(w, http.StatusBadRequest, []string{"error checking author"}, fmt.Sprintf("error checking author: %v", err))
		return
	}
	if !accept {
		handlers.ErrorHandler(w, http.StatusForbidden, []string{"user is not the author"}, fmt.Sprintf("user email: %d is not the author of chat %d", userID, chatID))
		return
	}
	handlers.SendRespond(w, http.StatusOK, &models.Response{
		Data:   userID,
		Errors: nil,
	})
}
