package utils

import (
	"chat_service/internal/consts"
	"encoding/json"

	"fmt"
)

func ParseBody[T any](body []byte, result *T) error {
	if err := json.Unmarshal(body, result); err != nil {
		return fmt.Errorf(consts.InternalFailedParseBody, err)
	}
	return nil
}
