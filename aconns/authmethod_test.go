package aconns

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthMethod_IsValid(t *testing.T) {
	tests := []struct {
		method   AuthMethod
		expected bool
	}{
		{AUTHMETHOD_PRIMARY, true},
		{AUTHMETHOD_MFA, true},
		{AUTHMETHOD_SSPR, true},
		{"invalid", false},
		{"  primary  ", true}, // TrimSpace should be used in logic
	}

	for _, tt := range tests {
		valid := tt.method.IsValid()
		if tt.expected != valid {
			t.Errorf("AuthMethod(%q).IsValid() = %v; want %v", tt.method, valid, tt.expected)
		}
	}
}

func TestAuthMethod_TrimSpace(t *testing.T) {
	raw := AuthMethod("  mfa  ")
	trimmed := raw.TrimSpace()
	assert.Equal(t, AuthMethod("mfa"), trimmed)
}

func TestAuthMethod_StringAndRaw(t *testing.T) {
	raw := AuthMethod("  sspr  ")
	assert.Equal(t, "sspr", raw.String())
	assert.Equal(t, "  sspr  ", raw.Raw())
}

func TestAuthMethods_Has(t *testing.T) {
	methods := AuthMethods{AUTHMETHOD_PRIMARY, AUTHMETHOD_SSPR}
	assert.True(t, methods.Has(AUTHMETHOD_SSPR))
	assert.False(t, methods.Has(AUTHMETHOD_MFA))
}

func TestAuthMethods_Validate(t *testing.T) {
	valid := AuthMethods{AUTHMETHOD_PRIMARY, AUTHMETHOD_MFA}
	assert.NoError(t, valid.Validate())

	invalid := AuthMethods{AUTHMETHOD_PRIMARY, "foobar"}
	err := invalid.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid")
}
