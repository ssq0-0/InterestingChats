package handlers

import (
	"InterestingChats/backend/user_services/internal/models"
	"InterestingChats/backend/user_services/internal/utils"
	"bytes"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	// Db *sql.DB
}

func NewHandler() *Handler {
	return &Handler{
		// Db: db,
	}
}

func (h *Handler) Registrations(w http.ResponseWriter, r *http.Request) {
	var u models.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println("Problems with generate password")
		return
	}

	accessToken, refreshToken, err := utils.GenerateJWT(u.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println("Problems with generate JWT")
		return
	}

	data := map[string]interface{}{
		"username": u.Username,
		"email":    u.Email,
		"password": hashPassword,
	}
	body, statusCode, err := sendRequestToDBServer(data)
	if err != nil {
		http.Error(w, "Failed to create user in database", statusCode)
		return
	}

	response := map[string]string{
		"message":      "Created new user!",
		"username":     u.Username,
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
		"body":         string(body),
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

func sendRequestToDBServer(data map[string]interface{}) ([]byte, int, error) {
	jsonReqData, err := json.Marshal(data)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to serialize data: %w", err)
	}

	resp, err := http.Post("http://localhost:8002/registration", "application/json", bytes.NewBuffer(jsonReqData))
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to read response: %w", err)
	}

	return body, resp.StatusCode, nil
}

// func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
// 	var u models.User
// 	fmt.Println("Login...")
// 	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 	}
// 	fmt.Println("Login...")

// 	userService := db.NewUserService(h.Db)
// 	boolCheck, accessToken, refreshToken, err := userService.LoginData(r.Context(), u)
// 	if err != nil && !boolCheck {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}
// 	fmt.Println("Login...")

// 	response := map[string]string{
// 		"message":      "Successful; login!",
// 		"accessToken":  accessToken,
// 		"refreshToken": refreshToken,
// 	}
// 	fmt.Println("Login...")

// 	// w.WriteHeader(http.StatusOK)
// 	// http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
// 	w.WriteHeader(http.StatusCreated)
// 	json.NewEncoder(w).Encode(response)
// 	fmt.Println("Login...")

// }

// func (h *Handler) RefreshAccessToken(w http.ResponseWriter, r *http.Request) {

// }

// func (h *Handler) UpdateAccessToken(w http.ResponseWriter, r *http.Request) {
// 	body, err := ioutil.ReadAll(r.Body)
// 	if err != nil {
// 		http.Error(w, "Unable to read request body", http.StatusBadRequest)
// 		return
// 	}

// 	var data map[string]string
// 	if err := json.Unmarshal(body, &data); err != nil {
// 		http.Error(w, "Invalid json format", http.StatusBadRequest)
// 		return
// 	}

// 	refreshToken, exists := data["refreshToken"]
// 	if !exists {
// 		http.Error(w, "Missing refreshToken in request body", http.StatusBadRequest)
// 		return
// 	}

// 	userService := db.NewUserService(h.Db)
// 	tokenExists, err := userService.CheckToken(r.Context(), refreshToken)
// 	if err != nil {
// 		http.Error(w, "Can't search token", http.StatusBadRequest)
// 		return
// 	}

// 	if !tokenExists {
// 		http.Error(w, "Invalid token", http.StatusUnauthorized)
// 		return
// 	}

// 	newAccessToken, err := utils.RefreshToken(refreshToken)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	response := map[string]string{
// 		"message":     "update token!",
// 		"accessToken": newAccessToken,
// 	}

// 	w.WriteHeader(http.StatusCreated)
// 	json.NewEncoder(w).Encode(response)
// }
