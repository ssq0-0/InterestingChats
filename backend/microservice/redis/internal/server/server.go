package server

import (
	"InterestingChats/backend/microservice/redis/internal/handlers"
	"InterestingChats/backend/microservice/redis/internal/logger"
	"InterestingChats/backend/microservice/redis/internal/rdb"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	rMux        *mux.Router
	rdb         *rdb.RedisClient
	log         logger.Logger
	userHandler *handlers.UserHandler
}

func NewServer(rdbClient *rdb.RedisClient) *Server {
	return &Server{
		rMux:        mux.NewRouter(),
		rdb:         rdbClient,
		log:         logger.New(logger.InfoLevel),
		userHandler: handlers.NewUserHandler(rdbClient, logger.New(logger.InfoLevel)),
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
	s.rMux.HandleFunc("/user", s.userHandler.GetUsersTokens).Methods("GET")
	s.rMux.HandleFunc("/setToken", s.userHandler.SetTokens).Methods("POST")
	s.rMux.HandleFunc("/deleteTokens", s.userHandler.DeleteTokens).Methods("DELETE")
}
