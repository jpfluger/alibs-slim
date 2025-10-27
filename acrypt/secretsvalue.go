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

// NewRandomSecret generates and sets a new random 32-byte secret key in decrypted format.
// Suitable for any secret type (e.g., JWT, BadgerDB encryption keys).
func (s *SecretsValue) NewRandomSecret() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Generate a new key.
	secretKey, err := GenerateSecretKey()
	if err != nil {
		return fmt.Errorf("failed to generate new secret key: %v", err)
	}

	// Encode the key in base64.
	encodedKey := base64.StdEncoding.EncodeToString(secretKey)

	// Update old value and expiration if MaxDuration is set.
	s.OldValue = s.Value
	if s.MaxDuration > 0 {
		expiresAt := time.Now().Add(time.Duration(s.MaxDuration) * time.Minute)
		s.OldValueExpiresAt = &expiresAt
	}

	// Assign new value with decrypted format.
	newValue := fmt.Sprintf("%s;%s;%s;%s", CRYPTMODE_DECRYPTED, ENCODINGTYPE_BASE64, ENCRYPTIONTYPE_AES256, encodedKey)
	s.Value = SecretsValueRaw(newValue)
	s.valueDecoded = secretKey
	return nil
}

// NewJWTSecretKey initializes a new JWT secret key with a default encrypted format and rotates old values.
func (s *SecretsValue) NewJWTSecretKey() error {
	return s.NewRandomSecret()
	//secretKey, err := GenerateSecretKey()
	//if err != nil {
	//	return fmt.Errorf("failed to generate new JWT secret key: %v", err)
	//}
	//
	//// Encode the key in base64.
	//encodedKey := base64.StdEncoding.EncodeToString(secretKey)
	//
	//// Update old value and assign new value with encrypted format.
	//s.OldValue = s.Value
	//if s.MaxDuration > 0 {
	//	expiresAt := time.Now().Add(time.Duration(s.MaxDuration) * time.Minute)
	//	s.OldValueExpiresAt = &expiresAt
	//}
	//newValue := fmt.Sprintf("%s;%s;%s;%s", CRYPTMODE_DECRYPTED, ENCODINGTYPE_BASE64, ENCRYPTIONTYPE_AES256, encodedKey)
	//s.Value = SecretsValueRaw(newValue)
	//s.valueDecoded = secretKey
	//return nil
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

	currentMode, encoding, encryption, rawValue, err := s.Value.Parse()
	if err != nil {
		return fmt.Errorf("failed to parse SecretsValue: %v", err)
	}

	if currentMode == targetMode {
		return nil
	}

	switch targetMode {
	case CRYPTMODE_ENCRYPTED:
		if currentMode != CRYPTMODE_DECRYPTED {
			return errors.New("cannot encrypt: current mode is not decrypted")
		}

		encryptedValue, err := AESGCMEncrypt([]byte(rawValue), password, encryption)
		if err != nil {
			return fmt.Errorf("failed to encrypt value: %w", err)
		}

		encodedValue := base64.StdEncoding.EncodeToString(encryptedValue)
		s.Value = SecretsValueRaw(fmt.Sprintf("%s;%s;%s;%s", CRYPTMODE_ENCRYPTED, encoding, encryption, encodedValue))
		s.valueDecoded = nil

	case CRYPTMODE_DECRYPTED:
		if currentMode != CRYPTMODE_ENCRYPTED {
			return errors.New("cannot decrypt: current mode is not encrypted")
		}

		decodedValue, err := s.Value.Decode(password) // Assuming Decode uses AESGCMDecrypt internally with encryption type
		if err != nil {
			return fmt.Errorf("failed to decode and decrypt value: %v", err)
		}

		s.Value = SecretsValueRaw(fmt.Sprintf("%s;%s;%s;%s", CRYPTMODE_DECRYPTED, encoding, encryption, string(decodedValue)))
		s.valueDecoded = decodedValue

	default:
		return errors.New("invalid target CryptMode")
	}

	return nil
}
