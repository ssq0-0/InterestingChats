package consts

import (
	"chat_service/internal/config"
)

var (
	GET_Method    = "GET"
	POST_Method   = "POST"
	PATCH_Method  = "PATCH"
	DELETE_Method = "DELETE"

	// Здесь будут храниться URL для обращения к сервису
	DB_GetChatHistory   string
	DB_GetAllChats      string
	DB_CreateChat       string
	DB_DeleteChat       string
	DB_GetAuthor        string
	DB_AddMembers       string
	DB_ChatMemberExists string
	DB_CheckUser        string
	DB_DeleteMember     string
	DB_GetUserChats     string
	DB_JoinToChat       string
	DB_SearchQueryChat  string
	DB_SetTag           string
	DB_GetTags          string
	DB_DeleteTags       string
	DB_ChangeChatName   string
)

// InitConstants инициализирует URL для базы данных
func InitConstants(cfg *config.Config) {
	dbService := cfg.Services.DBService

	DB_GetChatHistory = dbService + "/getChat?chatID=%s&userID=%d"
	DB_GetAllChats = dbService + "/getAllChats"
	DB_CreateChat = dbService + "/createChat"
	DB_DeleteChat = dbService + "/deleteChat?chatID=%d"
	DB_GetAuthor = dbService + "/getAuthor?userID=%v&chatID=%v"
	DB_AddMembers = dbService + "/addMembers"
	DB_ChatMemberExists = dbService + "/chatMemberExists?email=%s"
	DB_CheckUser = dbService + "/checkUser?userID=%d"
	DB_DeleteMember = dbService + "/deleteMember"
	DB_GetUserChats = dbService + "/getUserChats?userID=%d"
	DB_JoinToChat = dbService + "/joinToChat"
	DB_SearchQueryChat = dbService + "/searchChat?chatName=%s"
	DB_SetTag = dbService + "/setTag?chatID=%s"
	DB_GetTags = dbService + "/getTags?chatID=%s&userID=%d"
	DB_DeleteTags = dbService + "/deleteTags"
	DB_ChangeChatName = dbService + "/changeChatName?chatID=%s&chatName=%s"
}
