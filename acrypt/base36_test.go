package acrypt

import "testing"

func TestEncodeBase36(t *testing.T) {
	tests := []struct {
		input    int
		expected string
	}{
		{0, "0"},
		{1, "1"},
		{35, "Z"},
		{36, "10"},
		{12345, "9IX"},
	}
	for _, test := range tests {
		result := EncodeBase36(test.input)
		if result != test.expected {
			t.Errorf("EncodeBase36(%d) = %s, expected %s", test.input, result, test.expected)
		}
	}
}

func TestDecodeBase36(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"0", 0},
		{"1", 1},
		{"Z", 35},
		{"10", 36},
		{"9IX", 12345},
	}
	for _, test := range tests {
		result, err := DecodeBase36(test.input)
		if err != nil || result != test.expected {
			t.Errorf("DecodeBase36(%s) = %d, expected %d", test.input, result, test.expected)
		}
	}
}
