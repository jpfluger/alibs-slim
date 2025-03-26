package acrypt

import (
	"encoding/base64"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCryptKeyBase64(t *testing.T) {
	tests := []struct {
		name           string
		input          []byte
		expectedBase64 string
		expectedSize   int
		expectedError  string
	}{
		{
			name:           "Valid Ed25519 key",
			input:          make([]byte, Ed25519KeySize),
			expectedBase64: base64.StdEncoding.EncodeToString(make([]byte, Ed25519KeySize)),
			expectedSize:   Ed25519KeySize,
		},
		{
			name:           "Valid RSA key",
			input:          make([]byte, RSAKeySizeMin),
			expectedBase64: base64.StdEncoding.EncodeToString(make([]byte, RSAKeySizeMin)),
			expectedSize:   RSAKeySizeMin,
		},
		{
			name:           "Invalid key size",
			input:          make([]byte, 10),
			expectedBase64: base64.StdEncoding.EncodeToString(make([]byte, 10)),
			expectedSize:   Ed25519KeySize,
			expectedError:  "invalid key size",
		},
		{
			name:           "Invalid base64 key",
			input:          []byte("invalid-base64@@@"),
			expectedBase64: "",
			expectedSize:   Ed25519KeySize,
			expectedError:  "invalid key size: expected 32 bytes, got 17 bytes",
		},
		{
			name:           "Empty key",
			input:          []byte{},
			expectedBase64: "",
			expectedSize:   Ed25519KeySize,
			expectedError:  "key is empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := NewCryptKeyBase64(tt.input)

			// Check base64 encoding
			if tt.expectedBase64 != "" {
				assert.Equal(t, tt.expectedBase64, key.Encoded(), "Base64 encoding mismatch")
			}

			// Check IsEmpty
			if len(tt.input) == 0 {
				assert.True(t, key.IsEmpty(), "Key should be empty")
			} else {
				assert.False(t, key.IsEmpty(), "Key should not be empty")
			}

			// Validate key
			err := key.Validate(tt.expectedSize)
			if tt.expectedError != "" {
				assert.Error(t, err, "Validation should fail")
				assert.Contains(t, err.Error(), tt.expectedError, "Error message mismatch")
			} else {
				assert.NoError(t, err, "Validation should succeed")
			}

			// Decode key
			decoded, decodeErr := key.Decoded()
			if tt.expectedError == "key is not valid base64" {
				assert.Error(t, decodeErr, "Decoding should fail")
			} else if tt.expectedError == "" {
				assert.NoError(t, decodeErr, "Decoding should succeed")
				assert.Equal(t, tt.input, decoded, "Decoded key mismatch")
			}
		})
	}
}
