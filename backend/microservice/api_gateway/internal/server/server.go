package server

import (
	"InterestingChats/backend/api_gateway/internal/consts"
	"InterestingChats/backend/api_gateway/internal/logger"
	"InterestingChats/backend/api_gateway/internal/proxy"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type Server struct {
	log logger.Logger
	App *fiber.App
}

func NewServer() *Server {
	app := *fiber.New()
	return &Server{
		log: logger.New(logger.InfoLevel),
		App: &app,
	}
}

func (s *Server) Start() {
	s.App.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowHeaders:     "Authorization, Content-Type, X-User-ID",
		AllowMethods:     "GET, POST, PUT, PATCH, DELETE, OPTIONS",
		AllowCredentials: false,
	}))
	s.RegisterRoutes()

	s.log.Info("gateway server running on 800!")
	if err := s.App.Listen(":8000"); err != nil {
		s.log.Panic(err)
	}
}

func (s *Server) RegisterRoutes() {
	// --------------------------------------NO PROTECT------------------------------------------------------------------------------- //
	s.App.Post("/registration", proxy.GatewayProxyRequest(consts.SERVER_user_service, s.log))
	s.App.Post("/login", proxy.GatewayProxyRequest(consts.SERVER_user_service, s.log))
	s.App.Post("/refreshToken", proxy.GatewayProxyRequest(consts.SERVER_auth_service, s.log))

	protected := s.App.Group("", AuthMiddleware(s.log))
	// --------------------------------------USER------------------------------------------------------------------------------- //
	protected.Get("/my_profile", proxy.GatewayProxyRequest(consts.SERVER_user_service, s.log))
	protected.Get("/user_profile", proxy.GatewayProxyRequest(consts.SERVER_user_service, s.log))
	protected.Patch("/changeData", proxy.GatewayProxyRequest(consts.SERVER_user_service, s.log))
	protected.Get("/searchUsers", proxy.GatewayProxyRequest(consts.SERVER_user_service, s.log))

	// --------------------------------------TOKENS------------------------------------------------------------------------------- //
	protected.Post("/auth", proxy.GatewayProxyRequest(consts.SERVER_auth_service, s.log))

	// --------------------------------------FRIENDS------------------------------------------------------------------------------- //
	protected.Get("/getFriends", proxy.GatewayProxyRequest(consts.SERVER_user_service, s.log))
	protected.Get("/getSubscribers", proxy.GatewayProxyRequest(consts.SERVER_user_service, s.log))
	protected.Post("/requestToFriendShip", proxy.GatewayProxyRequest(consts.SERVER_user_service, s.log))
	protected.Post("/acceptFriendShip", proxy.GatewayProxyRequest(consts.SERVER_user_service, s.log))
	protected.Delete("/deleteFriend", proxy.GatewayProxyRequest(consts.SERVER_user_service, s.log))
	protected.Delete("/deleteFriendRequest", proxy.GatewayProxyRequest(consts.SERVER_user_service, s.log))

	// --------------------------------------NOTIFICATION------------------------------------------------------------------------------- //
	protected.Get("/getNotification", proxy.GatewayProxyRequest(consts.SERVER_notification_service, s.log))
	protected.Patch("/readNotification", proxy.GatewayProxyRequest(consts.SERVER_notification_service, s.log))

	// --------------------------------------FILES------------------------------------------------------------------------------- //
	protected.Post("/saveImage", proxy.GatewayProxyRequest(consts.SERVER_file_system, s.log))

	// --------------------------------------CHATS------------------------------------------------------------------------------- //
	protected.Get("/getChat", proxy.GatewayProxyRequest(consts.SERVER_chat_service, s.log))
	protected.Get("/getChatBySymbol", proxy.GatewayProxyRequest(consts.SERVER_chat_service, s.log))
	protected.Get("/getAllChats", proxy.GatewayProxyRequest(consts.SERVER_chat_service, s.log))
	protected.Get("/getUserChats", proxy.GatewayProxyRequest(consts.SERVER_chat_service, s.log))
	protected.Post("/joinToChat", proxy.GatewayProxyRequest(consts.SERVER_chat_service, s.log))
	protected.Post("/createChat", proxy.GatewayProxyRequest(consts.SERVER_chat_service, s.log))
	protected.Delete("/deleteChat", proxy.GatewayProxyRequest(consts.SERVER_chat_service, s.log))
	protected.Delete("/leaveChat", proxy.GatewayProxyRequest(consts.SERVER_chat_service, s.log))
	protected.Post("/addMember", proxy.GatewayProxyRequest(consts.SERVER_chat_service, s.log))
	protected.Delete("/deleteMember", proxy.GatewayProxyRequest(consts.SERVER_chat_service, s.log))
	protected.Patch("/changeChatName", proxy.GatewayProxyRequest(consts.SERVER_chat_service, s.log))
	protected.Patch("/setTag", proxy.GatewayProxyRequest(consts.SERVER_chat_service, s.log))
	protected.Get("/getTags", proxy.GatewayProxyRequest(consts.SERVER_chat_service, s.log))
	protected.Delete("/deleteTags", proxy.GatewayProxyRequest(consts.SERVER_chat_service, s.log))
}
