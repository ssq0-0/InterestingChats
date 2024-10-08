package handlers

import (
	"encoding/json"
	"net/http"
)

func SendRespond(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	// h.log.Infof()
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "error encoding json", http.StatusInternalServerError)
		return
	}
}
