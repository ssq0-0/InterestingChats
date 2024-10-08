package utils

import (
	"InterestingChats/backend/user_services/internal/consts"
	"fmt"
	"strconv"
)

// ConvertValue converts the interface value to the specified T type.
// If value is nil, an error is returned.
func ConvertValue[T any](value interface{}, result *T) (string, error) {
	if value == nil {
		return consts.ErrMissingParametr, fmt.Errorf(consts.ErrInternalServer)
	}

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
