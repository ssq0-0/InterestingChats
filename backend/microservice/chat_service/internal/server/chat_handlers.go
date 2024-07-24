package server

import (
	"chat_service/internal/consts"
	"chat_service/internal/models"
	"chat_service/internal/utils"
	"strconv"

	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func (s *Server) GetChatHistory(w http.ResponseWriter, r *http.Request) {
	chatName := r.URL.Query().Get("chatName")
	if chatName == "" {
		http.Error(w, "chatName is required", http.StatusBadRequest)
		return
	}

	body, statusCode, err := utils.ProxyRequest(consts.GET_Method, fmt.Sprintf(consts.DB_GetChatHistory, chatName), nil, http.StatusOK)
	if err != nil {
		log.Println("failed get info about chat", err)
		http.Error(w, "failed get info about chat", http.StatusBadRequest)
		return
	}
	if statusCode == http.StatusNotFound {
		log.Println("chat not found")
		http.Error(w, "chat not found", http.StatusBadRequest)
		return
	}

	chat := &models.Chat{}
	if err := json.Unmarshal(body, chat); err != nil {
		log.Println("failed to deserialize chat data", err)
		http.Error(w, "failed to deserialize chat data", http.StatusInternalServerError)
		return
	}

	log.Printf("Chat data: %+v\n", chat)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(chat); err != nil {
		http.Error(w, "error encode json", http.StatusInternalServerError)
		log.Println("error encoding json for response", err)
	}
}

func (s *Server) CreateChat(w http.ResponseWriter, r *http.Request) {
	authToken := r.Header.Get("access_token")
	// use email for "author chat"
	_, err := utils.CheckToken(authToken)
	if err != nil {
		log.Printf("access rejected: %v", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	chatInfo := &models.Chat{}
	if err := json.NewDecoder(r.Body).Decode(&chatInfo); err != nil {
		log.Printf("error decoding request body: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body, statusCode, err := utils.ProxyRequest(consts.POST_Method, consts.DB_CreateChat, chatInfo, http.StatusCreated)
	if err != nil {
		log.Printf("failed to create chat: %v. StatusCode: %d", err, statusCode)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.Unmarshal(body, chatInfo); err != nil {
		log.Printf("failed to deserialize chat data:%v. Chat data: %+v", err, body)
		http.Error(w, "failed to deserialize chat data", http.StatusInternalServerError)
		return
	}

	log.Println(chatInfo)
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(chatInfo); err != nil {
		log.Println("error encoding response:", err)
		http.Error(w, "error encoding response", http.StatusInternalServerError)
	}
}

func (s *Server) DeleteChat(w http.ResponseWriter, r *http.Request) {
	authToken := r.Header.Get("Authorization")
	email, err := utils.CheckToken(authToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		log.Printf("error check token: %v", err)
		return
	}

	chatIDstr := r.URL.Query().Get("chatID")
	if chatIDstr == "" {
		http.Error(w, "missing chat id", http.StatusBadRequest)
		log.Printf("missing chat id")
		return
	}

	chatIDint, err := strconv.Atoi(chatIDstr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("error convert to int: %v", err)
		return
	}

	_, authorStatusCode, err := utils.ProxyRequest(consts.GET_Method, fmt.Sprintf(consts.DB_GetAuthor, email, chatIDint), nil, http.StatusOK)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		log.Printf("access foribem: %v, status code: %d", err, authorStatusCode)
		return
	}

	_, deleteStatusCode, err := utils.ProxyRequest(consts.DELETE_Method, fmt.Sprintf(consts.DB_DeleteChat, chatIDstr), nil, http.StatusNoContent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("failed request: %v, statuscode: %d", err, deleteStatusCode)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) AddMember(w http.ResponseWriter, r *http.Request) {
	authToken := r.Header.Get("Authorization")
	email, err := utils.CheckToken(authToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		log.Printf("error check token: %v", err)
		return
	}

	var addMembers []models.AddMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&addMembers); err != nil {
		log.Printf("failed decode request: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if len(addMembers) == 0 {
		http.Error(w, "no members to add", http.StatusBadRequest)
		return
	}

	chatID := addMembers[0].ChatID
	_, authorStatusCode, err := utils.ProxyRequest(consts.GET_Method, fmt.Sprintf(consts.DB_GetAuthor, email, chatID), nil, http.StatusOK)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		log.Printf("access foribem: %v, status code: %d", err, authorStatusCode)
		return
	}

	var errors []string
	for _, member := range addMembers {
		_, _, err := utils.ProxyRequest(consts.GET_Method, fmt.Sprintf(consts.DB_CheckUser, member.UserID), nil, http.StatusOK)
		if err != nil {
			log.Printf("failed request: %v", err)
			errors = append(errors, fmt.Sprintf("failed to check user id: %v", member.UserID))
		}

		_, _, err = utils.ProxyRequest(consts.POST_Method, consts.DB_AddMembers, member, http.StatusAccepted)
		if err != nil {
			log.Printf("failed request: %v", err)
			errors = append(errors, fmt.Sprintf("failed to add user %d in chat %d", member.UserID, member.ChatID))
		}
	}

	if len(errors) > 0 {
		w.WriteHeader(http.StatusPartialContent)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	response := map[string]interface{}{
		"errors": errors,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "error encoding json", http.StatusInternalServerError)
		log.Println("error encoding json for response", err)
	}
}

func (s *Server) DeleteMember(w http.ResponseWriter, r *http.Request) {
	authToken := r.Header.Get("Authorization")
	email, err := utils.CheckToken(authToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		log.Printf("error check token: %v", err)
		return
	}

	var deleteMembers []models.AddMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&deleteMembers); err != nil {
		log.Printf("error decoding r body to json: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if len(deleteMembers) == 0 {
		http.Error(w, "no members to delete", http.StatusBadRequest)
		return
	}

	chatID := deleteMembers[0].ChatID
	_, authorStatusCode, err := utils.ProxyRequest(consts.GET_Method, fmt.Sprintf(consts.DB_GetAuthor, email, chatID), nil, http.StatusOK)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		log.Printf("access foribem: %v, status code: %d", err, authorStatusCode)
		return
	}

	var errors []string
	for _, member := range deleteMembers {
		_, _, err := utils.ProxyRequest(consts.GET_Method, fmt.Sprintf(consts.DB_CheckUser, member.UserID), nil, http.StatusOK)
		if err != nil {
			log.Printf("failed request: %v", err)
			errors = append(errors, fmt.Sprintf("failed to check user id: %v", member.UserID))
		}

		_, _, err = utils.ProxyRequest(consts.DELETE_Method, consts.DB_DeleteMember, member, http.StatusNoContent)
		if err != nil {
			log.Printf("failed request: %v", err)
			errors = append(errors, fmt.Sprintf("failed to add user %d in chat %d", member.UserID, member.ChatID))
		}
	}

	if len(errors) > 0 {
		w.WriteHeader(http.StatusPartialContent)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	response := map[string]interface{}{
		"errors": errors,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "error encoding json", http.StatusInternalServerError)
		log.Println("error encoding json for response", err)
	}
}
