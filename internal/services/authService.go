package services

import (
	"InterestingChats/internal/db"
	"InterestingChats/internal/utils"
	"encoding/json"
	"io/ioutil"
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

func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
	var u db.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	userService := db.NewUserService(s.Db)
	boolCheck, accessToken, refreshToken, err := userService.LoginData(r.Context(), &u)
	if err != nil && !boolCheck {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := map[string]string{
		"message":      "Successful; login!",
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	}
	// w.WriteHeader(http.StatusOK)
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	json.NewEncoder(w).Encode(response)
}

func (s *Server) RefreshAccessToken(w http.ResponseWriter, r *http.Request) {

}

func (s *Server) UpdateAccessToken(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}

	var data map[string]string
	if err := json.Unmarshal(body, &data); err != nil {
		http.Error(w, "Invalid json format", http.StatusBadRequest)
		return
	}

	refreshToken, exists := data["refreshToken"]
	if !exists {
		http.Error(w, "Missing refreshToken in request body", http.StatusBadRequest)
		return
	}

	userService := db.NewUserService(s.Db)
	tokenExists, err := userService.CheckToken(r.Context(), refreshToken)
	if err != nil {
		http.Error(w, "Can't search token", http.StatusBadRequest)
		return
	}

	if !tokenExists {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	newAccessToken, err := utils.RefreshToken(refreshToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"message":     "update token!",
		"accessToken": newAccessToken,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
