package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"notifications/internal/consts"
	"notifications/internal/models"
	"strings"
)

// createRequest creates an HTTP request with the specified method, URL, and data
func createRequest(method, url string, data interface{}) (*http.Request, error) {
	var jsonData []byte
	var err error
	if data != nil {
		jsonData, err = json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf(consts.InternalInvalidValueFormat, err)
		}
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf(consts.InternalFailedProxyRequest, err)
	}
	req.Header.Add("Content-Type", "application/json")
	return req, nil
}

// sendRequest sends an HTTP request and returns the response and body
func sendRequest(client *http.Client, req *http.Request) (*http.Response, []byte, error) {
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf(consts.InternalFailedProxyRequest, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp, nil, fmt.Errorf(consts.InternalFailedProxyRequest, err)
	}

	return resp, body, nil
}

// processResponse processes the HTTP response and checks for errors
func processResponse(resp *http.Response, body []byte, expectedStatusCode int) (*models.Response, int, string, error) {
	if resp.StatusCode == http.StatusNoContent {
		if resp.StatusCode != expectedStatusCode {
			return nil, resp.StatusCode, consts.ErrUnexpectedStatucCode, fmt.Errorf(consts.InternalUnexpectedStatucCode, resp.StatusCode)
		}
		return nil, resp.StatusCode, "", nil
	}

	var response models.Response
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, http.StatusInternalServerError, consts.ErrInternalServer, fmt.Errorf(consts.InternalInvalidValueFormat, err)
	}

	if resp.StatusCode != expectedStatusCode {
		errorMessage := consts.ErrUnexpectedStatucCode
		if len(response.Errors) > 0 {
			errorMessage = strings.Join(response.Errors, "; ")
		}
		return &response, resp.StatusCode, errorMessage, fmt.Errorf(errorMessage)
	}

	return &response, resp.StatusCode, "", nil
}

// ProxyRequest makes an HTTP request and processes the response
func ProxyRequest(method, url string, data interface{}, expectedStatusCode int) (*models.Response, int, string, error) {
	log.Printf("url: %s", url)
	client := &http.Client{}

	req, err := createRequest(method, url, data)
	if err != nil {
		return nil, http.StatusInternalServerError, consts.ErrInternalServer, err
	}

	resp, body, err := sendRequest(client, req)
	if err != nil {
		return nil, http.StatusInternalServerError, consts.ErrInternalServer, err
	}

	return processResponse(resp, body, expectedStatusCode)
}
