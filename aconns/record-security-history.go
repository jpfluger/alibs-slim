package aconns

import (
	"github.com/jpfluger/alibs-slim/auser"
	"time"
)

// RecordSecurityHistoryTime represents a single historical security record.
type RecordSecurityHistoryTime struct {
	User   auser.RecordUserIdentity `json:"user"`
	Action RecordActionType         `json:"action"`
	Event  string                   `json:"event"`
	Time   time.Time                `json:"time"`
}

// RecordSecurityHistoryTimes is a slice of RecordSecurityHistoryTime.
type RecordSecurityHistoryTimes []*RecordSecurityHistoryTime

// AddEntry adds a new entry to the RecordSecurityHistoryTimes.
func (rsh *RecordSecurityHistoryTimes) AddEntry(user auser.RecordUserIdentity, action RecordActionType, event string, time time.Time) {
	entry := &RecordSecurityHistoryTime{
		User:   user,
		Action: action,
		Event:  event,
		Time:   time,
	}
	*rsh = append(*rsh, entry)
}

// LatestEntry returns the most recent entry in the RecordSecurityHistoryTimes.
func (rsh RecordSecurityHistoryTimes) LatestEntry() *RecordSecurityHistoryTime {
	if len(rsh) == 0 {
		return nil
	}
	return rsh[len(rsh)-1]
}

// FindEntriesByUser returns all entries in the RecordSecurityHistoryTimes performed by a given user.
func (rsh RecordSecurityHistoryTimes) FindEntriesByUser(user auser.RecordUserIdentity) []*RecordSecurityHistoryTime {
	var entries []*RecordSecurityHistoryTime
	for _, entry := range rsh {
		if entry.User.HasMatch(user) {
			entries = append(entries, entry)
		}
	}
	return entries
}

// FilterByAction returns a new RecordSecurityHistoryTimes containing only the entries
// that match the given action.
func (rsh RecordSecurityHistoryTimes) FilterByAction(action RecordActionType) RecordSecurityHistoryTimes {
	var filtered RecordSecurityHistoryTimes
	for _, entry := range rsh {
		if entry.Action == action {
			filtered = append(filtered, entry)
		}
	}
	return filtered
}
