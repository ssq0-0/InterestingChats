package ws_service

import (
	"chat_service/internal/models"
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

func SendError(conn *websocket.Conn, msg string) {
	errMsg := &models.WSError{
		Type:    "error",
		Message: msg,
	}
	msgBytes, err := json.Marshal(errMsg)
	if err != nil {
		log.Printf("failed marshal ws error: %v", err)
		return
	}

	if err := conn.WriteMessage(websocket.TextMessage, msgBytes); err != nil {
		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			log.Printf("failed to write message: %v", err)
		}
		return
	}
}
