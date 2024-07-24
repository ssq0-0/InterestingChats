package consts

const (
	GET_Method    = "GET"
	POST_Method   = "POST"
	DELETE_Method = "DELETE"

	DB_GetChatHistory = "http://localhost:8002/getChat?chatName=%s"
	DB_CreateChat     = "http://localhost:8002/createChat"
	DB_DeleteChat     = "http://localhost:8002/deleteChat?chatID=%s"
	DB_GetAuthor      = "http://localhost:8002/getAuthor?email=%s&chatID=%d"
	DB_AddMembers     = "http://localhost:8002/addMembers"
	DB_CheckUser      = "http://localhost:8002/checkUser?userID=%d"
	DB_DeleteMember   = "http://localhost:8002/deleteMember"

	US_CheckToken = "http://localhost:8001/checkToken?token=%s"
)
