package atime

import (
	"github.com/teambition/rrule-go"
	"testing"
)

func TestTimeUnit_IsEmpty(t *testing.T) {
	tests := []struct {
		name string
		unit TimeUnit
		want bool
	}{
		{"Empty string", "", true},
		{"Non-empty", TIMEUNIT_HOURLY, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.unit.IsEmpty(); got != tt.want {
				t.Errorf("IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimeUnit_IsValid(t *testing.T) {
	validUnits := []TimeUnit{
		TIMEUNIT_SECONDLY,
		TIMEUNIT_MINUTELY,
		TIMEUNIT_HOURLY,
		TIMEUNIT_DAILY,
		TIMEUNIT_WEEKLY,
		TIMEUNIT_MONTHLY,
		TIMEUNIT_YEARLY,
	}
	for _, u := range validUnits {
		t.Run(string(u), func(t *testing.T) {
			if !u.IsValid() {
				t.Errorf("TimeUnit %q should be valid", u)
			}
		})
	}

	t.Run("Invalid unit", func(t *testing.T) {
		invalid := TimeUnit("decadely")
		if invalid.IsValid() {
			t.Errorf("TimeUnit %q should be invalid", invalid)
		}
	})
}

func TestTimeUnit_Default(t *testing.T) {
	tests := []struct {
		input TimeUnit
		want  TimeUnit
	}{
		{"", TIMEUNIT_DAILY},
		{TIMEUNIT_HOURLY, TIMEUNIT_HOURLY},
	}
	for _, tt := range tests {
		t.Run(tt.input.String(), func(t *testing.T) {
			if got := tt.input.Default(); got != tt.want {
				t.Errorf("Default() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimeUnit_ToFrequency(t *testing.T) {
	tests := []struct {
		unit TimeUnit
		want rrule.Frequency
	}{
		{TIMEUNIT_SECONDLY, rrule.SECONDLY},
		{TIMEUNIT_MINUTELY, rrule.MINUTELY},
		{TIMEUNIT_HOURLY, rrule.HOURLY},
		{TIMEUNIT_DAILY, rrule.DAILY},
		{TIMEUNIT_WEEKLY, rrule.WEEKLY},
		{TIMEUNIT_MONTHLY, rrule.MONTHLY},
		{TIMEUNIT_YEARLY, rrule.YEARLY},
		{"", rrule.DAILY}, // Default fallback
	}
	for _, tt := range tests {
		t.Run(tt.unit.String(), func(t *testing.T) {
			if got := tt.unit.ToFrequency(); got != tt.want {
				t.Errorf("ToFrequency() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFromFrequency(t *testing.T) {
	tests := []struct {
		freq rrule.Frequency
		want TimeUnit
	}{
		{rrule.SECONDLY, TIMEUNIT_SECONDLY},
		{rrule.MINUTELY, TIMEUNIT_MINUTELY},
		{rrule.HOURLY, TIMEUNIT_HOURLY},
		{rrule.DAILY, TIMEUNIT_DAILY},
		{rrule.WEEKLY, TIMEUNIT_WEEKLY},
		{rrule.MONTHLY, TIMEUNIT_MONTHLY},
		{rrule.YEARLY, TIMEUNIT_YEARLY},
		{999, TIMEUNIT_DAILY}, // Unknown fallback
	}
	for _, tt := range tests {
		t.Run(tt.want.String(), func(t *testing.T) {
			if got := FromFrequency(int(tt.freq)); got != tt.want {
				t.Errorf("FromFrequency(%d) = %v, want %v", tt.freq, got, tt.want)
			}
		})
	}
}
