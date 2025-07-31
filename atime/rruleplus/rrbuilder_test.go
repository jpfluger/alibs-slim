package rruleplus

import (
	"fmt"
	"github.com/jpfluger/alibs-slim/atime"
	"github.com/stretchr/testify/require"
	"github.com/teambition/rrule-go"
	"testing"
	"time"
)

func TestRRBuilderMonthDay_Build(t *testing.T) {
	rule := NewRRBuilderMonthDay(
		time.Date(2025, 12, 25, 0, 0, 0, 0, time.UTC),
		atime.TIMEUNIT_YEARLY,
	).
		AddJWBeforeDaily("reminder", 7, 1).
		WithDuration(1, atime.TIMEUNIT_DAILY).
		Build()

	require.Equal(t, 25, rule.ROptions.ByMonthDay[0])
	require.Equal(t, 12, rule.ROptions.ByMonth[0])
	require.Equal(t, atime.TIMEUNIT_YEARLY, rule.ROptions.Freq)
	require.Len(t, rule.JoinWindows, 2)
	require.Equal(t, "reminder", rule.JoinWindows[0].Tag)
}

func TestToRRuleNth_Positive(t *testing.T) {
	tests := []struct {
		n      int
		day    time.Weekday
		expect rrule.Weekday
	}{
		{1, time.Monday, rrule.MO.Nth(1)},
		{2, time.Tuesday, rrule.TU.Nth(2)},
		{3, time.Friday, rrule.FR.Nth(3)},
		{4, time.Sunday, rrule.SU.Nth(4)},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%d-%s", tt.n, tt.day), func(t *testing.T) {
			rr := toRRuleNth(tt.n, tt.day)
			require.Equal(t, tt.expect.String(), rr.String())
			require.Equal(t, tt.expect.N(), rr.N())
		})
	}
}

func TestToRRuleNth_Negative(t *testing.T) {
	tests := []struct {
		n      int
		day    time.Weekday
		expect rrule.Weekday
	}{
		{-1, time.Monday, rrule.MO.Nth(-1)},
		{-2, time.Wednesday, rrule.WE.Nth(-2)},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%d-%s", tt.n, tt.day), func(t *testing.T) {
			rr := toRRuleNth(tt.n, tt.day)
			require.Equal(t, tt.expect.String(), rr.String())
			require.Equal(t, tt.expect.N(), rr.N())
		})
	}
}

func TestToRRuleNth_DefaultFallback(t *testing.T) {
	invalidWeekday := time.Weekday(99)
	result := toRRuleNth(1, invalidWeekday)

	require.Equal(t, rrule.MO.Day(), result.Day()) // fallback to Monday
	require.Equal(t, 1, result.N())                // retain original N value
}

func TestRRBuilder_MonthDay(t *testing.T) {
	rule := NewRRBuilderMonthDay(
		time.Date(2025, 12, 25, 0, 0, 0, 0, time.UTC),
		atime.TIMEUNIT_YEARLY,
	).
		AddJWBeforeDaily("reminder", 30, 15, 1).
		WithDuration(1, atime.TIMEUNIT_DAILY).
		Allow().
		Build()

	require.Equal(t, []int{12}, rule.ROptions.ByMonth)
	require.Equal(t, []int{25}, rule.ROptions.ByMonthDay)
	require.Equal(t, 1, rule.Duration)
	require.False(t, rule.IsDeny)
	require.Len(t, rule.JoinWindows, 3)
	require.Equal(t, "reminder", rule.JoinWindows[0].Tag)
}

func TestRRBuilder_AnyTime_Allow(t *testing.T) {
	rule := NewRRBuilderAnyTime().
		Allow().
		WithPriority(1).
		Build()

	require.True(t, rule.IsAnyTime)
	require.False(t, rule.IsDeny)
	require.Equal(t, 1, rule.Priority)
}

func TestRRBuilder_AnyTime_Deny(t *testing.T) {
	rule := NewRRBuilderAnyTime().
		Deny().
		WithPriority(10).
		Build()

	require.True(t, rule.IsAnyTime)
	require.True(t, rule.IsDeny)
	require.Equal(t, 10, rule.Priority)
}

func TestRRBuilder_Weekday_BusinessHours(t *testing.T) {
	rule := NewRRBuilderWeekday(2025,
		[]rrule.Weekday{rrule.MO, rrule.TU, rrule.WE, rrule.TH, rrule.FR},
		9, 17,
	).
		Allow().
		WithPriority(20).
		Build()

	require.Equal(t, atime.TIMEUNIT_DAILY, rule.ROptions.Freq)
	require.ElementsMatch(t, []rrule.Weekday{rrule.MO, rrule.TU, rrule.WE, rrule.TH, rrule.FR}, rule.ROptions.ByDay)
	require.ElementsMatch(t, []int{9, 10, 11, 12, 13, 14, 15, 16}, rule.ROptions.ByHour)
	require.False(t, rule.IsDeny)
}

func TestRRBuilder_Weekday_TuesdayAfternoon(t *testing.T) {
	rule := NewRRBuilderWeekday(2025,
		[]rrule.Weekday{rrule.TU},
		14, 18,
	).
		Allow().
		Build()

	require.ElementsMatch(t, []rrule.Weekday{rrule.TU}, rule.ROptions.ByDay)
	require.ElementsMatch(t, []int{14, 15, 16, 17}, rule.ROptions.ByHour)
	require.False(t, rule.IsDeny)
}

func TestRRBuilder_SpecificDate_July4(t *testing.T) {
	rule := NewRRBuilderSpecificDate(time.Date(0, 7, 4, 0, 0, 0, 0, time.UTC)).
		Deny().
		WithDuration(1, atime.TIMEUNIT_DAILY).
		AddJWBeforeDaily("alert", 1).
		Build()

	require.True(t, rule.IsDeny)
	require.Equal(t, []int{7}, rule.ROptions.ByMonth)
	require.Equal(t, []int{4}, rule.ROptions.ByMonthDay)
	require.Len(t, rule.JoinWindows, 1)
	require.Equal(t, "alert", rule.JoinWindows[0].Tag)
}

func TestRRBuilder_NthWeekday_FirstMonday(t *testing.T) {
	rule := NewRRBuilderNthWeekday(0, 1, time.Monday).
		WithDuration(1, atime.TIMEUNIT_DAILY).
		AddJWBeforeDaily("notify", 7).
		Build()

	require.Len(t, rule.ROptions.ByDay, 1)
	require.Equal(t, rrule.MO.Day(), rule.ROptions.ByDay[0].Day()) // 0 for Monday
	require.Equal(t, 1, rule.ROptions.ByDay[0].N())                // 1st Monday
	require.Len(t, rule.JoinWindows, 1)
	require.Equal(t, "notify", rule.JoinWindows[0].Tag)
}

func TestRRBuilder_NthWeekday_LastFriday(t *testing.T) {
	rule := NewRRBuilderNthWeekday(0, -1, time.Friday).
		Allow().
		WithPriority(5).
		Build()

	require.Len(t, rule.ROptions.ByDay, 1)
	require.Equal(t, rrule.FR.Day(), rule.ROptions.ByDay[0].Day()) // 4 for Friday
	require.Equal(t, -1, rule.ROptions.ByDay[0].N())               // Last Friday
	require.False(t, rule.IsDeny)
	require.Equal(t, 5, rule.Priority)
}

func TestRRBuilder_Weekday_BusinessHours_StackedWindows_IsBetween(t *testing.T) {
	// This rule anchors every hour from 9AM–5PM, and each anchor allows a 2-hour access window.
	// That means matches occur at:
	// - 09:00–11:00
	// - 10:00–12:00
	// - 11:00–13:00
	// ...
	// A time like 11:15 will fall into the 11:00–13:00 window.
	//
	// Compare with:
	// * TestRRBuilder_Weekday_BusinessHours_SingleAnchor_IsBetween_2Hours
	// * TestRRBuilder_Weekday_BusinessHours_SingleAnchor_IsBetween_8Hours

	rule := NewRRBuilderWeekday(2025,
		[]rrule.Weekday{rrule.MO, rrule.TU, rrule.WE, rrule.TH, rrule.FR},
		9, 17, // Anchors: 9AM to 4PM inclusive
	).WithDuration(2, atime.TIMEUNIT_HOURLY).
		Allow().
		WithPriority(20).
		Build()

	inside := time.Date(2025, 6, 23, 9, 30, 0, 0, time.UTC)       // Monday 9:30AM (within 9–11 window)
	alsoInside := time.Date(2025, 6, 23, 11, 15, 0, 0, time.UTC)  // Monday 11:15AM (within 11–13 window)
	alsoInside2 := time.Date(2025, 6, 23, 17, 15, 0, 0, time.UTC) // Monday 17:15 (within 16–18 window)
	outside := time.Date(2025, 6, 23, 18, 15, 0, 0, time.UTC)     // Monday 18:15 (outside 16–18 window)

	in, err := rule.IsBetween(inside)
	require.NoError(t, err)
	require.True(t, in, "Expected 9:30 to match 9–11 window")

	in2, err := rule.IsBetween(alsoInside)
	require.NoError(t, err)
	require.True(t, in2, "Expected 11:15 to match 11–13 window")

	in3, err := rule.IsBetween(alsoInside2)
	require.NoError(t, err)
	require.True(t, in3, "Expected 17:15 to match 16-18 window")

	out, err := rule.IsBetween(outside)
	require.NoError(t, err)
	require.False(t, out, "Expected 18:15 to falls outside 16-18 window")
}

func TestRRBuilder_Weekday_BusinessHours_SingleAnchor_IsBetween_2Hours(t *testing.T) {
	// This rule only anchors a single time at 9AM (via WithHourRange(9, -1)),
	// and allows a 2-hour window (9AM–11AM). Any time outside that won't match.

	rule := NewRRBuilderWeekday(2025,
		[]rrule.Weekday{rrule.MO, rrule.TU, rrule.WE, rrule.TH, rrule.FR},
		9, -1, // Anchor ONLY at 9AM
	).WithDuration(2, atime.TIMEUNIT_HOURLY).
		Allow().
		WithPriority(20).
		Build()

	inside := time.Date(2025, 6, 23, 9, 30, 0, 0, time.UTC)   // Monday 9:30AM (inside 9–11 window)
	outside := time.Date(2025, 6, 23, 11, 15, 0, 0, time.UTC) // Monday 11:15AM (outside window)

	in, err := rule.IsBetween(inside)
	require.NoError(t, err)
	require.True(t, in, "Expected 9:30 to match single 9–11 window")

	out, err := rule.IsBetween(outside)
	require.NoError(t, err)
	require.False(t, out, "Expected 11:15 to be outside 9–11 window")
}

// Tests rule where the only anchor is 9am (Mon–Fri), and the duration extends 8 hours from that point.
// This simulates rules like “allow access starting at 9am for the duration of a full workday.”
//
// Unlike the previous test, this does NOT use a range of hours, but a single starting point with fixed duration.
func TestRRBuilder_Weekday_BusinessHours_SingleAnchor_IsBetween_8Hours(t *testing.T) {
	rule := NewRRBuilderWeekday(2025,
		[]rrule.Weekday{rrule.MO, rrule.TU, rrule.WE, rrule.TH, rrule.FR},
		9, -1, // Only 9AM anchor
	).WithDuration(8, atime.TIMEUNIT_HOURLY). // Full workday
							Allow().
							WithPriority(20).
							Build()

	inside := time.Date(2025, 6, 23, 9, 30, 0, 0, time.UTC)   // Within 8-hour span
	outside := time.Date(2025, 6, 23, 17, 15, 0, 0, time.UTC) // Just beyond window

	in, err := rule.IsBetween(inside)
	require.NoError(t, err)
	require.True(t, in)

	out, err := rule.IsBetween(outside)
	require.NoError(t, err)
	require.False(t, out)
}

func TestRRBuilder_Weekday_TuesdayAfternoon_StackedAnchors(t *testing.T) {
	rule := NewRRBuilderWeekday(2025,
		[]rrule.Weekday{rrule.TU},
		14, 18, // 2PM to 6PM
	).WithDuration(1, atime.TIMEUNIT_HOURLY).
		Allow().
		Build()

	// Matches 14:00–15:00
	inside := time.Date(2025, 6, 24, 14, 15, 0, 0, time.UTC)

	// Matches 15:00–16:00 anchor
	alsoInside := time.Date(2025, 6, 24, 15, 30, 0, 0, time.UTC)

	// After final window: 17:00–18:00
	tooLate := time.Date(2025, 6, 24, 18, 15, 0, 0, time.UTC)

	in, err := rule.IsBetween(inside)
	require.NoError(t, err)
	require.True(t, in)

	in2, err := rule.IsBetween(alsoInside)
	require.NoError(t, err)
	require.True(t, in2)

	out, err := rule.IsBetween(tooLate)
	require.NoError(t, err)
	require.False(t, out)
}

func TestRRBuilder_SpecificDate_July4_IsBetween_And_MatchJoinWindow(t *testing.T) {
	// Allow rule for July 4 with join window
	allowRule := NewRRBuilderSpecificDate(time.Date(2025, 7, 4, 0, 0, 0, 0, time.UTC)).
		WithDuration(1, atime.TIMEUNIT_DAILY).
		AddJWBeforeDaily("alert", 1).
		Build()

	// July 4 — should match
	match := time.Date(2025, 7, 4, 10, 0, 0, 0, time.UTC)
	in, err := allowRule.IsBetween(match)
	require.NoError(t, err)
	require.True(t, in)

	// Deny rule for July 4 with same join window
	denyRule := NewRRBuilderSpecificDate(time.Date(2025, 7, 4, 0, 0, 0, 0, time.UTC)).
		Deny().
		WithDuration(1, atime.TIMEUNIT_DAILY).
		AddJWBeforeDaily("alert", 1).
		Build()

	in, err = denyRule.IsBetween(match)
	require.NoError(t, err)
	require.False(t, in) // Deny rule blocks

	// Join window match on July 3 (1 day before)
	jwTime := time.Date(2025, 7, 3, 9, 0, 0, 0, time.UTC)
	jw, err := denyRule.MatchJoinWindow(jwTime, true, true)
	require.NoError(t, err)
	require.NotNil(t, jw)
	require.Equal(t, "alert", jw.Tag)

	// No match on July 5
	notMatch := time.Date(2025, 7, 5, 0, 0, 0, 0, time.UTC)
	out, err := denyRule.IsBetween(notMatch)
	require.NoError(t, err)
	require.True(t, out)

	jw, err = denyRule.MatchJoinWindow(notMatch, true, true)
	require.NoError(t, err)
	require.Nil(t, jw)
}

func TestRRBuilder_NthWeekday_FirstMonday_IsBetween_And_MatchJoinWindow(t *testing.T) {
	rule := NewRRBuilderNthWeekday(2025, 1, time.Monday).
		WithHourRange(9, -1). // Anchor to 9AM
		WithDuration(3, atime.TIMEUNIT_HOURLY).
		AddJWBeforeDaily("notify", 7).
		Allow().
		Build()

	// Sept 1, 2025 is the 1st Monday
	match := time.Date(2025, 9, 1, 9, 0, 0, 0, time.UTC)
	later := time.Date(2025, 9, 1, 13, 0, 0, 0, time.UTC)

	in, err := rule.IsBetween(match)
	require.NoError(t, err)
	require.True(t, in)

	out, err := rule.IsBetween(later)
	require.NoError(t, err)
	require.False(t, out)

	// Match 7-day join window prior (Aug 25)
	jwTime := time.Date(2025, 8, 25, 9, 0, 0, 0, time.UTC)
	jw, err := rule.MatchJoinWindow(jwTime, true, true)
	require.NoError(t, err)
	require.NotNil(t, jw)
	require.Equal(t, "notify", jw.Tag)

	// Outside join window
	jwMiss := time.Date(2025, 8, 23, 10, 0, 0, 0, time.UTC)
	jw, err = rule.MatchJoinWindow(jwMiss, true, true)
	require.NoError(t, err)
	require.Nil(t, jw)
}

func TestRRBuilder_NthWeekday_LastFriday_IsBetween(t *testing.T) {
	rule := NewRRBuilderNthWeekday(2025, -1, time.Friday).
		WithHourRange(10, -1).                  // Start at 10AM
		WithDuration(4, atime.TIMEUNIT_HOURLY). // 10AM–2PM
		Allow().
		WithPriority(5).
		Build()

	start := time.Date(2025, 6, 27, 10, 0, 0, 0, time.UTC)
	mid := time.Date(2025, 6, 27, 13, 30, 0, 0, time.UTC)
	tooLate := time.Date(2025, 6, 27, 15, 0, 0, 0, time.UTC)

	match, err := rule.IsBetween(start)
	require.NoError(t, err)
	require.True(t, match)

	alsoMatch, err := rule.IsBetween(mid)
	require.NoError(t, err)
	require.True(t, alsoMatch)

	notMatch, err := rule.IsBetween(tooLate)
	require.NoError(t, err)
	require.False(t, notMatch)
}

func TestRRBuilder_NthWeekday_FirstMonday_NewYorkTime(t *testing.T) {
	loc, err := time.LoadLocation("America/New_York")
	require.NoError(t, err)

	// Rule: First Monday of month at 9AM New York time, duration 8 hours
	rule := NewRRBuilderNthWeekday(2025, 1, time.Monday).
		WithTimeZone("America/New_York").
		WithBeginTime(time.Date(2025, 9, 1, 9, 0, 0, 0, time.Local)). // local ignored; zone reset
		WithHourRange(9, -1).                                         // 9AM local time
		WithDuration(8, atime.TIMEUNIT_HOURLY).
		Allow().
		Build()

	// September 1, 2025 is a Monday
	localStart := time.Date(2025, 9, 1, 9, 0, 0, 0, loc) // 9AM NY (EDT)
	localLate := time.Date(2025, 9, 1, 18, 0, 0, 0, loc) // 6PM NY

	// Convert to UTC for evaluation (since RRule uses UTC internally)
	in, err := rule.IsBetween(localStart.UTC())
	require.NoError(t, err)
	require.True(t, in, "Expected 9AM NY (UTC) to fall within window")

	out, err := rule.IsBetween(localLate.UTC())
	require.NoError(t, err)
	require.False(t, out, "Expected 6PM NY (UTC) to be outside window")
}

func TestRRBuilder_NthWeekday_DST_Start_SpringForward(t *testing.T) {
	loc, err := time.LoadLocation("America/New_York")
	require.NoError(t, err)

	// First test: using NthWeekday
	rule := NewRRBuilderNthWeekday(2025, 2, time.Monday). // 2nd Monday of March 2025
								WithTimeZone("America/New_York").
								WithBeginTime(time.Date(2025, 3, 10, 9, 0, 0, 0, time.Local)). // anchor to local time
								WithHourRange(9, -1).                                          // 9AM local
								WithDuration(2, atime.TIMEUNIT_HOURLY).                        // 9–11AM
								Allow().
								Build()

	localTime := time.Date(2025, 3, 10, 9, 30, 0, 0, loc) // 9:30AM NY (EDT, post-DST)

	in, err := rule.IsBetween(localTime.UTC())
	require.NoError(t, err)
	require.True(t, in, "Expected time to fall within post-DST window")

	// Second test: using MonthDay + FirstNthWeekday
	anchor := FirstNthWeekday(2025, time.March, time.Monday, 2)

	rule2 := NewRRBuilderMonthDay(anchor, atime.TIMEUNIT_YEARLY).
		WithTimeZone("America/New_York").
		WithBeginTime(anchor).
		WithHourRange(9, -1).
		WithDuration(2, atime.TIMEUNIT_HOURLY).
		Allow().
		Build()

	in2, err := rule2.IsBetween(localTime.UTC())
	require.NoError(t, err)
	require.True(t, in2, "MonthDay builder should also match post-DST window")
}

func TestRRBuilder_NthWeekday_DST_End_FallBack(t *testing.T) {
	loc, err := time.LoadLocation("America/New_York")
	require.NoError(t, err)

	// In 2025, DST ends on Nov 2 (Sunday), so Nov 3 is the first Monday after
	rule := NewRRBuilderNthWeekday(2025, 1, time.Monday).
		WithTimeZone("America/New_York").
		WithHourRange(9, -1).                   // 9AM local
		WithDuration(3, atime.TIMEUNIT_HOURLY). // 9AM–12PM
		Allow().
		Build()

	localValid := time.Date(2025, 11, 3, 10, 0, 0, 0, loc) // Within window
	localLate := time.Date(2025, 11, 3, 13, 30, 0, 0, loc) // After window

	in, err := rule.IsBetween(localValid.UTC())
	require.NoError(t, err)
	require.True(t, in, "Expected valid time during fallback to match")

	out, err := rule.IsBetween(localLate.UTC())
	require.NoError(t, err)
	require.False(t, out, "Expected time after window to not match")
}

// For working versions, see:
// * TestRRBuilderFiscalCycle_Quarterly_WithHolidayObservance
// * TestRRBuilderFiscalCycle_US_TaxDay_WithObservance
//
// DEPRECATION NOTE: This test demonstrates a manually assembled rule set for handling a fixed-date policy
// scenario (e.g., April 15 as Tax Day), with fallback and weekend handling. However, this
// approach is fragile and prone to bugs.
//
// ❌ Limitations of Manual Rule Composition:
// - The fallback rule (e.g., April 16) must be manually constructed with an exact BeginTime
//   corresponding to the evaluation year. If omitted or misaligned, join windows will anchor
//   to incorrect dates (e.g., the wrong year or an unexpected reference).
// - Deny rules for weekends are global but don’t automatically suppress specific date rules
//   unless all timing is carefully orchestrated.
// - Hardcoded static dates (like `time.Date(0, 4, 15, ...)`) lose contextual meaning unless
//   BeginTime is explicitly set for each rule and updated per year.
// - Subtle bugs arise if Duration or DurationUnit are not explicitly declared.
//
// ✅ Recommended Alternative:
// Use `RRBuilder.BuildSpecificDateStack()` instead of manual stacking. It encapsulates:
// - Proper anchoring of the base date,
// - Smart fallback to next valid weekday,
// - Optional holiday exclusion,
// - Automatic rule priority assignment and denial suppression.
//
// This DSL method reduces human error and ensures the generated rules remain valid across years.
//
// See: `TestRRBuilder_BuildSpecificDateStack_April15Scenarios()` for a more robust approach.
//func TestTaxDayRuleSet(t *testing.T) {
//	rules := RRuleExtends{
//		NewRRBuilderWeekday(2025, []rrule.Weekday{rrule.SA, rrule.SU}, 0, 24).
//			Deny().
//			WithDuration(1, TIMEUNIT_DAILY).
//			WithPriority(100).
//			Build(),
//
//		NewRRBuilderSpecificDate(time.Date(0, 4, 16, 0, 0, 0, 0, time.UTC)). // ⬅️ year 0
//											Allow().
//											WithDuration(1, TIMEUNIT_DAILY).
//											AddJWBeforeDaily("tax_notice", 1).
//											WithPriority(90).
//											Build(),
//
//		NewRRBuilderSpecificDate(time.Date(0, 4, 15, 0, 0, 0, 0, time.UTC)). // ⬅️ year 0
//											Allow().
//											WithDuration(1, TIMEUNIT_DAILY).
//											AddJWBeforeDaily("tax_notice", 7, 1).
//											WithPriority(10).
//											Build(),
//	}
//
//	// April 15, 2026 → Saturday → denied
//	result, err := rules.Evaluate(time.Date(2026, 4, 15, 10, 0, 0, 0, time.UTC))
//	require.NoError(t, err)
//	require.False(t, result)
//
//	// April 16, 2026 → Sunday fallback → allowed
//	result, err = rules.Evaluate(time.Date(2026, 4, 16, 10, 0, 0, 0, time.UTC))
//	require.NoError(t, err)
//	require.True(t, result)
//}

// This is imperfect too.
//func TestRRBuilder_BuildSpecificDateStack_April15Scenarios(t *testing.T) {
//	tests := []struct {
//		name        string
//		april15     time.Time
//		holidays    map[int]bool
//		expectAllow time.Time
//		expectDeny  *time.Time
//	}{
//		{
//			name:        "Weekday April 15",
//			april15:     time.Date(2025, 4, 15, 0, 0, 0, 0, time.UTC), // Tuesday
//			holidays:    map[int]bool{},
//			expectAllow: time.Date(2025, 4, 15, 10, 0, 0, 0, time.UTC),
//			expectDeny:  nil,
//		},
//		{
//			name:        "Weekend April 15, fallback to April 17 (Monday)",
//			april15:     time.Date(2028, 4, 15, 0, 0, 0, 0, time.UTC), // Saturday
//			holidays:    map[int]bool{},
//			expectAllow: time.Date(2028, 4, 17, 10, 0, 0, 0, time.UTC),
//			expectDeny:  ToPointer(time.Date(2028, 4, 15, 10, 0, 0, 0, time.UTC)),
//		},
//		{
//			name:        "Holiday April 16, fallback to April 17",
//			april15:     time.Date(2029, 4, 15, 0, 0, 0, 0, time.UTC), // Sunday
//			holidays:    map[int]bool{20290416: true},
//			expectAllow: time.Date(2029, 4, 17, 10, 0, 0, 0, time.UTC),
//			expectDeny:  ToPointer(time.Date(2029, 4, 15, 10, 0, 0, 0, time.UTC)),
//		},
//		{
//			name:        "Deny on Saturday, fallback to Monday",
//			april15:     time.Date(2034, 4, 15, 0, 0, 0, 0, time.UTC), // Saturday
//			holidays:    map[int]bool{},
//			expectAllow: time.Date(2034, 4, 17, 10, 0, 0, 0, time.UTC),
//			expectDeny:  ToPointer(time.Date(2034, 4, 15, 10, 0, 0, 0, time.UTC)),
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			builder := NewRRBuilderSpecificDate(tt.april15).
//				Allow().
//				WithDuration(1, TIMEUNIT_DAILY).
//				WithWeekendFallback(5).
//				ExcludeHolidaysMap(tt.holidays)
//
//			stack := builder.BuildSpecificDateStack()
//
//			// Confirm allow on fallback or original
//			allowed, err := stack.Evaluate(tt.expectAllow)
//			require.NoError(t, err)
//			require.True(t, allowed, "Expected to allow on %v", tt.expectAllow)
//
//			// If specified, check original date is denied
//			if tt.expectDeny != nil {
//				denied, err := stack.Evaluate(*tt.expectDeny)
//				require.NoError(t, err)
//				require.False(t, denied, "Expected to deny on %v", tt.expectDeny)
//			}
//		})
//	}
//}

func TestRRBuilderFiscalCycle_Quarterly_WithHolidayObservance(t *testing.T) {
	// July 1, 2025 @ 9:00am UTC is our cycle anchor
	start := time.Date(2025, 7, 1, 9, 0, 0, 0, time.UTC)

	// Build a quarterly rule using RRBuilder
	builder := NewRRBuilderMonthDay(start, atime.TIMEUNIT_MONTHLY).
		WithInterval(3).
		WithCount(4).
		WithShiftOffHolidays().
		WithShiftOffWeekend().
		WithObservance(ObservanceNextBizDay).
		WithISOCode("US")

	// Convert to RRulePlus
	rp, err := builder.rule.ToRRule()
	require.NoError(t, err)

	// Define expected values after holiday/weekend logic
	expected := []time.Time{
		time.Date(2025, 7, 1, 9, 0, 0, 0, time.UTC),  // Not a holiday
		time.Date(2025, 10, 1, 9, 0, 0, 0, time.UTC), // Not a holiday
		time.Date(2026, 1, 2, 9, 0, 0, 0, time.UTC),  // Jan 1 is holiday → shift
		time.Date(2026, 4, 1, 9, 0, 0, 0, time.UTC),  // Not a holiday
	}

	// Evaluate cycle
	var actual []time.Time
	cursor := start.Add(-time.Second)

	for i := 0; i < len(expected); i++ {
		next := rp.After(cursor, false)
		require.False(t, next.IsZero(), "Unexpected nil occurrence at index %d", i)
		actual = append(actual, next)
		cursor = next.Add(time.Second)
	}

	// Assert correctness
	require.Equal(t, expected, actual)
}

func TestRRBuilderFiscalCycle_US_TaxDay_WithObservance(t *testing.T) {
	// Anchor the cycle to April 15, 2025 at 9:00am UTC
	start := time.Date(2025, 4, 15, 9, 0, 0, 0, time.UTC)

	// Build a yearly recurrence on April 15 with shifting/observance rules
	builder := NewRRBuilderMonthDay(start, atime.TIMEUNIT_YEARLY).
		OnApril().
		OnDayOfMonth(15).
		WithCount(12).
		WithShiftOffHolidays().
		WithShiftOffWeekend().
		WithObservance(ObservanceNextBizDay).
		WithISOCode("US")

	// Convert to RRulePlus
	rp, err := builder.rule.ToRRule()
	require.NoError(t, err)

	// Define expected valid dates (business day adjusted)
	expected := []time.Time{
		time.Date(2025, 4, 15, 9, 0, 0, 0, time.UTC), // Tue (ok)
		time.Date(2026, 4, 15, 9, 0, 0, 0, time.UTC), // Wed (ok)
		time.Date(2027, 4, 15, 9, 0, 0, 0, time.UTC), // Thu (ok)
		time.Date(2028, 4, 17, 9, 0, 0, 0, time.UTC), // 15 = Sat → shift to Mon
		time.Date(2029, 4, 16, 9, 0, 0, 0, time.UTC), // 15 = Sun → shift to Mon
		time.Date(2030, 4, 15, 9, 0, 0, 0, time.UTC), // Mon (ok)
		time.Date(2031, 4, 15, 9, 0, 0, 0, time.UTC), // Tue (ok)
		time.Date(2032, 4, 15, 9, 0, 0, 0, time.UTC), // Thu (ok)
		time.Date(2033, 4, 15, 9, 0, 0, 0, time.UTC), // Fri (ok)
		time.Date(2034, 4, 17, 9, 0, 0, 0, time.UTC), // 15 = Sat → shift
		time.Date(2035, 4, 16, 9, 0, 0, 0, time.UTC), // 15 = Sun → shift
		time.Date(2036, 4, 15, 9, 0, 0, 0, time.UTC), // Tue (ok)
	}

	// Evaluate actual recurrence set
	var actual []time.Time
	cursor := start.Add(-time.Second)

	for i := 0; i < len(expected); i++ {
		next := rp.After(cursor, false)
		require.False(t, next.IsZero(), "Unexpected nil occurrence at index %d", i)
		actual = append(actual, next)
		cursor = next.Add(time.Second)
	}

	// Verify actual vs expected
	require.Equal(t, expected, actual)

	// Compare difference with 9am
	// Anchor the cycle to April 15, 2025 at 12:00am UTC
	start = time.Date(2025, 4, 15, 0, 0, 0, 0, time.UTC)

	// Build a yearly recurrence on April 15 with shifting/observance rules
	builder = NewRRBuilderMonthDay(start, atime.TIMEUNIT_YEARLY).
		OnApril().
		OnDayOfMonth(15).
		WithCount(12).
		WithShiftOffHolidays().
		WithShiftOffWeekend().
		WithObservance(ObservanceNextBizDay).
		WithISOCode("US")

	// Convert to RRulePlus
	rp, err = builder.rule.ToRRule()
	require.NoError(t, err)

	for ii, t := range expected {
		expected[ii] = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	}

	actual = []time.Time{}
	cursor = start.Add(-time.Second)

	for i := 0; i < len(expected); i++ {
		next := rp.After(cursor, false)
		require.False(t, next.IsZero(), "Unexpected nil occurrence at index %d", i)
		actual = append(actual, next)
		cursor = next.Add(time.Second)
	}

	// Verify actual vs expected
	require.Equal(t, expected, actual)
}
