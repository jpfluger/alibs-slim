package acrypt

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

// Notes:
// - SHA-256 and SHA-512 are FIPS 140-3 compliant hash algorithms (per FIPS 180-4).
// - For FIPS mode in Go, build with GOEXPERIMENT=systemcrypto to use certified implementations.
// - These functions are for data integrity/checksums, not password storage (use PBKDF2 for passwords).

// FromStringToSHA256CheckSum generates a SHA-256 hash of the given string.
func FromStringToSHA256CheckSum(target string) []byte {
	hash := sha256.Sum256([]byte(target))
	return hash[:]
}

// FromBytesToSHA256Checksum generates a SHA-256 hash of the given bytes.
func FromBytesToSHA256Checksum(target []byte) []byte {
	hashArray := sha256.Sum256(target)
	return hashArray[:] // Convert the array to a slice
}

// FromFileToSHA256Checksum calculates the SHA-256 checksum of the file at the given filepath.
func FromFileToSHA256Checksum(filepath string) ([]byte, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, f); err != nil {
		return nil, fmt.Errorf("failed to calculate SHA-256 checksum: %v", err)
	}
	return hash.Sum(nil), nil
}

// FormatSHA256ChecksumHex converts a SHA-256 checksum into a hexadecimal string.
func FormatSHA256ChecksumHex(checksum []byte) string {
	return fmt.Sprintf("%x", checksum)
}

// FromStringToSHA512CheckSum generates a SHA-512 hash of the given string.
func FromStringToSHA512CheckSum(target string) []byte {
	hash := sha512.Sum512([]byte(target))
	return hash[:]
}

// FromBytesToSHA512Checksum generates a SHA-512 hash of the given bytes.
func FromBytesToSHA512Checksum(target []byte) []byte {
	hashArray := sha512.Sum512(target)
	return hashArray[:] // Convert the array to a slice
}

// FromFileToSHA512Checksum calculates the SHA-512 checksum of the file at the given filepath.
func FromFileToSHA512Checksum(filepath string) ([]byte, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	hash := sha512.New()
	if _, err := io.Copy(hash, f); err != nil {
		return nil, fmt.Errorf("failed to calculate SHA-512 checksum: %v", err)
	}
	return hash.Sum(nil), nil
}

// ToCheckSumSHA256 generates a hex-encoded SHA-256 hash of the given target.
func ToCheckSumSHA256(b []byte) string {
	return fmt.Sprintf("%x", sha256.Sum256(b))
}

// ToCheckSumSHA512 generates a hex-encoded SHA-512 hash of the given target.
func ToCheckSumSHA512(b []byte) string {
	return fmt.Sprintf("%x", sha512.Sum512(b))
}

// FormatSHA512ChecksumHex converts a SHA-512 checksum into a hexadecimal string.
func FormatSHA512ChecksumHex(checksum []byte) string {
	return fmt.Sprintf("%x", checksum)
}

// prependCryptFormatBase64 prepends a format label to a base64-encoded string.
func prependCryptFormatBase64(format string, b []byte) string {
	return fmt.Sprintf("%s%s", format, base64.RawStdEncoding.EncodeToString(b))
}

// ToHashSHA256Base64 generates a base64-encoded SHA-256 hash of the given target, optionally prepending a format label.
func ToHashSHA256Base64(target string, prependFormat bool) string {
	hash := sha256.Sum256([]byte(target))
	format := ""
	if prependFormat {
		format = "{sha256}"
	}
	return prependCryptFormatBase64(format, hash[:])
}

// ToHashSHA512Base64 generates a base64-encoded SHA-512 hash of the given target, optionally prepending a format label.
func ToHashSHA512Base64(target string, prependFormat bool) string {
	hash := sha512.Sum512([]byte(target))
	format := ""
	if prependFormat {
		format = "{sha512}"
	}
	return prependCryptFormatBase64(format, hash[:])
}

// FromFileToSHA256ChecksumBase64 calculates a base64-encoded SHA-256 checksum of the file at the given filepath, optionally prepending a format label.
func FromFileToSHA256ChecksumBase64(filepath string, prependFormat bool) (string, error) {
	// Use the existing function to calculate the SHA-256 checksum.
	checksum, err := FromFileToSHA256Checksum(filepath)
	if err != nil {
		return "", err
	}

	// Use the new function to encode the checksum to Base64.
	checksumBase64 := EncodeToBase64(checksum)

	// Prepend the format label if required.
	if prependFormat {
		checksumBase64 = "{sha256}" + checksumBase64
	}

	return checksumBase64, nil
}

// FromFileToSHA512ChecksumBase64 calculates a base64-encoded SHA-512 checksum of the file at the given filepath, optionally prepending a format label.
func FromFileToSHA512ChecksumBase64(filepath string, prependFormat bool) (string, error) {
	// Use the existing function to calculate the SHA-512 checksum.
	checksum, err := FromFileToSHA512Checksum(filepath)
	if err != nil {
		return "", err
	}

	// Use the new function to encode the checksum to Base64.
	checksumBase64 := EncodeToBase64(checksum)

	// Prepend the format label if required.
	if prependFormat {
		checksumBase64 = "{sha512}" + checksumBase64
	}

	return checksumBase64, nil
}

// EncodeToBase64 encodes the given bytes to a base64-encoded string.
func EncodeToBase64(b []byte) string {
	return base64.RawStdEncoding.EncodeToString(b)
}

// Sum256ToSlice computes SHA-256 hash as a slice (helper for array-to-slice conversion).
func Sum256ToSlice(b []byte) []byte {
	hash := sha256.Sum256(b)
	return hash[:]
}

// Sum512ToSlice computes SHA-512 hash as a slice (helper for array-to-slice conversion).
func Sum512ToSlice(b []byte) []byte {
	hash := sha512.Sum512(b)
	return hash[:]
}

// ToHexSHA256 generates a hex-encoded SHA-256 hash (using helper).
func ToHexSHA256(b []byte) string {
	return hex.EncodeToString(Sum256ToSlice(b))
}

// ToHexSHA512 generates a hex-encoded SHA-512 hash (using helper).
func ToHexSHA512(b []byte) string {
	return hex.EncodeToString(Sum512ToSlice(b))
}

// ToBase64SHA256 generates a base64-encoded SHA-256 hash (using helper).
func ToBase64SHA256(b []byte) string {
	return base64.RawStdEncoding.EncodeToString(Sum256ToSlice(b))
}

// ToBase64SHA512 generates a base64-encoded SHA-512 hash (using helper).
func ToBase64SHA512(b []byte) string {
	return base64.RawStdEncoding.EncodeToString(Sum512ToSlice(b))
}
