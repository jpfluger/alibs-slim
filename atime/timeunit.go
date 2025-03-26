package atime

import "strings"

const (
	TIMEUNIT_HOURLY  TimeUnit = "hourly"
	TIMEUNIT_DAILY   TimeUnit = "daily"
	TIMEUNIT_WEEKLY  TimeUnit = "weekly"
	TIMEUNIT_MONTHLY TimeUnit = "monthly"
	TIMEUNIT_YEARLY  TimeUnit = "yearly"
)

// TimeUnit defines the unit of recurrence for the time ranges.
type TimeUnit string

func (t TimeUnit) IsEmpty() bool { return string(t) == "" }

func (t TimeUnit) String() string {
	return strings.ToLower(string(t))
}

func (t TimeUnit) IsValid() bool {
	switch t {
	case TIMEUNIT_HOURLY, TIMEUNIT_DAILY, TIMEUNIT_WEEKLY, TIMEUNIT_MONTHLY, TIMEUNIT_YEARLY:
		return true
	default:
		return false
	}
}

func (t TimeUnit) Default() TimeUnit {
	if t.IsEmpty() {
		return TIMEUNIT_DAILY // Assume daily is the default
	}
	return t
}
