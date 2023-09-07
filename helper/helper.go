package helper

import "errors"

func ConvertToFloat64Slice(val interface{}) ([]float64, error) {
	// Attempt to type assert the value to a slice of float64
	if floatSlice, ok := val.([]interface{}); ok {
		result := make([]float64, len(floatSlice))
		for i, v := range floatSlice {
			if floatValue, ok := v.(float64); ok {
				result[i] = floatValue
			} else {
				return nil, errors.New("Invalid data type in slice")
			}
		}
		return result, nil
	}
	return nil, errors.New("Invalid data type")
}
