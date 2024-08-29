package proxy

import (
	"InterestingChats/backend/api_gateway/internal/logger"
	"encoding/json"
	"net/http"
)

func SendRespond(w http.ResponseWriter, statusCode int, log logger.Logger, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	// w.Header().Set("Access-Control-Allow-Origin", "*") // Разрешите ваш домен здесь, если нужно
	// w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	// w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
	w.WriteHeader(statusCode)

	if statusCode != 204 {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, "error encode json", http.StatusInternalServerError)
			log.Warn("error encoding json for response", err)
		}
	}
}
