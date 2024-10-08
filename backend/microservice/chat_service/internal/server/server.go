package server

import (
	"chat_service/internal/config"
	"chat_service/internal/logger"
	"chat_service/internal/models"
	"chat_service/internal/services/ws_service"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// Server represents the main server structure, holding the router, WebSocket manager, logger, and port.
type Server struct {
	RMux *mux.Router
	WS   *WS
	log  logger.Logger
	Port string
}

// NewServer creates a new instance of the server based on the configuration file, setting up WebSocket and Kafka producer.
func NewServer(cfg *config.Config) *Server {
	log := logger.New(logger.InfoLevel)
	kfp, err := ws_service.NewProducer(cfg.Kafka)
	if err != nil {
		log.Errorf("cannot create producer: %v", err)
		return nil
	}
	return &Server{
		WS: &WS{
			Upgrader: &websocket.Upgrader{
				ReadBufferSize:  1024,
				WriteBufferSize: 1024,
				CheckOrigin:     func(r *http.Request) bool { return true },
			},
			Chats:    make(map[string]*models.Chat),
			Mu:       &sync.RWMutex{},
			Producer: kfp,
			log:      log,
		},
		RMux: mux.NewRouter(),
		log:  log,
		Port: cfg.Server.Port,
	}
}

// Start begins the HTTP server on the specified port and registers all routes.
func (s *Server) Start() {
	s.RegisterHandlers()

	s.log.Infof("server has been started on port %s", s.Port)

	if err := http.ListenAndServe(fmt.Sprintf(":%s", s.Port), s.RMux); err != nil {
		s.log.Fatalf("could not start server: %v", err)
	}
}

// RegisterHandlers sets up the routes for handling different HTTP requests related to chats and users.
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
	s.RMux.HandleFunc("/setTag", s.SetTag).Methods("PATCH")
	s.RMux.HandleFunc("/getTags", s.GetTags).Methods("GET")
	s.RMux.HandleFunc("/deleteTags", s.DeleteTags).Methods("DELETE")
	s.RMux.HandleFunc("/leaveChat", s.LeaveChat).Methods("DELETE")
	s.RMux.HandleFunc("/changeChatName", s.ChangeChatName).Methods("PATCH")
}
