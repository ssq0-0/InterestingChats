package handlers

import (
	"log"
	"net/http"
	"notifications/internal/services"
)

// GetNotifications handles HTTP GET requests for notifications
func GetNotifications(w http.ResponseWriter, r *http.Request) {
	response, statusCode, clientErr, err := services.GetNotifications(r)
	if err != nil {
		HandleError(w, statusCode, []string{clientErr}, err)
		return
	}
	log.Printf("notifications: %v", response)
	SendRespond(w, http.StatusOK, response)
}

// ReadNotification handles HTTP PATCH requests to mark notifications as read
func ReadNotification(w http.ResponseWriter, r *http.Request) {
	response, statusCode, clientErr, err := services.ReadNotification(r)
	if err != nil {
		HandleError(w, statusCode, []string{clientErr}, err)
		return
	}

	SendRespond(w, http.StatusOK, response)
}
