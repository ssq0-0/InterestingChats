package services

import (
	"InterestingChats/backend/user_services/internal/logger"
	"InterestingChats/backend/user_services/internal/models"
	"net/http"
)

type UserService interface {
	UserRequest(method, url string, user models.User, expectedStatusCode int) (*models.UserTokens, *models.User, string, error)
	ChangeUserData(r *http.Request) (*models.Response, int, string, error)
	AddUserAvatar(r *http.Request) (*models.Response, int, string, error)
	FriendsOperations(r *http.Request, typeRequest int) (*models.Response, int, string, error)
	GetMyProfile(r *http.Request, log logger.Logger) (*models.Response, int, string, error)
	GetUserInfo(r *http.Request, log logger.Logger) (*models.Response, int, string, error)
	GetSearchUserResult(r *http.Request) (*models.Response, int, string, error)
}

type userService struct {
	producer Producer
}

func NewUserService(producer Producer) UserService {
	return &userService{
		producer: producer,
	}
}
