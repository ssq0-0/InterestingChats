package handlers

import (
	"InterestingChats/backend/microservice/redis/internal/logger"
	"InterestingChats/backend/microservice/redis/internal/rdb"
	"net/http"
)

// For future scaling and connection of other various handlers, for example - handlers for caching chats or files.
type Handler interface {
	GetSession(w http.ResponseWriter, r *http.Request)
	GetFriends(w http.ResponseWriter, r *http.Request)
	GetSubscribers(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	rdb *rdb.RedisClient
	log *logger.Logger
}

func NewHandler(rdbClient *rdb.RedisClient, log *logger.Logger) Handler {
	return &handler{
		rdb: rdbClient,
		log: log,
	}
}
