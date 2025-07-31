package rruleplus

import (
	"fmt"
	"github.com/jpfluger/alibs-slim/atime"
	"github.com/teambition/rrule-go"
	"strings"
	"time"
)

type ROptionPlus struct {
	rrule.ROption `json:"-"` // Embedded options not re-tagged directly here

	// Policy extensions
	ShiftOffWeekend     bool                 `json:"shiftOffWeekend"`     // Shift Sat/Sun to Monday
	ShiftOffHolidays    bool                 `json:"shiftOffHolidays"`    // Shift holidays to next weekday
	ValidOnlyOnHolidays bool                 `json:"validOnlyOnHolidays"` // Only match on holiday dates
	ValidOnlyOnWeekends bool                 `json:"validOnlyOnWeekends"` // Only match on Saturday/Sunday
	ISOCode             string               `json:"isoCode"`             // e.g., "US", "US-NY" for region-aware calendars
	Observance          ObservanceMode       `json:"observance"`          // e.g., ObservanceNextBizDay
	CustomFilter        func(time.Time) bool `json:"-"`                   // Functions should not be serialized
	Calendar            ICalendar            `json:"-"`                   // Interface types typically excluded from JSON
}

type ROptionExtend struct {
	// RRule: Core recurrence rule config
	Freq      atime.TimeUnit `json:"freq"`                // e.g. "daily", "weekly"
	Interval  int            `json:"interval"`            // Interval between occurrences
	Count     int            `json:"count,omitempty"`     // Max number of recurrences
	BeginTime *time.Time     `json:"beginTime,omitempty"` // Start time of recurrence
	UntilTime *time.Time     `json:"untilTime,omitempty"` // End time

	// RRule: BYxxx rule filters (RFC 5545) – narrow the recurrence set
	BySecond   []int           `json:"bySecond,omitempty"`   // 0–59
	ByMinute   []int           `json:"byMinute,omitempty"`   // 0–59
	ByHour     []int           `json:"byHour,omitempty"`     // 0–23
	ByDay      []rrule.Weekday `json:"byDay,omitempty"`      // e.g. MO, -1FR
	ByMonthDay []int           `json:"byMonthDay,omitempty"` // 1–31 or -31–-1
	ByYearDay  []int           `json:"byYearDay,omitempty"`  // 1–366 or -366–-1
	ByWeekNo   []int           `json:"byWeekNo,omitempty"`   // 1–53 or -53–-1
	ByMonth    []int           `json:"byMonth,omitempty"`    // 1–12
	BySetPos   []int           `json:"bySetPos,omitempty"`   // select N-th from other BY rules
	ByEaster   []int           `json:"byEaster,omitempty"`   // days before/after Easter (rare)

	// Policy extensions
	ShiftOffWeekend     bool                 `json:"shiftOffWeekend"`     // Shift Sat/Sun to Monday
	ShiftOffHolidays    bool                 `json:"shiftOffHolidays"`    // Shift holidays to next weekday
	ValidOnlyOnHolidays bool                 `json:"validOnlyOnHolidays"` // Only match on holiday dates
	ValidOnlyOnWeekends bool                 `json:"validOnlyOnWeekends"` // Only match on Saturday/Sunday
	ISOCode             string               `json:"isoCode"`             // e.g., "US", "US-NY" for region-aware calendars
	Observance          ObservanceMode       `json:"observance"`          // e.g., ObservanceNextBizDay
	CustomFilter        func(time.Time) bool `json:"-"`                   // Functions should not be serialized
	Calendar            ICalendar            `json:"-"`                   // Interface types typically excluded from JSON

	// Optional in RRuleExtend
	RRIncType RRIncType `json:"rrIncType,omitempty"`
}

func (opt *ROptionExtend) WithCustomFilter(fn func(time.Time) bool) *ROptionExtend {
	opt.CustomFilter = fn
	return opt
}

func (opt *ROptionExtend) WithCalendar(cal ICalendar) *ROptionExtend {
	opt.Calendar = cal
	return opt
}

func (opt *ROptionExtend) ValidateRecurrence() error {
	if !opt.Freq.IsValid() {
		return fmt.Errorf("invalid frequency: %q", opt.Freq)
	}
	if opt.Interval < 1 {
		return fmt.Errorf("interval must be ≥ 1")
	}

	// Ensure BeginTime is before or equal to UntilTime if both are set
	if opt.BeginTime != nil && opt.UntilTime != nil {
		if opt.BeginTime.After(*opt.UntilTime) {
			return fmt.Errorf("beginTime (%v) must not be after untilTime (%v)", *opt.BeginTime, *opt.UntilTime)
		}
	}

	// Internal helper to check bounds on BYxxx values
	check := func(name string, vals []int, min, max int, allowNeg bool) error {
		for _, v := range vals {
			if v >= min && v <= max {
				continue
			}
			if allowNeg && v <= -min && v >= -max {
				continue
			}
			return fmt.Errorf("%s: value %d out of bounds (%d to %d)", name, v, min, max)
		}
		return nil
	}

	if err := check("bySecond", opt.BySecond, 0, 59, false); err != nil {
		return err
	}
	if err := check("byMinute", opt.ByMinute, 0, 59, false); err != nil {
		return err
	}
	if err := check("byHour", opt.ByHour, 0, 23, false); err != nil {
		return err
	}
	if err := check("byMonthDay", opt.ByMonthDay, 1, 31, true); err != nil {
		return err
	}
	if err := check("byYearDay", opt.ByYearDay, 1, 366, true); err != nil {
		return err
	}
	if err := check("byWeekNo", opt.ByWeekNo, 1, 53, true); err != nil {
		return err
	}
	if err := check("byMonth", opt.ByMonth, 1, 12, false); err != nil {
		return err
	}
	if err := check("bySetPos", opt.BySetPos, 1, 366, true); err != nil {
		return err
	}

	// Validate weekday positions
	for _, wd := range opt.ByDay {
		if wd.N() > 53 || wd.N() < -53 {
			return fmt.Errorf("byDay: weekday position must be between -53 and 53, got %d", wd.N())
		}
	}

	return nil
}

func (ext *ROptionExtend) ToROptionPlus() ROptionPlus {
	return ROptionPlus{
		ROption: rrule.ROption{
			Freq:       ext.Freq.ToFrequency(),
			Interval:   ext.Interval,
			Count:      ext.Count,
			Dtstart:    atime.EnsureDateTimeUTC(ext.BeginTime),
			Until:      atime.EnsureDateTimeUTC(ext.UntilTime),
			Bysecond:   ext.BySecond,
			Byminute:   ext.ByMinute,
			Byhour:     ext.ByHour,
			Byweekday:  ext.ByDay,
			Bymonthday: ext.ByMonthDay,
			Byyearday:  ext.ByYearDay,
			Byweekno:   ext.ByWeekNo,
			Bymonth:    ext.ByMonth,
			Bysetpos:   ext.BySetPos,
			Byeaster:   ext.ByEaster,
		},
		ShiftOffWeekend:     ext.ShiftOffWeekend,
		ShiftOffHolidays:    ext.ShiftOffHolidays,
		ValidOnlyOnHolidays: ext.ValidOnlyOnHolidays,
		ValidOnlyOnWeekends: ext.ValidOnlyOnWeekends,
		ISOCode:             ext.ISOCode,
		Observance:          ext.Observance,
		CustomFilter:        ext.CustomFilter,
		Calendar:            ext.Calendar,
	}
}

func (opt *ROptionExtend) Clone() *ROptionExtend {
	if opt == nil {
		return nil
	}

	clone := *opt // shallow copy

	// Deep copy time pointers
	if opt.BeginTime != nil {
		t := *opt.BeginTime
		clone.BeginTime = &t
	}
	if opt.UntilTime != nil {
		t := *opt.UntilTime
		clone.UntilTime = &t
	}

	// Deep copy slices
	clone.BySecond = append([]int{}, opt.BySecond...)
	clone.ByMinute = append([]int{}, opt.ByMinute...)
	clone.ByHour = append([]int{}, opt.ByHour...)
	clone.ByDay = append([]rrule.Weekday{}, opt.ByDay...)
	clone.ByMonthDay = append([]int{}, opt.ByMonthDay...)
	clone.ByYearDay = append([]int{}, opt.ByYearDay...)
	clone.ByWeekNo = append([]int{}, opt.ByWeekNo...)
	clone.ByMonth = append([]int{}, opt.ByMonth...)
	clone.BySetPos = append([]int{}, opt.BySetPos...)
	clone.ByEaster = append([]int{}, opt.ByEaster...)

	// CustomFilter and Calendar are function/interface types — intentionally not cloned
	return &clone
}

func describeROptions(opt ROptionExtend) []string {
	var out []string

	// Frequency and Interval
	if opt.Interval > 1 {
		out = append(out, fmt.Sprintf("Every %d %s", opt.Interval, strings.ToLower(opt.Freq.String())))
	} else {
		out = append(out, fmt.Sprintf("Every %s", strings.ToLower(opt.Freq.String())))
	}

	// Count / Until
	if opt.Count > 0 {
		out = append(out, fmt.Sprintf("Up to %d times", opt.Count))
	}
	if opt.UntilTime != nil && !opt.UntilTime.IsZero() {
		out = append(out, fmt.Sprintf("Until %s", opt.UntilTime.Format("2006-01-02")))
	}

	// Anchor start time
	if opt.BeginTime != nil && !opt.BeginTime.IsZero() {
		out = append(out, fmt.Sprintf("Start: %s", opt.BeginTime.Format("2006-01-02 15:04")))
	}

	// Time of day
	if len(opt.ByHour) > 0 || len(opt.ByMinute) > 0 {
		h := 0
		m := 0
		if len(opt.ByHour) > 0 {
			h = opt.ByHour[0]
		}
		if len(opt.ByMinute) > 0 {
			m = opt.ByMinute[0]
		}
		out = append(out, fmt.Sprintf("At %02d:%02d", h, m))
	}

	// ByDay
	if len(opt.ByDay) > 0 {
		dayLabels := make([]string, len(opt.ByDay))
		for i, d := range opt.ByDay {
			dayLabels[i] = d.String()
		}
		out = append(out, fmt.Sprintf("On days: %s", strings.Join(dayLabels, ", ")))
	}

	// ByMonthDay
	if len(opt.ByMonthDay) > 0 {
		out = append(out, fmt.Sprintf("On month days: %v", opt.ByMonthDay))
	}

	// ByMonth
	if len(opt.ByMonth) > 0 {
		out = append(out, fmt.Sprintf("In months: %v", opt.ByMonth))
	}

	// Additional BYxxx filters (if needed, can be gated with flags or verbosity level)
	if len(opt.ByYearDay) > 0 {
		out = append(out, fmt.Sprintf("Year days: %v", opt.ByYearDay))
	}
	if len(opt.ByWeekNo) > 0 {
		out = append(out, fmt.Sprintf("Week numbers: %v", opt.ByWeekNo))
	}
	if len(opt.BySetPos) > 0 {
		out = append(out, fmt.Sprintf("Set positions: %v", opt.BySetPos))
	}
	if len(opt.ByEaster) > 0 {
		out = append(out, fmt.Sprintf("Easter offsets: %v", opt.ByEaster))
	}

	// Policy extensions
	if opt.ShiftOffWeekend {
		out = append(out, "Shift off weekends")
	}
	if opt.ShiftOffHolidays {
		out = append(out, "Shift off holidays")
	}
	if opt.ValidOnlyOnHolidays {
		out = append(out, "Only on holidays")
	}
	if opt.ValidOnlyOnWeekends {
		out = append(out, "Only on weekends")
	}
	if opt.ISOCode != "" {
		out = append(out, fmt.Sprintf("Region: %s", opt.ISOCode))
	}
	if !opt.Observance.IsEmpty() {
		out = append(out, fmt.Sprintf("Observance: %s", opt.Observance))
	}

	return out
}
