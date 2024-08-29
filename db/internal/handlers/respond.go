package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

func SendRespond(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if statusCode != 204 {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, "error encode json", http.StatusInternalServerError)
			log.Println("error encoding json for response", err)
		}
	}
}
