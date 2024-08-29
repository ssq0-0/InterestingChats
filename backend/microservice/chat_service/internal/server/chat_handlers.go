package server

import (
	"chat_service/internal/consts"
	"chat_service/internal/models"
	services "chat_service/internal/services/authService"
	chatservice "chat_service/internal/services/chatService"
	userservice "chat_service/internal/services/userService"
	"chat_service/internal/utils"
	"fmt"
	"net/http"
)

func (s *Server) GetAllChats(w http.ResponseWriter, r *http.Request) {
	repsonse, statusCode, clientErr, err := chatservice.GetAllChats(r)
	if err != nil {
		chatservice.ErrorHandler(w, statusCode, s.log, []string{clientErr}, err.Error())
		return
	}

	chatservice.SendRespond(w, http.StatusOK, s.log, repsonse)
}

func (s *Server) GetUserChats(w http.ResponseWriter, r *http.Request) {
	userID, statusCode, clientErr, err := userservice.TokenVerification(r.Header.Get("Authorization"))
	if err != nil || userID == 0 {
		chatservice.ErrorHandler(w, statusCode, s.log, []string{clientErr}, err.Error())
		return
	}

	userChatLists, statusCode, clientErr, err := chatservice.UserChats(userID)
	if err != nil {
		chatservice.ErrorHandler(w, statusCode, s.log, []string{clientErr}, err.Error())
		return
	}

	chatservice.SendRespond(w, http.StatusOK, s.log, userChatLists)
}

func (s *Server) GetChatHistory(w http.ResponseWriter, r *http.Request) {
	chatInfo, statusCode, clientErr, err := chatservice.GetChatHistory(r)
	if err != nil {
		chatservice.ErrorHandler(w, statusCode, s.log, []string{clientErr}, err.Error())
		return
	}

	chatservice.SendRespond(w, http.StatusOK, s.log, &models.Response{
		Data:   chatInfo,
		Errors: nil,
	})
}

func (s *Server) SearchChat(w http.ResponseWriter, r *http.Request) {
	userID, statusCode, clientErr, err := userservice.TokenVerification(r.Header.Get("Authorization"))
	if err != nil || userID == 0 {
		chatservice.ErrorHandler(w, statusCode, s.log, []string{clientErr}, err.Error())
		return
	}

	response, clientErr, err := chatservice.GetQueryChat(r)
	if err != nil {
		chatservice.ErrorHandler(w, http.StatusBadRequest, s.log, []string{clientErr}, err.Error())
		return
	}

	chatservice.SendRespond(w, http.StatusOK, s.log, response)
}

func (s *Server) CreateChat(w http.ResponseWriter, r *http.Request) {
	userID, statusCode, clientErr, err := userservice.TokenVerification(r.Header.Get("Authorization"))
	if err != nil || userID == 0 {
		chatservice.ErrorHandler(w, statusCode, s.log, []string{clientErr}, err.Error())
		return
	}

	chatInfo, clientErr, err := chatservice.Validation(r, userID)
	if err != nil {
		chatservice.ErrorHandler(w, http.StatusBadRequest, s.log, []string{clientErr}, err.Error())
		return
	}

	body, statusCode, clientErr, err := utils.ProxyRequest(consts.POST_Method, consts.DB_CreateChat, chatInfo, http.StatusCreated)
	if err != nil {
		chatservice.ErrorHandler(w, statusCode, s.log, []string{clientErr}, err.Error())
		return
	}

	var chat models.GetAndCreateChat
	if err := utils.ParseBody(body, &chat); err != nil {
		chatservice.ErrorHandler(w, http.StatusBadRequest, s.log, []string{consts.ErrInternalServer}, err.Error())
		return
	}

	chatservice.SendRespond(w, http.StatusCreated, s.log, &models.Response{
		Errors: nil,
		Data:   chat,
	})
}

func (s *Server) DeleteChat(w http.ResponseWriter, r *http.Request) {
	userID, statusCode, clientErr, err := userservice.TokenVerification(r.Header.Get("Authorization"))
	if err != nil || userID == 0 {
		chatservice.ErrorHandler(w, statusCode, s.log, []string{clientErr}, err.Error())
		return
	}

	chatID, err := chatservice.GetChatID(r)
	if err != nil {
		chatservice.ErrorHandler(w, http.StatusBadRequest, s.log, []string{consts.ErrChatNotFound}, err.Error())
		return
	}

	// 1
	if _, authorStatusCode, clientErr, err := utils.ProxyRequest(consts.GET_Method, fmt.Sprintf(consts.DB_GetAuthor, userID, chatID), nil, http.StatusOK); err != nil {
		chatservice.ErrorHandler(w, authorStatusCode, s.log, []string{clientErr}, err.Error())
		return
	}

	// 2
	if _, deleteStatusCode, clientErr, err := utils.ProxyRequest(consts.DELETE_Method, fmt.Sprintf(consts.DB_DeleteChat, chatID), nil, http.StatusNoContent); err != nil {
		chatservice.ErrorHandler(w, deleteStatusCode, s.log, []string{clientErr}, err.Error())
		return
	}

	// TODO: change data to 'response' by 1 & 2
	chatservice.SendRespond(w, http.StatusOK, s.log, &models.Response{
		Data:   "successful delete",
		Errors: nil,
	})
}

func (s *Server) AddMember(w http.ResponseWriter, r *http.Request) {
	memberList, statusCode, clientErr, err := services.AuthenticateAndAuthorize(r)
	if err != nil {
		chatservice.ErrorHandler(w, statusCode, s.log, []string{clientErr}, err.Error())
		return
	}

	response := &models.Response{}
	for _, member := range memberList {
		if clientErr, _, err := userservice.UserExists(member.UserID); err != nil {
			response.Errors = append(response.Errors, clientErr)
			continue
		}

		if clientErr, err = userservice.ManageMember(consts.POST_Method, consts.DB_AddMembers, member, http.StatusAccepted); err != nil {
			response.Errors = append(response.Errors, clientErr)
		}
	}

	chatservice.SendRespond(w, http.StatusOK, s.log, response)
}

func (s *Server) DeleteMember(w http.ResponseWriter, r *http.Request) {
	memberList, statusCode, clientErr, err := services.AuthenticateAndAuthorize(r)
	if err != nil {
		chatservice.ErrorHandler(w, statusCode, s.log, []string{clientErr}, err.Error())
		return
	}

	response := &models.Response{}
	for _, member := range memberList {
		if _, _, clientErr, err := utils.ProxyRequest(consts.GET_Method, fmt.Sprintf(consts.DB_CheckUser, member.UserID), nil, http.StatusOK); err != nil {
			response.Errors = append(response.Errors, clientErr)
			continue
		}

		if clientErr, err = userservice.ManageMember(consts.DELETE_Method, consts.DB_DeleteMember, member, http.StatusNoContent); err != nil {
			response.Errors = append(response.Errors, clientErr)
		}
	}

	chatservice.SendRespond(w, http.StatusOK, s.log, response)
}

func (s *Server) JoinToChat(w http.ResponseWriter, r *http.Request) {
	response, statusCode, clientErr, err := chatservice.JoinToChat(r)
	if err != nil {
		chatservice.ErrorHandler(w, statusCode, s.log, []string{clientErr}, err.Error())
		return
	}

	chatservice.SendRespond(w, http.StatusOK, s.log, response)
}
