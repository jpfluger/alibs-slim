package aconns

import (
	"encoding/json"
	"fmt"
	"github.com/jpfluger/alibs-slim/aerr"
	"github.com/jpfluger/alibs-slim/alog"
	"github.com/jpfluger/alibs-slim/atags"
	"github.com/jpfluger/alibs-slim/auser"
	"github.com/rs/zerolog"
	"strings"
	"time"
)

// nowUpdateFuncUTC is a configurable function for retrieving the current UTC time.
var nowUpdateFuncUTC = func() time.Time {
	return time.Now().UTC()
}

// RecordSecurity is used for tracking the security-related changes and actions
// performed on a database record. It includes the user who performed the action,
// the action itself, and a history of past actions.
type RecordSecurity struct {
	User    auser.RecordUserIdentity   `json:"user,omitempty"`
	Action  RecordActionType           `json:"action,omitempty"`
	Event   string                     `json:"event,omitempty"`
	Time    time.Time                  `json:"time,omitempty"`
	History RecordSecurityHistoryTimes `json:"history,omitempty"`
	Error   *aerr.Error                `json:"error,omitempty"`
	RecIds  atags.TagArrStrings        `json:"recIds,omitempty"`
	Meta    json.RawMessage            `json:"meta,omitempty"`
}

// NewRecordSecurity creates a new RecordSecurity instance with the current time.
func NewRecordSecurity(user auser.RecordUserIdentity, action RecordActionType, event string) *RecordSecurity {
	return &RecordSecurity{
		User:    user,
		Action:  action,
		Event:   event,
		Time:    time.Now().UTC(),
		History: make(RecordSecurityHistoryTimes, 0),
	}
}

// UpdateFrom updates the current RecordSecurity with information from another instance,
// adding the previous state to the history.
func (rs *RecordSecurity) UpdateFrom(target *RecordSecurity) error {
	if target == nil {
		return fmt.Errorf("target RecordSecurity is nil")
	}

	if err := target.IsValid(); err != nil {
		return fmt.Errorf("target RecordSecurity is invalid: %v", err)
	}

	// Add the current state to history before updating
	rs.History.AddEntry(rs.User, rs.Action, rs.Event, rs.Time)

	// Update fields from the target
	rs.User = target.User
	rs.Action = target.Action
	rs.Event = target.Event
	rs.Time = nowUpdateFuncUTC() // Use the configurable time function

	return nil
}

// IsValid checks if the RecordSecurity instance is valid by ensuring that
// required fields are set and not empty.
func (rs *RecordSecurity) IsValid() error {
	if rs.Action.IsEmpty() {
		return fmt.Errorf("action is required")
	}

	if rs.User.IsEmpty() && rs.Action != REC_EVENT_IMPORT && rs.Action != REC_EVENT_ADMINFIX {
		return fmt.Errorf("user is required unless the action is import or adminfix")
	}

	if rs.Time.IsZero() {
		return fmt.Errorf("time must be set")
	}

	return nil
}

// EnsureEventDefaults sets a default event description based on the provided column names
// if the event is not already set.
func (rs *RecordSecurity) EnsureEventDefaults(cols DbColumnNames) {
	if rs == nil || rs.Event != "" {
		return
	}

	var sb strings.Builder
	sb.WriteString("Updating ")
	for i, col := range cols {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(col.String())
	}
	rs.Event = sb.String()
}

func (rs RecordSecurity) MarshalJSON_Log_NoTimeNoHistory() ([]byte, error) {
	// Define a struct excluding Time and History for logging
	type loggableRecordSecurity struct {
		User   auser.RecordUserIdentity `json:"user,omitempty"`
		Action RecordActionType         `json:"action,omitempty"`
		Event  string                   `json:"event,omitempty"`
		Error  *aerr.Error              `json:"error,omitempty"`
		RecIds atags.TagArrStrings      `json:"recIds,omitempty"`
		Meta   json.RawMessage          `json:"meta,omitempty"`
	}

	return json.Marshal(loggableRecordSecurity{
		User:   rs.User,
		Action: rs.Action,
		Event:  rs.Event,
		Error:  rs.Error,
		RecIds: rs.RecIds,
		Meta:   rs.Meta,
	})
}

const (
	LOGGER_SECURITYRECORD alog.ChannelLabel = "secrec" // Security Record Channel
)

func LogSecurityRecord(sec RecordSecurity) error {
	return LogSecurityRecordWithOptions(sec, true)
}

func LogSecurityRecordWithOptions(sec RecordSecurity, doValidate bool) error {
	// Optional validation
	if doValidate {
		if err := sec.IsValid(); err != nil {
			return err
		}
	}

	var eventLog *zerolog.Event
	if sec.Error != nil {
		eventLog = alog.LOGGER(LOGGER_SECURITYRECORD).Err(sec.Error.ToError())
	} else {
		eventLog = alog.LOGGER(LOGGER_SECURITYRECORD).Info()
	}

	if eventLog == nil {
		return fmt.Errorf("event log is nil")
	}

	// Serialize the RecordSecurity excluding Time and History
	bRS, err := sec.MarshalJSON_Log_NoTimeNoHistory()
	if err != nil {
		return fmt.Errorf("failed to marshal RecordSecurity: %w", err)
	}

	// Add fields to the log event
	eventLog = eventLog.
		Time("time", sec.Time). // Explicit timestamp
		RawJSON("rec", bRS)     // Include serialized RecordSecurity

	// Log the event with a message
	eventLog.Msg(LOGGER_SECURITYRECORD.String())

	return nil
}
