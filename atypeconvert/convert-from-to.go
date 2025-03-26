package atypeconvert

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

// ConvertFromTo attempts to convert a value to a specified type.
func ConvertFromTo(value interface{}, toType string) (interface{}, error) {
	if toType == "" {
		return nil, fmt.Errorf("the type parameter is empty")
	}
	if value == nil {
		return value, nil
	}

	valueType, isArray := DetectType(value)

	if valueType == "" {
		vType := reflect.TypeOf(value)
		return nil, fmt.Errorf("the value type (`%v`) is unrecognized", vType)
	}

	if isArray {
		switch valueType {
		case "obj", "arr":
			return ConvertArrayFromInterfaceTo(value, toType)
		case "string":
			return ConvertArrayFromStringTo(value, toType)
		case "int":
			return ConvertArrayFromIntTo(value, toType)
		case "float":
			return ConvertArrayFromFloatTo(value, toType)
		case "bool":
			return ConvertArrayFromBoolTo(value, toType)
		}
	} else {
		switch valueType {
		case "obj", "arr":
			return value, nil
		case "string":
			return ConvertFromStringTo(value, toType)
		case "int":
			return ConvertFromIntTo(value, toType)
		case "float":
			return ConvertFromFloatTo(value, toType)
		case "bool":
			return ConvertFromBoolTo(value, toType)
		}
	}

	return nil, fmt.Errorf("unsupported value type (`%s`)", valueType)
}

// ConvertFromStringTo attempts to convert a string value to a specified type.
func ConvertFromStringTo(value interface{}, toType string) (interface{}, error) {
	if toType == "" {
		return nil, fmt.Errorf("the type parameter is empty")
	}
	if value == nil {
		return value, nil
	}

	switch v := value.(type) {
	case string:
		return convertString(v, toType)
	case *string:
		if v == nil {
			return nil, nil
		}
		return convertString(*v, toType)
	default:
		return nil, fmt.Errorf("type `%v` has unsupported conversion to `%s`", v, toType)
	}
}

// convertString is a helper function to handle string conversions.
func convertString(cv string, toType string) (interface{}, error) {
	switch toType {
	case "obj":
		var obj map[string]interface{}
		if err := json.Unmarshal([]byte(cv), &obj); err != nil {
			return nil, fmt.Errorf("json unmarshal failed: %v", err)
		}
		return obj, nil
	case "arr":
		var arr []interface{}
		if err := json.Unmarshal([]byte(cv), &arr); err != nil {
			return nil, fmt.Errorf("json unmarshal failed: %v", err)
		}
		return arr, nil
	case "int":
		return strconv.Atoi(cv)
	case "float":
		return strconv.ParseFloat(cv, 64)
	case "bool":
		return strconv.ParseBool(cv)
	default:
		return nil, fmt.Errorf("unsupported conversion to `%s`", toType)
	}
}

// ConvertFromIntTo attempts to convert an integer value to a specified type.
func ConvertFromIntTo(value interface{}, toType string) (interface{}, error) {
	if toType == "" {
		return nil, fmt.Errorf("the type parameter is empty")
	}
	if value == nil {
		return nil, nil
	}

	cv, err := convertToInt(value)
	if err != nil {
		return nil, err
	}

	switch toType {
	case "string":
		return fmt.Sprintf("%d", cv), nil
	case "float":
		return float64(cv), nil
	case "bool":
		return cv != 0, nil
	default:
		return nil, fmt.Errorf("unsupported conversion to `%s`", toType)
	}
}

// convertToInt is a helper function to handle integer conversions.
func convertToInt(value interface{}) (int, error) {
	switch v := value.(type) {
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
	case *int:
		if v == nil {
			return 0, nil
		}
		return *v, nil
	case *int64:
		if v == nil {
			return 0, nil
		}
		return int(*v), nil
	case *int32:
		if v == nil {
			return 0, nil
		}
		return int(*v), nil
	case *int16:
		if v == nil {
			return 0, nil
		}
		return int(*v), nil
	case *int8:
		if v == nil {
			return 0, nil
		}
		return int(*v), nil
	case *uint:
		if v == nil {
			return 0, nil
		}
		return int(*v), nil
	case *uint64:
		if v == nil {
			return 0, nil
		}
		return int(*v), nil
	case *uint32:
		if v == nil {
			return 0, nil
		}
		return int(*v), nil
	case *uint16:
		if v == nil {
			return 0, nil
		}
		return int(*v), nil
	case *uint8:
		if v == nil {
			return 0, nil
		}
		return int(*v), nil
	default:
		return 0, fmt.Errorf("type `%v` is not an integer", v)
	}
}

// ConvertFromFloatTo attempts to convert a float value to a specified type.
func ConvertFromFloatTo(value interface{}, toType string) (interface{}, error) {
	if toType == "" {
		return nil, fmt.Errorf("the type parameter is empty")
	}
	if value == nil {
		return nil, nil
	}

	cv, err := convertToFloat(value)
	if err != nil {
		return nil, err
	}

	switch toType {
	case "string":
		return fmt.Sprintf("%f", cv), nil
	case "int":
		return int(cv), nil
	case "bool":
		return int(cv) != 0, nil
	default:
		return nil, fmt.Errorf("unsupported conversion to `%s`", toType)
	}
}

// convertToFloat is a helper function to handle float conversions.
func convertToFloat(value interface{}) (float64, error) {
	switch v := value.(type) {
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	case *float64:
		if v == nil {
			return 0, nil
		}
		return *v, nil
	case *float32:
		if v == nil {
			return 0, nil
		}
		return float64(*v), nil
	default:
		return 0, fmt.Errorf("type `%v` is not a float", v)
	}
}

// ConvertFromBoolTo attempts to convert a boolean value to a specified type.
func ConvertFromBoolTo(value interface{}, toType string) (interface{}, error) {
	if toType == "" {
		return nil, fmt.Errorf("the type parameter is empty")
	}
	if value == nil {
		return value, nil
	}

	cv, err := convertToBool(value)
	if err != nil {
		return nil, err
	}

	switch toType {
	case "string":
		return fmt.Sprintf("%t", cv), nil
	case "int":
		if cv {
			return 1, nil
		}
		return 0, nil
	case "float":
		if cv {
			return float64(1), nil
		}
		return float64(0), nil
	default:
		return nil, fmt.Errorf("unsupported conversion to `%s`", toType)
	}
}

// convertToBool is a helper function to handle boolean conversions.
func convertToBool(value interface{}) (bool, error) {
	switch v := value.(type) {
	case bool:
		return v, nil
	case *bool:
		if v == nil {
			return false, nil
		}
		return *v, nil
	default:
		return false, fmt.Errorf("type `%v` is not a boolean", v)
	}
}
