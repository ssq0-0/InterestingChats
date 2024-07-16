package server

import (
	"InterestingChats/backend/microservice/db/internal/db"
	chat "InterestingChats/backend/microservice/db/internal/handlers/chat_handlers"
	"InterestingChats/backend/microservice/db/internal/handlers/user"
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	rMux        *mux.Router
	Db          *sql.DB
	UserHandler *user.UserHandler
	ChatHandler *chat.ChatHandler
}

func NewServer(DB *sql.DB) *Server {
	UserHandler := user.NewHandler(DB)
	chatService := db.NewChatService(DB)
	chatHandler := chat.NewChatHandler(chatService)
	return &Server{
		rMux:        mux.NewRouter(),
		Db:          DB,
		UserHandler: UserHandler,
		ChatHandler: chatHandler,
	}
}

func (s *Server) Start() {
	s.RegisterHandler()
	log.Println("Successfully connected to the database!")
	log.Println("Starting server on :8002")
	if err := http.ListenAndServe(":8002", s.rMux); err != nil {
		log.Fatalf("Unable start server: %v\n", err)
	}
}

func (s *Server) RegisterHandler() {
	s.rMux.HandleFunc("/registration", s.UserHandler.Registrations).Methods("POST")
	s.rMux.HandleFunc("/login", s.UserHandler.Login).Methods("POST")
	s.rMux.HandleFunc("/getChat", s.ChatHandler.GetChat).Methods("GET")
	s.rMux.HandleFunc("/createChat", s.ChatHandler.CreateChat).Methods("POST")
	log.Println("Continue...")
}
