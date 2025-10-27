package acrypt

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/opencontainers/go-digest"
)

// Notes:
// - SHA-256, SHA-384, and SHA-512 are FIPS 140-3 compliant hash algorithms (per FIPS 180-4).
// - For FIPS mode in Go, build with GOEXPERIMENT=systemcrypto to use certified implementations.
// - These functions are for data integrity/checksums, not password storage (use PBKDF2 for passwords).

// prependCryptFormatBase64 prepends a format label to a base64-encoded string.
func prependCryptFormatBase64(format string, b []byte) string {
	return fmt.Sprintf("%s%s", format, base64.RawStdEncoding.EncodeToString(b))
}

// EncodeToBase64 encodes the given bytes to a base64-encoded string using raw standard encoding (URL-safe, no padding).
func EncodeToBase64(b []byte) string {
	return base64.RawStdEncoding.EncodeToString(b)
}

// SHA-256 Functions

// Sum256ToSlice computes the SHA-256 hash of the input bytes and returns it as a slice.
func Sum256ToSlice(b []byte) []byte {
	hash := sha256.Sum256(b)
	return hash[:]
}

// HashSHA256Bytes computes the SHA-256 hash of the input bytes.
func HashSHA256Bytes(b []byte) []byte {
	return Sum256ToSlice(b)
}

// HashSHA256String computes the SHA-256 hash of the input string.
func HashSHA256String(s string) []byte {
	return HashSHA256Bytes([]byte(s))
}

// HashSHA256File computes the SHA-256 hash of the file at the given filepath.
func HashSHA256File(filepath string) ([]byte, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, f); err != nil {
		return nil, fmt.Errorf("failed to calculate SHA-256 hash: %w", err)
	}
	return hash.Sum(nil), nil
}

// ToHexSHA256 computes the SHA-256 hash of the input bytes and returns it as a hexadecimal string.
func ToHexSHA256(b []byte) string {
	return hex.EncodeToString(HashSHA256Bytes(b))
}

// ToHexSHA256String computes the SHA-256 hash of the input string and returns it as a hexadecimal string.
func ToHexSHA256String(s string) string {
	return ToHexSHA256([]byte(s))
}

// ToHexSHA256File computes the SHA-256 hash of the file and returns it as a hexadecimal string.
func ToHexSHA256File(filepath string) (string, error) {
	checksum, err := HashSHA256File(filepath)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(checksum), nil
}

// ToBase64SHA256 computes the SHA-256 hash of the input bytes and returns it as a base64-encoded string.
func ToBase64SHA256(b []byte) string {
	return EncodeToBase64(HashSHA256Bytes(b))
}

// ToBase64SHA256String computes the SHA-256 hash of the input string and returns it as a base64-encoded string.
func ToBase64SHA256String(s string) string {
	return ToBase64SHA256([]byte(s))
}

// ToBase64SHA256File computes the SHA-256 hash of the file and returns it as a base64-encoded string.
func ToBase64SHA256File(filepath string) (string, error) {
	checksum, err := HashSHA256File(filepath)
	if err != nil {
		return "", err
	}
	return EncodeToBase64(checksum), nil
}

// ToBase64SHA256WithFormat computes the SHA-256 hash of the input string and returns it as a base64-encoded string,
// optionally prepending a "{sha256}" format label.
func ToBase64SHA256WithFormat(s string, prependFormat bool) string {
	hash := HashSHA256String(s)
	format := ""
	if prependFormat {
		format = "{sha256}"
	}
	return prependCryptFormatBase64(format, hash)
}

// ToBase64SHA256FileWithFormat computes the SHA-256 hash of the file and returns it as a base64-encoded string,
// optionally prepending a "{sha256}" format label.
func ToBase64SHA256FileWithFormat(filepath string, prependFormat bool) (string, error) {
	checksum, err := HashSHA256File(filepath)
	if err != nil {
		return "", err
	}
	checksumBase64 := EncodeToBase64(checksum)
	if prependFormat {
		checksumBase64 = "{sha256}" + checksumBase64
	}
	return checksumBase64, nil
}

// SHA-384 Functions

// Sum384ToSlice computes the SHA-384 hash of the input bytes and returns it as a slice.
func Sum384ToSlice(b []byte) []byte {
	hash := sha512.Sum384(b)
	return hash[:]
}

// HashSHA384Bytes computes the SHA-384 hash of the input bytes.
func HashSHA384Bytes(b []byte) []byte {
	return Sum384ToSlice(b)
}

// HashSHA384String computes the SHA-384 hash of the input string.
func HashSHA384String(s string) []byte {
	return HashSHA384Bytes([]byte(s))
}

// HashSHA384File computes the SHA-384 hash of the file at the given filepath.
func HashSHA384File(filepath string) ([]byte, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	hash := sha512.New384()
	if _, err := io.Copy(hash, f); err != nil {
		return nil, fmt.Errorf("failed to calculate SHA-384 hash: %w", err)
	}
	return hash.Sum(nil), nil
}

// ToHexSHA384 computes the SHA-384 hash of the input bytes and returns it as a hexadecimal string.
func ToHexSHA384(b []byte) string {
	return hex.EncodeToString(HashSHA384Bytes(b))
}

// ToHexSHA384String computes the SHA-384 hash of the input string and returns it as a hexadecimal string.
func ToHexSHA384String(s string) string {
	return ToHexSHA384([]byte(s))
}

// ToHexSHA384File computes the SHA-384 hash of the file and returns it as a hexadecimal string.
func ToHexSHA384File(filepath string) (string, error) {
	checksum, err := HashSHA384File(filepath)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(checksum), nil
}

// ToBase64SHA384 computes the SHA-384 hash of the input bytes and returns it as a base64-encoded string.
func ToBase64SHA384(b []byte) string {
	return EncodeToBase64(HashSHA384Bytes(b))
}

// ToBase64SHA384String computes the SHA-384 hash of the input string and returns it as a base64-encoded string.
func ToBase64SHA384String(s string) string {
	return ToBase64SHA384([]byte(s))
}

// ToBase64SHA384File computes the SHA-384 hash of the file and returns it as a base64-encoded string.
func ToBase64SHA384File(filepath string) (string, error) {
	checksum, err := HashSHA384File(filepath)
	if err != nil {
		return "", err
	}
	return EncodeToBase64(checksum), nil
}

// ToBase64SHA384WithFormat computes the SHA-384 hash of the input string and returns it as a base64-encoded string,
// optionally prepending a "{sha384}" format label.
func ToBase64SHA384WithFormat(s string, prependFormat bool) string {
	hash := HashSHA384String(s)
	format := ""
	if prependFormat {
		format = "{sha384}"
	}
	return prependCryptFormatBase64(format, hash)
}

// ToBase64SHA384FileWithFormat computes the SHA-384 hash of the file and returns it as a base64-encoded string,
// optionally prepending a "{sha384}" format label.
func ToBase64SHA384FileWithFormat(filepath string, prependFormat bool) (string, error) {
	checksum, err := HashSHA384File(filepath)
	if err != nil {
		return "", err
	}
	checksumBase64 := EncodeToBase64(checksum)
	if prependFormat {
		checksumBase64 = "{sha384}" + checksumBase64
	}
	return checksumBase64, nil
}

// SHA-512 Functions

// Sum512ToSlice computes the SHA-512 hash of the input bytes and returns it as a slice.
func Sum512ToSlice(b []byte) []byte {
	hash := sha512.Sum512(b)
	return hash[:]
}

// HashSHA512Bytes computes the SHA-512 hash of the input bytes.
func HashSHA512Bytes(b []byte) []byte {
	return Sum512ToSlice(b)
}

// HashSHA512String computes the SHA-512 hash of the input string.
func HashSHA512String(s string) []byte {
	return HashSHA512Bytes([]byte(s))
}

// HashSHA512File computes the SHA-512 hash of the file at the given filepath.
func HashSHA512File(filepath string) ([]byte, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	hash := sha512.New()
	if _, err := io.Copy(hash, f); err != nil {
		return nil, fmt.Errorf("failed to calculate SHA-512 hash: %w", err)
	}
	return hash.Sum(nil), nil
}

// ToHexSHA512 computes the SHA-512 hash of the input bytes and returns it as a hexadecimal string.
func ToHexSHA512(b []byte) string {
	return hex.EncodeToString(HashSHA512Bytes(b))
}

// ToHexSHA512String computes the SHA-512 hash of the input string and returns it as a hexadecimal string.
func ToHexSHA512String(s string) string {
	return ToHexSHA512([]byte(s))
}

// ToHexSHA512File computes the SHA-512 hash of the file and returns it as a hexadecimal string.
func ToHexSHA512File(filepath string) (string, error) {
	checksum, err := HashSHA512File(filepath)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(checksum), nil
}

// ToBase64SHA512 computes the SHA-512 hash of the input bytes and returns it as a base64-encoded string.
func ToBase64SHA512(b []byte) string {
	return EncodeToBase64(HashSHA512Bytes(b))
}

// ToBase64SHA512String computes the SHA-512 hash of the input string and returns it as a base64-encoded string.
func ToBase64SHA512String(s string) string {
	return ToBase64SHA512([]byte(s))
}

// ToBase64SHA512File computes the SHA-512 hash of the file and returns it as a base64-encoded string.
func ToBase64SHA512File(filepath string) (string, error) {
	checksum, err := HashSHA512File(filepath)
	if err != nil {
		return "", err
	}
	return EncodeToBase64(checksum), nil
}

// ToBase64SHA512WithFormat computes the SHA-512 hash of the input string and returns it as a base64-encoded string,
// optionally prepending a "{sha512}" format label.
func ToBase64SHA512WithFormat(s string, prependFormat bool) string {
	hash := HashSHA512String(s)
	format := ""
	if prependFormat {
		format = "{sha512}"
	}
	return prependCryptFormatBase64(format, hash)
}

// ToBase64SHA512FileWithFormat computes the SHA-512 hash of the file and returns it as a base64-encoded string,
// optionally prepending a "{sha512}" format label.
func ToBase64SHA512FileWithFormat(filepath string, prependFormat bool) (string, error) {
	checksum, err := HashSHA512File(filepath)
	if err != nil {
		return "", err
	}
	checksumBase64 := EncodeToBase64(checksum)
	if prependFormat {
		checksumBase64 = "{sha512}" + checksumBase64
	}
	return checksumBase64, nil
}

// ParseChecksum parses a checksum string (e.g., "sha256:abcdef...") into the algorithm (lowercase) and raw hash bytes.
// It uses go-digest for OCI-compliant validation and extraction.
// Only FIPS 140-3 compliant algorithms are allowed: sha256, sha384, sha512.
// Returns an error for invalid format, unsupported algorithm, or decoding issues.
func ParseChecksum(d string) (alg string, rawBytes []byte, err error) {
	parsedDigest, err := digest.Parse(d)
	if err != nil {
		return "", nil, fmt.Errorf("invalid checksum format: %w", err)
	}

	alg = strings.ToLower(parsedDigest.Algorithm().String())

	// Allowlist: FIPS-compliant algorithms
	allowedAlgs := map[string]struct{}{
		"sha256": {},
		"sha384": {},
		"sha512": {},
	}
	if _, ok := allowedAlgs[alg]; !ok {
		return "", nil, fmt.Errorf("unsupported algorithm: %s (allowed: sha256, sha384, sha512)", alg)
	}

	hexEncoded := parsedDigest.Encoded()
	rawBytes, err = hex.DecodeString(hexEncoded)
	if err != nil {
		return "", nil, fmt.Errorf("failed to decode hex: %w", err)
	}

	return alg, rawBytes, nil
}

// IsValidChecksumFingerprint checks if the given hex string is a valid representation for the specified algorithm.
// Supports sha256 (64 chars), sha384 (96 chars), sha512 (128 chars) for FIPS/OCI compliance.
// Verifies length and ensures only hex characters (0-9, a-f, A-F) are present.
// Case-insensitive for hex digits.
func IsValidChecksumFingerprint(alg string, s string) bool {
	alg = strings.ToLower(alg)
	var expectedLen int
	switch alg {
	case "sha256":
		expectedLen = 64
	case "sha384":
		expectedLen = 96
	case "sha512":
		expectedLen = 128
	default:
		return false // Unsupported alg
	}

	if len(s) != expectedLen {
		return false
	}

	s = strings.ToLower(s) // Normalize for check
	for _, c := range s {
		if (c < '0' || c > '9') && (c < 'a' || c > 'f') {
			return false
		}
	}
	return true
}
