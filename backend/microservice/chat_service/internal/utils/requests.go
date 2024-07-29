package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func ProxyRequest(method, url string, data interface{}, expectedStatusCode int) ([]byte, int, error) {
	log.Printf("request going to: %s", url)
	var jsonData []byte
	var err error
	if data != nil {
		jsonData, err = json.Marshal(data)
		if err != nil {
			log.Printf("failed to marshal data: %v", err)
			return nil, http.StatusInternalServerError, fmt.Errorf("failed to marshal data: %w", err)
		}
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("failed to create new request: %v", err)
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to create new request: %w", err)
	}
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("failed to send request: %v", err)
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to read response body: %w", err)
	}
	log.Printf("response body: %s", string(body))

	if resp.StatusCode != expectedStatusCode {
		return body, resp.StatusCode, fmt.Errorf("received unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	return body, resp.StatusCode, nil
}
