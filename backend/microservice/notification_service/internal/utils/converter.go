package utils

import (
	"fmt"
	"notifications/internal/consts"
	"strconv"
)

// ConvertValue converts a value of any type to the specified type T
func ConvertValue[T any](value interface{}, result *T) (string, error) {
	if value == nil {
		return consts.ErrMissingParametr, fmt.Errorf(consts.InternalMissingParametr)
	}

	switch v := value.(type) {
	case string:
		if _, ok := any(*result).(int); ok {
			intValue, err := strconv.Atoi(v)
			if err != nil {
				return consts.ErrInvalidValueFormat, fmt.Errorf(consts.InternalErrorConvertValue, err)
			}
			*result = any(intValue).(T)
		} else {
			return consts.ErrUnsupportedType, fmt.Errorf(consts.InternalUnsopertedString, *result)
		}
	case int:
		if _, ok := any(*result).(string); ok {
			*result = any(strconv.Itoa(v)).(T)
		} else {
			return consts.ErrUnsupportedType, fmt.Errorf(consts.InternalErrorConvertValue, value)
		}
	default:
		return consts.ErrUnsupportedType, fmt.Errorf(consts.InternalUnsuportedType, value)
	}

	return "", nil
}
