package asessions

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/jpfluger/alibs-slim/aimage"
	"github.com/jpfluger/alibs-slim/atags"
)

// TestNewUserSessionBase checks the creation of a new UserSessionBase with default values.
func TestNewUserSessionBase(t *testing.T) {
	session := NewUserSessionBase()
	assert.NotNil(t, session.SID)
	assert.Empty(t, session.Actions)
	assert.Empty(t, session.Options)
}

// TestNewUserSessionBaseWithLoginStatus checks the creation of a new UserSessionBase with a specified login status.
func TestNewUserSessionBaseWithLoginStatus(t *testing.T) {
	statusType := LOGIN_SESSION_STATUS_OK
	session := NewUserSessionBaseWithLoginStatus(statusType)
	assert.NotNil(t, session.SID)
	assert.Equal(t, statusType, session.Status)
	assert.Empty(t, session.Actions)
	assert.Empty(t, session.Options)
}

// TestUserSessionBaseGetters checks the getters of the UserSessionBase struct.
func TestUserSessionBaseGetters(t *testing.T) {
	session := NewUserSessionBase()
	session.Status = LOGIN_SESSION_STATUS_OK
	session.Username = "testuser"
	session.LastLogin = time.Now()
	session.LanguageType = "en-US"
	session.DisplayName = "Test User"
	session.Avatar = aimage.CreateImageCircleAvatar("A")

	assert.Equal(t, session.SID, session.GetSID())
	assert.Equal(t, session.Status, session.GetStatusType())
	assert.Equal(t, session.Username, session.GetUsername())
	assert.Equal(t, session.LastLogin, session.GetLastLogin())
	assert.Equal(t, session.LanguageType, session.GetLanguageType())
	assert.Equal(t, session.DisplayName, session.GetDisplayName())
	assert.Equal(t, session.Avatar, session.GetAvatar())
	assert.Equal(t, session.Options, session.GetOptions())
}

// TestUserSessionBaseIsLoggedIn checks if the IsLoggedIn method correctly identifies a logged-in user.
func TestUserSessionBaseIsLoggedIn(t *testing.T) {
	session := NewUserSessionBase()
	session.Status = LOGIN_SESSION_STATUS_OK
	session.Username = "testuser"

	assert.True(t, session.IsLoggedIn())
}

// TestUserSessionBaseActionMethods checks the action-related methods of the UserSessionBase struct.
func TestUserSessionBaseActionMethods(t *testing.T) {
	session := NewUserSessionBase()
	actionKey := ActionKey("2factor")

	assert.False(t, session.HasAction(actionKey))

	session.AddAction(actionKey)
	assert.True(t, session.HasAction(actionKey))

	session.RemoveAction(actionKey)
	assert.False(t, session.HasAction(actionKey))
}

// TestUserSessionBaseOptionMethods checks the option-related methods of the UserSessionBase struct.
func TestUserSessionBaseOptionMethods(t *testing.T) {
	session := NewUserSessionBase()
	optionKey := atags.TagKey("theme")
	optionValue := "dark"

	assert.False(t, session.HasOption(optionKey))

	session.OptionSet(optionKey, optionValue)
	assert.True(t, session.HasOption(optionKey))
	assert.Equal(t, optionValue, session.OptionValue(optionKey))

	session.OptionRemove(optionKey)
	assert.False(t, session.HasOption(optionKey))

	session.OptionsClear()
	assert.Empty(t, session.GetOptions())
}

// TestUserSessionBaseGetAvatarDataBase64 checks the GetAvatarDataBase64 method.
func TestUserSessionBaseGetAvatarDataBase64(t *testing.T) {
	session := NewUserSessionBase()
	session.DisplayName = "Test User"
	avatarData := session.GetAvatarDataBase64()
	assert.NotEmpty(t, avatarData)
}

func TestUserSessionBase_Meta(t *testing.T) {
	type TestMeta struct {
		Field1 string `json:"field1"`
		Field2 int    `json:"field2"`
	}

	session := &UserSessionBase{}

	// Test saving metadata
	metaToSave := TestMeta{
		Field1: "test_value",
		Field2: 42,
	}
	err := session.SaveMeta(metaToSave)
	assert.NoError(t, err)

	// Validate the saved Meta field
	expectedMeta := `{"field1":"test_value","field2":42}`
	assert.JSONEq(t, expectedMeta, string(session.Meta))

	// Test loading metadata
	var loadedMeta TestMeta
	err = session.LoadMeta(&loadedMeta)
	assert.NoError(t, err)
	assert.Equal(t, metaToSave, loadedMeta)

	// Test clearing metadata
	err = session.SaveMeta(nil)
	assert.NoError(t, err)
	assert.Nil(t, session.Meta)
}
