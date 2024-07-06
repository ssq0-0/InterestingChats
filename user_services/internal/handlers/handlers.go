package handlers

import (
	"InterestingChats/backend/user_services/internal/db"
	"InterestingChats/backend/user_services/internal/utils"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Handler struct {
	Db *sql.DB
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{
		Db: db,
	}
}

func (h *Handler) Registrations(w http.ResponseWriter, r *http.Request) {
	var u db.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userService := db.NewUserService(h.Db)
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

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var u db.User
	fmt.Println("Login...")
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	fmt.Println("Login...")

	userService := db.NewUserService(h.Db)
	boolCheck, accessToken, refreshToken, err := userService.LoginData(r.Context(), &u)
	if err != nil && !boolCheck {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("Login...")

	response := map[string]string{
		"message":      "Successful; login!",
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	}
	fmt.Println("Login...")

	// w.WriteHeader(http.StatusOK)
	// http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
	fmt.Println("Login...")

}

func (h *Handler) RefreshAccessToken(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) UpdateAccessToken(w http.ResponseWriter, r *http.Request) {
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

	userService := db.NewUserService(h.Db)
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
