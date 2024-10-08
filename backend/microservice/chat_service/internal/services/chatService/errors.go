package chatservice

import (
	"chat_service/internal/logger"
	"chat_service/internal/models"
	"net/http"
)

// ErrorHandler handles the error response for the HTTP requests.
// It sends an error response and logs the error message.
func ErrorHandler(w http.ResponseWriter, statusCode int, log logger.Logger, errMsg []string, logMsg string) {
	SendRespond(w, statusCode, log, &models.Response{
		Data:   nil,
		Errors: errMsg,
	})
	log.Error(logMsg)
}
