package atime

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
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
