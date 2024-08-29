package chatservice

import (
	"chat_service/internal/consts"
	"chat_service/internal/models"
	userservice "chat_service/internal/services/userService"
	"chat_service/internal/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func Validation(r *http.Request, creator int) (*models.GetAndCreateChat, string, error) {
	var chatInfo *models.GetAndCreateChat
	if err := json.NewDecoder(r.Body).Decode(&chatInfo); err != nil {
		log.Println(err)
		return nil, consts.ErrInternalServer, fmt.Errorf(consts.InternalInvalidValueFormat, err)
	}

	if chatInfo.ChatName == "" || chatInfo.Creator == 0 || len(chatInfo.Members) == 0 {
		return nil, consts.ErrMissingChatInfo, fmt.Errorf(consts.InternalMissingParametr)
	}

	// Create a regular expression
	for _, userID := range chatInfo.Members {
		if userID == 0 {
			return nil, consts.ErrIncompleteUserData, fmt.Errorf(consts.InternalIncompleteUserData, userID)
		}
	}

	return chatInfo, "", nil
}

func GetChatID(r *http.Request) (int, error) {
	chatIDstr := r.URL.Query().Get("chatID")
	if chatIDstr == "" {
		return 0, fmt.Errorf("missing chat id")
	}

	chatIDint, err := strconv.Atoi(chatIDstr)
	if err != nil {
		return 0, fmt.Errorf("error convert to int: %v", err)
	}

	return chatIDint, nil
}

func GetChatHistory(r *http.Request) (*models.Chat, int, string, error) {
	chatID := r.URL.Query().Get("chatID")
	if chatID == "" {
		return nil, http.StatusInternalServerError, consts.ErrMissingParametr, fmt.Errorf(consts.InternalMissingParametr)
	}

	// TODO: Status code!!! √
	userID, statusCode, clientErr, err := userservice.TokenVerification(r.URL.Query().Get("Authorization"))
	if err != nil {
		return nil, statusCode, clientErr, fmt.Errorf(consts.InternalErrUserUnauthorized)
	}

	response, statusCode, clientErr, err := utils.ProxyRequest(consts.GET_Method, fmt.Sprintf(consts.DB_GetChatHistory, chatID, userID), nil, http.StatusOK)
	if err != nil {
		return nil, statusCode, clientErr, err
	}

	var chat models.Chat
	if err := json.Unmarshal(response, &chat); err != nil {
		return nil, http.StatusInternalServerError, consts.ErrInternalServer, err
	}

	return &chat, http.StatusOK, "", nil
}

func GetQueryChat(r *http.Request) (*models.Response, string, error) {
	symbols := r.URL.Query().Get("chatName")
	if strings.ReplaceAll(symbols, " ", "") == "" {
		return nil, consts.ErrMissingChatInfo, fmt.Errorf(consts.InternalMissingParametr)
	}

	body, _, clientErr, err := utils.ProxyRequest(consts.GET_Method, fmt.Sprintf(consts.DB_SearchQueryChat, symbols), nil, http.StatusOK)
	if err != nil {
		return nil, clientErr, err
	}

	var chats models.Response
	if err := json.Unmarshal(body, &chats); err != nil {
		return nil, consts.ErrInternalServer, err
	}

	return &chats, "", nil
}

func GetAllChats(r *http.Request) (*models.Response, int, string, error) {
	_, statusCode, clientErr, err := userservice.TokenVerification(r.Header.Get("Authorization"))
	if err != nil {
		log.Println("here")
		return nil, statusCode, clientErr, fmt.Errorf(consts.InternalErrUserUnauthorized)
	}

	body, statusCode, clientErr, err := utils.ProxyRequest(consts.GET_Method, consts.DB_GetAllChats, nil, http.StatusOK)
	if err != nil {
		return nil, statusCode, clientErr, err
	}

	var chats models.Response
	if err := json.Unmarshal(body, &chats); err != nil {
		return nil, http.StatusInternalServerError, consts.ErrInternalServer, err
	}

	return &chats, http.StatusOK, "", nil
}

func UserChats(userID int) (*models.Response, int, string, error) {
	response, statusCode, clientErr, err := utils.ProxyRequest(consts.GET_Method, fmt.Sprintf(consts.DB_GetUserChats, userID), nil, http.StatusOK)
	if err != nil {
		return nil, statusCode, clientErr, err
	}

	var chatLists *models.Response
	if err := json.Unmarshal(response, &chatLists); err != nil {
		log.Print(err)

		return nil, http.StatusInternalServerError, consts.ErrInternalServer, err
	}

	return chatLists, http.StatusOK, "", nil
}

func JoinToChat(r *http.Request) (*models.Response, int, string, error) {
	userID, statusCode, clientErr, err := userservice.TokenVerification(r.Header.Get("Authorization"))
	if err != nil || userID == 0 {
		return nil, statusCode, clientErr, err
	}

	var chatID int
	if clientErr, err := utils.ConvertValue(r.URL.Query().Get("chatID"), &chatID); err != nil {
		return nil, http.StatusInternalServerError, clientErr, err
	}

	response := &models.Response{}
	if clientErr, err = userservice.ManageMember(consts.POST_Method, consts.DB_AddMembers, models.MemberRequest{
		UserID: userID,
		ChatID: chatID,
	}, http.StatusAccepted); err != nil {
		response.Errors = append(response.Errors, clientErr)
	}

	return response, http.StatusOK, "", nil
}
