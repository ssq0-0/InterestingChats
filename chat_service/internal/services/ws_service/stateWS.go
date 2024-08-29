package ws_service

import (
	"chat_service/internal/models"
	chatservice "chat_service/internal/services/chatService"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

func PrepareWS(w http.ResponseWriter, r *http.Request, chats map[string]*models.Chat) (*models.Chat, int, string, error) {
	chatHistory, statusCode, clientErr, err := chatservice.GetChatHistory(r)
	if err != nil {
		return nil, statusCode, clientErr, err
	}

	if _, exists := chats[chatHistory.ChatName]; !exists {
		if chatHistory.Mu == nil {
			chatHistory.Mu = &sync.RWMutex{}
		}
		if chatHistory.Clients == nil {
			chatHistory.Clients = make(map[*websocket.Conn]bool)
		}
		if chatHistory.Broadcast == nil {
			chatHistory.Broadcast = make(chan models.Message)
		}
		if chatHistory.Members == nil {
			chatHistory.Members = make(map[int]models.User)
		}

		chats[chatHistory.ChatName] = chatHistory
	} else {
		chatHistory = chats[chatHistory.ChatName]
	}
	return chatHistory, http.StatusOK, "", nil
}

func OpenWS(w http.ResponseWriter, r *http.Request, chat *models.Chat, upgrader *websocket.Upgrader) (*websocket.Conn, string, error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("failed open ws connection: %v", err)
		return nil, fmt.Sprintf("failed open ws: %v", err), fmt.Errorf("failed connection")
	}

	chat.Mu.Lock()
	if _, exists := chat.Clients[conn]; !exists {
		chat.Clients[conn] = true
	}
	chat.Mu.Unlock()

	return conn, "", nil
}

func CloseWS(chat *models.Chat, conn *websocket.Conn) {
	chat.Mu.Lock()
	delete(chat.Clients, conn)
	chat.Mu.Unlock()
}
