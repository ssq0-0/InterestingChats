package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"notifications/internal/consts"
	"notifications/internal/models"
	"notifications/internal/utils"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// SendNotification processes and sends a notification message from Kafka
func SendNotification(message *kafka.Message) {
	var notification *models.Notification
	if err := json.Unmarshal(message.Value, &notification); err != nil {
		log.Printf("Failed to decode notification message: %v", err)
		return
	}

	response, statusCode, clientErr, err := utils.ProxyRequest(consts.POST_Method, consts.DB_AddNotification, notification, http.StatusOK)
	if err != nil {
		log.Printf("respnse %+v, statusCode: %d, ClientErr: %s, error: %v", response, statusCode, clientErr, err)
		return
	}

	log.Printf("respnse %+v, statusCode: %d, ClientErr: %s, error: %v", response, statusCode, clientErr, err)
}

// GetNotifications retrieves notifications for a user based on their ID
func GetNotifications(r *http.Request) (*models.Response, int, string, error) {
	var userID int
	if clientErr, err := utils.ConvertValue(r.Header.Get("X-User-ID"), &userID); err != nil {
		return nil, http.StatusBadRequest, clientErr, err
	}

	response, statusCode, clientErr, err := utils.ProxyRequest(consts.GET_Method, fmt.Sprintf(consts.DB_GetNotifications, userID), nil, http.StatusOK)
	if err != nil {
		return nil, statusCode, clientErr, err
	}
	return response, http.StatusOK, "", nil
}

// ReadNotification marks a notification as read
func ReadNotification(r *http.Request) (*models.Response, int, string, error) {
	var notification []*models.Notification
	if err := json.NewDecoder(r.Body).Decode(&notification); err != nil {
		log.Printf("Failed to decode notification message: %v", err)
		return nil, http.StatusBadRequest, "clientErr", err
	}

	response, statusCode, clientErr, err := utils.ProxyRequest(consts.PATCH_Method, consts.DB_ReadNotification, notification, http.StatusOK)
	if err != nil {
		return nil, statusCode, clientErr, err
	}

	return response, http.StatusOK, "", nil
}
