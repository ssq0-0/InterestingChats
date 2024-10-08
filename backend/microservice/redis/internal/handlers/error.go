package handlers

import (
	"InterestingChats/backend/microservice/redis/internal/models"
	"log"
	"net/http"
)

// HandleError sends a JSON error response to the client and logs the error message.
func HandleError(w http.ResponseWriter, statusCode int, errors []string, logMsg error) {
	SendRespond(w, statusCode, &models.Response{
		Errors: errors,
		Data:   nil,
	})
	log.Println(logMsg)
}
