package db_queries

import (
	"InterestingChats/backend/microservice/db/internal/consts"
	"InterestingChats/backend/microservice/db/internal/models"
	"context"
	"database/sql"
	"log"

	"github.com/lib/pq"
)

func AddNotification(userID, senderID int, notifType, message string, tx *sql.Tx) error {
	if _, err := tx.ExecContext(context.Background(), consts.QUERY_AddNotification, userID, senderID, notifType, message); err != nil {
		log.Printf("Error adding notification: %v", err)
		return err
	}
	return nil
}

func GetNotifications(userID int, tx *sql.DB) ([]*models.Notification, error) {
	rows, err := tx.QueryContext(context.Background(), consts.QUERY_GetUserNotifications, userID)
	if err != nil {
		log.Printf("Error retrieving notifications: %v", err)
		return nil, err
	}
	defer rows.Close()

	var notifications []*models.Notification
	for rows.Next() {
		notif := &models.Notification{}
		if err := rows.Scan(&notif.ID, &notif.UserID, &notif.SenderID, &notif.Type, &notif.Message, &notif.Time, &notif.IsRead); err != nil {
			log.Printf("Error scanning notification: %v", err)
			return nil, err
		}
		notifications = append(notifications, notif)
	}
	log.Printf("notify: %+v", notifications)
	return notifications, nil
}

func MarkNotificationAsRead(notifications []*models.Notification, tx *sql.Tx) (string, error) {
	if len(notifications) == 0 {
		return "", nil
	}

	var ids []int
	for _, notifi := range notifications {
		ids = append(ids, notifi.ID)
	}
	log.Printf("Attempting to mark notifications as read for IDs: %v", ids)

	if _, err := tx.ExecContext(context.Background(), consts.QUERY_MarkNotificationIsRead, pq.Array(ids)); err != nil {
		log.Printf("Error marking notifications as read: %v", err)
		return consts.ErrFailedInsertNotification, err
	}
	return "", nil
}
