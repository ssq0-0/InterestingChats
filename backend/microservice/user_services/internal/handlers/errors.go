package handlers

import (
	"InterestingChats/backend/user_services/internal/models"
	"log"
	"net/http"
)

func (h *Handler) HandleError(w http.ResponseWriter, statusCode int, errMsg []string, logMsg string) {
	h.SendRespond(w, statusCode, &models.Response{
		Errors: errMsg,
		Data:   nil,
	})
	log.Println(logMsg)
}
