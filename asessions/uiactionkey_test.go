package asessions

import (
	"github.com/jpfluger/alibs-slim/auser"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/jpfluger/alibs-slim/atags"
)

func TestNewUIActionKey(t *testing.T) {
	action := ActionKey("test-action")
	username := auser.Username("test-user")
	urlPost := "http://example.com/post"
	formDate := time.Now().UTC()
	title := "Test Title"
	lead := "Test Lead"
	footer := "Test Footer"
	tags := atags.TagMapString{"key1": "value1"}

	uiActionKey := NewUIActionKey(action, username, urlPost, formDate, title, lead, footer, tags)

	assert.Equal(t, action, uiActionKey.GetAction())
	assert.Equal(t, username, uiActionKey.GetUsername())
	assert.Equal(t, urlPost, uiActionKey.GetUrlPost())
	assert.Equal(t, formDate, uiActionKey.GetFormDate())
	assert.Equal(t, title, uiActionKey.GetTitle())
	assert.Equal(t, lead, uiActionKey.GetLead())
	assert.Equal(t, footer, uiActionKey.GetFooter())
	assert.Equal(t, tags, uiActionKey.GetTags())
}

func TestSetFormDate(t *testing.T) {
	uiActionKey := &UIActionKey{}
	newFormDate := time.Date(2024, 8, 22, 12, 0, 0, 0, time.UTC)
	uiActionKey.SetFormDate(newFormDate)
	assert.Equal(t, newFormDate, uiActionKey.GetFormDate())
}
