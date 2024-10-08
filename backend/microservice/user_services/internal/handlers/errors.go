package handlers

import (
	"InterestingChats/backend/user_services/internal/models"
	"net/http"
)

// HandleError sends a JSON error response to the client and logs the error message.
func (h *Handler) HandleError(w http.ResponseWriter, statusCode int, errMsg []string, logMsg error) {
	h.SendRespond(w, statusCode, &models.Response{
		Errors: errMsg,
		Data:   nil,
	})
	h.log.Warn(logMsg)
}
