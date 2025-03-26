package anode

import "time"

// AdminLock represents the lock status of an admin account.
type AdminLock struct {
	// IsPasswordLocked indicates whether the password is locked.
	// This can be set to true in two scenarios:
	// 1. The admin manually sets this value to true when managing the user account.
	// 2. The password is empty or the account requires a password reset.
	IsPasswordLocked bool `json:"isPasswordLocked,omitempty"`

	// Date is the actual date when the account was locked.
	Date *time.Time `json:"lockedDate,omitempty"`

	// Message provides the reason for the account lock.
	Message string `json:"message,omitempty"`

	// RequestResetPassword is the date when the admin requested the user to create/reset the password.
	// If this is set, the user is forced to set the password on login, regardless of the IsPasswordLocked value.
	// If IsPasswordLocked is true and RequestResetPassword is valid, the password can only be reset via an email link,
	// whether initiated by the admin or through the forgot-login process.
	RequestResetPassword *time.Time `json:"requestResetPassword,omitempty"`
}
