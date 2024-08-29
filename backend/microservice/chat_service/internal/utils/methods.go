package utils

import (
	"chat_service/internal/consts"
	"chat_service/internal/models"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func CheckToken(token string) (int, int, string, error) {
	if token == "" {
		return 0, http.StatusInternalServerError, consts.ErrTokenHeader, fmt.Errorf(consts.InternalTokenHeader)
	}

	tokenString := strings.TrimPrefix(token, "Bearer ")
	if tokenString == "" {
		return 0, http.StatusUnauthorized, consts.ErrTokenHeader, fmt.Errorf(consts.InternalTokenInvalidFormat)
	}

	body, statusCode, clientErr, err := ProxyRequest(consts.GET_Method, fmt.Sprintf(consts.US_CheckToken, tokenString), nil, http.StatusOK)
	if err != nil {
		return 0, statusCode, clientErr, err
	}

	var response models.Response
	if err := json.Unmarshal(body, &response); err != nil {
		return 0, http.StatusInternalServerError, consts.ErrInvalidToken, fmt.Errorf("error unmarshaling response: %v", err)
	}

	userID, ok := response.Data.(float64)
	if !ok {
		return 0, http.StatusInternalServerError, consts.ErrInvalidToken, fmt.Errorf("error asserting response.Data to float64")
	}

	return int(userID), http.StatusOK, "", nil
}
