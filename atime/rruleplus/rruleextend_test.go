package rruleplus

import (
	"errors"
	"fmt"
	"github.com/jpfluger/alibs-slim/ageo"
	"github.com/jpfluger/alibs-slim/atime"
	"github.com/stretchr/testify/require"
	"github.com/teambition/rrule-go"
	"testing"
	"time"
)

func TestRRuleExtend_IsBetween_YEARLY(t *testing.T) {
	type timeCase struct {
		label     string
		now       time.Time
		expect    bool
		expectMsg string
	}

	tests := []struct {
		name        string
		rrule       RRuleExtend
		description string
		cases       []timeCase
	}{
		{
			name: "Yearly full span of 2025",
			rrule: RRuleExtend{
				ROptions: ROptionExtend{
					Freq:       atime.TIMEUNIT_YEARLY,
					Interval:   1,
					BeginTime:  atime.MustParsePtrRFC3339("2025-01-01T00:00:00Z"),
					UntilTime:  atime.MustParsePtrRFC3339("2025-12-31T23:59:59Z"),
					ByMonth:    []int{1},
					ByMonthDay: []int{1},
					ByHour:     []int{0},
				},
				Duration:     8760, // 365 days
				DurationUnit: atime.TIMEUNIT_HOURLY,
			},
			description: "One yearly event valid through all of 2025",
			cases: []timeCase{
				{
					label:     "preNow — before window",
					now:       atime.MustParseRFC3339("2024-12-31T23:59:59Z"),
					expect:    false,
					expectMsg: "Should not match before Jan 1, 2025",
				},
				{
					label:     "between — middle of 2025",
					now:       atime.MustParseRFC3339("2025-06-15T12:00:00Z"),
					expect:    true,
					expectMsg: "Should match mid-year",
				},
				{
					label:     "postNow — just after 2025",
					now:       atime.MustParseRFC3339("2026-01-01T00:00:00Z"),
					expect:    false,
					expectMsg: "Should not match in next year",
				},
			},
		},
		{
			name: "Every other year from 2025",
			rrule: RRuleExtend{
				ROptions: ROptionExtend{
					Freq:       atime.TIMEUNIT_YEARLY,
					Interval:   2, // 👈 key change here
					BeginTime:  atime.MustParsePtrRFC3339("2025-01-01T00:00:00Z"),
					ByMonth:    []int{1},
					ByMonthDay: []int{1},
					ByHour:     []int{0},
				},
				Duration:     8760, // entire year
				DurationUnit: atime.TIMEUNIT_HOURLY,
			},
			description: "Valid only on odd years starting 2025",
			cases: []timeCase{
				{
					label:     "2024 — before start",
					now:       atime.MustParseRFC3339("2024-06-01T00:00:00Z"),
					expect:    false,
					expectMsg: "Too early",
				},
				{
					label:     "2025 — valid year",
					now:       atime.MustParseRFC3339("2025-07-01T12:00:00Z"),
					expect:    true,
					expectMsg: "Every 2 years from 2025 — match",
				},
				{
					label:     "2026 — should skip",
					now:       atime.MustParseRFC3339("2026-06-01T00:00:00Z"),
					expect:    false,
					expectMsg: "Skipped year (interval=2)",
				},
				{
					label:     "2027 — valid again",
					now:       atime.MustParseRFC3339("2027-03-15T00:00:00Z"),
					expect:    true,
					expectMsg: "Every 2 years from 2025 — match",
				},
			},
		},
		{
			name: "Leap day recurrence",
			rrule: RRuleExtend{
				ROptions: ROptionExtend{
					Freq:       atime.TIMEUNIT_YEARLY,
					Interval:   1,
					BeginTime:  atime.MustParsePtrRFC3339("2024-02-29T00:00:00Z"),
					ByMonth:    []int{2},
					ByMonthDay: []int{29},
					ByHour:     []int{0},
				},
				Duration:     60,
				DurationUnit: atime.TIMEUNIT_MINUTELY,
			},
			description: "Occurs only in leap years",
			cases: []timeCase{
				{"Valid in 2024", atime.MustParseRFC3339("2024-02-29T00:30:00Z"), true, "Leap day match"},
				{"Invalid in 2025", atime.MustParseRFC3339("2025-02-28T00:30:00Z"), false, "Skipped in non-leap year"},
				{"Invalid in 2025 actual", atime.MustParseRFC3339("2025-02-29T00:00:00Z"), false, "Feb 29 doesn't exist"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, tc := range tt.cases {
				t.Run(tc.label, func(t *testing.T) {
					ok, err := tt.rrule.IsBetween(tc.now)
					require.NoError(t, err)
					require.Equal(t, tc.expect, ok, fmt.Sprintf("[%s] %s", tc.label, tt.description))

					if ok != tc.expect {
						t.Errorf(
							"Failed case: %s\nExpected: %v\nGot: %v\nRule logic: %s\nNow: %s\n",
							tc.label, tc.expect, ok, tc.expectMsg, tc.now.Format(time.RFC3339),
						)
					}
				})
			}
		})
	}
}

func TestRRuleExtend_IsBetween_MONTHLY(t *testing.T) {
	type timeCase struct {
		label     string
		now       time.Time
		expect    bool
		expectMsg string
	}

	tests := []struct {
		name        string
		rrule       RRuleExtend
		description string
		cases       []timeCase
	}{
		{
			name: "Monthly 24h event on the 15th",
			rrule: RRuleExtend{
				ROptions: ROptionExtend{
					Freq:       atime.TIMEUNIT_MONTHLY,
					Interval:   1,
					BeginTime:  atime.MustParsePtrRFC3339("2025-06-15T00:00:00Z"),
					ByMonthDay: []int{15},
					ByHour:     []int{0},
				},
				Duration:     24, // 24 hours
				DurationUnit: atime.TIMEUNIT_HOURLY,
			},
			description: "One event per month on the 15th for 24 hours",
			cases: []timeCase{
				{
					label:     "preNow — before window",
					now:       atime.MustParseRFC3339("2025-06-14T23:59:59Z"),
					expect:    false,
					expectMsg: "Should not match before June 15",
				},
				{
					label:     "between — middle of window",
					now:       atime.MustParseRFC3339("2025-06-15T12:00:00Z"),
					expect:    true,
					expectMsg: "Should match during 24h span on June 15",
				},
				{
					label:     "postNow — one second after end",
					now:       atime.MustParseRFC3339("2025-06-16T00:00:01Z"),
					expect:    false,
					expectMsg: "Should not match outside 24h range",
				},
				{
					label:     "nextMonth — future recurrence",
					now:       atime.MustParseRFC3339("2025-07-15T00:30:00Z"),
					expect:    true,
					expectMsg: "Should match July 15 recurrence",
				},
			},
		},
		{
			name: "Bi-Monthly full month event starting June",
			rrule: RRuleExtend{
				ROptions: ROptionExtend{
					Freq:       atime.TIMEUNIT_MONTHLY,
					Interval:   2,
					BeginTime:  atime.MustParsePtrRFC3339("2025-06-01T00:00:00Z"),
					ByMonthDay: []int{1},
					ByHour:     []int{0},
				},
				Duration:     744, // maximum monthly hours
				DurationUnit: atime.TIMEUNIT_HOURLY,
			},
			description: "Event valid entire month, only every other month",
			cases: []timeCase{
				{
					label:     "June — valid",
					now:       atime.MustParseRFC3339("2025-06-15T12:00:00Z"),
					expect:    true,
					expectMsg: "June should be active month",
				},
				{
					label:     "July — skipped month",
					now:       atime.MustParseRFC3339("2025-07-15T12:00:00Z"),
					expect:    false,
					expectMsg: "July should be skipped (every 2 months)",
				},
				{
					label:     "August — valid again",
					now:       atime.MustParseRFC3339("2025-08-01T00:01:00Z"),
					expect:    true,
					expectMsg: "August should be active",
				},
			},
		},
		{
			name: "Quarterly event on 1st of starting month",
			rrule: RRuleExtend{
				ROptions: ROptionExtend{
					Freq:       atime.TIMEUNIT_MONTHLY,
					Interval:   3,
					BeginTime:  atime.MustParsePtrRFC3339("2025-01-01T00:00:00Z"),
					ByMonthDay: []int{1},
					ByHour:     []int{0},
				},
				Duration:     24,
				DurationUnit: atime.TIMEUNIT_HOURLY,
			},
			description: "Occurs every 3 months on the 1st",
			cases: []timeCase{
				{
					label:     "Jan — valid",
					now:       atime.MustParseRFC3339("2025-01-01T01:00:00Z"),
					expect:    true,
					expectMsg: "Should match Jan 1",
				},
				{
					label:     "Feb — skip",
					now:       atime.MustParseRFC3339("2025-02-01T01:00:00Z"),
					expect:    false,
					expectMsg: "Should skip February",
				},
				{
					label:     "April — valid again",
					now:       atime.MustParseRFC3339("2025-04-01T00:00:00Z"),
					expect:    true,
					expectMsg: "Quarterly recurrence valid",
				},
			},
		},
		{
			name: "Implicit day from BeginTime without ByMonthDay",
			rrule: RRuleExtend{
				ROptions: ROptionExtend{
					Freq:      atime.TIMEUNIT_MONTHLY,
					Interval:  1,
					BeginTime: atime.MustParsePtrRFC3339("2025-06-10T00:00:00Z"),
					ByHour:    []int{0},
				},
				Duration:     24,
				DurationUnit: atime.TIMEUNIT_HOURLY,
			},
			description: "Should match on the day of BeginTime each month",
			cases: []timeCase{
				{
					label:     "June 10 — match",
					now:       atime.MustParseRFC3339("2025-06-10T12:00:00Z"),
					expect:    true,
					expectMsg: "Matches implicit day",
				},
				{
					label:     "June 11 — no match",
					now:       atime.MustParseRFC3339("2025-06-11T12:00:00Z"),
					expect:    false,
					expectMsg: "Outside recurrence",
				},
			},
		},
		{
			name: "Monthly event on last day",
			rrule: RRuleExtend{
				ROptions: ROptionExtend{
					Freq:       atime.TIMEUNIT_MONTHLY,
					Interval:   1,
					BeginTime:  atime.MustParsePtrRFC3339("2025-01-31T00:00:00Z"), // any last-day start
					ByMonthDay: []int{-1},                                         // <- precise & correct!
					ByHour:     []int{0},
				},
				Duration:     24,
				DurationUnit: atime.TIMEUNIT_HOURLY,
			},
			description: "Occurs on the last day of each month using ByMonthDay=-1",
			cases: []timeCase{
				{
					label:     "Feb (non-leap) — valid on Feb 28",
					now:       atime.MustParseRFC3339("2025-02-28T10:00:00Z"),
					expect:    true,
					expectMsg: "Should match Feb 28 as last day (non-leap year)",
				},
				{
					label:     "Apr 30 — valid",
					now:       atime.MustParseRFC3339("2025-04-30T10:00:00Z"),
					expect:    true,
					expectMsg: "Should match April 30 as last day",
				},
				{
					label:     "May 30 — invalid",
					now:       atime.MustParseRFC3339("2025-05-30T10:00:00Z"),
					expect:    false,
					expectMsg: "Should skip day before last",
				},
				{
					label:     "May 31 — valid",
					now:       atime.MustParseRFC3339("2025-05-31T10:00:00Z"),
					expect:    true,
					expectMsg: "Should match May 31",
				},
			},
		},
		{
			name: "Last day of month after quarter end",
			rrule: RRuleExtend{
				ROptions: ROptionExtend{
					Freq:       atime.TIMEUNIT_MONTHLY,
					Interval:   1,
					ByMonth:    []int{1, 4, 7, 10}, // Jan, Apr, Jul, Oct
					ByMonthDay: []int{-1},          // last day of those months
					ByHour:     []int{0},
					BeginTime:  atime.MustParsePtrRFC3339("2025-01-31T00:00:00Z"),
				},
				Duration:     24,
				DurationUnit: atime.TIMEUNIT_HOURLY,
			},
			description: "Runs at midnight on the last day of the month after quarter end",
			cases: []timeCase{
				{
					label:     "Jan 31 — Q4 closing action",
					now:       atime.MustParseRFC3339("2025-01-31T10:00:00Z"),
					expect:    true,
					expectMsg: "Valid: last day of Jan (after Q4)",
				},
				{
					label:     "Apr 30 — Q1 closing action",
					now:       atime.MustParseRFC3339("2025-04-30T10:00:00Z"),
					expect:    true,
					expectMsg: "Valid: last day of Apr (after Q1)",
				},
				{
					label:     "May 31 — Not a target month",
					now:       atime.MustParseRFC3339("2025-05-31T10:00:00Z"),
					expect:    false,
					expectMsg: "Should not match — not a post-quarter month",
				},
				{
					label:     "July 31 — Q2 closing action",
					now:       atime.MustParseRFC3339("2025-07-31T10:00:00Z"),
					expect:    true,
					expectMsg: "Valid: post-Q2 month end",
				},
				{
					label:     "October 31 — Q3 closing action",
					now:       atime.MustParseRFC3339("2025-10-31T10:00:00Z"),
					expect:    true,
					expectMsg: "Valid: end of October",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, tc := range tt.cases {
				t.Run(tc.label, func(t *testing.T) {
					ok, err := tt.rrule.IsBetween(tc.now)
					require.NoError(t, err)
					require.Equal(t, tc.expect, ok, fmt.Sprintf("[%s] %s", tc.label, tt.description))

					if ok != tc.expect {
						t.Errorf(
							"Failed case: %s\nExpected: %v\nGot: %v\nRule logic: %s\nNow: %s\n",
							tc.label, tc.expect, ok, tc.expectMsg, tc.now.Format(time.RFC3339),
						)
					}
				})
			}
		})
	}
}

func TestRRuleExtend_IsBetween_DAILY(t *testing.T) {
	type timeCase struct {
		label     string
		now       time.Time
		expect    bool
		expectMsg string
	}

	tests := []struct {
		name        string
		rrule       RRuleExtend
		description string
		cases       []timeCase
	}{
		{
			name: "Daily event 2-hour window",
			rrule: RRuleExtend{
				ROptions: ROptionExtend{
					Freq:      atime.TIMEUNIT_DAILY,
					Interval:  1,
					BeginTime: atime.MustParsePtrRFC3339("2025-06-01T09:00:00Z"),
					ByHour:    []int{9},
				},
				Duration:     2,
				DurationUnit: atime.TIMEUNIT_HOURLY,
			},
			description: "Every day from 9:00 to 11:00 UTC",
			cases: []timeCase{
				{
					label:     "before window",
					now:       atime.MustParseRFC3339("2025-06-10T08:59:59Z"),
					expect:    false,
					expectMsg: "Too early",
				},
				{
					label:     "in window",
					now:       atime.MustParseRFC3339("2025-06-10T10:00:00Z"),
					expect:    true,
					expectMsg: "During 2-hour window",
				},
				{
					label:     "after window",
					now:       atime.MustParseRFC3339("2025-06-10T11:01:00Z"),
					expect:    false,
					expectMsg: "Too late",
				},
			},
		},
		{
			name: "Daily event every other day",
			rrule: RRuleExtend{
				ROptions: ROptionExtend{
					Freq:      atime.TIMEUNIT_DAILY,
					Interval:  2,
					BeginTime: atime.MustParsePtrRFC3339("2025-06-01T06:00:00Z"),
					ByHour:    []int{6},
				},
				Duration:     1,
				DurationUnit: atime.TIMEUNIT_HOURLY,
			},
			description: "Occurs every 2 days",
			cases: []timeCase{
				{
					label:     "even day match",
					now:       atime.MustParseRFC3339("2025-06-03T06:15:00Z"),
					expect:    true,
					expectMsg: "June 3 is every 2nd day from June 1",
				},
				{
					label:     "odd day miss",
					now:       atime.MustParseRFC3339("2025-06-04T06:15:00Z"),
					expect:    false,
					expectMsg: "June 4 should not be valid (not every 2nd day from June 1)",
				},
				{
					label:     "next even day match",
					now:       atime.MustParseRFC3339("2025-06-05T06:15:00Z"),
					expect:    true,
					expectMsg: "June 5 is every 2nd day from June 1",
				},
				{
					label:     "next odd day miss",
					now:       atime.MustParseRFC3339("2025-06-06T06:15:00Z"),
					expect:    false,
					expectMsg: "June 6 should not be valid (not every 2nd day from June 1)",
				},
			},
		},
		{
			name: "Daily event with StartDate/EndDate bounds",
			rrule: RRuleExtend{
				ROptions: ROptionExtend{
					Freq:      atime.TIMEUNIT_DAILY,
					Interval:  1,
					BeginTime: atime.MustParsePtrRFC3339("2025-06-10T07:00:00Z"),
					ByHour:    []int{7},
				},
				StartDate:    atime.MustParsePtrRFC3339("2025-06-10T00:00:00Z"),
				EndDate:      atime.MustParsePtrRFC3339("2025-06-15T00:00:00Z"),
				Duration:     2,
				DurationUnit: atime.TIMEUNIT_HOURLY,
			},
			description: "Event bounded between June 10–14 inclusive",
			cases: []timeCase{
				{
					label:     "before StartDate",
					now:       atime.MustParseRFC3339("2025-06-09T07:30:00Z"),
					expect:    false,
					expectMsg: "Should not match before StartDate",
				},
				{
					label:     "inside window",
					now:       atime.MustParseRFC3339("2025-06-11T08:00:00Z"),
					expect:    true,
					expectMsg: "Valid match inside period",
				},
				{
					label:     "after EndDate",
					now:       atime.MustParseRFC3339("2025-06-15T07:00:00Z"),
					expect:    false,
					expectMsg: "Should not match on or after EndDate",
				},
			},
		},
		{
			name: "Daily event with recurrence Count limit",
			rrule: RRuleExtend{
				ROptions: ROptionExtend{
					Freq:      atime.TIMEUNIT_DAILY,
					Interval:  1,
					BeginTime: atime.MustParsePtrRFC3339("2025-06-01T05:00:00Z"),
					ByHour:    []int{5},
					Count:     3,
				},
				Duration:     1,
				DurationUnit: atime.TIMEUNIT_HOURLY,
			},
			description: "Runs only first 3 days",
			cases: []timeCase{
				{
					label:     "first day",
					now:       atime.MustParseRFC3339("2025-06-01T05:30:00Z"),
					expect:    true,
					expectMsg: "1st occurrence",
				},
				{
					label:     "third day",
					now:       atime.MustParseRFC3339("2025-06-03T05:30:00Z"),
					expect:    true,
					expectMsg: "3rd occurrence",
				},
				{
					label:     "fourth day — exceeds count",
					now:       atime.MustParseRFC3339("2025-06-04T05:30:00Z"),
					expect:    false,
					expectMsg: "Should not match after count limit",
				},
			},
		},
		{
			name: "Multiple daily windows (8–10 AM and 4–6 PM)",
			rrule: RRuleExtend{
				ROptions: ROptionExtend{
					Freq:      atime.TIMEUNIT_DAILY,
					Interval:  1,
					BeginTime: atime.MustParsePtrRFC3339("2025-06-01T08:00:00Z"),
					ByHour:    []int{8, 9, 16, 17},
				},
				Duration:     1,
				DurationUnit: atime.TIMEUNIT_HOURLY,
			},
			description: "Allow access during either AM or PM window",
			cases: []timeCase{
				{"early AM window", atime.MustParseRFC3339("2025-06-10T08:30:00Z"), true, "Within AM window"},
				{"late PM window", atime.MustParseRFC3339("2025-06-10T17:15:00Z"), true, "Within PM window"},
				{"noon disallowed", atime.MustParseRFC3339("2025-06-10T12:00:00Z"), false, "Outside defined hours"},
			},
		},
		{
			name: "Daily blackout window from 00:00–04:00",
			rrule: RRuleExtend{
				ROptions: ROptionExtend{
					Freq:      atime.TIMEUNIT_DAILY,
					Interval:  1,
					BeginTime: atime.MustParsePtrRFC3339("2025-06-01T00:00:00Z"),
					ByHour:    []int{0, 1, 2, 3},
				},
				Duration:     1,
				DurationUnit: atime.TIMEUNIT_HOURLY,
				IsDeny:       true,
			},
			description: "Deny access during blackout hours every day",
			cases: []timeCase{
				{"1 AM access denied", atime.MustParseRFC3339("2025-06-11T01:15:00Z"), false, "Blackout rule should deny"},
				{"5 AM access allowed", atime.MustParseRFC3339("2025-06-11T05:00:00Z"), true, "Outside blackout window"},
			},
		},
		{
			name: "Grace period after 9 AM login (10 mins)",
			rrule: RRuleExtend{
				ROptions: ROptionExtend{
					Freq:      atime.TIMEUNIT_DAILY,
					Interval:  1,
					BeginTime: atime.MustParsePtrRFC3339("2025-06-01T09:00:00Z"),
					ByHour:    []int{9},
				},
				Duration:     1,
				DurationUnit: atime.TIMEUNIT_MINUTELY,
				JoinWindows: JoinWindows{
					&JoinWindow{IsBefore: false, Duration: 10, DurationUnit: atime.TIMEUNIT_MINUTELY},
				},
				//RRIncType: RRULE_INC_TYPE_INCLUSIVE,
			},
			description: "Allow access up to 10 minutes after scheduled login",
			cases: []timeCase{
				{"8:59 AM too early", atime.MustParseRFC3339("2025-06-11T08:59:00Z"), false, "Before window"},
				{"9:00 AM on time", atime.MustParseRFC3339("2025-06-11T09:00:00Z"), true, "Start of window"},
				{"9:09 AM in grace but not IsBetween, thus invalid", atime.MustParseRFC3339("2025-06-11T09:09:59Z"), false, "Valid during grace but not for IsBetween func"},
				{"9:11 AM too late", atime.MustParseRFC3339("2025-06-11T09:11:00Z"), false, "Past grace period"},
			},
		},
		{
			name: "Daily 2 AM email report",
			rrule: RRuleExtend{
				ROptions: ROptionExtend{
					Freq:      atime.TIMEUNIT_DAILY,
					Interval:  1,
					BeginTime: atime.MustParsePtrRFC3339("2025-06-01T02:00:00Z"),
					ByHour:    []int{2},
				},
				Duration:     15,
				DurationUnit: atime.TIMEUNIT_MINUTELY,
			},
			description: "Valid only shortly around 2 AM",
			cases: []timeCase{
				{"1:59 AM too early", atime.MustParseRFC3339("2025-06-12T01:59:00Z"), false, "Before window"},
				{"2:05 AM in range", atime.MustParseRFC3339("2025-06-12T02:05:00Z"), true, "Valid for short report time"},
				{"2:16 AM too late", atime.MustParseRFC3339("2025-06-12T02:16:00Z"), false, "Expired window"},
			},
		},
		{
			name: "Recurring every 3 days",
			rrule: RRuleExtend{
				ROptions: ROptionExtend{
					Freq:      atime.TIMEUNIT_DAILY,
					Interval:  3,
					BeginTime: atime.MustParsePtrRFC3339("2025-06-01T00:00:00Z"),
					ByHour:    []int{0},
				},
				Duration:     24,
				DurationUnit: atime.TIMEUNIT_HOURLY,
			},
			description: "Access every 3 days starting June 1",
			cases: []timeCase{
				{"June 1 — yes", atime.MustParseRFC3339("2025-06-01T12:00:00Z"), true, "Exact recurrence"},
				{"June 2 — no", atime.MustParseRFC3339("2025-06-02T12:00:00Z"), false, "Non-matching day"},
				{"June 4 — yes", atime.MustParseRFC3339("2025-06-04T12:00:00Z"), true, "Third day recurrence"},
			},
		},
		{
			name: "End-of-day expiry window (access ends at 23:59)",
			rrule: RRuleExtend{
				ROptions: ROptionExtend{
					Freq:      atime.TIMEUNIT_DAILY,
					Interval:  1,
					BeginTime: atime.MustParsePtrRFC3339("2025-06-11T00:00:00Z"),
					ByHour:    []int{0},
					ByMinute:  []int{0},
					RRIncType: RRULE_INC_TYPE_EXCLUSIVE, // needed for "00:00 next day — no"
				},
				Duration:     1440,
				DurationUnit: atime.TIMEUNIT_MINUTELY, // 1440 minutes = full 24 hours
			},
			description: "Each hour lasts 1 minute; access ends at 23:59",
			cases: []timeCase{
				{"23:58 OK", atime.MustParseRFC3339("2025-06-11T23:58:00Z"), true, "valid moment"},
				{"23:58:30 OK", atime.MustParseRFC3339("2025-06-11T23:58:30Z"), true, "valid moment with seconds"},
				{"23:59 OK", atime.MustParseRFC3339("2025-06-11T23:59:00Z"), true, "Last valid moment"},
				{"23:59:59 OK", atime.MustParseRFC3339("2025-06-11T23:59:59Z"), true, "Last valid moment with seconds"},
				{"00:00 next day — no", atime.MustParseRFC3339("2025-06-12T00:00:00Z"), false, "Out of window"},
			},
		},
		{
			name: "Exact EndDate match",
			rrule: RRuleExtend{
				ROptions: ROptionExtend{
					Freq:      atime.TIMEUNIT_DAILY,
					Interval:  1,
					BeginTime: atime.MustParsePtrRFC3339("2025-06-01T00:00:00Z"),
					ByHour:    []int{0},
					RRIncType: RRULE_INC_TYPE_EXCLUSIVE, // needed for "At EndDate exact"
				},
				EndDate:      atime.MustParsePtrRFC3339("2025-06-15T00:00:00Z"),
				Duration:     60,
				DurationUnit: atime.TIMEUNIT_MINUTELY,
			},
			description: "Valid exactly until EndDate",
			cases: []timeCase{
				{"Just before EndDate", atime.MustParseRFC3339("2025-06-14T00:59:00Z"), true, "Final valid moment"},
				{"At EndDate exact", atime.MustParseRFC3339("2025-06-15T00:00:00Z"), false, "Just after EndDate"},
			},
		},
		{
			name: "ByMinute overlaps midnight",
			rrule: RRuleExtend{
				ROptions: ROptionExtend{
					Freq:      atime.TIMEUNIT_DAILY,
					Interval:  1,
					BeginTime: atime.MustParsePtrRFC3339("2025-06-01T00:00:00Z"),
					ByHour:    []int{23},
					ByMinute:  []int{59},
					RRIncType: RRULE_INC_TYPE_INCLUSIVE,
				},
				Duration:     2,
				DurationUnit: atime.TIMEUNIT_MINUTELY,
			},
			description: "Ensure 23:59 occurrence doesn't match at 00:01 next day",
			cases: []timeCase{
				{"At 23:59 OK", atime.MustParseRFC3339("2025-06-11T23:59:00Z"), true, "Start of window"},
				{"At 00:00 next day OK", atime.MustParseRFC3339("2025-06-12T00:00:00Z"), true, "Within 2-min window"},
				{"At 00:01 next day NO", atime.MustParseRFC3339("2025-06-12T00:01:00Z"), false, "Outside window"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, tc := range tt.cases {
				t.Run(tc.label, func(t *testing.T) {
					ok, err := tt.rrule.IsBetween(tc.now)
					require.NoError(t, err)
					require.Equal(t, tc.expect, ok, fmt.Sprintf("[%s] %s", tc.label, tt.description))
					if ok != tc.expect {
						t.Errorf(
							"Failed case: %s\nExpected: %v\nGot: %v\nRule logic: %s\nNow: %s\n",
							tc.label, tc.expect, ok, tc.expectMsg, tc.now.Format(time.RFC3339),
						)
					}
				})
			}
		})
	}
}

func TestRRuleExtend_IsBetween_WEEKLY(t *testing.T) {
	type timeCase struct {
		label     string
		now       time.Time
		expect    bool
		expectMsg string
	}

	tests := []struct {
		name        string
		rrule       RRuleExtend
		description string
		cases       []timeCase
	}{
		{
			name: "Weekly Tuesday event at 9 AM for 30 minutes",
			rrule: RRuleExtend{
				ROptions: ROptionExtend{
					Freq:      atime.TIMEUNIT_WEEKLY,
					Interval:  1,
					ByHour:    []int{9},
					ByDay:     atime.TimeWeekdaysToRRuleWeekdays(time.Tuesday),
					BeginTime: atime.MustParsePtrRFC3339("2025-06-01T00:00:00Z"),
				},
				Duration:     30,
				DurationUnit: atime.TIMEUNIT_MINUTELY,
			},
			description: "Basic weekly matching on Tuesday",
			cases: []timeCase{
				{"Tuesday match", atime.MustParseRFC3339("2025-06-24T09:10:00Z"), true, "Within window"},
				{"Monday miss", atime.MustParseRFC3339("2025-06-23T09:10:00Z"), false, "Wrong day"},
				{"Tuesday too early", atime.MustParseRFC3339("2025-06-24T08:59:00Z"), false, "Before window"},
				{"Tuesday too late", atime.MustParseRFC3339("2025-06-24T09:31:00Z"), false, "After window"},
			},
		},
		{
			name: "Multiple weekdays (Mon, Wed, Fri)",
			rrule: RRuleExtend{
				ROptions: ROptionExtend{
					Freq:      atime.TIMEUNIT_WEEKLY,
					Interval:  1,
					ByHour:    []int{8},
					ByDay:     atime.TimeWeekdaysToRRuleWeekdays(time.Monday, time.Wednesday, time.Friday),
					BeginTime: atime.MustParsePtrRFC3339("2025-06-01T00:00:00Z"),
				},
				Duration:     60,
				DurationUnit: atime.TIMEUNIT_MINUTELY,
			},
			description: "Match only on specified weekdays",
			cases: []timeCase{
				{"Monday match", atime.MustParseRFC3339("2025-06-23T08:30:00Z"), true, "Within window"},
				{"Tuesday miss", atime.MustParseRFC3339("2025-06-24T08:30:00Z"), false, "Not scheduled"},
				{"Friday match", atime.MustParseRFC3339("2025-06-27T08:15:00Z"), true, "Within window"},
			},
		},
		{
			name: "Weekly with JoinWindow 5m before and after",
			rrule: RRuleExtend{
				ROptions: ROptionExtend{
					Freq:      atime.TIMEUNIT_WEEKLY,
					Interval:  1,
					ByHour:    []int{10},
					ByDay:     atime.TimeWeekdaysToRRuleWeekdays(time.Thursday),
					BeginTime: atime.MustParsePtrRFC3339("2025-06-01T00:00:00Z"),
				},
				Duration:     10,
				DurationUnit: atime.TIMEUNIT_MINUTELY,
				JoinWindows: JoinWindows{
					{IsBefore: true, Duration: 5, DurationUnit: atime.TIMEUNIT_MINUTELY},
					{IsBefore: false, Duration: 5, DurationUnit: atime.TIMEUNIT_MINUTELY},
				},
			},
			description: "Expand 10m window with JoinWindow",
			cases: []timeCase{
				{"Before grace window", atime.MustParseRFC3339("2025-06-26T09:54:00Z"), false, "Too early"},
				{"Inside pre-window", atime.MustParseRFC3339("2025-06-26T09:56:00Z"), false, "Within JoinWindow before but not in between"},
				{"Valid between", atime.MustParseRFC3339("2025-06-26T10:01:00Z"), true, "Within between"},
				{"End of post-window", atime.MustParseRFC3339("2025-06-26T10:15:00Z"), false, "Just after"},
			},
		},
		{
			name: "Weekly with EndDate",
			rrule: RRuleExtend{
				ROptions: ROptionExtend{
					Freq:      atime.TIMEUNIT_WEEKLY,
					Interval:  1,
					ByHour:    []int{12},
					ByDay:     atime.TimeWeekdaysToRRuleWeekdays(time.Sunday),
					BeginTime: atime.MustParsePtrRFC3339("2025-06-01T00:00:00Z"),
				},
				Duration:     30,
				DurationUnit: atime.TIMEUNIT_MINUTELY,
				EndDate:      atime.MustParsePtrRFC3339("2025-06-23T00:00:00Z"),
			},
			description: "Ensure EndDate is respected",
			cases: []timeCase{
				{"Before EndDate", atime.MustParseRFC3339("2025-06-22T12:15:00Z"), true, "Last valid week"},
				{"After EndDate", atime.MustParseRFC3339("2025-06-29T12:00:00Z"), false, "Too late"},
			},
		},
		{
			name: "Weekly deny rule",
			rrule: RRuleExtend{
				ROptions: ROptionExtend{
					Freq:      atime.TIMEUNIT_WEEKLY,
					Interval:  1,
					ByHour:    []int{15},
					ByDay:     atime.TimeWeekdaysToRRuleWeekdays(time.Friday),
					BeginTime: atime.MustParsePtrRFC3339("2025-06-01T00:00:00Z"),
				},
				Duration:     60,
				DurationUnit: atime.TIMEUNIT_MINUTELY,
				IsDeny:       true,
			},
			description: "Access should be denied in this window",
			cases: []timeCase{
				{"Friday match", atime.MustParseRFC3339("2025-06-27T15:30:00Z"), false, "Deny rule applies"},
				{"Other day OK", atime.MustParseRFC3339("2025-06-26T15:30:00Z"), true, "Thursday is not denied"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, c := range tt.cases {
				t.Run(c.label, func(t *testing.T) {
					result, err := tt.rrule.IsBetween(c.now)
					require.NoError(t, err)
					require.Equal(t, c.expect, result, c.expectMsg)
				})
			}
		})
	}
}

func TestRRuleExtend_IsBetween_HOURLY(t *testing.T) {
	type timeCase struct {
		label     string
		now       time.Time
		expect    bool
		expectMsg string
	}

	tests := []struct {
		name        string
		rrule       RRuleExtend
		description string
		cases       []timeCase
	}{
		{
			name: "Simple hourly rule at minute 00 for 10 minutes",
			rrule: RRuleExtend{
				ROptions: ROptionExtend{
					Freq:      atime.TIMEUNIT_HOURLY,
					Interval:  1,
					ByMinute:  []int{0},
					BeginTime: atime.MustParsePtrRFC3339("2025-06-21T00:00:00Z"),
					// RRIncType:    RRULE_INC_TYPE_INCLUSIVE, // makes needed for "Exactly 9:00"
				},
				Duration:     10,
				DurationUnit: atime.TIMEUNIT_MINUTELY,
			},
			description: "Match at every hour's start",
			cases: []timeCase{
				{"Exactly 9:00", atime.MustParseRFC3339("2025-06-21T09:00:00Z"), true, "Start of window"},
				{"9:05 OK", atime.MustParseRFC3339("2025-06-21T09:05:00Z"), true, "Inside window"},
				{"9:10 miss", atime.MustParseRFC3339("2025-06-21T09:10:00Z"), false, "Just outside"},
			},
		},
		{
			name: "Hourly with JoinWindow before and after",
			rrule: RRuleExtend{
				ROptions: ROptionExtend{
					Freq:      atime.TIMEUNIT_HOURLY,
					Interval:  1,
					ByMinute:  []int{15},
					BeginTime: atime.MustParsePtrRFC3339("2025-06-21T00:00:00Z"),
					// RRIncType: RRULE_INC_TYPE_INCLUSIVE, // needed for Valid between at "2025-06-26T10:15:00Z"
				},
				Duration:     5,
				DurationUnit: atime.TIMEUNIT_MINUTELY,
				JoinWindows: JoinWindows{
					{IsBefore: true, Duration: 2, DurationUnit: atime.TIMEUNIT_MINUTELY},
					{IsBefore: false, Duration: 3, DurationUnit: atime.TIMEUNIT_MINUTELY},
				},
			},
			description: "Grace on both sides",
			cases: []timeCase{
				{"Before grace window", atime.MustParseRFC3339("2025-06-26T09:54:00Z"), false, "Too early — join window not included"},
				{"Inside pre-window", atime.MustParseRFC3339("2025-06-26T09:56:00Z"), false, "Pre-window not part of IsBetween"},
				{"Valid between", atime.MustParseRFC3339("2025-06-26T10:15:00Z"), true, "Inside scheduled interval"},
				{"End of post-window", atime.MustParseRFC3339("2025-06-26T10:30:00Z"), false, "Post-window excluded"},
			},
		},
		{
			name: "Hourly with EndDate",
			rrule: RRuleExtend{
				ROptions: ROptionExtend{
					Freq:      atime.TIMEUNIT_HOURLY,
					Interval:  1,
					ByMinute:  []int{45},
					BeginTime: atime.MustParsePtrRFC3339("2025-06-21T00:00:00Z"),
					// RRIncType:    RRULE_INC_TYPE_INCLUSIVE, // works with or without
				},
				Duration:     10,
				DurationUnit: atime.TIMEUNIT_MINUTELY,
				EndDate:      atime.MustParsePtrRFC3339("2025-06-21T12:00:00Z"),
			},
			description: "No matches past EndDate",
			cases: []timeCase{
				{"11:50 OK", atime.MustParseRFC3339("2025-06-21T11:50:00Z"), true, "Before EndDate"},
				{"12:45 too late", atime.MustParseRFC3339("2025-06-21T12:45:00Z"), false, "After EndDate"},
			},
		},
		{
			name: "Hourly deny rule",
			rrule: RRuleExtend{
				ROptions: ROptionExtend{
					Freq:      atime.TIMEUNIT_HOURLY,
					Interval:  1,
					ByMinute:  []int{0},
					BeginTime: atime.MustParsePtrRFC3339("2025-06-21T00:00:00Z"),
					// RRIncType:    RRULE_INC_TYPE_INCLUSIVE, // needed for "9:00 hit"
				},
				Duration:     15,
				DurationUnit: atime.TIMEUNIT_MINUTELY,
				IsDeny:       true,
			},
			description: "Should return false due to deny",
			cases: []timeCase{
				{"9:00 hit", atime.MustParseRFC3339("2025-06-21T09:00:00Z"), false, "Denied at start"},
				{"9:10 hit", atime.MustParseRFC3339("2025-06-21T09:10:00Z"), false, "Still denied"},
				{"9:16 safe", atime.MustParseRFC3339("2025-06-21T09:16:00Z"), true, "Outside deny"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, c := range tt.cases {
				t.Run(c.label, func(t *testing.T) {
					result, err := tt.rrule.IsBetween(c.now)
					require.NoError(t, err)
					require.Equal(t, c.expect, result, c.expectMsg)
				})
			}
		})
	}
}

func TestRRuleExtend_IsBetween_MINUTELY(t *testing.T) {
	type timeCase struct {
		label     string
		now       time.Time
		expect    bool
		expectMsg string
	}

	tests := []struct {
		name        string
		rrule       RRuleExtend
		description string
		cases       []timeCase
	}{
		{
			name: "Every 15 minutes starting at 10:00",
			rrule: RRuleExtend{
				ROptions: ROptionExtend{
					Freq:      atime.TIMEUNIT_MINUTELY,
					Interval:  15,
					BeginTime: atime.MustParsePtrRFC3339("2025-06-21T10:00:00Z"),
				},
				Duration:     5,
				DurationUnit: atime.TIMEUNIT_MINUTELY,
			},
			description: "Should match every 15 minutes with a 5-minute active window",
			cases: []timeCase{
				{"10:00 hit", atime.MustParseRFC3339("2025-06-21T10:00:00Z"), true, "Exact match"},
				{"10:04 in window", atime.MustParseRFC3339("2025-06-21T10:04:00Z"), true, "Within window"},
				{"10:05 outside", atime.MustParseRFC3339("2025-06-21T10:05:00Z"), false, "Just after window"},
				{"10:15 hit", atime.MustParseRFC3339("2025-06-21T10:15:00Z"), true, "Next interval"},
			},
		},
		{
			name: "Minutely recurrence with JoinWindow — still only match core",
			rrule: RRuleExtend{
				ROptions: ROptionExtend{
					Freq:      atime.TIMEUNIT_MINUTELY,
					Interval:  30,
					BeginTime: atime.MustParsePtrRFC3339("2025-06-21T12:00:00Z"),
				},
				Duration:     1,
				DurationUnit: atime.TIMEUNIT_MINUTELY,
				JoinWindows: JoinWindows{
					{IsBefore: true, Duration: 5, DurationUnit: atime.TIMEUNIT_MINUTELY},
					{IsBefore: false, Duration: 5, DurationUnit: atime.TIMEUNIT_MINUTELY},
				},
			},
			description: "JoinWindow present but not affecting IsBetween",
			cases: []timeCase{
				{"Too early", atime.MustParseRFC3339("2025-06-21T11:54:00Z"), false, "Before recurrence"},
				{"Inside join window", atime.MustParseRFC3339("2025-06-21T11:56:00Z"), false, "Join window not part of IsBetween"},
				{"Exact match", atime.MustParseRFC3339("2025-06-21T12:00:00Z"), true, "Inside core recurrence"},
				{"After post window", atime.MustParseRFC3339("2025-06-21T12:06:00Z"), false, "Out of core window"},
			},
		},
		{
			name: "Minutely recurrence ending at fixed time",
			rrule: RRuleExtend{
				ROptions: ROptionExtend{
					Freq:      atime.TIMEUNIT_MINUTELY,
					Interval:  10,
					BeginTime: atime.MustParsePtrRFC3339("2025-06-21T14:00:00Z"),
				},
				Duration:     5,
				DurationUnit: atime.TIMEUNIT_MINUTELY,
				EndDate:      atime.MustParsePtrRFC3339("2025-06-21T14:30:00Z"),
			},
			description: "Should not match after EndDate",
			cases: []timeCase{
				{"Final valid window", atime.MustParseRFC3339("2025-06-21T14:24:59Z"), true, "Still valid"},
				{"Final invalid window", atime.MustParseRFC3339("2025-06-21T14:25:00Z"), false, "Invalid"},
				{"After EndDate", atime.MustParseRFC3339("2025-06-21T14:31:00Z"), false, "Out of range"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, c := range tt.cases {
				t.Run(c.label, func(t *testing.T) {
					ok, err := tt.rrule.IsBetween(c.now)
					require.NoError(t, err)
					require.Equal(t, c.expect, ok, c.expectMsg)
				})
			}
		})
	}
}

func TestRRuleExtend_IsBetween_SECONDLY(t *testing.T) {
	type timeCase struct {
		label     string
		now       time.Time
		expect    bool
		expectMsg string
	}

	tests := []struct {
		name        string
		rrule       RRuleExtend
		description string
		cases       []timeCase
	}{
		{
			name: "Stock tick every 15 seconds",
			rrule: RRuleExtend{
				ROptions: ROptionExtend{
					Freq:      atime.TIMEUNIT_SECONDLY,
					Interval:  15,
					BeginTime: atime.MustParsePtrRFC3339("2025-06-21T14:00:00Z"),
				},
				Duration:     5,
				DurationUnit: atime.TIMEUNIT_SECONDLY,
			},
			description: "Allow only within 5-second window every 15 seconds",
			cases: []timeCase{
				{"Exact match", atime.MustParseRFC3339("2025-06-21T14:00:00Z"), true, "Start of tick window"},
				{"Within tick window", atime.MustParseRFC3339("2025-06-21T14:00:04Z"), true, "Within 5s window"},
				{"Too late", atime.MustParseRFC3339("2025-06-21T14:00:06Z"), false, "Outside tick window"},
			},
		},
		{
			name: "Security check every 30s with 10s allowance",
			rrule: RRuleExtend{
				ROptions: ROptionExtend{
					Freq:      atime.TIMEUNIT_SECONDLY,
					Interval:  30,
					BeginTime: atime.MustParsePtrRFC3339("2025-06-21T15:00:00Z"),
				},
				Duration:     10,
				DurationUnit: atime.TIMEUNIT_SECONDLY,
			},
			description: "Verify within 10 seconds of each ping",
			cases: []timeCase{
				{"First window match", atime.MustParseRFC3339("2025-06-21T15:00:05Z"), true, "In allowance"},
				{"Missed window", atime.MustParseRFC3339("2025-06-21T15:00:11Z"), false, "Too late"},
				{"Second tick match", atime.MustParseRFC3339("2025-06-21T15:00:30Z"), true, "Second cycle"},
			},
		},
		{
			name: "Short-lived access burst",
			rrule: RRuleExtend{
				ROptions: ROptionExtend{
					Freq:      atime.TIMEUNIT_SECONDLY,
					Interval:  60,
					BeginTime: atime.MustParsePtrRFC3339("2025-06-21T16:00:00Z"),
				},
				Duration:     3,
				DurationUnit: atime.TIMEUNIT_SECONDLY,
				EndDate:      atime.MustParsePtrRFC3339("2025-06-21T16:01:30Z"),
			},
			description: "Access for 3 seconds every minute",
			cases: []timeCase{
				{"First burst", atime.MustParseRFC3339("2025-06-21T16:00:01Z"), true, "Valid burst window"},
				{"Expired range", atime.MustParseRFC3339("2025-06-21T16:01:35Z"), false, "After end date"},
			},
		},
		{
			name: "Edge alignment: 5s pulse every 20s",
			rrule: RRuleExtend{
				ROptions: ROptionExtend{
					Freq:      atime.TIMEUNIT_SECONDLY,
					Interval:  20,
					BeginTime: atime.MustParsePtrRFC3339("2025-06-21T17:00:00Z"),
				},
				Duration:     5,
				DurationUnit: atime.TIMEUNIT_SECONDLY,
			},
			description: "Match only at edge of pulses",
			cases: []timeCase{
				{"Start exact", atime.MustParseRFC3339("2025-06-21T17:00:00Z"), true, "Exact alignment"},
				{"End edge", atime.MustParseRFC3339("2025-06-21T17:00:05Z"), false, "Just outside window"},
				{"Mid pulse", atime.MustParseRFC3339("2025-06-21T17:00:03Z"), true, "Within window"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, c := range tt.cases {
				t.Run(c.label, func(t *testing.T) {
					result, err := tt.rrule.IsBetween(c.now)
					require.NoError(t, err)
					require.Equal(t, c.expect, result, c.expectMsg)
				})
			}
		})
	}
}

// TestRRuleExtend_IsBetween has the original base-line tests before Freq-specific.
func TestRRuleExtend_IsBetween(t *testing.T) {
	tests := []struct {
		name        string
		description string
		rrule       RRuleExtend
		now         time.Time
		expect      bool
		expectMsg   string
	}{
		{
			name:      "Exact match at scheduled time",
			expectMsg: "Match occurs exactly at 10:00 AM UTC every day; 'now' is 10:00 AM → match.",
			rrule: RRuleExtend{
				ROptions: ROptionExtend{
					Freq: atime.TIMEUNIT_DAILY, Interval: 1,
					BeginTime: atime.MustParsePtrRFC3339("2025-06-20T10:00:00Z"),
					ByHour:    []int{10},
				},
				Duration: 1, DurationUnit: atime.TIMEUNIT_MINUTELY,
			},
			now:    atime.MustParseRFC3339("2025-06-20T10:00:00Z"),
			expect: true,
		},
		{
			name:      "One minute after exact time — should fail",
			expectMsg: "Match at 10:00 AM UTC with 1-minute window; 'now' is 10:01 AM → outside window.",
			rrule: RRuleExtend{
				ROptions: ROptionExtend{
					Freq: atime.TIMEUNIT_DAILY, Interval: 1,
					BeginTime: atime.MustParsePtrRFC3339("2025-06-20T10:00:00Z"),
					ByHour:    []int{10},
				},
				Duration: 1, DurationUnit: atime.TIMEUNIT_MINUTELY,
			},
			now:    atime.MustParseRFC3339("2025-06-20T10:01:00Z"),
			expect: false,
		},
		{
			name:      "Within join window before",
			expectMsg: "Join window of 5 minutes before start; 09:55 AM is too early to count as a match.",
			rrule: RRuleExtend{
				ROptions: ROptionExtend{
					Freq: atime.TIMEUNIT_DAILY, Interval: 1,
					BeginTime: atime.MustParsePtrRFC3339("2025-06-20T10:00:00Z"),
					ByHour:    []int{10},
				},
				JoinWindows: JoinWindows{
					{IsBefore: true, Duration: 5, DurationUnit: atime.TIMEUNIT_MINUTELY},
				},
				Duration: 1, DurationUnit: atime.TIMEUNIT_MINUTELY,
			},
			now:    atime.MustParseRFC3339("2025-06-20T09:55:00Z"),
			expect: false,
		},
		{
			name:      "Outside join window (too late)",
			expectMsg: "Start at 10:00 + 10 min after window ends at 10:10; 10:16 is outside → no match.",
			rrule: RRuleExtend{
				ROptions: ROptionExtend{
					Freq: atime.TIMEUNIT_DAILY, Interval: 1,
					BeginTime: atime.MustParsePtrRFC3339("2025-06-20T10:00:00Z"),
					ByHour:    []int{10},
				},
				JoinWindows: JoinWindows{
					{IsBefore: true, Duration: 5, DurationUnit: atime.TIMEUNIT_MINUTELY},
					{IsBefore: false, Duration: 10, DurationUnit: atime.TIMEUNIT_MINUTELY},
				},
				Duration: 1, DurationUnit: atime.TIMEUNIT_MINUTELY,
			},
			now:    atime.MustParseRFC3339("2025-06-20T10:16:00Z"),
			expect: false,
		},
		{
			name:      "StartDate boundary hit",
			expectMsg: "StartDate = 6/1, recurrence starts 6/15; 6/15 is within bounds → match.",
			rrule: RRuleExtend{
				ROptions: ROptionExtend{
					Freq: atime.TIMEUNIT_DAILY, Interval: 1,
					BeginTime: atime.MustParsePtrRFC3339("2025-06-15T10:00:00Z"),
					ByHour:    []int{10},
				},
				StartDate: atime.MustParsePtrRFC3339("2025-06-01T00:00:00Z"),
				Duration:  1, DurationUnit: atime.TIMEUNIT_MINUTELY,
			},
			now:    atime.MustParseRFC3339("2025-06-15T10:00:00Z"),
			expect: true,
		},
		{
			name:      "Before StartDate",
			expectMsg: "StartDate is 6/10; 6/9 is too early, so should not match.",
			rrule: RRuleExtend{
				ROptions: ROptionExtend{
					Freq: atime.TIMEUNIT_DAILY, Interval: 1,
					BeginTime: atime.MustParsePtrRFC3339("2025-06-10T10:00:00Z"),
					ByHour:    []int{10},
				},
				StartDate: atime.MustParsePtrRFC3339("2025-06-10T00:00:00Z"),
				Duration:  1, DurationUnit: atime.TIMEUNIT_MINUTELY,
			},
			now:    atime.MustParseRFC3339("2025-06-09T10:00:00Z"),
			expect: false,
		},
		{
			name:      "After EndDate",
			expectMsg: "EndDate = 6/15; 'now' is 6/16, which is after window → no match.",
			rrule: RRuleExtend{
				ROptions: ROptionExtend{
					Freq: atime.TIMEUNIT_DAILY, Interval: 1,
					BeginTime: atime.MustParsePtrRFC3339("2025-06-15T10:00:00Z"),
					ByHour:    []int{10},
				},
				EndDate:  atime.MustParsePtrRFC3339("2025-06-15T00:00:00Z"),
				Duration: 1, DurationUnit: atime.TIMEUNIT_MINUTELY,
			},
			now:    atime.MustParseRFC3339("2025-06-16T10:00:00Z"),
			expect: false,
		},
		{
			name:      "Count limit in range",
			expectMsg: "3 daily recurrences allowed starting 6/1; 6/3 is still within → match.",
			rrule: RRuleExtend{
				ROptions: ROptionExtend{
					Freq: atime.TIMEUNIT_DAILY, Interval: 1,
					BeginTime: atime.MustParsePtrRFC3339("2025-06-01T10:00:00Z"),
					ByHour:    []int{10},
					Count:     3,
				},
				Duration: 1, DurationUnit: atime.TIMEUNIT_MINUTELY,
			},
			now:    atime.MustParseRFC3339("2025-06-03T10:00:00Z"),
			expect: true,
		},
		{
			name:      "Count limit exceeded",
			expectMsg: "Count=3 from 6/1; 6/5 is beyond the third recurrence → no match.",
			rrule: RRuleExtend{
				ROptions: ROptionExtend{
					Freq: atime.TIMEUNIT_DAILY, Interval: 1,
					BeginTime: atime.MustParsePtrRFC3339("2025-06-01T10:00:00Z"),
					ByHour:    []int{10},
					Count:     3,
				},
				Duration: 1, DurationUnit: atime.TIMEUNIT_MINUTELY,
			},
			now:    atime.MustParseRFC3339("2025-06-05T10:00:00Z"),
			expect: false,
		},
		{
			name:      "Weekly Monday match",
			expectMsg: "Weekly on Monday at 9 AM; 6/23 is a Monday at 9 → match.",
			rrule: RRuleExtend{
				ROptions: ROptionExtend{
					Freq: atime.TIMEUNIT_WEEKLY, Interval: 1,
					BeginTime: atime.MustParsePtrRFC3339("2025-06-23T09:00:00Z"),
					ByHour:    []int{9},
					ByDay:     []rrule.Weekday{rrule.MO},
				},
				Duration: 1, DurationUnit: atime.TIMEUNIT_MINUTELY,
			},
			now:    atime.MustParseRFC3339("2025-06-23T09:00:00Z"),
			expect: true,
		},
		{
			name:      "Weekly mismatch (Tuesday)",
			expectMsg: "Rule is Monday only; 6/24 is Tuesday → mismatch.",
			rrule: RRuleExtend{
				ROptions: ROptionExtend{
					Freq: atime.TIMEUNIT_WEEKLY, Interval: 1,
					BeginTime: atime.MustParsePtrRFC3339("2025-06-23T09:00:00Z"),
					ByHour:    []int{9},
					ByDay:     []rrule.Weekday{rrule.MO},
				},
				Duration: 1, DurationUnit: atime.TIMEUNIT_MINUTELY,
			},
			now:    atime.MustParseRFC3339("2025-06-24T09:00:00Z"),
			expect: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, err := tt.rrule.IsBetween(tt.now)
			require.NoError(t, err)
			require.Equal(t, tt.expect, ok, tt.description)
			if ok != tt.expect {
				t.Errorf(
					"Failed test: %s\nExpected: %v\nGot: %v\nRule logic: %s\nNow: %s\n",
					tt.name, tt.expect, ok, tt.expectMsg, tt.now.Format(time.RFC3339),
				)
			}
		})
	}
}

func TestRRuleExtend_Validate(t *testing.T) {
	tests := []struct {
		name        string
		rrule       RRuleExtend
		expectError bool
	}{
		{
			name: "Valid daily rule",
			rrule: RRuleExtend{
				ROptions: ROptionExtend{
					Freq:     atime.TIMEUNIT_DAILY,
					Interval: 1,
				},
			},
			expectError: false,
		},
		{
			name: "IsAnyTime skips validation",
			rrule: RRuleExtend{
				IsAnyTime: true,
				ROptions: ROptionExtend{
					Freq:     "INVALID", // ignored
					Interval: 0,         // ignored
				},
			},
			expectError: false,
		},
		{
			name: "Invalid frequency",
			rrule: RRuleExtend{
				ROptions: ROptionExtend{
					Freq:     "banana",
					Interval: 1,
				},
			},
			expectError: true,
		},
		{
			name: "Interval too low",
			rrule: RRuleExtend{
				ROptions: ROptionExtend{
					Freq:     atime.TIMEUNIT_DAILY,
					Interval: 0,
				},
			},
			expectError: true,
		},
		{
			name: "Valid byhour, byminute, bysecond",
			rrule: RRuleExtend{
				ROptions: ROptionExtend{
					Freq:     atime.TIMEUNIT_DAILY,
					Interval: 1,
					ByHour:   []int{0, 23},
					ByMinute: []int{0, 59},
					BySecond: []int{0, 59},
				},
			},
			expectError: false,
		},
		{
			name: "Invalid byhour",
			rrule: RRuleExtend{
				ROptions: ROptionExtend{
					Freq:     atime.TIMEUNIT_DAILY,
					Interval: 1,
					ByHour:   []int{24}, // invalid
				},
			},
			expectError: true,
		},
		{
			name: "Invalid bymonthday",
			rrule: RRuleExtend{
				ROptions: ROptionExtend{
					Freq:       atime.TIMEUNIT_DAILY,
					Interval:   1,
					ByMonthDay: []int{32, -32},
				},
			},
			expectError: true,
		},
		{
			name: "Valid bymonthday negative",
			rrule: RRuleExtend{
				ROptions: ROptionExtend{
					Freq:       atime.TIMEUNIT_DAILY,
					Interval:   1,
					ByMonthDay: []int{-1, -31, 1, 15},
				},
			},
			expectError: false,
		},
		{
			name: "Invalid byyearday",
			rrule: RRuleExtend{
				ROptions: ROptionExtend{
					Freq:      atime.TIMEUNIT_YEARLY,
					Interval:  1,
					ByYearDay: []int{367},
				},
			},
			expectError: true,
		},
		{
			name: "Invalid byweekno",
			rrule: RRuleExtend{
				ROptions: ROptionExtend{
					Freq:     atime.TIMEUNIT_WEEKLY,
					Interval: 1,
					ByWeekNo: []int{-54},
				},
			},
			expectError: true,
		},
		{
			name: "Invalid bymonth",
			rrule: RRuleExtend{
				ROptions: ROptionExtend{
					Freq:     atime.TIMEUNIT_YEARLY,
					Interval: 1,
					ByMonth:  []int{0, 13},
				},
			},
			expectError: true,
		},
		{
			name: "Invalid bysetpos",
			rrule: RRuleExtend{
				ROptions: ROptionExtend{
					Freq:     atime.TIMEUNIT_MONTHLY,
					Interval: 1,
					BySetPos: []int{-400},
				},
			},
			expectError: true,
		},
		{
			name: "Valid byday with position",
			rrule: RRuleExtend{
				ROptions: ROptionExtend{
					Freq:     atime.TIMEUNIT_MONTHLY,
					Interval: 1,
					ByDay: []rrule.Weekday{
						rrule.MO.Nth(1),
						rrule.FR.Nth(-1),
					},
				},
			},
			expectError: false,
		},
		{
			name: "Invalid byday with position out of range",
			rrule: RRuleExtend{
				ROptions: ROptionExtend{
					Freq:     atime.TIMEUNIT_MONTHLY,
					Interval: 1,
					ByDay: []rrule.Weekday{
						rrule.TU.Nth(100),
					},
				},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.rrule.Validate()
			if tt.expectError {
				require.Error(t, err, "expected an error but got nil")
			} else {
				require.NoError(t, err, "expected no error but got one")
			}
		})
	}
}

func TestRRuleExtends_Evaluate(t *testing.T) {
	type timeCase struct {
		label     string
		now       time.Time
		expect    bool
		expectMsg string
		//expectJoin bool // true if we expect MatchJoinWindow to match
	}

	tests := []struct {
		name        string
		rrules      RRuleExtends
		description string
		cases       []timeCase
	}{
		{
			name:        "Single allow rule match",
			description: "Allow rule with matching window at 10:05",
			rrules: RRuleExtends{
				{
					ROptions: ROptionExtend{
						Freq:      atime.TIMEUNIT_DAILY,
						BeginTime: atime.MustParsePtrRFC3339("2025-06-20T10:00:00Z"),
						ByHour:    []int{10},
					},
					IsDeny:       false,
					Priority:     1,
					Duration:     1,
					DurationUnit: atime.TIMEUNIT_MINUTELY,
					JoinWindows: JoinWindows{
						{IsBefore: true, Duration: 10, DurationUnit: atime.TIMEUNIT_MINUTELY},
						{IsBefore: false, Duration: 10, DurationUnit: atime.TIMEUNIT_MINUTELY},
					},
				},
			},
			cases: []timeCase{
				{"Hit window", atime.MustParseRFC3339("2025-06-20T10:00:30Z"), true, "Within core window"},
				{"Hit window within join", atime.MustParseRFC3339("2025-06-20T10:05:00Z"), false, "Should allow within join window"},
				{"Too late", atime.MustParseRFC3339("2025-06-20T10:20:01Z"), false, "Outside window"},
			},
		},
		{
			name:        "Allow and deny rule; deny wins",
			description: "Conflict with deny of higher priority",
			rrules: RRuleExtends{
				{
					ROptions: ROptionExtend{
						Freq:      atime.TIMEUNIT_DAILY,
						BeginTime: atime.MustParsePtrRFC3339("2025-06-20T10:00:00Z"),
						ByHour:    []int{10},
					},
					IsDeny:       false,
					Priority:     5,
					Duration:     1,
					DurationUnit: atime.TIMEUNIT_MINUTELY,
					JoinWindows: JoinWindows{
						{IsBefore: true, Duration: 10, DurationUnit: atime.TIMEUNIT_MINUTELY},
						{IsBefore: false, Duration: 10, DurationUnit: atime.TIMEUNIT_MINUTELY},
					},
				},
				{
					ROptions: ROptionExtend{
						Freq:      atime.TIMEUNIT_DAILY,
						BeginTime: atime.MustParsePtrRFC3339("2025-06-20T10:00:00Z"),
						ByHour:    []int{10},
					},
					IsDeny:       true,
					Priority:     10,
					Duration:     1,
					DurationUnit: atime.TIMEUNIT_MINUTELY,
					JoinWindows: JoinWindows{
						{IsBefore: true, Duration: 10, DurationUnit: atime.TIMEUNIT_MINUTELY},
						{IsBefore: false, Duration: 10, DurationUnit: atime.TIMEUNIT_MINUTELY},
					},
				},
			},
			cases: []timeCase{
				{"Hit window", atime.MustParseRFC3339("2025-06-20T10:05:00Z"), false, "Should deny due to higher priority"},
			},
		},
		{
			name:        "No matching rules",
			description: "Default deny when no rule matches",
			rrules: RRuleExtends{
				{
					ROptions: ROptionExtend{
						Freq:      atime.TIMEUNIT_DAILY,
						BeginTime: atime.MustParsePtrRFC3339("2025-06-20T07:00:00Z"),
						ByHour:    []int{7},
					},
					IsDeny:       false,
					Priority:     1,
					Duration:     1,
					DurationUnit: atime.TIMEUNIT_MINUTELY,
				},
			},
			cases: []timeCase{
				{"Unmatched time", atime.MustParseRFC3339("2025-06-20T10:05:00Z"), false, "Should deny when unmatched"},
			},
		},
		{
			name:        "Empty rules deny by default",
			description: "Should deny if no rules are present",
			rrules:      RRuleExtends{},
			cases: []timeCase{
				{"Any time", atime.MustParseRFC3339("2025-06-20T10:00:00Z"), false, "Empty rules = deny"},
			},
		},
		{
			name:        "Equal priority, deny precedes allow",
			description: "With equal priority, deny evaluated first wins",
			rrules: RRuleExtends{
				{
					IsDeny:   true,
					Priority: 10,
					ROptions: ROptionExtend{
						Freq:      atime.TIMEUNIT_DAILY,
						BeginTime: atime.MustParsePtrRFC3339("2025-06-20T10:00:00Z"),
						ByHour:    []int{10},
					},
					Duration:     1,
					DurationUnit: atime.TIMEUNIT_MINUTELY,
				},
				{
					IsDeny:   false,
					Priority: 10,
					ROptions: ROptionExtend{
						Freq:      atime.TIMEUNIT_DAILY,
						BeginTime: atime.MustParsePtrRFC3339("2025-06-20T10:00:00Z"),
						ByHour:    []int{10},
					},
					Duration:     1,
					DurationUnit: atime.TIMEUNIT_MINUTELY,
				},
			},
			cases: []timeCase{
				{"Exact match", atime.MustParseRFC3339("2025-06-20T10:00:00Z"), false, "Deny wins when equal priority"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, c := range tt.cases {
				t.Run(c.label, func(t *testing.T) {
					result, err := tt.rrules.Evaluate(c.now)
					require.NoError(t, err)
					require.Equal(t, c.expect, result, c.expectMsg)
				})
			}
		})
	}
}

// Mock evaluator
type mockEvaluator struct {
	preAllow    RREvaluatorResultType
	preAllowerr error
	allow       bool
	err         error
}

func (m *mockEvaluator) IsPreAllowed(_ time.Time, _ ageo.GeoInfo) (RREvaluatorResultType, error) {
	return m.preAllow, m.preAllowerr
}

func (m *mockEvaluator) IsAllowed(_ time.Time, _ ageo.GeoInfo) error {
	if m.err != nil {
		return m.err
	}
	if m.allow {
		return nil
	}
	return errors.New("rejected by evaluator")
}

func TestRRuleExtends_EvaluateWithOptions(t *testing.T) {
	type timeCase struct {
		label     string
		now       time.Time
		expect    bool
		expectMsg string
	}

	tests := []struct {
		name        string
		rrules      RRuleExtends
		description string
		eval        IRRuleEvaluator
		geo         ageo.GeoInfo
		cases       []timeCase
	}{
		{
			name:        "Allow match with no evaluator",
			description: "Simple allow rule within join window",
			rrules: RRuleExtends{
				{
					ROptions: ROptionExtend{
						Freq:      atime.TIMEUNIT_DAILY,
						BeginTime: atime.MustParsePtrRFC3339("2025-06-20T10:00:00Z"),
						ByHour:    []int{10},
					},
					IsDeny:       false,
					Priority:     1,
					Duration:     1,
					DurationUnit: atime.TIMEUNIT_MINUTELY,
					JoinWindows: JoinWindows{
						{IsBefore: true, Duration: 10, DurationUnit: atime.TIMEUNIT_MINUTELY},
						{IsBefore: false, Duration: 10, DurationUnit: atime.TIMEUNIT_MINUTELY},
					},
				},
			},
			cases: []timeCase{
				{"Within window begin", atime.MustParseRFC3339("2025-06-20T10:00:00Z"), true, "Allow match inside join window-begin"},
				{"Within window end", atime.MustParseRFC3339("2025-06-20T10:00:59Z"), true, "Allow match inside join window-end"},
				{"Too late", atime.MustParseRFC3339("2025-06-20T10:21:00Z"), false, "Miss after post window"},
			},
		},
		{
			name:        "Deny match with no evaluator",
			description: "Simple deny rule blocking access",
			rrules: RRuleExtends{
				{
					ROptions: ROptionExtend{
						Freq:      atime.TIMEUNIT_DAILY,
						BeginTime: atime.MustParsePtrRFC3339("2025-06-20T10:00:00Z"),
						ByHour:    []int{10},
					},
					IsDeny:       true,
					Priority:     1,
					Duration:     1,
					DurationUnit: atime.TIMEUNIT_MINUTELY,
					JoinWindows: JoinWindows{
						{IsBefore: true, Duration: 10, DurationUnit: atime.TIMEUNIT_MINUTELY},
						{IsBefore: false, Duration: 10, DurationUnit: atime.TIMEUNIT_MINUTELY},
					},
				},
			},
			cases: []timeCase{
				{"Hit deny", atime.MustParseRFC3339("2025-06-20T10:05:00Z"), false, "Denied inside join window"},
			},
		},
		{
			name:        "Evaluator denies even though rule matches",
			description: "Evaluator overrides matching allow rule",
			rrules: RRuleExtends{
				{
					ROptions: ROptionExtend{
						Freq:      atime.TIMEUNIT_DAILY,
						BeginTime: atime.MustParsePtrRFC3339("2025-06-20T10:00:00Z"),
						ByHour:    []int{10},
					},
					IsDeny:       false,
					Priority:     1,
					Duration:     1,
					DurationUnit: atime.TIMEUNIT_MINUTELY,
					JoinWindows: JoinWindows{
						{IsBefore: true, Duration: 10, DurationUnit: atime.TIMEUNIT_MINUTELY},
						{IsBefore: false, Duration: 10, DurationUnit: atime.TIMEUNIT_MINUTELY},
					},
				},
			},
			eval: &mockEvaluator{allow: false},
			cases: []timeCase{
				{"Match but evaluator blocks", atime.MustParseRFC3339("2025-06-20T10:05:00Z"), false, "Evaluator override to deny"},
			},
		},
		{
			name:        "Evaluator allows access",
			description: "Evaluator grants access even without match",
			rrules: RRuleExtends{
				{
					ROptions: ROptionExtend{
						Freq:      atime.TIMEUNIT_DAILY,
						BeginTime: atime.MustParsePtrRFC3339("2025-06-20T10:00:00Z"),
						ByHour:    []int{10},
					},
					IsDeny:       false,
					Priority:     1,
					Duration:     1,
					DurationUnit: atime.TIMEUNIT_MINUTELY,
				},
			},
			eval: &mockEvaluator{preAllow: RREVALUATOR_RESULTTYPE_ALLOW},
			cases: []timeCase{
				{"Pre-evaluator wins", atime.MustParseRFC3339("2025-06-20T10:05:00Z"), true, "Evaluator grants access"},
			},
		},
		{
			name:        "Allow and deny, deny wins",
			description: "Matching allow and deny, deny higher priority",
			rrules: RRuleExtends{
				{
					ROptions: ROptionExtend{
						Freq:      atime.TIMEUNIT_DAILY,
						BeginTime: atime.MustParsePtrRFC3339("2025-06-20T10:00:00Z"),
						ByHour:    []int{10},
					},
					IsDeny:       false,
					Priority:     1,
					Duration:     1,
					DurationUnit: atime.TIMEUNIT_MINUTELY,
					JoinWindows: JoinWindows{
						{IsBefore: true, Duration: 10, DurationUnit: atime.TIMEUNIT_MINUTELY},
						{IsBefore: false, Duration: 10, DurationUnit: atime.TIMEUNIT_MINUTELY},
					},
				},
				{
					ROptions: ROptionExtend{
						Freq:      atime.TIMEUNIT_DAILY,
						BeginTime: atime.MustParsePtrRFC3339("2025-06-20T10:00:00Z"),
						ByHour:    []int{10},
					},
					IsDeny:       true,
					Priority:     2,
					Duration:     1,
					DurationUnit: atime.TIMEUNIT_MINUTELY,
					JoinWindows: JoinWindows{
						{IsBefore: true, Duration: 10, DurationUnit: atime.TIMEUNIT_MINUTELY},
						{IsBefore: false, Duration: 10, DurationUnit: atime.TIMEUNIT_MINUTELY},
					},
				},
			},
			cases: []timeCase{
				{"Both match, deny wins", atime.MustParseRFC3339("2025-06-20T10:05:00Z"), false, "Deny overrides allow"},
			},
		},
		{
			name:        "Allow wins with higher priority",
			description: "Allow has higher priority, overrides deny",
			rrules: RRuleExtends{
				{
					ROptions: ROptionExtend{
						Freq:      atime.TIMEUNIT_DAILY,
						BeginTime: atime.MustParsePtrRFC3339("2025-06-20T10:00:00Z"),
						ByHour:    []int{10},
					},
					IsDeny:       true,
					Priority:     1,
					Duration:     1,
					DurationUnit: atime.TIMEUNIT_MINUTELY,
					JoinWindows: JoinWindows{
						{IsBefore: true, Duration: 10, DurationUnit: atime.TIMEUNIT_MINUTELY},
						{IsBefore: false, Duration: 10, DurationUnit: atime.TIMEUNIT_MINUTELY},
					},
				},
				{
					ROptions: ROptionExtend{
						Freq:      atime.TIMEUNIT_DAILY,
						BeginTime: atime.MustParsePtrRFC3339("2025-06-20T10:00:00Z"),
						ByHour:    []int{10},
					},
					IsDeny:       false,
					Priority:     10,
					Duration:     1,
					DurationUnit: atime.TIMEUNIT_MINUTELY,
					JoinWindows: JoinWindows{
						{IsBefore: true, Duration: 10, DurationUnit: atime.TIMEUNIT_MINUTELY},
						{IsBefore: false, Duration: 10, DurationUnit: atime.TIMEUNIT_MINUTELY},
					},
				},
			},
			cases: []timeCase{
				{"Allow wins", atime.MustParseRFC3339("2025-06-20T10:00:30Z"), true, "Allow overrides deny"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, c := range tt.cases {
				t.Run(c.label, func(t *testing.T) {
					result, err := tt.rrules.EvaluateWithOptions(c.now, tt.geo, tt.eval)
					require.NoError(t, err)
					require.Equal(t, c.expect, result, c.expectMsg)
				})
			}
		})
	}
}

func TestRRuleExtend_geoCheck_GeoFilters(t *testing.T) {
	polygon := ageo.GISPoints{
		{Latitude: 0.0, Longitude: 0.0},
		{Latitude: 0.0, Longitude: 10.0},
		{Latitude: 10.0, Longitude: 10.0},
		{Latitude: 10.0, Longitude: 0.0},
		{Latitude: 0.0, Longitude: 0.0},
	}

	tests := []struct {
		name     string
		rule     RRuleExtend
		geo      ageo.GeoInfo
		expected bool
	}{
		{
			name: "Inside polygon",
			rule: RRuleExtend{
				GeoFilters: ageo.GeoFilters{
					&ageo.GeoFilter{GISPolygon: polygon},
				},
			},
			geo: ageo.GeoInfo{
				City: "TestCity",
				GISPoint: ageo.GISPoint{
					Latitude:  5.0,
					Longitude: 5.0,
				},
			},
			expected: true,
		},
		{
			name: "Near polygon edge (within radius)",
			rule: RRuleExtend{
				GeoFilters: ageo.GeoFilters{
					&ageo.GeoFilter{GISPolygon: polygon},
				},
			},
			geo: ageo.GeoInfo{
				GISPoint: ageo.GISPoint{
					Latitude:  0.0,
					Longitude: 10.05, // ~5.5km from east edge
				},
			},
			expected: true,
		},
		{
			name: "Far from polygon",
			rule: RRuleExtend{
				GeoFilters: ageo.GeoFilters{
					&ageo.GeoFilter{GISPolygon: polygon},
				},
			},
			geo: ageo.GeoInfo{
				GISPoint: ageo.GISPoint{
					Latitude:  20.0,
					Longitude: 20.0,
				},
			},
			expected: false,
		},
		{
			name: "Valid city but fails polygon",
			rule: RRuleExtend{
				GeoFilters: ageo.GeoFilters{
					&ageo.GeoFilter{
						Cities:     []string{"TestCity"},
						GISPolygon: polygon,
					},
				},
			},
			geo: ageo.GeoInfo{
				City: "TestCity",
				GISPoint: ageo.GISPoint{
					Latitude:  50.0,
					Longitude: 50.0,
				},
			},
			expected: false,
		},
		{
			name: "Empty polygon, city matches",
			rule: RRuleExtend{
				GeoFilters: ageo.GeoFilters{
					&ageo.GeoFilter{
						Cities: []string{"TestCity"},
					},
				},
			},
			geo: ageo.GeoInfo{
				City: "TestCity",
				GISPoint: ageo.GISPoint{
					Latitude:  45.0,
					Longitude: 99.0,
				},
			},
			expected: true,
		},
		{
			name: "Deny filter matches polygon → block",
			rule: RRuleExtend{
				GeoFilters: ageo.GeoFilters{
					&ageo.GeoFilter{
						IsDeny:     true,
						GISPolygon: polygon,
					},
				},
			},
			geo: ageo.GeoInfo{
				GISPoint: ageo.GISPoint{
					Latitude:  5.0,
					Longitude: 5.0,
				},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok := tt.rule.GeoFilters.Evaluate(tt.geo)
			if ok != tt.expected {
				t.Errorf("GeoFilters.Evaluate() = %v, want %v", ok, tt.expected)
			}
		})
	}
}

func TestRRuleExtend_MatchJoinWindow_DAILY(t *testing.T) {
	type timeCase struct {
		label     string
		now       time.Time
		expect    bool
		expectTag string
	}

	tests := []struct {
		name        string
		rrule       RRuleExtend
		description string
		cases       []timeCase
	}{
		{
			name: "Grace period after 9 AM login (10 mins)",
			rrule: RRuleExtend{
				ROptions: ROptionExtend{
					Freq:      atime.TIMEUNIT_DAILY,
					Interval:  1,
					BeginTime: atime.MustParsePtrRFC3339("2025-06-01T09:00:00Z"),
					ByHour:    []int{9},
				},
				Duration:     1,
				DurationUnit: atime.TIMEUNIT_MINUTELY,
				JoinWindows: JoinWindows{
					&JoinWindow{
						IsBefore:     false,
						Duration:     10,
						DurationUnit: atime.TIMEUNIT_MINUTELY,
						Tag:          "grace_after_login",
					},
				},
			},
			description: "Allow access up to 10 minutes after scheduled login",
			cases: []timeCase{
				{"8:59 AM too early", atime.MustParseRFC3339("2025-06-11T08:59:00Z"), false, ""},
				{"9:00 AM on time", atime.MustParseRFC3339("2025-06-11T09:00:00Z"), true, "grace_after_login"},
				{"9:09 AM still valid", atime.MustParseRFC3339("2025-06-11T09:09:59Z"), true, "grace_after_login"},
				{"9:11 AM too late", atime.MustParseRFC3339("2025-06-11T09:11:00Z"), false, ""},
			},
		},
		{
			name: "Multiple JoinWindows (pre and post)",
			rrule: RRuleExtend{
				ROptions: ROptionExtend{
					Freq:      atime.TIMEUNIT_DAILY,
					Interval:  1,
					BeginTime: atime.MustParsePtrRFC3339("2025-06-20T10:00:00Z"),
					ByHour:    []int{10},
				},
				Duration:     1,
				DurationUnit: atime.TIMEUNIT_MINUTELY,
				JoinWindows: JoinWindows{
					&JoinWindow{IsBefore: true, Duration: 10, DurationUnit: atime.TIMEUNIT_MINUTELY, Tag: "pre_10min"},
					&JoinWindow{IsBefore: false, Duration: 10, DurationUnit: atime.TIMEUNIT_MINUTELY, Tag: "post_10min"},
				},
			},
			description: "Check matching against both JoinWindows",
			cases: []timeCase{
				{"Too early", atime.MustParseRFC3339("2025-06-20T09:49:00Z"), false, ""},
				{"Within pre-window", atime.MustParseRFC3339("2025-06-20T09:55:00Z"), true, "pre_10min"},
				{"Exact start", atime.MustParseRFC3339("2025-06-20T10:00:00Z"), true, "post_10min"},
				{"Within post-window", atime.MustParseRFC3339("2025-06-20T10:05:00Z"), true, "post_10min"},
				{"Too late", atime.MustParseRFC3339("2025-06-20T10:11:00Z"), false, ""},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, tc := range tt.cases {
				t.Run(tc.label, func(t *testing.T) {
					jw, err := tt.rrule.MatchJoinWindow(tc.now, true, true)
					require.NoError(t, err)
					if tc.expect {
						require.NotNil(t, jw, "[%s] Expected to match JoinWindow", tc.label)
						require.Equal(t, tc.expectTag, jw.Tag, "[%s] Tag mismatch", tc.label)
					} else {
						require.Nil(t, jw, "[%s] Expected no match", tc.label)
					}
				})
			}
		})
	}
}

func TestRRuleExtend_MatchJoinWindow(t *testing.T) {
	beginTime := atime.MustParseRFC3339("2020-12-12T00:00:00Z")

	rule := &RRuleExtend{
		ROptions: ROptionExtend{
			Freq:       atime.TIMEUNIT_YEARLY,
			Interval:   1,
			ByMonth:    []int{12},
			ByMonthDay: []int{12},
			BeginTime:  atime.ToPointer(beginTime),
		},
		Duration:     1,
		DurationUnit: atime.TIMEUNIT_DAILY,
		JoinWindows: JoinWindows{
			{IsBefore: true, Duration: 30, DurationUnit: atime.TIMEUNIT_DAILY, Tag: "billing_notice_30"},
			{IsBefore: true, Duration: 15, DurationUnit: atime.TIMEUNIT_DAILY, Tag: "billing_notice_15"},
			{IsBefore: true, Duration: 1, DurationUnit: atime.TIMEUNIT_DAILY, Tag: "billing_notice_1"},
		},
	}

	tests := []struct {
		name     string
		now      string
		expected string // tag
	}{
		{"Match 30-day window", "2025-11-12T10:00:00Z", "billing_notice_30"},
		{"Match 15-day window", "2025-11-27T10:00:00Z", "billing_notice_15"},
		{"Match 1-day window", "2025-12-11T10:00:00Z", "billing_notice_1"},
		{"No match before", "2025-10-01T00:00:00Z", ""},
		{"No match after", "2025-12-13T00:00:00Z", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			now := atime.MustParseRFC3339(tt.now)
			jw, err := rule.MatchJoinWindow(now, true, true)
			require.NoError(t, err)
			if tt.expected == "" {
				require.Nil(t, jw)
			} else {
				require.NotNil(t, jw)
				require.Equal(t, tt.expected, jw.Tag)
			}
		})
	}

	// Replaced NewRRuleSingleDate with builder-based version
	rule = NewRRBuilderMonthDay(beginTime, atime.TIMEUNIT_YEARLY).
		AddJWBeforeDaily("billing_notice", 30, 15, 1).
		WithDuration(1, atime.TIMEUNIT_DAILY).
		Build()

	tests = []struct {
		name     string
		now      string
		expected string // tag
	}{
		{"Match 30-day window", "2025-11-12T10:00:00Z", "billing_notice"},
		{"Match 15-day window", "2025-11-27T10:00:00Z", "billing_notice"},
		{"Match 1-day window", "2025-12-11T10:00:00Z", "billing_notice"},
		{"No match before", "2025-10-01T00:00:00Z", ""},
		{"No match after", "2025-12-13T00:00:00Z", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			now := atime.MustParseRFC3339(tt.now)
			jw, err := rule.MatchJoinWindow(now, true, true)
			require.NoError(t, err)
			if tt.expected == "" {
				require.Nil(t, jw)
			} else {
				require.NotNil(t, jw)
				require.Equal(t, tt.expected, jw.Tag)
			}
		})
	}
}

func TestRRuleExtend_GetNextTimes(t *testing.T) {
	tests := []struct {
		name        string
		rrule       RRuleExtend
		now         string
		count       int
		expectTimes []string
	}{
		{
			name: "Daily recurrence, 3 times",
			rrule: RRuleExtend{
				ROptions: ROptionExtend{
					Freq:      atime.TIMEUNIT_DAILY,
					BeginTime: atime.MustParsePtrRFC3339("2025-06-21T09:00:00Z"),
					ByHour:    []int{9},
				},
				Duration:     1,
				DurationUnit: atime.TIMEUNIT_HOURLY,
			},
			now:   "2025-06-21T08:00:00Z",
			count: 3,
			expectTimes: []string{
				"2025-06-21T09:00:00Z",
				"2025-06-22T09:00:00Z",
				"2025-06-23T09:00:00Z",
			},
		},
		{
			name: "Hourly recurrence, starting now (inclusive)",
			rrule: RRuleExtend{
				ROptions: ROptionExtend{
					Freq:      atime.TIMEUNIT_HOURLY,
					BeginTime: atime.MustParsePtrRFC3339("2025-06-21T09:00:00Z"),
				},
				Duration:     15,
				DurationUnit: atime.TIMEUNIT_MINUTELY,
			},
			now:   "2025-06-21T09:00:00Z",
			count: 2,
			expectTimes: []string{
				"2025-06-21T09:00:00Z",
				"2025-06-21T10:00:00Z",
			},
		},
		{
			name: "Zero count returns empty",
			rrule: RRuleExtend{
				ROptions: ROptionExtend{
					Freq:      atime.TIMEUNIT_DAILY,
					BeginTime: atime.MustParsePtrRFC3339("2025-06-21T00:00:00Z"),
					ByHour:    []int{0},
				},
				Duration:     1,
				DurationUnit: atime.TIMEUNIT_HOURLY,
			},
			now:         "2025-06-21T00:00:00Z",
			count:       0,
			expectTimes: []string{},
		},
		{
			name: "No future occurrences (Count limit)",
			rrule: RRuleExtend{
				ROptions: ROptionExtend{
					Freq:      atime.TIMEUNIT_DAILY,
					BeginTime: atime.MustParsePtrRFC3339("2025-06-01T00:00:00Z"),
					ByHour:    []int{0},
					Count:     2,
				},
				Duration:     1,
				DurationUnit: atime.TIMEUNIT_HOURLY,
			},
			now:         "2025-06-10T00:00:00Z", // after all occurrences
			count:       3,
			expectTimes: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			now := atime.MustParseRFC3339(tt.now)
			got, err := tt.rrule.GetNextTimes(now, tt.count)
			require.NoError(t, err)
			require.Len(t, got, len(tt.expectTimes))

			for i, exp := range tt.expectTimes {
				expectedTime := atime.MustParseRFC3339(exp)
				require.True(t, got[i].Equal(expectedTime),
					"Expected %v, got %v", expectedTime, got[i])
			}
		})
	}
}

func TestRRuleExtend_GetNextOccurrences(t *testing.T) {
	now := atime.MustParseRFC3339("2025-06-21T00:00:00Z")

	tests := []struct {
		name        string
		rrule       RRuleExtend
		count       int
		expectLen   int
		expectNames []string
		expectDeny  bool
	}{
		{
			name: "Daily recurrence",
			rrule: RRuleExtend{
				ROptions: ROptionExtend{
					Freq:      atime.TIMEUNIT_DAILY,
					Interval:  1,
					BeginTime: atime.MustParsePtrRFC3339("2025-06-20T08:00:00Z"),
					ByHour:    []int{8},
				},
				Duration:     1,
				DurationUnit: atime.TIMEUNIT_HOURLY,
				Name:         "daily-rule",
				Priority:     10,
			},
			count:       3,
			expectLen:   3,
			expectNames: []string{"daily-rule", "daily-rule", "daily-rule"},
			expectDeny:  false,
		},
		{
			name: "IsAnyTime returns now once",
			rrule: RRuleExtend{
				IsAnyTime: true,
				Name:      "anytime-rule",
			},
			count:       5,
			expectLen:   1, // Should always return exactly one result
			expectNames: []string{"anytime-rule"},
		},
		{
			name: "IsDeny rule still returns occurrences",
			rrule: RRuleExtend{
				ROptions: ROptionExtend{
					Freq:      atime.TIMEUNIT_DAILY,
					Interval:  1,
					BeginTime: atime.MustParsePtrRFC3339("2025-06-20T10:00:00Z"),
					ByHour:    []int{10},
				},
				Duration:     1,
				DurationUnit: atime.TIMEUNIT_HOURLY,
				IsDeny:       true,
				Name:         "deny-rule",
			},
			count:       2,
			expectLen:   2,
			expectNames: []string{"deny-rule", "deny-rule"},
			expectDeny:  true,
		},
		{
			name: "Zero count yields no occurrences",
			rrule: RRuleExtend{
				IsAnyTime: true,
				Name:      "anytime-zero",
			},
			count:       0,
			expectLen:   0,
			expectNames: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			occs, err := tt.rrule.GetNextOccurrences(now, tt.count)
			require.NoError(t, err)
			require.Len(t, occs, tt.expectLen)

			for i, occ := range occs {
				require.NotZero(t, occ.Time, "Occurrence time should not be zero")
				require.Equal(t, tt.expectDeny, occ.IsDeny, "IsDeny mismatch")
				if len(tt.expectNames) > i {
					require.Equal(t, tt.expectNames[i], occ.Name, "Rule name mismatch")
				}
			}
		})
	}
}

func TestGetNextOccurrencesStacked(t *testing.T) {
	now := atime.MustParseRFC3339("2025-06-20T00:00:00Z")

	rules := RRuleExtends{
		&RRuleExtend{
			Name: "Early Morning",
			ROptions: ROptionExtend{
				Freq:      atime.TIMEUNIT_DAILY,
				Interval:  1,
				BeginTime: atime.MustParsePtrRFC3339("2025-06-01T06:00:00Z"),
				ByHour:    []int{6},
			},
			Duration:     30,
			DurationUnit: atime.TIMEUNIT_MINUTELY,
		},
		&RRuleExtend{
			Name: "Late Morning",
			ROptions: ROptionExtend{
				Freq:      atime.TIMEUNIT_DAILY,
				Interval:  1,
				BeginTime: atime.MustParsePtrRFC3339("2025-06-01T10:00:00Z"),
				ByHour:    []int{10},
			},
			Duration:     30,
			DurationUnit: atime.TIMEUNIT_MINUTELY,
		},
	}

	result, err := rules.GetNextOccurrencesStacked(now, 2)
	require.NoError(t, err)
	require.Len(t, result, 2)

	// Check the first rule's occurrences
	early := result[0]
	require.Len(t, early, 2)
	require.Equal(t, "Early Morning", early[0].Name)
	require.True(t, early[0].Time.After(now))

	// Check the second rule's occurrences
	late := result[1]
	require.Len(t, late, 2)
	require.Equal(t, "Late Morning", late[0].Name)
	require.True(t, late[0].Time.After(now))

	// Ensure sorted order and distinct rule outputs
	require.NotEqual(t, early[0].Time, late[0].Time)
}

// Utility: creates a UTC time
func mustTime(y int, m time.Month, d, h, min int) *time.Time {
	t := time.Date(y, m, d, h, min, 0, 0, time.UTC)
	return &t
}

func TestRRuleExtend_String_NameOnly(t *testing.T) {
	r := RRuleExtend{
		Name: "Quarterly Check",
	}
	require.Equal(t, "Quarterly Check", r.String())
}

func TestRRuleExtend_String_Descriptor_Generation(t *testing.T) {
	r := RRuleExtend{
		IsDeny:   true,
		Priority: 5,
		ROptions: ROptionExtend{
			Freq:            atime.TIMEUNIT_MONTHLY,
			Interval:        2,
			BeginTime:       mustTime(2025, 1, 1, 9, 0),
			ByMonthDay:      []int{15},
			ByHour:          []int{9},
			ByMinute:        []int{30},
			ShiftOffWeekend: true,
			ISOCode:         "US",
			Observance:      ObservanceNextBizDay,
		},
	}
	str := r.String()
	require.Contains(t, str, "Schedule Rule")
	require.Contains(t, str, "DENY rule")
	require.Contains(t, str, "Every 2 monthly")
	require.Contains(t, str, "Start: 2025-01-01 09:00")
	require.Contains(t, str, "On month days: [15]")
	require.Contains(t, str, "At 09:30")
	require.Contains(t, str, "Shift off weekends")
	require.Contains(t, str, "Region: US")
	require.Contains(t, str, "Observance: next-business-day")
}

func TestRRuleExtend_String_IsAnyTime(t *testing.T) {
	r := RRuleExtend{
		IsAnyTime: true,
	}
	require.Contains(t, r.String(), "Any time")
}

func TestRRuleExtend_String_WithDurationAndJoin(t *testing.T) {
	r := RRuleExtend{
		ROptions: ROptionExtend{
			Freq:     atime.TIMEUNIT_WEEKLY,
			Interval: 1,
		},
		Duration:     2,
		DurationUnit: atime.TIMEUNIT_DAILY,
		JoinWindows: JoinWindows{
			{Duration: 2, DurationUnit: atime.TIMEUNIT_HOURLY},
		},
	}
	desc := r.ToDescriptor()
	require.Contains(t, desc, "Every weekly")
	require.Contains(t, desc, "Duration: 2 daily")
	require.Contains(t, desc, "Join window(s) enabled")
}

func TestRRuleExtend_String_WithGeoFilters(t *testing.T) {
	r := RRuleExtend{
		ROptions: ROptionExtend{
			Freq:     atime.TIMEUNIT_DAILY,
			Interval: 1,
		},
		GeoFilters: ageo.GeoFilters{
			{Countries: []string{"US"}, Regions: []string{"NY"}},
		},
	}
	desc := r.ToDescriptor()
	require.Contains(t, desc, "Geo filter(s) enabled")
}
