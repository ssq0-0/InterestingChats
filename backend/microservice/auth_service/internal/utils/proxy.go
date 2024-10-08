package utils

import (
	"auth_service/internal/consts"
	"auth_service/internal/models"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// ProxyRequest makes an HTTP request to the given URL with the specified method and data.
// It expects a specific status code and returns the response or an error if something goes wrong.
func ProxyRequest(method, url string, data interface{}, expectedStatusCode int) (*models.Response, int, string, error) {
	var jsonData []byte
	var err error
	if data != nil {
		jsonData, err = json.Marshal(data)
		if err != nil {
			return nil, http.StatusInternalServerError, consts.ErrInvalidValueFormat, fmt.Errorf(consts.InternalInvalidValueFormat, err)
		}
	}
	log.Printf("url: %s", url)
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
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, http.StatusInternalServerError, consts.ErrInternalServer, fmt.Errorf(consts.InternalInvalidValueFormat, err)
	}
	var response *models.Response
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, http.StatusInternalServerError, consts.ErrInternalServer, fmt.Errorf(consts.InternalInvalidValueFormat, err)
	}

	if resp.StatusCode != expectedStatusCode {
		if len(response.Errors) > 0 {
			return response, resp.StatusCode, strings.Join(response.Errors, "; "), fmt.Errorf(strings.Join(response.Errors, "; "))
		}
		return response, resp.StatusCode, consts.ErrUnexpectedStatucCode, fmt.Errorf(consts.InternalUnexpectedStatucCode, resp.StatusCode)
	}

	return response, resp.StatusCode, "", nil
}
