package asessions

import (
	"github.com/jpfluger/alibs-slim/auser"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/jpfluger/alibs-slim/atags"
	"github.com/jpfluger/alibs-slim/azb"
)

// TestUiLogin_NewUILoginUser checks the creation of a UILoginUser with default options.
func TestUiLogin_NewUILoginUser(t *testing.T) {
	username := auser.Username("testuser")
	urlPost := "http://example.com/login"
	loginType := LOGINTYPE_SIMPLEAUTH

	user := NewUILoginUser(loginType, username, urlPost)

	assert.Equal(t, loginType, user.GetLoginType())
	assert.Equal(t, username, user.GetUsername())
	assert.Equal(t, urlPost, user.GetUrlPost())
	assert.WithinDuration(t, time.Now().UTC(), user.GetFormDate(), time.Second)
	assert.False(t, user.GetIsOnRememberMe())
	assert.Empty(t, user.GetTitle())
	assert.Empty(t, user.GetLead())
	assert.Empty(t, user.GetFooter())
	assert.Empty(t, user.GetActiveTypes())
	assert.Empty(t, user.GetTags())
}

// TestUiLogin_NewUILoginUserWithOptions checks the creation of a UILoginUser with specified options.
func TestUiLogin_NewUILoginUserWithOptions(t *testing.T) {
	username := auser.Username("testuser")
	urlPost := "http://example.com/login"
	loginType := LOGINTYPE_SIMPLEAUTH
	title := "Login"
	lead := "Please enter your credentials"
	footer := "Footer text"
	formDate := time.Now().UTC()
	activeTypes := azb.ZBTypes{LOGINTYPE_SIMPLEAUTH}
	tags := atags.TagMapString{"key": "value"}
	isOnRememberMe := true

	user := NewUILoginUserWithOptions(loginType, username, urlPost, title, lead, footer, formDate, activeTypes, isOnRememberMe, tags)

	assert.Equal(t, loginType, user.GetLoginType())
	assert.Equal(t, username, user.GetUsername())
	assert.Equal(t, urlPost, user.GetUrlPost())
	assert.Equal(t, formDate, user.GetFormDate())
	assert.Equal(t, title, user.GetTitle())
	assert.Equal(t, lead, user.GetLead())
	assert.Equal(t, footer, user.GetFooter())
	assert.Equal(t, activeTypes, user.GetActiveTypes())
	assert.Equal(t, isOnRememberMe, user.GetIsOnRememberMe())
	assert.Equal(t, tags, user.GetTags())
}

// TestSetFormDate checks if the form date is set correctly.
func TestUiLogin_SetFormDate(t *testing.T) {
	user := &UILoginUser{}
	newFormDate := time.Date(2024, time.April, 23, 0, 0, 0, 0, time.UTC)
	user.SetFormDate(newFormDate)

	assert.Equal(t, newFormDate, user.GetFormDate())
}
