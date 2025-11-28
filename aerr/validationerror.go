package aerr

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/jpfluger/alibs-slim/acrypt"
	"github.com/nbutton23/zxcvbn-go"
)

// ValidationError represents an error that occurs during validation of input.
// It includes a human-readable message, the field that caused the error,
// the validation tag that failed, and an optional system error for internal use.
type ValidationError struct {
	Message  string `json:"message,omitempty"` // User-friendly error message.
	Field    string `json:"field,omitempty"`   // The input field associated with the error.
	Tag      string `json:"tag,omitempty"`     // The validation rule that was violated.
	SysError error  `json:"-"`                 // System error, not to be sent to the client.
}

// Error returns the error message.
func (ve *ValidationError) Error() string {
	return ve.Message
}

// ErrorLowercase returns the error message in lowercase.
func (ve *ValidationError) ErrorLowercase() string {
	return strings.ToLower(ve.Message)
}

// GetSysError returns the system error if present; otherwise, it returns a new error based on the message.
func (ve *ValidationError) GetSysError() error {
	if ve.SysError != nil {
		return ve.SysError
	}
	return errors.New(ve.ErrorLowercase())
}

// MarshalJSON customizes the JSON marshaling to exclude the SysError field.
func (ve *ValidationError) MarshalJSON() ([]byte, error) {
	type Alias ValidationError
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(ve),
	})
}

// ValidationErrors is a slice of ValidationError pointers, used to aggregate multiple validation errors.
type ValidationErrors []*ValidationError

// Add appends a new ValidationError to the slice.
func (ves *ValidationErrors) Add(ve *ValidationError) {
	*ves = append(*ves, ve)
}

// Error implements the error interface for ValidationErrors.
// It returns a concatenated message of all validation errors.
func (ves ValidationErrors) Error() string {
	var messages []string
	for _, ve := range ves {
		messages = append(messages, ve.Error())
	}
	return strings.Join(messages, "; ")
}

// MarshalJSON customizes the JSON marshaling for ValidationErrors to provide a clean error array.
func (ves ValidationErrors) MarshalJSON() ([]byte, error) {
	return json.Marshal([]*ValidationError(ves))
}

// ValidatePassword performs validation on a new password and its confirmation.
// It checks for non-empty values, matching passwords, and uses fnValidatePassword
// to enforce password strength requirements. It returns a ValidationError if any
// validation fails.
func ValidatePassword(newPassword, verifyPassword string, fnValidatePassword acrypt.FNValidatePassword) *ValidationError {

	// Trim whitespace from inputs
	newPassword = strings.TrimSpace(newPassword)
	verifyPassword = strings.TrimSpace(verifyPassword)

	// Check for empty newPassword
	if newPassword == "" {
		return &ValidationError{
			Message: "New password is empty.",
			Field:   "newPassword",
			Tag:     "empty_password",
		}
	}

	// Check for empty verifyPassword
	if verifyPassword == "" {
		return &ValidationError{
			Message: "Verify password is empty.",
			Field:   "verifyPassword",
			Tag:     "empty_password",
		}
	}

	// Check for password match
	if verifyPassword != newPassword {
		return &ValidationError{
			Message: "Verify password does not match new password.",
			Field:   "verifyPassword",
			Tag:     "password_mismatch",
		}
	}

	// Use default password validator if none is provided
	isCustomFunc := true
	if fnValidatePassword == nil {
		isCustomFunc = false
		fnValidatePassword = acrypt.ValidatePasswordComplex
	}

	// Validate password strength with the custom or default validator
	if err := fnValidatePassword(newPassword); err != nil {
		tag := "complexity"
		if isCustomFunc {
			tag = "custom"
		}
		return &ValidationError{
			Message: err.Error(),
			Field:   "newPassword",
			Tag:     tag,
		}
	}

	// Password validated successfully
	return nil
}

// EvaluatePasswordStrengthAndErrors returns both the password strength score and any validation errors.
func EvaluatePasswordStrengthAndErrors(password string, customValidator acrypt.FNValidatePassword) (int, *ValidationErrors) {
	validationErrors := ValidationErrors{}

	// Step 1: Basic length and structure requirements
	if err := acrypt.ValidatePasswordComplex(password); err != nil {
		validationErrors.Add(&ValidationError{
			Message: err.Error(),
			Field:   "password",
			Tag:     "complexity",
		})
	}

	// Step 2: Check password strength score with zxcvbn
	strength := zxcvbn.PasswordStrength(password, nil)
	if strength.Score < 3 {
		validationErrors.Add(&ValidationError{
			Message: fmt.Sprintf("Password strength is too weak: score %d (minimum required score is 3)", strength.Score),
			Field:   "password",
			Tag:     "strength",
		})
	}

	// Step 3: Run custom validation if provided
	if customValidator != nil {
		if err := customValidator(password); err != nil {
			validationErrors.Add(&ValidationError{
				Message: err.Error(),
				Field:   "password",
				Tag:     "custom",
			})
		}
	}

	// Return the score and any accumulated validation errors
	if len(validationErrors) > 0 {
		return strength.Score, &validationErrors
	}
	return strength.Score, nil
}
