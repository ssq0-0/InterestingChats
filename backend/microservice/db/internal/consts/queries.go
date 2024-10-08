package consts

// Database query
const (
	QUERY_SelectUser            = "SELECT id FROM users WHERE email = $1"
	QUERY_SelectUsername        = "SELECT username FROM users WHERE id=$1"
	QUERY_SelectEmail           = "SELECT email FROM users WHERE id=$1"
	QUERY_SelectUserExists      = "SELECT EXISTS(SELECT 1 FROM users WHERE id=$1)"
	QUERY_SelectUserData        = "SELECT id, username, password FROM users WHERE email = $1"
	QUERY_SelectProfileUserInfo = "SELECT username, email, avatar FROM users WHERE id = $1"
	QUERY_ChangeUsername        = "UPDATE users SET username =$1 WHERE id=$2"
	QUERY_ChangeEmail           = "UPDATE users SET email =$1 WHERE id=$2"
	QUERY_CheckSingleEmail      = "SELECT id FROM users WHERE email=$1 AND id<>$2"
	QUERY_SelectIndexUser       = "SELECT id, username, email FROM users WHERE to_tsvector('russian', username) @@ to_tsquery($1 || ':*') OR to_tsvector('russian', email) @@ to_tsquery($1 || ':*') OR to_tsvector('english', username) @@ to_tsquery($1 || ':*') OR to_tsvector('english', email) @@ to_tsquery($1 || ':*')"
	QUERY_AddAvatar             = "UPDATE users SET avatar = $1 WHERE id = $2"

	QUERY_InsertToFriendshipRequests = "INSERT INTO friend_requests (requester_id, receiver_id) VALUES ($1::integer, $2::integer) ON CONFLICT (requester_id, receiver_id) DO NOTHING"
	QUERY_GetInitator                = "SELECT id AS user_id, username FROM users WHERE id = $1;"

	QUERY_SelectFromFriendships         = "SELECT EXISTS (SELECT 1 FROM friendships WHERE (user_id1 = LEAST($1::integer, $2::integer) AND user_id2 = GREATEST($1::integer, $2::integer)))"
	QUERY_SelectFromFriendshipsRequests = "SELECT EXISTS (SELECT 1 FROM friend_requests WHERE requester_id = $1 AND receiver_id = $2)"
	QUERY_SelectUsersFriendships        = "SELECT u.id AS friend_id, u.username AS friend_name, u.email AS friend_email, u.avatar AS friend_avatar FROM friendships f JOIN users u ON u.id = CASE WHEN f.user_id1 = $1 THEN f.user_id2 ELSE f.user_id1 END WHERE $1 IN (f.user_id1, f.user_id2)"
	QUERY_SelectSubscribers             = "SELECT u.id AS subscriber_id, u.username AS subscriber_username, u.email AS subscriber_email, u.avatar AS subscriber_avatar FROM friend_requests fr JOIN users u ON u.id = fr.requester_id WHERE fr.receiver_id = $1 AND fr.status = 'pending'"
	QUERY_InsertFriendships             = "INSERT INTO friendships (user_id1, user_id2) VALUES (LEAST($1::integer, $2::integer), GREATEST($1::integer, $2::integer)) ON CONFLICT (user_id1, user_id2) DO NOTHING"

	QUERY_InsertFriendshipRequest = "INSERT INTO friend_requests (requester_id, receiver_id) VALUES ($1::integer, $2::integer) ON CONFLICT (requester_id, receiver_id) DO NOTHING"

	QUERY_DeleteFriendshipRequest = "DELETE FROM friend_requests WHERE requester_id = $1 AND receiver_id = $2"
	QUERY_DeleteFriendship        = "DELETE FROM friendships WHERE user_id1 = LEAST($1::integer, $2::integer) AND user_id2 = GREATEST($1::integer, $2::integer)"

	QUERY_SelectChatInfo    = "SELECT name, author_id FROM chats WHERE id = $1"
	QUERY_SelectAllChats    = "SELECT id, name, author_id FROM chats"
	QUERY_SelectIndexChat   = "SELECT c.id, c.name, c.author_id FROM chats c WHERE to_tsvector('russian', c.name) @@ to_tsquery($1) OR to_tsvector('english', c.name) @@ to_tsquery($1) UNION SELECT c.id, c.name, c.author_id FROM chats c JOIN chat_hashtags ch ON c.id = ch.chat_id JOIN hashtags h ON ch.hashtag_id = h.id WHERE to_tsvector('russian', h.name) @@ to_tsquery($1)  OR to_tsvector('english', h.name) @@ to_tsquery($1)"
	QUERY_SelectUsersChat   = "SELECT c.id, c.name, c.author_id FROM chats c JOIN chat_members cm ON c.id = cm.chat_id WHERE cm.user_id = $1;"
	QUERY_UserChatExists    = "SELECT EXISTS(SELECT 1 FROM chat_members WHERE chat_id = $1 AND user_id = $2)"
	QUERY_SelectChatMembers = "SELECT u.id, u.username, u.email, u.avatar as avatar FROM chat_members cm JOIN users u ON cm.user_id = u.id WHERE cm.chat_id = $1;"

	QUERY_SelectChatMessages = "SELECT messages.id, messages.body, messages.created_at, users.id, users.username, users.email FROM messages JOIN users ON messages.user_id = users.id WHERE messages.chat_id = $1 ORDER BY messages.created_at"
	QUERY_SelectAuthor       = "SELECT u.id FROM chats c JOIN users u ON c.author_id = u.id WHERE c.id=$1 LIMIT 1"

	QUERY_InsertUser    = "INSERT INTO users(username, email, password) VALUES ($1, $2, $3) RETURNING id"
	QUERY_InsertChat    = "INSERT INTO chats(name, created_at, author_id) VALUES($1, $2, $3) RETURNING id"
	QUERY_InsertMembers = "INSERT INTO chat_members(chat_id, user_id) VALUES($1, $2)"
	QUERY_InsertMessage = "INSERT INTO messages(chat_id, user_id, body, created_at) VALUES($1, $2, $3, $4)"

	QUERY_DeleteChat           = "DELETE FROM chats WHERE id = $1"
	QUERY_DeleteMemberFromChat = "DELETE FROM chat_members WHERE chat_id = $1 AND user_id=$2"

	QUERY_ChangeChatName = "UPDATE chats SET name = $1 WHERE id = $2"

	QUERY_AddNotification        = "INSERT INTO notifications (user_id, sender_id, type, message) VALUES ($1, $2, $3, $4)"
	QUERY_GetUserNotifications   = "SELECT id, user_id, sender_id, type, message, created_at, is_read FROM notifications WHERE user_id = $1 ORDER BY created_at DESC"
	QUERY_MarkNotificationIsRead = "UPDATE notifications SET is_read = TRUE WHERE id = ANY($1)"

	QUERY_InsertInHashtag      = "INSERT INTO hashtags (name) VALUES ($1) ON CONFLICT (name) DO NOTHING RETURNING id"
	QUERY_InsertInChatHashtags = "INSERT INTO chat_hashtags (chat_id, hashtag_id) VALUES ($1, $2) ON CONFLICT (chat_id, hashtag_id) DO NOTHING"
	QUERY_SelectTags           = "SELECT h.id, h.name FROM chat_hashtags ch JOIN hashtags h ON ch.hashtag_id = h.id WHERE ch.chat_id = $1"
	QUERY_DeleteTags           = "DELETE FROM chat_hashtags WHERE chat_id = $1 AND hashtag_id = ANY($2::integer[]);"
)
