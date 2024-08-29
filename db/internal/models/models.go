package models

import "time"

type BaseUser struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type User struct {
	ID int `json:"id"`
	BaseUser
	Password string `json:"password,omitempty"`
}

type Message struct {
	ID   int       `json:"id"`
	Body string    `json:"body"`
	Time time.Time `json:"time"`
	User User      `json:"user"`
}

type Chat struct {
	ID       int          `json:"id,omitempty"`
	Creator  int          `json:"creator"`
	ChatName string       `json:"chat_name"`
	Members  map[int]User `json:"members"`
	Messages []Message    `json:"messages"`
}

type CreateChat struct {
	ID       int       `json:"id"`
	Creator  int       `json:"creator"`
	ChatName string    `json:"chat_name"`
	Members  []int     `json:"members"`
	Messages []Message `json:"messages"`
}

type MemberRequest struct {
	UserID int `json:"user_id"`
	ChatID int `json:"chat_id"`
}

type Response struct {
	Errors []string    `json:"Errors"`
	Data   interface{} `json:"Data"`
}

type ChangeUserData struct {
	Type   string `json:"type"`
	Data   string `json:"Data"`
	UserID int    `json:"user_id"`
}
