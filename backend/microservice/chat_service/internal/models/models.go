package models

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// User represents a user in the chat system, containing their ID, username, email, avatar, and chats.
type User struct {
	ID       int                  `json:"id"`
	Username string               `json:"username"`
	Email    string               `json:"email"`
	Avatar   *string              `json:"avatar"`
	Chats    map[int]UserChatInfo `json:"chats,omitempty"`
}

// UserChatInfo represents information about a user's participation in a specific chat, including the chat ID and messages.
type UserChatInfo struct {
	ID       int `json:"id,omitempty"`
	Messages []Message
}

// Message represents a chat message, including its ID, body, timestamp, user, and chat ID.
type Message struct {
	ID     int       `json:"id"`
	Body   string    `json:"body"`
	Time   time.Time `json:"time"`
	User   User      `json:"user"`
	ChatID int       `json:"chat_id,omitempty"`
}

// Hashtag represents a hashtag used in the chat system.
type Hashtag struct {
	ID      int    `json:"id,omitempty"`
	Hashtag string `json:"hashtag"`
}

// Chat represents a chat in the system, including its ID, creator, name, members, messages, and WebSocket clients.
type Chat struct {
	ID        int                      `json:"id"`
	Creator   int                      `json:"creator"`
	ChatName  string                   `json:"chat_name"`
	Members   map[int]User             `json:"members"`
	Messages  []Message                `json:"messages"`
	Clients   map[*websocket.Conn]bool `json:"-"` // WebSocket clients, excluded from JSON serialization.
	Broadcast chan Message             `json:"-"` // Channel for broadcasting messages.
	Mu        *sync.RWMutex            `json:"-"` // Mutex to synchronize chat operations.
}

// GetAndCreateChat represents the data structure for creating or retrieving a chat, including members, messages, and hashtags.
type GetAndCreateChat struct {
	ID       int       `json:"id"`
	Creator  int       `json:"creator"`
	ChatName string    `json:"chat_name"`
	Members  []int     `json:"members"`
	Messages []Message `json:"messages"`
	Hashtags []Hashtag `json:"hashtags"`
}

// MemberRequest represents a request to add or remove a member from a chat, including the user ID and options.
type MemberRequest struct {
	UserID  int      `json:"user_id"`
	ChatID  int      `json:"chat_id"`
	Options []string `json:"options,omitempty"`
}

// Response represents a generic response from the server, containing errors and data.
type Response struct {
	Errors []string    `json:"Errors"`
	Data   interface{} `json:"Data"`
}

// WSError represents a WebSocket error message.
type WSError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}
