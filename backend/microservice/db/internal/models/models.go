package models

import "time"

type User struct {
	ID       int    `json:"id,omitempty"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
}

type Message struct {
	ID     int       `json:"id"`
	Body   string    `json:"body"`
	Time   time.Time `json:"time"`
	UserID int       `json:"user_id"`
}

type Chat struct {
	ID       int          `json:"id"`
	ChatName string       `json:"chat_name"`
	Members  map[int]User `json:"members"`
	Messages []Message    `json:"messages"`
}
