package db

import (
	"InterestingChats/backend/microservice/db/internal/models"
	"context"
	"database/sql"
	"fmt"
	"time"
)

type UserService struct {
	Db *sql.DB
}

func NewUserService(db *sql.DB) *UserService {
	return &UserService{Db: db}
}

func (us *UserService) CreateNewUser(ctx context.Context, user models.User) error {
	ctx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()

	selectQuery := "SELECT id FROM users WHERE email = $1"
	var existingUserID int
	err := us.Db.QueryRowContext(ctx, selectQuery, user.Email).Scan(&existingUserID)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("could not check if user exists: %v", err)
	}
	if existingUserID != 0 {
		return fmt.Errorf("user with email %s already exists", user.Email)
	}

	insertQuery := "INSERT INTO users(username, email, password) VALUES ($1, $2, $3) RETURNING id"
	err = us.Db.QueryRowContext(ctx, insertQuery, user.Username, user.Email, user.Pasword).Scan(&user.ID)
	if err != nil {
		return fmt.Errorf("cloud not execute query: %v", err)
	}

	return nil
}

func (us *UserService) LoginData(ctx context.Context, user models.User) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()

	selectQuery := "SELECT password FROM users WHERE email = $1"
	var DBPassword string
	err := us.Db.QueryRowContext(ctx, selectQuery, user.Email).Scan(&DBPassword)
	if err != nil {
		return "", fmt.Errorf("user with email %s not found: %v", user.Email, err)
	}

	return DBPassword, nil
}
