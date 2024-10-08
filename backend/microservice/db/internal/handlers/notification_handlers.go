package handlers

import (
	"InterestingChats/backend/microservice/db/internal/models"
	"net/http"
)

// AddNotification to add a new notification for a user.
func (h *handler) AddNotification(w http.ResponseWriter, r *http.Request) {
	statusCode, clientErr, err := h.DBService.AddNotification(r)
	if err != nil {
		ErrorHandler(w, statusCode, h.log, []string{clientErr}, err)
		return
	}

	SendRespond(w, http.StatusOK, &models.Response{
		Errors: nil,
		Data:   "notification add",
	})
}

// GetNotification to retrieve a user's notifications.
func (h *handler) GetNotification(w http.ResponseWriter, r *http.Request) {
	notify, statusCode, clientErr, err := h.DBService.GetNotification(r.URL.Query().Get("userID"))
	if err != nil {
		ErrorHandler(w, statusCode, h.log, []string{clientErr}, err)
		return
	}

	SendRespond(w, http.StatusOK, &models.Response{
		Errors: nil,
		Data:   notify,
	})
}

// ReadNotifications to mark notifications as read.
func (h *handler) ReadNotifications(w http.ResponseWriter, r *http.Request) {
	if statusCode, clientErr, err := h.DBService.ReadNotifications(r); err != nil {
		ErrorHandler(w, statusCode, h.log, []string{clientErr}, err)
		return
	}

	SendRespond(w, http.StatusOK, &models.Response{
		Errors: nil,
		Data:   "seccessful read notifications!",
	})
}
