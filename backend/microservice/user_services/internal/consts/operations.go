package consts

// Operation types for services
const (
	VALDIDATION_RegistrationType = iota
	VALDIDATION_LoginType

	FRIENDSHIP_RequestType
	FRIENDSHIP_AcceptType
	FRIENDSHIP_GetFriendsType
	FRIENDSHIP_DeleteFriendType
	FRIENDSHIP_DeleteFriendRequestType
	FRIENDSHIP_GetSubsType

	NOTIFICATION_AddFriend
	NOTIFICATION_DeleteFriend
	NOTIFICATION_AcceptFriendship

	KAFKA_Friendship
	KAFKA_Session
	KAFKA_PushFriends
	KAFKA_PushSubscribers
	KAFKA_UpdateSubscribers
	KAFKA_RemoveFriendAndAddSubscriber
	KAFKA_RemoveSubscriberAndAddFriend
	KAFKA_session_UPDATE
)
