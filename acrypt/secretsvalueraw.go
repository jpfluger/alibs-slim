package acrypt

import (
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

type SecretsValueRaw string

// IsEmpty checks if SecretsValueRaw is empty after trimming whitespace.
func (s SecretsValueRaw) IsEmpty() bool {
	return strings.TrimSpace(string(s)) == ""
}

// Parse splits the SecretsValueRaw into its components.
func (s SecretsValueRaw) Parse() (CryptMode, EncodingType, EncryptionType, string, error) {
	parts := strings.SplitN(string(s), ";", 4)
	if len(parts) != 4 {
		return "", "", "", "", errors.New("invalid format: must have 4 parts")
	}

	mode := CryptMode(parts[0])
	if mode != CRYPTMODE_ENCRYPTED && mode != CRYPTMODE_DECRYPTED {
		return "", "", "", "", errors.New("invalid crypt mode")
	}

	encoding := EncodingType(parts[1])
	if encoding != ENCODINGTYPE_BASE64 && encoding != ENCODINGTYPE_HEX && encoding != ENCODINGTYPE_PLAIN {
		return "", "", "", "", errors.New("invalid encoding type")
	}

	encryption := EncryptionType(parts[2])
	if encryption != ENCRYPTIONTYPE_AES128 && encryption != ENCRYPTIONTYPE_AES256 {
		return "", "", "", "", errors.New("invalid encryption type")
	}

	return mode, encoding, encryption, parts[3], nil
}

// IsBase64Encoded checks if the given byte slice is valid base64-encoded data.
func IsBase64Encoded(data []byte) bool {
	// Attempt to decode the data
	_, err := base64.StdEncoding.DecodeString(string(data))
	return err == nil
}

func NewSecretsValueRawBase64Decrypted(encryptionType EncryptionType, value []byte) SecretsValueRaw {
	var encodedValue string

	// Check if the value is already base64-encoded
	if IsBase64Encoded(value) {
		encodedValue = string(value)
	} else {
		encodedValue = base64.StdEncoding.EncodeToString(value)
	}

	if encryptionType.IsEmpty() {
		encryptionType = ENCRYPTIONTYPE_AES256
	}

	return SecretsValueRaw(fmt.Sprintf("%s;%s;%s;%s", CRYPTMODE_DECRYPTED, ENCODINGTYPE_BASE64, encryptionType, encodedValue))
}

// Validate ensures the raw value is in the proper format. If not, it defaults to a clear format.
func (s *SecretsValueRaw) Validate(value string) {
	parts := strings.SplitN(value, ";", 4)
	if len(parts) != 4 {
		*s = SecretsValueRaw(fmt.Sprintf("%s;%s;%s;%s", CRYPTMODE_DECRYPTED, ENCODINGTYPE_PLAIN, ENCRYPTIONTYPE_AES256, value))
		return
	}
	*s = SecretsValueRaw(value)
}

// Decode decrypts and decodes the value based on the encoding and encryption types.
func (s SecretsValueRaw) Decode(password string) ([]byte, error) {
	mode, encoding, encryption, value, err := s.Parse()
	if err != nil {
		return nil, err
	}

	if mode == CRYPTMODE_ENCRYPTED {
		var decrypted []byte
		cipherText, err := base64.StdEncoding.DecodeString(value)
		if err != nil {
			return nil, err
		}

		switch encryption {
		case ENCRYPTIONTYPE_AES128:
			decrypted, err = AESGCM128Decrypt(cipherText, password)
		case ENCRYPTIONTYPE_AES256:
			decrypted, err = AESGCM256Decrypt(cipherText, password)
		default:
			return nil, errors.New("unsupported encryption type")
		}

		if err != nil {
			return nil, err
		}
		value = string(decrypted)
	}

	var decoded []byte
	if encoding == ENCODINGTYPE_BASE64 {
		decoded, err = base64.StdEncoding.DecodeString(value)
		if err != nil {
			return nil, err
		}
	} else if encoding == ENCODINGTYPE_HEX {
		decoded, err = hex.DecodeString(value)
		if err != nil {
			return nil, err
		}
	} else {
		decoded = []byte(value)
	}

	return decoded, nil
}

// Encode changes the raw value to a new encrypted or plain value based on the current mode.
func (s *SecretsValueRaw) Encode(rawValue []byte, masterPassword string) error {
	if rawValue == nil || len(rawValue) == 0 {
		return fmt.Errorf("rawValue cannot be empty")
	}
	if strings.TrimSpace(masterPassword) == "" {
		return fmt.Errorf("masterPassword cannot be empty")
	}

	// Parse the existing value or use defaults if parsing fails.
	mode, encoding, encryption, _, err := s.Parse()
	if err != nil {
		if mode.IsEmpty() {
			mode = CRYPTMODE_ENCRYPTED
		}
		if encoding.IsEmpty() {
			encoding = ENCODINGTYPE_BASE64
		}
		if encryption.IsEmpty() {
			encryption = ENCRYPTIONTYPE_AES128
		}
	}

	// Encode the raw value based on the specified encoding type.
	var encoded string
	switch encoding {
	case ENCODINGTYPE_BASE64:
		encoded = base64.StdEncoding.EncodeToString(rawValue)
	case ENCODINGTYPE_HEX:
		encoded = hex.EncodeToString(rawValue)
	case ENCODINGTYPE_PLAIN:
		encoded = string(rawValue)
	default:
		return fmt.Errorf("unsupported encoding type: %s", encoding)
	}

	// Handle encryption or plain formatting based on the mode.
	switch mode {
	case CRYPTMODE_ENCRYPTED:
		// Encrypt the encoded value using the specified encryption type.
		var encryptedValue []byte
		switch encryption {
		case ENCRYPTIONTYPE_AES128:
			encryptedValue, err = AESGCM128Encrypt([]byte(encoded), masterPassword)
		case ENCRYPTIONTYPE_AES256:
			encryptedValue, err = AESGCM256Encrypt([]byte(encoded), masterPassword)
		default:
			return fmt.Errorf("unsupported encryption type: %s", encryption)
		}
		if err != nil {
			return fmt.Errorf("failed to encrypt value: %w", err)
		}

		// Base64 encode the encrypted value and update the SecretsValueRaw.
		encryptedEncodedValue := base64.StdEncoding.EncodeToString(encryptedValue)
		*s = SecretsValueRaw(fmt.Sprintf("%s;%s;%s;%s", CRYPTMODE_ENCRYPTED, encoding, encryption, encryptedEncodedValue))

	case CRYPTMODE_DECRYPTED:
		// Store the plain encoded value without encryption.
		*s = SecretsValueRaw(fmt.Sprintf("%s;%s;%s;%s", CRYPTMODE_DECRYPTED, encoding, encryption, encoded))

	default:
		return fmt.Errorf("invalid crypt mode: %s", mode)
	}

	return nil
}
