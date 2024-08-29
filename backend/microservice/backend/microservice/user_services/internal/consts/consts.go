package consts

const (
	POST_Method = "POST"
	GET_Method  = "GET"

	Redis_SetToken = "http://localhost:8003/setToken"
	Redis_GetToken = "http://localhost:8003/user?email=%s"

	DB_Registration = "http://localhost:8002/registration"
	DB_Login        = "http://localhost:8002/login"
)
