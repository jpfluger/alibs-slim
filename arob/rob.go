package arob

import (
	"fmt"
)

// ROB (Return/Request/Response Object) is a struct that encapsulates
// various types of responses for API or internal function calls.
type ROB struct {
	Message      ROBMessage  `json:"message,omitempty" xml:"message,omitempty"`           // The main message to be conveyed.
	MessageType  ROBType     `json:"messageType,omitempty" xml:"messageType,omitempty"`   // The type of message being conveyed.
	MessageTitle string      `json:"messageTitle,omitempty" xml:"messageTitle,omitempty"` // The title of the message.
	RedirectUrl  string      `json:"redirect,omitempty" xml:"redirect,omitempty"`         // URL to redirect to, if applicable.
	Errs         ROBErrors   `json:"errs,omitempty" xml:"errs,omitempty"`                 // Collection of errors.
	Recs         interface{} `json:"recs,omitempty" xml:"recs,omitempty"`                 // Records or data to be returned.
	Columns      interface{} `json:"columns,omitempty" xml:"columns,omitempty"`           // Column definitions, if returning tabular data.
	// Deprecated
	// Paginate     *Paginate   `json:"paginate,omitempty" xml:"paginate,omitempty"`       // Pagination information, if applicable.

	// Type was added especially for API responses to not confuse the WebUI.
	// For example, the WebUI could have MessageType populated and this triggers an event but
	// this Type is API specific only.
	Type ROBType `json:"type,omitempty" xml:"type,omitempty"`
	// Status was added especially for API responses, in lieu of simply `return c.String(http.StatusOK, "ok")`
	Status string `json:"status,omitempty" xml:"status,omitempty"`
}

// NewROB creates a new instance of ROB.
func NewROB() *ROB {
	return &ROB{}
}

// NewROBWithMessage creates a new ROB with a specified message.
func NewROBWithMessage(message ROBMessage) *ROB {
	return &ROB{Message: message}
}

// NewROBWithMessagef creates a new ROB with a formatted message.
func NewROBWithMessagef(message ROBMessage, v ...interface{}) *ROB {
	return &ROB{Message: ROBMessage(fmt.Sprintf(message.String(), v...))}
}

// NewROBWithRedirect creates a new ROB with a redirect URL.
func NewROBWithRedirect(redirectUrl string) *ROB {
	return &ROB{RedirectUrl: redirectUrl}
}

// NewROBRedirectWithMessage creates a new ROB with a redirect URL and a message.
func NewROBRedirectWithMessage(redirectUrl string, message ROBMessage) *ROB {
	return &ROB{RedirectUrl: redirectUrl, Message: message}
}

// NewROBRedirectWithMessagef creates a new ROB with a redirect URL and a formatted message.
func NewROBRedirectWithMessagef(redirectUrl string, message ROBMessage, v ...interface{}) *ROB {
	return &ROB{RedirectUrl: redirectUrl, Message: ROBMessage(fmt.Sprintf(message.String(), v...))}
}

// NewROBWithRecs creates a new ROB with records/data.
func NewROBWithRecs(recs interface{}) *ROB {
	return &ROB{Recs: recs}
}

// NewROBWithMessageTypeTitle creates a new ROB with its MessageType and MessageTitle populated.
func NewROBWithMessageTypeTitle(messageType ROBType, messageTitle string) *ROB {
	return &ROB{MessageType: messageType, MessageTitle: messageTitle}
}

// NewROBWithStatus creates a new ROB with its StatusType and Status populated.
func NewROBWithStatus(statusType ROBType, status string) *ROB {
	return &ROB{Type: statusType, Status: status}
}

// NewROBWithStatusRecs creates a new ROB with its StatusType and Status populated with valid Recs.
func NewROBWithStatusRecs(statusType ROBType, status string, recs interface{}) *ROB {
	return &ROB{Type: statusType, Status: status, Recs: recs}
}

// NewROBWithStatusError creates a new ROB with its StatusType and Status populated with an Error.
func NewROBWithStatusError(statusType ROBType, status string, err error) *ROB {
	rob := &ROB{Type: statusType, Status: status}
	if err != nil {
		return rob
	}
	rob.AddError(ROBERRORFIELD_SYSTEM, NewROBMessage(err.Error()))
	return rob
}

//// NewROBWithRecsPaginate creates a new ROB with records/data and pagination information.
//func NewROBWithRecsPaginate(recs interface{}, paginate *Paginate) *ROB {
//	return &ROB{Recs: recs, Paginate: paginate}
//}
//
//// NewROBWithRecsPaginateColumns creates a new ROB with records/data, pagination information, and column definitions.
//func NewROBWithRecsPaginateColumns(recs interface{}, paginate *Paginate, columns interface{}) *ROB {
//	return &ROB{Recs: recs, Paginate: paginate, Columns: columns}
//}
//
//// NewROBWithRecsPaginateMessage creates a new ROB with records/data, pagination information, and a message.
//func NewROBWithRecsPaginateMessage(recs interface{}, paginate *Paginate, message ROBMessage) *ROB {
//	return &ROB{Recs: recs, Paginate: paginate, Message: message}
//}

// NewROBMessageWithOptionSysError creates a new ROB with a system error or a message.
func NewROBMessageWithOptionSysError(message ROBMessage, isSysError bool) *ROB {
	var rob *ROB
	if !isSysError {
		// If it's not a system error, create a new ROB with the provided message.
		rob = NewROBWithMessage(message)
	} else {
		// If it is a system error, create a new ROB and add the system error message.
		rob = NewROB()
		rob.AddError(ROBERRORFIELD_SYSTEM, message)
	}
	return rob
}

// NewROBWithSysError creates a new ROB with an error message associated with a system field.
func NewROBWithSysError(message ROBMessage) *ROB {
	rob := NewROB()
	rob.AddError(ROBERRORFIELD_SYSTEM, message)
	return rob
}

// NewROBWithError creates a new ROB with an error message associated with a specific field.
func NewROBWithError(field ROBErrorField, message ROBMessage) *ROB {
	rob := NewROB()
	rob.AddError(field, message)
	return rob
}

// NewROBWithErrorf creates a new ROB with a formatted error message associated with a specific field.
func NewROBWithErrorf(field ROBErrorField, message ROBMessage, v ...interface{}) *ROB {
	rob := NewROB()
	rob.AddErrorf(field, message, v...)
	return rob
}

// NewROBWithErrorNum creates a new ROB with an error message associated with a specific field and a number.
func NewROBWithErrorNum(field ROBErrorField, message ROBMessage, num int) *ROB {
	rob := NewROB()
	rob.AddError(field, message.AppendE(num))
	return rob
}

// NewROBWithErrorNumf creates a new ROB with a formatted error message associated with a specific field and a number.
func NewROBWithErrorNumf(field ROBErrorField, message ROBMessage, num int, v ...interface{}) *ROB {
	rob := NewROB()
	message = ROBMessage(fmt.Sprintf(message.String(), v...))
	rob.AddError(field, message.AppendE(num))
	return rob
}

// HasErrors checks if the ROB contains any errors.
func (rob *ROB) HasErrors() bool {
	return rob.Errs != nil && len(rob.Errs) > 0
}

// AddCritical adds a critical error message to the ROB.
func (rob *ROB) AddCritical(field ROBErrorField, message ROBMessage) {
	rob.Errs = append(rob.Errs, &ROBError{
		Type:    ROBTYPE_CRITICAL,
		Message: message,
		Field:   field,
	})
}

// AddCriticalf adds a formatted critical error message to the ROB.
func (rob *ROB) AddCriticalf(field ROBErrorField, message ROBMessage, v ...interface{}) {
	rob.Errs = append(rob.Errs, &ROBError{
		Type:    ROBTYPE_CRITICAL,
		Message: ROBMessage(fmt.Sprintf(message.String(), v...)),
		Field:   field,
	})
}

// AddError adds an error message to the ROB.
func (rob *ROB) AddError(field ROBErrorField, message ROBMessage) {
	rob.Errs = append(rob.Errs, &ROBError{
		Type:    ROBTYPE_ERROR,
		Message: message,
		Field:   field,
	})
}

// AddErrorf adds a formatted error message to the ROB.
func (rob *ROB) AddErrorf(field ROBErrorField, message ROBMessage, v ...interface{}) {
	rob.Errs = append(rob.Errs, &ROBError{
		Type:    ROBTYPE_ERROR,
		Message: ROBMessage(fmt.Sprintf(message.String(), v...)),
		Field:   field,
	})
}

// AddWarning adds a warning message to the ROB.
func (rob *ROB) AddWarning(field ROBErrorField, message ROBMessage) {
	rob.Errs = append(rob.Errs, &ROBError{
		Type:    ROBTYPE_WARNING,
		Message: message,
		Field:   field,
	})
}

// AddWarningf adds a formatted warning message to the ROB.
func (rob *ROB) AddWarningf(field ROBErrorField, message ROBMessage, v ...interface{}) {
	rob.Errs = append(rob.Errs, &ROBError{
		Type:    ROBTYPE_WARNING,
		Message: ROBMessage(fmt.Sprintf(message.String(), v...)),
		Field:   field,
	})
}

// AddNotice adds a notice message to the ROB.
func (rob *ROB) AddNotice(field ROBErrorField, message ROBMessage) {
	rob.Errs = append(rob.Errs, &ROBError{
		Type:    ROBTYPE_NOTICE,
		Message: message,
		Field:   field,
	})
}

// AddNoticef adds a formatted notice message to the ROB.
func (rob *ROB) AddNoticef(field ROBErrorField, message ROBMessage, v ...interface{}) {
	rob.Errs = append(rob.Errs, &ROBError{
		Type:    ROBTYPE_NOTICE,
		Message: ROBMessage(fmt.Sprintf(message.String(), v...)),
		Field:   field,
	})
}

// AddInfo adds an informational message to the ROB.
func (rob *ROB) AddInfo(field ROBErrorField, message ROBMessage) {
	rob.Errs = append(rob.Errs, &ROBError{
		Type:    ROBTYPE_INFO,
		Message: message,
		Field:   field,
	})
}

// AddInfof adds a formatted informational message to the ROB.
func (rob *ROB) AddInfof(field ROBErrorField, message ROBMessage, v ...interface{}) {
	rob.Errs = append(rob.Errs, &ROBError{
		Type:    ROBTYPE_INFO,
		Message: ROBMessage(fmt.Sprintf(message.String(), v...)),
		Field:   field,
	})
}

// AddDebug adds a debug message to the ROB.
func (rob *ROB) AddDebug(field ROBErrorField, message ROBMessage) {
	rob.Errs = append(rob.Errs, &ROBError{
		Type:    ROBTYPE_DEBUG,
		Message: message,
		Field:   field,
	})
}

// AddDebugf adds a formatted debug message to the ROB.
func (rob *ROB) AddDebugf(field ROBErrorField, message ROBMessage, v ...interface{}) {
	rob.Errs = append(rob.Errs, &ROBError{
		Type:    ROBTYPE_DEBUG,
		Message: ROBMessage(fmt.Sprintf(message.String(), v...)),
		Field:   field,
	})
}
