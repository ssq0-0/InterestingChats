package proxy

import (
	"InterestingChats/backend/api_gateway/internal/consts"
	"InterestingChats/backend/api_gateway/internal/logger"
	"InterestingChats/backend/api_gateway/internal/models"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

func GatewayProxyRequest(target string, log logger.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return proxyRequest(target, c, log)
	}
}

func proxyRequest(target string, c *fiber.Ctx, log logger.Logger) error {
	req := c.Request()
	log.Infof("Proxying request to: %s", target+string(req.RequestURI()))

	proxyReq, err := http.NewRequest(string(req.Header.Method()), target+string(req.RequestURI()), bytes.NewReader(req.Body()))
	if err != nil {
		log.Errorf("Failed to create proxy request: %v", err)
		return c.Status(fiber.StatusBadRequest).SendString("Failed to create proxy request")
	}

	req.Header.VisitAll(func(key, value []byte) {
		proxyReq.Header.Set(string(key), string(value))
	})

	if userID := c.Locals("X-User-ID"); userID != "" {
		proxyReq.Header.Set("X-User-ID", fmt.Sprintf("%d", userID))
		log.Infof("Added X-User-ID to proxy request: %s", userID)
	} else {
		log.Warn("X-User-ID not set in context")
	}

	resp, err := http.DefaultClient.Do(proxyReq)
	if err != nil {
		log.Errorf("Failed to proxy request: %v", err)
		return c.Status(fiber.StatusBadRequest).SendString("Failed to proxy request")
	}
	defer resp.Body.Close()

	for key, values := range resp.Header {
		for _, value := range values {
			c.Set(key, value)
		}
	}

	c.Status(resp.StatusCode)
	_, err = io.Copy(c.Response().BodyWriter(), resp.Body)
	if err != nil {
		log.Errorf("Failed to copy response body: %v", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to copy response body")
	}
	return nil
}

func AuthenticateUser(req *fasthttp.Request, body []byte, log logger.Logger) (int, error) {
	authReq, err := http.NewRequest("POST", consts.SERVER_auth_service+"/auth", bytes.NewBuffer(body))
	if err != nil {
		return 0, fmt.Errorf("failed to create auth request: %w", err)
	}

	req.Header.VisitAll(func(key, value []byte) {
		authReq.Header.Set(string(key), string(value))
	})

	client := &http.Client{}
	resp, err := client.Do(authReq)
	if err != nil {
		return 0, fmt.Errorf("auth request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("authentication failed, status code: %d", resp.StatusCode)
	}

	var response models.Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return 0, fmt.Errorf("failed to decode response: %w", err)
	}

	userIDFloat, ok := response.Data.(float64)
	if !ok {
		return 0, fmt.Errorf("invalid user ID format")
	}

	userID := int(userIDFloat)
	if userID == 0 {
		return 0, fmt.Errorf("invalid user ID")
	}
	log.Infof("Authenticated user with ID: %d", userID)
	return userID, nil
}
