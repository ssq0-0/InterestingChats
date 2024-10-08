package handlers

import (
	"file_service/internal/models"
	"net/http"
)

func (is *ImageService) HandleError(w http.ResponseWriter, statusCode int, errMsg []string, logMsg error) {
	is.SendRespond(w, statusCode, &models.Response{
		Errors: errMsg,
		Data:   nil,
	})
	is.log.Warn(logMsg)
}
