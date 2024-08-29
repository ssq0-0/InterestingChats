package handlers

import (
	"InterestingChats/backend/microservice/redis/internal/models"
	"log"
	"net/http"
)

func HandleError(w http.ResponseWriter, statusCode int, errors []string, logMsg string) {
	SendRespond(w, statusCode, &models.Response{
		Errors: errors,
		Data:   nil,
	})
	log.Println(logMsg)
}
