package acrypt

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/pbkdf2"
	"hash"
	"strconv"
	"strings"
)

// PBKDF2Presets holds the configuration parameters for PBKDF2.
type PBKDF2Presets struct {
	Iterations int    // Number of iterations (OWASP min: 600,000 for SHA-256)
	KeyLen     int    // Derived key length (32 bytes default)
	HashFunc   string // Hash function: "sha256" or "sha512"
}

// NewPBKDF2Presets creates a new PBKDF2Presets with OWASP-recommended defaults (SHA-256, 600,000 iterations).
func NewPBKDF2Presets() *PBKDF2Presets {
	return &PBKDF2Presets{
		Iterations: 600000,
		KeyLen:     32,
		HashFunc:   "sha256",
	}
}

// HashPBKDF2 hashes a string using PBKDF2 and returns the hash in PHC format.
// It returns an error if the input string is empty or if hashing fails.
func HashPBKDF2(target string, presets *PBKDF2Presets) (string, error) {
	if target == "" {
		return "", fmt.Errorf("cannot hash an empty string")
	}
	if presets == nil {
		presets = NewPBKDF2Presets()
	}
	if presets.Iterations < 600000 {
		return "", fmt.Errorf("iterations too low; minimum 600,000 for security")
	}
	if presets.KeyLen < 32 {
		return "", fmt.Errorf("key length too short; minimum 32 bytes")
	}
	if presets.HashFunc != "sha256" && presets.HashFunc != "sha512" {
		return "", fmt.Errorf("unsupported hash function: %s (use sha256 or sha512)", presets.HashFunc)
	}

	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	var hashFunc func() hash.Hash
	switch presets.HashFunc {
	case "sha256":
		hashFunc = sha256.New
	case "sha512":
		hashFunc = sha512.New
	}

	hash := pbkdf2.Key([]byte(target), salt, presets.Iterations, presets.KeyLen, hashFunc)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	format := "$pbkdf2-%s$i=%d$%s$%s"
	full := fmt.Sprintf(format, presets.HashFunc, presets.Iterations, b64Salt, b64Hash)
	return full, nil
}

// VerifyPBKDF2 compares a PBKDF2 hashed string with a plaintext string.
// It returns true if they match, false otherwise, along with an error if any occurs.
func VerifyPBKDF2(hashed string, plain string) (bool, error) {
	if hashed == "" || plain == "" {
		return false, fmt.Errorf("empty parameters 'hashed' or 'plain'")
	}

	parts := strings.Split(hashed, "$")
	if len(parts) != 5 || !strings.HasPrefix(parts[1], "pbkdf2-") {
		return false, fmt.Errorf("invalid PBKDF2 hash format")
	}

	hashFuncStr := strings.TrimPrefix(parts[1], "pbkdf2-")
	iterations, err := strconv.Atoi(parts[2][2:]) // Skip "i="
	if err != nil {
		return false, fmt.Errorf("failed to parse iterations: %w", err)
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[3])
	if err != nil {
		return false, fmt.Errorf("failed to decode salt: %w", err)
	}

	decodedHash, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, fmt.Errorf("failed to decode hash: %w", err)
	}

	var hashFunc func() hash.Hash
	switch hashFuncStr {
	case "sha256":
		hashFunc = sha256.New
	case "sha512":
		hashFunc = sha512.New
	default:
		return false, fmt.Errorf("unsupported hash function in hash: %s", hashFuncStr)
	}

	comparisonHash := pbkdf2.Key([]byte(plain), salt, iterations, len(decodedHash), hashFunc)

	return subtle.ConstantTimeCompare(decodedHash, comparisonHash) == 1, nil
}

// IsPBKDF2Hash checks if a given hash is a PBKDF2 hash.
func IsPBKDF2Hash(hash string) bool {
	parts := strings.Split(hash, "$")
	return len(parts) == 5 && strings.HasPrefix(parts[1], "pbkdf2-")
}

// IsFIPSModeFunc is the function type for checking FIPS mode.
type IsFIPSModeFunc func() bool

// DefaultIsFIPSMode is the default implementation.
func DefaultIsFIPSMode() bool {
	return false // fips140.Enabled()
}

// IsFIPSMode is the global variable for function IsFIPSModeFunc (override in tests).
var IsFIPSMode IsFIPSModeFunc = DefaultIsFIPSMode

// HashPassword hashes the password using Argon2id (default) or PBKDF2 (if FIPS mode).
// Optionally accepts presets; uses defaults if nil.
func HashPassword(password string, presets interface{}) (string, error) {
	if password == "" {
		return "", fmt.Errorf("cannot hash an empty password")
	}

	if IsFIPSMode() {
		pbkdf2Presets, ok := presets.(*PBKDF2Presets)
		if !ok && presets != nil {
			return "", fmt.Errorf("invalid presets for FIPS mode; expect *PBKDF2Presets")
		}
		return HashPBKDF2(password, pbkdf2Presets)
	}

	argon2Presets, ok := presets.(*Argon2Presets)
	if !ok && presets != nil {
		return "", fmt.Errorf("invalid presets; expect *Argon2Presets")
	}
	return HashArgon2id(password, argon2Presets)
}

// MatchPassword verifies if the password matches the stored hash.
// Automatically detects and routes to PBKDF2 or Argon2id verifier based on hash format.
func MatchPassword(hashed string, password string) (bool, error) {
	if hashed == "" || password == "" {
		return false, fmt.Errorf("empty parameters 'hashed' or 'password'")
	}

	if IsPBKDF2Hash(hashed) {
		return VerifyPBKDF2(hashed, password)
	}
	if IsArgon2idHash(hashed) {
		return VerifyArgon2id(hashed, password)
	}
	return false, fmt.Errorf("unsupported hash format")
}
