package models

type Response struct {
	Errors []string    `json:"Errors"`
	Data   interface{} `json:"Data"`
}

type Session struct {
	ID       string  `json:"session_id"`
	UserID   int     `json:"id"`
	Username string  `json:"username"`
	Email    string  `json:"email"`
	Avatar   *string `json:"avatar"`
}

type ChangeUserData struct {
	Type   string `json:"type"`
	Data   string `json:"Data"`
	UserID int    `json:"user_id"`
}

type Friend struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
}

type FriendOperation struct {
	UserID   int      `json:"user_id"`
	Friend   []Friend `json:"friend_list"`
	FriendID int      `json:"friend_id"`
}
