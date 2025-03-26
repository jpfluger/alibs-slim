package asessions

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var mySigningKey = []byte("AllYourBase")

func TestJWT5(t *testing.T) {
	// Create the claims
	claims := JWTClaim{
		Username: "bar",
		Action:   "totp",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "test",
			Subject:   "somebody",
			ID:        "1",
			Audience:  jwt.ClaimStrings{"somebody_else"},
		},
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(mySigningKey)
	assert.NoError(t, err)
	assert.NotEmpty(t, ss)

	// Parse the token
	token2, err := jwt.ParseWithClaims(ss, &JWTClaim{}, func(token *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	})
	assert.NoError(t, err)

	// Assert the claims
	if claims2, ok := token2.Claims.(*JWTClaim); ok && token2.Valid {
		assert.Equal(t, claims.Username, claims2.Username)
		assert.Equal(t, claims.Action, claims2.Action)
		assert.Equal(t, claims.Issuer, claims2.Issuer)
		assert.Equal(t, claims.Subject, claims2.Subject)
		assert.Equal(t, claims.ID, claims2.ID)

		// Verify the time-based claims
		clickTime := time.Now().Add(5 * time.Hour)
		assert.True(t, claims2.ExpiresAt.Time.After(clickTime))
		assert.True(t, claims2.IssuedAt.Time.Before(clickTime))
		assert.True(t, claims2.NotBefore.Time.Before(clickTime))

		clickTime = time.Now().Add(-1 * time.Hour)
		assert.True(t, claims2.ExpiresAt.Time.After(clickTime))
		assert.False(t, claims2.IssuedAt.Time.Before(clickTime))
		assert.False(t, claims2.NotBefore.Time.Before(clickTime))

		clickTime = time.Now().Add(25 * time.Hour)
		assert.False(t, claims2.ExpiresAt.Time.After(clickTime))
		assert.True(t, claims2.IssuedAt.Time.Before(clickTime))
		assert.True(t, claims2.NotBefore.Time.Before(clickTime))
	} else {
		t.Errorf("Token claims are not valid")
	}
}
