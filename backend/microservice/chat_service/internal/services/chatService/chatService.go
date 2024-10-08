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
	"strings"
)

// Validation validates the incoming request for creating or getting chat information.
// It decodes the request body into a GetAndCreateChat model and checks for required fields.
func Validation(r *http.Request) (*models.GetAndCreateChat, int, string, error) {
	var chatInfo *models.GetAndCreateChat
	if err := json.NewDecoder(r.Body).Decode(&chatInfo); err != nil {
		log.Println(err)
		return nil, http.StatusBadRequest, consts.ErrInternalServer, fmt.Errorf(consts.InternalInvalidValueFormat, err)
	}

	if clientErr, err := utils.ConvertValue(r.Header.Get("X-User-ID"), &chatInfo.Creator); err != nil {
		return nil, http.StatusUnauthorized, clientErr, err
	}

	log.Printf("chat info: %+v", chatInfo)
	if chatInfo.ChatName == "" || chatInfo.Creator == 0 || len(chatInfo.Members) == 0 {
		return nil, http.StatusBadRequest, consts.ErrMissingChatInfo, fmt.Errorf(consts.InternalMissingParametr)
	}

	// Create a regular expression
	for _, userID := range chatInfo.Members {
		if userID == 0 {
			return nil, http.StatusBadRequest, consts.ErrIncompleteUserData, fmt.Errorf(consts.InternalIncompleteUserData, userID)
		}
	}

	return chatInfo, http.StatusOK, "", nil
}

// GetChatHistory retrieves the chat history for a given chatID.
// It validates the request and calls the appropriate utility function to get the chat history.
func GetChatHistory(r *http.Request, typeGet int) (*models.Chat, int, string, error) {
	chatID := r.URL.Query().Get("chatID")
	if chatID == "" {
		return nil, http.StatusInternalServerError, consts.ErrMissingParametr, fmt.Errorf(consts.InternalMissingParametr)
	}

	var userID int
	switch typeGet {
	case 0:
		if clientErr, err := utils.ConvertValue(r.Header.Get("X-User-ID"), &userID); err != nil {
			return nil, http.StatusUnauthorized, clientErr, err
		}
		log.Printf("r.Header.Get('X-User-ID') %v", r.Header.Get("X-User-ID"))
	case 1:
		if clientErr, err := utils.ConvertValue(r.URL.Query().Get("userID"), &userID); err != nil {
			return nil, http.StatusUnauthorized, clientErr, err
		}
		log.Printf("url:%d", userID)
	}

	response, statusCode, clientErr, err := utils.ProxyRequest(consts.GET_Method, fmt.Sprintf(consts.DB_GetChatHistory, chatID, userID), nil, http.StatusOK)
	if err != nil {
		return nil, statusCode, clientErr, err
	}

	var chat models.Chat
	if err := utils.ConvertToModel(response.Data, &chat); err != nil {
		return nil, http.StatusInternalServerError, consts.ErrInternalServer, err
	}

	return &chat, http.StatusOK, "", nil
}

// GetQueryChat searches for a chat by its name.
// It returns the response and appropriate status code based on the request.
func GetQueryChat(r *http.Request) (*models.Response, int, string, error) {
	symbols := r.URL.Query().Get("chatName")
	if strings.ReplaceAll(symbols, " ", "") == "" {
		return nil, http.StatusBadRequest, consts.ErrMissingChatInfo, fmt.Errorf(consts.InternalMissingParametr)
	}

	response, statusCode, clientErr, err := utils.ProxyRequest(consts.GET_Method, fmt.Sprintf(consts.DB_SearchQueryChat, symbols), nil, http.StatusOK)
	if err != nil {
		return nil, statusCode, clientErr, err
	}

	return response, http.StatusOK, "", nil
}

// GetAllChats retrieves all chats available.
// It returns the response containing all chats and the corresponding status code.
func GetAllChats(r *http.Request) (*models.Response, int, string, error) {
	response, statusCode, clientErr, err := utils.ProxyRequest(consts.GET_Method, consts.DB_GetAllChats, nil, http.StatusOK)
	if err != nil {
		return nil, statusCode, clientErr, err
	}

	return response, http.StatusOK, "", nil
}

// UserChats retrieves the chats associated with a specific user.
// It validates the user ID and fetches the corresponding chats from the database.
func UserChats(r *http.Request) (*models.Response, int, string, error) {
	var userID int
	if clientErr, err := utils.ConvertValue(r.Header.Get("X-User-ID"), &userID); err != nil {
		return nil, http.StatusUnauthorized, clientErr, err
	}

	response, statusCode, clientErr, err := utils.ProxyRequest(consts.GET_Method, fmt.Sprintf(consts.DB_GetUserChats, userID), nil, http.StatusOK)
	if err != nil {
		return nil, statusCode, clientErr, err
	}

	return response, http.StatusOK, "", nil
}

// SetTagToChat updates the tag for a chat based on the request body.
// It sends the update request to the database and returns the response.
func SetTagToChat(r *http.Request) (*models.Response, int, string, error) {
	var changeInfo models.MemberRequest
	if err := json.NewDecoder(r.Body).Decode(&changeInfo); err != nil {
		return nil, http.StatusBadRequest, consts.ErrInternalServer, err
	}

	if clientErr, err := utils.ConvertValue(r.Header.Get("X-User-ID"), &changeInfo.UserID); err != nil {
		return nil, http.StatusUnauthorized, clientErr, err
	}

	response, statusCode, clientErr, err := utils.ProxyRequest(consts.PATCH_Method, consts.DB_SetTag, changeInfo, http.StatusOK)
	if err != nil {
		return nil, statusCode, clientErr, err
	}

	return response, http.StatusOK, "", nil
}

// GetTags retrieves the tags associated with a chat for a specific user.
// It validates the user ID and calls the appropriate utility function to fetch tags.
func GetTags(r *http.Request) (*models.Response, int, string, error) {
	var userID int
	if clientErr, err := utils.ConvertValue(r.Header.Get("X-User-ID"), &userID); err != nil {
		return nil, http.StatusUnauthorized, clientErr, err
	}

	response, statusCode, clientErr, err := utils.ProxyRequest(consts.GET_Method, fmt.Sprintf(consts.DB_GetTags, r.URL.Query().Get("chatID"), userID), nil, http.StatusOK)
	if err != nil {
		return nil, statusCode, clientErr, err
	}

	return response, http.StatusOK, "", nil
}

// DeleteTags removes tags from a chat based on the request body.
// It sends the delete request to the database and returns the response.
func DeleteTags(r *http.Request) (*models.Response, int, string, error) {
	var data *models.MemberRequest
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return nil, http.StatusBadRequest, consts.ErrInvalidValueFormat, err
	}

	if clientErr, err := utils.ConvertValue(r.Header.Get("X-User-ID"), &data.UserID); err != nil {
		return nil, http.StatusUnauthorized, clientErr, err
	}

	response, statusCode, clientErr, err := utils.ProxyRequest(consts.DELETE_Method, consts.DB_DeleteTags, data, http.StatusOK)
	if err != nil {
		return nil, statusCode, clientErr, err
	}

	return response, http.StatusOK, "", nil
}

// JoinToChat allows a user to join a specific chat based on the chat ID from the URL query.
// It sends a request to add the user to the chat and returns the response.
func JoinToChat(r *http.Request) (*models.Response, int, string, error) {
	var userID int
	if clientErr, err := utils.ConvertValue(r.Header.Get("X-User-ID"), &userID); err != nil {
		return nil, http.StatusUnauthorized, clientErr, err
	}

	var chatID int
	if clientErr, err := utils.ConvertValue(r.URL.Query().Get("chatID"), &chatID); err != nil {
		return nil, http.StatusInternalServerError, clientErr, err
	}

	response := &models.Response{}
	if clientErr, err := userservice.ManageMember(consts.POST_Method, consts.DB_AddMembers, models.MemberRequest{
		UserID: userID,
		ChatID: chatID,
	}, http.StatusAccepted); err != nil {
		response.Errors = append(response.Errors, clientErr)
	}

	return response, http.StatusOK, "", nil
}

// DeleteChat removes a chat identified by chatID from the database.
// It returns the status code and any error messages.
func DeleteChat(r *http.Request, chatID int) (int, string, error) {
	var userID int
	if clientErr, err := utils.ConvertValue(r.Header.Get("X-User-ID"), &userID); err != nil {
		return http.StatusUnauthorized, clientErr, err
	}

	if _, deleteStatusCode, clientErr, err := utils.ProxyRequest(consts.DELETE_Method, fmt.Sprintf(consts.DB_DeleteChat, chatID), nil, http.StatusNoContent); err != nil {
		return deleteStatusCode, clientErr, err
	}

	return http.StatusOK, "", nil
}

// CheckAuthor verifies if the user is the author of the specified chat.
// It returns the user ID, chat ID, status code, and any error messages.
func CheckAuthor(r *http.Request) (int, int, int, string, error) {
	var chatID int
	if clientErr, err := utils.ConvertValue(r.URL.Query().Get("chatID"), &chatID); err != nil {
		return 0, 0, http.StatusBadRequest, clientErr, err
	}

	var userID int
	if clientErr, err := utils.ConvertValue(r.Header.Get("X-User-ID"), &userID); err != nil {
		return 0, 0, http.StatusUnauthorized, clientErr, err
	}

	if _, authorStatusCode, clientErr, err := utils.ProxyRequest(consts.GET_Method, fmt.Sprintf(consts.DB_GetAuthor, userID, chatID), nil, http.StatusOK); err != nil {
		return 0, 0, authorStatusCode, clientErr, err
	}

	return userID, chatID, http.StatusOK, "", nil
}

// MemberList retrieves a list of members from the request body.
// It validates the list and returns it along with the status code and any error messages.
func MemberList(r *http.Request) ([]models.MemberRequest, int, string, error) {
	var membersList []models.MemberRequest
	if err := json.NewDecoder(r.Body).Decode(&membersList); err != nil {
		return nil, http.StatusInternalServerError, consts.ErrInternalServer, fmt.Errorf(consts.InternalInvalidValueFormat, err)
	}
	if len(membersList) == 0 {
		return nil, http.StatusInternalServerError, consts.ErrUsersNotFoundInList, fmt.Errorf(consts.InternalNoFoundUsers)

	}

	return membersList, http.StatusOK, "", nil
}
