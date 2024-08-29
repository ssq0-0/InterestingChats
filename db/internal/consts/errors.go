package consts

const (
	ErrInternalServerError  = "Internal server error"
	ErrInvalidValueFormat   = "Invalid value format"
	ErrUnsupportedType      = "Unsupported type for conversion"
	ErrInvalidRequestFormat = "Invalid request format. Please check your input data."
	ErrMissingParametr      = "Missing parametr in request"

	ErrMissingChatName      = "Missing chat name in URL"
	ErrMissingChatID        = "Missing chat id"
	ErrGetChatInfo          = "Failed to get chat info"
	ErrScanChatMemberResult = "Failed to retrieve chat members"
	ErrScanMessagesResult   = "Failed to retrieve chat messages"
	ErrFailedDeleteChat     = "Failed deleted chat. Internal server error"
	ErrNoAuthorFound        = "No author found for chat"
	ErrChatAlreadyExists    = "Failed to create. Chat with that name already exist.)"
	ErrUserNoChatMember     = "The user is not a member of the chat room"
	ErrUserChatMember       = "The user already chat member"

	ErrUserNoAuthor             = "User no author"
	ErrMissingUserID            = "Missing user id in URL"
	ErrUserNoExists             = "Can't found user."
	ErrUserNotFound             = "User not found."
	ErrUserEmailAlreadyExists   = "A user with this email address already exists"
	ErrIncorrectEmailOrPassword = "Incorrect email or password. Please check you input data"
	ErrFailedVerifyUser         = "Failed to verify if the user exists"
	ErrUserChatExists           = "Failed to verify the existence of the user in the chat room"
	ErrUsernameAlredyExists     = "Failed to change username. This is your current username"
	ErrEmailAlredyExists        = "Failed to change email. This is your current email"
)

const (
	InternalErrSelectInfo      = "failed to select info from db: %w"
	InternalErrDelete          = "failed to delete from db: %w"
	InternalErrFailedInsert    = "failed to insert into db: %w"
	InternalErrChangedUserInfo = "failed to change user info"

	InternalErrCheckUser         = "failed to verify the existence of a user in the chat room: %w"
	InternalErrUserNotFoud       = "user not found"
	InternalErrUserAlreadyExists = "user with email %s already exists"
	InternalErrPasswordNotFound  = "user password not found"

	InternalErrAuthorNoFound         = "no author was found for this chat room: %w"
	InternalErrUserNoAuthor          = "user %d no author chat %d"
	InternalErrUserAlreadyChatMember = "user %d already in chat %d"
	InternalErrUserNoChatMember      = "user no chat member"

	InternalErrScanResult    = "error of scanning query results from the database: %w"
	InternalErrIterateErrors = "failed to iterate on scan errors: %w"
	InternalErrPrepareTx     = "unsuccessful transaction preparation: %w"

	InternalErrFailedRequest   = "invalid value format: %w"
	InternalErrMissingURLVal   = "missing URL parametr"
	InternalErrMissingParametr = "missing  parametr"
)
