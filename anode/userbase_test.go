package anode

import (
	"encoding/base64"
	"encoding/json"
	"github.com/jpfluger/alibs-slim/auser"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/jpfluger/alibs-slim/aemail"
)

var jwtSecretKey []byte

func init() {
	jwtSecretKey, _ = base64.StdEncoding.DecodeString("EiWboODvRtZ+t3culHdQCnj7skow+M6f+VO6zSEYapA=")
}

// UserConfig represents a user configuration.
// Many various styles of UserConfig can exist, dependent upon the business purpose.
// The struct here is general-purpose that you can copy and adapt.
type UserConfig struct {
	UserBase

	// Account contains protected user details.
	// The Email is found under Account. Email is required for forgot-login functionality.
	Account UserAccount `json:"account"`

	// Vault contains secrets that should not be shared without a higher degree of authorization.
	// The Username, associated with the password/tokens, is found under Vault.Credential.Username.
	// Username is extracted upon new session creation (e.g., login to website).
	Vault UserVault `json:"vault"`

	mu sync.RWMutex // Added mutex for thread-safety
}

// Validate checks the UserConfig for validity.
func (uc *UserConfig) Validate() error {
	if err := uc.Vault.Credential.ValidateWithSecretKey(jwtSecretKey); err != nil {
		return err
	}
	return nil
}

// GetEmail returns the email address from the user account.
func (uc *UserConfig) GetEmail() aemail.EmailAddress {
	return uc.Account.Email
}

// GetUsername returns the username from the user vault.
func (uc *UserConfig) GetUsername() auser.Username {
	return uc.Vault.Credential.Username
}

// RefreshAccessToken refreshes the access token using the secret key.
func (uc *UserConfig) RefreshAccessToken() (string, error) {
	return uc.Vault.Credential.RefreshAccessTokenWithSecretKey(uc.Vault.Credential.GetUsername(), jwtSecretKey)
}

// VerifyToken verifies the token using the secret key.
func (uc *UserConfig) VerifyToken() (bool, error) {
	return uc.Vault.Credential.VerifyTokenWithSecretKey(jwtSecretKey)
}

// VerifyTokenTime verifies the token and returns the expiration time.
func (uc *UserConfig) VerifyTokenTime(token string) (bool, *time.Time, error) {
	return uc.Vault.Credential.VerifyTokenTimeWithSecretKey(token, jwtSecretKey)
}

// GenerateTokens generates access and refresh tokens with specified expiration times.
func (uc *UserConfig) GenerateTokens(accessExpiresInHours, refreshExpiresInDays int) (string, string, error) {
	return uc.Vault.Credential.GenerateTokensWithSecretKey(accessExpiresInHours, refreshExpiresInDays, jwtSecretKey)
}

// CheckAuthorizationHeader checks the authorization header using the secret key.
func (uc *UserConfig) CheckAuthorizationHeader(authHeader string) (bool, error) {
	return uc.Vault.Credential.CheckAuthorizationHeaderWithSecretKey(authHeader, jwtSecretKey)
}

var testJSON = `{
  "account": {
    "email": "bob@example.com"
  },
  "vault": {
    "credential": {
      "username": "testuser",
      "password": "$argon2id$v=19$m=65536,t=1,p=4$VX77b4N9XFMLgELAy7ZhRw==$GzZYY3bA86LFmorUUs85up1UydiMYqhEk2D8HlL1mRc=",
      "whitelistIPs": ["192.168.1.1", "10.0.0.1", "127.0.0.1", "::1/128"],
      "accessToken": "",
      "refreshToken": ""
    }
  }
}`

func TestUserConfig_Validate(t *testing.T) {
	var uc UserConfig
	if err := json.Unmarshal([]byte(testJSON), &uc); err != nil {
		t.Error(err)
		return
	}
	err := uc.Validate()
	assert.NoError(t, err)
	assert.False(t, uc.Vault.Credential.GetIsTokenValid())
	assert.Equal(t, 4, len(uc.Vault.Credential.GetWhitelistIPs()))
	assert.True(t, uc.Vault.Credential.MatchUsernamePassword("testuser", "password123"))
}

func TestUserConfig_GetEmail(t *testing.T) {
	var uc UserConfig
	if err := json.Unmarshal([]byte(testJSON), &uc); err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, "bob@example.com", uc.GetEmail().String())
}

func TestUserConfig_GetUsername(t *testing.T) {
	var uc UserConfig
	if err := json.Unmarshal([]byte(testJSON), &uc); err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, "testuser", uc.GetUsername().String())
}

func TestUserConfig_GenerateTokens(t *testing.T) {
	var uc UserConfig
	if err := json.Unmarshal([]byte(testJSON), &uc); err != nil {
		t.Error(err)
		return
	}
	accessToken, refreshToken, err := uc.GenerateTokens(1, 7)
	assert.NoError(t, err)
	assert.NotEmpty(t, accessToken)
	assert.NotEmpty(t, refreshToken)

	valid, err := uc.VerifyToken()
	assert.NoError(t, err)
	assert.True(t, valid)

	// Expire the token
	uc.Vault.Credential.SetTokenExpiration(time.Now().Add(-1 * time.Hour))
	_, _, err = uc.GenerateTokens(-1, 7) // Generate an expired token
	assert.NoError(t, err)

	valid, err = uc.VerifyToken()
	assert.Error(t, err)
	assert.False(t, valid)

	newAccessToken, err := uc.RefreshAccessToken()
	assert.NoError(t, err)
	assert.NotEmpty(t, newAccessToken)

	valid, err = uc.VerifyToken()
	assert.NoError(t, err)
	assert.True(t, valid)
}
