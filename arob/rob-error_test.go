package arob

import (
	"testing"
)

// TestNewROBError tests the NewROBError function to ensure it creates a new ROBError instance with the correct type and message.
func TestNewROBError(t *testing.T) {
	errType := ROBTYPE_ERROR
	message := ROBMessage("error occurred")
	robError := NewROBError(errType, message)
	if robError.Type != errType {
		t.Errorf("NewROBError() Type = %v, want %v", robError.Type, errType)
	}
	if robError.Message != message {
		t.Errorf("NewROBError() Message = %v, want %v", robError.Message, message)
	}
}

// TestNewROBErrorWithField tests the NewROBErrorWithField function to ensure it creates a new ROBError instance with the correct type, message, and field.
func TestNewROBErrorWithField(t *testing.T) {
	errType := ROBTYPE_ERROR
	message := ROBMessage("error occurred")
	field := ROBErrorField("field")
	robError := NewROBErrorWithField(errType, message, field)
	if robError.Field != field {
		t.Errorf("NewROBErrorWithField() Field = %v, want %v", robError.Field, field)
	}
}

// TestIsErr tests the IsErr method to ensure it correctly identifies if the ROBError is an error type.
func TestIsErr(t *testing.T) {
	robError := NewROBError(ROBTYPE_ERROR, "error occurred")
	if !robError.IsErr() {
		t.Error("IsErr() should return true for ROBTYPE_ERROR")
	}
}

// TestToStringArray tests the ToStringArray method to ensure it correctly converts a slice of ROBErrors to a slice of string messages.
func TestToStringArray(t *testing.T) {
	errType := ROBTYPE_ERROR
	message := ROBMessage("error occurred")
	robErrors := ROBErrors{
		NewROBError(errType, message),
	}
	strArray := robErrors.ToStringArray(errType)
	if len(strArray) != 1 || strArray[0] != message.String() {
		t.Errorf("ToStringArray() = %v, want [%v]", strArray, message.String())
	}
}

// Additional tests for other methods and error types can be added similarly.
