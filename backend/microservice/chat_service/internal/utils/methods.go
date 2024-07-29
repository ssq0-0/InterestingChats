package utils

import (
	"chat_service/internal/consts"
	"chat_service/internal/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func CheckToken(token string) (int, error) {
	if token == "" {
		return 0, fmt.Errorf("token not found in header")
	}

	tokenString := strings.TrimPrefix(token, "Bearer ")
	if tokenString == "" {
		return 0, fmt.Errorf("invalid Authorization header format")
	}

	body, _, err := ProxyRequest(consts.GET_Method, fmt.Sprintf(consts.US_CheckToken, tokenString), nil, http.StatusOK)
	if err != nil {
		return 0, fmt.Errorf("error sending token request: %v", err)
	}

	var response models.Response
	if err := json.Unmarshal(body, &response); err != nil {
		return 0, fmt.Errorf("failed decode body response: %v", err)
	}
	userID, ok := response.Data.(float64)
	if !ok {
		log.Println("failed to assert response data to int")
		return 0, fmt.Errorf("failed to assert response data to int")
	}

	return int(userID), nil
}
