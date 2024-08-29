package db_queries

import (
	"InterestingChats/backend/microservice/db/internal/consts"
	"InterestingChats/backend/microservice/db/internal/models"
	"context"
	"database/sql"
	"fmt"
	"log"
)

func SelectUserFromDB(email string, db *sql.DB) (string, error) {
	var existingUserID int
	err := db.QueryRowContext(context.Background(), consts.QUERY_SelectUser, email).Scan(&existingUserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return consts.ErrFailedVerifyUser, fmt.Errorf(consts.InternalErrSelectInfo, err)
	}
	if existingUserID != 0 {
		return consts.ErrUserEmailAlreadyExists, fmt.Errorf(consts.InternalErrUserAlreadyExists, email)
	}

	return "", nil
}

func InsertUser(user models.User, db *sql.DB) (*models.User, string, error) {
	if err := db.QueryRowContext(context.Background(), consts.QUERY_InsertUser, user.Username, user.Email, user.Password).Scan(&user.ID); err != nil {
		return nil, consts.ErrInternalServerError, fmt.Errorf(consts.InternalErrFailedInsert, err)
	}

	return &user, "", nil
}

func SelectLoginUserInfo(user *models.User, db *sql.DB) (string, string, error) {
	var password string
	if err := db.QueryRowContext(context.Background(), consts.QUERY_SelectUserData, user.Email).Scan(&user.ID, &user.Username, &password); err != nil {
		return "", consts.ErrUserNotFound, fmt.Errorf(consts.InternalErrUserNotFoud)
	}

	if len(password) == 0 {
		return "", consts.ErrInvalidRequestFormat, fmt.Errorf(consts.InternalErrPasswordNotFound)
	}

	return password, "", nil
}

func ChechUserUxistsById(userID int, db *sql.DB) (string, error) {
	var exists bool
	if err := db.QueryRowContext(context.Background(), consts.QUERY_SelectUserExists, userID).Scan(&exists); err != nil {
		log.Printf("error transaction: %v", err)
		return consts.ErrInternalServerError, err
	}
	if !exists {
		return consts.ErrUserNotFound, fmt.Errorf(consts.InternalErrUserNotFoud)
	}
	return "", nil
}

func SelectAllUserInfo(userID int, db *sql.DB) (*models.BaseUser, error) {
	var user models.BaseUser
	if err := db.QueryRowContext(context.Background(), consts.QUERY_SelectProfileUserInfo, userID).Scan(&user.Username, &user.Email); err != nil {
		return nil, err
	}

	return &user, nil
}

func ChangeUsername(username string, userID int, db *sql.DB) (string, error) {
	var scanResult string
	if err := db.QueryRowContext(context.Background(), consts.QUERY_SelectUsername, userID).Scan(&scanResult); err != nil {
		return consts.ErrInternalServerError, err
	}

	if scanResult == username {
		return consts.ErrUsernameAlredyExists, fmt.Errorf(consts.InternalErrChangedUserInfo)
	}

	if _, err := db.ExecContext(context.Background(), consts.QUERY_ChangeUsername, username, userID); err != nil {
		return "", err
	}

	return "", nil
}

func SelectUserBySymbols(symbols string, db *sql.DB) ([]*models.User, string, error) {
	// queryString := symbols + ":*"
	rows, err := db.Query(consts.QUERY_SelectIndexUser, symbols)
	if err != nil {
		return nil, consts.ErrInternalServerError, fmt.Errorf(consts.InternalErrFailedRequest, err)
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email); err != nil {
			return nil, consts.ErrInternalServerError, err
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, consts.ErrInternalServerError, err
	}

	return users, "", nil
}

func ChangeEmail(email string, userID int, db *sql.DB) (string, error) {
	var existingUserID int
	if err := db.QueryRowContext(context.Background(), consts.QUERY_CheckSingleEmail, email, userID).Scan(&existingUserID); err != nil {
		if err != sql.ErrNoRows {
			return "", err
		}
	} else {
		return consts.ErrEmailAlredyExists, fmt.Errorf(consts.InternalErrChangedUserInfo)
	}

	var scanResult string
	if err := db.QueryRowContext(context.Background(), consts.QUERY_SelectEmail, userID).Scan(&scanResult); err != nil {
		return consts.ErrMissingParametr, err
	}
	if scanResult == email {
		return consts.ErrEmailAlredyExists, fmt.Errorf(consts.InternalErrChangedUserInfo)
	}

	if _, err := db.ExecContext(context.Background(), consts.QUERY_ChangeEmail, email, userID); err != nil {
		return "", err
	}

	return "", nil
}
