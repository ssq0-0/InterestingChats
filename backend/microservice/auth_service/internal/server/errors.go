package server

import (
	"auth_service/internal/models"
	"net/http"
)

// HandleError sends a JSON error response to the client and logs the error message.
func (s *Server) HandleError(w http.ResponseWriter, statusCode int, errMsg []string, logMsg error) {
	s.SendRespond(w, statusCode, &models.Response{
		Errors: errMsg,
		Data:   nil,
	})
	s.log.Warn(logMsg)
}
