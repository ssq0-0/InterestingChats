package consts

// Error constants that represent various server-side error messages.
const (
	ErrInternalServer     = "Internal error server"
	ErrInvalidValueFormat = "Invalid value format"
	ErrMissingParametr    = "Missing parametr in request"
	ErrUnsupportedType    = "request type unsuported"
)

// Internal error messages for server use.
const (
	InternalServerError     = "internal server error"
	InternalMissingParametr = "missing parametr"

	InternalUnknowParametr = "unknow parametr"
	InternalUnknowTopic    = "recived message from unknow topic: %s %s"
)
