package consts

const (
	ErrInternalServer     = "Internal error server"
	ErrInvalidValueFormat = "Invalid value format"
	ErrMissingParametr    = "Missing parametr in request"

	ErrTokenMalformed    = "token is malformed"
	ErrTokenExpired      = "token is expired"
	ErrTokenNotValid     = "token is not valid yet"
	ErrNotHandleForToken = "could not handle this token"
	ErrInvalidToken      = "invalid token"
	ErrUserUnathorized   = "User unauthorized"

	ErrUnexpectedStatucCode   = "Unexpected status code received"
	ErrUnexpectedValueFormat  = "Unexpected value format"
	ErrUnexpectedRecivedEmail = "Unexpected recived email"
)

const (
	InternalMissingParametr      = "missing parametr"
	InternalFailedProxyRequest   = "proxy request error: %w"
	InternalUnexpectedStatucCode = "unexpected status code received: %d"
	InternalInvalidValueFormat   = "invalid value format: %w"

	InternalFailedParseBody = "failed to parse data: %w"
	InternalGhostType       = "unfortunate ghost of types"

	InternalGenerateHash   = "failed to generate hash from password: %v"
	InternalGenerateJWT    = "failed generate jwt token: %w"
	InternalFailedSetToken = "failed to set token in redis: %w"

	InternalUserEmailInToken = "failed to find user email in token"
	InternalTokenInvalid     = "invalid token in request"
)
