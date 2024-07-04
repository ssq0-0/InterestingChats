package services

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
)

type Server struct {
	Db *sql.DB
}

func NewServer(db *sql.DB) *Server {
	return &Server{Db: db}
}

func (s *Server) Start() {
	fmt.Println("Successfully connected to the database!")
	log.Println("Starting server on :8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Unable start server: %v\n", err)
	}
}
