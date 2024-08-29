package server

import (
	"InterestingChats/backend/api_gateway/internal/consts"
	"InterestingChats/backend/api_gateway/internal/logger"
	"InterestingChats/backend/api_gateway/internal/proxy"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Server struct {
	rMux *mux.Router
	log  logger.Logger
}

func NewServer() *Server {
	return &Server{
		rMux: mux.NewRouter(),
		log:  logger.New(logger.InfoLevel),
	}
}

func (s *Server) Start() {
	s.RegisterRoutes()

	headersOk := handlers.AllowedHeaders([]string{"Authorization", "Content-Type"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS", "DELETE"})

	s.log.Info("gateway server running on 800!")
	if err := http.ListenAndServe(":8000", handlers.CORS(originsOk, headersOk, methodsOk)(s.rMux)); err != nil {
		// if err := http.ListenAndServe(":8000", s.rMux); err != nil {
		s.log.Panic(err)
	}
}

func (s *Server) RegisterRoutes() {
	s.rMux.HandleFunc("/registration", proxy.ProxyRequest(consts.SERVER_user_service, s.log)).Methods("POST")
	s.rMux.HandleFunc("/login", proxy.ProxyRequest(consts.SERVER_user_service, s.log)).Methods("POST")
	s.rMux.HandleFunc("/user_profile", proxy.ProxyRequest(consts.SERVER_user_service, s.log)).Methods("GET")
	s.rMux.HandleFunc("/my_profile", proxy.ProxyRequest(consts.SERVER_user_service, s.log)).Methods("GET")
	s.rMux.HandleFunc("/changeData", proxy.ProxyRequest(consts.SERVER_user_service, s.log)).Methods("POST")
	s.rMux.HandleFunc("/getTokens", proxy.ProxyRequest(consts.SERVER_user_service, s.log)).Methods("GET")
	s.rMux.HandleFunc("/checkToken", proxy.ProxyRequest(consts.SERVER_user_service, s.log)).Methods("GET")
	s.rMux.HandleFunc("/searchUsers", proxy.ProxyRequest(consts.SERVER_user_service, s.log)).Methods("GET")
	s.rMux.HandleFunc("/setToken", proxy.ProxyRequest(consts.SERVER_user_service, s.log)).Methods("POST")
	s.rMux.HandleFunc("/refreshToken", proxy.ProxyRequest(consts.SERVER_user_service, s.log)).Methods("POST")

	s.rMux.HandleFunc("/saveImage", proxy.ProxyRequest(consts.SERVER_file_system, s.log)).Methods("POST")

	s.rMux.HandleFunc("/deleteTokens", proxy.ProxyRequest(consts.SERVER_redis, s.log)).Methods("DELETE")

	s.rMux.HandleFunc("/ws", proxy.ProxyRequest(consts.SERVER_chat_service, s.log)).Methods("GET")
	s.rMux.HandleFunc("/wsOpen", proxy.ProxyRequest(consts.SERVER_chat_service, s.log)).Methods("GET")
	s.rMux.HandleFunc("/getChat", proxy.ProxyRequest(consts.SERVER_chat_service, s.log)).Methods("GET")
	s.rMux.HandleFunc("/getChatBySymbol", proxy.ProxyRequest(consts.SERVER_chat_service, s.log)).Methods("GET")
	s.rMux.HandleFunc("/getAllChats", proxy.ProxyRequest(consts.SERVER_chat_service, s.log)).Methods("GET")
	s.rMux.HandleFunc("/getUserChats", proxy.ProxyRequest(consts.SERVER_chat_service, s.log)).Methods("GET")
	s.rMux.HandleFunc("/joinToChat", proxy.ProxyRequest(consts.SERVER_chat_service, s.log)).Methods("POST")
	s.rMux.HandleFunc("/createChat", proxy.ProxyRequest(consts.SERVER_chat_service, s.log)).Methods("POST")
	s.rMux.HandleFunc("/deleteChat", proxy.ProxyRequest(consts.SERVER_chat_service, s.log)).Methods("DELETE")
	s.rMux.HandleFunc("/addMember", proxy.ProxyRequest(consts.SERVER_chat_service, s.log)).Methods("POST")
	s.rMux.HandleFunc("/deleteMember", proxy.ProxyRequest(consts.SERVER_chat_service, s.log)).Methods("DELETE")
	s.rMux.HandleFunc("/saveMessage", proxy.ProxyRequest(consts.SERVER_chat_service, s.log)).Methods("POST")
}
