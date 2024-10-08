package utils

import (
	"file_service/internal/consts"
	"log"
	"mime/multipart"
	"net/http"
)

// ParseRequestFormFile parses the multipart form file from the request.
func ParseRequestFormFile(r *http.Request) (multipart.File, *multipart.FileHeader, string, string, string, error) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		log.Printf("Ошибка парсинга формы: %v", err)
		return nil, nil, "", "", consts.ErrInvalidRequestData, err
	}

	file, handler, err := r.FormFile("image")
	if err != nil {
		log.Printf("Ошибка получения файла: %v", err)
		return nil, nil, "", "", consts.ErrInternalServer, err
	}
	defer file.Close()

	newFileName, clientErr, err := RenameFile(handler.Filename)
	if err != nil {
		return nil, nil, "", "", clientErr, err
	}

	mimeType, err := DetectMimeType(file)
	if err != nil {
		log.Printf("Ошибка определения MIME-типа: %v", err)
		mimeType = "application/octet-stream"
	}
	log.Printf("Определен MIME-тип: %s для файла: %s", mimeType, newFileName)

	return file, handler, newFileName, mimeType, "", nil
}
