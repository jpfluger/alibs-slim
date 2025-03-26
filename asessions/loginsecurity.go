package asessions

import (
	"encoding/json"
	"fmt"
	"github.com/jpfluger/alibs-slim/atags"
	"github.com/jpfluger/alibs-slim/auser"
	"github.com/jpfluger/alibs-slim/azb"
	"time"
)

// LOGIN_SECURITY is a constant tag key used for managing login security options.
const LOGIN_SECURITY = atags.TagKey("lsec")

// LoginSecurity struct holds information about a user's login security.
type LoginSecurity struct {
	Usernames   auser.Usernames `json:"username,omitempty"`    // Slice of usernames associated with the security.
	HasPassword bool            `json:"hasPassword,omitempty"` // Indicates if a password exists.
	FormDate    *time.Time      `json:"formDate,omitempty"`    // Date when the form was filled.
}

// NewLoginSecurityFromJSON creates a new LoginSecurity instance from a JSON string.
func NewLoginSecurityFromJSON(target string) (*LoginSecurity, error) {
	lsec := NewLoginSecurity()
	if target != "" {
		if err := json.Unmarshal([]byte(target), lsec); err != nil {
			return nil, fmt.Errorf("failed to unmarshal JSON for %s: %v", LOGIN_SECURITY, err)
		}
	}
	return lsec, nil
}

// NewLoginSecurity creates a new instance of LoginSecurity with default values.
func NewLoginSecurity() *LoginSecurity {
	return &LoginSecurity{
		Usernames: auser.Usernames{},
	}
}

// AddUsername adds a new username to the LoginSecurity.
func (ls *LoginSecurity) AddUsername(username auser.Username) {
	ls.Usernames = append(ls.Usernames, username)
}

// IsExcessiveUsernamesAttempts checks if the number of username attempts exceeds the maximum allowed.
func (ls *LoginSecurity) IsExcessiveUsernamesAttempts(max int) error {
	if max <= 1 {
		max = 10 // Default to 10 if an invalid max value is provided.
	}
	if len(ls.Usernames) > max {
		return fmt.Errorf("number of usernames exceeded the maximum limit of %d", max)
	}
	return nil
}

// IsUsernameValid checks if the provided username is valid based on default validation rules.
func (ls *LoginSecurity) IsUsernameValid(username auser.Username) error {
	return ls.IsUsernameValidWithOptions(username, auser.USERNAMEVALIDATETYPE_EMAIL_OR_USER, nil)
}

// IsUsernameValidWithOptions checks if the provided username is valid based on specified validation rules and custom function.
func (ls *LoginSecurity) IsUsernameValidWithOptions(username auser.Username, uvType azb.ZBType, fn auser.FNUsernameValidate) error {
	if !username.IsValid(uvType, fn) {
		return fmt.Errorf("invalid username")
	}
	return nil
}
