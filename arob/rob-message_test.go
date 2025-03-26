package arob

import (
	"testing"
)

// TestNewROBMessage tests the NewROBMessage function to ensure it creates a new ROBMessage instance.
func TestNewROBMessage(t *testing.T) {
	message := "test message"
	robMessage := NewROBMessage(message)
	if string(robMessage) != message {
		t.Errorf("NewROBMessage() = %v, want %v", robMessage, message)
	}
}

// TestNewROBMessagef tests the NewROBMessagef function for creating a new formatted ROBMessage.
func TestNewROBMessagef(t *testing.T) {
	format := "test %s"
	arg := "message"
	want := ROBMessage("test message")
	robMessage := NewROBMessagef(ROBMessage(format), arg)
	if robMessage != want {
		t.Errorf("NewROBMessagef() = %v, want %v", robMessage, want)
	}
}

// TestROBMessage_IsEmpty tests the IsEmpty method to ensure it correctly identifies empty messages.
func TestROBMessage_IsEmpty(t *testing.T) {
	tests := []struct {
		name string
		rmt  ROBMessage
		want bool
	}{
		{"Empty", "", true},
		{"WhitespaceOnly", "   ", true},
		{"NonEmpty", "message", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.rmt.IsEmpty(); got != tt.want {
				t.Errorf("ROBMessage.IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestROBMessage_TrimSpace tests the TrimSpace method to ensure it correctly trims whitespace.
func TestROBMessage_TrimSpace(t *testing.T) {
	tests := []struct {
		name string
		rmt  ROBMessage
		want ROBMessage
	}{
		{"NoTrim", "message", "message"},
		{"TrimLeading", "  message", "message"},
		{"TrimTrailing", "message  ", "message"},
		{"TrimBoth", "  message  ", "message"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.rmt.TrimSpace(); got != tt.want {
				t.Errorf("ROBMessage.TrimSpace() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestROBMessage_AppendE tests the AppendE method to ensure it correctly appends an error number to the message.
func TestROBMessage_AppendE(t *testing.T) {
	message := ROBMessage("error")
	num := 1
	want := ROBMessage("error (e1)")
	got := message.AppendE(num)
	if got != want {
		t.Errorf("ROBMessage.AppendE() = %v, want %v", got, want)
	}
}

// Additional tests for AppendR, AppendA, and Append can be added similarly.
