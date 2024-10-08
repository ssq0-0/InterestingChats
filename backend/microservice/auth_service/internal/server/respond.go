package server

import (
	"encoding/json"
	"log"
	"net/http"
)

// SendRespond sends a JSON response to the client with the specified status code and data.
func (s *Server) SendRespond(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	// h.log.Infof()
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("here")
		http.Error(w, "error encoding json", http.StatusInternalServerError)
		s.log.Warn("error encoding json for response", err)
		return
	}
}
