package acrypt

import (
	"encoding/base64"
	"fmt"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/argon2"
	"strings"
	"testing"
)

func TestEncryptBCrypt(t *testing.T) {
	// Test with a non-empty string
	hashed, err := EncryptBCrypt("password123")
	assert.NoError(t, err)
	assert.True(t, strings.HasPrefix(hashed, "$2a$") || strings.HasPrefix(hashed, "$2b$"))

	// Test with an empty string
	_, err = EncryptBCrypt("")
	assert.Error(t, err)
}

func TestIsEqualBCrypt(t *testing.T) {
	hashed, _ := EncryptBCrypt("password123")
	assert.True(t, IsEqualBCrypt(hashed, "password123"))
	assert.False(t, IsEqualBCrypt(hashed, "wrongpassword"))
}

func TestEncryptSCrypt(t *testing.T) {
	presets := NewScryptPresets()
	hashed, err := EncryptSCrypt("password123", presets)
	assert.NoError(t, err)
	assert.True(t, strings.HasPrefix(hashed, "$scrypt$"))
}

func TestIsEqualSCrypt(t *testing.T) {
	presets := NewScryptPresets()
	hashed, _ := EncryptSCrypt("password123", presets)
	match, err := IsEqualSCrypt(hashed, "password123")
	assert.NoError(t, err)
	assert.True(t, match)

	match, err = IsEqualSCrypt(hashed, "wrongpassword")
	assert.NoError(t, err)
	assert.False(t, match)
}

func TestEncryptArgon2id(t *testing.T) {
	presets := NewArgon2Presets()
	hashed, err := EncryptArgon2id("password123", presets)
	assert.NoError(t, err)
	assert.True(t, strings.HasPrefix(hashed, "$argon2id$"))
}

func TestIsEqualArgon2id(t *testing.T) {
	presets := NewArgon2Presets()
	hashed, _ := EncryptArgon2id("password123", presets)
	match, err := IsEqualArgon2id(hashed, "password123")
	assert.NoError(t, err)
	assert.True(t, match)

	match, err = IsEqualArgon2id(hashed, "wrongpassword")
	assert.NoError(t, err)
	assert.False(t, match)
}

// TestIsArgon2idHash checks the IsArgon2idHash function.
func TestIsArgon2idHash(t *testing.T) {
	// Test cases
	tests := []struct {
		hash     string
		expected bool
	}{
		{"$argon2id$v=19$m=65536,t=1,p=4$someSalt$someHash", true},
		{"$argon2i$v=19$m=65536,t=1,p=4$someSalt$someHash", false},        // Incorrect identifier
		{"$argon2id$v=19$m=65536,t=1,p=4$someSalt", false},                // Missing hash part
		{"$argon2id$v=19$m=65536,t=1,p=4$", false},                        // Missing salt and hash
		{"$argon2id$v=19$m=65536,t=1,p=4$someSalt$someHash$extra", false}, // Extra part
		{"$bcrypt$someSalt$someHash", false},                              // Different algorithm
		{"", false},                                                       // Empty string
	}

	for _, test := range tests {
		result := IsArgon2idHash(test.hash)
		assert.Equal(t, test.expected, result, "Hash: %s", test.hash)
	}
}

func TestArgon2HashEncoding_Raw_vs_Std(t *testing.T) {
	// Test input, which is RawStdEncoding.
	hash := "$argon2id$v=19$m=65536,t=1,p=4$W0slJwWEDTWj14RWKx73QQ$VjG08hh5d4Lj9CQyrd7vaHeOYGXZm1TEGXyYYsIGl9g"
	password := "password123"

	// Split the hash into components
	parts := strings.Split(hash, "$")
	assert.Equal(t, 6, len(parts), "Invalid hash format")

	// Extract components
	saltEncodedRaw := parts[4]
	hashEncodedRaw := parts[5]

	// Decode salt and hash using RawStdEncoding
	saltRaw, err := base64.RawStdEncoding.DecodeString(saltEncodedRaw)
	assert.NoError(t, err, "Error decoding salt with RawStdEncoding")

	hashRaw, err := base64.RawStdEncoding.DecodeString(hashEncodedRaw)
	assert.NoError(t, err, "Error decoding hash with RawStdEncoding")

	// Convert RawStdEncoding to StdEncoding (add padding)
	saltStd := base64.StdEncoding.EncodeToString(saltRaw)
	hashStd := base64.StdEncoding.EncodeToString(hashRaw)

	// Print the converted hash
	convertedHash := fmt.Sprintf("$argon2id$v=19$m=65536,t=1,p=4$%s$%s", saltStd, hashStd)
	fmt.Println("Converted Hash with StdEncoding:", convertedHash)

	// Validate that the StdEncoding version decodes correctly
	decodedSalt, err := base64.StdEncoding.DecodeString(saltStd)
	assert.NoError(t, err, "Error decoding converted salt with StdEncoding")
	assert.Equal(t, saltRaw, decodedSalt, "Decoded salt does not match original")

	decodedHash, err := base64.StdEncoding.DecodeString(hashStd)
	assert.NoError(t, err, "Error decoding converted hash with StdEncoding")
	assert.Equal(t, hashRaw, decodedHash, "Decoded hash does not match original")

	// Compute Argon2 hash using the decoded salt
	params := &Argon2Presets{
		Memory:  65536,
		Time:    1,
		Threads: 4,
		KeyLen:  uint32(len(hashRaw)),
	}

	computedHash := argon2.IDKey([]byte(password), decodedSalt, params.Time, params.Memory, params.Threads, params.KeyLen)
	computedHashEncoded := base64.StdEncoding.EncodeToString(computedHash)

	// Compare the computed hash with the StdEncoded version
	assert.Equal(t, hashStd, computedHashEncoded, "Converted hash does not match computed hash")
}
