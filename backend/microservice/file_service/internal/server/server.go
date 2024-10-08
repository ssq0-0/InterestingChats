package server

import (
	"file_service/internal/config"
	"file_service/internal/handlers"
	"file_service/internal/logger"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// Server represents the HTTP server that handles incoming requests.
type Server struct {
	log          logger.Logger
	RMux         *mux.Router
	ImageService *handlers.ImageService
}

// NewServer initializes a new Server instance with logger and MinIO client.
func NewServer(cfg *config.Config) *Server {
	log := logger.New(logger.InfoLevel)

	minioClient, err := minio.New(cfg.Minio.Endpoint, &minio.Options{
		Creds:  credentials.NewStatic(cfg.Minio.AccessKeyID, cfg.Minio.SecretAccessKey, "", credentials.SignatureV4),
		Secure: true,
	})
	if err != nil {
		log.Fatalf("Error creating MinIO client: %v", err)
	}

	return &Server{
		log:          log,
		RMux:         mux.NewRouter(),
		ImageService: handlers.NewImageService(log, minioClient, cfg.Minio.BucketName),
	}
}

// Start starts the HTTP server and registers routes.
func (s *Server) Start(cfg *config.Config) {
	s.RegisterRoutes()

	s.log.Infof("File system starting server on %s:%s!", cfg.Server.Host, cfg.Server.Port)

	if err := http.ListenAndServe(cfg.Server.Host+":"+cfg.Server.Port, s.RMux); err != nil {
		s.log.Panicf("Failed to start file system: %w", err)
	}
}

// RegisterRoutes defines the routes for the server.
func (s *Server) RegisterRoutes() {
	s.RMux.HandleFunc("/saveImage", s.ImageService.SaveImage).Methods("POST")
	s.RMux.HandleFunc("/deleteImage", s.ImageService.DeleteImage).Methods("DELETE")
}
