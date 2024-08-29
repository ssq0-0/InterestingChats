package handlers

import (
	"InterestingChats/backend/microservice/db/internal/logger"
	"InterestingChats/backend/microservice/db/internal/models"
	"net/http"
)

func ErrorHandler(w http.ResponseWriter, statusCode int, log logger.Logger, errMsg []string, logMsg error) {
	SendRespond(w, statusCode, &models.Response{
		Data:   nil,
		Errors: errMsg,
	})
	log.Warn(logMsg)
}
