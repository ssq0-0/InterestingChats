package services

import (
	"InterestingChats/internal/db"
	"encoding/json"
	"net/http"
)

func (s *Server) Registrations(w http.ResponseWriter, r *http.Request) {
	var u db.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userService := db.NewUserService(s.Db)
	accessToken, refreshToken, err := userService.CreateNewUser(r.Context(), &u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response := map[string]string{
		"message":      "Created new user!",
		"username":     u.Username,
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
