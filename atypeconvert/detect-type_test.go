package atypeconvert

import (
	"testing"
)

// TestDetectType tests the DetectType function for various types and arrays.
func TestDetectType(t *testing.T) {
	tests := []struct {
		input         interface{}
		expectedType  string
		expectedArray bool
	}{
		{"hello", "string", false},
		{123, "int", false},
		{123.456, "float", false},
		{true, "bool", false},
		{map[string]interface{}{"key": "value"}, "obj", false},
		{[]interface{}{1, 2, 3}, "arr", true},
		{[]string{"a", "b", "c"}, "string", true},
		{[]int{1, 2, 3}, "int", true},
		{[]float64{1.1, 2.2, 3.3}, "float", true},
		{[]bool{true, false, true}, "bool", true},
		{nil, "", false},
		{(*string)(nil), "string", false},
	}

	for _, test := range tests {
		gotType, gotArray := DetectType(test.input)
		if gotType != test.expectedType || gotArray != test.expectedArray {
			t.Errorf("DetectType(%v) = %v, %v; want %v, %v", test.input, gotType, gotArray, test.expectedType, test.expectedArray)
		}
	}
}
