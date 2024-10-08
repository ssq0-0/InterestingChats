package utils

import (
	"InterestingChats/backend/user_services/internal/consts"
	"InterestingChats/backend/user_services/internal/models"
	"encoding/json"
	"fmt"
	"time"
)

// MessageConstruct creates a notification based on the response data and notification type.
// Returns a notification or error if something went wrong.
func MessageConstruct(response *models.Response, typeNotification int) (*models.Notification, error) {
	if response == nil {
		return nil, fmt.Errorf("response is nil")
	}
	if response.Data == nil {
		return nil, fmt.Errorf("response.Data is nil")
	}

	dataBytes, err := json.Marshal(response.Data)
	if err != nil {
		return nil, err
	}

	var notifRequest *models.NotificationRequest
	if err := json.Unmarshal(dataBytes, &notifRequest); err != nil {
		return nil, err
	}

	var message string
	var typeMessage string
	switch typeNotification {
	case consts.NOTIFICATION_AddFriend:
		message = fmt.Sprintf("User %s sent a friend request :)", notifRequest.Sender.Username)
		typeMessage = "friend_request"
	case consts.NOTIFICATION_DeleteFriend:
		message = fmt.Sprintf("User %s deleted you from friends :(", notifRequest.Sender.Username)
		typeMessage = "friend_deleted"
	case consts.NOTIFICATION_AcceptFriendship:
		message = fmt.Sprintf("User %s accept you friendship request ;)", notifRequest.Sender.Username)
		typeMessage = "friend_accept"
	default:
		return nil, fmt.Errorf("unsupported notification type: %d", typeNotification)
	}

	return &models.Notification{
		UserID:   notifRequest.Receiver,
		SenderID: notifRequest.Sender.ID,
		Message:  message,
		Time:     time.Now(),
		Type:     typeMessage,
		IsRead:   false,
	}, nil
}
