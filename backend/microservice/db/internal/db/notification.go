package db

import (
	"InterestingChats/backend/microservice/db/internal/consts"
	"InterestingChats/backend/microservice/db/internal/db_queries"
	"InterestingChats/backend/microservice/db/internal/models"
	"InterestingChats/backend/microservice/db/internal/utils"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// ReadNotifications marks notifications as read for a user based on the incoming HTTP request.
func (db *dbService) AddNotification(r *http.Request) (int, string, error) {
	tx, err := db.db.BeginTx(context.Background(), nil)
	if err != nil {
		return http.StatusBadRequest, consts.ErrInternalServerError, err
	}

	defer utils.RollbackOnError(tx, err)

	var notif *models.Notification
	if err := json.NewDecoder(r.Body).Decode(&notif); err != nil {
		return http.StatusBadRequest, consts.ErrInvalidValueFormat, err
	}

	if err := db_queries.AddNotification(notif.UserID, notif.SenderID, notif.Type, notif.Message, tx); err != nil {
		return http.StatusBadRequest, consts.ErrInternalServerError, err
	}

	if err := tx.Commit(); err != nil {
		return http.StatusBadRequest, consts.ErrInternalServerError, err
	}

	return http.StatusOK, "", nil
}

// GetNotification retrieves notifications for a user based on their userID.
func (db *dbService) GetNotification(userID string) ([]*models.Notification, int, string, error) {
	if strings.ReplaceAll(userID, " ", "") == "" {
		return nil, 400, consts.ErrMissingParametr, fmt.Errorf(consts.InternalErrMissingParametr)
	}

	var userIDint int
	if clientErr, err := utils.ConvertValue(userID, &userIDint); err != nil {
		return nil, 400, clientErr, err
	}

	notify, err := db_queries.GetNotifications(userIDint, db.db)
	if err != nil {
		return nil, 400, consts.ErrInternalServerError, err
	}

	return notify, 200, "", nil
}

// AddNotification creates a new notification based on the incoming HTTP request.
func (db *dbService) ReadNotifications(r *http.Request) (int, string, error) {
	tx, err := db.db.BeginTx(context.Background(), nil)
	if err != nil {
		return http.StatusBadRequest, consts.ErrInternalServerError, err
	}
	defer utils.RollbackOnError(tx, err)

	var notifications []*models.Notification
	if err := json.NewDecoder(r.Body).Decode(&notifications); err != nil {
		return http.StatusBadRequest, consts.ErrInternalServerError, err
	}

	if clientErr, err := db_queries.MarkNotificationAsRead(notifications, tx); err != nil {
		return http.StatusBadRequest, clientErr, err
	}

	if err := tx.Commit(); err != nil {
		return http.StatusBadRequest, consts.ErrInternalServerError, err
	}

	return http.StatusOK, "", nil
}
