package consts

const (
	QUERY_SelectUser            = "SELECT id FROM users WHERE email = $1"
	QUERY_SelectUsername        = "SELECT username FROM users WHERE id=$1"
	QUERY_SelectEmail           = "SELECT email FROM users WHERE id=$1"
	QUERY_SelectUserExists      = "SELECT EXISTS(SELECT 1 FROM users WHERE id=$1)"
	QUERY_SelectUserData        = "SELECT id, username, password FROM users WHERE email = $1"
	QUERY_SelectProfileUserInfo = "SELECT username, email FROM users WHERE id = $1"
	QUERY_ChangeUsername        = "UPDATE users SET username =$1 WHERE id=$2"
	QUERY_ChangeEmail           = "UPDATE users SET email =$1 WHERE id=$2"
	QUERY_CheckSingleEmail      = "SELECT id FROM users WHERE email=$1 AND id<>$2"
	QUERY_SelectIndexUser       = "SELECT id, username, email FROM users WHERE to_tsvector('russian', username) @@ to_tsquery($1 || ':*') OR to_tsvector('russian', email) @@ to_tsquery($1 || ':*') OR to_tsvector('english', username) @@ to_tsquery($1 || ':*') OR to_tsvector('english', email) @@ to_tsquery($1 || ':*')"

	QUERY_SelectChatInfo     = "SELECT name, author_id FROM chats WHERE id = $1"
	QUERY_SelectAllChats     = "SELECT id, name, author_id FROM chats"
	QUERY_SelectIndexChat    = "SELECT id, name, author_id FROM chats WHERE to_tsvector('russian', name) @@ to_tsquery($1) OR to_tsvector('english', name) @@ to_tsquery($1)"
	QUERY_SelectUsersChat    = "SELECT c.id, c.name, c.author_id FROM chats c JOIN chat_members cm ON c.id = cm.chat_id WHERE cm.user_id = $1;"
	QUERY_UserChatExists     = "SELECT EXISTS(SELECT 1 FROM chat_members WHERE chat_id = $1 AND user_id = $2)"
	QUERY_SelectChatMembers  = "SELECT u.id, u.username, u.email FROM chat_members cm JOIN users u ON cm.user_id = u.id WHERE cm.chat_id = $1"
	QUERY_SelectChatMessages = "SELECT messages.id, messages.body, messages.created_at, users.id, users.username, users.email FROM messages JOIN users ON messages.user_id = users.id WHERE messages.chat_id = $1 ORDER BY messages.created_at"
	QUERY_SelectAuthor       = "SELECT u.id FROM chats c JOIN users u ON c.author_id = u.id WHERE c.id=$1 LIMIT 1"

	QUERY_InsertUser    = "INSERT INTO users(username, email, password) VALUES ($1, $2, $3) RETURNING id"
	QUERY_InsertChat    = "INSERT INTO chats(name, created_at, author_id) VALUES($1, $2, $3) RETURNING id"
	QUERY_InsertMembers = "INSERT INTO chat_members(chat_id, user_id) VALUES($1, $2)"
	QUERY_InsertMessage = "INSERT INTO messages(chat_id, user_id, body, created_at) VALUES($1, $2, $3, $4)"

	QUERY_DeleteChat           = "DELETE FROM chats WHERE id = $1"
	QUERY_DeleteMemberFromChat = "DELETE FROM chat_members WHERE chat_id = $1 AND user_id=$2"
)
