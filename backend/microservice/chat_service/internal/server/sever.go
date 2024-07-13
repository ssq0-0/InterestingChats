package server

import (
	"log"
	"net/http"
)

type Server struct {
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Start() {
	s.Registerhandlers()
	log.Println("Server start on 8004 port!")
	if err := http.ListenAndServe(":8004", nil); err != nil { // inject 'rMux'
		panic(err)
	}
}

func (s *Server) Registerhandlers() {

}
