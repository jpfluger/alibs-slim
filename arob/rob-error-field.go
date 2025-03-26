package arob

import (
	"strings"
)

// ROBErrorField represents a field in an error object, which can be used to identify where an error occurred.
type ROBErrorField string

// Predefined constants for system-level error fields.
const (
	ROBERRORFIELD_SYSTEM               = ROBErrorField("_system")              // Default system field.
	ROBERRORFIELD_SYSTEM_DISPLAY       = ROBErrorField("_system:display")      // System field for display errors.
	ROBERRORFIELD_SYSTEM_DISPLAY_MODAL = ROBErrorField("_system:displayModal") // System field for modal display errors.
)

// IsEmpty checks if the ROBErrorField is empty after trimming whitespace.
func (fld ROBErrorField) IsEmpty() bool {
	return strings.TrimSpace(string(fld)) == ""
}

// TrimSpace trims leading and trailing whitespace from the ROBErrorField.
func (fld ROBErrorField) TrimSpace() ROBErrorField {
	return ROBErrorField(strings.TrimSpace(string(fld)))
}

// String returns the string representation of the ROBErrorField.
func (fld ROBErrorField) String() string {
	return string(fld)
}

// GetType extracts the type part of the ROBErrorField, before any colon.
func (fld ROBErrorField) GetType() string {
	parts := strings.SplitN(string(fld), ":", 2)
	return parts[0]
}

// GetSubType extracts the subtype part of the ROBErrorField, after the colon.
func (fld ROBErrorField) GetSubType() string {
	parts := strings.SplitN(string(fld), ":", 2)
	if len(parts) == 2 {
		return parts[1]
	}
	return ""
}
