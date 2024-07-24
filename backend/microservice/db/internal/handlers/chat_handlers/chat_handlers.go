package handlers

import (
	"InterestingChats/backend/microservice/db/internal/db"
	"InterestingChats/backend/microservice/db/internal/models"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
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

	chatID, err := ch.ChatService.CreateChat(&chat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("failed create chat: %v", err)
		return
	}
	chat.ID = chatID

	log.Printf("chat id in handler after func: %d. Func return: %d", chat.ID, chatID)
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(chat); err != nil {
		http.Error(w, "error encode json", http.StatusInternalServerError)
		log.Println("error encoding json for response", err)
	}
	log.Printf("chat struct: %+v", chat)
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
	var deleteMember models.AddMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&deleteMember); err != nil {
		log.Printf("error decoding from request body: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := ch.ChatService.DeleteMember(deleteMember.ChatID, deleteMember.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("failed delete from db: %v", err)
		return
	}

	log.Println("user deleted from chat")
	w.WriteHeader(http.StatusNoContent)
}

func (ch *ChatHandler) AddMembers(w http.ResponseWriter, r *http.Request) {
	var addMember models.AddMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&addMember); err != nil {
		log.Printf("error decode json to array models: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := ch.ChatService.AddMember(addMember.ChatID, addMember.UserID)
	if err != nil {
		log.Printf("failed operation: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Println("user add to chat")
	w.WriteHeader(http.StatusAccepted)
}

func (ch *ChatHandler) GetAuthor(w http.ResponseWriter, r *http.Request) {
	if ch.ChatService == nil {
		log.Fatal("ChatService is nil")
	}

	userEmail := r.URL.Query().Get("email")
	chatIDstr := r.URL.Query().Get("chatID")
	if userEmail == "" || chatIDstr == "" {
		http.Error(w, "missing email or chat name", http.StatusBadRequest)
		log.Printf("missing email or chat name")
		return
	}

	chatIDint, err := strconv.Atoi(chatIDstr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("failed conver chat id to int: %v", err)
		return
	}

	accept, err := ch.ChatService.CheckAuthor(userEmail, chatIDint)
	if err != nil {
		http.Error(w, "error checking author", http.StatusForbidden)
		log.Printf("error checking author: %v", err)
		return
	}
	if !accept {
		http.Error(w, "user is not the author", http.StatusForbidden)
		log.Printf("user email: %s is not the author of chat %d", userEmail, chatIDint)
		return
	}

	w.WriteHeader(http.StatusOK)
}
