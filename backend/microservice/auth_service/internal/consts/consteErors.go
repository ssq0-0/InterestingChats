package consts

// Error messages used for client-facing errors.
const (
	ErrInternalServer     = "Internal error server"
	ErrInvalidValueFormat = "Invalid value format"

	ErrTokenMalformed    = "token is malformed"
	ErrTokenExpired      = "token is expired"
	ErrTokenNotValid     = "token is not valid yet"
	ErrNotHandleForToken = "could not handle this token"
	ErrInvalidToken      = "invalid token"
	ErrUserUnathorized   = "User unauthorized"

	ErrUnexpectedStatucCode = "Unexpected status code received"
)

// Internal error messages used by the server for logging and debugging.
const (
	InternalServerError          = "internal server error"
	InternalFailedProxyRequest   = "proxy request error: %w"
	InternalUnexpectedStatucCode = "unexpected status code received: %d"
	InternalInvalidValueFormat   = "invalid value format: %w"

	InternalTokenInvalid = "invalid token in request"
	InternalTokenError   = "failed pase refresh token from headers"
)
