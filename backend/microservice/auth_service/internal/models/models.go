package models

import "github.com/dgrijalva/jwt-go"

// BaseUser represents core user attributes
type BaseUser struct {
	Email    string  `json:"email"`
	Username string  `json:"username"`
	Avatar   *string `json:"avatar"`
}

// User extends BaseUser by adding an ID and Password
type User struct {
	ID int `json:"id"`
	BaseUser
	Password string `json:"password"`
}

// Tokens holds access and refresh tokens for the user
type Tokens struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// UserTokens is a wrapper for user-related tokens
type UserTokens struct {
	Tokens Tokens `json:"tokens"`
}

// Response is used for API responses with errors and data
type Response struct {
	Errors []string    `json:"Errors"`
	Data   interface{} `json:"Data"`
}

// Claims defines the payload of a JWT token, extending the standard claims
type Claims struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	jwt.StandardClaims
}
