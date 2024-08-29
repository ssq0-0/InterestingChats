package server

import (
	"chat_service/internal/models"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type Server struct {
	Upgrader *websocket.Upgrader
	RMux     *mux.Router
	Chats    map[string]*models.Chat
	Mu       *sync.Mutex
}

func NewServer() *Server {
	return &Server{
		Upgrader: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		RMux:  mux.NewRouter(),
		Chats: make(map[string]*models.Chat),
		Mu:    &sync.Mutex{},
	}
}

func (s *Server) Start() {
	s.RegisterHandlers()
	log.Println("Server start on 8004 port!")
	if err := http.ListenAndServe(":8004", s.RMux); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}

func (s *Server) RegisterHandlers() {
	s.RMux.HandleFunc("/ws", s.GetChatHistory).Methods("GET")
	s.RMux.HandleFunc("/wsOpen", s.OpenWS).Methods("GET")
	s.RMux.HandleFunc("/getChat", s.GetChatHistory).Methods("GET")
	s.RMux.HandleFunc("/createChat", s.CreateChat).Methods("POST")
	s.RMux.HandleFunc("/deleteChat", s.DeleteChat).Methods("DELETE")
	s.RMux.HandleFunc("/addMember", s.AddMember).Methods("POST")
	s.RMux.HandleFunc("/deleteMember", s.DeleteMember).Methods("DELETE")
}
