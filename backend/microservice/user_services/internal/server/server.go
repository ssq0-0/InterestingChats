package server

import (
	"InterestingChats/backend/user_services/internal/config"
	"InterestingChats/backend/user_services/internal/handlers"
	"InterestingChats/backend/user_services/internal/logger"
	"InterestingChats/backend/user_services/internal/services"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	rMux    *mux.Router
	Handler *handlers.Handler
	log     logger.Logger
	config  *config.Config
}

// Теперь NewServer принимает конфиг
func NewServer(cfg *config.Config) *Server {
	log := logger.New(logger.InfoLevel)

	// Передаем конфиг с топиками в продюсер
	producer, err := services.NewProducer(cfg.Kafka.BootstrapServers, cfg.Kafka.Topics)
	if err != nil {
		fmt.Printf("Ошибка создания продюсера Kafka: %v\n", err)
		return nil
	}

	userService := services.NewUserService(*producer)
	handler := handlers.NewService(log, userService)

	return &Server{
		rMux:    mux.NewRouter(),
		Handler: handler,
		log:     log,
		config:  cfg,
	}
}

func (s *Server) Start() {
	s.RegisterHandler()

	s.log.Info("успешно подключено к базе данных!")
	address := fmt.Sprintf("%s:%s", s.config.Server.Host, s.config.Server.Port)
	s.log.Infof("сервер запускается на %s", address)
	if err := http.ListenAndServe(address, s.rMux); err != nil {
		s.log.Fatalf("Не удалось запустить сервер: %v\n", err)
	}
}

// RegisterRoutes registers routes for API Gateway.
func (s *Server) RegisterHandler() {
	s.rMux.HandleFunc("/registration", s.Handler.Registrations).Methods("POST")
	s.rMux.HandleFunc("/login", s.Handler.Login).Methods("POST")
	s.rMux.HandleFunc("/my_profile", s.Handler.GetMyProfile).Methods("GET")
	s.rMux.HandleFunc("/user_profile", s.Handler.GetUserProfile).Methods("GET")
	s.rMux.HandleFunc("/changeData", s.Handler.ChangeUserData).Methods("PATCH")
	s.rMux.HandleFunc("/searchUsers", s.Handler.GetSearchUserResult).Methods("GET")

	s.rMux.HandleFunc("/requestToFriendShip", s.Handler.RequestToFriendShip).Methods("POST")
	s.rMux.HandleFunc("/acceptFriendShip", s.Handler.AcceptFriendShip).Methods("POST")
	s.rMux.HandleFunc("/getFriends", s.Handler.GetFriends).Methods("GET")
	s.rMux.HandleFunc("/getSubscribers", s.Handler.GetSubscribers).Methods("GET")
	s.rMux.HandleFunc("/deleteFriend", s.Handler.DeleteFriend).Methods("DELETE")
	s.rMux.HandleFunc("/deleteFriendRequest", s.Handler.DeleteFriendRequest).Methods("DELETE")

	s.rMux.HandleFunc("/saveImage", s.Handler.UploadAvatar).Methods("POST")
}
