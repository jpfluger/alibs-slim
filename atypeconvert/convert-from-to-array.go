package atypeconvert

import (
	"fmt"
	"reflect"
	"strconv"
)

// ConvertArrayFromTo attempts to convert an array value to a specified type.
func ConvertArrayFromTo(value interface{}, toType string) (interface{}, error) {
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

	if !isArray {
		return nil, fmt.Errorf("the value type (`%s`) is not an array", valueType)
	}

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

	return nil, fmt.Errorf("unrecognized value type (`%s`)", valueType)
}

///*
//All the array conversions in a single function.
//*/
//func ConvertArrayFromTo(value interface{}, toType string) (interface{}, error) {
//	if toType == "" {
//		return nil, fmt.Errorf("the type paramter is empty")
//	}
//	if value == nil {
//		return value, nil
//	}
//
//	valueType, isArray := DetectType(value)
//
//	if valueType == "" {
//		vType := reflect.TypeOf(value)
//		return nil, fmt.Errorf("the value type (`%v`) is unrecognized", vType)
//	}
//
//	if !isArray {
//		return nil, fmt.Errorf("the value type (`%s`) is not array", valueType)
//	}
//
//	switch valueType {
//	case "obj", "arr":
//		return ConvertArrayFromInterfaceTo(value, toType)
//	case "string":
//		return ConvertArrayFromStringTo(value, toType)
//	case "int":
//		return ConvertArrayFromIntTo(value, toType)
//	case "float":
//		return ConvertArrayFromFloatTo(value, toType)
//	case "bool":
//		return ConvertArrayFromBoolTo(value, toType)
//	}
//
//	return nil, fmt.Errorf("unrecognized value type (`%s`)", valueType)
//}

// ConvertArrayFromInterfaceTo attempts to convert an array of interfaces to a specified type.
func ConvertArrayFromInterfaceTo(value interface{}, toType string) (interface{}, error) {
	if toType == "" {
		return nil, fmt.Errorf("the type paramter is empty")
	}
	if value == nil {
		return value, nil
	}
	switch v := value.(type) {
	case []interface{}:
		if toType == "obj" || toType == "arr" {
			return value, nil
		} else {
			cv, ok := value.([]interface{})
			if !ok {
				return nil, fmt.Errorf("type `%v` failed type-assert", v)
			}
			switch toType {
			case "string":
				carr := []string{}
				for _, cvi := range cv {
					if ncvi, err := ConvertToStringFrom(cvi); err != nil {
						return nil, fmt.Errorf("type `%v` failed ConvertToStringFrom; %v", v, err)
					} else {
						carr = append(carr, ncvi)
					}
				}
				return carr, nil
			case "int":
				carr := []int{}
				for _, cvi := range cv {
					if ncvi, err := ConvertToIntFrom(cvi); err != nil {
						return nil, fmt.Errorf("type `%v` failed ConvertToIntFrom; %v", v, err)
					} else {
						carr = append(carr, ncvi)
					}
				}
				return carr, nil
			case "float":
				carr := []float64{}
				for _, cvi := range cv {
					if ncvi, err := ConvertToFloatFrom(cvi); err != nil {
						return nil, fmt.Errorf("type `%v` failed ConvertToFloatFrom; %v", v, err)
					} else {
						carr = append(carr, ncvi)
					}
				}
				return carr, nil
			case "bool":
				carr := []bool{}
				for _, cvi := range cv {
					if ncvi, err := ConvertToBoolFrom(cvi); err != nil {
						return nil, fmt.Errorf("type `%v` failed ConvertToBoolFrom; %v", v, err)
					} else {
						carr = append(carr, ncvi)
					}
				}
				return carr, nil
			default:
				return nil, fmt.Errorf("type `%v` has unsupported conversion to `%s`", v, toType)
			}
		}
	default:
		return nil, fmt.Errorf("type `%v` has unsupported conversion to `%s`", v, toType)
	}
}

// ConvertArrayFromStringTo attempts to convert an array of strings to a specified type.
func ConvertArrayFromStringTo(value interface{}, toType string) (interface{}, error) {
	if toType == "" {
		return nil, fmt.Errorf("the type paramter is empty")
	}
	if value == nil {
		return value, nil
	}
	switch v := value.(type) {
	case []string:
		if toType == "string" {
			return value, nil
		} else {
			cv, ok := value.([]string)
			if !ok {
				return nil, fmt.Errorf("type `%v` failed type-assert", v)
			}
			switch toType {
			case "obj":
				return nil, fmt.Errorf("type `%v` has unsupported conversion to `%s`", v, toType)
			case "arr":
				return nil, fmt.Errorf("type `%v` has unsupported conversion to `%s`", v, toType)
			case "int":
				carr := []int{}
				for _, val := range cv {
					if vint, err := strconv.Atoi(val); err != nil {
						return false, fmt.Errorf("arr type `%v` failed conversion to `%s`", v, toType)
					} else {
						carr = append(carr, vint)
					}
				}
				return carr, nil
			case "float":
				carr := []float64{}
				for _, val := range cv {
					if vint, err := strconv.ParseFloat(val, 64); err != nil {
						return false, fmt.Errorf("arr type `%v` failed conversion to `%s`", v, toType)
					} else {
						carr = append(carr, vint)
					}
				}
				return carr, nil
			case "bool":
				carr := []bool{}
				for _, val := range cv {
					if vint, err := strconv.ParseBool(val); err != nil {
						return false, fmt.Errorf("arr type `%v` failed conversion to `%s`", v, toType)
					} else {
						carr = append(carr, vint)
					}
				}
				return carr, nil
			default:
				return nil, fmt.Errorf("type `%v` has unsupported conversion to `%s`", v, toType)
			}
		}
	default:
		return nil, fmt.Errorf("type `%v` has unsupported conversion to `%s`", v, toType)
	}
}

// ConvertArrayFromIntTo attempts to convert an array of integers to a specified type.
func ConvertArrayFromIntTo(value interface{}, toType string) (interface{}, error) {
	if toType == "" {
		return nil, fmt.Errorf("the type paramter is empty")
	}
	if value == nil {
		return nil, nil
	}
	var cv []int
	var ok bool
	switch value.(type) {
	case []int:
		cv, ok = value.([]int)
	case []int64:
		var cv64 []int64
		if cv64, ok = value.([]int64); ok {
			for _, cvi := range cv64 {
				cv = append(cv, int(cvi))
			}
		}
	case []int32:
		var cv32 []int32
		if cv32, ok = value.([]int32); ok {
			for _, cvi := range cv32 {
				cv = append(cv, int(cvi))
			}
		}
	case []int16:
		var cv16 []int16
		if cv16, ok = value.([]int16); ok {
			for _, cvi := range cv16 {
				cv = append(cv, int(cvi))
			}
		}
	case []int8:
		var cv8 []int8
		if cv8, ok = value.([]int8); ok {
			for _, cvi := range cv8 {
				cv = append(cv, int(cvi))
			}
		}
	case []uint:
		var cvUI []uint
		if cvUI, ok = value.([]uint); ok {
			for _, cvi := range cvUI {
				cv = append(cv, int(cvi))
			}
		}
	case []uint64:
		var cvU64 []uint64
		if cvU64, ok = value.([]uint64); ok {
			for _, cvi := range cvU64 {
				cv = append(cv, int(cvi))
			}
		}
	case []uint32:
		var cvU32 []uint32
		if cvU32, ok = value.([]uint32); ok {
			for _, cvi := range cvU32 {
				cv = append(cv, int(cvi))
			}
		}
	case []uint16:
		var cvU16 []uint16
		if cvU16, ok = value.([]uint16); ok {
			for _, cvi := range cvU16 {
				cv = append(cv, int(cvi))
			}
		}
	case []uint8:
		var cvU8 []uint8
		if cvU8, ok = value.([]uint8); ok {
			for _, cvi := range cvU8 {
				cv = append(cv, int(cvi))
			}
		}
	}

	if !ok {
		return nil, fmt.Errorf("type `%v` has unsupported conversion to `%s`", cv, toType)
	}

	if toType == "int" {
		return cv, nil
	}

	switch toType {
	case "obj":
		return nil, fmt.Errorf("type `%v` has unsupported conversion to `%s`", reflect.TypeOf(value), toType)
	case "arr":
		return nil, fmt.Errorf("type `%v` has unsupported conversion to `%s`", reflect.TypeOf(value), toType)
	case "string":
		carr := []string{}
		for _, val := range cv {
			carr = append(carr, fmt.Sprintf("%d", val))
		}
		return carr, nil
	case "float":
		carr := []float64{}
		for _, val := range cv {
			carr = append(carr, float64(val))
		}
		return carr, nil
	case "bool":
		carr := []bool{}
		for _, val := range cv {
			carr = append(carr, val != 0)
		}
		return carr, nil
	default:
		return nil, fmt.Errorf("type `%v` has unsupported conversion to `%s`", reflect.TypeOf(value), toType)
	}
}

// ConvertArrayFromFloatTo attempts to convert an array of floats to a specified type.
func ConvertArrayFromFloatTo(value interface{}, toType string) (interface{}, error) {
	if toType == "" {
		return nil, fmt.Errorf("the type paramter is empty")
	}
	if value == nil {
		return nil, nil
	}
	switch v := value.(type) {
	case []float64, []float32:
		cv, ok := value.([]float64)
		if !ok {
			cv = []float64{}
			if cv32, ok := value.([]float32); ok {
				for _, cvi := range cv32 {
					cv = append(cv, float64(cvi))
				}
			} else {
				return nil, fmt.Errorf("type `%v` failed type-assert", v)
			}
		}
		if toType == "float" {
			return value, nil
		} else {
			switch toType {
			case "obj":
				return nil, fmt.Errorf("type `%v` has unsupported conversion to `%s`", v, toType)
			case "arr":
				return nil, fmt.Errorf("type `%v` has unsupported conversion to `%s`", v, toType)
			case "string":
				carr := []string{}
				for _, val := range cv {
					carr = append(carr, fmt.Sprintf("%f", val))
				}
				return carr, nil
			case "int":
				carr := []int{}
				for _, val := range cv {
					carr = append(carr, int(val))
				}
				return carr, nil
			case "bool":
				carr := []bool{}
				for _, val := range cv {
					carr = append(carr, int(val) != 0)
				}
				return carr, nil
			default:
				return nil, fmt.Errorf("type `%v` has unsupported conversion to `%s`", v, toType)
			}
		}
	default:
		return nil, fmt.Errorf("type `%v` has unsupported conversion to `%s`", v, toType)
	}
}

// ConvertArrayFromBoolTo attempts to convert an array of booleans to a specified type.
func ConvertArrayFromBoolTo(value interface{}, toType string) (interface{}, error) {
	if toType == "" {
		return nil, fmt.Errorf("the type paramter is empty")
	}
	if value == nil {
		return value, nil
	}
	switch v := value.(type) {
	case []bool:
		if toType == "bool" {
			return value, nil
		} else {
			cv, ok := value.([]bool)
			if !ok {
				return nil, fmt.Errorf("type `%v` failed type-assert", v)
			}
			switch toType {
			case "obj":
				return nil, fmt.Errorf("type `%v` has unsupported conversion to `%s`", v, toType)
			case "arr":
				return nil, fmt.Errorf("type `%v` has unsupported conversion to `%s`", v, toType)
			case "string":
				carr := []string{}
				for _, val := range cv {
					carr = append(carr, fmt.Sprintf("%t", val))
				}
				return carr, nil
			case "int":
				carr := []int{}
				for _, val := range cv {
					if val {
						carr = append(carr, 1)
					} else {
						carr = append(carr, 0)
					}
				}
				return carr, nil
			case "float":
				carr := []float64{}
				for _, val := range cv {
					if val {
						carr = append(carr, float64(1))
					} else {
						carr = append(carr, float64(0))
					}
				}
				return carr, nil
			default:
				return nil, fmt.Errorf("type `%v` has unsupported conversion to `%s`", v, toType)
			}
		}
	default:
		return nil, fmt.Errorf("type unknown for casted-type of `%v`", v)
	}
}
