package server

import (
	"chat_service/internal/logger"
	"chat_service/internal/models"
	chatservice "chat_service/internal/services/chatService"
	"chat_service/internal/services/ws_service"
	"crypto/sha256"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type WS struct {
	Upgrader *websocket.Upgrader
	Chats    map[string]*models.Chat
	Mu       *sync.RWMutex
	log      logger.Logger
}

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
	ws.log.Infof("Connection open, addr: %s", conn.RemoteAddr())

	go ws.Receiver(chat, conn)
	ws.log.Infof("connections now: %+v", chat.Clients)
	for {
		msg, err := ws_service.MessageRecording(conn, chat, ws.log)
		if err != nil {
			ws_service.SendError(conn, fmt.Sprintf("failed read message from websocket: %v", err))
			break
		}

		if !ws.isMessageValid(msg, chat, conn) {
			continue
		}

		if response, err := ws_service.SaveMessage(msg, chat.ID); err != nil {
			ws_service.SendError(conn, fmt.Sprintf("error save message: %v", response.Errors))
			continue
		}

		chat.Broadcast <- *msg
		ws.log.Infof("Message sent to Broadcast: %+v", msg)
	}
	ws_service.CloseWS(chat, conn)
	ws.log.Infof("connections without cycle: %+v", chat.Clients)
}

func (ws *WS) isMessageValid(msg *models.Message, chat *models.Chat, conn *websocket.Conn) bool {
	ws.Mu.RLock()
	defer ws.Mu.RUnlock()

	if accept := ws_service.IsValidMessage(msg, chat); !accept {
		ws_service.SendError(conn, "user is not member chat")
		return false
	}
	return true
}

func (ws *WS) Receiver(chat *models.Chat, conn *websocket.Conn) {
	for msg := range chat.Broadcast {
		ws.log.Infof("Message received from Broadcast: %+v", msg)
		hash := sha256.Sum256([]byte(msg.Body))
		ws.log.Infof("Message ID: %x", hash)
		ws.log.Infof("сообщение в ресивере")

		clients := ws.filterClients(chat, conn)
		log.Printf("After filtering, clients count: %d", len(clients))

		for _, client := range clients {
			if err := ws_service.MessageReading(client, chat, &msg); err != nil {
				ws.handleClientError(chat, client, err)
			}
			ws.log.Infof("сообщение прочитано")
		}
	}
}

func (ws *WS) filterClients(chat *models.Chat, conn *websocket.Conn) []*websocket.Conn {
	ws.Mu.RLock()
	defer ws.Mu.RUnlock()
	log.Println("функция filterClients вызвана из Receiver")
	clients := make([]*websocket.Conn, 0, len(chat.Clients))
	for client := range chat.Clients {
		log.Printf("Before filtering, clients count: %d", len(chat.Clients))
		// TODO
		// if client != conn {
		clients = append(clients, client)
		ws.log.Infof("Client added: %+v", ws)
		// }
	}
	log.Println("клиенты возвращаются в Receiver")
	return clients
}

func (ws *WS) handleClientError(chat *models.Chat, client *websocket.Conn, err error) {
	ws_service.SendError(client, fmt.Sprintf("failed reading message: %v", err))
	ws.Mu.Lock()
	defer ws.Mu.Unlock()
	delete(chat.Clients, client)
	client.Close()
}
