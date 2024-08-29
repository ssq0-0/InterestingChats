package user

import (
	"InterestingChats/backend/microservice/db/internal/db"
	"InterestingChats/backend/microservice/db/internal/handlers"
	"InterestingChats/backend/microservice/db/internal/models"
	"InterestingChats/backend/microservice/db/internal/utils"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

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
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		handlers.ErrorHandler(w, http.StatusBadRequest, []string{"Can't parse user data from request"}, fmt.Sprintf("Can't parse user data from request: %v", err))
		return
	}

	userService := db.NewUserService(h.Db)
	if userInfo, err := userService.CreateNewUser(r.Context(), user); err != nil {
		handlers.ErrorHandler(w, http.StatusBadRequest, []string{"Can't exec user in database"}, fmt.Sprintf("Can't exec user in database: %v", err))
		return

	} else {
		handlers.SendRespond(w, http.StatusCreated, &models.Response{
			Errors: nil,
			Data:   userInfo,
		})
	}
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		handlers.ErrorHandler(w, http.StatusBadRequest, []string{"Failed to decode user login data"}, fmt.Sprintf("Error decoding user login data: %v", err))
		return
	}

	userService := db.NewUserService(h.Db)
	dbUser, err := userService.LoginData(r.Context(), user)
	if err != nil {
		handlers.ErrorHandler(w, http.StatusBadRequest, []string{"Error fetching user data from database"}, fmt.Sprintf("Error fetching user data from database: %v", err))
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password)); err != nil {
		handlers.ErrorHandler(w, http.StatusBadRequest, []string{"Incorrect email or password"}, fmt.Sprintf("Incorrect email or password: %v", err))
		return
	}
	// interim solution, rewrite
	dbUser.Password = ""

	handlers.SendRespond(w, http.StatusOK, &models.Response{
		Errors: nil,
		Data:   dbUser,
	})
}

func (h *UserHandler) CheckUser(w http.ResponseWriter, r *http.Request) {
	userIDstr := r.URL.Query().Get("userID")
	if userIDstr == "" {
		handlers.ErrorHandler(w, http.StatusBadRequest, []string{"missing user id"}, "missing user id in request")
		return
	}

	var userID int
	if err := utils.ConvertValue(userIDstr, &userID); err != nil {
		handlers.ErrorHandler(w, http.StatusBadRequest, []string{"error in user service"}, fmt.Sprintf("error in user service: %v", err))
		return
	}

	// Change to flip to Handler
	userService := db.NewUserService(h.Db)
	accept, err := userService.CheckUser(userID)
	if err != nil {
		handlers.ErrorHandler(w, http.StatusBadRequest, []string{"failed request"}, fmt.Sprintf("failed request: %v", err))
		return
	}
	if !accept {
		handlers.ErrorHandler(w, http.StatusBadRequest, []string{"user not found"}, "user not found")
		return
	}

	handlers.SendRespond(w, http.StatusOK, &models.Response{
		Errors: nil,
		Data:   userID,
	})
}
