package server

import (
	"net/http"

	"InterestingChats/backend/api_gateway/internal/consts"
	"InterestingChats/backend/api_gateway/internal/logger"
	"InterestingChats/backend/api_gateway/internal/proxy"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

func AuthMiddleware(log logger.Logger) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		token := ctx.Get("Authorization")
		if token == "" {
			return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error": "user unauthorized",
			})
		}

		req := &fasthttp.Request{}
		req.Header.SetMethod("POST")
		req.SetRequestURI(consts.SERVER_auth_service)
		req.Header.Set("Authorization", token)

		userID, err := proxy.AuthenticateUser(req, nil, log)
		if err != nil {
			log.Errorf("Authentication failed: %v", err)
			return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error": "authentication failed",
			})
		}

		ctx.Locals("X-User-ID", userID)
		return ctx.Next()
	}
}
