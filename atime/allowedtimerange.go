package atime

import (
	"fmt"
	"time"
)

type AllowedTimeRange struct {
	Start           time.Time      `json:"start,omitempty" jsonschema:"type=string,format=date-time,title=Start Time,description=The starting time of the range"`
	End             time.Time      `json:"end,omitempty" jsonschema:"type=string,format=date-time,title=End Time,description=The ending time of the range"`
	Unit            TimeUnit       `json:"unit,omitempty" jsonschema:"type=string,title=Recurrence Unit,description=The recurrence unit (e.g., hourly, daily)"`
	IncludeWeekdays bool           `json:"includeWeekdays,omitempty" jsonschema:"type=boolean,title=Include Weekdays,description=Whether weekdays are included"`
	IncludeWeekends bool           `json:"includeWeekends,omitempty" jsonschema:"type=boolean,title=Include Weekends,description=Whether weekends are included"`
	SpecificDays    []time.Weekday `json:"specificDays,omitempty" jsonschema:"type=array,items=integer,title=Specific Days,description=Specific weekdays allowed (e.g., Sunday=0, Monday=1)"`
	SpecificDates   []time.Time    `json:"specificDates,omitempty" jsonschema:"type=array,items=string,format=date-time,title=Specific Dates,description=Specific calendar dates allowed"`
	RecurringMonths []time.Month   `json:"recurringMonths,omitempty" jsonschema:"type=array,items=integer,title=Recurring Months,description=Specific months allowed (e.g., January=1, February=2)"`
	ExcludeDates    []time.Time    `json:"excludeDates,omitempty" jsonschema:"type=array,items=string,format=date-time,title=Excluded Dates,description=Specific calendar dates to exclude"`
	TimeZone        *time.Location `json:"timeZone,omitempty" jsonschema:"type=string,title=Time Zone,description=Time zone for the range"`
}

func (atr *AllowedTimeRange) ConvertToUTC() {
	if atr.TimeZone == nil {
		atr.TimeZone = time.UTC
	}

	atr.Start = atr.Start.In(atr.TimeZone).UTC()
	atr.End = atr.End.In(atr.TimeZone).UTC()

	for i, date := range atr.SpecificDates {
		atr.SpecificDates[i] = date.In(atr.TimeZone).UTC()
	}

	for i, date := range atr.ExcludeDates {
		atr.ExcludeDates[i] = date.In(atr.TimeZone).UTC()
	}
}

func (atr *AllowedTimeRange) GetLocation() *time.Location {
	if atr.TimeZone != nil {
		return atr.TimeZone
	}
	// Default to UTC if no timezone is specified
	return time.UTC
}

func (atr *AllowedTimeRange) SetTimeZone(location *time.Location) {
	if location == nil {
		location = time.UTC
	}
	atr.TimeZone = location
	atr.ConvertToUTC()
}

// FNTimeNow is a function type that returns the current time.
type FNTimeNow func() time.Time

// IsNowAllowed determines if the current time is within the allowed time range.
func (atr *AllowedTimeRange) IsNowAllowed() bool {
	return atr.IsNowAllowedWithTime(time.Now)
}

// IsNowAllowedWithTime determines if a given "now" function satisfies the allowed time range.
func (atr *AllowedTimeRange) IsNowAllowedWithTime(fn FNTimeNow) bool {
	now := fn().In(atr.GetLocation())

	// Exclude specific dates
	for _, excludeDate := range atr.ExcludeDates {
		excludeDate = excludeDate.In(atr.GetLocation()).Truncate(24 * time.Hour)
		if now.Truncate(24 * time.Hour).Equal(excludeDate) {
			return false
		}
	}

	// Check weekdays and weekends
	weekday := now.Weekday()
	if atr.IncludeWeekdays && weekday >= time.Monday && weekday <= time.Friday {
		return true
	}
	if atr.IncludeWeekends && (weekday == time.Saturday || weekday == time.Sunday) {
		return true
	}

	// Additional checks for SpecificDays
	for _, day := range atr.SpecificDays {
		if day == weekday {
			return true
		}
	}

	// Additional checks for SpecificDates
	for _, date := range atr.SpecificDates {
		date = date.In(atr.GetLocation()).Truncate(24 * time.Hour)
		if now.Truncate(24 * time.Hour).Equal(date) {
			return true
		}
	}

	return false
}

func (atr *AllowedTimeRange) Validate() error {
	if atr.End.Before(atr.Start) {
		return fmt.Errorf("end time must be after start time")
	}

	if atr.TimeZone == nil {
		atr.TimeZone = time.UTC
	}

	if len(atr.SpecificDays) == 0 && len(atr.SpecificDates) == 0 && len(atr.RecurringMonths) == 0 && !atr.IncludeWeekdays && !atr.IncludeWeekends {
		return fmt.Errorf("at least one condition must be set for allowed time ranges")
	}

	return nil
}

func (atr *AllowedTimeRange) Normalize() {
	atr.ConvertToUTC()
	if atr.TimeZone == nil {
		atr.TimeZone = time.UTC
	}
}

func (atr *AllowedTimeRange) IsAllowedAt(t time.Time) bool {
	t = t.In(atr.GetLocation()) // Ensure the time is converted to the proper time zone

	// Check if the specific date is excluded
	for _, exclude := range atr.ExcludeDates {
		if t.Year() == exclude.Year() && t.Month() == exclude.Month() && t.Day() == exclude.Day() {
			return false
		}
	}

	// Check if the specific date is allowed
	if len(atr.SpecificDates) > 0 {
		for _, specific := range atr.SpecificDates {
			if t.Year() == specific.Year() && t.Month() == specific.Month() && t.Day() == specific.Day() {
				return true
			}
		}
		return false // If specific dates are defined, time must match one of them
	}

	// Check weekday inclusion
	weekday := t.Weekday()
	if atr.IncludeWeekdays && weekday >= time.Monday && weekday <= time.Friday {
		// Weekday is allowed
	} else if atr.IncludeWeekends && (weekday == time.Saturday || weekday == time.Sunday) {
		// Weekend is allowed
	} else if len(atr.SpecificDays) > 0 {
		for _, specificDay := range atr.SpecificDays {
			if specificDay == weekday {
				return true
			}
		}
		return false // Weekday doesn't match specific allowed days
	} else {
		return false // Neither weekday nor weekend is allowed
	}

	// Check if the time falls within the start and end times
	start := time.Date(0, 1, 1, atr.Start.Hour(), atr.Start.Minute(), atr.Start.Second(), 0, atr.Start.Location())
	end := time.Date(0, 1, 1, atr.End.Hour(), atr.End.Minute(), atr.End.Second(), 0, atr.End.Location())
	now := time.Date(0, 1, 1, t.Hour(), t.Minute(), t.Second(), 0, t.Location())

	if !now.Before(start) && !now.After(end) {
		return true
	}

	return false
}

type AllowedTimeRanges []*AllowedTimeRange

func (atrs AllowedTimeRanges) IsNowAllowed(fnTimeNow FNTimeNow) bool {
	for _, atr := range atrs {
		if atr.IsNowAllowedWithTime(fnTimeNow) {
			return true
		}
	}
	return false
}

func (atrs AllowedTimeRanges) IsAllowedAt(t time.Time) bool {
	for _, atr := range atrs {
		if atr.IsAllowedAt(t) {
			return true
		}
	}
	return false
}

func (atrs AllowedTimeRanges) Validate() error {
	for i, atr := range atrs {
		if err := atr.Validate(); err != nil {
			return fmt.Errorf("Time range at index %d failed validation: %v", i, err)
		}
	}
	return nil
}
