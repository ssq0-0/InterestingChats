package consts

import "InterestingChats/backend/api_gateway/internal/config"

// Server constants for various microservices.
var (
	SERVER_user_service         string
	SERVER_chat_service         string
	SERVER_redis                string
	SERVER_file_system          string
	SERVER_notification_service string
	SERVER_auth_service         string
)

// InitConstants initializes service URLs from config data.
func InitConstants(cfg *config.Config) {
	SERVER_user_service = cfg.Services.UserService
	SERVER_chat_service = cfg.Services.ChatService
	SERVER_redis = cfg.Services.RedisService
	SERVER_file_system = cfg.Services.FileSystemService
	SERVER_notification_service = cfg.Services.NotificationService
	SERVER_auth_service = cfg.Services.AuthService
}
