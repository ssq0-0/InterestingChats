package consts

const (
	POST_Method = "POST"
	GET_Method  = "GET"

	Redis_SetToken = "http://localhost:8003/setToken"
	Redis_GetToken = "http://localhost:8003/user?email=%s"

	DB_Registration         = "http://localhost:8002/registration"
	DB_Login                = "http://localhost:8002/login"
	DB_GetUserProfileInfo   = "http://localhost:8002/profileInfo?userID=%d"
	DB_ChangeUserData       = "http://localhost:8002/changeUserData"
	DB_SearchUsersBySymbols = "http://localhost:8002/searchUsers?symbols=%s"

	VALDIDATION_RegistrationType = 0
	VALDIDATION_LoginType        = 1
)
