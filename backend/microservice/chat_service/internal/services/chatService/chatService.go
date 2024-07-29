package chatservice

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func GetChatID(r *http.Request) (int, error) {
	chatIDstr := r.URL.Query().Get("chatID")
	if chatIDstr == "" {
		log.Printf("missing chat id")
		return 0, fmt.Errorf("missing chat id")
	}

	chatIDint, err := strconv.Atoi(chatIDstr)
	if err != nil {
		log.Printf("error convert to int: %v", err)
		return 0, fmt.Errorf("error convert to int: %v", err)
	}

	return chatIDint, nil
}
