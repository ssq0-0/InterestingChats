package handlers

import (
	"InterestingChats/backend/microservice/db/internal/db"
	"InterestingChats/backend/microservice/db/internal/models"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
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

func (ch *ChatHandler) DeleteChat(w http.ResponseWriter, r *http.Request) {
	chatIDStr := r.URL.Query().Get("chatID")
	if chatIDStr == "" {
		http.Error(w, "missing chat id in request", http.StatusBadGateway)
		log.Printf("missing chat id in request")
		return
	}

	chatID, err := strconv.Atoi(chatIDStr)
	if err != nil {
		http.Error(w, "failed convet id to int value", http.StatusBadRequest)
		log.Printf("failed convet id to int value: %v", err)
	}
	err = ch.ChatService.DeleteChat(chatID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("failed delete chat: %v", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (ch *ChatHandler) DeleteMember(w http.ResponseWriter, r *http.Request) {
	userIDsStr := r.URL.Query().Get("userID")
	chatIDstr := r.URL.Query().Get("chatID")
	if userIDsStr == "" || chatIDstr == "" {
		http.Error(w, "missing user id in request", http.StatusBadRequest)
		log.Printf("missing user id")
		return
	}

	chatIDint, err := strconv.Atoi(chatIDstr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println("failed chatID convert to int")
		return
	}

	userIDs := strings.Split(userIDsStr, ",")
	for _, userIDStr := range userIDs {
		userIDint, err := strconv.Atoi(userIDStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Println("failed userID convert to int")
			return
		}

		err = ch.ChatService.DeleteMember(chatIDint, userIDint)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Println("failed delete from db")
			return
		}
	}

	log.Println("user deleted from chat")
	w.WriteHeader(http.StatusNoContent)
}

func (ch *ChatHandler) AddMembers(w http.ResponseWriter, r *http.Request) {
	userIDsStr := r.URL.Query().Get("userID")
	chatIDstr := r.URL.Query().Get("chatID")

	if userIDsStr == "" || chatIDstr == "" {
		http.Error(w, "missing chat id or user id's", http.StatusBadRequest)
		log.Printf("missing user id or chat id")
		return
	}

	chatIDint, err := strconv.Atoi(chatIDstr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println("failed chatID convert to int")
		return
	}

	userIDArray := strings.Split(userIDsStr, ",")
	for _, userIDstr := range userIDArray {
		userIDint, err := strconv.Atoi(userIDstr)
		if err != nil {
			log.Printf("failed convert to int: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = ch.ChatService.AddMember(chatIDint, userIDint)
		if err != nil {
			log.Printf("failed operation: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	log.Println("user add to chat")
	w.WriteHeader(http.StatusAccepted)
}
