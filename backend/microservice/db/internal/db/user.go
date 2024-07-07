package db

import (
	"InterestingChats/backend/microservice/db/internal/models"
	"context"
	"database/sql"
	"fmt"
	"time"
	// "golang.org/x/crypto/bcrypt"
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
	err = us.Db.QueryRowContext(ctx, insertQuery, user.Username, user.Email, user.HashPasword).Scan(&user.ID)
	if err != nil {
		return fmt.Errorf("cloud not execute query: %v", err)
	}

	return nil
}

// func (us *UserService) LoginData(ctx context.Context, user *User) (bool, string, string, error) {
// 	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
// 	defer cancel()

// 	query := "SELECT password FROM users WHERE email = $1"
// 	var storedPassword string
// 	err := us.Db.QueryRowContext(ctx, query, user.Email).Scan(&storedPassword)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return false, "", "", nil
// 		}
// 		return false, "", "", fmt.Errorf("incorrect email or password")
// 	}

// 	if err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(user.Password)); err != nil {
// 		return false, "", "", fmt.Errorf("incorrect email or password")
// 	}

// 	accessToken, refreshToken, err := utils.GenerateJWT(user.Username)
// 	if err != nil {
// 		return false, "", "", err
// 	}

// 	updateTokenQuery := "UPDATE users SET token = $1 WHERE email = $2"
// 	_, err = us.Db.ExecContext(ctx, updateTokenQuery, refreshToken, user.Email)
// 	if err != nil {
// 		return false, "", "", fmt.Errorf("incorrect email or password")
// 	}

// 	return false, accessToken, refreshToken, nil
// }

// func (us *UserService) CheckToken(ctx context.Context, token string) (bool, error) {
// 	ctx, cancel := context.WithTimeout(ctx, 4*time.Second)
// 	defer cancel()

// 	query := "SELECT token FROM users WHERE token = $1"
// 	var rtoken string
// 	err := us.Db.QueryRowContext(ctx, query, token).Scan(&rtoken)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return false, nil
// 		}
// 		return false, fmt.Errorf("could not execute query: %v", err)
// 	}

// 	return true, nil
// }
