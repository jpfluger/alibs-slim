package anode

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/jpfluger/alibs-slim/acrypt"
	"github.com/jpfluger/alibs-slim/anetwork"
)

func TestMatchUsernamePassword(t *testing.T) {
	uc := &UserCredential{
		Username: "testuser",
	}
	err := uc.EncryptPassword("password123")
	assert.NoError(t, err)

	assert.True(t, uc.MatchUsernamePassword("testuser", "password123"))
	assert.False(t, uc.MatchUsernamePassword("testuser", "wrongpassword"))
	assert.False(t, uc.MatchUsernamePassword("wronguser", "password123"))

	// By default, same results as before, when user match is "case-insensitive"
	uc.Username = "TestUser"
	assert.True(t, uc.MatchUsernamePassword("testuser", "password123"))
	assert.False(t, uc.MatchUsernamePassword("testuser", "wrongpassword"))
	assert.False(t, uc.MatchUsernamePassword("wronguser", "password123"))

	// All are not false because user match is "case-sensitive"
	assert.False(t, uc.MatchUsernamePasswordWithCaseSensitive("testuser", "password123", true))
	assert.False(t, uc.MatchUsernamePasswordWithCaseSensitive("testuser", "wrongpassword", true))
	assert.False(t, uc.MatchUsernamePasswordWithCaseSensitive("wronguser", "password123", true))
}

func TestCheckAuthorizationHeaderWithSecretKey(t *testing.T) {
	uc := &UserCredential{
		Username: "testuser",
	}
	uc.EncryptPassword("password123")

	// Basic Auth
	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte("testuser:password123"))
	ok, err := uc.CheckAuthorizationHeaderWithSecretKey(authHeader, secretKey)
	assert.NoError(t, err)
	assert.True(t, ok)

	// Bearer Token Auth
	accessToken, _, _ := uc.GenerateTokensWithUsername("testuser", 1, 7, secretKey)
	authHeader = "Bearer " + accessToken
	ok, err = uc.CheckAuthorizationHeaderWithSecretKey(authHeader, secretKey)
	assert.NoError(t, err)
	assert.True(t, ok)

	// Invalid Auth
	authHeader = "Invalid " + base64.StdEncoding.EncodeToString([]byte("testuser:password123"))
	ok, err = uc.CheckAuthorizationHeaderWithSecretKey(authHeader, secretKey)
	assert.Error(t, err)
	assert.False(t, ok)
}

func TestCheckClientIP(t *testing.T) {
	uc := &UserCredential{
		WhitelistIPs: []string{"192.168.1.0/24"},
	}
	netIPs, _ := anetwork.ToNetIP(uc.WhitelistIPs)
	uc.parsedIPs = netIPs

	assert.True(t, uc.CheckClientIP("192.168.1.100"))
	assert.False(t, uc.CheckClientIP("10.0.0.1"))
}

func TestEncryptPassword(t *testing.T) {
	uc := &UserCredential{}
	err := uc.EncryptPassword("password123")
	assert.NoError(t, err)
	assert.NotEmpty(t, uc.Password)
	assert.True(t, acrypt.IsArgon2idHash(uc.Password))
}

func TestCheckPassword(t *testing.T) {
	uc := &UserCredential{}
	uc.EncryptPassword("password123")

	isEqual, err := uc.CheckPassword("password123")
	assert.NoError(t, err)
	assert.True(t, isEqual)

	isEqual, err = uc.CheckPassword("wrongpassword")
	assert.NoError(t, err)
	assert.False(t, isEqual)
}

func TestValidateWithSecretKey(t *testing.T) {
	uc := &UserCredential{
		Username:     "testuser",
		Password:     "password123",
		WhitelistIPs: []string{"192.168.1.0/24"},
	}
	uc.EncryptPassword("password123")

	err := uc.ValidateWithSecretKey(secretKey)
	assert.NoError(t, err)
	assert.NotEmpty(t, uc.parsedIPs)
	assert.True(t, acrypt.IsArgon2idHash(uc.Password))
}

func TestWhitelistIPContains(t *testing.T) {
	uc := &UserCredential{
		WhitelistIPs: []string{"192.168.1.0/24"},
	}
	netIPs, _ := anetwork.ToNetIP(uc.WhitelistIPs)
	uc.parsedIPs = netIPs

	hasParsedIPs, containsIP := uc.WhitelistIPContains("192.168.1.100")
	assert.True(t, hasParsedIPs)
	assert.True(t, containsIP)

	hasParsedIPs, containsIP = uc.WhitelistIPContains("10.0.0.1")
	assert.True(t, hasParsedIPs)
	assert.False(t, containsIP)
}

func TestGetParsedWhitelistIPs(t *testing.T) {
	uc := &UserCredential{
		WhitelistIPs: []string{"192.168.1.0/24"},
	}
	netIPs, _ := anetwork.ToNetIP(uc.WhitelistIPs)
	uc.parsedIPs = netIPs

	parsedIPs := uc.GetParsedWhitelistIPs()
	assert.Equal(t, netIPs, parsedIPs)
}
