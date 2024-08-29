package db

import (
	"InterestingChats/backend/microservice/db/internal/consts"
	"InterestingChats/backend/microservice/db/internal/db_queries"
	"InterestingChats/backend/microservice/db/internal/models"
	"database/sql"
	"fmt"
	"strings"
)

type UserService struct {
	Db *sql.DB
}

func NewUserService(db *sql.DB) *UserService {
	return &UserService{Db: db}
}

func (us *UserService) CreateNewUser(user models.User) (*models.User, string, error) {
	clientErr, err := db_queries.SelectUserFromDB(user.Email, us.Db)
	if err != nil {
		if err.Error() == fmt.Sprintf("user with email %s already exists", user.Email) {
			return nil, clientErr, err
		}
		return nil, clientErr, fmt.Errorf("could not check if user exists: %v", err)
	}

	userInfo, clientErr, err := db_queries.InsertUser(user, us.Db)
	if err != nil {
		return nil, clientErr, err
	}

	return userInfo, "", nil
}

func (us *UserService) LoginData(user *models.User) (string, string, error) {
	password, clietnErr, err := db_queries.SelectLoginUserInfo(user, us.Db)
	if err != nil {
		return "", clietnErr, err
	}

	return password, "", nil
}

func (us *UserService) CheckUser(userID int) (bool, string, error) {
	if clientErr, err := db_queries.ChechUserUxistsById(userID, us.Db); err != nil {
		return false, clientErr, err
	}

	return true, "", nil
}

func (us *UserService) GetUserInfo(userID int) (*models.BaseUser, string, error) {
	userInfo, err := db_queries.SelectAllUserInfo(userID, us.Db)
	if err != nil {
		return nil, consts.ErrUserNotFound, err
	}

	return userInfo, "", nil
}

func (us *UserService) ChangeUsername(username string, userID int) (string, error) {
	if username == "" || len(strings.ReplaceAll(username, " ", "")) == 0 || userID == 0 {
		return consts.ErrMissingParametr, fmt.Errorf(consts.InternalErrMissingParametr)
	}

	if clientErr, err := db_queries.ChangeUsername(username, userID, us.Db); err != nil {
		return clientErr, err
	}

	return "", nil
}

func (us *UserService) SearchUsers(symbols string) ([]*models.User, string, error) {
	userList, clientErr, err := db_queries.SelectUserBySymbols(symbols, us.Db)
	if err != nil {
		return nil, clientErr, err
	}

	return userList, "", nil
}

func (us *UserService) ChangeEmail(email string, userID int) (string, error) {
	if email == "" || len(strings.ReplaceAll(email, " ", "")) == 0 || userID == 0 {
		return consts.ErrMissingParametr, fmt.Errorf(consts.InternalErrMissingParametr)
	}

	if clientErr, err := db_queries.ChangeEmail(email, userID, us.Db); err != nil {
		return clientErr, err
	}

	return "", nil
}
