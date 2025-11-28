package acrypt

import (
	"fmt"
	"strings"
)

// HashValueMap is a simple map for storing hashed values (e.g., one-way password hashes) without decoding or rotation.
type HashValueMap map[SecretsKey]string

// Set sets or updates the hash value for the specified key.
// It validates that the key is not empty and the value is not blank.
func (hvm HashValueMap) Set(key SecretsKey, value string) error {
	if key.IsEmpty() {
		return fmt.Errorf("key is empty")
	}
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("value is empty")
	}
	hvm[key] = value
	return nil
}

// Get returns the hash value for a key, or an empty string if not found.
func (hvm HashValueMap) Get(key SecretsKey) string {
	return hvm[key]
}

// Has checks if the key exists and has a non-empty value.
func (hvm HashValueMap) Has(key SecretsKey) bool {
	value, exists := hvm[key]
	return exists && strings.TrimSpace(value) != ""
}

// Validate checks all required keys are present and non-empty.
func (hvm HashValueMap) Validate(requiredKeys []SecretsKey) error {
	for _, key := range requiredKeys {
		if !hvm.Has(key) {
			return fmt.Errorf("%s missing or empty", key)
		}
	}
	return nil
}

// Delete removes a key.
func (hvm HashValueMap) Delete(key SecretsKey) {
	delete(hvm, key)
}
