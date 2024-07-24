package utils

import (
	"chat_service/internal/consts"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func CheckToken(token string) (string, error) {
	if token == "" {
		log.Println("token not found")
		return "", fmt.Errorf("token not found in header")
	}

	tokenString := strings.TrimPrefix(token, "Bearer ")
	log.Printf("token in method in chat service: %s", tokenString)
	if tokenString == "" {
		log.Println("Invalid Authorization header format")
		return "", fmt.Errorf("invalid Authorization header format")
	}

	body, statusCodeToken, err := ProxyRequest(consts.GET_Method, fmt.Sprintf(consts.US_CheckToken, tokenString), nil, http.StatusOK)
	log.Printf("body after proxyreq: %s", body)
	log.Printf("proxy after string method: %s", err)
	if err != nil {
		log.Printf("error sending or nil email: %v", err)
		return "", fmt.Errorf("error sending or nil email")
	}
	log.Println(body)
	if statusCodeToken != http.StatusOK {
		log.Println("chat not found")
		return "", fmt.Errorf("status code invalid")
	}
	log.Printf("body before trim: %s", string(body))
	email := string(body)
	email = strings.TrimSpace(email)
	email = strings.Trim(email, `"`)
	log.Printf("email after trim: %s", email)
	return email, nil
}
