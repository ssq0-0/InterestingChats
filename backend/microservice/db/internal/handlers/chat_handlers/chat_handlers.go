package handlers

import (
	"InterestingChats/backend/microservice/db/internal/db"
	"InterestingChats/backend/microservice/db/internal/models"
	"encoding/json"
	"log"
	"net/http"
)

type ChatHandler struct {
	ChatService *db.ChatService
}

func NewChatHandler(chatService *db.ChatService) *ChatHandler {
	return &ChatHandler{
		ChatService: chatService,
	}
}

func (ch *ChatHandler) GetChat(w http.ResponseWriter, r *http.Request) {
	chatName := r.URL.Query().Get("chatName")
	log.Println(chatName)
	chat, err := ch.ChatService.GetChatInfo(chatName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println("Failed get info about chat", err)
		return
	}
	log.Println(chat)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(chat); err != nil {
		http.Error(w, "error encode json", http.StatusInternalServerError)
		log.Println("error encoding json for response", err)
	}
}

func (ch *ChatHandler) CreateChat(w http.ResponseWriter, r *http.Request) {
	var chat models.Chat
	if err := json.NewDecoder(r.Body).Decode(&chat); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("failed decode json to chat: %v", err)
		return
	}

	err := ch.ChatService.CreateChat(&chat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("failed create chat: %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(chat); err != nil {
		http.Error(w, "error encode json", http.StatusInternalServerError)
		log.Println("error encoding json for response", err)
	}
}
