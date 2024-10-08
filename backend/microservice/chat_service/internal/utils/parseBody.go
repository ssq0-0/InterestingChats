package utils

import (
	"chat_service/internal/consts"
	"encoding/json"

	"fmt"
)

// ParseBody parses the JSON body and unmarshals it into the provided result type T.
// Returns an error if parsing fails.
func ParseBody[T any](body []byte, result *T) error {
	if err := json.Unmarshal(body, result); err != nil {
		return fmt.Errorf(consts.InternalFailedParseBody, err)
	}
	return nil
}
