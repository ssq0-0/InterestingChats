package user

import (
	"InterestingChats/backend/microservice/db/internal/db"
	"InterestingChats/backend/microservice/db/internal/models"
	"database/sql"
	"encoding/json"
	"fmt"
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
	var u models.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println("Can't parse user data from request")
		return
	}

	userService := db.NewUserService(h.Db)
	err := userService.CreateNewUser(r.Context(), u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println("Can't exec user in database")
		return
	}
	response := map[string]string{
		"message": "User registered successfully",
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
