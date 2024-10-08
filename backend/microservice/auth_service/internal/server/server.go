package server

import (
	"auth_service/internal/config"
	"auth_service/internal/logger"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Server struct contains the router, logger, and config settings
type Server struct {
	rMux   *mux.Router
	log    logger.Logger
	config *config.Config
}

// NewServer initializes a new server instance with routing and logging
func NewServer(cfg *config.Config) *Server {
	log := logger.New(logger.InfoLevel)
	return &Server{
		rMux:   mux.NewRouter(),
		log:    log,
		config: cfg,
	}
}

// Start launches the server on the specified host and port
func (s *Server) Start() {
	s.RegisterHandler()

	s.log.Info("successfully connected to the database!")
	// Используем параметры из конфигурации
	address := fmt.Sprintf("%s:%s", s.config.Server.Host, s.config.Server.Port)
	log.Printf("starting server on %s", address)
	if err := http.ListenAndServe(address, s.rMux); err != nil {
		s.log.Fatalf("Unable start server: %v\n", err)
	}
}

// RegisterHandler defines routes and handlers for the server
func (s *Server) RegisterHandler() {
	s.rMux.HandleFunc("/auth", s.Authorization).Methods("POST")
	s.rMux.HandleFunc("/refreshToken", s.RefreshToken).Methods("POST")
	s.rMux.HandleFunc("/generate_tokens", s.GenerateTokens).Methods("POST")
}
