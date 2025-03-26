package acrypt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestSecretsKeyIsEmpty tests the IsEmpty method for SecretsKey.
func TestSecretsKeyIsEmpty(t *testing.T) {
	var sk SecretsKey

	// Test with an empty key
	sk = ""
	assert.True(t, sk.IsEmpty(), "Key should be empty")

	// Test with a key that contains only whitespace
	sk = "   "
	assert.True(t, sk.IsEmpty(), "Key should be empty")

	// Test with a non-empty key
	sk = SECRETSKEY_APPHOSTEK
	assert.False(t, sk.IsEmpty(), "Key should not be empty")
}

// TestSecretsKeyTrimSpace tests the TrimSpace method for SecretsKey.
func TestSecretsKeyTrimSpace(t *testing.T) {
	sk := SecretsKey("  apphostek  ")
	expected := SECRETSKEY_APPHOSTEK
	assert.Equal(t, expected, sk.TrimSpace(), "Key should be trimmed")
}

// TestSecretsKeyString tests the String method for SecretsKey.
func TestSecretsKeyString(t *testing.T) {
	sk := SECRETSKEY_APPHOSTEK
	assert.Equal(t, "apphostek", sk.String(), "String method should return the correct string")
}

// TestSecretsKeyToStringTrimLower tests the ToStringTrimLower method for SecretsKey.
func TestSecretsKeyToStringTrimLower(t *testing.T) {
	sk := SecretsKey("  APPHOSTEK  ")
	expected := "apphostek"
	assert.Equal(t, expected, sk.ToStringTrimLower(), "Key should be trimmed and lowercase")
}
