package asessions

import (
	"fmt"
	"github.com/jpfluger/alibs-slim/auser"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewLoginSecurityFromJSON(t *testing.T) {
	jsonString := `{"username":["user1","user2"],"hasPassword":true,"formDate":"2024-08-22T15:04:05Z"}`
	expectedTime, _ := time.Parse(time.RFC3339, "2024-08-22T15:04:05Z")

	lsec, err := NewLoginSecurityFromJSON(jsonString)
	assert.NoError(t, err)
	assert.NotNil(t, lsec)
	assert.Equal(t, auser.Usernames{"user1", "user2"}, lsec.Usernames)
	assert.True(t, lsec.HasPassword)
	assert.Equal(t, &expectedTime, lsec.FormDate)
}

func TestNewLoginSecurity(t *testing.T) {
	lsec := NewLoginSecurity()
	assert.NotNil(t, lsec)
	assert.Empty(t, lsec.Usernames)
}

func TestAddUsername(t *testing.T) {
	lsec := NewLoginSecurity()
	lsec.AddUsername("newuser")
	assert.Contains(t, lsec.Usernames, auser.Username("newuser"))
}

func TestIsExcessiveUsernamesAttempts(t *testing.T) {
	lsec := NewLoginSecurity()
	for i := 0; i < 10; i++ {
		lsec.AddUsername(auser.Username("user" + fmt.Sprintf("%d", i)))
	}
	err := lsec.IsExcessiveUsernamesAttempts(10)
	assert.NoError(t, err)

	lsec.AddUsername("oneTooMany")
	err = lsec.IsExcessiveUsernamesAttempts(10)
	assert.Error(t, err)
}

func TestIsUsernameValid(t *testing.T) {
	lsec := NewLoginSecurity()
	err := lsec.IsUsernameValid("validUsername")
	assert.NoError(t, err)

	err = lsec.IsUsernameValid("")
	assert.Error(t, err)
}

func TestIsUsernameValidWithOptions(t *testing.T) {
	lsec := NewLoginSecurity()
	err := lsec.IsUsernameValidWithOptions("validUsername", auser.USERNAMEVALIDITYTYPE_EMAIL_OR_USER, nil)
	assert.NoError(t, err)

	err = lsec.IsUsernameValidWithOptions("", auser.USERNAMEVALIDITYTYPE_EMAIL_OR_USER, nil)
	assert.Error(t, err)
}
