package aconns

import (
	"strings"
)

// RecordIdentifierType represents the type of action performed on a database record.
type RecordIdentifierType string

// Constants for RecordIdentifierType values.
const (
	REC_ID_USER  RecordIdentifierType = "user"
	REC_ID_EMAIL RecordIdentifierType = "email"
	REC_ID_PHONE RecordIdentifierType = "phone"
)

var validRecordIdentifierTypes = RecordIdentifierTypes{
	REC_ID_USER,
	REC_ID_EMAIL,
	REC_ID_PHONE,
}

func RECIDENTIFIERS() RecordIdentifierTypes {
	return validRecordIdentifierTypes
}

// IsEmpty checks if the RecordIdentifierType is empty after trimming whitespace.
func (rt RecordIdentifierType) IsEmpty() bool {
	return strings.TrimSpace(string(rt)) == ""
}

// TrimSpace returns a new RecordIdentifierType with leading and trailing whitespace removed.
func (rt RecordIdentifierType) TrimSpace() RecordIdentifierType {
	return RecordIdentifierType(strings.TrimSpace(string(rt)))
}

// String converts the RecordIdentifierType to a regular string.
func (rt RecordIdentifierType) String() string {
	return strings.TrimSpace(string(rt))
}

type RecordIdentifierTypes []RecordIdentifierType

func (rts RecordIdentifierTypes) HasMatch(recIdType RecordIdentifierType) bool {
	if rts == nil || len(rts) == 0 {
		return false
	}
	for _, rt := range rts {
		if rt == recIdType {
			return true
		}
	}
	return false
}

// IsValid checks if the RecordIdentifierType is one of the allowed constants.
func (rt RecordIdentifierType) IsValid() bool {
	return RECIDENTIFIERS().HasMatch(rt)
}
