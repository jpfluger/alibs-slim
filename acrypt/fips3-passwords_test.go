package acrypt

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewPBKDF2Presets(t *testing.T) {
	presets := NewPBKDF2Presets()
	assert.Equal(t, 600000, presets.Iterations)
	assert.Equal(t, 32, presets.KeyLen)
	assert.Equal(t, "sha256", presets.HashFunc)
}

func TestHashPBKDF2_SuccessSHA256(t *testing.T) {
	target := "password123"
	presets := NewPBKDF2Presets()

	hash, err := HashPBKDF2(target, presets)
	assert.NoError(t, err)
	assert.True(t, IsPBKDF2Hash(hash))
	assert.Contains(t, hash, "$pbkdf2-sha256$")
}

func TestHashPBKDF2_SuccessSHA512(t *testing.T) {
	target := "password123"
	presets := NewPBKDF2Presets()
	presets.HashFunc = "sha512"

	hash, err := HashPBKDF2(target, presets)
	assert.NoError(t, err)
	assert.True(t, IsPBKDF2Hash(hash))
	assert.Contains(t, hash, "$pbkdf2-sha512$")
}

func TestHashPBKDF2_EmptyTarget(t *testing.T) {
	_, err := HashPBKDF2("", nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot hash an empty string")
}

func TestHashPBKDF2_InvalidPresets(t *testing.T) {
	presets := NewPBKDF2Presets()
	presets.Iterations = 599999 // Below minimum
	_, err := HashPBKDF2("pass", presets)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "iterations too low")

	presets.Iterations = 600000
	presets.KeyLen = 31 // Below minimum
	_, err = HashPBKDF2("pass", presets)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "key length too short")

	presets.KeyLen = 32
	presets.HashFunc = "invalid"
	_, err = HashPBKDF2("pass", presets)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported hash function")
}

func TestVerifyPBKDF2_Success(t *testing.T) {
	target := "password123"
	presets := NewPBKDF2Presets()
	hash, err := HashPBKDF2(target, presets)
	assert.NoError(t, err)

	match, err := VerifyPBKDF2(hash, target)
	assert.NoError(t, err)
	assert.True(t, match)
}

func TestVerifyPBKDF2_WrongPassword(t *testing.T) {
	target := "password123"
	wrong := "wrongpass"
	presets := NewPBKDF2Presets()
	hash, err := HashPBKDF2(target, presets)
	assert.NoError(t, err)

	match, err := VerifyPBKDF2(hash, wrong)
	assert.NoError(t, err)
	assert.False(t, match)
}

func TestVerifyPBKDF2_EmptyParams(t *testing.T) {
	match, err := VerifyPBKDF2("", "pass")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "empty parameters")
	assert.False(t, match)

	match, err = VerifyPBKDF2("$pbkdf2-sha256$i=600000$abc$def", "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "empty parameters")
	assert.False(t, match)
}

func TestVerifyPBKDF2_InvalidFormat(t *testing.T) {
	match, err := VerifyPBKDF2("invalid", "pass")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid PBKDF2 hash format")
	assert.False(t, match)

	match, err = VerifyPBKDF2("$pbkdf2-invalid$i=600000$abc$def", "pass")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported hash function")
	assert.False(t, match)
}

func TestIsPBKDF2Hash(t *testing.T) {
	validSHA256 := "$pbkdf2-sha256$i=600000$MTIzNDU2Nzg5MDEyMzQ1Ng==$MTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTI="
	validSHA512 := "$pbkdf2-sha512$i=600000$MTIzNDU2Nzg5MDEyMzQ1Ng==$MTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTI="
	invalid := "$argon2id$v=19$m=19456,t=2,p=1$MTIzNDU2Nzg5MDEyMzQ1Ng==$MTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTI="

	assert.True(t, IsPBKDF2Hash(validSHA256))
	assert.True(t, IsPBKDF2Hash(validSHA512))
	assert.False(t, IsPBKDF2Hash(invalid))
	assert.False(t, IsPBKDF2Hash(""))
	assert.False(t, IsPBKDF2Hash("$pbkdf2-sha256$invalid"))
}

func TestHashAndVerifyWithCustomPresets(t *testing.T) {
	target := "securepass"
	presets := &PBKDF2Presets{
		Iterations: 700000,
		KeyLen:     64,
		HashFunc:   "sha512",
	}

	hash, err := HashPBKDF2(target, presets) // Use PBKDF2 directly for this test
	assert.NoError(t, err)

	match, err := VerifyPBKDF2(hash, target) // Use PBKDF2 verifier
	assert.NoError(t, err)
	assert.True(t, match)
}

func TestHashPassword_NonFIPS_Argon2(t *testing.T) {
	// Mock non-FIPS mode
	oldIsFIPSMode := IsFIPSMode
	IsFIPSMode = func() bool { return false }
	defer func() { IsFIPSMode = oldIsFIPSMode }()

	password := "testpass"
	hash, err := HashPassword(password, nil) // Default presets
	assert.NoError(t, err)
	assert.True(t, IsArgon2idHash(hash))

	match, err := MatchPassword(hash, password)
	assert.NoError(t, err)
	assert.True(t, match)

	// Wrong password
	match, err = MatchPassword(hash, "wrong")
	assert.NoError(t, err)
	assert.False(t, match)
}

func TestHashPassword_FIPS_PBKDF2(t *testing.T) {
	// Mock FIPS mode
	oldIsFIPSMode := IsFIPSMode
	IsFIPSMode = func() bool { return true }
	defer func() { IsFIPSMode = oldIsFIPSMode }()

	password := "testpass"
	hash, err := HashPassword(password, nil) // Default presets
	assert.NoError(t, err)
	assert.True(t, IsPBKDF2Hash(hash))

	match, err := MatchPassword(hash, password)
	assert.NoError(t, err)
	assert.True(t, match)

	// Wrong password
	match, err = MatchPassword(hash, "wrong")
	assert.NoError(t, err)
	assert.False(t, match)
}

func TestHashPassword_CustomPresets_NonFIPS(t *testing.T) {
	// Mock non-FIPS
	oldIsFIPSMode := IsFIPSMode
	IsFIPSMode = func() bool { return false }
	defer func() { IsFIPSMode = oldIsFIPSMode }()

	password := "secure"
	presets := &Argon2Presets{
		Time:    3,
		Memory:  32768,
		Threads: 2,
		KeyLen:  64,
	}

	hash, err := HashPassword(password, presets)
	assert.NoError(t, err)
	assert.True(t, IsArgon2idHash(hash))

	match, err := MatchPassword(hash, password)
	assert.NoError(t, err)
	assert.True(t, match)
}

func TestHashPassword_CustomPresets_FIPS(t *testing.T) {
	// Mock FIPS
	oldIsFIPSMode := IsFIPSMode
	IsFIPSMode = func() bool { return true }
	defer func() { IsFIPSMode = oldIsFIPSMode }()

	password := "secure"
	presets := &PBKDF2Presets{
		Iterations: 700000,
		KeyLen:     64,
		HashFunc:   "sha512",
	}

	hash, err := HashPassword(password, presets)
	assert.NoError(t, err)
	assert.True(t, IsPBKDF2Hash(hash))

	match, err := MatchPassword(hash, password)
	assert.NoError(t, err)
	assert.True(t, match)
}

func TestHashPassword_InvalidPresets(t *testing.T) {
	// Non-FIPS with wrong presets type
	oldIsFIPSMode := IsFIPSMode
	IsFIPSMode = func() bool { return false }
	defer func() { IsFIPSMode = oldIsFIPSMode }()

	_, err := HashPassword("pass", &PBKDF2Presets{}) // Wrong type for non-FIPS
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid presets")

	// FIPS with wrong presets type
	IsFIPSMode = func() bool { return true }
	_, err = HashPassword("pass", &Argon2Presets{}) // Wrong type for FIPS
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid presets for FIPS mode")
}

func TestHashPassword_EmptyPassword(t *testing.T) {
	_, err := HashPassword("", nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot hash an empty password")
}

func TestMatchPassword_EmptyParams(t *testing.T) {
	match, err := MatchPassword("", "pass")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "empty parameters")
	assert.False(t, match)

	match, err = MatchPassword("hash", "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "empty parameters")
	assert.False(t, match)
}

func TestMatchPassword_UnsupportedFormat(t *testing.T) {
	match, err := MatchPassword("invalid", "pass")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported hash format")
	assert.False(t, match)
}

func TestMatchPassword_CrossAlgorithm(t *testing.T) {
	// Generate Argon2 hash (non-FIPS mock)
	oldIsFIPSMode := IsFIPSMode
	IsFIPSMode = func() bool { return false }
	defer func() { IsFIPSMode = oldIsFIPSMode }()
	hashArgon, err := HashPassword("pass", nil)
	assert.NoError(t, err)

	// Try matching with PBKDF2 verifier (should fail detection)
	match, err := VerifyPBKDF2(hashArgon, "pass")
	assert.Error(t, err) // Invalid format for PBKDF2
	assert.False(t, match)

	// Generate PBKDF2 hash (FIPS mock)
	IsFIPSMode = func() bool { return true }
	hashPBKDF2, err := HashPassword("pass", nil)
	assert.NoError(t, err)

	// Try matching with Argon2 verifier (should fail detection)
	match, err = VerifyArgon2id(hashPBKDF2, "pass")
	assert.Error(t, err) // Invalid format for Argon2
	assert.False(t, match)
}
