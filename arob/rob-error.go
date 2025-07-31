package arob

// ROBErrors is a slice of pointers to ROBError, representing a collection of errors.
type ROBErrors []*ROBError

// ROBError represents an individual error, warning, notice, info, or debug message.
type ROBError struct {
	Type    ROBType       `json:"type,omitempty" xml:"type,omitempty"`       // The severity type of the error.
	Message ROBMessage    `json:"message,omitempty" xml:"message,omitempty"` // The error message content.
	Field   ROBErrorField `json:"field,omitempty" xml:"field,omitempty"`     // The field associated with the error, if any.
	Stack   string        `json:"stack,omitempty" xml:"stack,omitempty"`     // The stack trace, if available.
}

// NewROBError creates a new ROBError with the specified type and message.
func NewROBError(errType ROBType, message ROBMessage) *ROBError {
	return &ROBError{Type: NormalizeROBType(errType), Message: message}
}

// NewROBErrorWithField creates a new ROBError with the specified type, message, and field.
func NewROBErrorWithField(errType ROBType, message ROBMessage, field ROBErrorField) *ROBError {
	robe := NewROBError(errType, message)
	robe.Field = field
	return robe
}

// IsErr determines if the ROBError is of an error type (as opposed to warning, info, etc.).
func (robError *ROBError) IsErr() bool {
	// Only return true for types that are considered errors.
	return robError.Type == ROBTYPE_ERROR || robError.Type == ROBTYPE_CRITICAL || robError.Type == ROBTYPE_EMERGENCY
}

// ToStringArray converts a slice of ROBErrors to a slice of their string messages, filtered by type.
func (robes ROBErrors) ToStringArray(robType ROBType) []string {
	var arr []string
	for _, robe := range robes {
		if robe.Type == robType {
			arr = append(arr, robe.Message.String())
		}
	}
	return arr
}
