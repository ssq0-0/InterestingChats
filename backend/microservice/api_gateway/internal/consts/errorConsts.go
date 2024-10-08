package consts

// Error messages used for client-facing errors.
const (
	ErrInternalServerError  = "Internal server error"
	ErrUnexpectedStatucCode = "Unexpected status code"
	ErrMissingParametr      = "Missing parametr"
	ErrInvalidValueFormat   = "Invalid value format"
	ErrUnsupportedType      = "Unsuported format"
)

// Internal error messages used by the server for logging and debugging.
const (
	InternalServerError          = "internal server error: %w"
	InternalInvalidValueFormat   = "invalid value format: %v"
	InternalFailedProxyRequest   = "failed proxy request: %v"
	InternalUnexpectedStatucCode = "unexpectetd status code: %d"
	InternalFailedConvert        = "failed convert value"
)
