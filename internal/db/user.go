package db

import (
	"InterestingChats/internal/utils"
	"context"
	"database/sql"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserService struct {
	Db *sql.DB
}

func NewUserService(db *sql.DB) *UserService {
	return &UserService{Db: db}
}

func (us *UserService) CreateNewUser(ctx context.Context, user *User) (string, string, error) {
	hashPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", "", err
	}

	accessToken, refreshToken, err := utils.GenerateJWT(user.Username)
	if err != nil {
		return "", "", err
	}

	query := "INSERT INTO users(username, email, password, token) VALUES ($1, $2, $3, $4) RETURNING id"
	err = us.Db.QueryRowContext(ctx, query, user.Username, user.Email, hashPass, refreshToken).Scan(&user.ID)
	if err != nil {
		return "", "", fmt.Errorf("cloud not execute query: %v", err)
	}

	return accessToken, refreshToken, nil
}
