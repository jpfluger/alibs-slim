package rruleplus

import "strings"

type ObservanceMode string

const (
	ObservanceNone           ObservanceMode = "none"
	ObservanceNextBizDay     ObservanceMode = "next-business-day"
	ObservancePreviousBizDay ObservanceMode = "previous-business-day"
)

// IsEmpty returns true if the ObservanceMode is empty or only whitespace.
func (om ObservanceMode) IsEmpty() bool {
	return strings.TrimSpace(string(om)) == ""
}

// TrimSpace trims whitespace from the ObservanceMode string.
func (om ObservanceMode) TrimSpace() ObservanceMode {
	return ObservanceMode(strings.TrimSpace(string(om)))
}
