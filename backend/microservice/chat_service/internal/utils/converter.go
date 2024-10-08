package utils

import (
	"chat_service/internal/consts"
	"fmt"
	"log"
	"strconv"
)

// ConvertValue converts a given interface{} value to a specified type T.
// It returns an error if the conversion fails or if the value is nil.
func ConvertValue[T any](value interface{}, result *T) (string, error) {
	if value == nil {
		log.Printf("here: %T", value)
		return consts.ErrMissingParametr, fmt.Errorf(consts.InternalFailedConvertValie)
	}
	log.Printf("here1: %T", value)

	switch v := value.(type) {
	case string:
		if _, ok := any(*result).(int); ok {
			intValue, err := strconv.Atoi(v)
			if err != nil {
				log.Printf("here3: %T", value)

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
