package utils

import (
	"InterestingChats/backend/microservice/db/internal/consts"
	"fmt"
	"strconv"
)

func ConvertValue[T any](value interface{}, result *T) (string, error) {
	switch v := value.(type) {
	case string:
		if _, ok := any(*result).(int); ok {
			intValue, err := strconv.Atoi(v)
			if err != nil {
				return consts.ErrInvalidValueFormat, fmt.Errorf("failed to convert value:%v", err)
			}
			*result = any(intValue).(T)
		} else {
			return consts.ErrUnsupportedType, fmt.Errorf("unsupported type for string conversion: %T", *result)
		}
	case int:
		if _, ok := any(*result).(string); ok {
			*result = any(strconv.Itoa(v)).(T)
		} else {
			return consts.ErrUnsupportedType, fmt.Errorf("failed to convert value:%v", value)
		}
	default:
		return consts.ErrUnsupportedType, fmt.Errorf("unsupported type: %T", value)
	}

	return "", nil
}
