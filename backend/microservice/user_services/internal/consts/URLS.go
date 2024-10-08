package consts

import "InterestingChats/backend/user_services/internal/config"

// URL ohter services
var (
	DB_Registration         string
	DB_Login                string
	DB_GetUserProfileInfo   string
	DB_ChangeUserData       string
	DB_SearchUsersBySymbols string
	DB_UploadPhoto          string
	DB_AddFriend            string
	DB_AcceptFriend         string
	DB_GetFriendList        string
	DB_DeleteFriend         string
	DB_GetSubs              string
	DB_DeleteFriendRequest  string

	DB_ProfileInfo string

	FS_UploadFile string

	AS_GenerateTokens string

	CACHE_GetUserProfileInfo string
	CACHE_GetFriendList      string
	CACHE_UpdateData         string
	CACHE_GetSubs            string
)

// InitConstants init URL for other services
func InitConstants(cfg *config.Config) {
	dbService := cfg.Services.DBService
	fileService := cfg.Services.FileService
	authService := cfg.Services.AuthService
	cacheService := cfg.Services.CacheService

	DB_Registration = dbService + "/registration"
	DB_Login = dbService + "/login"
	DB_GetUserProfileInfo = dbService + "/profileInfo?userID=%v"
	DB_ChangeUserData = dbService + "/changeUserData"
	DB_SearchUsersBySymbols = dbService + "/searchUsers?symbols=%s"
	DB_UploadPhoto = dbService + "/uploadPhoto"
	DB_AddFriend = dbService + "/requestToFriendShip?"
	DB_AcceptFriend = dbService + "/acceptFriendShip?requestID=%v&receiverID=%v"
	DB_GetFriendList = dbService + "/getFriends?userID=%d"
	DB_DeleteFriend = dbService + "/deleteFriend?remover=%d&removing=%s"
	DB_GetSubs = dbService + "/getSubs?userID=%d"
	DB_DeleteFriendRequest = dbService + "/deleteFriendRequest?remover=%v&removing=%v"

	DB_ProfileInfo = dbService + "/profileInfo?userID=%v"

	FS_UploadFile = fileService + "/saveImage"

	AS_GenerateTokens = authService + "/generate_tokens"

	CACHE_GetUserProfileInfo = cacheService + "/getSession?userID=%v"
	CACHE_GetFriendList = cacheService + "/getFriendList?userID=%v"
	CACHE_UpdateData = cacheService + "/updateData?userID=%v"
	CACHE_GetSubs = cacheService + "/getSubscribers?userID=%v"
}
