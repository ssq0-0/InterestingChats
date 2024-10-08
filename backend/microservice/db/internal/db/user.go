package db

import (
	"InterestingChats/backend/microservice/db/internal/consts"
	"InterestingChats/backend/microservice/db/internal/db_queries"
	"InterestingChats/backend/microservice/db/internal/models"
	"InterestingChats/backend/microservice/db/internal/utils"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// CreateNewUser takes a User model as input and attempts to create a new user in the database.
func (db *dbService) CreateNewUser(user models.User) (*models.User, string, error) {
	clientErr, err := db_queries.SelectUserFromDB(user.Email, db.db)
	if err != nil {
		if err.Error() == fmt.Sprintf("user with email %s already exists", user.Email) {
			return nil, clientErr, err
		}
		return nil, clientErr, fmt.Errorf("could not check if user exists: %v", err)
	}

	userInfo, clientErr, err := db_queries.InsertUser(user, db.db)
	if err != nil {
		return nil, clientErr, err
	}

	return userInfo, "", nil
}

// LoginData retrieves the login credentials of a User from the database. It takes a pointer to
func (db *dbService) LoginData(user *models.User) (string, string, error) {
	password, clietnErr, err := db_queries.SelectLoginUserInfo(user, db.db)
	if err != nil {
		return "", clietnErr, err
	}

	return password, "", nil
}

// CheckUser verifies if a user with the specified userID exists in the database.
func (db *dbService) CheckUser(userID int) (bool, string, error) {
	if clientErr, err := db_queries.ChechUserUxistsById(userID, db.db); err != nil {
		return false, clientErr, err
	}

	return true, "", nil
}

// GetUserInfo fetches detailed information about a user specified by userID.
func (db *dbService) GetUserInfo(userID int) (*models.User, string, error) {
	userInfo, err := db_queries.SelectAllUserInfo(userID, db.db)
	if err != nil {
		return nil, consts.ErrUserNotFound, err
	}

	return userInfo, "", nil
}

// ChangeUsername updates the username for the user identified by userID.
func (db *dbService) ChangeUsername(username string, userID int) (string, error) {
	if username == "" || len(strings.ReplaceAll(username, " ", "")) == 0 || userID == 0 {
		return consts.ErrMissingParametr, fmt.Errorf(consts.InternalErrMissingParametr)
	}

	if clientErr, err := db_queries.ChangeUsername(username, userID, db.db); err != nil {
		return clientErr, err
	}

	return "", nil
}

// SearchUsers searches for users based on a specified string of symbols.
func (db *dbService) SearchUsers(symbols string) ([]*models.User, string, error) {
	userList, clientErr, err := db_queries.SelectUserBySymbols(symbols, db.db)
	if err != nil {
		return nil, clientErr, err
	}

	return userList, "", nil
}

// ChangeEmail updates the email address for the user identified by userID.
func (db *dbService) ChangeEmail(email string, userID int) (string, error) {
	if email == "" || len(strings.ReplaceAll(email, " ", "")) == 0 || userID == 0 {
		return consts.ErrMissingParametr, fmt.Errorf(consts.InternalErrMissingParametr)
	}

	if clientErr, err := db_queries.ChangeEmail(email, userID, db.db); err != nil {
		return clientErr, err
	}

	return "", nil
}

// UploadAvatar uploads an avatar image for a user.
func (db *dbService) UploadAvatar(req *models.UserFile) (int, string, error) {
	log.Printf("db: %+v", req)
	if strings.ReplaceAll(req.URL.StaticLink, " ", "") == "" {
		return http.StatusBadRequest, consts.ErrMissingParametr, fmt.Errorf(consts.InternalErrMissingParametr)
	}

	tx, err := db.db.BeginTx(context.Background(), nil)
	if err != nil {
		return http.StatusBadRequest, consts.ErrInternalServerError, err
	}
	defer utils.RollbackOnError(tx, err)

	// TODO: change temporary url to static.
	if clientErr, err := db_queries.InsertAvatar(req.URL.TemporaryLink, req.UserID, tx); err != nil {
		return http.StatusBadRequest, clientErr, err
	}

	if err := tx.Commit(); err != nil {
		return http.StatusBadRequest, consts.ErrInternalServerError, err
	}

	return http.StatusOK, "", nil
}

// FriendshipOperation handles friendship requests and operations based on the request type.
func (db *dbService) FriendshipOperation(r *http.Request, requestType int) (interface{}, int, string, error) {
	var friendRequest models.FriendRequest
	if r.Method != http.MethodGet {
		if err := json.NewDecoder(r.Body).Decode(&friendRequest); err != nil {
			return nil, 400, consts.ErrInternalServerError, err
		}
	} else {
		var userID int
		if clientErr, err := utils.ConvertValue(r.URL.Query().Get("userID"), &userID); err != nil {
			return nil, 400, clientErr, err
		}
		friendRequest = models.FriendRequest{
			UserID: userID,
		}
	}

	tx, err := db.db.BeginTx(context.Background(), nil)
	if err != nil {
		return nil, http.StatusBadRequest, consts.ErrInternalServerError, err
	}
	defer utils.RollbackOnError(tx, err)

	var userValue interface{}
	var clientErr string
	switch requestType {
	case consts.FRIENDSHIP_RequestType:
		userValue, clientErr, err = db_queries.InsertFriendshipRequest(friendRequest.UserID, friendRequest.FriendID, tx)
		if err != nil {
			return nil, http.StatusBadRequest, clientErr, err
		}
	case consts.FRIENDSHIP_AcceptType:
		userValue, clientErr, err = db_queries.AcceptFriendshipRequest(friendRequest.UserID, friendRequest.FriendID, tx)
		if err != nil {
			return nil, http.StatusBadRequest, clientErr, err
		}
	case consts.FRIENDSHIP_DeleteFriend:
		userValue, clientErr, err = db_queries.DeleteFriend(friendRequest.UserID, friendRequest.FriendID, tx)
		if err != nil {
			return nil, http.StatusBadRequest, clientErr, err
		}
	case consts.FRIENDSHIP_DeleteFriendRequest:
		clientErr, err = db_queries.DeleteFriendRequest(friendRequest.UserID, friendRequest.FriendID, tx)
		if err != nil {
			return nil, http.StatusBadRequest, clientErr, err
		}
	case consts.FRIENDSHIP_GetUsersFriends:
		userValue, clientErr, err = db_queries.GetFriendList(friendRequest.UserID, db.db)
		if err != nil {
			return nil, http.StatusBadRequest, clientErr, err
		}
	case consts.FRIENDSHIP_GetSubs:
		userValue, clientErr, err = db_queries.GetSubList(friendRequest.UserID, db.db)
		if err != nil {
			return nil, http.StatusBadRequest, clientErr, err
		}
	default:
		return nil, http.StatusBadRequest, consts.ErrInvalidRequestFormat, fmt.Errorf("unsupported request type: %d", requestType)
	}

	if err := tx.Commit(); err != nil {
		return nil, http.StatusBadRequest, consts.ErrInternalServerError, err
	}

	return userValue, http.StatusOK, "", nil
}
