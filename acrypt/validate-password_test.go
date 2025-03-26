package acrypt

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestValidatePasswordWithScore tests the password validation with score.
func TestValidatePasswordWithScore(t *testing.T) {
	// Test a strong password
	password := "ThisIsAVeryStrongPass1!"
	score, err := ValidatePasswordWithScore(password)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, score, 3)

	// Test a weak password
	password = "weak"
	score, err = ValidatePasswordWithScore(password)
	assert.Error(t, err)
	assert.Less(t, score, 3)
}

// TestValidatePasswordWithOptions tests the password validation with custom options.
func TestValidatePasswordWithOptions(t *testing.T) {
	// Test a strong password with custom validator
	password := "ThisIsAVeryStrongPass1!"
	score, err := ValidatePasswordWithOptions(password, ValidatePasswordComplex)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, score, 3)

	// Test a weak password with custom validator
	password = "weak"
	score, err = ValidatePasswordWithOptions(password, ValidatePasswordComplex)
	assert.Error(t, err)
}

// TestValidatePasswordComplex tests the complex password validation.
func TestValidatePasswordComplex(t *testing.T) {
	// Test a password that meets the complexity requirements
	password := "StrongPass1!"
	err := ValidatePasswordComplex(password)
	assert.NoError(t, err)

	// Test a password that does not meet the complexity requirements
	password = "weak"
	err = ValidatePasswordComplex(password)
	assert.Error(t, err)
}

// TestValidateCustomCharacterTypes tests the custom character type validation.
func TestValidateCustomCharacterTypes(t *testing.T) {
	// Test a password that meets the character type requirements
	password := "StrongPass1!"
	err := validateCustomCharacterTypes(password)
	assert.NoError(t, err)

	// Test a password that does not meet the character type requirements
	password = "weak"
	err = validateCustomCharacterTypes(password)
	assert.Error(t, err)
}

// TestValidatePasswordComplexWithOptions tests the complex password validation with custom options.
func TestValidatePasswordComplexWithOptions(t *testing.T) {
	// Test a strong password with custom validator
	password := "StrongPass1!"
	err := ValidatePasswordComplexWithOptions(password, ValidatePasswordComplex)
	assert.NoError(t, err)

	// Test a weak password with custom validator
	password = "weak"
	err = ValidatePasswordComplexWithOptions(password, ValidatePasswordComplex)
	assert.Error(t, err)
}
