package aconns

import (
	"strings"
)

// RecordActionType represents the type of action performed on a database record.
type RecordActionType string

// Constants for RecordActionType values.
const (
	// Actions that typically require a username.
	REC_ACTION_INSERT RecordActionType = "insert" // Insert a new record.
	REC_ACTION_UPDATE RecordActionType = "update" // Update an existing record.
	REC_ACTION_DELETE RecordActionType = "delete" // Delete a record.
	REC_ACTION_UPSERT RecordActionType = "upsert" // Insert or update a record.

	// Events that do not require a username.
	REC_EVENT_IMPORT   RecordActionType = "import"   // Import records in bulk.
	REC_EVENT_ADMINFIX RecordActionType = "adminfix" // Administrative fixes to records.
)

// IsEmpty checks if the RecordActionType is empty after trimming whitespace.
func (rt RecordActionType) IsEmpty() bool {
	return strings.TrimSpace(string(rt)) == ""
}

// TrimSpace returns a new RecordActionType with leading and trailing whitespace removed.
func (rt RecordActionType) TrimSpace() RecordActionType {
	return RecordActionType(strings.TrimSpace(string(rt)))
}

// String converts the RecordActionType to a regular string.
func (rt RecordActionType) String() string {
	return strings.TrimSpace(string(rt))
}
