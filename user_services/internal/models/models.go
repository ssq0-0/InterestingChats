package models

type BaseUser struct {
	Email    string `json:"email"`
	Username string `json:"username"`
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
