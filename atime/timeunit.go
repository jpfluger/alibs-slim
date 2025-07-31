package atime

import (
	"github.com/teambition/rrule-go"
	"strings"
	"time"
)

// TimeUnit defines the unit of recurrence for time ranges.
type TimeUnit string

const (
	TIMEUNIT_SECONDLY TimeUnit = "secondly"
	TIMEUNIT_MINUTELY TimeUnit = "minutely"
	TIMEUNIT_HOURLY   TimeUnit = "hourly"
	TIMEUNIT_DAILY    TimeUnit = "daily"
	TIMEUNIT_WEEKLY   TimeUnit = "weekly"
	TIMEUNIT_MONTHLY  TimeUnit = "monthly"
	TIMEUNIT_YEARLY   TimeUnit = "yearly"
)

// IsEmpty returns true if the TimeUnit is not set.
func (t TimeUnit) IsEmpty() bool {
	return string(t) == ""
}

// String returns the lowercase string value of the TimeUnit.
func (t TimeUnit) String() string {
	return strings.ToLower(string(t))
}

// IsValid checks if the TimeUnit is one of the supported values.
func (t TimeUnit) IsValid() bool {
	switch t {
	case TIMEUNIT_SECONDLY,
		TIMEUNIT_MINUTELY,
		TIMEUNIT_HOURLY,
		TIMEUNIT_DAILY,
		TIMEUNIT_WEEKLY,
		TIMEUNIT_MONTHLY,
		TIMEUNIT_YEARLY:
		return true
	default:
		return false
	}
}

func (t TimeUnit) CalcDuration(duration int) time.Duration {
	if t.IsEmpty() {
		return 0
	}
	if duration <= 0 {
		duration = 1
	}
	switch t {
	case TIMEUNIT_SECONDLY:
		return time.Duration(duration) * time.Second
	case TIMEUNIT_MINUTELY:
		return time.Duration(duration) * time.Minute
	case TIMEUNIT_HOURLY:
		return time.Duration(duration) * time.Hour
	case TIMEUNIT_DAILY:
		return time.Duration(duration) * 24 * time.Hour
	case TIMEUNIT_WEEKLY:
		return time.Duration(duration) * 7 * 24 * time.Hour
	case TIMEUNIT_MONTHLY, TIMEUNIT_YEARLY:
		// Variable-length durations — caller must handle these differently
		// You could also panic or return error depending on usage context
		return 0
	default:
		return 0
	}
}

func (t TimeUnit) IsSameInterval(last, now time.Time) bool {
	switch t {
	case TIMEUNIT_SECONDLY:
		return last.Equal(now)
	case TIMEUNIT_MINUTELY:
		return last.Year() == now.Year() &&
			last.Month() == now.Month() &&
			last.Day() == now.Day() &&
			last.Hour() == now.Hour() &&
			last.Minute() == now.Minute()
	case TIMEUNIT_HOURLY:
		return last.Year() == now.Year() &&
			last.Month() == now.Month() &&
			last.Day() == now.Day() &&
			last.Hour() == now.Hour()
	case TIMEUNIT_DAILY:
		return last.Year() == now.Year() &&
			last.Month() == now.Month() &&
			last.Day() == now.Day()
	case TIMEUNIT_WEEKLY:
		ay, aw := last.ISOWeek()
		by, bw := now.ISOWeek()
		return ay == by && aw == bw
	case TIMEUNIT_MONTHLY:
		return last.Year() == now.Year() &&
			last.Month() == now.Month()
	case TIMEUNIT_YEARLY:
		return last.Year() == now.Year()
	default:
		return false
	}
}

// Default returns the default TimeUnit (daily) if empty.
func (t TimeUnit) Default() TimeUnit {
	if t.IsEmpty() {
		return TIMEUNIT_DAILY
	}
	return t
}

// ToFrequency converts the custom TimeUnit to the corresponding rrule.Frequency value.
func (t TimeUnit) ToFrequency() rrule.Frequency {
	switch t.Default() {
	case TIMEUNIT_SECONDLY:
		return rrule.SECONDLY
	case TIMEUNIT_MINUTELY:
		return rrule.MINUTELY
	case TIMEUNIT_HOURLY:
		return rrule.HOURLY
	case TIMEUNIT_DAILY:
		return rrule.DAILY
	case TIMEUNIT_WEEKLY:
		return rrule.WEEKLY
	case TIMEUNIT_MONTHLY:
		return rrule.MONTHLY
	case TIMEUNIT_YEARLY:
		return rrule.YEARLY
	default:
		return rrule.DAILY
	}
}

// FromFrequency converts an rrule.Frequency constant to a corresponding TimeUnit.
func FromFrequency(freq int) TimeUnit {
	switch rrule.Frequency(freq) {
	case rrule.SECONDLY:
		return TIMEUNIT_SECONDLY
	case rrule.MINUTELY:
		return TIMEUNIT_MINUTELY
	case rrule.HOURLY:
		return TIMEUNIT_HOURLY
	case rrule.DAILY:
		return TIMEUNIT_DAILY
	case rrule.WEEKLY:
		return TIMEUNIT_WEEKLY
	case rrule.MONTHLY:
		return TIMEUNIT_MONTHLY
	case rrule.YEARLY:
		return TIMEUNIT_YEARLY
	default:
		return TIMEUNIT_DAILY
	}
}
