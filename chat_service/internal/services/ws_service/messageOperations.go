package ws_service

import (
	"chat_service/internal/consts"
	"chat_service/internal/logger"
	"chat_service/internal/models"
	"chat_service/internal/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

func MessageReading(client *websocket.Conn, chat *models.Chat, msg *models.Message) error {
	log.Printf("сообщение читается: %+v", msg)
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
		Body: message.Body,
		User: message.User,
		Time: time.Now(),
	}, nil
}

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

func SaveMessage(msg *models.Message, chatID int) (*models.Response, error) {
	// TODO: if response....
	if _, _, _, err := utils.ProxyRequest(consts.POST_Method, fmt.Sprintf(consts.DB_SaveMessage, chatID), msg, http.StatusOK); err != nil {
		return nil, err
	}
	log.Println("сообщение отправлено в базу данных")
	return nil, nil
}
