package acrypt

import (
	"encoding/base64"
	"errors"
	"fmt"
	"sync"
	"time"
)

// SecretsValue represents a single item in the secrets management system.
type SecretsValue struct {
	Value             SecretsValueRaw `json:"value"` // The current raw value of the secret.
	valueDecoded      []byte          // Decoded value of the secret (private field).
	ExpiresAt         *time.Time      `json:"expiresAt,omitempty"`         // Expiration time for the current value.
	OldValue          SecretsValueRaw `json:"oldValue,omitempty"`          // The previous raw value of the secret being rotated out.
	OldValueExpiresAt *time.Time      `json:"oldValueExpiresAt,omitempty"` // Expiration time for the old value.
	MaxDuration       int             `json:"maxDuration,omitempty"`       // Maximum duration (in minutes) that the secret is valid.
	mu                sync.RWMutex    // Mutex to protect concurrent access to the SecretsValue.
}

// GetDecoded returns valueDecoded.
func (s *SecretsValue) GetDecoded() []byte {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.valueDecoded
}

// HasValue returns true if there is content in the Value field.
func (s *SecretsValue) HasValue() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return !s.Value.IsEmpty()
}

// IsValidParsedValue returns true if there is content in the Value field that is not empty and parseable.
func (s *SecretsValue) IsValidParsedValue() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.Value.IsEmpty() {
		return false
	}
	_, _, _, _, err := s.Value.Parse()
	return err == nil
}

// NewJWTSecretKey initializes a new JWT secret key with a default encrypted format and rotates old values.
func (s *SecretsValue) NewJWTSecretKey() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Generate a new ChaCha20 key.
	secretKey, err := GenerateChaCha20Key()
	if err != nil {
		return fmt.Errorf("failed to generate new JWT secret key: %v", err)
	}

	// Encode the key in base64.
	encodedKey := base64.StdEncoding.EncodeToString(secretKey)

	// Update old value and assign new value with encrypted format.
	s.OldValue = s.Value
	if s.MaxDuration > 0 {
		expiresAt := time.Now().Add(time.Duration(s.MaxDuration) * time.Minute)
		s.OldValueExpiresAt = &expiresAt
	}
	newValue := fmt.Sprintf("%s;%s;%s;%s", CRYPTMODE_DECRYPTED, ENCODINGTYPE_BASE64, ENCRYPTIONTYPE_AES256, encodedKey)
	s.Value = SecretsValueRaw(newValue)
	s.valueDecoded = secretKey
	return nil
}

// Decode decrypts and decodes the value of the secret, with an option to cache the decoded value.
func (s *SecretsValue) Decode(password string, cacheDecoded bool) ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	decoded, err := s.Value.Decode(password)
	if err != nil {
		return nil, err
	}

	if cacheDecoded {
		s.valueDecoded = decoded
	}
	return decoded, nil
}

// Rotate updates the secret value and optionally extends expiration.
func (s *SecretsValue) Rotate(newValue string, duration time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.OldValue = s.Value
	if duration > 0 {
		expiresAt := time.Now().Add(duration)
		s.OldValueExpiresAt = &expiresAt
	}
	s.Value.Validate(newValue)
	s.valueDecoded, _ = s.Value.Decode("")
}

// HasExpired checks whether the secret has expired.
func (s *SecretsValue) HasExpired() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.ExpiresAt != nil && time.Now().After(*s.ExpiresAt) {
		return true
	}
	if s.OldValueExpiresAt != nil && time.Now().After(*s.OldValueExpiresAt) {
		return true
	}
	return false
}

// IsDecoded checks if the valueDecoded is set.
func (s *SecretsValue) IsDecoded() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.valueDecoded) > 0
}

// EnsureCryptMode ensures the CryptMode is set to the specified mode (either "e" or "d").
// It switches between encryption and decryption based on the current mode and the provided password.
func (s *SecretsValue) EnsureCryptMode(password string, targetMode CryptMode) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Parse the current mode, encoding, encryption type, and value.
	// The rawValue is already encoded.
	currentMode, encoding, encryption, rawValue, err := s.Value.Parse()
	if err != nil {
		return fmt.Errorf("failed to parse SecretsValue: %v", err)
	}

	// If the current mode matches the target mode, no changes are needed.
	if currentMode == targetMode {
		return nil
	}

	switch targetMode {
	case CRYPTMODE_ENCRYPTED:
		// Ensure the current mode is decrypted, so we can encrypt the value.
		if currentMode != CRYPTMODE_DECRYPTED {
			return errors.New("cannot encrypt: current mode is not decrypted")
		}

		// Encrypt the raw value using the specified encryption type.
		var encryptedValue []byte
		switch encryption {
		case ENCRYPTIONTYPE_AES128:
			encryptedValue, err = AESGCM128Encrypt([]byte(rawValue), password)
		case ENCRYPTIONTYPE_AES256:
			encryptedValue, err = AESGCM256Encrypt([]byte(rawValue), password)
		default:
			return errors.New("unsupported encryption type")
		}
		if err != nil {
			return fmt.Errorf("failed to encrypt value: %v", err)
		}

		// Always encode with the encrypted value in the specified encoding type.
		encodedValue := base64.StdEncoding.EncodeToString(encryptedValue)

		// Update the SecretsValue with the encrypted format.
		s.Value = SecretsValueRaw(fmt.Sprintf("%s;%s;%s;%s", CRYPTMODE_ENCRYPTED, encoding, encryption, encodedValue))
		s.valueDecoded = nil // Clear the decoded value since it has changed.

	case CRYPTMODE_DECRYPTED:
		// Ensure the current mode is encrypted, so we can decrypt the value.
		if currentMode != CRYPTMODE_ENCRYPTED {
			return errors.New("cannot decrypt: current mode is not encrypted")
		}

		// Decode and decrypt the value.
		decodedValue, err := s.Value.Decode(password)
		if err != nil {
			return fmt.Errorf("failed to decode and decrypt value: %v", err)
		}

		// Update the SecretsValue with the decrypted format.
		s.Value = SecretsValueRaw(fmt.Sprintf("%s;%s;%s;%s", CRYPTMODE_DECRYPTED, encoding, encryption, string(decodedValue)))
		s.valueDecoded = decodedValue // Cache the decoded value.

	default:
		return errors.New("invalid target CryptMode")
	}

	return nil
}
