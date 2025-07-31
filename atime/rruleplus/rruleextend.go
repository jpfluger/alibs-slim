package rruleplus

import (
	"fmt"
	"github.com/jpfluger/alibs-slim/ageo"
	"github.com/jpfluger/alibs-slim/atime"
	"sort"
	"strings"
	"time"
)

type IRRuleEvaluator interface {
	IsPreAllowed(now time.Time, geo ageo.GeoInfo) (RREvaluatorResultType, error)
	IsAllowed(now time.Time, geo ageo.GeoInfo) error
}

type RRuleExtend struct {
	IsDeny   bool `json:"isDeny,omitempty"`   // If true, this rule acts as a DENY policy. If it matches, access is immediately denied.
	Priority int  `json:"priority,omitempty"` // Lower = evaluated earlier

	Name string `json:"name,omitempty"` // Optional name for the rule

	IsAnyTime bool `json:"isAnyTime,omitempty"` // If true, time checks automatically succeed.

	// Hard begin/stop.
	StartDate *time.Time `json:"startDate,omitempty"` // Optional hard window start
	EndDate   *time.Time `json:"endDate,omitempty"`   // Optional hard window end

	// RRule: Core recurrence rule config
	ROptions ROptionExtend `json:"rOptions,omitempty"` // ROptions + ROptionsPlus

	// Join window config
	JoinWindows JoinWindows `json:"joinWindows,omitempty"`

	// The effective window of the occurrence.
	Duration     int            `json:"duration,omitempty"`
	DurationUnit atime.TimeUnit `json:"durationUnit,omitempty"`

	// Optional geo filtering
	GeoFilters ageo.GeoFilters `json:"geoFilters,omitempty"`
}

func (rre *RRuleExtend) Validate() error {
	rre.Name = strings.TrimSpace(rre.Name)

	if rre.IsAnyTime {
		return nil
	}

	if err := rre.ROptions.ValidateRecurrence(); err != nil {
		return err
	}

	if rre.Duration > 0 && rre.DurationUnit.IsEmpty() {
		return fmt.Errorf("rrule duration unit must be set when duration > 0")
	}

	rre.JoinWindows = rre.JoinWindows.Sanitize()

	return nil
}

// Clone returns a deep copy of the provided RRuleExtend.
func (src *RRuleExtend) Clone() *RRuleExtend {
	if src == nil {
		return nil
	}

	clone := *src // shallow copy of everything except pointer/slice fields

	// Deep copy time pointers
	if src.StartDate != nil {
		t := *src.StartDate
		clone.StartDate = &t
	}
	if src.EndDate != nil {
		t := *src.EndDate
		clone.EndDate = &t
	}

	clone.ROptions = *src.ROptions.Clone()

	// Deep copy JoinWindows
	if len(src.JoinWindows) > 0 {
		clone.JoinWindows = make(JoinWindows, len(src.JoinWindows))
		for i, jw := range src.JoinWindows {
			if jw != nil {
				jwCopy := *jw
				clone.JoinWindows[i] = &jwCopy
			}
		}
	}

	// Deep copy GeoFilters if needed
	if gf := src.GeoFilters; gf != nil {
		clone.GeoFilters = gf.Clone() // assuming Clone() exists for ageo.GeoFilters
	}

	return &clone
}

// IsBetween applies recurrence logic to the current time and honors `IsDeny`.
func (rre *RRuleExtend) IsBetween(now time.Time) (bool, error) {
	match, err := rre.matchesOccurrenceWindow(now)
	if err != nil {
		return false, err
	}
	if rre.IsDeny {
		return !match, nil
	}
	return match, nil
}

// matchesOccurrenceWindow checks whether `now` is within any recurrence instance,
// including JoinWindowBefore and JoinWindowAfter windows.
// It does not apply `IsDeny`.
func (rre *RRuleExtend) matchesOccurrenceWindow(now time.Time) (bool, error) {
	if rre.IsAnyTime {
		return !rre.IsDeny, nil
	}

	now = now.UTC()

	// Check hard boundaries first
	if rre.StartDate != nil && now.Before(atime.EnsureDateTimeUTC(rre.StartDate)) {
		return false, nil
	}
	if rre.EndDate != nil && now.After(atime.EnsureDateTimeUTC(rre.EndDate)) {
		return false, nil
	}

	if rre.DurationUnit.IsEmpty() {
		return false, fmt.Errorf("DurationUnit must be set")
	}

	// If time anchors are present, use refined matching
	if rre.shouldUseTimeAnchor() {
		return rre.matchesTimeAnchoredWindow(now)
	}

	// Standard recurrence-based handling
	rule, err := rre.ToRRule()
	if err != nil {
		return false, err
	}

	last := rule.Before(now, true)
	if last.IsZero() {
		return false, nil
	}

	end := rre.calculateEndTime(last)
	return !now.Before(last) && now.Before(end), nil
}

func (rre *RRuleExtend) matchesTimeAnchoredWindow(now time.Time) (bool, error) {
	rule, err := rre.ToRRule()
	if err != nil {
		return false, err
	}

	last := rule.Before(now, rre.shouldUseInclusiveMatch())
	if last.IsZero() {
		return false, nil
	}

	anchor := time.Date(
		last.Year(), last.Month(), last.Day(),
		last.Hour(), last.Minute(), last.Second(), 0, time.UTC,
	)
	start := anchor
	end := rre.calculateEndTime(anchor)

	// SKIPPING HERE!!! (see MatchJoinWindow)
	// Optionally apply JoinWindows if this is a primary matcher
	// (skip this if JoinWindow should only influence join-check)
	//for _, jw := range rre.JoinWindow {
	//	if jw == nil || jw.DurationUnit.IsEmpty() {
	//		continue
	//	}
	//	dur := jw.DurationUnit.CalcDuration(jw.Duration)
	//	if jw.IsBefore {
	//		start = start.Add(-dur)
	//	} else {
	//		end = end.Add(dur)
	//	}
	//}

	fmt.Printf("TimeAnchored match: now=%v start=%v end=%v (anchor=%v)\n", now, start, end, anchor)
	return !now.Before(start) && now.Before(end), nil
}

// shouldUseTimeAnchor determines whether time-based anchors are in use.
func (rre *RRuleExtend) shouldUseTimeAnchor() bool {
	return len(rre.ROptions.ByHour) > 0 || len(rre.ROptions.ByMinute) > 0 || len(rre.ROptions.BySecond) > 0
}

func (rre *RRuleExtend) calculateEndTime(start time.Time) time.Time {
	switch rre.DurationUnit {
	case atime.TIMEUNIT_SECONDLY, atime.TIMEUNIT_MINUTELY, atime.TIMEUNIT_HOURLY, atime.TIMEUNIT_DAILY, atime.TIMEUNIT_WEEKLY:
		return start.Add(rre.DurationUnit.CalcDuration(rre.Duration))
	case atime.TIMEUNIT_MONTHLY:
		d := rre.Duration
		if d <= 0 {
			d = 1
		}
		return start.AddDate(0, d, 0)
	case atime.TIMEUNIT_YEARLY:
		d := rre.Duration
		if d <= 0 {
			d = 1
		}
		return start.AddDate(d, 0, 0)
	default:
		// Safe fallback — acts like 0 duration
		return start
	}
}

func (rre *RRuleExtend) shouldUseInclusiveMatch() bool {
	if !rre.ROptions.RRIncType.IsEmpty() {
		return rre.ROptions.RRIncType == RRULE_INC_TYPE_INCLUSIVE
	}
	return true // default to inclusive for all frequencies
}

func (rre *RRuleExtend) WithCustomFilter(fn func(time.Time) bool) *RRuleExtend {
	rre.ROptions.WithCustomFilter(fn)
	return rre
}

func (rre *RRuleExtend) WithCalendar(cal ICalendar) *RRuleExtend {
	rre.ROptions.WithCalendar(cal)
	return rre
}

// ToRRule builds the underlying rrule.RRule from RRuleExtend returning and enhanced ToRRulePlus.
func (rre *RRuleExtend) ToRRule() (*RRulePlus, error) {
	rop := rre.ROptions.ToROptionPlus()
	return NewRRulePlus(rop)
}

func (rre *RRuleExtend) matchesOccurrenceWindowWithOptions(now time.Time, geo ageo.GeoInfo, eval IRRuleEvaluator) (bool, error) {
	// Evaluator check (Pre)
	if eval != nil {
		result, err := eval.IsPreAllowed(now, geo)
		if err != nil {
			return false, err
		}
		switch result {
		case RREVALUATOR_RESULTTYPE_ALLOW:
			return true, nil
		case RREVALUATOR_RESULTTYPE_DENY:
			return false, nil
		case RREVALUATOR_RESULTTYPE_CONTINUE:
			// fall through to normal IsBetween/match logic
		default:
			break
		}
	}

	// Do NOT apply IsDeny here
	ok, err := rre.matchesOccurrenceWindow(now)
	if err != nil || !ok {
		return false, err
	}

	// Evaluator check
	if eval != nil {
		if err := eval.IsAllowed(now, geo); err != nil {
			return false, nil
		}
	}

	// Geo filtering via stacked rules
	if len(rre.GeoFilters) > 0 {
		if !rre.GeoFilters.Evaluate(geo) {
			return false, nil
		}
	}

	return true, nil
}

func (rre *RRuleExtend) IsBetweenWithOptions(now time.Time, geo ageo.GeoInfo, eval IRRuleEvaluator) (bool, error) {
	// Do NOT apply IsDeny here
	match, err := rre.matchesOccurrenceWindowWithOptions(now, geo, eval)
	if err != nil {
		return false, err
	}
	if rre.IsDeny {
		return !match, nil
	}
	return match, nil
}

// MatchJoinWindowNext returns the JoinWindow that matches 'now' for the next recurrence.
// This is useful for upcoming notifications or grace periods after a scheduled time.
func (rre *RRuleExtend) MatchJoinWindowNext(now time.Time) (*JoinWindow, error) {
	return rre.MatchJoinWindow(now, false, true)
}

// MatchJoinWindowPrevious returns the JoinWindow that matches 'now' for the previous recurrence.
// This is useful for catching late actions or retroactive matching (e.g., user missed window).
func (rre *RRuleExtend) MatchJoinWindowPrevious(now time.Time) (*JoinWindow, error) {
	return rre.MatchJoinWindow(now, true, false)
}

// MatchJoinWindow returns a JoinWindow that matches the current time based on
// the closest recurrence either before or after 'now', depending on flags.
// - includePrevious: considers the previous occurrence (rule.Before)
// - includeNext:     considers the next occurrence (rule.After)
// Returns the first matching JoinWindow found (prefers previous match if both).
func (rre *RRuleExtend) MatchJoinWindow(now time.Time, includePrevious, includeNext bool) (*JoinWindow, error) {
	if rre.IsAnyTime || len(rre.JoinWindows) == 0 {
		return nil, nil
	}

	rule, err := rre.ToRRule()
	if err != nil {
		return nil, err
	}

	if includePrevious {
		prev := rule.Before(now, false) // false: exclusive of now
		if !prev.IsZero() {
			if jw := rre.JoinWindows.Matches(now, prev); jw != nil {
				return jw, nil
			}
		}
	}

	if includeNext {
		next := rule.After(now, true) // true: inclusive of now
		if !next.IsZero() {
			if jw := rre.JoinWindows.Matches(now, next); jw != nil {
				return jw, nil
			}
		}
	}

	return nil, nil
}

// GetNextTimes returns the next `count` recurrence times starting from `now`.
// The times returned do not evaluate negation (.IsDeny) of times.
// See GetNextOccurrences for a more enhanced analyzer.
// Uses RRule.After repeatedly since the iterator does not expose .Next().
func (rre *RRuleExtend) GetNextTimes(now time.Time, count int) ([]time.Time, error) {
	if count <= 0 {
		return nil, nil
	}

	// If it's an IsAnyTime rule, simulate a single time (or none)
	if rre.IsAnyTime {
		return []time.Time{now}, nil
	}

	rule, err := rre.ToRRule()
	if err != nil {
		return nil, err
	}

	var results []time.Time
	cursor := now

	for i := 0; i < count; i++ {
		next := rule.After(cursor, i == 0) // inclusive on first iteration
		if next.IsZero() {
			break
		}
		results = append(results, next)
		cursor = next.Add(time.Second) // step forward
	}

	return results, nil
}

// GetNextOccurrences returns the next N recurrence instances as structured RROccurrences,
// preserving rule metadata (IsDeny, Priority, Name).
func (rre *RRuleExtend) GetNextOccurrences(now time.Time, count int) (RROccurrences, error) {
	if count <= 0 {
		return nil, nil
	}

	// Special case: IsAnyTime returns current time once as a pseudo-occurrence.
	if rre.IsAnyTime {
		return RROccurrences{&RROccurrence{
			Time:     now,
			IsDeny:   rre.IsDeny,
			Priority: rre.Priority,
			Name:     rre.Name,
		}}, nil
	}

	rule, err := rre.ToRRule()
	if err != nil {
		return nil, err
	}

	var results RROccurrences
	cursor := now

	for i := 0; i < count; i++ {
		next := rule.After(cursor, i == 0) // inclusive on first
		if next.IsZero() {
			break
		}
		results = append(results, &RROccurrence{
			Time:     next,
			IsDeny:   rre.IsDeny,
			Priority: rre.Priority,
			Name:     rre.Name,
		})
		cursor = next.Add(time.Second)
	}

	return results, nil
}

func (rre *RRuleExtend) String() string {
	name := strings.TrimSpace(rre.Name)
	if name != "" {
		return name
	}
	parts := rre.ToDescriptor()
	if len(parts) == 0 {
		return ""
	}
	return strings.Join(parts, " · ")
}

func (rre *RRuleExtend) ToDescriptor() []string {
	var parts []string

	// 0: Prefer .Name if present
	name := strings.TrimSpace(rre.Name)
	if name != "" {
		parts = append(parts, name)
	} else {
		parts = append(parts, "Schedule Rule")
	}

	// Rule type
	if rre.IsDeny {
		parts = append(parts, "DENY rule")
	}

	// Time window
	if rre.IsAnyTime {
		parts = append(parts, "Any time (always valid)")
	} else {
		if rre.StartDate != nil && !rre.StartDate.IsZero() {
			parts = append(parts, fmt.Sprintf("Starts %s", rre.StartDate.Format("2006-01-02")))
		}
		if rre.EndDate != nil && !rre.EndDate.IsZero() {
			parts = append(parts, fmt.Sprintf("Ends %s", rre.EndDate.Format("2006-01-02")))
		}
		parts = append(parts, describeROptions(rre.ROptions)...)
	}

	// Duration
	if rre.Duration > 0 && !rre.DurationUnit.IsEmpty() {
		parts = append(parts, fmt.Sprintf("Duration: %d %s", rre.Duration, rre.DurationUnit.String()))
	}

	// Join Windows
	if len(rre.JoinWindows) > 0 {
		parts = append(parts, "Join window(s) enabled")
	}

	// Geo filters
	if rre.GeoFilters != nil && len(rre.GeoFilters) > 0 {
		parts = append(parts, "Geo filter(s) enabled")
	}

	return parts
}

type RRuleExtends []*RRuleExtend

// Validate validates each RRuleExtend in the slice.
// It returns the first encountered error, including the index where it occurred.
func (rres RRuleExtends) Validate() error {
	for i, rule := range rres {
		if rule == nil {
			return fmt.Errorf("rule at index %d is nil", i)
		}
		if err := rule.Validate(); err != nil {
			return fmt.Errorf("rule at index %d failed validation: %w", i, err)
		}
	}
	return nil
}

// sortByPriority returns a new slice of rules sorted by ascending Priority.
func (rres RRuleExtends) sorted() RRuleExtends {
	sorted := make(RRuleExtends, len(rres))
	copy(sorted, rres)
	sort.SliceStable(sorted, func(i, j int) bool {
		return sorted[i].Priority > sorted[j].Priority // Highest priority first
	})
	return sorted
}

// Evaluate evaluates a stack of RRuleExtend policies against a given time.
//
// Evaluation order:
// - Sorted by priority (descending)
// - If any deny rule matches, access is denied immediately
// - If any allow rule matches, access is granted
// - If no rules match, access is denied by default
func (rres RRuleExtends) Evaluate(now time.Time) (bool, error) {
	sorted := rres.sorted()

	for _, r := range sorted {
		match, err := r.matchesOccurrenceWindow(now)
		if err != nil {
			return false, err
		}
		if r.IsDeny && match {
			return false, nil // matched a deny rule
		}
		if !r.IsDeny && match {
			return true, nil // matched an allow rule
		}
	}
	return false, nil // default deny
}

// EvaluateWithOptions evaluates time-based access using optional geo or custom evaluator.
// Uses `IsBetweenWithOptions`, which applies extended logic (e.g. geo/IP filters).
func (rres RRuleExtends) EvaluateWithOptions(now time.Time, geo ageo.GeoInfo, eval IRRuleEvaluator) (bool, error) {
	sorted := rres.sorted()

	var bestAllow *RRuleExtend
	var bestDeny *RRuleExtend

	for _, r := range sorted {
		match, err := r.matchesOccurrenceWindowWithOptions(now, geo, eval)
		if err != nil {
			return false, err
		}
		if !match {
			continue
		}

		if r.IsDeny {
			if bestDeny == nil || r.Priority > bestDeny.Priority {
				bestDeny = r
			}
		} else {
			if bestAllow == nil || r.Priority > bestAllow.Priority {
				bestAllow = r
			}
		}
	}

	switch {
	case bestDeny != nil && (bestAllow == nil || bestDeny.Priority >= bestAllow.Priority):
		return false, nil
	case bestAllow != nil:
		return true, nil
	default:
		return false, nil
	}
}

// GetNextOccurrencesStacked returns a map of upcoming occurrences,
// keyed by the evaluation order (based on sorted priority).
// Each entry contains the upcoming RROccurrences from one rule.
func (rres RRuleExtends) GetNextOccurrencesStacked(now time.Time, count int) (RROccurrenceMap, error) {
	sorted := rres.sorted()
	result := make(RROccurrenceMap)

	for ii, rule := range sorted {
		occs, err := rule.GetNextOccurrences(now, count)
		if err != nil {
			return nil, err
		}
		if len(occs) == 0 {
			continue
		}
		result[ii] = append(result[ii], occs...)
	}

	return result, nil
}
