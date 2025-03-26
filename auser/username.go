package auser

import (
	"fmt"
	"github.com/jpfluger/alibs-slim/azb"
	"net/mail"
	"strings"
)

// Username represents a username string with associated methods for validation.
type Username string

// IsEmpty checks if the Username is empty after trimming whitespace.
func (e Username) IsEmpty() bool {
	return strings.TrimSpace(string(e)) == ""
}

// TrimSpace trims leading and trailing whitespace from the Username.
func (e Username) TrimSpace() Username {
	return Username(strings.TrimSpace(string(e)))
}

// ToStringTrimLower returns the Username as a trimmed and lowercase string.
func (e Username) ToStringTrimLower() string {
	return strings.ToLower(strings.TrimSpace(string(e)))
}

// IsEmail checks if the Username is a valid email address.
func (e Username) IsEmail() bool {
	_, err := mail.ParseAddress(string(e))
	return err == nil
}

// ExpectEmail checks if the Username contains an "@" symbol, suggesting it's an email.
func (e Username) ExpectEmail() bool {
	return strings.Contains(string(e), "@")
}

// Name extracts the name part of an email address, before the "@" symbol.
func (e Username) Name() string {
	parts := strings.Split(string(e), "@")
	if len(parts) >= 1 {
		return parts[0]
	}
	return ""
}

// Domain extracts the domain part of an email address, after the "@" symbol.
func (e Username) Domain() string {
	parts := strings.Split(string(e), "@")
	if len(parts) == 2 {
		return parts[1]
	}
	return ""
}

// String returns the Username as a standard string.
func (e Username) String() string {
	return string(e)
}

// IsValid checks if the Username is valid according to the provided validation type and function.
func (e Username) IsValid(uvType azb.ZBType, fn FNUsernameValidate) bool {
	return e.IsValidElseError(uvType, fn) == nil
}

// IsValidElseError returns an error if the Username is not valid according to the provided validation type and function.
func (e Username) IsValidElseError(uvType azb.ZBType, fn FNUsernameValidate) error {
	// Helper function to check email format.
	fnCheckEmail := func() error {
		if e.IsEmail() {
			return nil
		}
		return fmt.Errorf("invalid email format")
	}

	// Validate based on the type of username validation required.
	switch uvType {
	case USERNAMEVALIDATETYPE_EMAIL_OR_USER:
		if e.ExpectEmail() {
			return fnCheckEmail()
		}
		fallthrough // If not an email, validate as a regular username.
	case USERNAMEVALIDATETYPE_USER:
		if fn != nil {
			return fn(string(e))
		}
		return ValidateUsername(string(e))
	case USERNAMEVALIDATETYPE_EMAIL:
		return fnCheckEmail()
	default:
		return fmt.Errorf("unknown username validation type")
	}
}

// Constants for different types of username validation.
const (
	USERNAMEVALIDATETYPE_EMAIL_OR_USER = azb.ZBType("email-or-user")
	USERNAMEVALIDATETYPE_EMAIL         = azb.ZBType("email")
	USERNAMEVALIDATETYPE_USER          = azb.ZBType("user")
)

// FNUsernameValidate defines a function type for custom username validation.
type FNUsernameValidate func(target string) error

// ValidateUsername checks if the target string is a valid username according to specific rules.
func ValidateUsername(target string) error {
	if len(target) < 4 || len(target) > 49 {
		return fmt.Errorf("must have 4 to 49 characters")
	}
	if strings.HasPrefix(target, "-") || strings.Contains(target, "--") || strings.HasSuffix(target, "-") {
		return fmt.Errorf("single hyphens allowed but not at the start or end")
	}
	for _, r := range target {
		if !isAlnumOrHyphen(r) {
			return fmt.Errorf("only alphanumeric characters or single hyphens allowed")
		}
	}
	return nil
}

// isAlnumOrHyphen checks if a rune is an ASCII alphanumeric character or a hyphen.
func isAlnumOrHyphen(r rune) bool {
	return ('a' <= r && r <= 'z') || ('A' <= r && r <= 'Z') || ('0' <= r && r <= '9') || r == '-'
}

// Usernames is a slice of Username, allowing for methods that operate on multiple usernames.
type Usernames []Username
