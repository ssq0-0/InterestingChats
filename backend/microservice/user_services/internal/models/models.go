package models

import "time"

type BaseUser struct {
	Email    string  `json:"email"`
	Username string  `json:"username"`
	Avatar   *string `json:"avatar"`
}

type User struct {
	ID int `json:"id"`
	BaseUser
	Password string `json:"password"`
}

type Tokens struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type UserTokens struct {
	Tokens Tokens `json:"tokens"`
}

type Response struct {
	Errors []string    `json:"Errors"`
	Data   interface{} `json:"Data"`
}

type AuthResponse struct {
	Tokens Tokens `json:"tokens"`
	User   User   `json:"user"`
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
