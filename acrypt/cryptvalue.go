package acrypt

import (
	"encoding/base64"
	"fmt"
	"strings"
	"sync"
	"time"
)

// CryptValue represents a simple rotatable secret value (no encryption, just decoding and rotation).
type CryptValue struct {
	Value             string       `json:"value"` // Formatted: "base64;<base64-encoded-key>"
	valueDecoded      []byte       // Cached decoded bytes.
	OldValue          string       `json:"oldValue,omitempty"`          // Previous value during rotation.
	OldValueExpiresAt *time.Time   `json:"oldValueExpiresAt,omitempty"` // Expiration for old value.
	MaxDuration       int          `json:"maxDuration,omitempty"`       // Max validity in minutes.
	mu                sync.RWMutex // Concurrency protection.
}

// GetDecoded returns the cached decoded value.
func (cv *CryptValue) GetDecoded() []byte {
	cv.mu.RLock()
	defer cv.mu.RUnlock()
	return cv.valueDecoded
}

// HasValue checks if the value is set.
func (cv *CryptValue) HasValue() bool {
	cv.mu.RLock()
	defer cv.mu.RUnlock()
	return strings.TrimSpace(cv.Value) != ""
}

// IsValid checks if the value is parseable (base64 format).
func (cv *CryptValue) IsValid() bool {
	cv.mu.RLock()
	defer cv.mu.RUnlock()
	parts := strings.SplitN(cv.Value, ";", 2)
	if len(parts) != 2 || parts[0] != "base64" {
		return false
	}
	_, err := base64.StdEncoding.DecodeString(parts[1])
	return err == nil
}

// Decode decodes the base64 value (caches if not already).
func (cv *CryptValue) Decode() ([]byte, error) {
	cv.mu.Lock()
	defer cv.mu.Unlock()
	return cv.decode()
}

// decode decodes the base64 value (caches if not already).
func (cv *CryptValue) decode() ([]byte, error) {
	if len(cv.valueDecoded) > 0 {
		return cv.valueDecoded, nil
	}
	parts := strings.SplitN(cv.Value, ";", 2)
	if len(parts) != 2 || parts[0] != "base64" {
		return nil, fmt.Errorf("invalid format")
	}
	decoded, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("failed to decode: %v", err)
	}
	cv.valueDecoded = decoded
	return decoded, nil
}

// Rotate sets a new value, moves current to old, and sets expiration.
func (cv *CryptValue) Rotate(newValue string, duration time.Duration) error {
	cv.mu.Lock()
	defer cv.mu.Unlock()
	cv.OldValue = cv.Value
	if duration > 0 {
		expiresAt := time.Now().Add(duration)
		cv.OldValueExpiresAt = &expiresAt
	}
	cv.Value = newValue
	// Decode new value
	_, err := cv.decode()
	return err
}

// HasExpired checks if current or old value is expired.
func (cv *CryptValue) HasExpired() bool {
	cv.mu.RLock()
	defer cv.mu.RUnlock()
	// Note: No ExpiresAt on current; use MaxDuration or add if needed.
	if cv.MaxDuration > 0 && time.Since(time.Time{}) > time.Duration(cv.MaxDuration)*time.Minute { // Placeholder; adjust if needed
		return true
	}
	if cv.OldValueExpiresAt != nil && time.Now().After(*cv.OldValueExpiresAt) {
		return true
	}
	return false
}

// CryptValueMap is a map of secrets with helper methods.
type CryptValueMap map[SecretsKey]*CryptValue

// Initialize generates new secrets for missing keys.
func (cvm CryptValueMap) Initialize(requiredKeys []SecretsKey) error {
	for _, key := range requiredKeys {
		if _, exists := cvm[key]; !exists {
			cvm[key] = &CryptValue{}
		}
		cv := cvm[key]
		if !cv.HasValue() {
			genKey, err := GenerateSecretKey()
			if err != nil {
				return fmt.Errorf("failed to generate %s: %v", key, err)
			}
			encoded := base64.StdEncoding.EncodeToString(genKey)
			cv.Value = fmt.Sprintf("base64;%s", encoded)
			_, err = cv.Decode() // Cache
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// SetCryptValueClearBytes sets or updates the CryptValue for the specified key using the provided clear bytes.
// The bytes are encoded in base64 and stored in the format "base64;<encoded>".
// It validates the key and bytes, decodes to cache the value, and returns an error if the key is empty,
// the bytes are nil or empty, or if decoding fails (though decoding failure is unlikely since the value is freshly encoded).
func (cvm CryptValueMap) SetCryptValueClearBytes(key SecretsKey, clearBytes []byte) error {
	if key.IsEmpty() {
		return fmt.Errorf("key is empty")
	}
	if clearBytes == nil || len(clearBytes) == 0 {
		return fmt.Errorf("clear bytes is empty")
	}
	encoded := base64.StdEncoding.EncodeToString(clearBytes)
	cv := &CryptValue{
		Value: fmt.Sprintf("base64;%s", encoded),
	}
	_, err := cv.Decode() // Cache
	if err != nil {
		return err
	}
	cvm[key] = cv
	return nil
}

// Validate checks all required keys are present and valid.
func (cvm CryptValueMap) Validate(requiredKeys []SecretsKey) error {
	for _, key := range requiredKeys {
		cv, exists := cvm[key]
		if !exists {
			return fmt.Errorf("%s missing", key)
		}
		if !cv.IsValid() {
			return fmt.Errorf("%s invalid", key)
		}
	}
	return nil
}

// Rotate rotates secrets for required keys with grace duration.
func (cvm CryptValueMap) Rotate(requiredKeys []SecretsKey, graceDuration time.Duration) error {
	for _, key := range requiredKeys {
		cv, exists := cvm[key]
		if !exists {
			return fmt.Errorf("%s missing", key)
		}
		genKey, err := GenerateSecretKey()
		if err != nil {
			return fmt.Errorf("failed to rotate %s: %v", key, err)
		}
		encoded := base64.StdEncoding.EncodeToString(genKey)
		newValue := fmt.Sprintf("base64;%s", encoded)
		if err := cv.Rotate(newValue, graceDuration); err != nil {
			return err
		}
	}
	return nil
}

// HasAnyExpired checks if any required secret expired.
func (cvm CryptValueMap) HasAnyExpired(requiredKeys []SecretsKey) bool {
	for _, key := range requiredKeys {
		if cv := cvm[key]; cv != nil && cv.HasExpired() {
			return true
		}
	}
	return false
}

// GetDecoded returns decoded value for a key.
func (cvm CryptValueMap) GetDecoded(key SecretsKey) ([]byte, error) {
	cv, exists := cvm[key]
	if !exists {
		return nil, fmt.Errorf("%s missing", key)
	}
	return cv.Decode()
}

// Set sets or updates a CryptValue.
func (cvm CryptValueMap) Set(key SecretsKey, value string) error {
	if key.IsEmpty() {
		return fmt.Errorf("key empty")
	}
	cv := &CryptValue{Value: value}
	_, err := cv.Decode() // Validate and cache
	if err != nil {
		return err
	}
	cvm[key] = cv
	return nil
}

// Delete removes a key.
func (cvm CryptValueMap) Delete(key SecretsKey) {
	delete(cvm, key)
}
