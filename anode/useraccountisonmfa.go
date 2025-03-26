package anode

import "time"

// UserAccountIsOnMFA represents the Multi-Factor Authentication (MFA) status of a user account.
type UserAccountIsOnMFA struct {
	// IsOn indicates whether MFA is enabled for the user account.
	IsOn bool `json:"isOn,omitempty"`

	// Created is the timestamp when MFA was enabled.
	Created *time.Time `json:"created,omitempty"`

	// Verified is the timestamp when MFA was verified.
	Verified *time.Time `json:"verified,omitempty"`
}

// IsVerified checks if the MFA has been verified.
// Returns true if the Verified timestamp is set and not zero.
func (uat *UserAccountIsOnMFA) IsVerified() bool {
	return uat != nil && uat.Verified != nil && !uat.Verified.IsZero()
}
