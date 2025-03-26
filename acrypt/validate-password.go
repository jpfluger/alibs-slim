package acrypt

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/nbutton23/zxcvbn-go"
	"unicode"
)

// PasswordRequirements struct defines the requirements for a valid password.
type PasswordRequirements struct {
	Password string `validate:"required,min=8,max=100"`
}

// FNValidatePassword is a function type that takes a password string and returns an error.
type FNValidatePassword func(target string) error

// ValidatePasswordWithScore validates a password with basic requirements and checks its strength score.
func ValidatePasswordWithScore(password string) (int, error) {
	validate := validator.New()
	reqs := &PasswordRequirements{Password: password}

	// Check basic length requirements
	err := validate.Struct(reqs)
	if err != nil {
		return 0, err
	}

	// Check password strength using zxcvbn
	strength := zxcvbn.PasswordStrength(password, nil)
	if strength.Score < 3 {
		return strength.Score, fmt.Errorf("password is too weak: score %d (min required score is 3)", strength.Score)
	}

	return strength.Score, nil
}

// ValidatePasswordWithOptions validates a password with additional custom checks.
func ValidatePasswordWithOptions(password string, customValidator FNValidatePassword) (int, error) {
	score, err := ValidatePasswordWithScore(password)
	if err != nil {
		return score, err
	}

	if customValidator != nil {
		if err := customValidator(password); err != nil {
			return score, err
		}
	}

	return score, nil
}

// ValidatePasswordComplex checks if the password meets the specified criteria.
func ValidatePasswordComplex(password string) error {
	validate := validator.New()
	passwordValidator := PasswordRequirements{Password: password}

	// Validate basic length requirements using struct tags
	err := validate.Struct(passwordValidator)
	if err != nil {
		return err
	}

	// Custom validation for character types
	if err = validateCustomCharacterTypes(password); err != nil {
		return err
	}

	return nil
}

// validateCustomCharacterTypes checks for at least one lowercase, uppercase, digit, and special character.
func validateCustomCharacterTypes(password string) error {
	var hasLower, hasUpper, hasNumber, hasSpecial bool
	for _, char := range password {
		switch {
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsDigit(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasLower {
		return fmt.Errorf("password must contain at least one lowercase letter")
	}
	if !hasUpper {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}
	if !hasNumber {
		return fmt.Errorf("password must contain at least one digit")
	}
	if !hasSpecial {
		return fmt.Errorf("password must contain at least one special character")
	}

	return nil
}

// ValidatePasswordComplexWithOptions checks if the password meets the specified criteria.
func ValidatePasswordComplexWithOptions(password string, customValidator FNValidatePassword) error {
	err := ValidatePasswordComplex(password)
	if err != nil {
		return err
	}

	if customValidator != nil {
		if err := customValidator(password); err != nil {
			return err
		}
	}

	return nil
}
