package server

import (
	"chat_service/internal/consts"
	"chat_service/internal/models"
	services "chat_service/internal/services/authService"
	chatservice "chat_service/internal/services/chatService"
	userservice "chat_service/internal/services/userService"
	"chat_service/internal/utils"

	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

func (s *Server) GetChatHistory(w http.ResponseWriter, r *http.Request) {
	chatName := r.URL.Query().Get("chatName")
	if chatName == "" {
		ErrorHandler(w, http.StatusBadRequest, []string{"chatName is required"}, "chatName is required")
		return
	}

	encodedChatName := url.QueryEscape(chatName)
	body, _, err := utils.ProxyRequest(consts.GET_Method, fmt.Sprintf(consts.DB_GetChatHistory, encodedChatName), nil, http.StatusOK)
	if err != nil {
		ErrorHandler(w, http.StatusBadRequest, []string{"failed get info about chat"}, fmt.Sprintf("failed get info about chat: %v", err))
		return
	}

	var chatInfo models.Chat
	if err := utils.ParseBody(body, &chatInfo); err != nil {
		ErrorHandler(w, http.StatusBadRequest, []string{"failed to deserialize chat data"}, fmt.Sprintf("failed to deserialize chat data: %v", err))
		return
	}

	SendRespond(w, http.StatusOK, &models.Response{
		Data:   chatInfo,
		Errors: nil,
	})
}

func (s *Server) CreateChat(w http.ResponseWriter, r *http.Request) {
	userID, err := userservice.UserVerification(r)
	if err != nil || userID == 0 {
		ErrorHandler(w, http.StatusBadRequest, []string{"access rejected"}, fmt.Sprintf("access rejected: %v", err))
		return
	}

	chatInfo := &models.Chat{Creator: userID}
	if err := json.NewDecoder(r.Body).Decode(&chatInfo); err != nil {
		ErrorHandler(w, http.StatusBadRequest, []string{"error decoding request body"}, fmt.Sprintf("error decoding request body: %v", err))
		return
	}

	body, statusCode, err := utils.ProxyRequest(consts.POST_Method, consts.DB_CreateChat, chatInfo, http.StatusCreated)
	if err != nil {
		ErrorHandler(w, statusCode, []string{"failed to create chat"}, fmt.Sprintf("failed to create chat: %v. StatusCode: %d", err, statusCode))
		return
	}

	if err := utils.ParseBody(body, &chatInfo); err != nil {
		ErrorHandler(w, http.StatusBadRequest, []string{"failed to deserialize chat data"}, fmt.Sprintf("failed to deserialize chat data: %v", err))
		return
	}

	SendRespond(w, http.StatusCreated, &models.Response{
		Data:   chatInfo,
		Errors: nil,
	})
}

func (s *Server) DeleteChat(w http.ResponseWriter, r *http.Request) {
	userID, err := userservice.UserVerification(r)
	if err != nil || userID == 0 {
		ErrorHandler(w, http.StatusBadRequest, []string{"access rejected"}, fmt.Sprintf("access rejected: %v", err))
		return
	}

	chatID, err := chatservice.GetChatID(r)
	if err != nil {
		ErrorHandler(w, http.StatusBadRequest, []string{"failed get info about chat"}, fmt.Sprintf("failed get info about chat: %v", err))
		return
	}

	_, authorStatusCode, err := utils.ProxyRequest(consts.GET_Method, fmt.Sprintf(consts.DB_GetAuthor, userID, chatID), nil, http.StatusOK)
	if err != nil {
		ErrorHandler(w, authorStatusCode, []string{"failed get info about author"}, fmt.Sprintf("access foribem: %v, status code: %d", err, authorStatusCode))
		return
	}

	_, deleteStatusCode, err := utils.ProxyRequest(consts.DELETE_Method, fmt.Sprintf(consts.DB_DeleteChat, chatID), nil, http.StatusNoContent)
	if err != nil {
		ErrorHandler(w, deleteStatusCode, []string{"failed delete chat"}, fmt.Sprintf("failed request: %v, statuscode: %d", err, deleteStatusCode))
		return
	}

	SendRespond(w, http.StatusNoContent, nil)
}

func (s *Server) AddMember(w http.ResponseWriter, r *http.Request) {
	memberList, err := services.AuthenticateAndAuthorize(r)
	if err != nil {
		ErrorHandler(w, http.StatusBadRequest, []string{err.Error()}, fmt.Sprintf("error in auth service: %v", err))
		return
	}

	response := &models.Response{}
	for _, member := range memberList {
		if err := userservice.UserExists(member.UserID); err != nil {
			response.Errors = append(response.Errors, err.Error())
			continue
		}

		if err = userservice.ManageMember(consts.POST_Method, consts.DB_AddMembers, member, http.StatusAccepted); err != nil {
			response.Errors = append(response.Errors, "failed")
		}
	}

	SendRespond(w, http.StatusOK, response)
}

func (s *Server) DeleteMember(w http.ResponseWriter, r *http.Request) {
	memberList, err := services.AuthenticateAndAuthorize(r)
	if err != nil {
		ErrorHandler(w, http.StatusBadRequest, []string{err.Error()}, fmt.Sprintf("error in auth service: %v", err))
		return
	}

	response := &models.Response{}
	for _, member := range memberList {
		if _, _, err := utils.ProxyRequest(consts.GET_Method, fmt.Sprintf(consts.DB_CheckUser, member.UserID), nil, http.StatusOK); err != nil {
			log.Printf("failed request: %v", err)
			response.Errors = append(response.Errors, err.Error())
			continue
		}

		if err = userservice.ManageMember(consts.DELETE_Method, consts.DB_DeleteMember, member, http.StatusNoContent); err != nil {
			log.Printf("failed delete member: %v", err)
			response.Errors = append(response.Errors, "failed")
		}
	}

	SendRespond(w, http.StatusOK, response)
}
