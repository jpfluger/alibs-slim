package acrypt

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"github.com/jpfluger/alibs-slim/autils"
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"strings"
	"testing"
)

func TestFromStringToSHA256CheckSum(t *testing.T) {
	key := "test-key"
	expectedLength := 32 // SHA-256 hashes are 256 bits, i.e., 32 bytes.
	hash := FromStringToSHA256CheckSum(key)
	if len(hash) != expectedLength {
		t.Errorf("Expected SHA-256 hash length of %d, got %d", expectedLength, len(hash))
	}
}

func TestFromStringToSHA512CheckSum(t *testing.T) {
	key := "test-key"
	expectedLength := 64 // SHA-512 hashes are 512 bits, i.e., 64 bytes.
	hash := FromStringToSHA512CheckSum(key)
	if len(hash) != expectedLength {
		t.Errorf("Expected SHA-512 hash length of %d, got %d", expectedLength, len(hash))
	}
}

func TestFromFileToSHA256Checksum(t *testing.T) {
	dir, file1, err := initCryptFileTests(t)
	if err != nil {
		t.Fatalf("Failed to initialize test: %v", err)
	}
	defer deleteDir(t, []string{dir})

	// The expected checksum should be the known hexadecimal string of the file's checksum.
	expectedChecksumHex := "ed489a6bee50871f3aa7e10ac35612ce49a7178fc780dba0225fd4c132252ff6"

	// Calculate the actual checksum of the file.
	checksum, err := FromFileToSHA256Checksum(file1)
	if err != nil {
		t.Fatalf("Failed to create SHA-256 file checksum: %v", err)
	}

	// Convert the checksum to a hexadecimal string.
	actualChecksumHex := FormatSHA256ChecksumHex(checksum)

	// Compare the actual checksum with the expected checksum.
	if actualChecksumHex != expectedChecksumHex {
		t.Errorf("Checksum does not match expected value. Got: %s, Want: %s", actualChecksumHex, expectedChecksumHex)
	}
}

func TestFromFileToSHA512Checksum(t *testing.T) {
	dir, file1, err := initCryptFileTests(t)
	if err != nil {
		t.Fatalf("Failed to initialize test: %v", err)
	}
	defer deleteDir(t, []string{dir})

	// The expected checksum should be the known hexadecimal string of the file's checksum.
	expectedChecksumHex := "683e84c66a1bb6b660a3c2ec274242965020c05715c129c7e86cb3ff0c25ff898c9e257786bee0b462b36916756c1b8dfd78f50c621427ce16d7474b1d5fd976"

	// Calculate the actual checksum of the file.
	checksum, err := FromFileToSHA512Checksum(file1)
	if err != nil {
		t.Fatalf("Failed to create SHA-512 file checksum: %v", err)
	}

	// Convert the checksum to a hexadecimal string.
	actualChecksumHex := FormatSHA512ChecksumHex(checksum)

	// Compare the actual checksum with the expected checksum.
	if actualChecksumHex != expectedChecksumHex {
		t.Errorf("Checksum does not match expected value. Got: %s, Want: %s", actualChecksumHex, expectedChecksumHex)
	}
}

func TestPrependCryptFormatBase64(t *testing.T) {
	format := "{md5}"
	data := []byte("test-data")
	expectedResult := format + base64.RawStdEncoding.EncodeToString(data)
	result := prependCryptFormatBase64(format, data)
	if result != expectedResult {
		t.Errorf("Expected %s, got %s", expectedResult, result)
	}
}

func TestToHashSHA256Base64(t *testing.T) {
	key := "test-key"
	hash := ToHashSHA256Base64(key, true)
	if !strings.HasPrefix(hash, "{sha256}") {
		t.Errorf("Expected hash to be prefixed with {sha256}")
	}
}

func TestToHashSHA512Base64(t *testing.T) {
	key := "test-key"
	hash := ToHashSHA512Base64(key, true)
	if !strings.HasPrefix(hash, "{sha512}") {
		t.Errorf("Expected hash to be prefixed with {sha512}")
	}
}

func initCryptFileTests(t *testing.T) (dir string, file1 string, err error) {
	dir, err = autils.CreateTempDir()
	if err != nil {
		t.Fatalf("cannot create temp directory; %v", err)
	}

	file1 = path.Join(dir, "test.txt")
	if err := os.WriteFile(file1, []byte("hello crypto"), autils.PATH_CHMOD_FILE); err != nil {
		deleteDir(t, []string{dir})
		t.Fatalf("cannot create test file; %v", err)
	}

	return dir, file1, nil
}

func deleteDir(t *testing.T, dir []string) {
	for _, d := range dir {
		if err := os.RemoveAll(d); err != nil {
			t.Fatalf("failed to remove test directory at %s; %v", d, err)
		}
	}
}

// TestFromBytesToSHA256Checksum tests the FromBytesToSHA256Checksum function.
func TestFromBytesToSHA256Checksum(t *testing.T) {
	// Define a byte slice to hash.
	data := []byte("hello world")

	// Generate the SHA-256 checksum using the standard library for comparison.
	expectedChecksumArray := sha256.Sum256(data)

	// Generate the SHA-256 checksum using our function.
	actualChecksumSlice := FromBytesToSHA256Checksum(data)

	// Convert the expected checksum array to a slice for comparison.
	expectedChecksumSlice := expectedChecksumArray[:]

	// Compare the checksums.
	assert.Equal(t, expectedChecksumSlice, actualChecksumSlice)
}

// TestFromBytesToSHA512Checksum tests the FromBytesToSHA512Checksum function.
func TestFromBytesToSHA512Checksum(t *testing.T) {
	// Define a byte slice to hash.
	data := []byte("hello world")

	// Generate the SHA-512 checksum using the standard library for comparison.
	expectedChecksumArray := sha512.Sum512(data)

	// Generate the SHA-512 checksum using our function.
	actualChecksumSlice := FromBytesToSHA512Checksum(data)

	// Convert the expected checksum array to a slice for comparison.
	expectedChecksumSlice := expectedChecksumArray[:]

	// Compare the checksums.
	assert.Equal(t, expectedChecksumSlice, actualChecksumSlice)
}

// TestFromFileToSHA256ChecksumBase64 tests the FromFileToSHA256ChecksumBase64 function.
func TestFromFileToSHA256ChecksumBase64(t *testing.T) {
	dir, file1, err := initCryptFileTests(t)
	if err != nil {
		t.Fatalf("Failed to initialize test: %v", err)
	}
	defer deleteDir(t, []string{dir})

	// Create a temporary file with known content.
	content := []byte("hello crypto")

	// Calculate the expected SHA-256 checksum using the standard library.
	expectedChecksum := sha256.Sum256(content)
	expectedChecksumBase64 := base64.RawStdEncoding.EncodeToString(expectedChecksum[:])

	// Test without prepending format label.
	checksumBase64, err := FromFileToSHA256ChecksumBase64(file1, false)
	if err != nil {
		t.Errorf("Failed to create SHA-256 file checksum base64: %v", err)
	}
	if checksumBase64 != expectedChecksumBase64 {
		t.Errorf("Checksum base64 does not match expected value. Got: %s, Want: %s", checksumBase64, expectedChecksumBase64)
	}

	// Test with prepending format label.
	checksumBase64, err = FromFileToSHA256ChecksumBase64(file1, true)
	if err != nil {
		t.Errorf("Failed to create SHA-256 file checksum base64 with format label: %v", err)
	}
	if checksumBase64 != "{sha256}"+expectedChecksumBase64 {
		t.Errorf("Checksum base64 with format label does not match expected value. Got: %s, Want: %s", checksumBase64, "{sha256}"+expectedChecksumBase64)
	}
}

// TestFromFileToSHA512ChecksumBase64 tests the FromFileToSHA512ChecksumBase64 function.
func TestFromFileToSHA512ChecksumBase64(t *testing.T) {
	dir, file1, err := initCryptFileTests(t)
	if err != nil {
		t.Fatalf("Failed to initialize test: %v", err)
	}
	defer deleteDir(t, []string{dir})

	// Create a temporary file with known content.
	content := []byte("hello crypto")

	// Calculate the expected SHA-512 checksum using the standard library.
	expectedChecksum := sha512.Sum512(content)
	expectedChecksumBase64 := base64.RawStdEncoding.EncodeToString(expectedChecksum[:])

	// Test without prepending format label.
	checksumBase64, err := FromFileToSHA512ChecksumBase64(file1, false)
	if err != nil {
		t.Errorf("Failed to create SHA-512 file checksum base64: %v", err)
	}
	if checksumBase64 != expectedChecksumBase64 {
		t.Errorf("Checksum base64 does not match expected value. Got: %s, Want: %s", checksumBase64, expectedChecksumBase64)
	}

	// Test with prepending format label.
	checksumBase64, err = FromFileToSHA512ChecksumBase64(file1, true)
	if err != nil {
		t.Errorf("Failed to create SHA-512 file checksum base64 with format label: %v", err)
	}
	if checksumBase64 != "{sha512}"+expectedChecksumBase64 {
		t.Errorf("Checksum base64 with format label does not match expected value. Got: %s, Want: %s", checksumBase64, "{sha512}"+expectedChecksumBase64)
	}
}

func TestToCheckSumSHA256(t *testing.T) {
	tests := []struct {
		name      string
		input     []byte
		expectHex string
	}{
		{
			name:      "Empty Input",
			input:     []byte(""),
			expectHex: ToHexSHA256([]byte("")),
		},
		{
			name:      "Short String",
			input:     []byte("hello"),
			expectHex: ToHexSHA256([]byte("hello")),
		},
		{
			name:      "Long String",
			input:     []byte("The quick brown fox jumps over the lazy dog"),
			expectHex: ToHexSHA256([]byte("The quick brown fox jumps over the lazy dog")),
		},
		{
			name:      "Binary Data",
			input:     []byte{0x00, 0xFF, 0xAA, 0xBB, 0xCC},
			expectHex: ToHexSHA256([]byte{0x00, 0xFF, 0xAA, 0xBB, 0xCC}),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := ToCheckSumSHA256(tc.input)
			assert.Equal(t, tc.expectHex, actual)
		})
	}
}

func TestToCheckSumSHA512(t *testing.T) {
	tests := []struct {
		name      string
		input     []byte
		expectHex string
	}{
		{
			name:      "Empty Input",
			input:     []byte(""),
			expectHex: ToHexSHA512([]byte("")),
		},
		{
			name:      "Short String",
			input:     []byte("hello"),
			expectHex: ToHexSHA512([]byte("hello")),
		},
		{
			name:      "Long String",
			input:     []byte("The quick brown fox jumps over the lazy dog"),
			expectHex: ToHexSHA512([]byte("The quick brown fox jumps over the lazy dog")),
		},
		{
			name:      "Binary Data",
			input:     []byte{0x00, 0xFF, 0xAA, 0xBB, 0xCC},
			expectHex: ToHexSHA512([]byte{0x00, 0xFF, 0xAA, 0xBB, 0xCC}),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := ToCheckSumSHA512(tc.input)
			assert.Equal(t, tc.expectHex, actual)
		})
	}
}
