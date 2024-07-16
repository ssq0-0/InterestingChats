package server

import (
	"chat_service/internal/models"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

func (s *Server) GetChatHistory(w http.ResponseWriter, r *http.Request) {
	chatName := r.URL.Query().Get("chatName")
	if chatName == "" {
		http.Error(w, "chatName is required", http.StatusBadRequest)
		return
	}
	body, statusCode, err := SendRequestToDB(chatName)
	if err != nil {
		log.Println("failed get info about chat", err)
		http.Error(w, "failed get info about chat", http.StatusBadRequest)
		return
	}
	if statusCode == http.StatusNotFound {
		log.Println("chat not found")
		http.Error(w, "chat not found", http.StatusBadRequest)
		return
	}

	chat := &models.Chat{}
	if err := json.Unmarshal(body, chat); err != nil {
		log.Println("failed to deserialize chat data", err)
		http.Error(w, "failed to deserialize chat data", http.StatusInternalServerError)
		return
	}

	log.Printf("Chat data: %+v\n", chat)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(chat); err != nil {
		http.Error(w, "error encode json", http.StatusInternalServerError)
		log.Println("error encoding json for response", err)
	}
}

func (s *Server) OpenWS(w http.ResponseWriter, r *http.Request) {
	chatName := r.URL.Query().Get("chatName")
	if chatName == "" {
		http.Error(w, "chatName is required", http.StatusBadRequest)
		return
	}

	authHeader := r.Header.Get("access_token")
	if authHeader == "" {
		log.Println("token not found")
		http.Error(w, "token not found", http.StatusBadRequest)
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == "" {
		log.Println("Invalid Authorization header format")
		http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
		return
	}

	email, err := SendToUserService(tokenString)
	if err != nil || email == "" {
		log.Printf("error sending or nil email: %v", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	body, statusCode, err := SendRequestToDB(chatName)
	if err != nil {
		log.Println("failed get info about chat", err)
		http.Error(w, "failed get info about chat", http.StatusBadRequest)
		return
	}
	if statusCode == http.StatusNotFound {
		log.Println("chat not found")
		http.Error(w, "chat not found", http.StatusBadRequest)
		return
	}

	chat := &models.Chat{}
	if err := json.Unmarshal(body, chat); err != nil {
		log.Println("failed to deserialize chat data", err)
		http.Error(w, "failed to deserialize chat data", http.StatusInternalServerError)
		return
	}

	log.Printf("email: %s", email)
	log.Printf("chat details: ID: %d, Name: %s, Members: %v, Messages: %v", chat.ID, chat.ChatName, chat.Members, chat.Messages)

	userFound := false
	for _, user := range chat.Members {
		if user.Email == email {
			userFound = true
			break
		}
	}

	if !userFound {
		log.Println("member not in chat. Rejected to open")
		http.Error(w, "member not in chat. Rejected to open", http.StatusUnauthorized)
		return
	}

	conn, err := s.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("failed open ws connection: %v", err)
		return
	}

	if chat.Mu == nil {
		chat.Mu = &sync.Mutex{}
	}
	if chat.Clients == nil {
		chat.Clients = make(map[*websocket.Conn]bool)
	}
	if chat.Broadcast == nil {
		chat.Broadcast = make(chan models.Message)
	}

	chat.Mu.Lock()
	chat.Clients[conn] = true
	chat.Mu.Unlock()

	defer func() {
		chat.Mu.Lock()
		delete(chat.Clients, conn)
		chat.Mu.Unlock()
		conn.Close()
	}()
	go s.Reciver(chat)
	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Printf("failed read message in ws: %v", err)
			return
		}

		message := string(p)
		parts := strings.SplitN(message, ":", 2)
		if len(parts) != 2 {
			log.Printf("len parts < 2: %d", len(parts))
			continue
		}
		userID, err := strconv.Atoi(parts[0])
		if err != nil {
			log.Printf("error convert to int: %v", err)
			continue
		}

		chat.Mu.Lock()
		_, userExists := chat.Members[userID]
		chat.Mu.Unlock()
		if !userExists {
			log.Printf("user %d not a member of chat", userID)
			continue
		}

		msg := models.Message{Body: parts[1], UserID: userID, Time: time.Now()}
		chat.Broadcast <- msg
	}
}

func (s *Server) Reciver(chat *models.Chat) {
	for {
		msg := <-chat.Broadcast
		responseMessage := msg.Body + ":" + msg.Time.String() + strconv.Itoa(msg.UserID)

		chat.Mu.Lock()
		for client := range chat.Clients {
			err := client.WriteMessage(websocket.TextMessage, []byte(responseMessage))
			if err != nil {
				log.Printf("failed send message: %v", err)
				client.Close()
				delete(chat.Clients, client)
			}
		}
		chat.Mu.Unlock()
	}
}
