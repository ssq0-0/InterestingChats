package consts

const (
	ErrInternalServerError  = "Internal server error"
	ErrUnexpectedStatucCode = "Unexpected status code"
	ErrMissingParametr      = "Missing parametr"
	ErrInvalidValueFormat   = "Invalid value format"
	ErrUnsupportedType      = "Unsuported format"
)

const (
	InternalServerError          = "internal server error: %w"
	InternalInvalidValueFormat   = "invalid value format: %v"
	InternalFailedProxyRequest   = "failed proxy request: %v"
	InternalUnexpectedStatucCode = "unexpectetd status code: %d"
	InternalFailedConvert        = "failed convert value"
)
