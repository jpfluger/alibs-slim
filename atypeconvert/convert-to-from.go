package atypeconvert

import (
	"fmt"
	"strconv"
)

// ConvertToStringFrom converts a value of any basic type to a string.
func ConvertToStringFrom(value interface{}) (string, error) {
	if value == nil {
		return "", nil
	}

	switch v := value.(type) {
	case string:
		return v, nil
	case int, int64, int32, int16, int8, uint, uint64, uint32, uint16, uint8:
		return fmt.Sprint(value), nil
	case float64, float32:
		return fmt.Sprintf("%f", value), nil
	case bool:
		return fmt.Sprintf("%t", value), nil
	default:
		return "", fmt.Errorf("unsupported conversion from %T to string", v)
	}
}

// ConvertToIntFrom converts a value of any basic type to an int.
func ConvertToIntFrom(value interface{}) (int, error) {
	if value == nil {
		return 0, nil
	}

	switch v := value.(type) {
	case string:
		return strconv.Atoi(v)
	case int:
		return v, nil
	case int64:
		return int(v), nil
	case int32:
		return int(v), nil
	case int16:
		return int(v), nil
	case int8:
		return int(v), nil
	case uint:
		return int(v), nil
	case uint64:
		return int(v), nil
	case uint32:
		return int(v), nil
	case uint16:
		return int(v), nil
	case uint8:
		return int(v), nil
	case float64:
		return int(v), nil
	case float32:
		return int(v), nil
	case bool:
		if v {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, fmt.Errorf("unsupported conversion from %T to int", v)
	}
}

// ConvertToFloatFrom converts a value of any basic type to a float64.
func ConvertToFloatFrom(value interface{}) (float64, error) {
	if value == nil {
		return 0.0, nil
	}

	switch v := value.(type) {
	case string:
		return strconv.ParseFloat(v, 64)
	case int:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int16:
		return float64(v), nil
	case int8:
		return float64(v), nil
	case uint:
		return float64(v), nil
	case uint64:
		return float64(v), nil
	case uint32:
		return float64(v), nil
	case uint16:
		return float64(v), nil
	case uint8:
		return float64(v), nil
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	case bool:
		if v {
			return 1.0, nil
		}
		return 0.0, nil
	default:
		return 0.0, fmt.Errorf("unsupported conversion from %T to float64", v)
	}
}

// ConvertToBoolFrom converts a value of any basic type to a bool.
func ConvertToBoolFrom(value interface{}) (bool, error) {
	if value == nil {
		return false, nil
	}

	switch v := value.(type) {
	case string:
		return strconv.ParseBool(v)
	case int, int64, int32, int16, int8, uint, uint64, uint32, uint16, uint8:
		return value != 0, nil
	case float64, float32:
		return value != 0.0, nil
	case bool:
		return v, nil
	default:
		return false, fmt.Errorf("unsupported conversion from %T to bool", v)
	}
}
