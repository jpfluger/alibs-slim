package rruleplus

import (
	"github.com/jpfluger/alibs-slim/atime"
	"github.com/teambition/rrule-go"
	"time"
)

type IRRule interface {
	After(time.Time, bool) time.Time
	Before(time.Time, bool) time.Time
	Between(time.Time, time.Time, bool) []time.Time
}

type RRulePlus struct {
	base     rrule.RRule
	calendar ICalendar // optional
	options  ROptionPlus
}

func NewRRulePlus(opt ROptionPlus) (*RRulePlus, error) {
	base, err := rrule.NewRRule(opt.ROption)
	if err != nil {
		return nil, err
	}

	var cal ICalendar

	if opt.Calendar != nil {
		cal = opt.Calendar
	} else if opt.ISOCode != "" {
		// Resolved from ISOCode inside options
		c, err := GetCalendar(opt.ISOCode)
		if err != nil || c == nil {
			if c, err = NewCalendar(opt.ISOCode); err != nil {
				return nil, err
			}
			SetCalendar(opt.ISOCode, c)
		}
		cal = c
	}

	return &RRulePlus{
		base:     *base,
		calendar: cal,
		options:  opt,
	}, nil
}

func (rp *RRulePlus) IsPlusMode() bool {
	opt := rp.options
	return opt.ShiftOffWeekend ||
		opt.ShiftOffHolidays ||
		opt.ValidOnlyOnHolidays ||
		opt.ValidOnlyOnWeekends ||
		!opt.Observance.IsEmpty() ||
		opt.CustomFilter != nil ||
		opt.ISOCode != ""
}

func (rp *RRulePlus) isValid(t time.Time) bool {
	opt := rp.options

	// Custom filter
	if opt.CustomFilter != nil && !opt.CustomFilter(t) {
		return false
	}

	// Weekend check
	isWeekend := atime.IsWeekendByTime(t)

	// Holiday check
	var isHoliday bool
	if rp.calendar != nil {
		actual, observed, _ := rp.calendar.IsHoliday(t)
		isHoliday = actual || observed
	}

	// Constraints: Only on holidays/weekends
	if opt.ValidOnlyOnWeekends && !isWeekend {
		return false
	}
	if opt.ValidOnlyOnHolidays && !isHoliday {
		return false
	}

	// Filtering logic: reject unwanted times unless shifted
	if isWeekend && !opt.ShiftOffWeekend && !opt.ValidOnlyOnWeekends {
		return false
	}
	if isHoliday && !opt.ShiftOffHolidays && !opt.ValidOnlyOnHolidays {
		return false
	}

	return true
}

// applyShift adjusts time based on weekend/holiday observance.
func (rp *RRulePlus) applyShift(t time.Time) time.Time {
	opt := rp.options

	// Shift off weekends
	if opt.ShiftOffWeekend {
		switch t.Weekday() {
		case time.Saturday:
			t = t.AddDate(0, 0, 2)
		case time.Sunday:
			t = t.AddDate(0, 0, 1)
		}
	}

	// Shift off holidays
	if opt.ShiftOffHolidays && rp.calendar != nil {
		for {
			actual, observed, _ := rp.calendar.IsHoliday(t)
			if !actual && !observed {
				break
			}
			t = t.AddDate(0, 0, 1)
		}
	}

	// Observance fallback
	if !opt.Observance.IsEmpty() && rp.calendar != nil {
		switch opt.Observance {
		case ObservanceNextBizDay:
			for {
				actual, observed, _ := rp.calendar.IsHoliday(t)
				if !actual && !observed && !atime.IsWeekendByTime(t) {
					break
				}
				t = t.AddDate(0, 0, 1)
			}
		case ObservancePreviousBizDay:
			for {
				actual, observed, _ := rp.calendar.IsHoliday(t)
				if !actual && !observed && !atime.IsWeekendByTime(t) {
					break
				}
				t = t.AddDate(0, 0, -1)
			}
		}
	}

	return t
}

func (rp *RRulePlus) scan(forward bool, t time.Time, inclusive bool) time.Time {
	cursor := t
	step := time.Second
	if !forward {
		step = -step
	}

	for attempts := 0; attempts < 1000; attempts++ {
		var next time.Time
		if forward {
			next = rp.base.After(cursor, inclusive)
		} else {
			next = rp.base.Before(cursor, inclusive)
		}
		if next.IsZero() {
			return time.Time{}
		}
		adjusted := rp.applyShift(next)
		if rp.isValid(adjusted) {
			return adjusted
		}
		cursor = next.Add(step)
		inclusive = false
	}
	return time.Time{}
}

func (rp *RRulePlus) After(t time.Time, inclusive bool) time.Time {
	if !rp.IsPlusMode() {
		return rp.base.After(t, inclusive)
	}
	return rp.scan(true, t, inclusive)
}

func (rp *RRulePlus) Before(t time.Time, inclusive bool) time.Time {
	if !rp.IsPlusMode() {
		return rp.base.Before(t, inclusive)
	}
	return rp.scan(false, t, inclusive)
}

// Between returns all valid occurrences in a range, applying filtering and shifts.
func (rp *RRulePlus) Between(after, before time.Time, inclusive bool) []time.Time {
	if !rp.IsPlusMode() {
		return rp.base.Between(after, before, inclusive)
	}

	results := []time.Time{}
	raw := rp.base.Between(after, before, inclusive)

	for _, t := range raw {
		adjusted := rp.applyShift(t)
		if rp.isValid(adjusted) && adjusted.After(after) && adjusted.Before(before) {
			results = append(results, adjusted)
		}
	}

	return results
}
