package acrypt

import (
	"encoding/base64"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSecretsValue_GetDecoded(t *testing.T) {
	s := SecretsValue{}
	s.valueDecoded = []byte("decoded_value")

	decoded := s.GetDecoded()

	assert.Equal(t, []byte("decoded_value"), decoded, "GetDecoded should return the decoded value")
}

func TestSecretsValue_NewJWTSecretKey(t *testing.T) {
	s := SecretsValue{
		MaxDuration: 60, // 1 hour max duration
	}

	err := s.NewJWTSecretKey()

	assert.NoError(t, err, "NewJWTSecretKey should not return an error")
	assert.NotEmpty(t, s.Value, "NewJWTSecretKey should set a new value")
	assert.NotNil(t, s.valueDecoded, "NewJWTSecretKey should decode the key")
	assert.NotNil(t, s.OldValueExpiresAt, "NewJWTSecretKey should set OldValueExpiresAt")
	assert.WithinDuration(t, time.Now().Add(1*time.Hour), *s.OldValueExpiresAt, time.Minute)
}

func TestSecretsValue_Decode(t *testing.T) {
	plainText := "Hello world"
	password := "testpassword"

	// Generate an AES256 encrypted value for testing
	cipherText, err := AESGCM256Encrypt([]byte(plainText), password)
	if err != nil {
		t.Fatalf("Failed to encrypt data: %v", err)
	}

	// Encode the cipherText to base64
	encodedCipherText := base64.StdEncoding.EncodeToString(cipherText)

	// Create a SecretsValue with the encoded encrypted value
	s := SecretsValue{
		Value: SecretsValueRaw("e;plain;aes256;" + encodedCipherText),
	}

	// Decode the value
	decoded, err := s.Decode(password, true)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	// Assert the decoded value matches the original plaintext
	assert.Equal(t, []byte(plainText), decoded, "Decode should return the correct decoded value")
	assert.Equal(t, []byte(plainText), s.GetDecoded(), "Decode should cache the decoded value when cacheDecoded is true")
}

func TestSecretsValue_Rotate(t *testing.T) {
	s := SecretsValue{
		Value: SecretsValueRaw("d;plain;aes256;old_value"),
	}

	newValue := "new_value"
	rotationDuration := 30 * time.Minute

	s.Rotate(newValue, rotationDuration)

	assert.Equal(t, SecretsValueRaw("d;plain;aes256;old_value"), s.OldValue, "Rotate should set the old value")
	assert.Equal(t, SecretsValueRaw("d;plain;aes256;new_value"), s.Value, "Rotate should update the new value")
	assert.NotNil(t, s.OldValueExpiresAt, "Rotate should set OldValueExpiresAt when a duration is provided")
	assert.WithinDuration(t, time.Now().Add(rotationDuration), *s.OldValueExpiresAt, time.Minute)
}

func TestSecretsValue_HasExpired(t *testing.T) {
	now := time.Now()
	expired := now.Add(-1 * time.Hour)
	notExpired := now.Add(1 * time.Hour)

	tests := []struct {
		name       string
		expires    *time.Time
		oldExpires *time.Time
		expected   bool
	}{
		{
			name:       "Neither expired",
			expires:    &notExpired,
			oldExpires: &notExpired,
			expected:   false,
		},
		{
			name:       "Current expired",
			expires:    &expired,
			oldExpires: &notExpired,
			expected:   true,
		},
		{
			name:       "Old expired",
			expires:    &notExpired,
			oldExpires: &expired,
			expected:   true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := SecretsValue{
				ExpiresAt:         test.expires,
				OldValueExpiresAt: test.oldExpires,
			}

			assert.Equal(t, test.expected, s.HasExpired(), "HasExpired should return the correct result")
		})
	}
}

func TestSecretsValue_IsDecoded(t *testing.T) {
	s := SecretsValue{}
	assert.False(t, s.IsDecoded(), "IsDecoded should return false when valueDecoded is nil")

	s.valueDecoded = []byte("decoded_value")
	assert.True(t, s.IsDecoded(), "IsDecoded should return true when valueDecoded is set")
}

func TestEnsureCryptMode(t *testing.T) {
	password := "strong_password"
	plainValue := "plain_secret_data"
	s := &SecretsValue{}

	// Initialize in decrypted mode.
	s.Value.Validate(plainValue)

	// Test switching to encrypted mode.
	t.Run("SwitchToEncrypted", func(t *testing.T) {
		err := s.EnsureCryptMode(password, CRYPTMODE_ENCRYPTED)
		if err != nil {
			t.Fatalf("failed to encrypt: %v", err)
		}

		// Verify the new mode is encrypted.
		mode, _, _, value, err := s.Value.Parse()
		if err != nil {
			t.Fatalf("failed to parse encrypted value: %v", err)
		}
		if mode != CRYPTMODE_ENCRYPTED {
			t.Fatalf("expected mode %s, got %s", CRYPTMODE_ENCRYPTED, mode)
		}

		// Verify the encrypted value is base64 encoded.
		if _, err := base64.StdEncoding.DecodeString(value); err != nil {
			t.Fatalf("encrypted value is not valid base64: %v", err)
		}
	})

	// Test switching back to decrypted mode.
	t.Run("SwitchToDecrypted", func(t *testing.T) {
		err := s.EnsureCryptMode(password, CRYPTMODE_DECRYPTED)
		if err != nil {
			t.Fatalf("failed to decrypt: %v", err)
		}

		// Verify the new mode is decrypted.
		mode, _, _, value, err := s.Value.Parse()
		if err != nil {
			t.Fatalf("failed to parse decrypted value: %v", err)
		}
		if mode != CRYPTMODE_DECRYPTED {
			t.Fatalf("expected mode %s, got %s", CRYPTMODE_DECRYPTED, mode)
		}

		// Verify the decrypted value matches the original plain text.
		if value != plainValue {
			t.Fatalf("expected value %s, got %s", plainValue, value)
		}
	})

	// Test error case: invalid target mode.
	t.Run("InvalidTargetMode", func(t *testing.T) {
		err := s.EnsureCryptMode(password, "invalid_mode")
		if err == nil {
			t.Fatal("expected error for invalid target mode, got nil")
		}
	})

	// Test error case: invalid encryption type (modify the SecretsValue manually).
	t.Run("InvalidEncryptionType", func(t *testing.T) {
		invalidValue := SecretsValueRaw("e;base64;invalid_encryption;some_value")
		s.Value = invalidValue
		err := s.EnsureCryptMode(password, CRYPTMODE_DECRYPTED)
		if err == nil {
			t.Fatal("expected error for invalid encryption type, got nil")
		}
	})

	// Test error case: invalid encoding type (modify the SecretsValue manually).
	t.Run("InvalidEncodingType", func(t *testing.T) {
		invalidValue := SecretsValueRaw("e;invalid_encoding;aes256;some_value")
		s.Value = invalidValue
		err := s.EnsureCryptMode(password, CRYPTMODE_DECRYPTED)
		if err == nil {
			t.Fatal("expected error for invalid encoding type, got nil")
		}
	})
}
