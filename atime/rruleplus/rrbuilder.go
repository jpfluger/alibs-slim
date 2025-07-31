package rruleplus

import (
	"github.com/jpfluger/alibs-slim/atime"
	"github.com/teambition/rrule-go"
	"time"
)

type RRBuilderType string

const (
	// RRBUILDERTYPE_MONTHDAY represents rules that trigger on a specific month and day each year or month.
	// Typically used for recurring holidays (e.g. Dec 25), reminders (e.g. April 1), or birthdays.
	RRBUILDERTYPE_MONTHDAY RRBuilderType = "monthday"

	// RRBUILDERTYPE_ANYTIME represents rules that match all times, like "always allow" or "deny all".
	// These are unconditional policies, and can serve as global overrides or base-case fallbacks.
	RRBUILDERTYPE_ANYTIME RRBuilderType = "anytime"

	// RRBUILDERTYPE_WEEKDAY represents rules that trigger on specific days of the week (e.g. Mon–Fri)
	// and may also include specific hour ranges (e.g. 9am–5pm). Used for office hours or scheduled access blocks.
	RRBUILDERTYPE_WEEKDAY RRBuilderType = "weekday"

	// RRBUILDERTYPE_SPECIFIC_DATE represents rules that match a fixed calendar date (e.g. April 15).
	// These are useful for events that occur once a year and do not follow a recurrence pattern like weekdays.
	RRBUILDERTYPE_SPECIFIC_DATE RRBuilderType = "specific-date"

	// RRBUILDERTYPE_NTHWEEKDAY represents ordinal weekday patterns like "first Monday" or "last Friday".
	// Commonly used for scheduling monthly meetings, billing cycles, or recurring governance events.
	RRBUILDERTYPE_NTHWEEKDAY RRBuilderType = "nth-weekday"
)

// RRBuilder is a fluent builder-style DSL for constructing complex recurrence rules using RRuleExtend.
// It encapsulates time-based logic such as repeat frequencies, window durations, and join triggers.
//
// # Supported Builder Types
//
// The builder supports a variety of recurrence strategies via the RRBuilderType enum:
//
//   - RRBUILDERTYPE_MONTHDAY – Recurs on a specific calendar day (e.g. December 25).
//   - RRBUILDERTYPE_WEEKDAY  – Matches days of the week with optional hour ranges (e.g. Mon–Fri 9–5).
//   - RRBUILDERTYPE_NTHWEEKDAY – Matches patterns like "first Monday of the month".
//   - RRBUILDERTYPE_SPECIFIC_DATE – Matches a fixed calendar date once per year.
//   - RRBUILDERTYPE_ANYTIME  – Matches all time, useful for unconditional allow/deny logic.
//
// # Key Configuration
//
// Builders allow optional chaining of modifiers:
//
//   - .WithDuration(d, unit)         – Sets how long each matched window remains active.
//   - .WithJoinWindowBefore(...)     – Defines pre-activation notification windows.
//   - .Allow() / .Deny()             – Controls access polarity.
//   - .WithPriority(p)               – Sorting priority (used when evaluating stacked rules).
//
// # Advanced Features
//
//   - .BuildSpecificDateStack() – Allows fallback logic for shifting dates to business days
//     and skipping holidays (e.g. April 15 → April 16 → April 17).
//
// # Example Usage
//
//	rule := NewRRBuilderMonthDay(time.Date(2025, 12, 25, 0, 0, 0, 0, time.UTC), TIMEUNIT_YEARLY).
//		WithDuration(1, TIMEUNIT_DAILY).
//		WithJoinWindowBefore([]int{30, 15, 1}, TIMEUNIT_DAILY, "reminder").
//		Allow().
//		WithPriority(10).
//		Build()
//
// This constructs a rule that triggers annually on Dec 25, runs for 1 day,
// and supports early notification join windows of 30, 15, and 1 day in advance.
//
// # Stacked Rule Evaluation
//
// Rules can be grouped into a `RRuleExtends` slice, and evaluated using `.Evaluate(now)`,
// where allow/deny decisions are determined by priority and match results.
//
// This makes RRBuilder ideal for implementing layered access control, calendar blocking,
// recurring notifications, or policy scheduling engines.
type RRBuilder struct {
	rule        RRuleExtend
	builderType RRBuilderType
	timeZone    *time.Location

	// Temporary config for JoinWindow construction
	joinWindowConfig *joinWindowParams

	// Behavior flags
	excludeHolidays     map[int]bool
	weekendFallbackDays int
	weekendShiftTarget  *time.Weekday
}

type joinWindowParams struct {
	durations []int
	unit      atime.TimeUnit
	tag       string
}

func (b *RRBuilder) WithDuration(d int, unit atime.TimeUnit) *RRBuilder {
	if d <= 0 {
		d = 1
	}
	if unit.IsEmpty() {
		unit = atime.TIMEUNIT_DAILY
	}
	b.rule.Duration = d
	b.rule.DurationUnit = unit
	return b
}

func (b *RRBuilder) AddJWBefore(unit atime.TimeUnit, tag string, durations ...int) *RRBuilder {
	if unit.IsEmpty() {
		unit = atime.TIMEUNIT_DAILY
	}
	b.joinWindowConfig = &joinWindowParams{
		durations: durations,
		unit:      unit,
		tag:       tag,
	}
	return b
}

func (b *RRBuilder) AddJWBeforeDaily(tag string, durations ...int) *RRBuilder {
	b.joinWindowConfig = &joinWindowParams{
		durations: durations,
		unit:      atime.TIMEUNIT_DAILY,
		tag:       tag,
	}
	return b
}

func (b *RRBuilder) AddJWBeforeHourly(tag string, durations ...int) *RRBuilder {
	b.joinWindowConfig = &joinWindowParams{
		durations: durations,
		unit:      atime.TIMEUNIT_HOURLY,
		tag:       tag,
	}
	return b
}

func (b *RRBuilder) AddJWBeforeMinutely(tag string, durations ...int) *RRBuilder {
	b.joinWindowConfig = &joinWindowParams{
		durations: durations,
		unit:      atime.TIMEUNIT_MINUTELY,
		tag:       tag,
	}
	return b
}

func (b *RRBuilder) WithPriority(p int) *RRBuilder {
	b.rule.Priority = p
	return b
}

func (b *RRBuilder) Allow() *RRBuilder {
	b.rule.IsDeny = false
	return b
}

func (b *RRBuilder) Deny() *RRBuilder {
	b.rule.IsDeny = true
	return b
}

func (b *RRBuilder) WithInterval(interval int) *RRBuilder {
	if interval < 0 {
		interval = 1
	}
	b.rule.ROptions.Interval = interval
	return b
}

func (b *RRBuilder) WithCount(count int) *RRBuilder {
	if count < 0 {
		count = 1
	}
	b.rule.ROptions.Count = count
	return b
}

func (b *RRBuilder) WithShiftOffHolidays() *RRBuilder {
	b.rule.ROptions.ShiftOffHolidays = true
	return b
}

func (b *RRBuilder) WithShiftOffWeekend() *RRBuilder {
	b.rule.ROptions.ShiftOffWeekend = true
	return b
}

func (b *RRBuilder) WithObservance(obs ObservanceMode) *RRBuilder {
	b.rule.ROptions.Observance = obs
	return b
}

func (b *RRBuilder) WithISOCode(iso string) *RRBuilder {
	b.rule.ROptions.ISOCode = iso
	return b
}

func (b *RRBuilder) SetBeginByYear(year int) *RRBuilder {
	if year <= 0 {
		year = time.Now().UTC().Year()
	}
	// Use Jan 1st at 00:00 UTC
	t := time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
	b.rule.ROptions.BeginTime = &t
	return b
}

func (b *RRBuilder) WithDays(days ...rrule.Weekday) *RRBuilder {
	b.rule.ROptions.ByDay = days
	return b
}

func (b *RRBuilder) WithHourRange(startHour, endHour int) *RRBuilder {
	if b.timeZone == nil {
		b.timeZone = time.UTC
	}

	// Pick reference date: prefer BeginTime if already set
	var reference time.Time
	if b.rule.ROptions.BeginTime != nil && !b.rule.ROptions.BeginTime.IsZero() {
		reference = (*b.rule.ROptions.BeginTime).In(b.timeZone)
	} else {
		// Fallback: safe default for anchoring (Jan 1 of current year)
		reference = time.Date(time.Now().Year(), time.January, 1, 0, 0, 0, 0, b.timeZone)
	}

	// Convert startHour to UTC-aligned hour using local offset
	startTime := time.Date(reference.Year(), reference.Month(), reference.Day(), startHour, 0, 0, 0, b.timeZone)
	startUTC := startTime.UTC().Hour()

	var hours []int
	if endHour < 0 {
		// Single anchor
		hours = []int{startUTC}
	} else {
		endTime := time.Date(reference.Year(), reference.Month(), reference.Day(), endHour, 0, 0, 0, b.timeZone)
		endUTC := endTime.UTC().Hour()
		hours = atime.HoursBetween(startUTC, endUTC, true)
	}

	b.rule.ROptions.ByHour = hours

	// If BeginTime hasn't been explicitly set, anchor it here
	if b.rule.ROptions.BeginTime == nil || b.rule.ROptions.BeginTime.IsZero() {
		utcAnchor := time.Date(reference.Year(), reference.Month(), reference.Day(), startUTC, 0, 0, 0, time.UTC)
		b.rule.ROptions.BeginTime = atime.ToPointer(utcAnchor)
	}

	return b
}

func (b *RRBuilder) OnJanuary() *RRBuilder   { return b.WithByMonth(1) }
func (b *RRBuilder) OnFebruary() *RRBuilder  { return b.WithByMonth(2) }
func (b *RRBuilder) OnMarch() *RRBuilder     { return b.WithByMonth(3) }
func (b *RRBuilder) OnApril() *RRBuilder     { return b.WithByMonth(4) }
func (b *RRBuilder) OnMay() *RRBuilder       { return b.WithByMonth(5) }
func (b *RRBuilder) OnJune() *RRBuilder      { return b.WithByMonth(6) }
func (b *RRBuilder) OnJuly() *RRBuilder      { return b.WithByMonth(7) }
func (b *RRBuilder) OnAugust() *RRBuilder    { return b.WithByMonth(8) }
func (b *RRBuilder) OnSeptember() *RRBuilder { return b.WithByMonth(9) }
func (b *RRBuilder) OnOctober() *RRBuilder   { return b.WithByMonth(10) }
func (b *RRBuilder) OnNovember() *RRBuilder  { return b.WithByMonth(11) }
func (b *RRBuilder) OnDecember() *RRBuilder  { return b.WithByMonth(12) }

func (b *RRBuilder) WithBySecond(seconds ...int) *RRBuilder {
	b.rule.ROptions.BySecond = append([]int{}, seconds...)
	return b
}

func (b *RRBuilder) WithByMinute(minutes ...int) *RRBuilder {
	b.rule.ROptions.ByMinute = append([]int{}, minutes...)
	return b
}

func (b *RRBuilder) WithByHour(hours ...int) *RRBuilder {
	b.rule.ROptions.ByHour = append([]int{}, hours...)
	return b
}

func (b *RRBuilder) WithByDay(days ...rrule.Weekday) *RRBuilder {
	b.rule.ROptions.ByDay = append([]rrule.Weekday{}, days...)
	return b
}

func (b *RRBuilder) WithByMonthDay(days ...int) *RRBuilder {
	b.rule.ROptions.ByMonthDay = append([]int{}, days...)
	return b
}

func (b *RRBuilder) WithByYearDay(days ...int) *RRBuilder {
	b.rule.ROptions.ByYearDay = append([]int{}, days...)
	return b
}

func (b *RRBuilder) WithByWeekNo(weeks ...int) *RRBuilder {
	b.rule.ROptions.ByWeekNo = append([]int{}, weeks...)
	return b
}

func (b *RRBuilder) WithByMonth(months ...int) *RRBuilder {
	b.rule.ROptions.ByMonth = append([]int{}, months...)
	return b
}

func (b *RRBuilder) WithBySetPos(positions ...int) *RRBuilder {
	b.rule.ROptions.BySetPos = append([]int{}, positions...)
	return b
}

func (b *RRBuilder) WithByEaster(offsets ...int) *RRBuilder {
	b.rule.ROptions.ByEaster = append([]int{}, offsets...)
	return b
}

func (b *RRBuilder) OnDayOfMonth(day int) *RRBuilder {
	return b.WithByMonthDay(day)
}

func (b *RRBuilder) attachJoinWindows() {
	if b.joinWindowConfig == nil {
		return
	}
	var jws JoinWindows
	for _, dur := range b.joinWindowConfig.durations {
		if dur <= 0 {
			continue
		}
		jws = append(jws, &JoinWindow{
			IsBefore:     true,
			Duration:     dur,
			DurationUnit: b.joinWindowConfig.unit,
			Tag:          b.joinWindowConfig.tag,
		})
	}
	b.rule.JoinWindows = jws
}

// WithTimeZone sets the desired IANA timezone (e.g. "America/New_York") for interpreting hour-based settings.
func (b *RRBuilder) WithTimeZone(tz string) *RRBuilder {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		// Fallback to UTC if invalid
		loc = time.UTC
	}
	b.timeZone = loc
	return b
}

func (b *RRBuilder) WithBeginTime(target time.Time) *RRBuilder {
	// Normalize to builder timeZone if it's set
	if b.timeZone != nil {
		target = target.In(b.timeZone)
	}
	// Store as UTC (rrule expects UTC)
	b.rule.ROptions.BeginTime = atime.ToPointer(target.UTC())
	return b
}

func (b *RRBuilder) WithStandardWorkWeek() *RRBuilder {
	// Traditional business week: 9AM–5PM, Mon–Fri
	return b.WithHourRange(9, 17).WithDays(rrule.MO, rrule.TU, rrule.WE, rrule.TH, rrule.FR)
}

func (b *RRBuilder) AvoidWeekends() *RRBuilder {
	b.rule.ROptions.ByDay = excludeDays(b.rule.ROptions.ByDay, rrule.SA, rrule.SU)
	return b
}

// excludeDays removes any of the specified daysToExclude from the input slice.
func excludeDays(input []rrule.Weekday, daysToExclude ...rrule.Weekday) []rrule.Weekday {
	if len(input) == 0 || len(daysToExclude) == 0 {
		return input
	}

	excludeSet := make(map[string]bool)
	for _, d := range daysToExclude {
		excludeSet[d.String()] = true
	}

	var result []rrule.Weekday
	for _, d := range input {
		if !excludeSet[d.String()] {
			result = append(result, d)
		}
	}
	return result
}

func (b *RRBuilder) WithWeekendFallback(days int) *RRBuilder {
	b.weekendFallbackDays = days
	return b
}

func (b *RRBuilder) WithShiftOnWeekend(target time.Weekday) *RRBuilder {
	b.weekendShiftTarget = &target
	return b
}

func (b *RRBuilder) ExcludeHolidays(dates ...time.Time) *RRBuilder {
	b.excludeHolidays = make(map[int]bool)
	for _, d := range dates {
		key := d.Year()*10000 + int(d.Month())*100 + d.Day()
		b.excludeHolidays[key] = true
	}
	return b
}

func (b *RRBuilder) ExcludeHolidaysMap(ymd map[int]bool) *RRBuilder {
	if ymd == nil {
		return b
	}
	b.excludeHolidays = ymd
	return b
}

// FirstNthWeekday returns the date of the Nth occurrence of a weekday in a given month and year.
// For example, the 2nd Monday of March 2025 would be: FirstNthWeekday(2025, time.March, time.Monday, 2)
func (b *RRBuilder) FirstNthWeekday(year int, month time.Month, weekday time.Weekday, nth int) time.Time {
	count := 0
	for day := 1; day <= 31; day++ {
		date := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
		if date.Month() != month {
			break
		}
		if date.Weekday() == weekday {
			count++
			if count == nth {
				return date
			}
		}
	}
	return time.Time{} // return zero if no match
}

// FirstNthWeekday returns the date of the Nth occurrence of a weekday in a given month and year.
// For example, the 2nd Monday of March 2025 would be: FirstNthWeekday(2025, time.March, time.Monday, 2)
func FirstNthWeekday(year int, month time.Month, weekday time.Weekday, nth int) time.Time {
	count := 0
	for day := 1; day <= 31; day++ {
		date := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
		if date.Month() != month {
			break
		}
		if date.Weekday() == weekday {
			count++
			if count == nth {
				return date
			}
		}
	}
	return time.Time{}
}

func (b *RRBuilder) Build() *RRuleExtend {
	b.attachJoinWindows()

	switch b.builderType {
	case RRBUILDERTYPE_MONTHDAY:
		return b.buildMonthDay()
	case RRBUILDERTYPE_ANYTIME:
		return b.buildAnyTime()
	case RRBUILDERTYPE_WEEKDAY:
		return b.buildWeekday()
	case RRBUILDERTYPE_SPECIFIC_DATE:
		return b.buildSpecificDate()
	case RRBUILDERTYPE_NTHWEEKDAY:
		return b.buildNthWeekday()
	default:
		panic("unsupported RRBuilderType: " + string(b.builderType))
	}
}

// NewRRBuilderMonthDay creates a recurring rule that triggers on a specific day-of-month.
// Common use: "Every Dec 25", "The 1st of each month", etc.
// Supports join windows and duration (e.g. alerts before billing date).
func NewRRBuilderMonthDay(date time.Time, freq atime.TimeUnit) *RRBuilder {
	if date.IsZero() {
		date = time.Now().UTC()
	} else {
		date = date.UTC()
	}
	if freq != atime.TIMEUNIT_YEARLY && freq != atime.TIMEUNIT_MONTHLY {
		freq = atime.TIMEUNIT_MONTHLY
	}
	return &RRBuilder{
		builderType: RRBUILDERTYPE_MONTHDAY,
		rule: RRuleExtend{
			ROptions: ROptionExtend{
				Freq:      freq,
				Interval:  1,
				BeginTime: &date,
			},
		},
	}
}

func (b *RRBuilder) buildMonthDay() *RRuleExtend {
	bt := b.rule.ROptions.BeginTime.UTC()
	b.rule.ROptions.ByMonth = []int{int(bt.Month())}
	b.rule.ROptions.ByMonthDay = []int{bt.Day()}
	return &b.rule
}

// NewRRBuilderAnyTime creates a rule that is always active, regardless of date or time.
// Use .Allow() or .Deny() to define intent.
// Does not support duration or join windows.
func NewRRBuilderAnyTime() *RRBuilder {
	return &RRBuilder{
		builderType: RRBUILDERTYPE_ANYTIME,
	}
}

func (b *RRBuilder) buildAnyTime() *RRuleExtend {
	b.rule.IsAnyTime = true
	return &b.rule
}

// NewRRBuilderWeekday creates a recurrence rule targeting specific weekdays with optional hour ranges.
//
// Usage patterns:
//   - Fixed hour range: e.g., Monday–Friday from 9am–5pm.
//   - Single start hour: e.g., Tuesday at 9am, with a custom duration.
//   - Defaults to TIMEUNIT_DAILY frequency, and supports chaining of duration and join windows.
//
// Parameters:
//   - `days`: slice of `rrule.Weekday` (e.g., `rrule.MO, rrule.TU, ...`)
//   - `startHour`: beginning hour (0–23); required
//   - `endHour`: end hour (0–23). If < 0, only `startHour` is used, and duration must be set explicitly.
//
// Examples:
//
//	NewRRBuilderWeekday([]rrule.Weekday{rrule.MO}, 9, 17) → 9am–5pm
//	NewRRBuilderWeekday([]rrule.Weekday{rrule.FR}, 9, -1).WithDuration(8, TIMEUNIT_HOURLY) → 9am to 5pm (single anchor)
func NewRRBuilderWeekday(year int, days []rrule.Weekday, startHour, endHour int) *RRBuilder {
	if year <= 0 {
		year = time.Now().UTC().Year()
	}

	base := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	var hours []int
	if endHour >= 0 {
		hours = atime.HoursBetween(startHour, endHour, true)
	} else {
		hours = []int{startHour}
	}
	return &RRBuilder{
		builderType: RRBUILDERTYPE_WEEKDAY,
		rule: RRuleExtend{
			ROptions: ROptionExtend{
				Freq:      atime.TIMEUNIT_DAILY,
				ByDay:     days,
				ByHour:    hours,
				BeginTime: &base,
				Interval:  1,
			},
		},
	}
}

func (b *RRBuilder) buildWeekday() *RRuleExtend {
	return &b.rule
}

// NewRRBuilderSpecificDate creates a recurring rule for an exact calendar date.
// This is useful for rules like "Every July 4" or "On Dec 25 each year".
// Supports duration and join windows.
func NewRRBuilderSpecificDate(date time.Time) *RRBuilder {
	date = date.UTC()
	begin := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)

	return &RRBuilder{
		builderType: RRBUILDERTYPE_SPECIFIC_DATE,
		rule: RRuleExtend{
			ROptions: ROptionExtend{
				Freq:       atime.TIMEUNIT_YEARLY,
				Interval:   1,
				BeginTime:  &begin,
				ByMonth:    []int{int(date.Month())},
				ByMonthDay: []int{date.Day()},
			},
		},
	}
}

func (b *RRBuilder) buildSpecificDate() *RRuleExtend {
	return &b.rule
}

// NewRRBuilderNthWeekday creates a rule for patterns like
// "the 2nd Tuesday of every month" or "last Friday each month".
// It accepts a start year and calculates the anchor BeginTime accordingly.
func NewRRBuilderNthWeekday(startYear int, nth int, weekday time.Weekday) *RRBuilder {
	if startYear <= 0 {
		startYear = time.Now().UTC().Year()
	}

	rrWDay := toRRuleNth(nth, weekday)
	begin := time.Date(startYear, 1, 1, 0, 0, 0, 0, time.UTC)

	return &RRBuilder{
		builderType: RRBUILDERTYPE_NTHWEEKDAY,
		rule: RRuleExtend{
			ROptions: ROptionExtend{
				Freq:      atime.TIMEUNIT_MONTHLY,
				Interval:  1,
				BeginTime: &begin,
				ByDay:     []rrule.Weekday{rrWDay},
			},
		},
	}
}

func (b *RRBuilder) buildNthWeekday() *RRuleExtend {
	// If hours are set, align BeginTime to first defined hour
	if len(b.rule.ROptions.ByHour) > 0 && b.rule.ROptions.BeginTime != nil {
		bt := *b.rule.ROptions.BeginTime
		b.rule.ROptions.BeginTime = atime.ToPointer(time.Date(
			bt.Year(), bt.Month(), bt.Day(),
			b.rule.ROptions.ByHour[0], 0, 0, 0, time.UTC,
		))
	}
	return &b.rule
}

func toRRuleNth(nth int, weekday time.Weekday) rrule.Weekday {
	wdConverted := atime.TimeWeekdayToRRuleWeekday(weekday)
	return wdConverted.Nth(nth)
}

func (b *RRBuilder) createWeekendFallbackRules(baseDate time.Time, basePriority, maxDays int) RRuleExtends {
	var fallbackRules RRuleExtends

	for i := 1; i <= maxDays; i++ {
		next := baseDate.AddDate(0, 0, i)

		// If shift target is defined, jump to that weekday
		if b.weekendShiftTarget != nil {
			offset := int(*b.weekendShiftTarget - baseDate.Weekday())
			if offset <= 0 {
				offset += 7
			}
			next = baseDate.AddDate(0, 0, offset)
			i = maxDays // break after 1 fallback
		}

		// Skip weekends
		if next.Weekday() == time.Saturday || next.Weekday() == time.Sunday {
			continue
		}

		// Skip excluded holidays
		if b.excludeHolidays != nil {
			ymd := next.Year()*10000 + int(next.Month())*100 + next.Day()
			if b.excludeHolidays[ymd] {
				continue
			}
		}

		builder := NewRRBuilderSpecificDate(next).
			Allow().
			WithDuration(b.rule.Duration, b.rule.DurationUnit).
			WithPriority(basePriority + (maxDays - i))

		if b.joinWindowConfig != nil {
			builder = builder.AddJWBefore(
				b.joinWindowConfig.unit,
				b.joinWindowConfig.tag,
				b.joinWindowConfig.durations...,
			)
		}

		fallbackRules = append(fallbackRules, builder.Build())
		break // use only the first valid fallback
	}

	return fallbackRules
}

// BuildSpecificDateStack creates a layered set of rules for a fixed date.
// Includes logic to shift or fallback if the date falls on a weekend or holiday.
func (b *RRBuilder) BuildSpecificDateStack() RRuleExtends {
	var rules RRuleExtends
	priority := 10

	// Validate base date
	if b.rule.ROptions.BeginTime == nil || b.rule.ROptions.BeginTime.IsZero() {
		return rules
	}
	baseDate := *b.rule.ROptions.BeginTime

	// 1. Base allow rule for specific date
	base := NewRRBuilderSpecificDate(baseDate).
		Allow().
		WithDuration(b.rule.Duration, b.rule.DurationUnit).
		WithPriority(priority)

	if b.joinWindowConfig != nil {
		base = base.AddJWBefore(
			b.joinWindowConfig.unit,
			b.joinWindowConfig.tag,
			b.joinWindowConfig.durations...,
		)
	}

	rules = append(rules, base.Build())

	// 2. Optional: Deny rule for weekends (to suppress base)
	if baseDate.Weekday() == time.Saturday || baseDate.Weekday() == time.Sunday {
		rules = append(rules,
			NewRRBuilderWeekday(baseDate.Year(), []rrule.Weekday{rrule.SA, rrule.SU}, 0, 24).
				Deny().
				WithPriority(priority+90).
				Build(),
		)
	}

	// 3. Optional: Add fallback rules (shift or next valid weekday)
	if b.weekendFallbackDays > 0 && (baseDate.Weekday() == time.Saturday || baseDate.Weekday() == time.Sunday) {
		fallbacks := b.createWeekendFallbackRules(baseDate, priority, b.weekendFallbackDays)
		rules = append(rules, fallbacks...)
	}

	return rules
}
