package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

func SendRequestToDB(chatName string) ([]byte, int, error) {
	escapedChatName := url.QueryEscape(chatName)
	reqURL := fmt.Sprintf("http://localhost:8002/getChat?chatName=%s", escapedChatName)
	// log.Printf("Sending request to URL: %s", reqURL)

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		log.Println("error create newrequest")
		return nil, http.StatusInternalServerError, fmt.Errorf("error creating new request: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("error create request")
		return nil, http.StatusInternalServerError, fmt.Errorf("error creating request: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Failed to read response.")
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to read response: %w", err)
	}

	return body, resp.StatusCode, nil
}

func SendToUserService(token string) (string, error) {
	url := fmt.Sprintf("http://localhost:8001/checkToken?token=%s", token)
	log.Printf("Sending request to URL: %s", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("error create new request: %v", err)
		return "", fmt.Errorf("error creating new request: %v", err)
	}
	req.Header.Set("Content-Type", "aplication/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("error sending request: %v", err)
		return "", fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("invalid token, status code: %d", resp.StatusCode)
		return "", fmt.Errorf("invalid token, status code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response: %v", err)
		return "", fmt.Errorf("error reading response: %v", err)
	}
	log.Printf("Response body: %s", body)
	var email string
	if err := json.Unmarshal(body, &email); err != nil {
		log.Printf("error decoding response: %v", err)
		return "", fmt.Errorf("error decoding response: %v", err)
	}

	log.Printf("email found: %s", email)
	return email, nil
}

// func SendInfoToDB(chat *models.Chat) ([]byte, int, error) {
// 	jsonData, err := json.Marshal(chat)
// 	if err != nil {
// 		log.Println("error marshaling chat data:", err)
// 		return nil, http.StatusInternalServerError, fmt.Errorf("error marshaling chat data: %w", err)
// 	}

// 	req, err := http.NewRequest("POST", "http://localhost:8002/createChat", bytes.NewBuffer(jsonData))
// 	if err != nil {
// 		log.Println("error create newrequest")
// 		return nil, http.StatusInternalServerError, fmt.Errorf("error creating new request: %w", err)
// 	}
// 	req.Header.Set("Content-Type", "application/json")

// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		log.Println("error create request")
// 		return nil, http.StatusInternalServerError, fmt.Errorf("error creating request: %w", err)
// 	}
// 	defer resp.Body.Close()

// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		log.Println("Failed to read response.")
// 		return nil, http.StatusInternalServerError, fmt.Errorf("failed to read response: %w", err)
// 	}

// 	return body, resp.StatusCode, nil
// }

// // func SendRequestToDB(method string, data interface{}) error {
// // 	var jsonReqData []byte
// // 	var err error
// // 	if data != "" {
// // 		jsonReqData, err = json.Marshal(data)
// // 		if err != nil {
// // 			log.Panicln("Error decode to json")
// // 			return err
// // 		}
// // 	}

// // 	requestURL := "http://localhost:8002/chats"
// // 	req, err := http.NewRequest(method, requestURL, bytes.NewBuffer(jsonReqData))
// // 	if err != nil {
// // 		log.Println("Failed to create request.")
// // 		return fmt.Errorf("failed to create request: %w", err)
// // 	}
// // 	req.Header.Add("Content-Type", "aplication/json")

// // 	client := http.Client{}
// // 	resp, err := client.Do(req)
// // 	if err != nil {
// // 		log.Println("Failed to send request")
// // 		return fmt.Errorf("failed to create request: %w", err)
// // 	}
// // 	defer resp.Body.Close()

// // 	body, err := ioutil.ReadAll(resp.Body)
// // 	if err != nil {
// // 		log.Println("Error decode body response", err, body) // delete body
// // 		return fmt.Errorf("failed to decode body response: %w", err)
// // 	}
// // 	return nil //add body
// // }
