package utils

import (
	"bytes"
	"chat_service/internal/consts"
	"chat_service/internal/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// ProxyRequest performs an HTTP request and proxies the response.
// It expects a specific status code and handles JSON marshaling/unmarshaling.
// Returns the response model, status code, error message, and an error if any.
func ProxyRequest(method, url string, data interface{}, expectedStatusCode int) (*models.Response, int, string, error) {
	var jsonData []byte
	var err error
	if data != nil {
		jsonData, err = json.Marshal(data)
		if err != nil {
			return nil, http.StatusInternalServerError, consts.ErrInternalServer, fmt.Errorf(consts.InternalInvalidValueFormat, err)
		}
	}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, http.StatusInternalServerError, consts.ErrInternalServer, fmt.Errorf(consts.InternalFailedProxyRequest, err)
	}
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, http.StatusInternalServerError, consts.ErrInternalServer, fmt.Errorf(consts.InternalFailedProxyRequest, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		if resp.StatusCode != expectedStatusCode {
			return nil, resp.StatusCode, consts.ErrUnexpectedStatucCode, fmt.Errorf(consts.InternalUnexpectedStatucCode, resp.StatusCode)
		}
		return nil, resp.StatusCode, "", nil
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, http.StatusInternalServerError, consts.ErrInternalServer, fmt.Errorf(consts.InternalFailedProxyRequest, err)
	}
	// log.Printf("resp: %s", string(body))

	var response *models.Response
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, http.StatusInternalServerError, consts.ErrInternalServer, fmt.Errorf(consts.InternalInvalidValueFormat, err)
	}
	if resp.StatusCode != expectedStatusCode {
		if len(response.Errors) > 0 {
			return nil, resp.StatusCode, strings.Join(response.Errors, "; "), fmt.Errorf(strings.Join(response.Errors, "; "))
		}
		return nil, resp.StatusCode, consts.ErrUnexpectedStatucCode, fmt.Errorf(consts.InternalUnexpectedStatucCode, resp.StatusCode)
	}

	return response, resp.StatusCode, "", nil
}
