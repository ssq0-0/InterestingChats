package utils

import (
	"bytes"
	"io"

	"github.com/valyala/fasthttp"
)

// BodySave saves the request body for reuse.
func BodySave(req *fasthttp.Request) ([]byte, io.Reader, error) {
	bodyBytes := req.Body()

	return bodyBytes, bytes.NewReader(bodyBytes), nil
}
