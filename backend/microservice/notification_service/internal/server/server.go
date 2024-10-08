package server

import (
	"net/http"
	"notifications/internal/handlers"
	"notifications/internal/logger"

	"github.com/gorilla/mux"
)

// Server represents the HTTP server that handles incoming requests
type Server struct {
	rMux *mux.Router
	log  logger.Logger
	host string
	port string
}

// NewServer initializes a new HTTP server with a logger and router
func NewServer(host, port string) *Server {
	log := logger.New(logger.InfoLevel)
	return &Server{
		rMux: mux.NewRouter(),
		log:  log,
		host: host,
		port: port,
	}
}

// Start begins the HTTP server and listens for incoming requests
func (s *Server) Start() {
	s.RegisterHandler()

	address := s.host + ":" + s.port
	s.log.Infof("starting server on %s", address) // Log the address to start the server
	if err := http.ListenAndServe(address, s.rMux); err != nil {
		s.log.Fatalf("Unable to start server: %v\n", err)
	}
}

// RegisterHandler sets up the HTTP routes and their corresponding handlers
func (s *Server) RegisterHandler() {
	s.rMux.HandleFunc("/getNotification", handlers.GetNotifications).Methods("GET")
	s.rMux.HandleFunc("/readNotification", handlers.ReadNotification).Methods("PATCH") // TODO: refactor
}
