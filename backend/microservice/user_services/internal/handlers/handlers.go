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

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Registrations(w http.ResponseWriter, r *http.Request) {
	var u models.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Problems with decode data", http.StatusBadRequest)
		fmt.Println("Problems with decode data")
		return
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Problems with generate password", http.StatusBadRequest)
		fmt.Println("Problems with generate password")
		return
	}
	fmt.Println(string(hashPassword))

	accessToken, refreshToken, err := utils.GenerateJWT(u.Username)
	if err != nil {
		http.Error(w, "Problems with generate JWT", http.StatusBadRequest)
		fmt.Println("Problems with generate JWT")
		return
	}

	data := map[string]interface{}{
		"username": u.Username,
		"email":    u.Email,
		"password": hashPassword,
	}
	body, statusCode, err := sendRequestToDBServer(data, "registration")
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

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var u models.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Failed to parse json data", http.StatusBadRequest)
		return
	}

	data := map[string]interface{}{
		"email":    u.Email,
		"password": u.Password,
	}
	body, statusCode, err := sendRequestToDBServer(data, "login")
	if err != nil {
		http.Error(w, "Failed to login.", statusCode)
		w.WriteHeader(statusCode)
		return
	}
	if statusCode != http.StatusAccepted {
		http.Error(w, string(body), statusCode)
		return
	}

	accessToken, refreshToken, err := utils.GenerateJWT(u.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println("Problems with generate JWT")
		return
	}

	response := map[string]string{
		"body":         string(body),
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}

func sendRequestToDBServer(data map[string]interface{}, target string) ([]byte, int, error) {
	jsonReqData, err := json.Marshal(data)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to serialize data: %w", err)
	}
	urlPOST := fmt.Sprintf("http://localhost:8002/%s", target)
	resp, err := http.Post(urlPOST, "application/json", bytes.NewBuffer(jsonReqData))
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
