package consts

const (
	GET_Method    = "GET"
	POST_Method   = "POST"
	DELETE_Method = "DELETE"

	DB_GetChatHistory   = "http://localhost:8002/getChat?chatID=%s&userID=%d"
	DB_GetAllChats      = "http://localhost:8002/getAllChats"
	DB_CreateChat       = "http://localhost:8002/createChat"
	DB_DeleteChat       = "http://localhost:8002/deleteChat?chatID=%d"
	DB_GetAuthor        = "http://localhost:8002/getAuthor?userID=%d&chatID=%d"
	DB_AddMembers       = "http://localhost:8002/addMembers"
	DB_ChatMemberExists = "http://localhost:8002/chatMemberExists?email=%s"
	DB_CheckUser        = "http://localhost:8002/checkUser?userID=%d"
	DB_DeleteMember     = "http://localhost:8002/deleteMember"
	DB_SaveMessage      = "http://localhost:8002/saveMessage?chatID=%d"
	DB_GetUserChats     = "http://localhost:8002/getUserChats?userID=%d"
	DB_JoinToChat       = "http://localhost:8002/joinToChat"
	DB_SearchQueryChat  = "http://localhost:8002/searchChat?chatName=%s"

	US_CheckToken = "http://localhost:8001/checkToken?token=%s"
)
