package server

import (
	"InterestingChats/backend/microservice/db/internal/db"
	"InterestingChats/backend/microservice/db/internal/handlers"
	"InterestingChats/backend/microservice/db/internal/logger"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// Server to represent the server structure containing routes, handlers, and database connection
type Server struct {
	rMux    *mux.Router
	log     logger.Logger
	port    int
	Db      *sql.DB
	Handler handlers.Handler
}

// NewServer to create a new server instance with the provided database and port
func NewServer(DB *sql.DB, port int) *Server {
	log := logger.New(logger.InfoLevel)
	DBService := db.NewDBService(DB)
	Handler := handlers.NewHandler(DBService, log)

	return &Server{
		rMux:    mux.NewRouter(),
		port:    port,
		Db:      DB,
		Handler: Handler,
		log:     log,
	}
}

// Start to start the server and register the HTTP handlers
func (s *Server) Start() {
	s.RegisterHandler()

	s.log.Info("successfully connected to the database!")
	if err := http.ListenAndServe(fmt.Sprintf(":%d", s.port), s.rMux); err != nil {
		s.log.Fatalf("Unable start server: %v\n", err)
	}

}

// RegisterHandler to register all HTTP routes and handlers
func (s *Server) RegisterHandler() {
	s.rMux.HandleFunc("/registration", s.Handler.Registrations).Methods("POST")
	s.rMux.HandleFunc("/login", s.Handler.Login).Methods("POST")
	s.rMux.HandleFunc("/profileInfo", s.Handler.GetUserProfileInfo).Methods("GET")
	s.rMux.HandleFunc("/checkUser", s.Handler.CheckUser).Methods("GET")
	s.rMux.HandleFunc("/changeUserData", s.Handler.ChangeUserData).Methods("POST")
	s.rMux.HandleFunc("/searchUsers", s.Handler.SearchUsers).Methods("GET")
	s.rMux.HandleFunc("/uploadPhoto", s.Handler.UploadPhoto).Methods("POST")

	s.rMux.HandleFunc("/requestToFriendShip", s.Handler.RequestToFriendShip).Methods("POST")
	s.rMux.HandleFunc("/acceptFriendShip", s.Handler.AcceptFriendShip).Methods("POST")
	s.rMux.HandleFunc("/getFriends", s.Handler.GetFriendList).Methods("GET")
	s.rMux.HandleFunc("/getSubs", s.Handler.GetSubList).Methods("GET")
	s.rMux.HandleFunc("/deleteFriend", s.Handler.DeleteFriend).Methods("DELETE")
	s.rMux.HandleFunc("/deleteFriendRequest", s.Handler.DeleteFriendRequest).Methods("DELETE")

	s.rMux.HandleFunc("/getChat", s.Handler.GetChat).Methods("GET")
	s.rMux.HandleFunc("/getUserChats", s.Handler.GetUserChats).Methods("GET")
	s.rMux.HandleFunc("/getAllChats", s.Handler.GetAllChats).Methods("GET")
	s.rMux.HandleFunc("/searchChat", s.Handler.GetChatBySymbols).Methods("GET")
	s.rMux.HandleFunc("/createChat", s.Handler.CreateChat).Methods("POST")
	s.rMux.HandleFunc("/deleteChat", s.Handler.DeleteChat).Methods("DELETE")
	s.rMux.HandleFunc("/deleteMember", s.Handler.DeleteMember).Methods("DELETE")
	s.rMux.HandleFunc("/addMembers", s.Handler.AddMembers).Methods("POST")
	s.rMux.HandleFunc("/joinToChat", s.Handler.AddMembers).Methods("POST")
	s.rMux.HandleFunc("/getAuthor", s.Handler.GetAuthor).Methods("GET")
	s.rMux.HandleFunc("/changeChatName", s.Handler.ChangeChatName).Methods("PATCH")

	s.rMux.HandleFunc("/getNotification", s.Handler.GetNotification).Methods("GET")
	s.rMux.HandleFunc("/readNotification", s.Handler.ReadNotifications).Methods("PATCH")
	s.rMux.HandleFunc("/addNotification", s.Handler.AddNotification).Methods("POST")

	s.rMux.HandleFunc("/setTag", s.Handler.SetTag).Methods("PATCH")
	s.rMux.HandleFunc("/getTags", s.Handler.GetTags).Methods("GET")
	s.rMux.HandleFunc("/deleteTags", s.Handler.DeleteTags).Methods("DELETE")
}
