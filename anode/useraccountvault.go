package anode

import (
	"time"

	"github.com/jpfluger/alibs-slim/acrypt"
	"github.com/jpfluger/alibs-slim/atime"
)

// UserVault represents a user's vault containing credentials, TOTP, and backup tokens.
type UserVault struct {
	// Credential holds the user's credentials.
	Credential UserCredential `json:"credential,omitempty"`

	// TOTP holds the current TOTP (Time-based One-Time Password) configuration.
	TOTP acrypt.TOTP `json:"totp,omitempty"`

	// TOTPNew holds the new TOTP configuration, if any.
	TOTPNew *acrypt.TOTP `json:"totpNew,omitempty"`

	// TokenBackups holds the backup tokens for the user.
	TokenBackupsDate *time.Time             `json:"tokenBackupsDate,omitempty"`
	TokenBackups     acrypt.MiniRandomCodes `json:"tokenBackups,omitempty"`

	// Support contains the support pin used to verify a user instead of a social security number.
	Support struct {
		Pin acrypt.IdBrief `json:"pin,omitempty"`
	} `json:"support,omitempty"`
}

// GetCredential returns a pointer to the user's credential.
func (uv *UserVault) GetCredential() *UserCredential {
	return &uv.Credential
}

// GetSupportPin returns the support pin.
func (uv *UserVault) GetSupportPin() acrypt.IdBrief {
	return uv.Support.Pin
}

// HasTOTP checks if the current TOTP configuration has a secret.
func (uv *UserVault) HasTOTP() bool {
	return uv.TOTP.HasSecret()
}

// HasTOTPNew checks if the new TOTP configuration exists and has a secret.
func (uv *UserVault) HasTOTPNew() bool {
	return uv.TOTPNew != nil && uv.TOTPNew.HasSecret()
}

// HasTokenBackups checks if the vault has tokens that can be used for backup and recover.
func (uv *UserVault) HasTokenBackups() bool {
	if uv.TokenBackupsDate == nil || uv.TokenBackupsDate.IsZero() || uv.TokenBackups == nil || len(uv.TokenBackups) == 0 {
		return false
	}
	return true
}

// GenerateTokenBackups generates backup tokens with default options.
func (uv *UserVault) GenerateTokenBackups() error {
	return uv.GenerateTokenBackupsWithOptions(16, 8)
}

// GenerateTokenBackupsWithOptions generates backup tokens with specified options.
func (uv *UserVault) GenerateTokenBackupsWithOptions(maxCount int, length int) error {
	if maxCount < 1 {
		maxCount = 16
	}
	if length < 1 {
		length = 8
	}
	codes := acrypt.MiniRandomCodes{}
	if err := codes.GenerateWithCharSet(maxCount, length, "", "ABCDEFGHJKMNPQRSTUVWXZ0123456789"); err != nil {
		return err
	}
	uv.TokenBackups = codes
	uv.TokenBackupsDate = atime.GetNowUTCPointer()
	return nil
}
