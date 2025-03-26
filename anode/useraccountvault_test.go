package anode

import (
	"github.com/jpfluger/alibs-slim/acrypt"
	"testing"
)

func TestUserVault_GetCredential(t *testing.T) {
	credential := UserCredential{Username: "testuser"}
	vault := UserVault{Credential: credential}

	if got := vault.GetCredential(); got.Username != "testuser" {
		t.Errorf("expected Username to be 'testuser', got %v", got.Username)
	}
}

func TestUserVault_GetSupportPin(t *testing.T) {
	pin := acrypt.IdBrief("1234")
	vault := UserVault{}
	vault.Support.Pin = pin

	if got := vault.GetSupportPin(); got != pin {
		t.Errorf("expected Support.Pin to be '1234', got %v", got)
	}
}

func TestUserVault_HasTOTP(t *testing.T) {
	totp := acrypt.TOTP{Secret: "secret"}
	vault := UserVault{TOTP: totp}

	if !vault.HasTOTP() {
		t.Errorf("expected HasTOTP to be true, got false")
	}

	vault.TOTP = acrypt.TOTP{}
	if vault.HasTOTP() {
		t.Errorf("expected HasTOTP to be false, got true")
	}
}

func TestUserVault_HasTOTPNew(t *testing.T) {
	totpNew := &acrypt.TOTP{Secret: "newsecret"}
	vault := UserVault{TOTPNew: totpNew}

	if !vault.HasTOTPNew() {
		t.Errorf("expected HasTOTPNew to be true, got false")
	}

	vault.TOTPNew = nil
	if vault.HasTOTPNew() {
		t.Errorf("expected HasTOTPNew to be false, got true")
	}
}

func TestUserVault_GenerateTokenBackups(t *testing.T) {
	vault := UserVault{}

	if err := vault.GenerateTokenBackups(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(vault.TokenBackups) != 16 {
		t.Errorf("expected 16 token backups, got %d", len(vault.TokenBackups))
	}
}

func TestUserVault_GenerateTokenBackupsWithOptions(t *testing.T) {
	vault := UserVault{}

	if err := vault.GenerateTokenBackupsWithOptions(10, 6); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(vault.TokenBackups) != 10 {
		t.Errorf("expected 10 token backups, got %d", len(vault.TokenBackups))
	}
}
