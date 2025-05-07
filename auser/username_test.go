package auser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestUsernameIsEmpty checks the IsEmpty method for Username.
func TestUsernameIsEmpty(t *testing.T) {
	assert.True(t, Username("").IsEmpty())
	assert.True(t, Username("   ").IsEmpty())
	assert.False(t, Username("user").IsEmpty())
}

// TestUsernameTrimSpace checks the TrimSpace method for Username.
func TestUsernameTrimSpace(t *testing.T) {
	assert.Equal(t, Username("user"), Username("  user  ").TrimSpace())
}

// TestUsernameToStringTrimLower checks the ToStringTrimLower method for Username.
func TestUsernameToStringTrimLower(t *testing.T) {
	assert.Equal(t, "user", Username("  USER  ").ToStringTrimLower())
}

// TestUsernameIsEmail checks the IsEmail method for Username.
func TestUsernameIsEmail(t *testing.T) {
	assert.True(t, Username("user@example.com").IsEmail())
	assert.False(t, Username("user").IsEmail())
}

// TestUsernameExpectEmail checks the ExpectEmail method for Username.
func TestUsernameExpectEmail(t *testing.T) {
	assert.True(t, Username("user@example.com").ExpectEmail())
	assert.False(t, Username("user").ExpectEmail())
}

// TestUsernameName checks the Name method for Username.
func TestUsernameName(t *testing.T) {
	assert.Equal(t, "user", Username("user@example.com").Name())
	assert.Equal(t, "user", Username("user").Name())
}

// TestUsernameDomain checks the Domain method for Username.
func TestUsernameDomain(t *testing.T) {
	assert.Equal(t, "example.com", Username("user@example.com").Domain())
	assert.Equal(t, "", Username("user").Domain())
}

// TestUsernameIsValid checks the IsValid method for Username.
func TestUsernameIsValid(t *testing.T) {
	assert.True(t, Username("user@example.com").IsValid(USERNAMEVALIDITYTYPE_EMAIL, nil))
	assert.False(t, Username("user").IsValid(USERNAMEVALIDITYTYPE_EMAIL, nil))
	assert.True(t, Username("user").IsValid(USERNAMEVALIDITYTYPE_USER, nil))
}

// TestUsernameIsValidElseError checks the IsValidElseError method for Username.
func TestUsernameIsValidElseError(t *testing.T) {
	assert.Nil(t, Username("user@example.com").IsValidElseError(USERNAMEVALIDITYTYPE_EMAIL, nil))
	assert.NotNil(t, Username("user").IsValidElseError(USERNAMEVALIDITYTYPE_EMAIL, nil))
	assert.Nil(t, Username("user").IsValidElseError(USERNAMEVALIDITYTYPE_USER, nil))
}

// TestValidateUsername checks the ValidateUsername function.
func TestValidateUsername(t *testing.T) {
	assert.Nil(t, ValidateUsername("valid-username"))
	assert.NotNil(t, ValidateUsername("invalid username"))
	assert.NotNil(t, ValidateUsername("user--name"))
	assert.NotNil(t, ValidateUsername("-username"))
	assert.NotNil(t, ValidateUsername("username-"))
}
