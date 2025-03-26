package atime

import (
	"github.com/dustin/go-humanize"
	"time"
)

// EnsureDateTime ensures that the input is a time.Time object.
// If the input is nil or not a time.Time, it returns the zero value of time.Time.
func EnsureDateTime(t interface{}) time.Time {
	switch v := t.(type) {
	case *time.Time:
		if v != nil {
			return *v
		}
	case time.Time:
		return v
	}
	return time.Time{}
}

// ToPointer converts a time.Time to a pointer to time.Time.
func ToPointer(t time.Time) *time.Time {
	return &t
}

// ToPointerNil converts a time.Time to a pointer to time.Time, returning nil if the time is zero.
func ToPointerNil(t time.Time) *time.Time {
	if t.IsZero() {
		return nil
	}
	return &t
}

// GetNowPointer returns a pointer to the current time.
func GetNowPointer() *time.Time {
	ts := time.Now()
	return &ts
}

// GetNowUTCPointer returns a pointer to the current UTC time.
func GetNowUTCPointer() *time.Time {
	ts := time.Now().UTC()
	return &ts
}

// FormatDateTime formats a time.Time object according to the provided layout.
// If the input is not a valid time.Time, it returns an empty string.
func FormatDateTime(t interface{}, layout string) string {
	dt := EnsureDateTime(t)
	if dt.IsZero() {
		return ""
	}
	return dt.Format(layout)
}

// FormatDateTimeRFC3339 formats a time.Time object according to RFC3339 layout.
func FormatDateTimeRFC3339(t interface{}) string {
	return FormatDateTime(t, time.RFC3339)
}

// FormatDateTimeRFC3339Nano formats a time.Time object according to RFC3339Nano layout.
func FormatDateTimeRFC3339Nano(t interface{}) string {
	return FormatDateTime(t, time.RFC3339Nano)
}

// FormatDateTimeRFC1123 formats a time.Time object according to RFC1123 layout.
func FormatDateTimeRFC1123(t interface{}) string {
	return FormatDateTime(t, time.RFC1123)
}

// FormatDateTimeRFC1123Z formats a time.Time object according to RFC1123Z layout.
func FormatDateTimeRFC1123Z(t interface{}) string {
	return FormatDateTime(t, time.RFC1123Z)
}

// FormatDateTimeElse formats a time.Time object according to the provided layout,
// returning elseThisString if the time is zero.
func FormatDateTimeElse(t interface{}, layout string, elseThisString string) string {
	dt := EnsureDateTime(t)
	if dt.IsZero() {
		return elseThisString
	}
	return dt.Format(layout)
}

// IsDateBeforeNow checks if the given time is before the current time.
func IsDateBeforeNow(u interface{}) bool {
	dt := EnsureDateTime(u)
	return !dt.IsZero() && time.Now().UTC().After(dt.UTC())
}

// IsDateAfterNow checks if the given time is after the current time.
func IsDateAfterNow(u interface{}) bool {
	dt := EnsureDateTime(u)
	return !dt.IsZero() && time.Now().UTC().Before(dt.UTC())
}

// IsDateBefore checks if the first time is before the second time.
func IsDateBefore(a interface{}, b interface{}) bool {
	dtA := EnsureDateTime(a)
	dtB := EnsureDateTime(b)
	return !dtA.IsZero() && dtA.Before(dtB)
}

// IsDateAfter checks if the first time is after the second time.
func IsDateAfter(a interface{}, b interface{}) bool {
	dtA := EnsureDateTime(a)
	dtB := EnsureDateTime(b)
	return !dtA.IsZero() && dtA.After(dtB)
}

// FormatDateTimeRelative formats two times relative to each other using humanize.RelTime.
func FormatDateTimeRelative(a, b interface{}, albl, blbl string) string {
	dt1 := EnsureDateTime(a)
	dt2 := EnsureDateTime(b)
	if dt1.IsZero() || dt2.IsZero() {
		return ""
	}
	return humanize.RelTime(dt1, dt2, albl, blbl)
}

// FormatDateTimeAgo formats a time relative to the current time using humanize.Time.
func FormatDateTimeAgo(then interface{}) string {
	dt := EnsureDateTime(then)
	if dt.IsZero() {
		return ""
	}
	return humanize.Time(dt)
}

// IfDateEmptyElse returns ifTrue if the target time is zero, otherwise returns ifFalse.
func IfDateEmptyElse(target interface{}, ifTrue, ifFalse string) string {
	dt := EnsureDateTime(target)
	if dt.IsZero() {
		return ifTrue
	}
	return ifFalse
}

// MustParse parses a time string according to the provided layout.
// If parsing fails, it returns the zero value of time.Time.
func MustParse(layout string, value string) time.Time {
	t, err := time.Parse(layout, value)
	if err != nil {
		return time.Time{}
	}
	return t
}

// MustParsePtr parses a time string according to the provided layout and returns a pointer to the result.
// If parsing fails, it returns a pointer to the zero value of time.Time.
func MustParsePtr(layout string, value string) *time.Time {
	t, err := time.Parse(layout, value)
	if err != nil {
		return &time.Time{}
	}
	return &t
}

// CurrentYear returns the current year.
func CurrentYear() int {
	return time.Now().Year()
}

// CurrentYearUTC returns the current year in UTC.
func CurrentYearUTC() int {
	return time.Now().UTC().Year()
}
