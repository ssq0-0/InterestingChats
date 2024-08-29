package server_test

// import (
// 	"chat_service/internal/models"
// 	"chat_service/internal/server"
// 	"net/http"
// 	"net/http/httptest"
// 	"sync"
// 	"testing"
// 	"time"

// 	"github.com/gorilla/websocket"
// 	"github.com/stretchr/testify/assert"
// )

// func TestWebSocketMessageExchange(t *testing.T) {
// 	// Создаем WebSocket сервер
// 	wsServer := &server.WS{
// 		Upgrader: &websocket.Upgrader{
// 			ReadBufferSize:  1024,
// 			WriteBufferSize: 1024,
// 			CheckOrigin: func(r *http.Request) bool {
// 				return true
// 			},
// 		},
// 		Chats: make(map[string]*models.Chat),
// 		Mu:    &sync.RWMutex{},
// 	}

// 	srv := httptest.NewServer(http.HandlerFunc(wsServer.ChatWebsocket))
// 	defer srv.Close()

// 	wsURL := "ws" + srv.URL[4:] + "/wsOpen?chatID=193&Authorization=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6NDAsImVtYWlsIjoibWF4QGdtYWlsLmNvbSIsInVzZXJuYW1lIjoibWF4IiwiZXhwIjoxNzI0NjkxNDA1fQ.-MMKW_GO7_2Cg1A0o140pE10XyEr8srEwEfBI1KWxrQ"

// 	senderWS, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
// 	if err != nil {
// 		t.Fatalf("Failed to dial WebSocket for sender: %v", err)
// 	}
// 	defer senderWS.Close()

// 	receiverWS, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
// 	if err != nil {
// 		t.Fatalf("Failed to dial WebSocket for receiver: %v", err)
// 	}
// 	defer receiverWS.Close()
// 	receiverWS.SetReadDeadline(time.Now().Add(10 * time.Second))

// 	t.Log("Sender and receiver connected")

// 	msg := models.Message{
// 		Body: "test message",
// 		User: models.User{
// 			ID:       1,
// 			Username: "testuser",
// 		},
// 		Time: time.Now(),
// 	}
// 	err = senderWS.WriteJSON(msg)
// 	if err != nil {
// 		t.Fatalf("Failed to write message: %v", err)
// 	}
// 	t.Log("Message sent")

// 	var receivedMsg models.Message
// 	err = receiverWS.ReadJSON(&receivedMsg)
// 	if err != nil {
// 		t.Fatalf("Failed to read message: %v", err)
// 	}
// 	t.Log("Message received")

// 	assert.Equal(t, msg.Body, receivedMsg.Body, "Received message body does not match sent message body")
// 	assert.Equal(t, msg.User.Username, receivedMsg.User.Username, "Received message username does not match sent message username")
// }
