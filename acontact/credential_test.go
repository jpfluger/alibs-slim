package acontact

import (
	"github.com/jpfluger/alibs-slim/acrypt"
	"testing"
	"time"
)

func TestCredential_Validate(t *testing.T) {
	t.Run("Valid Credential", func(t *testing.T) {
		cred := &Credential{
			Type:       "some_type",
			Label:      "Test Credential",
			Secret:     acrypt.SecretsValue{Value: acrypt.SecretsValueRaw("e;base64;aes256;encoded_value")},
			Modified:   time.Now(),
			SecretFile: "",
		}

		if err := cred.Validate(); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("Empty Credential", func(t *testing.T) {
		var cred *Credential
		if err := cred.Validate(); err == nil {
			t.Errorf("expected error for nil credential, got nil")
		}
	})

	t.Run("Missing Type", func(t *testing.T) {
		cred := &Credential{
			Label:    "Missing Type",
			Modified: time.Now(),
		}

		err := cred.Validate()
		if err == nil || err.Error() != "credential type is empty" {
			t.Errorf("expected error for missing type, got %v", err)
		}
	})

	t.Run("Missing Secret", func(t *testing.T) {
		cred := &Credential{
			Type:       "some_type",
			Label:      "Missing Secret",
			Secret:     acrypt.SecretsValue{},
			Modified:   time.Now(),
			SecretFile: "",
		}

		err := cred.Validate()
		if err == nil || err.Error() != "credential secret is empty" {
			t.Errorf("expected error for missing secret, got %v", err)
		}
	})
}

func TestCredential_SetSecret(t *testing.T) {
	t.Run("Set Valid Secret", func(t *testing.T) {
		cred := &Credential{
			Type:  "some_type",
			Label: "Test Set Secret",
		}
		secret := []byte("new_secret")
		err := cred.SetSecret("master_password", secret)

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if !cred.Secret.HasValue() {
			t.Errorf("expected secret to be set, but it's empty")
		}
	})
}

func TestCredential_GetSecret(t *testing.T) {
	t.Run("Get Decoded Secret", func(t *testing.T) {
		cred := &Credential{
			Type:   "some_type",
			Label:  "Test Get Secret",
			Secret: acrypt.SecretsValue{Value: acrypt.SecretsValueRaw("d;plain;aes256;encoded_value")},
		}

		decodedSecret, err := cred.GetSecret("master_password")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if len(decodedSecret) == 0 {
			t.Errorf("expected decoded secret, got empty result")
		}
	})
}

func TestCredential_Clean(t *testing.T) {
	t.Run("Clean Valid Credentials", func(t *testing.T) {
		creds := Credentials{
			&Credential{
				Type:   "valid_type",
				Secret: acrypt.SecretsValue{Value: acrypt.SecretsValueRaw("e;base64;aes256;valid_value")},
			},
			&Credential{
				Type:   "invalid_type",
				Secret: acrypt.SecretsValue{},
			},
		}

		cleaned, err := creds.Clean()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if len(cleaned) != 1 {
			t.Errorf("expected 1 valid credential, got %d", len(cleaned))
		}
	})
}
