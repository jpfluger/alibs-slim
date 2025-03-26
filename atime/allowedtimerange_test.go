package atime

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestAllowedTimeRange(t *testing.T) {
	location, _ := time.LoadLocation("America/New_York")
	atr := AllowedTimeRange{
		IncludeWeekdays: true,
		IncludeWeekends: true,
		TimeZone:        location,
	}

	t.Run("IsNowAllowed Weekday", func(t *testing.T) {
		fnTimeNow := func() time.Time {
			return time.Date(2023, 12, 27, 10, 0, 0, 0, location) // Wednesday
		}

		atr.IncludeWeekdays = true
		atr.IncludeWeekends = false
		assert.True(t, atr.IsNowAllowedWithTime(fnTimeNow), "Current time during weekday should be allowed")
	})

	t.Run("IsNowAllowed Weekend", func(t *testing.T) {
		fnTimeNow := func() time.Time {
			return time.Date(2023, 12, 30, 10, 0, 0, 0, location) // Saturday
		}

		atr.IncludeWeekdays = false
		atr.IncludeWeekends = true
		assert.True(t, atr.IsNowAllowedWithTime(fnTimeNow), "Current time during weekend should be allowed")
	})

	t.Run("IsNowAllowed SpecificDay", func(t *testing.T) {
		fnTimeNow := func() time.Time {
			return time.Date(2023, 12, 26, 10, 0, 0, 0, location) // Tuesday
		}

		atr.IncludeWeekdays = false
		atr.IncludeWeekends = false
		atr.SpecificDays = []time.Weekday{time.Tuesday}
		assert.True(t, atr.IsNowAllowedWithTime(fnTimeNow), "Specific day (Tuesday) should be allowed")
	})

	t.Run("IsNowAllowed ExcludeDate", func(t *testing.T) {
		fnTimeNow := func() time.Time {
			return time.Date(2023, 12, 25, 10, 0, 0, 0, location) // Christmas
		}

		atr.IncludeWeekdays = true
		atr.ExcludeDates = []time.Time{
			time.Date(2023, 12, 25, 0, 0, 0, 0, location), // Christmas
		}

		assert.False(t, atr.IsNowAllowedWithTime(fnTimeNow), "Excluded date (Christmas) should not be allowed")
	})
}

func TestAllowedTimeRange_IsAllowedAt(t *testing.T) {
	location := time.UTC
	atr := &AllowedTimeRange{
		IncludeWeekdays: true,
		Start:           time.Date(0, 1, 1, 9, 0, 0, 0, location),  // 9 AM
		End:             time.Date(0, 1, 1, 17, 0, 0, 0, location), // 5 PM
		TimeZone:        location,
	}

	t.Run("Allowed Weekday Time", func(t *testing.T) {
		now := time.Date(2023, 12, 29, 10, 0, 0, 0, location) // Friday
		assert.True(t, atr.IsAllowedAt(now), "Time within allowed range on a weekday should be allowed")
	})

	t.Run("Disallowed Weekday Time", func(t *testing.T) {
		now := time.Date(2023, 12, 29, 18, 0, 0, 0, location) // Friday, after 5 PM
		assert.False(t, atr.IsAllowedAt(now), "Time outside allowed range on a weekday should not be allowed")
	})

	t.Run("Allowed Weekend Time", func(t *testing.T) {
		atr.IncludeWeekends = true
		now := time.Date(2023, 12, 30, 10, 0, 0, 0, location) // Saturday
		assert.True(t, atr.IsAllowedAt(now), "Time within allowed range on a weekend should be allowed")
	})
}

func TestAllowedTimeRanges_IsAllowedAt(t *testing.T) {
	location := time.UTC
	atrs := AllowedTimeRanges{
		&AllowedTimeRange{
			IncludeWeekdays: true,
			Start:           time.Date(0, 1, 1, 9, 0, 0, 0, location),  // 9 AM
			End:             time.Date(0, 1, 1, 17, 0, 0, 0, location), // 5 PM
			TimeZone:        location,
		},
	}

	t.Run("Allowed Time in Range", func(t *testing.T) {
		now := time.Date(2023, 12, 29, 10, 0, 0, 0, location) // Friday
		assert.True(t, atrs.IsAllowedAt(now), "Time within allowed range should be allowed")
	})

	t.Run("Disallowed Time Out of Range", func(t *testing.T) {
		now := time.Date(2023, 12, 29, 18, 0, 0, 0, location) // Friday, after 5 PM
		assert.False(t, atrs.IsAllowedAt(now), "Time outside allowed range should not be allowed")
	})
}
