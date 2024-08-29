package user

import (
	"InterestingChats/backend/microservice/db/internal/consts"
	"InterestingChats/backend/microservice/db/internal/db"
	"InterestingChats/backend/microservice/db/internal/handlers"
	"InterestingChats/backend/microservice/db/internal/logger"
	"InterestingChats/backend/microservice/db/internal/models"
	"InterestingChats/backend/microservice/db/internal/utils"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type UserHandler struct {
	Db  *sql.DB
	log logger.Logger
}

func NewHandler(db *sql.DB, log logger.Logger) *UserHandler {
	return &UserHandler{
		Db:  db,
		log: log,
	}
}

func (h *UserHandler) Registrations(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		handlers.ErrorHandler(w, http.StatusBadRequest, h.log, []string{consts.ErrInvalidValueFormat}, err)
		return
	}

	userService := db.NewUserService(h.Db)
	if userInfo, clientErr, err := userService.CreateNewUser(user); err != nil {
		handlers.ErrorHandler(w, http.StatusBadRequest, h.log, []string{clientErr}, err)
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
		handlers.ErrorHandler(w, http.StatusBadRequest, h.log, []string{consts.ErrInvalidValueFormat}, err)
		return
	}

	userService := db.NewUserService(h.Db)
	userPassword, clientErr, err := userService.LoginData(&user)
	if err != nil {
		handlers.ErrorHandler(w, http.StatusBadRequest, h.log, []string{clientErr}, err)
		return
	}

	if err := utils.CompareHashAndPassword(userPassword, user.Password); err != nil {
		handlers.ErrorHandler(w, http.StatusBadRequest, h.log, []string{consts.ErrIncorrectEmailOrPassword}, err)
		return
	}
	// interim solution, rewrite
	user.Password = ""

	handlers.SendRespond(w, http.StatusOK, &models.Response{
		Errors: nil,
		Data:   user,
	})
}

func (h *UserHandler) CheckUser(w http.ResponseWriter, r *http.Request) {
	userIDstr := r.URL.Query().Get("userID")
	if userIDstr == "" {
		handlers.ErrorHandler(w, http.StatusBadRequest, h.log, []string{consts.ErrMissingUserID}, fmt.Errorf(consts.InternalErrMissingURLVal))
		return
	}

	var userID int
	if clientErr, err := utils.ConvertValue(userIDstr, &userID); err != nil {
		handlers.ErrorHandler(w, http.StatusBadRequest, h.log, []string{clientErr}, err)
		return
	}

	// Change to flip to Handler
	userService := db.NewUserService(h.Db)
	accept, clientErr, err := userService.CheckUser(userID)
	if err != nil {
		handlers.ErrorHandler(w, http.StatusBadRequest, h.log, []string{clientErr}, err)
		return
	}
	if !accept {
		handlers.ErrorHandler(w, http.StatusBadRequest, h.log, []string{consts.ErrUserNotFound}, fmt.Errorf(consts.InternalErrUserNotFoud))
		return
	}

	handlers.SendRespond(w, http.StatusOK, &models.Response{
		Errors: nil,
		Data:   userID,
	})
}

func (h *UserHandler) GetUserProfileInfo(w http.ResponseWriter, r *http.Request) {
	var userID int
	if clientErr, err := utils.ConvertValue(r.URL.Query().Get("userID"), &userID); err != nil {
		handlers.ErrorHandler(w, http.StatusBadRequest, h.log, []string{clientErr}, err)
		return
	}

	userService := db.NewUserService(h.Db)
	_, clientErr, err := userService.CheckUser(userID)
	if err != nil {
		handlers.ErrorHandler(w, http.StatusBadRequest, h.log, []string{clientErr}, err)
		return
	}

	userInfo, clientErr, err := userService.GetUserInfo(userID)
	if err != nil {
		handlers.ErrorHandler(w, http.StatusBadRequest, h.log, []string{clientErr}, err)
		return
	}

	handlers.SendRespond(w, http.StatusOK, &models.Response{
		Errors: nil,
		Data:   userInfo,
	})
}

func (h *UserHandler) SearchUsers(w http.ResponseWriter, r *http.Request) {
	log.Printf("req rec")
	userService := db.NewUserService(h.Db)
	userList, clientErr, err := userService.SearchUsers(r.URL.Query().Get("symbols"))
	if err != nil {
		handlers.ErrorHandler(w, http.StatusBadRequest, h.log, []string{clientErr}, err)
		return
	}

	handlers.SendRespond(w, http.StatusOK, &models.Response{
		Errors: nil,
		Data:   userList,
	})
}

func (h *UserHandler) ChangeUserData(w http.ResponseWriter, r *http.Request) {
	var data *models.ChangeUserData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		handlers.ErrorHandler(w, http.StatusBadRequest, h.log, []string{consts.ErrInternalServerError}, err)
		return
	}

	userService := db.NewUserService(h.Db)
	switch data.Type {
	case "username":
		if clientErr, err := userService.ChangeUsername(data.Data, data.UserID); err != nil {
			handlers.ErrorHandler(w, http.StatusBadRequest, h.log, []string{clientErr}, err)
			return
		}
	case "email":
		if clientErr, err := userService.ChangeEmail(data.Data, data.UserID); err != nil {
			handlers.ErrorHandler(w, http.StatusBadRequest, h.log, []string{clientErr}, err)
			return
		}
	default:
		handlers.ErrorHandler(w, http.StatusBadRequest, h.log, []string{consts.ErrInvalidRequestFormat}, fmt.Errorf(consts.InternalErrChangedUserInfo))
		return
	}

	handlers.SendRespond(w, http.StatusOK, &models.Response{
		Errors: nil,
		Data:   "successful changed",
	})
}
