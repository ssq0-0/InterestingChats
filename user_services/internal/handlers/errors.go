package handlers

import (
	"InterestingChats/backend/user_services/internal/models"
	"net/http"
)

func (h *UserService) HandleError(w http.ResponseWriter, statusCode int, errMsg []string, logMsg error) {
	h.SendRespond(w, statusCode, &models.Response{
		Errors: errMsg,
		Data:   nil,
	})
	h.log.Warn(logMsg)
}
