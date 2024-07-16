package models

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type User struct {
	ID       int                  `json:"id"`
	Username string               `json:"username"`
	Email    string               `json:"email"`
	Chats    map[int]UserChatInfo `json:"chats,omitempty"`
}

type UserChatInfo struct {
	Messages []Message
}

type Message struct {
	ID     int       `json:"id"`
	Body   string    `json:"body"`
	Time   time.Time `json:"time"`
	UserID int       `json:"user_id"`
}

type Chat struct {
	ID        int    `json:"id"`
	ChatName  string `json:"chat_name"`
	Members   map[int]User
	Messages  []Message
	Clients   map[*websocket.Conn]bool `json:"-"`
	Broadcast chan Message             `json:"-"`
	Mu        *sync.Mutex              `json:"-"`
}
