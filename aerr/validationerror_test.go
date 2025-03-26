package aerr

import (
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/jpfluger/alibs-slim/acrypt"
	"strings"
	"testing"
)

// TestValidationError_Error checks if the Error method returns the correct message.
func TestValidationError_Error(t *testing.T) {
	ve := ValidationError{Message: "field must be unique"}
	if ve.Error() != "field must be unique" {
		t.Errorf("Error() = %v, want %v", ve.Error(), "field must be unique")
	}
}

// TestValidationError_ErrorLowercase checks if the ErrorLowercase method returns the message in lowercase.
func TestValidationError_ErrorLowercase(t *testing.T) {
	ve := ValidationError{Message: "Field Must Be Unique"}
	if ve.ErrorLowercase() != "field must be unique" {
		t.Errorf("ErrorLowercase() = %v, want %v", ve.ErrorLowercase(), "field must be unique")
	}
}

// TestValidationError_GetSysError checks if GetSysError method returns the correct system error.
func TestValidationError_GetSysError(t *testing.T) {
	sysErr := errors.New("internal system error")
	ve := ValidationError{Message: "field must be unique", SysError: sysErr}
	if ve.GetSysError() != sysErr {
		t.Errorf("GetSysError() = %v, want %v", ve.GetSysError(), sysErr)
	}

	ve.SysError = nil
	if ve.GetSysError().Error() != "field must be unique" {
		t.Errorf("GetSysError() = %v, want %v", ve.GetSysError(), "field must be unique")
	}
}

// TestValidationError_MarshalJSON checks if the MarshalJSON method excludes the SysError field.
func TestValidationError_MarshalJSON(t *testing.T) {
	ve := ValidationError{Message: "field must be unique", Field: "username", Tag: "required"}
	bytes, err := json.Marshal(ve)
	if err != nil {
		t.Fatal(err)
	}
	jsonStr := string(bytes)
	if strings.Contains(jsonStr, "SysError") {
		t.Errorf("MarshalJSON() should not include SysError, got %v", jsonStr)
	}
}

// TestValidationErrors_Add checks if the Add method correctly appends a new ValidationError.
func TestValidationErrors_Add(t *testing.T) {
	ves := ValidationErrors{}
	ve := ValidationError{Message: "field must be unique"}
	ves.Add(&ve)
	if len(ves) != 1 || ves[0] != &ve {
		t.Errorf("Add() did not append ValidationError correctly, got %v", ves)
	}
}

// TestValidationErrors_Error checks if the Error method returns a concatenated message of all validation errors.
func TestValidationErrors_Error(t *testing.T) {
	ves := ValidationErrors{
		&ValidationError{Message: "field must be unique"},
		&ValidationError{Message: "field is required"},
	}
	want := "field must be unique; field is required"
	if ves.Error() != want {
		t.Errorf("Error() = %v, want %v", ves.Error(), want)
	}
}

// TestValidationErrors_MarshalJSON checks if the MarshalJSON method provides a clean error array.
func TestValidationErrors_MarshalJSON(t *testing.T) {
	ves := ValidationErrors{
		&ValidationError{Message: "field must be unique"},
		&ValidationError{Message: "field is required"},
	}
	bytes, err := json.Marshal(ves)
	if err != nil {
		t.Fatal(err)
	}
	jsonStr := string(bytes)
	if !strings.Contains(jsonStr, "field must be unique") || !strings.Contains(jsonStr, "field is required") {
		t.Errorf("MarshalJSON() did not return a clean error array, got %v", jsonStr)
	}
}

// Mock custom password validator for testing
func mockCustomValidator(password string) error {
	if password == "weakPassword!" {
		return errors.New("password is too weak")
	}
	return nil
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name            string
		newPassword     string
		verifyPassword  string
		customValidator acrypt.FNValidatePassword
		expectedError   *ValidationError
	}{
		{
			name:           "Empty new password",
			newPassword:    "",
			verifyPassword: "somePassword",
			expectedError: &ValidationError{
				Message: "New password is empty.",
				Field:   "newPassword",
				Tag:     "empty_password",
			},
		},
		{
			name:           "Empty verify password",
			newPassword:    "somePassword",
			verifyPassword: "",
			expectedError: &ValidationError{
				Message: "Verify password is empty.",
				Field:   "verifyPassword",
				Tag:     "empty_password",
			},
		},
		{
			name:           "Passwords do not match",
			newPassword:    "somePassword",
			verifyPassword: "differentPassword",
			expectedError: &ValidationError{
				Message: "Verify password does not match new password.",
				Field:   "verifyPassword",
				Tag:     "password_mismatch",
			},
		},
		{
			name:           "Password fails complexity requirements",
			newPassword:    "password", // Missing uppercase, digit, and special character
			verifyPassword: "password",
			expectedError: &ValidationError{
				Message: "password must contain at least one uppercase letter",
				Field:   "newPassword",
				Tag:     "complexity",
			},
		},
		{
			name:            "Weak password based on custom validator",
			newPassword:     "weakPassword!",
			verifyPassword:  "weakPassword!",
			customValidator: mockCustomValidator,
			expectedError: &ValidationError{
				Message: "password is too weak",
				Field:   "newPassword",
				Tag:     "custom",
			},
		},
		{
			name:           "Valid strong password",
			newPassword:    "StrongPassword123!",
			verifyPassword: "StrongPassword123!",
			expectedError:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePassword(tt.newPassword, tt.verifyPassword, tt.customValidator)

			if tt.expectedError != nil {
				// Ensure an error is returned and compare each field
				assert.NotNil(t, err)
				assert.Equal(t, tt.expectedError.Message, err.Message)
				assert.Equal(t, tt.expectedError.Field, err.Field)
				assert.Equal(t, tt.expectedError.Tag, err.Tag)
			} else {
				// Expecting no error
				assert.Nil(t, err)
			}
		})
	}
}

//func TestEvaluatePasswordStrengthAndErrors(t *testing.T) {
//	tests := []struct {
//		name            string
//		password        string
//		customValidator acrypt.FNValidatePassword
//		expectedScore   int
//		expectedErrors  []string
//	}{
//		{
//			name:          "Password too short",
//			password:      "Short1!",
//			expectedScore: 0,
//			expectedErrors: []string{
//				"Key: 'PasswordRequirements.Password' Error:Field validation for 'Password' failed on the 'min' tag",
//				"Password strength is too weak: score 0 (minimum required score is 3)",
//			},
//		},
//		{
//			name:          "Missing complexity (no uppercase)",
//			password:      "password123!",
//			expectedScore: 0, // `zxcvbn` may classify this as score 0 if it’s weak
//			expectedErrors: []string{
//				"password must contain at least one uppercase letter",
//				"Password strength is too weak: score 0 (minimum required score is 3)",
//			},
//		},
//		{
//			name:            "Custom validation failure",
//			password:        "ValidPass123!",
//			customValidator: mockCustomValidator,
//			expectedScore:   1, // Adjusted due to conservative `zxcvbn` scoring
//			expectedErrors:  []string{"password is too weak"},
//		},
//		{
//			name:           "Valid strong password",
//			password:       "ValidPass123!",
//			expectedScore:  1, // Adjusted due to `zxcvbn` possibly scoring this lower
//			expectedErrors: nil,
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			score, validationErrors := EvaluatePasswordStrengthAndErrors(tt.password, tt.customValidator)
//
//			// Check the expected score with an allowance for `zxcvbn` conservative scoring
//			if tt.expectedScore > score && score < 3 {
//				assert.Contains(t, validationErrors.Error(), "Password strength is too weak", "Expected password to be stronger")
//			} else {
//				assert.Equal(t, tt.expectedScore, score, "Unexpected strength score")
//			}
//
//			// Check for expected errors
//			if tt.expectedErrors == nil {
//				assert.Nil(t, validationErrors, "Expected no errors")
//			} else {
//				assert.NotNil(t, validationErrors, "Expected errors but got nil")
//				if validationErrors != nil {
//					for _, expectedErr := range tt.expectedErrors {
//						assert.Contains(t, validationErrors.Error(), expectedErr)
//					}
//				}
//			}
//		})
//	}
//}
