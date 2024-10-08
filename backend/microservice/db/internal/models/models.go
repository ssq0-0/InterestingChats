package models

import (
	"time"
)

type BaseUser struct {
	Username string  `json:"username"`
	Email    string  `json:"email"`
	Avatar   *string `json:"avatar"`
}

type User struct {
	ID int `json:"id"`
	BaseUser
	Password string `json:"password,omitempty"`
}

type Message struct {
	ID     int       `json:"id"`
	Body   string    `json:"body"`
	Time   time.Time `json:"time"`
	User   User      `json:"user"`
	ChatID int       `json:"chat_id,omitempty"`
}

type Hashtag struct {
	ID      int    `json:"id,omitempty"`
	Hashtag string `json:"hashtag"`
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
	Hashtags []Hashtag `json:"hashtags"`
}

type MemberRequest struct {
	UserID  int      `json:"user_id"`
	ChatID  int      `json:"chat_id"`
	Options []string `json:"options,omitempty"`
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

type UserFile struct {
	UserID int          `json:"user_id"`
	URL    FileResponse `json:"file_url"`
}

type FileResponse struct {
	Errors        []string `json:"Errors"`
	TemporaryLink string   `json:"temporary_url"`
	StaticLink    string   `json:"static_url"`
}

type Notification struct {
	ID       int
	UserID   int
	SenderID int
	Type     string
	Message  string
	Time     time.Time
	IsRead   bool
}

type NotificationRequest struct {
	Sender   User
	Receiver int
}

type FriendRequest struct {
	UserID   int         `json:"user_id"`
	FriendID int         `json:"friend_id"`
	Friends  interface{} `json:"friend_list"`
	Type     string      `json:"type"`
}
