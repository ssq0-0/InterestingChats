package handlers

import (
	"net/http"
	"notifications/internal/models"
)

func HandleError(w http.ResponseWriter, statusCode int, errMsg []string, logMsg error) {
	SendRespond(w, statusCode, &models.Response{
		Errors: errMsg,
		Data:   nil,
	})
	// log.Warn(logMsg)
}
