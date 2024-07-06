package server

import (
	"InterestingChats/backend/user_services/internal/handlers"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	rMux    *mux.Router
	Db      *sql.DB
	Handler *handlers.Handler
}

func NewServer(db *sql.DB) *Server {
	handler := handlers.NewHandler(db)
	return &Server{
		rMux:    mux.NewRouter(),
		Db:      db,
		Handler: handler,
	}
}

func (s *Server) Start() {
	s.RegisterHandler()
	fmt.Println("Successfully connected to the database!")
	log.Println("Starting server on :8001")
	if err := http.ListenAndServe(":8001", s.rMux); err != nil {
		log.Fatalf("Unable start server: %v\n", err)
	}
}

func (s *Server) RegisterHandler() {
	s.rMux.HandleFunc("/registration", s.Handler.Registrations).Methods("POST")
	s.rMux.HandleFunc("/login", s.Handler.Login).Methods("POST")
	s.rMux.HandleFunc("/updateToken", s.Handler.UpdateAccessToken).Methods("POST")
	fmt.Println("Continue...")
}
