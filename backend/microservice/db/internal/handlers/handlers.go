// This code represents the interface and its implementation in the `handlers` package. It demonstrates the composition of
// several modules to organise work with HTTP requests in a web application. The `Handler` interface
// describes a set of methods responsible for various actions in the system, such as authentication,
// user profile management, chats, friends, and notifications.
package handlers

import (
	"InterestingChats/backend/microservice/db/internal/db"
	"InterestingChats/backend/microservice/db/internal/logger"
	"net/http"
)

// The `Handler` interface
// This is the key element that defines the contracts for handling HTTP requests related to the functionality of the
// users and chats. Each method accepts an HTTP request and a response (typical signatures for HTTP handlers
// in Go). The methods cover the functionality of: registration, authorisation, profile handling, friendship, chats
// and notifications. This is a high-level interface that allows flexible replacement implementations.

// Composition via struct `handler`.
// The `handler` structure is a composition and encapsulates two components:
// 1. `DBService` - an interface for working with the database.
// 2. `logger.Logger` - interface for logging.
// This composition allows convenient access to the functionality of working with the database and logging in each handler.

// Factory `NewHandler`.
// The `NewHandler` function is a factory that creates a new instance of the handler, taking as arguments
// database service and logger. This is a creation pattern that makes further testing and code maintenance easier,
// since different implementations of these interfaces (e.g., test stubs) can be passed around.

// The `GetDBService` method
// This method implements a part of the `Handler` interface, providing access to the `DBService`. This is useful in cases,
// when other system components also need access to database services.

type Handler interface {
	Registrations(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	CheckUser(w http.ResponseWriter, r *http.Request)
	GetUserProfileInfo(w http.ResponseWriter, r *http.Request)
	SearchUsers(w http.ResponseWriter, r *http.Request)
	ChangeUserData(w http.ResponseWriter, r *http.Request)
	UploadPhoto(w http.ResponseWriter, r *http.Request)
	RequestToFriendShip(w http.ResponseWriter, r *http.Request)
	AcceptFriendShip(w http.ResponseWriter, r *http.Request)
	DeleteFriend(w http.ResponseWriter, r *http.Request)
	DeleteFriendRequest(w http.ResponseWriter, r *http.Request)
	GetFriendList(w http.ResponseWriter, r *http.Request)
	GetSubList(w http.ResponseWriter, r *http.Request)

	GetChat(w http.ResponseWriter, r *http.Request)
	GetAllChats(w http.ResponseWriter, r *http.Request)
	GetChatBySymbols(w http.ResponseWriter, r *http.Request)
	GetUserChats(w http.ResponseWriter, r *http.Request)
	CreateChat(w http.ResponseWriter, r *http.Request)
	DeleteChat(w http.ResponseWriter, r *http.Request)
	DeleteMember(w http.ResponseWriter, r *http.Request)
	AddMembers(w http.ResponseWriter, r *http.Request)
	GetAuthor(w http.ResponseWriter, r *http.Request)
	ChangeChatName(w http.ResponseWriter, r *http.Request)
	SetTag(w http.ResponseWriter, r *http.Request)
	GetTags(w http.ResponseWriter, r *http.Request)
	DeleteTags(w http.ResponseWriter, r *http.Request)

	AddNotification(w http.ResponseWriter, r *http.Request)
	GetNotification(w http.ResponseWriter, r *http.Request)
	ReadNotifications(w http.ResponseWriter, r *http.Request)

	GetDBService() db.DBService
}

type handler struct {
	DBService db.DBService
	log       logger.Logger
}

func NewHandler(dbService db.DBService, log logger.Logger) Handler {
	return &handler{
		DBService: dbService,
		log:       log,
	}
}

func (h *handler) GetDBService() db.DBService {
	return h.DBService
}
