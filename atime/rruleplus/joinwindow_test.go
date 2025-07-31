package rruleplus

import (
	"github.com/jpfluger/alibs-slim/atime"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestJoinWindow_Matches(t *testing.T) {
	occurrence := atime.MustParseRFC3339("2025-12-12T10:00:00Z")

	tests := []struct {
		name     string
		now      string
		window   *JoinWindow
		expected bool
	}{
		{
			name: "Within 30-day notice window",
			now:  "2025-11-20T10:00:00Z",
			window: &JoinWindow{
				IsBefore:     true,
				Duration:     30,
				DurationUnit: atime.TIMEUNIT_DAILY,
			},
			expected: true,
		},
		{
			name: "Outside 15-day window",
			now:  "2025-11-20T10:00:00Z",
			window: &JoinWindow{
				IsBefore:     true,
				Duration:     15,
				DurationUnit: atime.TIMEUNIT_DAILY,
			},
			expected: false,
		},
		{
			name: "Exactly at window start — exclude",
			now:  "2025-11-12T10:00:00Z", // exactly 30 days before
			window: &JoinWindow{
				IsBefore:     true,
				Duration:     30,
				DurationUnit: atime.TIMEUNIT_DAILY,
			},
			expected: true,
		},
		{
			name: "Inside 1-day window",
			now:  "2025-12-11T12:00:00Z",
			window: &JoinWindow{
				IsBefore:     true,
				Duration:     1,
				DurationUnit: atime.TIMEUNIT_DAILY,
			},
			expected: true,
		},
		{
			name: "After the occurrence — no match",
			now:  "2025-12-13T00:00:00Z",
			window: &JoinWindow{
				IsBefore:     true,
				Duration:     30,
				DurationUnit: atime.TIMEUNIT_DAILY,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			now := atime.MustParseRFC3339(tt.now)
			result := tt.window.Matches(now, occurrence)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestJoinWindows_Matches(t *testing.T) {
	occurrence := atime.MustParseRFC3339("2025-12-12T10:00:00Z")
	windows := JoinWindows{
		{IsBefore: true, Duration: 30, DurationUnit: atime.TIMEUNIT_DAILY, Tag: "billing_notice_30"},
		{IsBefore: true, Duration: 15, DurationUnit: atime.TIMEUNIT_DAILY, Tag: "billing_notice_15"},
		{IsBefore: true, Duration: 1, DurationUnit: atime.TIMEUNIT_DAILY, Tag: "billing_notice_1"},
	}

	tests := []struct {
		name     string
		now      string
		expected string // expected tag
	}{
		{
			name:     "Within 30-day but outside 15-day window",
			now:      "2025-11-20T10:00:00Z",
			expected: "billing_notice_30",
		},
		{
			name:     "Within 15-day window",
			now:      "2025-11-30T10:00:00Z",
			expected: "billing_notice_15",
		},
		{
			name:     "Within 1-day window",
			now:      "2025-12-11T11:00:00Z",
			expected: "billing_notice_1",
		},
		{
			name:     "No match before window starts",
			now:      "2025-11-01T10:00:00Z",
			expected: "",
		},
		{
			name:     "No match after occurrence",
			now:      "2025-12-13T00:00:00Z",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			now := atime.MustParseRFC3339(tt.now)
			matched := windows.Matches(now, occurrence)
			if tt.expected == "" {
				require.Nil(t, matched)
			} else {
				require.NotNil(t, matched)
				require.Equal(t, tt.expected, matched.Tag)
			}
		})
	}
}
