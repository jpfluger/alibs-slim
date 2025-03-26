package acrypt

import (
	"bytes"
	"encoding/base64"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSecretsValueRaw_Parse(t *testing.T) {
	tests := []struct {
		input    SecretsValueRaw
		expected struct {
			mode       CryptMode
			encoding   EncodingType
			encryption EncryptionType
			value      string
		}
		expectError bool
	}{
		{
			input: "e;base64;aes256;SGVsbG8gd29ybGQ=",
			expected: struct {
				mode       CryptMode
				encoding   EncodingType
				encryption EncryptionType
				value      string
			}{
				mode:       CRYPTMODE_ENCRYPTED,
				encoding:   ENCODINGTYPE_BASE64,
				encryption: ENCRYPTIONTYPE_AES256,
				value:      "SGVsbG8gd29ybGQ=",
			},
			expectError: false,
		},
		{
			input:       "x;base64;aes256;SGVsbG8gd29ybGQ=",
			expectError: true,
		},
		{
			input:       "e;base64;unknown;SGVsbG8gd29ybGQ=",
			expectError: true,
		},
		{
			input:       "e;base64;aes256",
			expectError: true,
		},
	}

	for _, test := range tests {
		mode, encoding, encryption, value, err := test.input.Parse()

		if test.expectError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expected.mode, mode)
			assert.Equal(t, test.expected.encoding, encoding)
			assert.Equal(t, test.expected.encryption, encryption)
			assert.Equal(t, test.expected.value, value)
		}
	}
}

func TestSecretsValueRaw_Validate(t *testing.T) {
	tests := []struct {
		input          string
		expectedOutput SecretsValueRaw
	}{
		{
			input:          "e;base64;aes256;SGVsbG8gd29ybGQ=",
			expectedOutput: "e;base64;aes256;SGVsbG8gd29ybGQ=",
		},
		{
			input:          "invalid_format",
			expectedOutput: "d;plain;aes256;invalid_format",
		},
	}

	for _, test := range tests {
		var raw SecretsValueRaw
		raw.Validate(test.input)
		assert.Equal(t, test.expectedOutput, raw)
	}
}

func TestSecretsValueRaw_Decode(t *testing.T) {
	password := "testpassword"
	// Generate an AES256 encrypted value for testing
	plainText := "Hello world"
	cipherText, _ := AESGCM256Encrypt([]byte(plainText), password)
	encodedCipherText := base64.StdEncoding.EncodeToString(cipherText)

	tests := []struct {
		input          SecretsValueRaw
		password       string
		expectedOutput []byte
		expectError    bool
	}{
		{
			input:          SecretsValueRaw("e;plain;aes256;" + encodedCipherText),
			password:       password,
			expectedOutput: []byte(plainText),
			expectError:    false,
		},
		{
			input:       SecretsValueRaw("e;base64;aes128;invalid_cipher_text"),
			password:    password,
			expectError: true,
		},
		{
			input:          SecretsValueRaw("d;base64;aes256;SGVsbG8gd29ybGQ="),
			password:       password,
			expectedOutput: []byte("Hello world"),
			expectError:    false,
		},
	}

	for _, test := range tests {
		decoded, err := test.input.Decode(test.password)

		if test.expectError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expectedOutput, decoded)
		}
	}
}

func TestSecretsValueRaw_Encode(t *testing.T) {
	tests := []struct {
		name           string
		rawValue       []byte
		masterPassword string
		initialValue   SecretsValueRaw
		expectedError  string
	}{
		{
			name:           "Valid encoding and encryption (AES128, Base64)",
			rawValue:       []byte("test_value"),
			masterPassword: "password123",
			initialValue:   SecretsValueRaw("e;base64;aes128;"),
			expectedError:  "",
		},
		{
			name:           "Valid plain encoding (Decrypted)",
			rawValue:       []byte("plain_value"),
			masterPassword: "password123",
			initialValue:   SecretsValueRaw("d;plain;aes256;"),
			expectedError:  "",
		},
		{
			name:           "Invalid encryption type",
			rawValue:       []byte("test_value"),
			masterPassword: "password123",
			initialValue:   SecretsValueRaw("e;base64;invalid_encryption;"),
			expectedError:  "", // e;base64;aes128;KH+9Hvtz5m9I4lrkOg24Rprw/1VKiSBAVmV6gvnyil0RjVdnJRpHCdSRIDo=
		},
		{
			name:           "Invalid encoding type",
			rawValue:       []byte("test_value"),
			masterPassword: "password123",
			initialValue:   SecretsValueRaw("e;invalid_encoding;aes128;"),
			expectedError:  "", //e;base64;aes128;0tHQ7cuvDux8UJ6uw5sihIW0Lmq6z8OyVkZzz7VJ5UpCmXVoBbV/FqFbtoQ=
		},
		{
			name:           "Invalid crypt mode",
			rawValue:       []byte("test_value"),
			masterPassword: "password123",
			initialValue:   SecretsValueRaw("x;base64;aes128;"),
			expectedError:  "", //e;base64;aes128;8HhlTMzMba/jCjhwncKmxcytdAN68+weyKQXg8U0UksdZHxJuF4aIuj0WtA=
		},
		{
			name:           "Empty raw value",
			rawValue:       nil,
			masterPassword: "password123",
			initialValue:   SecretsValueRaw("e;base64;aes128;"),
			expectedError:  "rawValue cannot be empty",
		},
		{
			name:           "Empty master password",
			rawValue:       []byte("test_value"),
			masterPassword: "",
			initialValue:   SecretsValueRaw("e;base64;aes128;"),
			expectedError:  "masterPassword cannot be empty",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := test.initialValue
			err := s.Encode(test.rawValue, test.masterPassword)

			if test.expectedError != "" {
				if err == nil || !strings.Contains(err.Error(), test.expectedError) {
					t.Errorf("expected error containing %q, got %v", test.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}

				rawValue, err := s.Decode(test.masterPassword)
				if err != nil {
					t.Errorf("unexpected error on decode check: %v", err)
				}
				// Compare two []byte slices
				if !bytes.Equal(rawValue, test.rawValue) {
					t.Errorf("expected decoded value to be %v, got %v", test.rawValue, rawValue)
				}
			}
		})
	}
}

func TestIsBase64Encoded(t *testing.T) {
	validBase64 := []byte("dGVzdC1zZWNyZXQtdmFsdWU=") // "test-secret-value" in base64
	invalidBase64 := []byte("not-a-base64-string")
	rawBytes := []byte("test-secret-value")

	assert.True(t, IsBase64Encoded(validBase64), "Should detect valid base64 encoding")
	assert.False(t, IsBase64Encoded(invalidBase64), "Should detect invalid base64 encoding")
	assert.False(t, IsBase64Encoded(rawBytes), "Should detect raw bytes are not base64")
}

func TestNewSecretsValueRawBase64Decrypted(t *testing.T) {
	tests := []struct {
		name          string
		inputValue    []byte
		inputEncType  EncryptionType
		expectedEnc   EncryptionType
		expectsBase64 bool
		expectedValue []byte
	}{
		{
			name:          "Raw bytes are encoded",
			inputValue:    []byte("test-secret-value"),
			inputEncType:  ENCRYPTIONTYPE_AES256,
			expectedEnc:   ENCRYPTIONTYPE_AES256,
			expectsBase64: true,
			expectedValue: []byte("test-secret-value"),
		},
		{
			name:          "Already base64-encoded",
			inputValue:    []byte("dGVzdC1zZWNyZXQtdmFsdWU="),
			inputEncType:  ENCRYPTIONTYPE_AES256,
			expectedEnc:   ENCRYPTIONTYPE_AES256,
			expectsBase64: true,
			expectedValue: []byte("test-secret-value"),
		},
		{
			name:          "Default encryption type",
			inputValue:    []byte("another-secret-value"),
			inputEncType:  "",
			expectedEnc:   ENCRYPTIONTYPE_AES256,
			expectsBase64: true,
			expectedValue: []byte("another-secret-value"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewSecretsValueRawBase64Decrypted(tt.inputEncType, tt.inputValue)

			// Parse the result to verify its structure
			parts := strings.Split(string(result), ";")
			assert.Len(t, parts, 4, "SecretsValueRaw should have 4 parts")

			// Validate parts
			assert.Equal(t, CRYPTMODE_DECRYPTED, CryptMode(parts[0]), "Mode should be DECRYPTED")
			assert.Equal(t, ENCODINGTYPE_BASE64, EncodingType(parts[1]), "Encoding type should be BASE64")
			assert.Equal(t, string(tt.expectedEnc), parts[2], "Encryption type should match expected")

			// Decode and compare value if applicable
			decodedValue, err := base64.StdEncoding.DecodeString(parts[3])
			assert.NoError(t, err, "Value should decode correctly")
			assert.Equal(t, tt.expectedValue, decodedValue, "Decoded value should match original")
		})
	}
}
