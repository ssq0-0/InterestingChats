package consts

// Error messages used for client-facing errors.
const (
	ErrInternalServer       = "Internal server error"
	ErrInvalidValueFormat   = "Invalid data format"
	ErrUnexpectedStatucCode = "Unexpercted status code"
	ErrMissingParametr      = "Missing parametr"
	ErrUnsupportedType      = "Unsupported value type"
)

// Internal error messages used by the server for logging and debugging.
const (
	InternalMissingParametr      = "missing parametr"
	InternalInvalidValueFormat   = "invalid value format: %v"
	InternalFailedProxyRequest   = "failed proxy request: %v"
	InternalUnexpectedStatucCode = "unexpected statusCode: %d"
	InternalErrorConvertValue    = "failed to convert value:%v"
	InternalUnsuportedType       = "unsupported type: %T"
	InternalUnsopertedString     = "unsupported type for string conversion: %T"
)
