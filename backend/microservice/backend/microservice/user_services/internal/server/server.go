package server

import (
	"InterestingChats/backend/user_services/internal/handlers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	rMux    *mux.Router
	Handler *handlers.Handler
}

func NewServer() *Server {
	handler := handlers.NewHandler()
	return &Server{
		rMux:    mux.NewRouter(),
		Handler: handler,
	}
}

func (s *Server) Start() {
	s.RegisterHandler()
	log.Println("Successfully connected to the database!")
	log.Println("Starting server on :8001")
	if err := http.ListenAndServe(":8001", s.rMux); err != nil {
		log.Fatalf("Unable start server: %v\n", err)
	}
}

func (s *Server) RegisterHandler() {
	s.rMux.HandleFunc("/registration", s.Handler.Registrations).Methods("POST")
	s.rMux.HandleFunc("/login", s.Handler.Login).Methods("POST")
	s.rMux.HandleFunc("/getTokens", s.Handler.GetTokens).Methods("GET")
	s.rMux.HandleFunc("/checkToken", s.Handler.CheckTokens).Methods("GET")
	log.Println("Continue...")
}
