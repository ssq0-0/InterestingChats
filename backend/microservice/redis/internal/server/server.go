package server

import (
	"InterestingChats/backend/microservice/redis/internal/handlers"
	"InterestingChats/backend/microservice/redis/internal/rdb"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	rMux *mux.Router
	rdb  *rdb.RedisClient
}

func NewServer(rdbClient *rdb.RedisClient) *Server {
	return &Server{
		rMux: mux.NewRouter(),
		rdb:  rdbClient,
	}
}

func (s *Server) Start() {
	s.RegisterHandler()
	log.Println("Redis DB running on 8003 port!")
	if err := http.ListenAndServe(":8003", s.rMux); err != nil {
		log.Fatalf("Unable to start redis server: %v\n", err)
	}
}

func (s *Server) RegisterHandler() {
	userHandler := handlers.NewUserHandler(s.rdb)
	s.rMux.HandleFunc("/user", userHandler.GetUsersTokens).Methods("GET")
	s.rMux.HandleFunc("/setToken", userHandler.SetTokens).Methods("POST")
}
