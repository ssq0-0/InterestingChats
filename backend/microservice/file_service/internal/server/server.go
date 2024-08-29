package server

import (
	"file_service/internal/handlers"
	"file_service/internal/logger"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Server struct {
	log          logger.Logger
	RMux         *mux.Router
	ImageService *handlers.ImageService
}

func NewServer() *Server {
	log := logger.New(logger.InfoLevel)
	endpoint := "storage.yandexcloud.net"
	accessKeyID := "YCAJEG1yAQpkZi45pIzd7Oxqj"
	secretAccessKey := "YCMc38di6R0T50M9Vhed9QF6oPiwLZ6gsU0-U62n"

	// Создание клиента MinIO
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStatic(accessKeyID, secretAccessKey, "", credentials.SignatureV4),
		Secure: true, // Используйте true для HTTPS
	})
	if err != nil {
		log.Fatalf("Ошибка создания клиента MinIO: %v", err)
	}

	return &Server{
		log:          log,
		RMux:         mux.NewRouter(),
		ImageService: handlers.NewImageService(log, minioClient, "chatsapp"),
	}
}

func (s *Server) Start() {
	s.RegisterRoutes()

	s.log.Infof("file system starting server on 8005 port!")

	if err := http.ListenAndServe(":8005", s.RMux); err != nil {
		s.log.Panicf("failed to start file system: %w", err)
	}
}

func (s *Server) RegisterRoutes() {
	s.RMux.HandleFunc("/saveImage", s.ImageService.SaveImage).Methods("POST")
}
