package handlers

import (
	"InterestingChats/backend/microservice/redis/internal/models"
	"InterestingChats/backend/microservice/redis/internal/rdb"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type UserHandler struct {
	rdb *rdb.RedisClient
}

func NewUserHandler(rdb *rdb.RedisClient) *UserHandler {
	return &UserHandler{rdb: rdb}
}

func (uh *UserHandler) GetUsersTokens(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		log.Println("email not found")
		http.Error(w, "email not found", http.StatusBadRequest)
		return
	}

	tokens, err := uh.rdb.GetTokens(email)
	if err != nil {
		http.Error(w, "user with this email not found", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tokens); err != nil {
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		return
	}
}

func (uh *UserHandler) SetTokens(w http.ResponseWriter, r *http.Request) {
	var userTokens map[string]models.UserTokens
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		log.Println("Failed to read request body:", err)
		return
	}
	log.Println("Raw request body:", string(body))

	if err := json.Unmarshal(body, &userTokens); err != nil {
		log.Println(userTokens)
		log.Println("Failed to parse json data:", err)
		http.Error(w, "Failed to parse json data", http.StatusBadRequest)
		return
	}
	log.Println("Received tokens data:", userTokens)

	if err := uh.rdb.SetToken(userTokens); err != nil {
		log.Println("Failed to set token to redis:", err)
		http.Error(w, "Failed to set token to redis", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("success")
}
