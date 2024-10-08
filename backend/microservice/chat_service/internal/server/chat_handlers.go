package server

import (
	"chat_service/internal/consts"
	"chat_service/internal/models"
	chatservice "chat_service/internal/services/chatService"
	userservice "chat_service/internal/services/userService"
	"chat_service/internal/utils"
	"fmt"
	"net/http"
)

// GetAllChats retrieves all chats and sends the response to the client.
func (s *Server) GetAllChats(w http.ResponseWriter, r *http.Request) {
	repsonse, statusCode, clientErr, err := chatservice.GetAllChats(r)
	if err != nil {
		chatservice.ErrorHandler(w, statusCode, s.log, []string{clientErr}, err.Error())
		return
	}

	chatservice.SendRespond(w, http.StatusOK, s.log, repsonse)
}

// GetUserChats retrieves chats for the authenticated user and sends the response to the client.
func (s *Server) GetUserChats(w http.ResponseWriter, r *http.Request) {
	userChatLists, statusCode, clientErr, err := chatservice.UserChats(r)
	if err != nil {
		chatservice.ErrorHandler(w, statusCode, s.log, []string{clientErr}, err.Error())
		return
	}

	chatservice.SendRespond(w, http.StatusOK, s.log, userChatLists)
}

// GetChatHistory retrieves the chat history for a specific chat and sends the response to the client.
func (s *Server) GetChatHistory(w http.ResponseWriter, r *http.Request) {
	chatInfo, statusCode, clientErr, err := chatservice.GetChatHistory(r, 0)
	if err != nil {
		chatservice.ErrorHandler(w, statusCode, s.log, []string{clientErr}, err.Error())
		return
	}

	chatservice.SendRespond(w, http.StatusOK, s.log, &models.Response{
		Errors: nil,
		Data:   chatInfo,
	})
}

// SearchChat performs a search for chats based on a query and sends the response to the client.
func (s *Server) SearchChat(w http.ResponseWriter, r *http.Request) {
	response, statusCode, clientErr, err := chatservice.GetQueryChat(r)
	if err != nil {
		chatservice.ErrorHandler(w, statusCode, s.log, []string{clientErr}, err.Error())
		return
	}

	chatservice.SendRespond(w, http.StatusOK, s.log, response)
}

// CreateChat creates a new chat with the provided information and sends the response to the client.
func (s *Server) CreateChat(w http.ResponseWriter, r *http.Request) {
	chatInfo, statusCode, clientErr, err := chatservice.Validation(r)
	if err != nil {
		chatservice.ErrorHandler(w, statusCode, s.log, []string{clientErr}, err.Error())
		return
	}

	response, statusCode, clientErr, err := utils.ProxyRequest(consts.POST_Method, consts.DB_CreateChat, chatInfo, http.StatusCreated)
	if err != nil {
		chatservice.ErrorHandler(w, statusCode, s.log, []string{clientErr}, err.Error())
		return
	}

	chatservice.SendRespond(w, http.StatusCreated, s.log, response)
}

// DeleteChat deletes a chat identified by its ID and sends the response to the client.
func (s *Server) DeleteChat(w http.ResponseWriter, r *http.Request) {
	_, chatID, statusCode, clientErr, err := chatservice.CheckAuthor(r)
	if err != nil {
		chatservice.ErrorHandler(w, statusCode, s.log, []string{clientErr}, err.Error())
		return
	}

	if _, deleteStatusCode, clientErr, err := utils.ProxyRequest(consts.DELETE_Method, fmt.Sprintf(consts.DB_DeleteChat, chatID), nil, http.StatusNoContent); err != nil {
		chatservice.ErrorHandler(w, deleteStatusCode, s.log, []string{clientErr}, err.Error())
		return
	}

	chatservice.SendRespond(w, http.StatusOK, s.log, &models.Response{
		Data:   "successful delete",
		Errors: nil,
	})
}

// ChangeChatName changes the name of a chat identified by its ID and sends the response to the client.
func (s *Server) ChangeChatName(w http.ResponseWriter, r *http.Request) {
	_, chatID, statusCode, clientErr, err := chatservice.CheckAuthor(r)
	if err != nil {
		chatservice.ErrorHandler(w, statusCode, s.log, []string{clientErr}, err.Error())
		return
	}

	if _, statusCode, clientErr, err := utils.ProxyRequest(consts.PATCH_Method, fmt.Sprintf(consts.DB_ChangeChatName, chatID, r.URL.Query().Get("chatName")), nil, http.StatusOK); err != nil {
		chatservice.ErrorHandler(w, statusCode, s.log, []string{clientErr}, err.Error())
		return
	}

	chatservice.SendRespond(w, http.StatusOK, s.log, &models.Response{
		Data:   "successful change name",
		Errors: nil,
	})
}

// AddMember adds members to a chat and sends the response to the client.
func (s *Server) AddMember(w http.ResponseWriter, r *http.Request) {
	if _, _, statusCode, clientErr, err := chatservice.CheckAuthor(r); err != nil {
		chatservice.ErrorHandler(w, statusCode, s.log, []string{clientErr}, err.Error())
		return
	}

	memberList, statusCode, clientErr, err := chatservice.MemberList(r)
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

// DeleteMember removes members from a chat and sends the response to the client.
func (s *Server) DeleteMember(w http.ResponseWriter, r *http.Request) {
	if _, _, statusCode, clientErr, err := chatservice.CheckAuthor(r); err != nil {
		chatservice.ErrorHandler(w, statusCode, s.log, []string{clientErr}, err.Error())
		return
	}

	memberList, statusCode, clientErr, err := chatservice.MemberList(r)
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

// LeaveChat allows a user to leave a chat and sends the response to the client.
func (s *Server) LeaveChat(w http.ResponseWriter, r *http.Request) {
	memberList, statusCode, clientErr, err := chatservice.MemberList(r)
	if err != nil {
		chatservice.ErrorHandler(w, statusCode, s.log, []string{clientErr}, err.Error())
		return
	}

	for _, member := range memberList {
		if _, statusCode, clientErr, err := utils.ProxyRequest(consts.GET_Method, fmt.Sprintf(consts.DB_CheckUser, member.UserID), nil, http.StatusOK); err != nil {
			chatservice.ErrorHandler(w, statusCode, s.log, []string{clientErr}, err.Error())
			return
		}

		if clientErr, err = userservice.ManageMember(consts.DELETE_Method, consts.DB_DeleteMember, member, http.StatusNoContent); err != nil {
			chatservice.ErrorHandler(w, http.StatusBadRequest, s.log, []string{clientErr}, err.Error())
			return
		}
	}

	chatservice.SendRespond(w, http.StatusOK, s.log, &models.Response{
		Errors: nil,
		Data:   "successful leave from chat!",
	})
}

// SetTag assigns a tag to a chat and sends the response to the client.
func (s *Server) SetTag(w http.ResponseWriter, r *http.Request) {
	response, statusCode, clientErr, err := chatservice.SetTagToChat(r)
	if err != nil {
		chatservice.ErrorHandler(w, statusCode, s.log, []string{clientErr}, err.Error())
		return
	}

	chatservice.SendRespond(w, http.StatusOK, s.log, response)
}

// GetTags retrieves tags associated with chats and sends the response to the client.
func (s *Server) GetTags(w http.ResponseWriter, r *http.Request) {
	tags, statusCode, clientErr, err := chatservice.GetTags(r)
	if err != nil {
		chatservice.ErrorHandler(w, statusCode, s.log, []string{clientErr}, err.Error())
		return
	}

	chatservice.SendRespond(w, http.StatusOK, s.log, tags)
}

// DeleteTags removes specified tags from a chat and sends the response to the client.
func (s *Server) DeleteTags(w http.ResponseWriter, r *http.Request) {
	response, statusCode, clientErr, err := chatservice.DeleteTags(r)
	if err != nil {
		chatservice.ErrorHandler(w, statusCode, s.log, []string{clientErr}, err.Error())
		return
	}

	chatservice.SendRespond(w, http.StatusOK, s.log, response)
}

// JoinToChat allows a user to join a chat and sends the response to the client.
func (s *Server) JoinToChat(w http.ResponseWriter, r *http.Request) {
	response, statusCode, clientErr, err := chatservice.JoinToChat(r)
	if err != nil {
		chatservice.ErrorHandler(w, statusCode, s.log, []string{clientErr}, err.Error())
		return
	}

	chatservice.SendRespond(w, http.StatusOK, s.log, response)
}
