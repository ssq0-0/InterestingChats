package chatservice

import (
	"chat_service/internal/logger"
	"encoding/json"
	"net/http"
)

func SendRespond(w http.ResponseWriter, statusCode int, log logger.Logger, data interface{}) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(statusCode)
	if statusCode != 204 {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, "error encode json", http.StatusInternalServerError)
			log.Errorf("error encoding json: %v", err)
		}
	}
}
