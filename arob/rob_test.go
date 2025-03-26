package arob

import (
	"testing"
)

// TestNewROB tests the NewROB function to ensure it creates a new ROB instance.
func TestNewROB(t *testing.T) {
	rob := NewROB()
	if rob == nil {
		t.Error("NewROB() should not return nil")
	}
}

// TestNewROBWithMessage tests the NewROBWithMessage function for creating a new ROB with a message.
func TestNewROBWithMessage(t *testing.T) {
	message := ROBMessage("test message")
	rob := NewROBWithMessage(message)
	if rob.Message != message {
		t.Errorf("NewROBWithMessage() message = %v, want %v", rob.Message, message)
	}
}

// TestNewROBWithMessagef tests the NewROBWithMessagef function for creating a new ROB with a formatted message.
func TestNewROBWithMessagef(t *testing.T) {
	format := "test %s"
	arg := "message"
	want := ROBMessage("test message")
	rob := NewROBWithMessagef(ROBMessage(format), arg)
	if rob.Message != want {
		t.Errorf("NewROBWithMessagef() message = %v, want %v", rob.Message, want)
	}
}

// TestNewROBWithRedirect tests the NewROBWithRedirect function for creating a new ROB with a redirect URL.
func TestNewROBWithRedirect(t *testing.T) {
	url := "http://example.com"
	rob := NewROBWithRedirect(url)
	if rob.RedirectUrl != url {
		t.Errorf("NewROBWithRedirect() redirectUrl = %v, want %v", rob.RedirectUrl, url)
	}
}

// TestNewROBWithRecs tests the NewROBWithRecs function for creating a new ROB with records.
func TestNewROBWithRecs(t *testing.T) {
	recs := []string{"record1", "record2"}
	rob := NewROBWithRecs(recs)
	if rob.Recs == nil {
		t.Error("NewROBWithRecs() recs should not be nil")
	}
}

// TestNewROBWithError tests the NewROBWithError function for creating a new ROB with an error.
func TestNewROBWithError(t *testing.T) {
	field := ROBErrorField("field")
	message := ROBMessage("error message")
	rob := NewROBWithError(field, message)
	if !rob.HasErrors() {
		t.Error("NewROBWithError() should have errors")
	}
}

// TestNewROBWithErrorf tests the NewROBWithErrorf function for creating a new ROB with a formatted error.
func TestNewROBWithErrorf(t *testing.T) {
	field := ROBErrorField("field")
	format := "error %s"
	arg := "message"
	want := ROBMessage("error message")
	rob := NewROBWithErrorf(field, ROBMessage(format), arg)
	if !rob.HasErrors() || rob.Errs[0].Message != want {
		t.Errorf("NewROBWithErrorf() message = %v, want %v", rob.Errs[0].Message, want)
	}
}
