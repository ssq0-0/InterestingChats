package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func SendRequest(method, url string, data interface{}) ([]byte, int, error) {
	var jsonReqData []byte
	var err error
	if data != nil {
		jsonReqData, err = json.Marshal(data)
		if err != nil {
			log.Println("Failed to serialize data.")
			return nil, http.StatusInternalServerError, fmt.Errorf("failed to serialize data: %w", err)
		}
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonReqData))
	if err != nil {
		log.Println("Failed to create request.")
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Failed to send request.")
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Failed to read response.")
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to read response: %w", err)
	}

	return body, resp.StatusCode, nil
}
