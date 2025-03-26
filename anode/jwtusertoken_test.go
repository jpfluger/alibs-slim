package anode

import (
	"github.com/jpfluger/alibs-slim/auser"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

var secretKey = []byte("my_secret_key")

func TestGenerateTokensWithUsername(t *testing.T) {
	jt := &JWTUserToken{}
	username := auser.Username("testuser")
	accessExpiresInHours := 1
	refreshExpiresInDays := 7

	accessToken, refreshToken, err := jt.GenerateTokensWithUsername(username, accessExpiresInHours, refreshExpiresInDays, secretKey)
	assert.NoError(t, err)
	assert.NotEmpty(t, accessToken)
	assert.NotEmpty(t, refreshToken)
	assert.True(t, jt.isTokenValid)
	assert.WithinDuration(t, time.Now().Add(time.Duration(accessExpiresInHours)*time.Hour), jt.tokenExpires, time.Minute)
	assert.WithinDuration(t, time.Now().Add(time.Duration(refreshExpiresInDays)*24*time.Hour), jt.refreshExpires, time.Minute)
}

func TestVerifyTokenWithSecretKey(t *testing.T) {
	jt := &JWTUserToken{}
	username := auser.Username("testuser")
	accessExpiresInHours := 1
	refreshExpiresInDays := 7

	accessToken, _, err := jt.GenerateTokensWithUsername(username, accessExpiresInHours, refreshExpiresInDays, secretKey)
	assert.NoError(t, err)
	assert.NotEmpty(t, accessToken)

	isValid, err := jt.VerifyTokenWithSecretKey(secretKey)
	assert.NoError(t, err)
	assert.True(t, isValid)

	// Test with invalid secret key
	isValid, err = jt.VerifyTokenWithSecretKey([]byte("invalid_secret_key"))
	assert.Error(t, err)
	assert.False(t, isValid)
}

func TestRefreshAccessTokenWithSecretKey(t *testing.T) {
	jt := &JWTUserToken{}
	username := auser.Username("testuser")
	accessExpiresInHours := 1
	refreshExpiresInDays := 7

	_, refreshToken, err := jt.GenerateTokensWithUsername(username, accessExpiresInHours, refreshExpiresInDays, secretKey)
	assert.NoError(t, err)
	assert.NotEmpty(t, refreshToken)

	newAccessToken, err := jt.RefreshAccessTokenWithSecretKey(username, secretKey)
	assert.NoError(t, err)
	assert.NotEmpty(t, newAccessToken)
	assert.True(t, jt.isTokenValid)
	assert.WithinDuration(t, time.Now().Add(1*time.Hour), jt.tokenExpires, time.Minute)

	// Test with invalid refresh token
	jt.RefreshToken = "invalid_refresh_token"
	newAccessToken, err = jt.RefreshAccessTokenWithSecretKey(username, secretKey)
	assert.Error(t, err)
	assert.Empty(t, newAccessToken)
}

func TestVerifyTokenTimeWithSecretKey(t *testing.T) {
	jt := &JWTUserToken{}
	username := auser.Username("testuser")
	accessExpiresInHours := 1
	refreshExpiresInDays := 7

	accessToken, _, err := jt.GenerateTokensWithUsername(username, accessExpiresInHours, refreshExpiresInDays, secretKey)
	assert.NoError(t, err)

	isValid, expTime, err := jt.VerifyTokenTimeWithSecretKey(accessToken, secretKey)
	assert.NoError(t, err)
	assert.True(t, isValid)
	assert.WithinDuration(t, jt.tokenExpires, *expTime, time.Minute)

	// Test with expired token
	expiredToken := jwt.NewWithClaims(jwt.SigningMethodHS256, NewClaims(username, time.Now().Add(-1*time.Hour)))
	expiredTokenString, _ := expiredToken.SignedString(secretKey)
	isValid, expTime, err = jt.VerifyTokenTimeWithSecretKey(expiredTokenString, secretKey)
	assert.Error(t, err)
	assert.False(t, isValid)
	assert.Nil(t, expTime)
}

func TestValidateTokenWithSecretKey(t *testing.T) {
	jt := &JWTUserToken{}
	username := auser.Username("testuser")
	accessExpiresInHours := 1
	refreshExpiresInDays := 7

	accessToken, refreshToken, err := jt.GenerateTokensWithUsername(username, accessExpiresInHours, refreshExpiresInDays, secretKey)
	assert.NoError(t, err)

	jt.AccessToken = accessToken
	jt.RefreshToken = refreshToken

	err = jt.ValidateTokenWithSecretKey(secretKey)
	assert.NoError(t, err)
	assert.True(t, jt.isTokenValid)
	assert.WithinDuration(t, time.Now().Add(time.Duration(accessExpiresInHours)*time.Hour), jt.tokenExpires, time.Minute)
	assert.WithinDuration(t, time.Now().Add(time.Duration(refreshExpiresInDays)*24*time.Hour), jt.refreshExpires, time.Minute)
}

func TestClearToken(t *testing.T) {
	jt := &JWTUserToken{}
	username := auser.Username("testuser")
	accessExpiresInHours := 1
	refreshExpiresInDays := 7

	accessToken, refreshToken, err := jt.GenerateTokensWithUsername(username, accessExpiresInHours, refreshExpiresInDays, secretKey)
	assert.NoError(t, err)
	assert.NotEmpty(t, accessToken)
	assert.NotEmpty(t, refreshToken)

	jt.ClearToken()
	assert.Empty(t, jt.AccessToken)
	assert.Empty(t, jt.RefreshToken)
	assert.False(t, jt.isTokenValid)
	assert.True(t, jt.tokenExpires.IsZero())
	assert.True(t, jt.refreshExpires.IsZero())
}
