package proxy

import (
	"InterestingChats/backend/api_gateway/internal/logger"
	"InterestingChats/backend/api_gateway/internal/models"
	"net/http"
)

func ErrorHandler(w http.ResponseWriter, statusCode int, log logger.Logger, errMsg []string, logMsg error) {
	SendRespond(w, statusCode, log, &models.Response{
		Data:   nil,
		Errors: errMsg,
	})
	log.Info(logMsg)
}
