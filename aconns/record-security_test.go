package aconns

import (
	"encoding/json"
	"fmt"
	"github.com/jpfluger/alibs-slim/alog"
	"github.com/jpfluger/alibs-slim/atags"
	"github.com/jpfluger/alibs-slim/atime"
	"github.com/jpfluger/alibs-slim/auser"
	"github.com/rs/zerolog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewRecordSecurity(t *testing.T) {
	user := auser.NewRecordUserIdentityByEmail("user@example.com")
	action := RecordActionType("CREATE")
	event := "Created a new record"

	rs := NewRecordSecurity(user, action, event)

	assert.Equal(t, user, rs.User)
	assert.Equal(t, action, rs.Action)
	assert.Equal(t, event, rs.Event)
	assert.False(t, rs.Time.IsZero(), "Expected Time to be set")
	assert.Empty(t, rs.History, "Expected History to be empty")
}

func TestRecordSecurity_UpdateFrom(t *testing.T) {
	// Mock the current time for consistency
	mockNow := time.Date(2025, time.January, 4, 15, 57, 42, 0, time.UTC)
	nowUpdateFuncUTC = func() time.Time {
		return mockNow
	}
	defer func() { nowUpdateFuncUTC = time.Now }() // Restore the original function

	user1 := auser.NewRecordUserIdentityByEmail("user1@example.com")
	user2 := auser.NewRecordUserIdentityByEmail("user2@example.com")
	action1 := RecordActionType("CREATE")
	action2 := RecordActionType("UPDATE")
	event1 := "Initial creation"
	event2 := "Updated record"

	rs1 := NewRecordSecurity(user1, action1, event1)
	rs1OriginalTime := rs1.Time

	// Ensure initial timestamps are correct
	assert.Equal(t, rs1OriginalTime, rs1.Time, "Expected Time to remain unchanged")

	// Perform the update
	err := rs1.UpdateFrom(NewRecordSecurity(user2, action2, event2))
	assert.NoError(t, err)

	// Verify history
	assert.Len(t, rs1.History, 1, "Expected one history entry")
	historyEntry := rs1.History[0]
	assert.Equal(t, user1, historyEntry.User, "History entry User mismatch")
	assert.Equal(t, action1, historyEntry.Action, "History entry Action mismatch")
	assert.Equal(t, event1, historyEntry.Event, "History entry Event mismatch")
	assert.Equal(t, rs1OriginalTime, historyEntry.Time, "History entry Time mismatch")

	// Verify updated fields
	assert.Equal(t, user2, rs1.User, "Expected User to be updated")
	assert.Equal(t, action2, rs1.Action, "Expected Action to be updated")
	assert.Equal(t, event2, rs1.Event, "Expected Event to be updated")
	assert.Equal(t, mockNow, rs1.Time, "Expected Time to match mockNow")
}

func TestRecordSecurity_UpdateFrom_InvalidTarget(t *testing.T) {
	user := auser.NewRecordUserIdentityByEmail("user@example.com")
	action := RecordActionType("CREATE")
	event := "Test event"

	rs1 := NewRecordSecurity(user, action, event)

	// Test with nil target
	err := rs1.UpdateFrom(nil)
	assert.Error(t, err)
	assert.Equal(t, "target RecordSecurity is nil", err.Error())

	// Test with invalid target
	rsInvalid := &RecordSecurity{}
	err = rs1.UpdateFrom(rsInvalid)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "target RecordSecurity is invalid")
}

func TestRecordSecurity_IsValid(t *testing.T) {
	user := auser.NewRecordUserIdentityByEmail("user@example.com")
	action := RecordActionType("CREATE")
	event := "Valid record"
	now := time.Now()

	validRS := &RecordSecurity{
		User:    user,
		Action:  action,
		Event:   event,
		Time:    now,
		History: make(RecordSecurityHistoryTimes, 0),
	}
	assert.NoError(t, validRS.IsValid())

	// Test missing action
	invalidRS := &RecordSecurity{
		User:    user,
		Event:   event,
		Time:    now,
		History: make(RecordSecurityHistoryTimes, 0),
	}
	err := invalidRS.IsValid()
	assert.Error(t, err)
	assert.Equal(t, "action is required", err.Error())

	// Test missing user
	invalidRS = &RecordSecurity{
		Action:  action,
		Event:   event,
		Time:    now,
		History: make(RecordSecurityHistoryTimes, 0),
	}
	err = invalidRS.IsValid()
	assert.Error(t, err)
	assert.Equal(t, "user is required unless the action is import or adminfix", err.Error())

	// Test missing time
	invalidRS = &RecordSecurity{
		User:    user,
		Action:  action,
		Event:   event,
		History: make(RecordSecurityHistoryTimes, 0),
	}
	err = invalidRS.IsValid()
	assert.Error(t, err)
	assert.Equal(t, "time must be set", err.Error())
}

func TestRecordSecurity_EnsureEventDefaults(t *testing.T) {
	user := auser.NewRecordUserIdentityByEmail("user@example.com")
	action := RecordActionType("UPDATE")
	rs := NewRecordSecurity(user, action, "")

	cols := DbColumnNames{"name", "email"}
	rs.EnsureEventDefaults(cols)

	expectedEvent := "Updating name, email"
	assert.Equal(t, expectedEvent, rs.Event)

	// Test when event is already set
	rs.Event = "Custom event"
	rs.EnsureEventDefaults(cols)
	assert.Equal(t, "Custom event", rs.Event)
}

func TestLogSecurityRecord_WithMockWriter(t *testing.T) {
	// Setup logger with mock provisioner
	prov, err := alog.SetupMockLogger(LOGGER_SECURITYRECORD, zerolog.InfoLevel)
	assert.NoError(t, err)

	// Log a security event
	user := auser.NewRecordUserIdentityByEmail("user@example.com")
	action := RecordActionType("UPDATE")
	event := "Updated user record"
	recIds := atags.TagArrStrings{
		&atags.TagKeyValueString{Key: "record_id", Value: "12345"},
	}
	sec := RecordSecurity{
		User:   user,
		Action: action,
		Event:  event,
		Time:   time.Now().UTC(),
		RecIds: recIds,
	}

	err = LogSecurityRecord(sec)
	assert.NoError(t, err)

	// Validate logs
	assert.NotNil(t, prov.Writer)
	assert.Greater(t, len(prov.Writer.Logs), 0)

	// Parse the first log entry
	var logOutput map[string]interface{}
	err = json.Unmarshal([]byte(prov.Writer.Logs[0]), &logOutput)
	assert.NoError(t, err)

	// Validate log fields
	assert.Equal(t, "secrec", logOutput["message"])
	assert.Equal(t, "info", logOutput["level"])

	// Validate time field
	timeString, ok := logOutput["time"].(string)
	assert.True(t, ok, "Expected 'time' to be a string")
	parsedTime, err := time.Parse(time.RFC3339Nano, timeString)
	assert.NoError(t, err, "Failed to parse 'time' field")

	// Convert time to different zones and validate
	utcTime := parsedTime.UTC()
	localTime := atime.ConvertToTimeZone(parsedTime, "America/New_York")
	if localTime.IsZero() {
		assert.Fail(t, "Failed to convert to America/New_York")
		return
	}

	fmt.Printf("Original Time (UTC): %s\n", utcTime)
	fmt.Printf("Converted Time (New York): %s\n", localTime)

	// Validate serialized RecordSecurity
	recData, ok := logOutput["rec"].(map[string]interface{})
	assert.True(t, ok, "Expected 'rec' to be a map")

	// Validate user structure
	userData, ok := recData["user"].(map[string]interface{})
	assert.True(t, ok, "Expected 'user' to be a map")
	assert.Equal(t, map[string]interface{}{
		"ids": map[string]interface{}{
			"email": "user@example.com",
		},
		"uid": nil, // Adjust as per the expected UID representation
	}, userData)

	// Validate Record IDs
	recIdsData, ok := recData["recIds"].([]interface{})
	assert.True(t, ok, "Expected 'recIds' to be an array")
	assert.Equal(t, map[string]interface{}{"key": "record_id", "value": "12345"}, recIdsData[0])
}

func TestLogSecurityRecord_WithMeta(t *testing.T) {
	// Set up the mock logger
	prov, err := alog.SetupMockLogger(LOGGER_SECURITYRECORD, zerolog.InfoLevel)
	assert.NoError(t, err)

	// Define metadata for the log
	meta := map[string]interface{}{
		"changes": []map[string]interface{}{
			{"field": "status", "from": "pending", "to": "approved"},
			{"field": "amount", "from": "100", "to": "200"},
		},
		"approver": "manager@example.com",
		"reason":   "Compliance audit",
	}

	// Serialize the Meta field
	metaBytes, err := json.Marshal(meta)
	assert.NoError(t, err)

	// Create a RecordSecurity instance
	user := auser.NewRecordUserIdentityByEmail("user@example.com")
	sec := RecordSecurity{
		User:   user,
		Action: RecordActionType("UPDATE"),
		Event:  "Updated financial record",
		Time:   time.Now().UTC(),
		RecIds: atags.TagArrStrings{
			&atags.TagKeyValueString{Key: "record_id", Value: "12345"},
		},
		Meta: metaBytes,
	}

	// Log the security record
	err = LogSecurityRecord(sec)
	assert.NoError(t, err)

	// Validate logs
	assert.NotNil(t, prov.Writer)
	assert.Greater(t, len(prov.Writer.Logs), 0)

	// Parse the first log entry
	var logOutput map[string]interface{}
	err = json.Unmarshal([]byte(prov.Writer.Logs[0]), &logOutput)
	assert.NoError(t, err)

	// Validate core log fields
	assert.Equal(t, "secrec", logOutput["message"])
	assert.Equal(t, "info", logOutput["level"])
	assert.NotNil(t, logOutput["time"])

	// Validate serialized RecordSecurity
	recData, ok := logOutput["rec"].(map[string]interface{})
	assert.True(t, ok, "Expected 'rec' to be a map")

	// Validate Meta field
	parsedMeta, ok := recData["meta"].(map[string]interface{})
	assert.True(t, ok, "Expected 'meta' to be a map")

	// Validate dynamic metadata content
	assert.Equal(t, "manager@example.com", parsedMeta["approver"])
	assert.Equal(t, "Compliance audit", parsedMeta["reason"])

	changes, ok := parsedMeta["changes"].([]interface{})
	assert.True(t, ok, "Expected 'changes' to be a list")
	assert.Len(t, changes, 2)

	firstChange := changes[0].(map[string]interface{})
	assert.Equal(t, "status", firstChange["field"])
	assert.Equal(t, "pending", firstChange["from"])
	assert.Equal(t, "approved", firstChange["to"])
}
