package utils

import (
	"mime/multipart"
	"net/http"
)

// DetectMimeType determines the MIME type of a file.
func DetectMimeType(file multipart.File) (string, error) {
	buffer := make([]byte, 512)
	_, err := file.Read(buffer)
	if err != nil {
		return "", err
	}

	file.Seek(0, 0)
	mimeType := http.DetectContentType(buffer)
	return mimeType, nil
}
