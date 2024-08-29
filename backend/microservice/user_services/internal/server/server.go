package server

import (
	"InterestingChats/backend/user_services/internal/handlers"
	"InterestingChats/backend/user_services/internal/logger"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	rMux    *mux.Router
	Handler *handlers.UserService
	log     logger.Logger
}

func NewServer() *Server {
	log := logger.New(logger.InfoLevel)
	handler := handlers.NewService(log)
	return &Server{
		rMux:    mux.NewRouter(),
		Handler: handler,
		log:     log,
	}
}

func (s *Server) Start() {
	s.RegisterHandler()
	// headersOk := hand.AllowedHeaders([]string{"Authorization", "Content-Type"})
	// originsOk := hand.AllowedOrigins([]string{"*"})
	// methodsOk := hand.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS", "DELETE"})

	s.log.Info("successfully connected to the database!")
	s.log.Info("starting server on :8001")
	// if err := http.ListenAndServe(":8001", hand.CORS(originsOk, headersOk, methodsOk)(s.rMux)); err != nil {
	if err := http.ListenAndServe(":8001", s.rMux); err != nil {
		s.log.Fatalf("Unable start server: %v\n", err)
	}
}

func (s *Server) RegisterHandler() {
	s.rMux.HandleFunc("/registration", s.Handler.Registrations).Methods("POST")
	s.rMux.HandleFunc("/login", s.Handler.Login).Methods("POST")
	s.rMux.HandleFunc("/getTokens", s.Handler.GetTokens).Methods("GET")
	s.rMux.HandleFunc("/checkToken", s.Handler.CheckTokens).Methods("GET")
	s.rMux.HandleFunc("/refreshToken", s.Handler.RefreshAccessToken).Methods("POST")
	s.rMux.HandleFunc("/my_profile", s.Handler.GetMyProfile).Methods("GET")
	s.rMux.HandleFunc("/user_profile", s.Handler.GetUserProfile).Methods("GET")
	s.rMux.HandleFunc("/changeData", s.Handler.ChangeUserData).Methods("POST")
	s.rMux.HandleFunc("/searchUsers", s.Handler.GetSearchUserResult).Methods("GET")
}
