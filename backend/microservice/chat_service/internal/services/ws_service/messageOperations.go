package ws_service

import (
	"chat_service/internal/logger"
	"chat_service/internal/models"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// MessageReading sends a message to the connected WebSocket client.
// It marshals the message to JSON and writes it to the WebSocket connection.
func MessageReading(client *websocket.Conn, chat *models.Chat, msg *models.Message) error {
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	if err = client.WriteMessage(websocket.TextMessage, msgBytes); err != nil {
		log.Printf("failed to write message to client: %v", err)
		client.Close()
		return err
	}

	return nil
}

// MessageRecording reads a message from the WebSocket connection.
// It unmarshals the received message and validates the user against the chat members.
func MessageRecording(conn *websocket.Conn, chat *models.Chat, log logger.Logger) (*models.Message, error) {
	_, p, err := conn.ReadMessage()
	if err != nil {
		log.Error("error in message recording on 33 line")
		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			return nil, fmt.Errorf("websocket closed: %v", err)
		}
		return nil, fmt.Errorf("failed read message in websocket: %v", err)
	}

	var message models.Message
	if err := json.Unmarshal(p, &message); err != nil {
		log.Error("error in message recording on 43 line")
		return nil, fmt.Errorf("failed unmarshal message: %v", err)
	}

	chat.Mu.RLock()
	defer chat.Mu.RUnlock()
	if _, userExists := chat.Members[message.User.ID]; !userExists {
		return nil, fmt.Errorf("failed unmarshal message: %v", err)
	}

	log.Infof("message was recording")
	return &models.Message{
		Body:   message.Body,
		User:   message.User,
		Time:   time.Now(),
		ChatID: message.ChatID,
	}, nil
}

// IsValidMessage validates a message before processing it.
// It checks that the message body is non-empty and that the user is a member of the chat.
func IsValidMessage(msg *models.Message, chat *models.Chat) bool {
	chat.Mu.RLock()
	defer chat.Mu.RUnlock()
	if msg != nil && msg.Body != "" && msg.User.ID != 0 {
		if _, ok := chat.Members[msg.User.ID]; ok {
			return true
		}
	}
	return false
}
