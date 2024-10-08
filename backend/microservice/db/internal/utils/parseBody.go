package utils

import (
	"encoding/json"
	"fmt"
	"log"
)

// ParseBody to parse the given JSON body into the provided result type.
func ParseBody[T any](body []byte, result *T) error {
	if err := json.Unmarshal(body, result); err != nil {
		log.Printf("Failed to parse response: %v", err)
		return fmt.Errorf("failed to parse response: %v", err)
	}
	return nil
}
