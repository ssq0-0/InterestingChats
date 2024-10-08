package server

import (
	"InterestingChats/backend/microservice/redis/internal/handlers"
	"InterestingChats/backend/microservice/redis/internal/logger"
	"InterestingChats/backend/microservice/redis/internal/rdb"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	rMux    *mux.Router
	rdb     *rdb.RedisClient
	log     logger.Logger
	Handler handlers.Handler
}

func NewServer(rdbClient *rdb.RedisClient) *Server {
	log := logger.New(logger.InfoLevel)
	return &Server{
		rMux:    mux.NewRouter(),
		rdb:     rdbClient,
		log:     logger.New(logger.InfoLevel),
		Handler: handlers.NewHandler(rdbClient, &log),
	}
}

func (s *Server) Start() {
	s.RegisterHandler()

	s.log.Info("redis DB running on 8003 port!")
	if err := http.ListenAndServe(":8003", s.rMux); err != nil {
		s.log.Fatalf("Unable to start redis server: %v\n", err)
	}
}

func (s *Server) RegisterHandler() {
	s.rMux.HandleFunc("/getSession", s.Handler.GetSession).Methods("GET")
	s.rMux.HandleFunc("/getFriendList", s.Handler.GetFriends).Methods("GET")
	s.rMux.HandleFunc("/getSubscribers", s.Handler.GetSubscribers).Methods("GET")
}
