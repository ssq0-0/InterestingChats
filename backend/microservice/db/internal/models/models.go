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
	ID       int          `json:"id,omitempty"`
	ChatName string       `json:"chat_name"`
	Members  map[int]User `json:"members"`
	Messages []Message    `json:"messages"`
}

type AddMemberRequest struct {
	UserID int `json:"user_id"`
	ChatID int `json:"chat_id"`
}
