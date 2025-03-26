package anode

import (
	"github.com/jpfluger/alibs-slim/acontact"
	"github.com/jpfluger/alibs-slim/aemail"
	"github.com/jpfluger/alibs-slim/alegal"
	"github.com/jpfluger/alibs-slim/asessions"
	"github.com/jpfluger/alibs-slim/atime"
	"strings"
	"time"
)

// UserAccount represents a user's account with various attributes and settings.
type UserAccount struct {
	// IsDeactivated indicates whether the user account is deactivated.
	IsDeactivated bool `json:"isDeactivated"`

	// MFA contains Multi-Factor Authentication settings.
	MFA struct {
		// TOTP holds the current TOTP (Time-based One-Time Password) configuration.
		TOTP UserAccountIsOnMFA `json:"totp,omitempty"`
		// TOTPNew holds the new TOTP configuration when regenerating a TOTP key.
		// This keeps the original TOTP in place to prevent a security leak.
		TOTPNew *UserAccountIsOnMFA `json:"totpNew,omitempty"`
	} `json:"mfa,omitempty"`

	// Email is required.
	// Depending on the site operation, email could parallel the username.
	// For Username, see UserVault.Credentials.Username.
	Email aemail.EmailAddress `json:"email,omitempty"`

	// Phone is optional but could be required based on the system implementing this struct.
	// Depending on the site operation, email could parallel the username.
	Phone acontact.Phone `json:"phone,omitempty"`

	// Logins holds the history of login sessions.
	Logins asessions.LoginSessionDeviceDates `json:"logins,omitempty"`

	// LDS holds the legal document signatures.
	LDS alegal.LegalDocSignatures `json:"lds"`

	// Roles holds the roles assigned to the user.
	Roles asessions.Roles `json:"roles"`

	// AdminLock holds the admin lock status of the account.
	AdminLock AdminLock `json:"adminLock,omitempty"`
}

// AddDeviceLogin adds a new device login to the user's login history.
// It trims the device name, sets the IP address, and records the current time.
// If the login history exceeds maxHistory, the oldest entry is removed.
func (ua *UserAccount) AddDeviceLogin(device string, realIP string, maxHistory int) *asessions.LoginSessionDeviceDate {
	if maxHistory < 1 {
		maxHistory = 10
	}

	lsdd := &asessions.LoginSessionDeviceDate{
		Device: strings.TrimSpace(device),
		IP:     realIP,
		Date:   time.Now().UTC(),
	}

	// Insert at the beginning
	ua.Logins = append(asessions.LoginSessionDeviceDates{lsdd}, ua.Logins...)

	// Trim the history to maxHistory
	if len(ua.Logins) > maxHistory {
		ua.Logins = ua.Logins[:maxHistory]
	}

	return lsdd
}

// AddLDSByKey adds a legal document signature by key.
// If appendIfFound is false and the key is found, no changes are made.
// If the key is not found or appendIfFound is true, a new signature is added.
func (ua *UserAccount) AddLDSByKey(key alegal.LegalDocSignatureKey, appendIfFound bool, effectiveDate *time.Time) (hasChanges bool) {
	if ua.LDS == nil {
		ua.LDS = alegal.LegalDocSignatures{}
	}

	ld := ua.LDS.Find(key)
	if ld != nil {
		if !appendIfFound {
			return false
		}
	}

	ua.LDS = append(ua.LDS, &alegal.LegalDocSignature{
		Key:           key,
		AcceptDate:    atime.GetNowUTCPointer(),
		EffectiveDate: effectiveDate,
	})

	return true
}
