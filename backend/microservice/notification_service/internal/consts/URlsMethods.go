package consts

import (
	"fmt"
)

// HTTP constants that define the methods for various http requests.
const (
	POST_Method  = "POST"
	GET_Method   = "GET"
	PATCH_Method = "PATCH"
)

// Variables to store URLs
var (
	DB_AddNotification  string
	DB_GetNotifications string
	DB_ReadNotification string
)

func InitServiceURLs(dbProtocol, dbHost, dbPort string) {
	DB_AddNotification = fmt.Sprintf("%s://%s:%s/addNotification", dbProtocol, dbHost, dbPort)
	DB_GetNotifications = fmt.Sprintf("%s://%s:%s/getNotification?userID=%%d", dbProtocol, dbHost, dbPort)
	DB_ReadNotification = fmt.Sprintf("%s://%s:%s/readNotification", dbProtocol, dbHost, dbPort)
}
