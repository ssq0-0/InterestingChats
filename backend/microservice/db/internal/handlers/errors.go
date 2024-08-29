package handlers

import (
	"InterestingChats/backend/microservice/db/internal/models"
	"log"
	"net/http"
)

func ErrorHandler(w http.ResponseWriter, statusCode int, errMsg []string, logMsg string) {
	SendRespond(w, statusCode, &models.Response{
		Data:   nil,
		Errors: errMsg,
	})
	log.Println(logMsg)
}
