package acrypt

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/scrypt"
	"math"
	"strconv"
	"strings"
)

// HashBCrypt hashes a string using bcrypt and returns the hash.
// It returns an error if the input string is empty or if the hashing fails.
func HashBCrypt(target string, cost int) (string, error) {
	if target == "" {
		return "", fmt.Errorf("cannot hash an empty string")
	}
	if cost == 0 {
		cost = bcrypt.DefaultCost
	}
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		return "", fmt.Errorf("invalid cost: must be between %d and %d", bcrypt.MinCost, bcrypt.MaxCost)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(target), cost)
	if err != nil {
		return "", fmt.Errorf("bcrypt.GenerateFromPassword failed: %w", err)
	}

	return string(hash), nil
}

// VerifyBCrypt compares a bcrypt hashed string with a plaintext string.
// It returns true if they match, false otherwise, along with an error if the hash is invalid or other issues occur.
func VerifyBCrypt(hashed string, plain string) (bool, error) {
	if hashed == "" || plain == "" {
		return false, fmt.Errorf("empty parameters 'hashed' or 'plain'")
	}

	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
	if err == nil {
		return true, nil
	}
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return false, nil
	}
	return false, fmt.Errorf("invalid bcrypt hash: %w", err)
}

// IsBCryptHash checks if a given hash is a bcrypt hash.
func IsBCryptHash(hash string) bool {
	return strings.HasPrefix(hash, "$2a$") || strings.HasPrefix(hash, "$2b$") || strings.HasPrefix(hash, "$2y$")
}

// ScryptPresets holds the configuration parameters for scrypt.
type ScryptPresets struct {
	N      int // CPU/memory cost parameter (must be power of 2)
	R      int // Block size parameter
	P      int // Parallelization parameter
	KeyLen int // Key length
}

// NewScryptPresets creates a new ScryptPresets with OWASP-recommended minimum values.
func NewScryptPresets() *ScryptPresets {
	return &ScryptPresets{
		N:      131072, // 2^17
		R:      8,
		P:      1,
		KeyLen: 32,
	}
}

// HashSCrypt hashes a string using scrypt and returns the hash in PHC format.
// It returns an error if the input string is empty or if the hashing fails.
func HashSCrypt(unencryptedTarget string, presets *ScryptPresets) (string, error) {
	if unencryptedTarget == "" {
		return "", fmt.Errorf("cannot hash an empty string")
	}
	if presets == nil {
		presets = NewScryptPresets()
	}
	if presets.N <= 1 || presets.R < 1 || presets.P < 1 || presets.KeyLen < 1 {
		return "", fmt.Errorf("invalid presets")
	}
	fl := math.Log2(float64(presets.N))
	ln := int(fl)
	if float64(ln) != fl {
		return "", fmt.Errorf("N must be a power of 2")
	}

	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash, err := scrypt.Key([]byte(unencryptedTarget), salt, presets.N, presets.R, presets.P, presets.KeyLen)
	if err != nil {
		return "", fmt.Errorf("failed to generate scrypt.Key: %w", err)
	}

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	format := "$scrypt$ln=%d,r=%d,p=%d$%s$%s"
	full := fmt.Sprintf(format, ln, presets.R, presets.P, b64Salt, b64Hash)
	return full, nil
}

// VerifySCrypt compares a scrypt hashed string with a plaintext string.
// It returns true if they match, false otherwise, along with an error if any occurs.
func VerifySCrypt(hashed string, plain string) (bool, error) {
	if hashed == "" || plain == "" {
		return false, fmt.Errorf("empty parameters 'hashed' or 'plain'")
	}

	parts := strings.Split(hashed, "$")
	if len(parts) != 5 || parts[1] != "scrypt" {
		return false, fmt.Errorf("invalid scrypt hash format")
	}

	presets := &ScryptPresets{}
	var ln int
	_, err := fmt.Sscanf(parts[2], "ln=%d,r=%d,p=%d", &ln, &presets.R, &presets.P)
	if err != nil {
		return false, fmt.Errorf("failed to parse parameters: %w", err)
	}
	if ln < 1 || presets.R < 1 || presets.P < 1 {
		return false, fmt.Errorf("invalid parameters")
	}
	presets.N = 1 << ln

	salt, err := base64.RawStdEncoding.DecodeString(parts[3])
	if err != nil {
		return false, fmt.Errorf("failed to decode salt: %w", err)
	}

	decodedHash, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, fmt.Errorf("failed to decode hash: %w", err)
	}
	presets.KeyLen = len(decodedHash)

	hash, err := scrypt.Key([]byte(plain), salt, presets.N, presets.R, presets.P, presets.KeyLen)
	if err != nil {
		return false, fmt.Errorf("failed to generate scrypt.Key: %w", err)
	}

	return subtle.ConstantTimeCompare(decodedHash, hash) == 1, nil
}

// IsSCryptHash checks if a given hash is a scrypt hash.
func IsSCryptHash(hash string) bool {
	parts := strings.Split(hash, "$")
	return len(parts) == 5 && parts[1] == "scrypt"
}

// Argon2Presets holds the configuration parameters for Argon2.
type Argon2Presets struct {
	Time    uint32 // Time cost parameter
	Memory  uint32 // Memory cost parameter (in KiB)
	Threads uint8  // Parallelism parameter
	KeyLen  uint32 // Key length
}

// NewArgon2Presets creates a new Argon2Presets with OWASP-recommended minimum values.
func NewArgon2Presets() *Argon2Presets {
	return &Argon2Presets{
		Time:    2,
		Memory:  19456, // 19 MiB
		Threads: 1,
		KeyLen:  32,
	}
}

// HashArgon2id hashes a string using Argon2id and returns the hash in PHC format.
// It returns an error if the input string is empty or if the hashing fails.
func HashArgon2id(unencryptedTarget string, presets *Argon2Presets) (string, error) {
	if unencryptedTarget == "" {
		return "", fmt.Errorf("cannot hash an empty string")
	}
	if presets == nil {
		presets = NewArgon2Presets()
	}

	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(unencryptedTarget), salt, presets.Time, presets.Memory, presets.Threads, presets.KeyLen)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	format := "$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s"
	full := fmt.Sprintf(format, argon2.Version, presets.Memory, presets.Time, presets.Threads, b64Salt, b64Hash)
	return full, nil
}

// VerifyArgon2id compares an Argon2id hashed string with a plaintext string.
// It returns true if they match, false otherwise, along with an error if any occurs.
func VerifyArgon2id(hashed string, plain string) (bool, error) {
	if hashed == "" || plain == "" {
		return false, fmt.Errorf("empty parameters 'hashed' or 'plain'")
	}

	parts := strings.Split(hashed, "$")
	if len(parts) != 6 {
		return false, fmt.Errorf("invalid hash format")
	}
	if parts[1] != "argon2id" {
		return false, fmt.Errorf("not an argon2id hash")
	}

	var version int
	if _, err := fmt.Sscanf(parts[2], "v=%d", &version); err != nil {
		return false, fmt.Errorf("failed to parse version: %w", err)
	}
	if version != argon2.Version {
		return false, fmt.Errorf("unsupported argon2 version %d", version)
	}

	var presets Argon2Presets
	params := strings.Split(parts[3], ",")
	for _, param := range params {
		kv := strings.Split(param, "=")
		if len(kv) != 2 {
			return false, fmt.Errorf("invalid param: %s", param)
		}
		val, err := strconv.Atoi(kv[1])
		if err != nil {
			return false, err
		}
		switch kv[0] {
		case "m":
			presets.Memory = uint32(val)
		case "t":
			presets.Time = uint32(val)
		case "p":
			presets.Threads = uint8(val)
		default:
			return false, fmt.Errorf("unknown param: %s", kv[0])
		}
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, fmt.Errorf("error decoding salt: %w", err)
	}

	decodedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
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
	if len(parts) != 6 || parts[1] != "argon2id" {
		return false
	}
	var version int
	if _, err := fmt.Sscanf(parts[2], "v=%d", &version); err != nil || version != argon2.Version {
		return false
	}
	return true
}
