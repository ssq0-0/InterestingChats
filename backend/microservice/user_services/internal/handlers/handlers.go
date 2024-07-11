package handlers

import (
	"InterestingChats/backend/user_services/internal/models"
	"InterestingChats/backend/user_services/internal/utils"
	"bytes"
	"fmt"
	"log"

	"encoding/json"
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
		log.Println("Problems with decode data")
		return
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Problems with generate password", http.StatusBadRequest)
		log.Println("Problems with generate password")
		return
	}

	accessToken, refreshToken, err := utils.GenerateJWT(u.Username)
	if err != nil {
		http.Error(w, "Problems with generate JWT", http.StatusBadRequest)
		log.Println("Problems with generate JWT")
		return
	}

	userTokens := map[string]models.UserTokens{
		u.Email: {
			Tokens: models.Tokens{
				AccessToken:  accessToken,
				RefreshToken: refreshToken,
			},
		},
	}

	targetRedis := "http://localhost:8003/setToken"
	redisBody, statusCodeByRedis, err := sendRequest("POST", targetRedis, userTokens)
	if err != nil || statusCodeByRedis != http.StatusOK {
		http.Error(w, "Failed to store tokens in Redis.", http.StatusInternalServerError)
		log.Println("Error:", err)
		return
	}
	log.Println("Status Code:", statusCodeByRedis)
	log.Println("Response Body:", string(redisBody))

	data := map[string]interface{}{
		"username": u.Username,
		"email":    u.Email,
		"password": hashPassword,
	}
	targetDB := "http://localhost:8002/registration"
	body, statusCode, err := sendRequest("POST", targetDB, data)
	if err != nil {
		http.Error(w, "Failed to create user in database", statusCode)
		log.Println("Failed to create user in database")
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
		log.Println("Failed to parse json data:", err)
		return
	}

	log.Println("Received login data:", u)

	data := map[string]interface{}{
		"email":    u.Email,
		"password": u.Password,
	}

	targetDB := "http://localhost:8002/login"
	body, statusCode, err := sendRequest("POST", targetDB, data)
	if err != nil {
		http.Error(w, "Failed to login.", statusCode)
		log.Println("Failed to login:", err)
		return
	}
	if statusCode != http.StatusAccepted {
		http.Error(w, string(body), statusCode)
		return
	}

	accessToken, refreshToken, err := utils.GenerateJWT(u.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println("Problems with generate JWT:", err)
		return
	}

	userTokens := map[string]models.UserTokens{
		u.Email: {
			Tokens: models.Tokens{
				AccessToken:  accessToken,
				RefreshToken: refreshToken,
			},
		},
	}

	targetRedis := "http://localhost:8003/setToken"
	redisBody, statusCodeByRedis, err := sendRequest("POST", targetRedis, userTokens)
	if err != nil {
		http.Error(w, "Failed to store tokens in Redis.", http.StatusInternalServerError)
		log.Println("Failed to store tokens in Redis:", err)
		return
	}
	if statusCode != http.StatusAccepted {
		http.Error(w, string(body), statusCode)
		log.Println("Login failed with status code:", statusCode)
		return
	}
	response := map[string]string{
		"body":         string(body),
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	}

	log.Println("Status Code:", statusCodeByRedis)
	log.Println("Response Body:", string(redisBody))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}

func (h *Handler) GetTokens(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "email not found", http.StatusBadRequest)
		return
	}

	target := fmt.Sprintf("http://localhost:8003/user?email=%s", email)
	redisResponse, statusCode, err := sendRequest("GET", target, nil)
	if err != nil {
		http.Error(w, "Failed get tokens.", statusCode)
		log.Println("Failed to get tokens:", err)
		return
	}

	if statusCode != http.StatusOK {
		http.Error(w, "Failed to get tokens.", statusCode)
		log.Println("Unexpected status code:", statusCode)
		return
	}

	var tokens models.Tokens
	err = json.Unmarshal(redisResponse, &tokens)
	if err != nil {
		http.Error(w, "Failed to decode tokens.", http.StatusInternalServerError)
		log.Println("Failed to decode tokens.")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tokens)
}

func sendRequest(method, url string, data interface{}) ([]byte, int, error) {
	var jsonReqData []byte
	var err error
	if data != nil {
		jsonReqData, err = json.Marshal(data)
		if err != nil {
			log.Println("Failed to serialize data.")
			return nil, http.StatusInternalServerError, fmt.Errorf("failed to serialize data: %w", err)
		}
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonReqData))
	if err != nil {
		log.Println("Failed to create request.")
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Failed to send request.")
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Failed to read response.")
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to read response: %w", err)
	}

	return body, resp.StatusCode, nil
}
