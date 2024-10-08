package server

import (
	"chat_service/internal/logger"
	"chat_service/internal/models"
	chatservice "chat_service/internal/services/chatService"
	"chat_service/internal/services/ws_service"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type WS struct {
	Upgrader *websocket.Upgrader
	Chats    map[string]*models.Chat
	Mu       *sync.RWMutex
	Producer *ws_service.Producer
	log      logger.Logger
}

// ChatWebsocket handles the WebSocket connection for a chat.
// It prepares the WebSocket, opens the connection, and processes incoming messages.
func (ws *WS) ChatWebsocket(w http.ResponseWriter, r *http.Request) {
	ws.Mu.Lock()
	chat, statusCode, userMsg, err := ws_service.PrepareWS(w, r, ws.Chats)
	ws.Mu.Unlock()
	if err != nil {
		chatservice.ErrorHandler(w, statusCode, ws.log, []string{userMsg}, err.Error())
		return
	}

	conn, logMsg, err := ws_service.OpenWS(w, r, chat, ws.Upgrader)
	if err != nil {
		chatservice.ErrorHandler(w, http.StatusBadGateway, ws.log, []string{err.Error()}, logMsg)
		return
	}
	defer conn.Close()

	go ws.Receiver(chat, conn)
	for {
		msg, err := ws_service.MessageRecording(conn, chat, ws.log)
		if err != nil {
			ws_service.SendError(conn, fmt.Sprintf("failed read message from websocket: %v", err))
			break
		}

		if !ws.isMessageValid(msg, chat, conn) {
			continue
		}

		if err := ws.Producer.Writer(msg); err != nil {
			ws_service.SendError(conn, "error save message")
			continue
		}

		chat.Broadcast <- *msg
	}
	ws_service.CloseWS(chat, conn)
}

// isMessageValid checks if a message is valid for the given chat and connection.
// It verifies whether the user is a member of the chat.
func (ws *WS) isMessageValid(msg *models.Message, chat *models.Chat, conn *websocket.Conn) bool {
	ws.Mu.RLock()
	defer ws.Mu.RUnlock()

	if accept := ws_service.IsValidMessage(msg, chat); !accept {
		ws_service.SendError(conn, "user is not member chat")
		return false
	}
	return true
}

// Receiver listens for broadcast messages from the chat and sends them to all connected clients.
// It also handles closed connections by removing clients from the chat.
func (ws *WS) Receiver(chat *models.Chat, conn *websocket.Conn) {
	for msg := range chat.Broadcast {
		ws.Mu.RLock()
		for client := range chat.Clients {
			if client == nil || client.WriteMessage(websocket.PingMessage, nil) != nil {
				ws.log.Infof("delete client with closed connection")
				ws.Mu.RUnlock()
				ws.handleClientError(chat, client, fmt.Errorf("connection closed"))
				ws.Mu.RLock()
				continue
			}

			if err := ws_service.MessageReading(client, chat, &msg); err != nil {
				ws.log.Errorf("failed sent message: %v", err)
				ws.handleClientError(chat, client, err)
			}
		}
		ws.Mu.RUnlock()
	}
}

// handleClientError handles errors related to a specific client.
// It sends an error message to the client and removes the client from the chat.
func (ws *WS) handleClientError(chat *models.Chat, client *websocket.Conn, err error) {
	ws_service.SendError(client, fmt.Sprintf("failed reading message: %v", err))
	ws.Mu.Lock()
	defer ws.Mu.Unlock()
	delete(chat.Clients, client)
	client.Close()
}
