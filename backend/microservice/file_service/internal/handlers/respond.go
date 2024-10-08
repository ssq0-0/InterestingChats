package handlers

import (
	"encoding/json"
	"net/http"
)

func (is *ImageService) SendRespond(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "error encoding json", http.StatusInternalServerError)
		is.log.Warn("error encoding json for response", err)
		return
	}
}
