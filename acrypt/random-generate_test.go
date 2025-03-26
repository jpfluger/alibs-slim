package acrypt

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestRandGenerate4Digits tests the RandGenerate4Digits function.
func TestRandGenerate4Digits(t *testing.T) {
	pass, err := RandGenerate4Digits()
	assert.NoError(t, err)
	assert.Len(t, pass, 4)

	pass, err = RandGenerate4Digits()
	if err != nil {
		t.Errorf("randGenerate4Digits should not return an error: %v", err)
	}
	if len(pass) != 4 {
		t.Errorf("randGenerate4Digits should return 4 digits, got: %s", pass)
	}
	if !regexp.MustCompile(`^\d{4}$`).MatchString(pass) {
		t.Errorf("randGenerate4Digits should return only digits, got: %s", pass)
	}
}

// TestRandGenerate16 tests the RandGenerate16 function.
func TestRandGenerate16(t *testing.T) {
	pass, err := RandGenerate16()
	assert.NoError(t, err)
	assert.Len(t, pass, 16)
}

// TestRandGenerate20 tests the RandGenerate20 function.
func TestRandGenerate20(t *testing.T) {
	pass, err := RandGenerate20()
	assert.NoError(t, err)
	assert.Len(t, pass, 20)
}

// TestRandGenerate32 tests the RandGenerate32 function.
func TestRandGenerate32(t *testing.T) {
	pass, err := RandGenerate32()
	assert.NoError(t, err)
	assert.Len(t, pass, 32)
}

// TestRandGenerate64 tests the RandGenerate64 function.
func TestRandGenerate64(t *testing.T) {
	pass, err := RandGenerate64()
	assert.NoError(t, err)
	assert.Len(t, pass, 64)
}

// TestRandGenerate tests the RandGenerate function with custom parameters.
func TestRandGenerate(t *testing.T) {
	length := 10
	numDigits := 2
	numSymbols := 2
	pass, err := RandGenerate(length, numDigits, numSymbols, false, false)
	assert.NoError(t, err)
	assert.Len(t, pass, length)
}

// TestRandomTextGenerator tests the Generate method of the RandomTextGenerator struct.
func TestRandomTextGenerator(t *testing.T) {
	rg := RandomTextGenerator{
		Length:      10,
		NumDigits:   2,
		NumSymbols:  2,
		NoUpper:     false,
		AllowRepeat: false,
	}
	pass, err := rg.Generate()
	assert.NoError(t, err)
	assert.Len(t, pass, rg.Length)
}

// TestTryRandGenerate4Digits tests the TryRandGenerate4Digits function.
func TestTryRandGenerate4Digits(t *testing.T) {
	pass := TryRandGenerate4Digits()
	assert.NotEmpty(t, pass)
}

// TestMustRandGenerate4Digits tests the MustRandGenerate4Digits function.
func TestMustRandGenerate4Digits(t *testing.T) {
	assert.NotPanics(t, func() {
		pass := MustRandGenerate4Digits()
		assert.Len(t, pass, 4)
	})
}

// TestGenerateRandomInt100KTo1B ensures the function generates numbers within the expected range.
func TestGenerateRandomInt100KTo1B(t *testing.T) {
	for i := 0; i < 1000; i++ { // Run multiple times to check randomness
		num := GenerateRandomInt100KTo1B()
		if num < 100000 || num > 999999999 {
			t.Errorf("Generated number %d is out of expected range (100000 - 999999999)", num)
		}
	}
}

// TestGenerateRandomInt100KTo1M ensures numbers are within the range 100,000 - 999,999.
func TestGenerateRandomInt100KTo1M(t *testing.T) {
	min, max := 100000, 999999

	for i := 0; i < 1000; i++ { // Run multiple times to check randomness
		num := GenerateRandomInt100KTo1M()
		if num < min || num > max {
			t.Errorf("Generated number %d is out of expected range (%d - %d)", num, min, max)
		}
	}
}

// TestGenerateRandomIntWithOptions ensures numbers are within the specified range.
func TestGenerateRandomIntWithOptions(t *testing.T) {
	tests := []struct {
		min         int
		max         int
		expectedMin int
		expectedMax int
	}{
		{100, 200, 100, 200},    // Normal range
		{500, 1000, 500, 1000},  // Normal range
		{-50, 50, 0, 50},        // Min is negative, should adjust to 0
		{100, 100, 100, 1099},   // Max == Min, should adjust max to min + 999
		{1000, 500, 1000, 1999}, // Max < Min, should adjust max to min + 999
	}

	for _, test := range tests {
		for i := 0; i < 100; i++ { // Run multiple times for randomness check
			num := GenerateRandomIntWithOptions(test.min, test.max)
			if num < test.expectedMin || num > test.expectedMax {
				t.Errorf("Generated number %d is out of expected range (%d - %d) for min=%d, max=%d",
					num, test.expectedMin, test.expectedMax, test.min, test.max)
			}
		}
	}
}
