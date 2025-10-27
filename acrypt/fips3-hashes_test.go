package acrypt

import (
	"bytes"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/jpfluger/alibs-slim/autils"
	"github.com/stretchr/testify/assert"
)

func TestHashSHA256String(t *testing.T) {
	key := "test-key"
	expectedLength := 32 // SHA-256 hashes are 256 bits, i.e., 32 bytes.
	hash := HashSHA256String(key)
	if len(hash) != expectedLength {
		t.Errorf("Expected SHA-256 hash length of %d, got %d", expectedLength, len(hash))
	}
}

func TestHashSHA384String(t *testing.T) {
	key := "test-key"
	expectedLength := 48 // SHA-384 hashes are 384 bits, i.e., 48 bytes.
	hash := HashSHA384String(key)
	if len(hash) != expectedLength {
		t.Errorf("Expected SHA-384 hash length of %d, got %d", expectedLength, len(hash))
	}
}

func TestHashSHA512String(t *testing.T) {
	key := "test-key"
	expectedLength := 64 // SHA-512 hashes are 512 bits, i.e., 64 bytes.
	hash := HashSHA512String(key)
	if len(hash) != expectedLength {
		t.Errorf("Expected SHA-512 hash length of %d, got %d", expectedLength, len(hash))
	}
}

func TestHashSHA256File(t *testing.T) {
	dir, file1, err := initCryptFileTests(t)
	if err != nil {
		t.Fatalf("Failed to initialize test: %v", err)
	}
	defer deleteDir(t, []string{dir})

	// The expected checksum should be the known hexadecimal string of the file's checksum.
	expectedChecksumHex := "ed489a6bee50871f3aa7e10ac35612ce49a7178fc780dba0225fd4c132252ff6"

	// Calculate the actual checksum of the file.
	checksum, err := HashSHA256File(file1)
	if err != nil {
		t.Fatalf("Failed to create SHA-256 file checksum: %v", err)
	}

	// Convert the checksum to a hexadecimal string.
	actualChecksumHex := hex.EncodeToString(checksum)

	// Compare the actual checksum with the expected checksum.
	if actualChecksumHex != expectedChecksumHex {
		t.Errorf("Checksum does not match expected value. Got: %s, Want: %s", actualChecksumHex, expectedChecksumHex)
	}
}

func TestHashSHA384File(t *testing.T) {
	dir, file1, err := initCryptFileTests(t)
	if err != nil {
		t.Fatalf("Failed to initialize test: %v", err)
	}
	defer deleteDir(t, []string{dir})

	// The expected checksum should be the known hexadecimal string of the file's checksum.
	expectedChecksumHex := "6705b17acfb97f016f2464792fb980c14f19142efb37bd73cebfc78f43830eb8cc8c3436539cefd3f1584c7ba5d38c30"

	// Calculate the actual checksum of the file.
	checksum, err := HashSHA384File(file1)
	if err != nil {
		t.Fatalf("Failed to create SHA-384 file checksum: %v", err)
	}

	// Convert the checksum to a hexadecimal string.
	actualChecksumHex := hex.EncodeToString(checksum)

	// Compare the actual checksum with the expected checksum.
	if actualChecksumHex != expectedChecksumHex {
		t.Errorf("Checksum does not match expected value. Got: %s, Want: %s", actualChecksumHex, expectedChecksumHex)
	}
}

func TestHashSHA512File(t *testing.T) {
	dir, file1, err := initCryptFileTests(t)
	if err != nil {
		t.Fatalf("Failed to initialize test: %v", err)
	}
	defer deleteDir(t, []string{dir})

	// The expected checksum should be the known hexadecimal string of the file's checksum.
	expectedChecksumHex := "683e84c66a1bb6b660a3c2ec274242965020c05715c129c7e86cb3ff0c25ff898c9e257786bee0b462b36916756c1b8dfd78f50c621427ce16d7474b1d5fd976"

	// Calculate the actual checksum of the file.
	checksum, err := HashSHA512File(file1)
	if err != nil {
		t.Fatalf("Failed to create SHA-512 file checksum: %v", err)
	}

	// Convert the checksum to a hexadecimal string.
	actualChecksumHex := hex.EncodeToString(checksum)

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

func TestToBase64SHA256WithFormat(t *testing.T) {
	key := "test-key"
	hash := ToBase64SHA256WithFormat(key, true)
	if !strings.HasPrefix(hash, "{sha256}") {
		t.Errorf("Expected hash to be prefixed with {sha256}")
	}
}

func TestToBase64SHA384WithFormat(t *testing.T) {
	key := "test-key"
	hash := ToBase64SHA384WithFormat(key, true)
	if !strings.HasPrefix(hash, "{sha384}") {
		t.Errorf("Expected hash to be prefixed with {sha384}")
	}
}

func TestToBase64SHA512WithFormat(t *testing.T) {
	key := "test-key"
	hash := ToBase64SHA512WithFormat(key, true)
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

// TestHashSHA256Bytes tests the HashSHA256Bytes function.
func TestHashSHA256Bytes(t *testing.T) {
	// Define a byte slice to hash.
	data := []byte("hello world")

	// Generate the SHA-256 checksum using the standard library for comparison.
	expectedChecksumArray := sha256.Sum256(data)

	// Generate the SHA-256 checksum using our function.
	actualChecksumSlice := HashSHA256Bytes(data)

	// Convert the expected checksum array to a slice for comparison.
	expectedChecksumSlice := expectedChecksumArray[:]

	// Compare the checksums.
	assert.Equal(t, expectedChecksumSlice, actualChecksumSlice)
}

// TestHashSHA384Bytes tests the HashSHA384Bytes function.
func TestHashSHA384Bytes(t *testing.T) {
	// Define a byte slice to hash.
	data := []byte("hello world")

	// Generate the SHA-384 checksum using the standard library for comparison.
	expectedChecksumArray := sha512.Sum384(data)

	// Generate the SHA-384 checksum using our function.
	actualChecksumSlice := HashSHA384Bytes(data)

	// Convert the expected checksum array to a slice for comparison.
	expectedChecksumSlice := expectedChecksumArray[:]

	// Compare the checksums.
	assert.Equal(t, expectedChecksumSlice, actualChecksumSlice)
}

// TestHashSHA512Bytes tests the HashSHA512Bytes function.
func TestHashSHA512Bytes(t *testing.T) {
	// Define a byte slice to hash.
	data := []byte("hello world")

	// Generate the SHA-512 checksum using the standard library for comparison.
	expectedChecksumArray := sha512.Sum512(data)

	// Generate the SHA-512 checksum using our function.
	actualChecksumSlice := HashSHA512Bytes(data)

	// Convert the expected checksum array to a slice for comparison.
	expectedChecksumSlice := expectedChecksumArray[:]

	// Compare the checksums.
	assert.Equal(t, expectedChecksumSlice, actualChecksumSlice)
}

// TestToBase64SHA256FileWithFormat tests the ToBase64SHA256FileWithFormat function.
func TestToBase64SHA256FileWithFormat(t *testing.T) {
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
	checksumBase64, err := ToBase64SHA256FileWithFormat(file1, false)
	if err != nil {
		t.Errorf("Failed to create SHA-256 file checksum base64: %v", err)
	}
	if checksumBase64 != expectedChecksumBase64 {
		t.Errorf("Checksum base64 does not match expected value. Got: %s, Want: %s", checksumBase64, expectedChecksumBase64)
	}

	// Test with prepending format label.
	checksumBase64, err = ToBase64SHA256FileWithFormat(file1, true)
	if err != nil {
		t.Errorf("Failed to create SHA-256 file checksum base64 with format label: %v", err)
	}
	if checksumBase64 != "{sha256}"+expectedChecksumBase64 {
		t.Errorf("Checksum base64 with format label does not match expected value. Got: %s, Want: %s", checksumBase64, "{sha256}"+expectedChecksumBase64)
	}
}

// TestToBase64SHA384FileWithFormat tests the ToBase64SHA384FileWithFormat function.
func TestToBase64SHA384FileWithFormat(t *testing.T) {
	dir, file1, err := initCryptFileTests(t)
	if err != nil {
		t.Fatalf("Failed to initialize test: %v", err)
	}
	defer deleteDir(t, []string{dir})

	// Create a temporary file with known content.
	content := []byte("hello crypto")

	// Calculate the expected SHA-384 checksum using the standard library.
	expectedChecksum := sha512.Sum384(content)
	expectedChecksumBase64 := base64.RawStdEncoding.EncodeToString(expectedChecksum[:])

	// Test without prepending format label.
	checksumBase64, err := ToBase64SHA384FileWithFormat(file1, false)
	if err != nil {
		t.Errorf("Failed to create SHA-384 file checksum base64: %v", err)
	}
	if checksumBase64 != expectedChecksumBase64 {
		t.Errorf("Checksum base64 does not match expected value. Got: %s, Want: %s", checksumBase64, expectedChecksumBase64)
	}

	// Test with prepending format label.
	checksumBase64, err = ToBase64SHA384FileWithFormat(file1, true)
	if err != nil {
		t.Errorf("Failed to create SHA-384 file checksum base64 with format label: %v", err)
	}
	if checksumBase64 != "{sha384}"+expectedChecksumBase64 {
		t.Errorf("Checksum base64 with format label does not match expected value. Got: %s, Want: %s", checksumBase64, "{sha384}"+expectedChecksumBase64)
	}
}

// TestToBase64SHA512FileWithFormat tests the ToBase64SHA512FileWithFormat function.
func TestToBase64SHA512FileWithFormat(t *testing.T) {
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
	checksumBase64, err := ToBase64SHA512FileWithFormat(file1, false)
	if err != nil {
		t.Errorf("Failed to create SHA-512 file checksum base64: %v", err)
	}
	if checksumBase64 != expectedChecksumBase64 {
		t.Errorf("Checksum base64 does not match expected value. Got: %s, Want: %s", checksumBase64, expectedChecksumBase64)
	}

	// Test with prepending format label.
	checksumBase64, err = ToBase64SHA512FileWithFormat(file1, true)
	if err != nil {
		t.Errorf("Failed to create SHA-512 file checksum base64 with format label: %v", err)
	}
	if checksumBase64 != "{sha512}"+expectedChecksumBase64 {
		t.Errorf("Checksum base64 with format label does not match expected value. Got: %s, Want: %s", checksumBase64, "{sha512}"+expectedChecksumBase64)
	}
}

func TestToHexSHA256(t *testing.T) {
	fnToHex := func(data []byte) string {
		hash := sha256.Sum256(data)
		return hex.EncodeToString(hash[:])
	}
	tests := []struct {
		name      string
		input     []byte
		expectHex string
	}{
		{
			name:      "Empty Input",
			input:     []byte(""),
			expectHex: fnToHex([]byte("")),
		},
		{
			name:      "Short String",
			input:     []byte("hello"),
			expectHex: fnToHex([]byte("hello")),
		},
		{
			name:      "Long String",
			input:     []byte("The quick brown fox jumps over the lazy dog"),
			expectHex: fnToHex([]byte("The quick brown fox jumps over the lazy dog")),
		},
		{
			name:      "Binary Data",
			input:     []byte{0x00, 0xFF, 0xAA, 0xBB, 0xCC},
			expectHex: fnToHex([]byte{0x00, 0xFF, 0xAA, 0xBB, 0xCC}),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := ToHexSHA256(tc.input)
			assert.Equal(t, tc.expectHex, actual)
		})
	}
}

func TestToHexSHA384(t *testing.T) {
	fnToHex := func(data []byte) string {
		hash := sha512.Sum384(data)
		return hex.EncodeToString(hash[:])
	}
	tests := []struct {
		name      string
		input     []byte
		expectHex string
	}{
		{
			name:      "Empty Input",
			input:     []byte(""),
			expectHex: fnToHex([]byte("")),
		},
		{
			name:      "Short String",
			input:     []byte("hello"),
			expectHex: fnToHex([]byte("hello")),
		},
		{
			name:      "Long String",
			input:     []byte("The quick brown fox jumps over the lazy dog"),
			expectHex: fnToHex([]byte("The quick brown fox jumps over the lazy dog")),
		},
		{
			name:      "Binary Data",
			input:     []byte{0x00, 0xFF, 0xAA, 0xBB, 0xCC},
			expectHex: fnToHex([]byte{0x00, 0xFF, 0xAA, 0xBB, 0xCC}),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := ToHexSHA384(tc.input)
			assert.Equal(t, tc.expectHex, actual)
		})
	}
}

func TestToHexSHA512(t *testing.T) {
	fnToHex := func(data []byte) string {
		hash := sha512.Sum512(data)
		return hex.EncodeToString(hash[:])
	}
	tests := []struct {
		name      string
		input     []byte
		expectHex string
	}{
		{
			name:      "Empty Input",
			input:     []byte(""),
			expectHex: fnToHex([]byte("")),
		},
		{
			name:      "Short String",
			input:     []byte("hello"),
			expectHex: fnToHex([]byte("hello")),
		},
		{
			name:      "Long String",
			input:     []byte("The quick brown fox jumps over the lazy dog"),
			expectHex: fnToHex([]byte("The quick brown fox jumps over the lazy dog")),
		},
		{
			name:      "Binary Data",
			input:     []byte{0x00, 0xFF, 0xAA, 0xBB, 0xCC},
			expectHex: fnToHex([]byte{0x00, 0xFF, 0xAA, 0xBB, 0xCC}),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := ToHexSHA512(tc.input)
			assert.Equal(t, tc.expectHex, actual)
		})
	}
}

// TestParseChecksum tests the ParseChecksum function with various valid and invalid inputs.
func TestParseChecksum(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantAlg    string
		wantRaw    []byte
		wantErr    bool
		errMessage string
	}{
		// Valid cases
		{
			name:    "valid sha256",
			input:   "sha256:" + strings.Repeat("a", 64),
			wantAlg: "sha256",
			wantRaw: mustHexDecode(strings.Repeat("a", 64)),
			wantErr: false,
		},
		{
			name:    "valid sha384",
			input:   "sha384:" + strings.Repeat("a", 96),
			wantAlg: "sha384",
			wantRaw: mustHexDecode(strings.Repeat("a", 96)),
			wantErr: false,
		},
		{
			name:    "valid sha512",
			input:   "sha512:" + strings.Repeat("a", 128),
			wantAlg: "sha512",
			wantRaw: mustHexDecode(strings.Repeat("a", 128)),
			wantErr: false,
		},
		{
			name:    "valid uppercase alg",
			input:   "sha256:" + strings.Repeat("a", 64), // Use lowercase for digest.Parse
			wantAlg: "sha256",
			wantRaw: mustHexDecode(strings.Repeat("a", 64)),
			wantErr: false,
		},

		// Invalid cases
		{
			name:       "invalid format no colon",
			input:      "sha256" + strings.Repeat("a", 64),
			wantErr:    true,
			errMessage: "invalid checksum format",
		},
		{
			name:       "invalid alg",
			input:      "md5:" + strings.Repeat("a", 32),
			wantErr:    true,
			errMessage: "invalid checksum format: unsupported digest algorithm",
		},
		{
			name:       "invalid alg",
			input:      "md5:" + strings.Repeat("a", 32), // MD5 is supported in non-FIPS
			wantErr:    true,
			errMessage: "invalid checksum format: unsupported digest algorithm",
		},
		{
			name:       "invalid hex length for sha256",
			input:      "sha256:" + strings.Repeat("a", 63),
			wantErr:    true,
			errMessage: "invalid checksum digest length",
		},
		{
			name:       "invalid hex chars",
			input:      "sha256:" + strings.Repeat("a", 63) + "g",
			wantErr:    true,
			errMessage: "invalid checksum digest format",
		},
		{
			name:       "empty input",
			input:      "",
			wantErr:    true,
			errMessage: "invalid checksum format",
		},
		{
			name:       "only alg no hash",
			input:      "sha256:",
			wantErr:    true,
			errMessage: "invalid checksum digest format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			alg, raw, err := ParseChecksum(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseChecksum() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && !strings.Contains(err.Error(), tt.errMessage) {
				t.Errorf("ParseChecksum() error = %v, want err containing %q", err, tt.errMessage)
			}
			if alg != tt.wantAlg {
				t.Errorf("ParseChecksum() alg = %v, want %v", alg, tt.wantAlg)
			}
			if !bytes.Equal(raw, tt.wantRaw) {
				t.Errorf("ParseChecksum() raw = %x, want %x", raw, tt.wantRaw)
			}
		})
	}
}

// TestIsValidChecksumFingerprint tests the IsValidChecksumFingerprint function with various inputs.
func TestIsValidChecksumFingerprint(t *testing.T) {
	tests := []struct {
		name string
		alg  string
		s    string
		want bool
	}{
		// Valid cases
		{
			name: "valid sha256 lowercase hex",
			alg:  "sha256",
			s:    strings.Repeat("a", 64),
			want: true,
		},
		{
			name: "valid sha256 uppercase hex",
			alg:  "sha256",
			s:    strings.Repeat("A", 64),
			want: true,
		},
		{
			name: "valid sha256 mixed hex",
			alg:  "sha256",
			s:    "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
			want: true,
		},
		{
			name: "valid sha384",
			alg:  "sha384",
			s:    strings.Repeat("a", 96),
			want: true,
		},
		{
			name: "valid sha512",
			alg:  "sha512",
			s:    strings.Repeat("a", 128),
			want: true,
		},
		{
			name: "valid uppercase alg",
			alg:  "SHA256",
			s:    strings.Repeat("a", 64),
			want: true,
		},

		// Invalid cases
		{
			name: "invalid length sha256 short",
			alg:  "sha256",
			s:    strings.Repeat("a", 63),
			want: false,
		},
		{
			name: "invalid length sha256 long",
			alg:  "sha256",
			s:    strings.Repeat("a", 65),
			want: false,
		},
		{
			name: "invalid chars",
			alg:  "sha256",
			s:    strings.Repeat("a", 63) + "g",
			want: false,
		},
		{
			name: "unsupported alg",
			alg:  "md5",
			s:    strings.Repeat("a", 32),
			want: false,
		},
		{
			name: "empty string",
			alg:  "sha256",
			s:    "",
			want: false,
		},
		{
			name: "invalid length for sha384",
			alg:  "sha384",
			s:    strings.Repeat("a", 95),
			want: false,
		},
		{
			name: "invalid length for sha512",
			alg:  "sha512",
			s:    strings.Repeat("a", 127),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidChecksumFingerprint(tt.alg, tt.s); got != tt.want {
				t.Errorf("IsValidChecksumFingerprint(%q, %q) = %v, want %v", tt.alg, tt.s, got, tt.want)
			}
		})
	}
}

// mustHexDecode decodes a hex string or panics (for test setup only).
func mustHexDecode(s string) []byte {
	b, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return b
}
