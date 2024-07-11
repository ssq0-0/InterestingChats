package user

import (
	"InterestingChats/backend/microservice/db/internal/db"
	"InterestingChats/backend/microservice/db/internal/models"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
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
	var u models.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println("Can't parse user data from request: ", err)
		return
	}

	userService := db.NewUserService(h.Db)
	err := userService.CreateNewUser(r.Context(), u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println("Can't exec user in database: ", err)
		return
	}
	response := map[string]string{
		"message": "User registered successfully",
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var u models.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		log.Println("error decode: ", err)
		http.Error(w, "Failed to decode user login data", http.StatusBadRequest)
		return
	}

	userService := db.NewUserService(h.Db)
	dbPassword, err := userService.LoginData(r.Context(), u)
	if err != nil {
		log.Println("Incorrect user data decode: ", err)
		http.Error(w, "Incorrect user data:", http.StatusBadRequest)
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(u.Password)); err != nil {
		log.Println("Incoerrect email or password:", err)
		http.Error(w, "Incorrect email or password", http.StatusBadRequest)
		return
	}
	response := map[string]string{
		"message": "Successful login!",
	}
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(response)
}
