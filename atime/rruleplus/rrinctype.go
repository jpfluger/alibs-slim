package rruleplus

import "strings"

const (
	RRULE_INC_TYPE_INCLUSIVE RRIncType = "inclusive"
	RRULE_INC_TYPE_EXCLUSIVE RRIncType = "exclusive"
)

type RRIncType string

// IsEmpty returns true if the RRIncType is not set.
func (t RRIncType) IsEmpty() bool {
	return string(t) == ""
}

// String returns the lowercase string value of the RRIncType.
func (t RRIncType) String() string {
	return strings.ToLower(string(t))
}
