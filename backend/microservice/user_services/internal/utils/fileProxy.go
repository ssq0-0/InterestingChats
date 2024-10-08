package utils

import (
	"InterestingChats/backend/user_services/internal/consts"
	"InterestingChats/backend/user_services/internal/models"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

// FileProxyRequest: HTTP request with type multipart/form-data.
func FileProxyRequest(r *http.Request, method, url string, expectedStatusCode int) (*models.FileResponse, int, string, error) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	file, handler, err := r.FormFile("image")
	if err != nil {
		log.Printf("Ошибка получения файла: %v", err)
		return nil, 400, "", err
	}
	defer file.Close()

	filePart, err := writer.CreateFormFile("image", handler.Filename)
	if err != nil {
		return nil, http.StatusInternalServerError, "Ошибка при создании части формы", fmt.Errorf("failed to create form file: %v", err)
	}
	if _, err = io.Copy(filePart, file); err != nil {
		return nil, http.StatusInternalServerError, "Ошибка при копировании файла", fmt.Errorf("failed to copy file: %v", err)
	}

	if err = writer.Close(); err != nil {
		return nil, http.StatusInternalServerError, "Ошибка при закрытии writer", fmt.Errorf("failed to close writer: %v", err)
	}
	fileExt := filepath.Ext(handler.Filename)
	mimeType := mime.TypeByExtension(fileExt)
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	req, err := http.NewRequest(method, url, &buf)
	if err != nil {
		return nil, http.StatusInternalServerError, "Ошибка создания запроса", fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, http.StatusInternalServerError, "Ошибка выполнения запроса", fmt.Errorf("failed to execute request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, http.StatusInternalServerError, "Ошибка чтения ответа", fmt.Errorf("failed to read response body: %v", err)
	}

	var response models.FileResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, http.StatusInternalServerError, "Ошибка декодирования JSON", fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	if resp.StatusCode != expectedStatusCode {
		if len(response.Errors) > 0 {
			return &response, resp.StatusCode, strings.Join(response.Errors, "; "), fmt.Errorf(strings.Join(response.Errors, "; "))
		}
		return &response, resp.StatusCode, consts.ErrUnexpectedStatucCode, fmt.Errorf(consts.InternalUnexpectedStatucCode, resp.StatusCode)
	}

	return &response, resp.StatusCode, "", nil
}
