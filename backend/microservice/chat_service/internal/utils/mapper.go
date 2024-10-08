package utils

import (
	"encoding/json"
	"fmt"
)

// ConvertToModel marshals the data into JSON and then unmarshals it into the specified model.
// Returns an error if any step of the conversion fails.
func ConvertToModel(data interface{}, model interface{}) error {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	if err := json.Unmarshal(dataBytes, model); err != nil {
		return fmt.Errorf("failed to unmarshal data into model: %w", err)
	}

	return nil
}
