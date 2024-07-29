package utils

import (
	"fmt"
	"strconv"
)

func ConvertValue[T any](value interface{}, result *T) error {
	switch v := value.(type) {
	case string:
		if _, ok := any(*result).(int); ok {
			intValue, err := strconv.Atoi(v)
			if err != nil {
				return fmt.Errorf("failed to convert value:%v", err)
			}
			*result = any(intValue).(T)
		} else {
			return fmt.Errorf("unsupported type for string conversion: %T", *result)
		}
	case int:
		if _, ok := any(*result).(string); ok {
			*result = any(strconv.Itoa(v)).(T)
		} else {
			return fmt.Errorf("failed to convert value:%v", value)
		}
	default:
		return fmt.Errorf("unsupported type: %T", value)
	}

	return nil
}
