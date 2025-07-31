package rruleplus

import (
	"github.com/jpfluger/alibs-slim/atime"
	"github.com/rickar/cal/v2"
	"github.com/rickar/cal/v2/us"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/teambition/rrule-go"
	"testing"
	"time"
)

func TestRRulePlus_ShiftWeekend(t *testing.T) {
	rule, err := NewRRulePlus(ROptionPlus{
		ROption: rrule.ROption{
			Freq:     rrule.DAILY,
			Interval: 1,
			Dtstart:  atime.MustParseRFC3339("2025-06-21T00:00:00Z"), // Saturday
		},
		ShiftOffWeekend: true,
	})
	assert.NoError(t, err)

	next := rule.After(atime.MustParseRFC3339("2025-06-21T00:00:00Z"), true)
	assert.Equal(t, atime.MustParseRFC3339("2025-06-23T00:00:00Z"), next) // Should shift to Monday
}

func TestRRulePlus_HolidaySkip(t *testing.T) {
	cal := cal.NewBusinessCalendar()
	cal.AddHoliday(us.ThanksgivingDay)

	rule, err := NewRRulePlus(ROptionPlus{
		ROption: rrule.ROption{
			Freq:     rrule.YEARLY,
			Interval: 1,
			Dtstart:  atime.MustParseRFC3339("2025-11-27T00:00:00Z"), // Thanksgiving 2025
		},
		ShiftOffHolidays: true,
		Calendar:         cal,
	})
	assert.NoError(t, err)

	next := rule.After(atime.MustParseRFC3339("2025-11-27T00:00:00Z"), true)
	assert.Equal(t, atime.MustParseRFC3339("2025-11-28T00:00:00Z"), next) // Should move to Friday
}

func TestRRulePlus_Observance(t *testing.T) {
	cal := cal.NewBusinessCalendar()
	cal.AddHoliday(us.NewYear)

	rule, err := NewRRulePlus(ROptionPlus{
		ROption: rrule.ROption{
			Freq:     rrule.YEARLY,
			Interval: 1,
			Dtstart:  atime.MustParseRFC3339("2028-01-01T00:00:00Z"), // New Year's on Saturday
		},
		Observance: ObservanceNextBizDay,
		Calendar:   cal,
	})
	assert.NoError(t, err)

	next := rule.After(atime.MustParseRFC3339("2028-01-01T00:00:00Z"), true)
	assert.Equal(t, atime.MustParseRFC3339("2028-01-03T00:00:00Z"), next) // Skips weekend + Monday observance
}

func TestRRulePlus_CustomFilter(t *testing.T) {
	// Only allow Tuesdays
	rule, err := NewRRulePlus(ROptionPlus{
		ROption: rrule.ROption{
			Freq:     rrule.DAILY,
			Interval: 1,
			Dtstart:  atime.MustParseRFC3339("2025-06-23T00:00:00Z"), // Monday
		},
		CustomFilter: func(t time.Time) bool {
			return t.Weekday() == time.Tuesday
		},
	})
	assert.NoError(t, err)

	next := rule.After(atime.MustParseRFC3339("2025-06-23T00:00:00Z"), true)
	assert.Equal(t, atime.MustParseRFC3339("2025-06-24T00:00:00Z"), next) // Tuesday
}

func TestFiscalCycle_Yearly_WithHolidayObservance(t *testing.T) {
	start := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC) // Jan 1 (New Year's Day, holiday)

	rule, err := NewRRulePlus(ROptionPlus{
		ROption: rrule.ROption{
			Freq:    rrule.YEARLY,
			Dtstart: start,
		},
		ShiftOffHolidays: true,
		ISOCode:          "US",
	})
	assert.NoError(t, err)

	next := rule.After(start.Add(-24*time.Hour), true)
	require.Equal(t, time.Date(2025, 1, 2, 10, 0, 0, 0, time.UTC), next)
}

func TestFiscalCycle_Quarterly_WithHolidayObservance(t *testing.T) {
	// Start at July 1, 2025 (a Tuesday), which is NOT a holiday.
	start := time.Date(2025, 7, 1, 9, 0, 0, 0, time.UTC)

	// Construct the base rrule with quarterly recurrence
	rp, err := NewRRulePlus(ROptionPlus{
		ROption: rrule.ROption{
			Freq:     rrule.MONTHLY,
			Interval: 3,
			Dtstart:  start,
			Count:    4,
		},
		ShiftOffHolidays: true,
		ShiftOffWeekend:  true,
		Observance:       ObservanceNextBizDay,
		ISOCode:          "US",
	})
	require.NoError(t, err)

	// Define the expected shifted dates (using U.S. holiday observance)
	expected := []time.Time{
		time.Date(2025, 7, 1, 9, 0, 0, 0, time.UTC),  // Not a holiday
		time.Date(2025, 10, 1, 9, 0, 0, 0, time.UTC), // Not a holiday
		time.Date(2026, 1, 2, 9, 0, 0, 0, time.UTC),  // Jan 1 is New Year's Day → shift to Jan 2
		time.Date(2026, 4, 1, 9, 0, 0, 0, time.UTC),  // Not a holiday
	}

	// Evaluate
	var actual []time.Time
	cursor := start.Add(-time.Second)

	for i := 0; i < len(expected); i++ {
		next := rp.After(cursor, false)
		require.False(t, next.IsZero(), "Unexpected nil occurrence")
		actual = append(actual, next)
		cursor = next.Add(time.Second)
	}

	// Compare results
	require.Equal(t, expected, actual)
}

func TestFiscalCycle_Monthly_WithHolidayObservance(t *testing.T) {
	start := time.Date(2025, 9, 1, 9, 0, 0, 0, time.UTC) // Labor Day

	rp, err := NewRRulePlus(ROptionPlus{
		ROption: rrule.ROption{
			Freq:    rrule.MONTHLY,
			Dtstart: start,
		},
		ShiftOffHolidays: true,
		ISOCode:          "US",
	})
	assert.NoError(t, err)

	next := rp.After(start.Add(-48*time.Hour), true)
	require.Equal(t, time.Date(2025, 9, 2, 9, 0, 0, 0, time.UTC), next) // Day after Labor Day
}

func TestFiscalCycle_Weekly_WithWeekendShift(t *testing.T) {
	start := time.Date(2025, 6, 28, 10, 0, 0, 0, time.UTC) // Saturday

	rp, err := NewRRulePlus(ROptionPlus{
		ROption: rrule.ROption{
			Freq:    rrule.WEEKLY,
			Dtstart: start,
		},
		ShiftOffWeekend: true,
		ISOCode:         "US",
	})

	assert.NoError(t, err)

	next := rp.After(start.Add(-72*time.Hour), true)
	require.Equal(t, time.Date(2025, 6, 30, 10, 0, 0, 0, time.UTC), next) // Monday
}
