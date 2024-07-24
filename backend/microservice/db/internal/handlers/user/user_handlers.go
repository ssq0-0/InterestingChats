package user

import (
	"InterestingChats/backend/microservice/db/internal/db"
	"InterestingChats/backend/microservice/db/internal/models"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	Db *sql.DB
}

func NewHandler(db *sql.DB) *UserHandler {
	return &UserHandler{
		Db: db,
	}
}

func (h *UserHandler) Registrations(w http.ResponseWriter, r *http.Request) {
	var u models.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println("Can't parse user data from request: ", err)
		return
	}

	userService := db.NewUserService(h.Db)
	err := userService.CreateNewUser(r.Context(), u)
	if err != nil {
		/*
			Figure out how to return not only StatusConfilct, but error by situation
		*/
		http.Error(w, err.Error(), http.StatusConflict)
		log.Println("Can't exec user in database: ", err)
		return
	}
	response := map[string]string{
		"message": "User registered successfully",
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var u models.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		log.Println("Error decoding user login data:", err)
		http.Error(w, "Failed to decode user login data", http.StatusBadRequest)
		return
	}

	userService := db.NewUserService(h.Db)
	dbPassword, err := userService.LoginData(r.Context(), u)
	if err != nil {
		log.Println("Error fetching user data from database:", err)
		http.Error(w, "Incorrect user data", http.StatusBadRequest)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(u.Password)); err != nil {
		log.Println("Incorrect email or password:", err)
		http.Error(w, "Incorrect email or password", http.StatusBadRequest)
		return
	}

	response := map[string]string{
		"message": "Successful login!",
	}
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(response)
}

func (h *UserHandler) CheckUser(w http.ResponseWriter, r *http.Request) {
	userIDstr := r.URL.Query().Get("userID")
	if userIDstr == "" {
		log.Println("missing user id")
		http.Error(w, "missing user id n request", http.StatusBadRequest)
		return
	}

	userIDint, err := strconv.Atoi(userIDstr)
	if err != nil {
		log.Printf("error in user service: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// срочно переделать на прокидывание сервиса в хендлер
	userService := db.NewUserService(h.Db)
	accept, err := userService.CheckUser(userIDint)
	if err != nil {
		log.Printf("failed request: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !accept {
		log.Println("user not found")
		http.Error(w, "user not found", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
