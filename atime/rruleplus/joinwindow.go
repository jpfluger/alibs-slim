package rruleplus

import (
	"fmt"
	"github.com/jpfluger/alibs-slim/atime"
	"strings"
	"time"
)

// JoinWindow defines a pre- or post-occurrence window relative to a scheduled event.
// It is used for proactive matching like reminders, alerts, or notifications.
type JoinWindow struct {
	IsBefore     bool           `json:"isBefore,omitempty"` // If true, the window occurs before the event; if false, it's after.
	Duration     int            `json:"duration,omitempty"` // How long the window lasts, expressed in DurationUnit.
	DurationUnit atime.TimeUnit `json:"durationUnit"`       // The time unit for the duration (e.g., daily, weekly).
	Label        string         `json:"label,omitempty"`    // Optional human-readable description.
	Tag          string         `json:"tag,omitempty"`      // Identifier returned on match (used to tag outcomes).
}

func (jw *JoinWindow) Validate() error {
	if jw == nil {
		return fmt.Errorf("join window is nil")
	}
	if jw.DurationUnit.IsEmpty() {
		return fmt.Errorf("join window duration unit is required")
	}
	if jw.Duration < 0 {
		jw.Duration = 0
	}
	jw.Label = strings.TrimSpace(jw.Label)
	jw.Tag = strings.TrimSpace(jw.Tag)
	return nil
}

func (jw *JoinWindow) Matches(now, occurrence time.Time) bool {
	if jw == nil || jw.Duration == 0 || jw.DurationUnit == "" {
		return false
	}
	if err := jw.Validate(); err != nil {
		return false
	}

	start, end := jw.WindowRange(occurrence)
	now = now.UTC()
	return (now.Equal(start) || now.After(start)) && now.Before(end)
}

//func (jw *JoinWindow) Matches(now, occurrence time.Time) bool {
//	if jw == nil {
//		return false
//	}
//	if err := jw.Validate(); err != nil {
//		return false
//	}
//
//	now = now.UTC()
//	occurrence = occurrence.UTC()
//
//	var start, end time.Time
//
//	switch jw.DurationUnit {
//	case TIMEUNIT_SECONDLY, TIMEUNIT_MINUTELY, TIMEUNIT_HOURLY, TIMEUNIT_DAILY, TIMEUNIT_WEEKLY:
//		dur := jw.DurationUnit.CalcDuration(jw.Duration)
//		if dur == 0 {
//			return false
//		}
//		if jw.IsBefore {
//			start = occurrence.Add(-dur)
//			end = occurrence
//		} else {
//			start = occurrence
//			end = occurrence.Add(dur)
//		}
//	case TIMEUNIT_MONTHLY:
//		d := jw.Duration
//		if d <= 0 {
//			d = 1
//		}
//		if jw.IsBefore {
//			start = occurrence.AddDate(0, -d, 0)
//			end = occurrence
//		} else {
//			start = occurrence
//			end = occurrence.AddDate(0, d, 0)
//		}
//	case TIMEUNIT_YEARLY:
//		d := jw.Duration
//		if d <= 0 {
//			d = 1
//		}
//		if jw.IsBefore {
//			start = occurrence.AddDate(-d, 0, 0)
//			end = occurrence
//		} else {
//			start = occurrence
//			end = occurrence.AddDate(d, 0, 0)
//		}
//	default:
//		return false
//	}
//
//	return (now.Equal(start) || now.After(start)) && now.Before(end)
//}

func (jw *JoinWindow) WindowRange(occurrence time.Time) (time.Time, time.Time) {
	if jw == nil || jw.Duration == 0 {
		return time.Time{}, time.Time{}
	}

	occurrence = occurrence.UTC()

	switch jw.DurationUnit {
	case atime.TIMEUNIT_SECONDLY, atime.TIMEUNIT_MINUTELY, atime.TIMEUNIT_HOURLY, atime.TIMEUNIT_DAILY, atime.TIMEUNIT_WEEKLY:
		dur := jw.DurationUnit.CalcDuration(jw.Duration)
		if jw.IsBefore {
			return occurrence.Add(-dur), occurrence
		}
		return occurrence, occurrence.Add(dur)

	case atime.TIMEUNIT_MONTHLY:
		if jw.IsBefore {
			return occurrence.AddDate(0, -jw.Duration, 0), occurrence
		}
		return occurrence, occurrence.AddDate(0, jw.Duration, 0)

	case atime.TIMEUNIT_YEARLY:
		if jw.IsBefore {
			return occurrence.AddDate(-jw.Duration, 0, 0), occurrence
		}
		return occurrence, occurrence.AddDate(jw.Duration, 0, 0)

	default:
		return time.Time{}, time.Time{}
	}
}

func (jw *JoinWindow) Clone() *JoinWindow {
	if jw == nil {
		return nil
	}
	clone := *jw // shallow copy is sufficient for value types
	return &clone
}

// JoinWindows is a collection of JoinWindow entries.
// It enables multiple lead/lag evaluation windows for a single recurrence rule.
type JoinWindows []*JoinWindow

func (jws JoinWindows) Clone() JoinWindows {
	if jws == nil {
		return nil
	}
	clone := make(JoinWindows, len(jws))
	for i, jw := range jws {
		if jw != nil {
			clone[i] = jw.Clone()
		}
	}
	return clone
}

func (jws JoinWindows) Validate() error {
	var errs []string
	for i, jw := range jws {
		if err := jw.Validate(); err != nil {
			errs = append(errs, fmt.Sprintf("JoinWindow[%d]: %v", i, err))
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("validation failed:\n%s", strings.Join(errs, "\n"))
	}
	return nil
}

func (jws JoinWindows) Sanitize() JoinWindows {
	if jws == nil {
		return nil
	}
	var sanitized JoinWindows
	for _, jw := range jws {
		if err := jw.Validate(); err == nil {
			sanitized = append(sanitized, jw)
		}
	}
	return sanitized
}

func (jws JoinWindows) MustSanitize() JoinWindows {
	sanitized := jws.Sanitize()
	if sanitized == nil {
		return JoinWindows{}
	}
	return sanitized
}

func (jws JoinWindows) Matches(now, occurrence time.Time) *JoinWindow {
	var matched *JoinWindow
	var smallest time.Duration

	now = now.UTC()

	for _, jw := range jws {
		if jw == nil || jw.Duration == 0 {
			continue
		}
		start, end := jw.WindowRange(occurrence)
		if (now.Equal(start) || now.After(start)) && now.Before(end) {
			dur := end.Sub(start)
			if matched == nil || dur < smallest {
				matched = jw
				smallest = dur
			}
		}
	}
	return matched
}
