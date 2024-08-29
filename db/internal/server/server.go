package server

import (
	"InterestingChats/backend/microservice/db/internal/db"
	chat "InterestingChats/backend/microservice/db/internal/handlers/chat_handlers"
	"InterestingChats/backend/microservice/db/internal/handlers/user"
	"InterestingChats/backend/microservice/db/internal/logger"
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	rMux        *mux.Router
	Db          *sql.DB
	UserHandler *user.UserHandler
	ChatHandler *chat.ChatHandler
	log         logger.Logger
}

func NewServer(DB *sql.DB) *Server {
	log := logger.New(logger.InfoLevel)
	UserHandler := user.NewHandler(DB, log)
	ChatService := db.NewChatService(DB)
	ChatHandler := chat.NewChatHandler(ChatService, log)
	return &Server{
		rMux:        mux.NewRouter(),
		Db:          DB,
		UserHandler: UserHandler,
		ChatHandler: ChatHandler,
		log:         log,
	}
}

func (s *Server) Start() {
	s.RegisterHandler()
	s.log.Info("successfully connected to the database!")
	s.log.Info("starting server on :8002")
	if err := http.ListenAndServe(":8002", s.rMux); err != nil {
		s.log.Fatalf("Unable start server: %v\n", err)
	}
}

func (s *Server) RegisterHandler() {
	s.rMux.HandleFunc("/registration", s.UserHandler.Registrations).Methods("POST")
	s.rMux.HandleFunc("/login", s.UserHandler.Login).Methods("POST")
	s.rMux.HandleFunc("/profileInfo", s.UserHandler.GetUserProfileInfo).Methods("GET")
	s.rMux.HandleFunc("/checkUser", s.UserHandler.CheckUser).Methods("GET")
	s.rMux.HandleFunc("/changeUserData", s.UserHandler.ChangeUserData).Methods("POST")
	s.rMux.HandleFunc("/searchUsers", s.UserHandler.SearchUsers).Methods("GET")

	s.rMux.HandleFunc("/getChat", s.ChatHandler.GetChat).Methods("GET")
	s.rMux.HandleFunc("/getUserChats", s.ChatHandler.GetUserChats).Methods("GET")
	s.rMux.HandleFunc("/getAllChats", s.ChatHandler.GetAllChats).Methods("GET")
	s.rMux.HandleFunc("/searchChat", s.ChatHandler.GetChatBySymbols).Methods("GET")
	s.rMux.HandleFunc("/createChat", s.ChatHandler.CreateChat).Methods("POST")
	s.rMux.HandleFunc("/deleteChat", s.ChatHandler.DeleteChat).Methods("DELETE")
	s.rMux.HandleFunc("/deleteMember", s.ChatHandler.DeleteMember).Methods("DELETE")
	s.rMux.HandleFunc("/addMembers", s.ChatHandler.AddMembers).Methods("POST")
	// s.rMux.HandleFunc("/joinToChat", s.ChatHandler.JoinToChat).Methods("POST")
	s.rMux.HandleFunc("/getAuthor", s.ChatHandler.GetAuthor).Methods("GET")
	s.rMux.HandleFunc("/saveMessage", s.ChatHandler.SaveMessage).Methods("POST")
}
