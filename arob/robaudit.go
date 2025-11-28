package arob

import (
	"fmt"
	"strings"
)

// ROBAudit tracks a section's validation status with results.
type ROBAudit struct {
	Label      ROBAuditLabel  `json:"label"`      // e.g., "paths"
	HasPassed  bool           `json:"hasPassed"`  // True if no errors
	TriedCheck bool           `json:"triedCheck"` // Attempted?
	Points     ROBAuditPoints `json:"points"`     // Slice of results (warnings/errors/infos)
}

func NewROBAudit(label ROBAuditLabel) *ROBAudit {
	return &ROBAudit{Label: label}
}

func NewROBAuditTriedTrue(label ROBAuditLabel) *ROBAudit {
	return &ROBAudit{Label: label, TriedCheck: true}
}

// WithPoint adds a point to the audit and returns the audit for chaining.
// It also updates HasPassed if the added point is an error.
func (c *ROBAudit) WithPoint(point *ROBAuditPoint) *ROBAudit {
	c.Points = append(c.Points, point)
	if point.IsErr() {
		c.HasPassed = false
	}
	return c
}

// HasErrors returns true if any result is error-level.
func (c *ROBAudit) HasErrors() bool {
	for _, r := range c.Points {
		if r.IsErr() {
			return true
		}
	}
	return false
}

// FilterErrors returns a slice containing only the error-level points from the audit.
func (c *ROBAudit) FilterErrors() ROBAuditPoints {
	arr := ROBAuditPoints{}
	for _, r := range c.Points {
		if r.IsErr() {
			arr = append(arr, r)
		}
	}
	return arr
}

// ErrorSummary returns a formatted string summarizing the errors in this audit.
func (c *ROBAudit) ErrorSummary() string {
	if !c.HasErrors() {
		return ""
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Errors in section '%s':\n", c.Label))
	errors := c.FilterErrors()
	for i, errPoint := range errors {
		sb.WriteString(fmt.Sprintf("  %d. %s (Type: %s, Field: %s)\n", i+1, errPoint.Message, errPoint.Type, errPoint.Field))
	}
	return sb.String()
}

// ROBAudits is a slice of ROBAudit for full validation output.
type ROBAudits []ROBAudit

// FilterByErrors returns only the audits that contain errors.
func (bas ROBAudits) FilterByErrors() ROBAudits {
	arr := ROBAudits{}
	for _, ba := range bas {
		if ba.HasErrors() {
			arr = append(arr, ba)
		}
	}
	return arr
}

// HasErrors returns true if any audit in the slice contains errors.
func (bas ROBAudits) HasErrors() bool {
	audits := bas.FilterByErrors()
	return len(audits) > 0
}

// ErrorSummary returns a formatted string summarizing all errors across audits.
func (bas ROBAudits) ErrorSummary() string {
	errAudits := bas.FilterByErrors()
	if len(errAudits) == 0 {
		return "No errors found."
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Found errors in %d sections:\n\n", len(errAudits)))
	for _, audit := range errAudits {
		sb.WriteString(audit.ErrorSummary() + "\n")
	}
	return sb.String()
}

// AllErrors returns a flat list of all error points across audits, with section labels prepended.
func (bas ROBAudits) AllErrors() []string {
	var errs []string
	errAudits := bas.FilterByErrors()
	for _, audit := range errAudits {
		errors := audit.FilterErrors()
		for _, errPoint := range errors {
			errs = append(errs, fmt.Sprintf("[%s] %s (Type: %s, Field: %s)", audit.Label, errPoint.Message, errPoint.Type, errPoint.Field))
		}
	}
	return errs
}

// ROBAuditPoint holds a single validation result with severity (via ROBType).
type ROBAuditPoint struct {
	Type    ROBType `json:"type,omitempty" xml:"type,omitempty"`       // Severity (e.g., error, warning)
	Message string  `json:"message,omitempty" xml:"message,omitempty"` // Details
	Field   string  `json:"field,omitempty" xml:"field,omitempty"`     // Optional field tie-in (e.g., "DIR_ROOT")
}

// NewROBAuditPoint creates a result with normalized type.
func NewROBAuditPoint(robType ROBType, message string) *ROBAuditPoint {
	return &ROBAuditPoint{
		Type:    NormalizeROBType(robType),
		Message: strings.TrimSpace(message),
	}
}

// NewROBAuditPointf creates a result with normalized type and type format specifier.
func NewROBAuditPointf(robType ROBType, format string, a ...any) *ROBAuditPoint {
	return &ROBAuditPoint{
		Type:    NormalizeROBType(robType),
		Message: fmt.Sprintf(format, a...),
	}
}

// NewROBAuditPointWithField adds a field association.
func NewROBAuditPointWithField(robType ROBType, field string, message string) *ROBAuditPoint {
	point := NewROBAuditPoint(robType, message)
	point.Field = strings.TrimSpace(field)
	return point
}

// NewROBAuditPointWithFieldf adds a field association.
func NewROBAuditPointWithFieldf(robType ROBType, field string, format string, a ...any) *ROBAuditPoint {
	point := NewROBAuditPoint(robType, fmt.Sprintf(format, a...))
	point.Field = strings.TrimSpace(field)
	return point
}

// IsErr returns true if this is an error-level point.
func (bap *ROBAuditPoint) IsErr() bool {
	return bap.Type.MatchesOne(ROBTYPE_ERROR, ROBTYPE_CRITICAL, ROBTYPE_EMERGENCY)
}

// ROBAuditPoints is a slice of ROBAuditPoint for ordered, multi-result feedback.
type ROBAuditPoints []*ROBAuditPoint

// ToStringArray filters messages by ROBType.
func (baps ROBAuditPoints) ToStringArray(robType ROBType) []string {
	var arr []string
	for _, bap := range baps {
		if bap.Type == robType {
			arr = append(arr, bap.Message)
		}
	}
	return arr
}

type ROBAuditLabel string

// String implements fmt.Stringer.
func (rl ROBAuditLabel) String() string { return string(rl) }

// IsEmpty trims whitespace first, then checks for emptiness.
func (rl ROBAuditLabel) IsEmpty() bool { return strings.TrimSpace(string(rl)) == "" }
