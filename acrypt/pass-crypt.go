package acrypt

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/scrypt"
	"strings"
)

// EncryptBCrypt hashes a string using bcrypt and returns the hash.
// It returns an error if the input string is empty or if the hashing fails.
func EncryptBCrypt(target string) (string, error) {
	if target == "" {
		return "", fmt.Errorf("cannot crypt an empty string")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(target), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("bcrypt.GenerateFromPassword failed; %v", err)
	}

	return string(hash), nil
}

// IsEqualBCrypt compares a bcrypt hashed string with a plaintext string.
// It returns true if they match, false otherwise.
func IsEqualBCrypt(hashed string, plain string) bool {
	if hashed == "" || plain == "" {
		return false
	}

	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
	return err == nil
}

// ScryptPresets holds the configuration parameters for scrypt.
type ScryptPresets struct {
	N      int // CPU/memory cost parameter
	R      int // Block size parameter
	P      int // Parallelization parameter
	KeyLen int // Key length
}

// NewScryptPresets creates a new ScryptPresets with default values.
func NewScryptPresets() *ScryptPresets {
	return &ScryptPresets{
		N:      32768,
		R:      8,
		P:      1,
		KeyLen: 32,
	}
}

// EncryptSCrypt hashes a string using scrypt and returns the hash.
// It returns an error if the input string is empty or if the hashing fails.
func EncryptSCrypt(unencryptedTarget string, presets *ScryptPresets) (string, error) {
	if presets == nil {
		presets = NewScryptPresets()
	}

	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash, err := scrypt.Key([]byte(unencryptedTarget), salt, presets.N, presets.R, presets.P, presets.KeyLen)
	if err != nil {
		return "", fmt.Errorf("failed to generate scrypt.Key; %v", err)
	}

	b64Salt := base64.StdEncoding.EncodeToString(salt)
	b64Hash := base64.StdEncoding.EncodeToString(hash)

	format := "$scrypt$n=%d,r=%d,p=%d$%s$%s"
	full := fmt.Sprintf(format, presets.N, presets.R, presets.P, b64Salt, b64Hash)
	return full, nil
}

// IsEqualSCrypt compares a scrypt hashed string with a plaintext string.
// It returns true if they match, false otherwise, along with an error if any occurs.
func IsEqualSCrypt(hashed string, plain string) (bool, error) {
	if hashed == "" || plain == "" {
		return false, fmt.Errorf("empty parameters 'hashed' and 'plain'")
	}

	parts := strings.Split(hashed, "$")
	presets := &ScryptPresets{}
	_, err := fmt.Sscanf(parts[2], "n=%d,r=%d,p=%d", &presets.N, &presets.R, &presets.P)
	if err != nil {
		return false, err
	}

	salt, err := base64.StdEncoding.DecodeString(parts[3])
	if err != nil {
		return false, err
	}

	decodedHash, err := base64.StdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}
	presets.KeyLen = len(decodedHash)

	hash, err := scrypt.Key([]byte(plain), salt, presets.N, presets.R, presets.P, presets.KeyLen)
	if err != nil {
		return false, fmt.Errorf("failed to decode scrypt.Key; %v", err)
	}

	return subtle.ConstantTimeCompare(decodedHash, hash) == 1, nil
}

// Argon2Presets holds the configuration parameters for Argon2.
type Argon2Presets struct {
	Time    uint32 // Time cost parameter
	Memory  uint32 // Memory cost parameter
	Threads uint8  // Parallelism parameter
	KeyLen  uint32 // Key length
}

// NewArgon2Presets creates a new Argon2Presets with default values.
func NewArgon2Presets() *Argon2Presets {
	return &Argon2Presets{
		Time:    1,
		Memory:  65536,
		Threads: 4,
		KeyLen:  32,
	}
}

// EncryptArgon2id hashes a string using Argon2id and returns the hash.
// It returns an error if the input string is empty or if the hashing fails.
func EncryptArgon2id(unencryptedTarget string, presets *Argon2Presets) (string, error) {
	if presets == nil {
		presets = NewArgon2Presets()
	}

	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(unencryptedTarget), salt, presets.Time, presets.Memory, presets.Threads, presets.KeyLen)

	b64Salt := base64.StdEncoding.EncodeToString(salt)
	b64Hash := base64.StdEncoding.EncodeToString(hash)

	format := "$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s"
	full := fmt.Sprintf(format, argon2.Version, presets.Memory, presets.Time, presets.Threads, b64Salt, b64Hash)
	return full, nil
}

// IsEqualArgon2id compares an Argon2id hashed string with a plaintext string.
// It returns true if they match, false otherwise, along with an error if any occurs.
func IsEqualArgon2id(hashed string, plain string) (bool, error) {
	if hashed == "" || plain == "" {
		return false, fmt.Errorf("empty parameters 'hashed' and 'plain'")
	}

	parts := strings.Split(hashed, "$")
	if len(parts) != 6 {
		return false, fmt.Errorf("invalid hash format")
	}

	var presets Argon2Presets
	_, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &presets.Memory, &presets.Time, &presets.Threads)
	if err != nil {
		return false, err
	}

	salt, err := base64.StdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, fmt.Errorf("error decoding salt: %w", err)
	}

	decodedHash, err := base64.StdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, fmt.Errorf("error decoding hash: %w", err)
	}
	presets.KeyLen = uint32(len(decodedHash))

	comparisonHash := argon2.IDKey([]byte(plain), salt, presets.Time, presets.Memory, presets.Threads, presets.KeyLen)

	return subtle.ConstantTimeCompare(decodedHash, comparisonHash) == 1, nil
}

// IsArgon2idHash checks if a given hash is an Argon2id hash.
func IsArgon2idHash(hash string) bool {
	parts := strings.Split(hash, "$")
	if len(parts) != 6 {
		return false
	}

	// Check if the hash starts with the Argon2id identifier
	if parts[1] != "argon2id" {
		return false
	}

	return true
}
