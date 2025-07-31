package atime

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/teambition/rrule-go"
	"testing"
	"time"
)

// TestHoursInYear tests the calculation of hours in a regular year.
func TestHoursInYear(t *testing.T) {
	expected := 8760
	if got := HoursInYear(false); got != expected {
		t.Errorf("HoursInYear(false) = %d, want %d", got, expected)
	}
}

// TestHoursInLeapYear tests the calculation of hours in a leap year.
func TestHoursInLeapYear(t *testing.T) {
	expected := 8784
	if got := HoursInYear(true); got != expected {
		t.Errorf("HoursInYear(true) = %d, want %d", got, expected)
	}
}

// HoursInYear calculates the total hours in a year, accounting for leap years.
func HoursInYear(isLeapYear bool) int {
	if isLeapYear {
		return 366 * 24 // 8784 hours in a leap year
	}
	return 365 * 24 // 8760 hours in a regular year
}

type testStruct struct {
	Time1 *time.Time `json:"time1,omitempty"`
	Time2 time.Time  `json:"time2,omitempty"`
}

func marshalThenUnmarshal(target *testStruct) *testStruct {
	b, _ := json.Marshal(target)
	ts := &testStruct{}
	_ = json.Unmarshal(b, ts)
	return ts
}

func TestEnsureDateTime(t *testing.T) {
	a := time.Now()
	b := EnsureDateTime(a)
	assert.Equal(t, a.String(), b.String())

	b = EnsureDateTime(&a)
	assert.Equal(t, a.String(), b.String())

	c := GetNowPointer()
	b = EnsureDateTime(c)
	assert.Equal(t, c.String(), b.String())

	b = EnsureDateTime(*c)
	assert.Equal(t, c.String(), b.String())
}

func TestDateCompare(t *testing.T) {
	target := EnsureDateTime(time.Now().AddDate(0, 0, 1))
	assert.Equal(t, true, IsDateAfterNow(target))
	assert.Equal(t, true, IsDateAfterNow(&target))

	target = EnsureDateTime(time.Now().AddDate(0, 0, -1))
	assert.Equal(t, true, IsDateBeforeNow(target))
	assert.Equal(t, true, IsDateBeforeNow(&target))

	a := time.Now()
	b := a.AddDate(0, 0, -1)
	assert.Equal(t, true, IsDateAfter(a, b))
	b = a.AddDate(0, 0, 1)
	assert.Equal(t, true, IsDateBefore(a, b))
}

func TestDateFormat(t *testing.T) {
	a := time.Now()
	assert.Equal(t, a.Format(time.RFC1123), FormatDateTime(a, time.RFC1123))
	assert.Equal(t, a.Format(time.RFC1123), FormatDateTime(&a, time.RFC1123))

	assert.Equal(t, "alternate", FormatDateTimeElse(time.Time{}, time.RFC1123, "alternate"))
	assert.Equal(t, "alternate", FormatDateTimeElse(&time.Time{}, time.RFC1123, "alternate"))
}

func TestTime(t *testing.T) {

	assert.NotNil(t, ToPointer(time.Now()))
	assert.NotNil(t, ToPointerNil(time.Now()))
	assert.Nil(t, ToPointerNil(time.Time{}))

	assert.NotNil(t, GetNowPointer())
	assert.NotNil(t, GetNowUTCPointer())

	ts := &testStruct{
		Time1: nil,
		Time2: time.Time{},
	}

	ts = marshalThenUnmarshal(ts)
	assert.Nil(t, ts.Time1)
	assert.NotNil(t, ts.Time2)
	assert.Equal(t, true, ts.Time2.IsZero())

	ts = &testStruct{
		Time1: GetNowPointer(),
		Time2: time.Now(),
	}

	ts = marshalThenUnmarshal(ts)
	assert.NotNil(t, ts.Time1)
	assert.Equal(t, false, ts.Time1.IsZero())
	assert.NotNil(t, ts.Time2)
	assert.Equal(t, false, ts.Time2.IsZero())
}

func TestIfDateEmptyElse(t *testing.T) {
	assert.Equal(t, "empty", IfDateEmptyElse(nil, "empty", "value"))
	assert.Equal(t, "empty", IfDateEmptyElse(time.Time{}, "empty", "value"))
	assert.Equal(t, "value", IfDateEmptyElse(ToPointer(time.Now()), "empty", "value"))
	assert.Equal(t, "value", IfDateEmptyElse(time.Now(), "empty", "value"))

	myt := time.Date(2000, 2, 0, 0, 0, 0, 0, time.UTC)
	assert.Equal(t, 2000, myt.Year())
}

func TestSanitizePtr(t *testing.T) {
	now := time.Now()
	zero := time.Time{}

	tests := []struct {
		name   string
		input  *time.Time
		expect *time.Time
	}{
		{
			name:   "Nil input returns nil",
			input:  nil,
			expect: nil,
		},
		{
			name:   "Zero time input returns nil",
			input:  &zero,
			expect: nil,
		},
		{
			name:   "Valid time input returns same pointer",
			input:  &now,
			expect: &now,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizePtr(tt.input)
			if tt.expect == nil {
				assert.Nil(t, result)
			} else {
				assert.Equal(t, tt.expect, result)
			}
		})
	}
}

func TestEnsureDateTimeUTC(t *testing.T) {
	t.Run("nil pointer returns zero time", func(t *testing.T) {
		var tptr *time.Time
		got := EnsureDateTimeUTC(tptr)
		require.True(t, got.IsZero(), "expected zero time, got: %v", got)
	})

	t.Run("non-nil pointer returns UTC time", func(t *testing.T) {
		local := time.Date(2025, 6, 19, 15, 30, 0, 0, time.FixedZone("Custom", -7*3600))
		tptr := &local
		got := EnsureDateTimeUTC(tptr)
		require.Equal(t, local.UTC(), got)
	})

	t.Run("non-pointer time returns UTC time", func(t *testing.T) {
		local := time.Date(2025, 6, 19, 12, 0, 0, 0, time.FixedZone("EDT", -4*3600))
		got := EnsureDateTimeUTC(local)
		require.Equal(t, local.UTC(), got)
	})

	t.Run("invalid type returns zero time", func(t *testing.T) {
		got := EnsureDateTimeUTC("not a time")
		require.True(t, got.IsZero(), "expected zero time for invalid input, got: %v", got)
	})
}

func TestMustParse(t *testing.T) {
	t.Run("valid RFC3339", func(t *testing.T) {
		input := "2025-06-21T14:00:00Z"
		expected, _ := time.Parse(time.RFC3339, input)

		result := MustParse(time.RFC3339, input)
		require.Equal(t, expected, result)
	})

	t.Run("invalid layout", func(t *testing.T) {
		result := MustParse("bad-layout", "2025-06-21T14:00:00Z")
		require.True(t, result.IsZero())
	})

	t.Run("invalid time string", func(t *testing.T) {
		result := MustParse(time.RFC3339, "not-a-date")
		require.True(t, result.IsZero())
	})
}

func TestMustParsePtr(t *testing.T) {
	t.Run("valid time string", func(t *testing.T) {
		input := "2025-06-21T14:00:00Z"
		expected, _ := time.Parse(time.RFC3339, input)

		ptr := MustParsePtr(time.RFC3339, input)
		require.NotNil(t, ptr)
		require.Equal(t, expected, *ptr)
	})

	t.Run("invalid time string", func(t *testing.T) {
		ptr := MustParsePtr(time.RFC3339, "bad-value")
		require.NotNil(t, ptr)
		require.True(t, ptr.IsZero())
	})
}

func TestMustParseRFC3339(t *testing.T) {
	t.Run("valid RFC3339", func(t *testing.T) {
		input := "2025-06-21T14:00:00Z"
		expected, _ := time.Parse(time.RFC3339, input)

		result := MustParseRFC3339(input)
		require.Equal(t, expected, result)
	})

	t.Run("invalid RFC3339", func(t *testing.T) {
		result := MustParseRFC3339("bad-rfc-date")
		require.True(t, result.IsZero())
	})
}

func TestMustParsePtrRFC3339(t *testing.T) {
	t.Run("valid RFC3339", func(t *testing.T) {
		input := "2025-06-21T14:00:00Z"
		expected, _ := time.Parse(time.RFC3339, input)

		ptr := MustParsePtrRFC3339(input)
		require.NotNil(t, ptr)
		require.Equal(t, expected, *ptr)
	})

	t.Run("invalid RFC3339", func(t *testing.T) {
		ptr := MustParsePtrRFC3339("bad-time")
		require.NotNil(t, ptr)
		require.True(t, ptr.IsZero())
	})
}

func TestHoursBetween(t *testing.T) {
	tests := []struct {
		name         string
		start        int
		end          int
		expectedWrap []int
		expectedSort []int
	}{
		{
			name:         "Simple increasing",
			start:        8,
			end:          11,
			expectedWrap: []int{8, 9, 10},
			expectedSort: []int{8, 9, 10},
		},
		{
			name:         "Same start and end at 0",
			start:        0,
			end:          0,
			expectedWrap: []int{},
			expectedSort: []int{},
		},
		{
			name:         "Wrap around midnight 22→2",
			start:        22,
			end:          2,
			expectedWrap: []int{22, 23, 0, 1},
			expectedSort: []int{0, 1, 22, 23},
		},
		{
			name:         "Wrap around midnight 23→1",
			start:        23,
			end:          1,
			expectedWrap: []int{23, 0},
			expectedSort: []int{0, 23},
		},
		{
			name:         "Same hour",
			start:        5,
			end:          5,
			expectedWrap: []int{},
			expectedSort: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name+"/WrapOrder", func(t *testing.T) {
			result := HoursBetween(tt.start, tt.end, false)
			require.Equal(t, tt.expectedWrap, result)
		})

		t.Run(tt.name+"/SortedOrder", func(t *testing.T) {
			result := HoursBetween(tt.start, tt.end, true)
			require.Equal(t, tt.expectedSort, result)
		})
	}
}

func TestHoursBetweenDates(t *testing.T) {
	tests := []struct {
		name     string
		start    string
		end      string
		expected []int
	}{
		{"Same day", "2025-06-21T08:00:00Z", "2025-06-21T11:00:00Z", []int{8, 9, 10}},
		{"Wrap midnight", "2025-06-21T22:00:00Z", "2025-06-22T02:00:00Z", []int{22, 23, 0, 1}}, // fixed ordering
		{"Same hour", "2025-06-21T05:00:00Z", "2025-06-21T05:00:00Z", []int{}},
		{"Multi-day", "2025-06-20T23:00:00Z", "2025-06-21T02:00:00Z", []int{23, 0, 1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start := MustParseRFC3339(tt.start)
			end := MustParseRFC3339(tt.end)
			result := HoursBetweenDates(start, end, false)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestHoursBetweenDates_Ordered(t *testing.T) {
	tests := []struct {
		name     string
		start    string
		end      string
		expected []int
	}{
		{
			name:     "Same day (ordered)",
			start:    "2025-06-21T08:00:00Z",
			end:      "2025-06-21T11:00:00Z",
			expected: []int{8, 9, 10},
		},
		{
			name:     "Wrap midnight (ordered)",
			start:    "2025-06-21T22:00:00Z",
			end:      "2025-06-22T02:00:00Z",
			expected: []int{0, 1, 22, 23},
		},
		{
			name:     "Same hour (empty)",
			start:    "2025-06-21T05:00:00Z",
			end:      "2025-06-21T05:00:00Z",
			expected: []int{},
		},
		{
			name:     "Multi-day (ordered)",
			start:    "2025-06-20T23:00:00Z",
			end:      "2025-06-21T02:00:00Z",
			expected: []int{0, 1, 23},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start := MustParseRFC3339(tt.start)
			end := MustParseRFC3339(tt.end)
			result := HoursBetweenDates(start, end, true)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestTimeWeekdayToRRuleWeekday(t *testing.T) {
	expected := []rrule.Weekday{
		rrule.SU, rrule.MO, rrule.TU, rrule.WE, rrule.TH, rrule.FR, rrule.SA,
	}
	for i := time.Sunday; i <= time.Saturday; i++ {
		require.Equal(t, expected[i], TimeWeekdayToRRuleWeekday(i))
	}
}

func TestRRuleWeekdayToTimeWeekday(t *testing.T) {
	cases := []struct {
		rr   rrule.Weekday
		want time.Weekday
	}{
		{rrule.SU, time.Sunday},
		{rrule.MO, time.Monday},
		{rrule.TU, time.Tuesday},
		{rrule.WE, time.Wednesday},
		{rrule.TH, time.Thursday},
		{rrule.FR, time.Friday},
		{rrule.SA, time.Saturday},
	}
	for _, c := range cases {
		require.Equal(t, c.want, RRuleWeekdayToTimeWeekday(c.rr))
	}
}

func TestTimeWeekdaysToInts(t *testing.T) {
	days := []time.Weekday{time.Sunday, time.Monday, time.Saturday}
	expected := []int{0, 1, 6}
	require.Equal(t, expected, TimeWeekdaysToInts(days...))
}

func TestTimeWeekdaysToRRuleWeekdays(t *testing.T) {
	days := []time.Weekday{time.Sunday, time.Monday, time.Saturday}
	ruleDays := TimeWeekdaysToRRuleWeekdays(days...)
	require.Equal(t, []rrule.Weekday{rrule.SU, rrule.MO, rrule.SA}, ruleDays)
}

func TestRRuleWeekdaysToInts(t *testing.T) {
	days := []rrule.Weekday{rrule.SU, rrule.MO, rrule.SA}
	expected := []int{6, 0, 5}
	require.Equal(t, expected, RRuleWeekdaysToInts(days...))
}

func TestRRuleWeekdaysToTimeWeekdays(t *testing.T) {
	ruleDays := []rrule.Weekday{rrule.SU, rrule.MO, rrule.SA}
	expected := []time.Weekday{time.Sunday, time.Monday, time.Saturday}
	require.Equal(t, expected, RRuleWeekdaysToTimeWeekdays(ruleDays...))
}
