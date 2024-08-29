package models

type BaseUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	BaseUser
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
