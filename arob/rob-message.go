package arob

import (
	"fmt"
	"strings"
)

// Predefined titles for various error messages.
var TitleErr = NewROBMessage(`Error`)
var TitleSystemErr = NewROBMessage(`System Error`)
var TitlePageNotFound = NewROBMessage(`Page Not Found`)
var TitleNotPermitted = NewROBMessage(`Not Permitted`)

// Predefined error messages for common scenarios.
var ErrServer = NewROBMessage(`Could not process request.`)
var ErrServerWithSupport = NewROBMessage(`Could not process request - contact support if this persists.`)
var ErrPageNotFound = NewROBMessage(`The requested page could not be found.`)
var ErrNotPermitted = NewROBMessage(`The requested action is not permitted.`)

// ROBMessage represents a message in the context of a ROB (Return/Request/Response Object).
type ROBMessage string

// NewROBMessage creates a new ROBMessage from a string.
func NewROBMessage(message string) ROBMessage {
	return ROBMessage(message)
}

// NewROBMessagef creates a new formatted ROBMessage.
func NewROBMessagef(message ROBMessage, v ...interface{}) ROBMessage {
	return ROBMessage(fmt.Sprintf(message.String(), v...))
}

// IsEmpty checks if the ROBMessage is empty after trimming whitespace.
func (rmt ROBMessage) IsEmpty() bool {
	return strings.TrimSpace(string(rmt)) == ""
}

// TrimSpace trims leading and trailing whitespace from the ROBMessage.
func (rmt ROBMessage) TrimSpace() ROBMessage {
	return ROBMessage(strings.TrimSpace(string(rmt)))
}

// String returns the string representation of the ROBMessage.
func (rmt ROBMessage) String() string {
	return string(rmt)
}

// AppendE appends an error number to the ROBMessage, prefixed with 'e'.
func (rmt ROBMessage) AppendE(num int) ROBMessage {
	return rmt.Append("e", num)
}

// AppendR appends a request/response number to the ROBMessage, prefixed with 'r'.
func (rmt ROBMessage) AppendR(num int) ROBMessage {
	return rmt.Append("r", num)
}

// AppendA appends an authentication number to the ROBMessage, prefixed with 'a'.
func (rmt ROBMessage) AppendA(num int) ROBMessage {
	return rmt.Append("a", num)
}

// Append adds a qualifier and number to the ROBMessage. If the message contains a placeholder for a number, it formats the message with the number.
func (rmt ROBMessage) Append(qualifier string, num int) ROBMessage {
	if num <= 0 {
		return rmt
	}
	if strings.Contains(rmt.String(), "%d") {
		return ROBMessage(fmt.Sprintf(rmt.String(), num))
	}
	if qualifier == "" {
		qualifier = "e"
	}
	if rmt.IsEmpty() {
		return ROBMessage(fmt.Sprintf("%s%d", qualifier, num))
	}
	return ROBMessage(fmt.Sprintf("%s (%s%d)", rmt.String(), qualifier, num))
}
