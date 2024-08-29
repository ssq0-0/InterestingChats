package server

import (
	"InterestingChats/backend/api_gateway/internal/proxy"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	rMux *mux.Router
}

func NewServer() *Server {
	return &Server{rMux: mux.NewRouter()}
}

func (s *Server) Start() {
	s.RegisterRoutes()
	fmt.Println("Gateway server running on 8000 port!")
	if err := http.ListenAndServe(":8000", s.rMux); err != nil {
		panic(err)
	}
}

func (s *Server) RegisterRoutes() {
	fmt.Println("Registering routes...")
	s.rMux.HandleFunc("/registration", proxy.ProxyRequest("http://localhost:8001")).Methods("POST")
	s.rMux.HandleFunc("/login", proxy.ProxyRequest("http://localhost:8001")).Methods("POST")
	s.rMux.HandleFunc("/getTokens", proxy.ProxyRequest("http://localhost:8001")).Methods("GET")
	log.Println("Route /registration registered for POST method")
	log.Println("Route /login registered for POST method")
	log.Println("Route /getTokens registered for POST method")
}
