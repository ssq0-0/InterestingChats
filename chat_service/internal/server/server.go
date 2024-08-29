package server

import (
	"chat_service/internal/logger"
	"chat_service/internal/models"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type Server struct {
	RMux *mux.Router
	WS   *WS
	log  logger.Logger
}

func NewServer() *Server {
	log := logger.New(logger.InfoLevel)
	return &Server{
		WS: &WS{
			Upgrader: &websocket.Upgrader{
				ReadBufferSize:  1024,
				WriteBufferSize: 1024,
				CheckOrigin:     func(r *http.Request) bool { return true },
			},
			Chats: make(map[string]*models.Chat),
			Mu:    &sync.RWMutex{},
			log:   log,
		},
		RMux: mux.NewRouter(),
		log:  log,
	}
}

func (s *Server) Start() {
	s.RegisterHandlers()

	s.log.Info("server has been started on 8004")

	if err := http.ListenAndServe(":8004", s.RMux); err != nil {
		s.log.Fatalf("could not start server: %v", err)
	}
}

func (s *Server) RegisterHandlers() {
	s.RMux.HandleFunc("/ws", s.GetChatHistory).Methods("GET")
	s.RMux.HandleFunc("/getAllChats", s.GetAllChats).Methods("GET")
	s.RMux.HandleFunc("/wsOpen", s.WS.ChatWebsocket).Methods("GET")
	s.RMux.HandleFunc("/getChat", s.GetChatHistory).Methods("GET")
	s.RMux.HandleFunc("/getUserChats", s.GetUserChats).Methods("GET")
	s.RMux.HandleFunc("/getChatBySymbol", s.SearchChat).Methods("GET")
	s.RMux.HandleFunc("/createChat", s.CreateChat).Methods("POST")
	s.RMux.HandleFunc("/deleteChat", s.DeleteChat).Methods("DELETE")
	s.RMux.HandleFunc("/addMember", s.AddMember).Methods("POST")
	s.RMux.HandleFunc("/deleteMember", s.DeleteMember).Methods("DELETE")
	s.RMux.HandleFunc("/joinToChat", s.JoinToChat).Methods("POST")
}
