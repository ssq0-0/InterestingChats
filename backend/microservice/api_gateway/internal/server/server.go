package server

import (
	"InterestingChats/backend/api_gateway/internal/proxy"
	"fmt"
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
	s.rMux.HandleFunc("/updateToken", proxy.ProxyRequest("http://localhost:8001")).Methods("POST")
	fmt.Println("Route /registration registered for POST method")
	fmt.Println("Route /login registered for POST method")
	fmt.Println("Route /updateToken registered for POST method")
}
