package proxy

import (
	"InterestingChats/backend/api_gateway/internal/logger"
	"encoding/json"
	"net/http"
)

// SendRespond sends a JSON response to the client with the specified status code and data.
func SendRespond(w http.ResponseWriter, statusCode int, log logger.Logger, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if statusCode != 204 {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, "error encode json", http.StatusInternalServerError)
			log.Warn("error encoding json for response", err)
		}
	}
}
