package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

// SendRespond sends a JSON response to the client with the specified status code and data.
func SendRespond(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "error encoding json", http.StatusInternalServerError)
		log.Println("error encoding json for response", err)
		return
	}
}
