package arob

import "fmt"

// ROBLog represents a single log message with a type and content.
// The `isNormalized` field ensures the log type is processed via NormalizeROBType only once.
type ROBLog struct {
	Type         ROBType `json:"type"`
	Message      string  `json:"message"`
	isNormalized bool
}

// GetType returns the normalized ROBType for this log.
// It memoizes the normalization to avoid repeated work.
func (rl *ROBLog) GetType() ROBType {
	if !rl.isNormalized {
		rl.Type = NormalizeROBType(rl.Type)
		rl.isNormalized = true
	}
	return rl.Type
}

// String returns a formatted string representation of the ROBLog.
func (rl ROBLog) String() string {
	return fmt.Sprintf("[%s] %s", rl.GetType().String(), rl.Message)
}

// ROBLogs represents a slice of ROBLog entries.
type ROBLogs []ROBLog

// HasLogType checks if any log in the collection matches one of the given ROBTypes.
func (rls ROBLogs) HasLogType(types ...ROBType) bool {
	for _, log := range rls {
		t := log.GetType()
		for _, check := range types {
			if t == check {
				return true
			}
		}
	}
	return false
}

// FilterByType returns a new ROBLogs slice containing only logs
// whose normalized type matches one of the specified ROBTypes.
func (rls ROBLogs) FilterByType(types ...ROBType) ROBLogs {
	var filtered ROBLogs
	for _, log := range rls {
		t := log.GetType()
		for _, filterType := range types {
			if t == filterType {
				filtered = append(filtered, log)
				break
			}
		}
	}
	return filtered
}

// ToStringArray returns a slice of formatted string representations of each ROBLog.
func (rls ROBLogs) ToStringArray() []string {
	result := make([]string, 0, len(rls))
	for _, log := range rls {
		result = append(result, log.String())
	}
	return result
}
