package atime

import (
	"github.com/dustin/go-humanize"
	"github.com/teambition/rrule-go"
	"sort"
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

// EnsureDateTimeUTC ensures that the input is a time.Time object in UTC.
// If the input is nil or not a time.Time, it returns the zero value of time.Time.
func EnsureDateTimeUTC(t interface{}) time.Time {
	switch v := t.(type) {
	case *time.Time:
		if v != nil {
			return (*v).UTC()
		}
	case time.Time:
		return v.UTC()
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

// MustParseRFC3339 parses a time string according to the RFC3339.
// If parsing fails, it returns the zero value of time.Time.
func MustParseRFC3339(value string) time.Time {
	return MustParse(time.RFC3339, value)
}

// MustParsePtrRFC3339 parses a time string according to RFC3339 and returns a pointer to the result.
// If parsing fails, it returns a pointer to the zero value of time.Time.
func MustParsePtrRFC3339(value string) *time.Time {
	return MustParsePtr(time.RFC3339, value)
}

// CurrentYear returns the current year.
func CurrentYear() int {
	return time.Now().Year()
}

// CurrentYearUTC returns the current year in UTC.
func CurrentYearUTC() int {
	return time.Now().UTC().Year()
}

// SanitizePtr ensures that a *time.Time pointer is either nil or a meaningful value.
// If the input pointer is nil or points to a zero-value time (`time.Time{}`),
// it returns nil. Otherwise, it returns the original pointer.
//
// This is useful for database serialization, JSON marshalling, or nullable semantics
// where `nil` should be used instead of meaningless zero-time values.
func SanitizePtr(target *time.Time) *time.Time {
	if target == nil || target.IsZero() {
		return nil
	}
	return target
}

// HoursBetween returns a slice of hour integers (0–23) between `start` and `end`,
// wrapping around midnight if needed. If isOrdered is true, results are sorted ascending.
// If start == end, returns empty list (no duration).
func HoursBetween(start, end int, isOrdered bool) []int {
	if start == end {
		return []int{}
	}

	var hours []int
	if isOrdered {
		hourSet := map[int]struct{}{}
		if start < end {
			for h := start; h < end; h++ {
				hourSet[h] = struct{}{}
			}
		} else {
			for h := start; h < 24; h++ {
				hourSet[h] = struct{}{}
			}
			for h := 0; h < end; h++ {
				hourSet[h] = struct{}{}
			}
		}
		for h := range hourSet {
			hours = append(hours, h)
		}
		sort.Ints(hours)
	} else {
		if start < end {
			for h := start; h < end; h++ {
				hours = append(hours, h)
			}
		} else {
			for h := start; h < 24; h++ {
				hours = append(hours, h)
			}
			for h := 0; h < end; h++ {
				hours = append(hours, h)
			}
		}
	}

	return hours
}

// HoursBetweenDates returns a list of unique hour values (0–23) between two timestamps.
// If isOrdered is true, the result is sorted in ascending order.
func HoursBetweenDates(start, end time.Time, isOrdered bool) []int {
	if !end.After(start) {
		return []int{}
	}

	hourSeen := make(map[int]struct{})
	var hours []int
	current := start.Truncate(time.Hour)

	for current.Before(end) {
		h := current.Hour()
		if _, exists := hourSeen[h]; !exists {
			hourSeen[h] = struct{}{}
			hours = append(hours, h)
		}
		current = current.Add(time.Hour)
	}

	if isOrdered {
		sort.Ints(hours)
	}

	return hours
}

// orderedHourKeys converts a map[int]struct{} to a sorted slice of int keys.
func orderedHourKeys(hourSet map[int]struct{}) []int {
	var hours []int
	for h := range hourSet {
		hours = append(hours, h)
	}
	sort.Ints(hours)
	return hours
}

// TimeWeekdayToRRuleWeekday converts a single time.Weekday to its corresponding rrule.Weekday
func TimeWeekdayToRRuleWeekday(d time.Weekday) rrule.Weekday {
	switch d {
	case time.Sunday:
		return rrule.SU
	case time.Monday:
		return rrule.MO
	case time.Tuesday:
		return rrule.TU
	case time.Wednesday:
		return rrule.WE
	case time.Thursday:
		return rrule.TH
	case time.Friday:
		return rrule.FR
	case time.Saturday:
		return rrule.SA
	default:
		return rrule.MO // fallback to Monday... otherwise `panic("invalid time.Weekday value")`
	}
}

// RRuleWeekdayToTimeWeekday converts a rrule.Weekday to time.Weekday
func RRuleWeekdayToTimeWeekday(d rrule.Weekday) time.Weekday {
	switch d.Day() {
	case 0:
		return time.Monday
	case 1:
		return time.Tuesday
	case 2:
		return time.Wednesday
	case 3:
		return time.Thursday
	case 4:
		return time.Friday
	case 5:
		return time.Saturday
	case 6:
		return time.Sunday
	default:
		panic("invalid rrule.Weekday value")
	}
}

// TimeWeekdaysToInts converts a variadic slice of time.Weekday into []int (0 = Sunday, ..., 6 = Saturday)
func TimeWeekdaysToInts(days ...time.Weekday) []int {
	result := make([]int, len(days))
	for i, d := range days {
		result[i] = int(d)
	}
	return result
}

// TimeWeekdaysToRRuleWeekdays converts a variadic slice of time.Weekday into []rrule.Weekday
func TimeWeekdaysToRRuleWeekdays(days ...time.Weekday) []rrule.Weekday {
	result := make([]rrule.Weekday, len(days))
	for i, d := range days {
		result[i] = TimeWeekdayToRRuleWeekday(d)
	}
	return result
}

// RRuleWeekdaysToInts converts a variadic slice of rrule.Weekday into []int (0 = Sunday, ..., 6 = Saturday)
func RRuleWeekdaysToInts(days ...rrule.Weekday) []int {
	result := make([]int, len(days))
	for i, d := range days {
		result[i] = d.Day()
	}
	return result
}

// RRuleWeekdaysToTimeWeekdays converts a slice of rrule.Weekday to []time.Weekday
func RRuleWeekdaysToTimeWeekdays(rrdays ...rrule.Weekday) []time.Weekday {
	result := make([]time.Weekday, len(rrdays))
	for i, d := range rrdays {
		result[i] = RRuleWeekdayToTimeWeekday(d)
	}
	return result
}

// IsWeekendByTime returns a true if the date falls on a Saturday or Sunday.
func IsWeekendByTime(t time.Time) bool {
	return t.Weekday() == time.Saturday || t.Weekday() == time.Sunday
}
