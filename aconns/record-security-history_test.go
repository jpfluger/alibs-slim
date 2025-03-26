package aconns

import (
	"github.com/jpfluger/alibs-slim/auser"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRecordSecurityHistoryTimes_AddEntry(t *testing.T) {
	var history RecordSecurityHistoryTimes

	user := auser.NewRecordUserIdentityByEmail("user@example.com")
	action := RecordActionType("LOGIN")
	event := "User logged in"
	eventTime := time.Now()

	history.AddEntry(user, action, event, eventTime)

	assert.Len(t, history, 1)
	assert.Equal(t, user, history[0].User)
	assert.Equal(t, action, history[0].Action)
	assert.Equal(t, event, history[0].Event)
	assert.WithinDuration(t, eventTime, history[0].Time, time.Second)
}

func TestRecordSecurityHistoryTimes_LatestEntry(t *testing.T) {
	var history RecordSecurityHistoryTimes

	// Test with no entries
	latest := history.LatestEntry()
	assert.Nil(t, latest)

	// Add entries and test
	user := auser.NewRecordUserIdentityByEmail("user@example.com")
	action := RecordActionType("LOGIN")
	history.AddEntry(user, action, "First Event", time.Now().Add(-time.Hour))
	history.AddEntry(user, action, "Second Event", time.Now())

	latest = history.LatestEntry()
	assert.NotNil(t, latest)
	assert.Equal(t, "Second Event", latest.Event)
}

func TestRecordSecurityHistoryTimes_FindEntriesByUser(t *testing.T) {
	var history RecordSecurityHistoryTimes

	user1 := auser.NewRecordUserIdentityByEmail("user1@example.com")
	user2 := auser.NewRecordUserIdentityByEmail("user2@example.com")
	action := RecordActionType("LOGIN")

	history.AddEntry(user1, action, "User1 Event1", time.Now())
	history.AddEntry(user2, action, "User2 Event1", time.Now())
	history.AddEntry(user1, action, "User1 Event2", time.Now())

	entries := history.FindEntriesByUser(user1)
	assert.Len(t, entries, 2)
	assert.Equal(t, "User1 Event1", entries[0].Event)
	assert.Equal(t, "User1 Event2", entries[1].Event)

	entries = history.FindEntriesByUser(user2)
	assert.Len(t, entries, 1)
	assert.Equal(t, "User2 Event1", entries[0].Event)
}

func TestRecordSecurityHistoryTimes_FilterByAction(t *testing.T) {
	var history RecordSecurityHistoryTimes

	user := auser.NewRecordUserIdentityByEmail("user@example.com")
	action1 := RecordActionType("LOGIN")
	action2 := RecordActionType("LOGOUT")

	history.AddEntry(user, action1, "Login Event", time.Now())
	history.AddEntry(user, action2, "Logout Event", time.Now())

	filtered := history.FilterByAction(action1)
	assert.Len(t, filtered, 1)
	assert.Equal(t, "Login Event", filtered[0].Event)

	filtered = history.FilterByAction(action2)
	assert.Len(t, filtered, 1)
	assert.Equal(t, "Logout Event", filtered[0].Event)

	filtered = history.FilterByAction("INVALID")
	assert.Empty(t, filtered)
}
