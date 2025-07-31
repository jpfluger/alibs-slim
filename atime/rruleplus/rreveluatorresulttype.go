package rruleplus

import "strings"

// RREvaluatorResultType defines the possible outcomes of an evaluator's decision logic.
// It is used by evaluators like IsPreAllowed to indicate how rule evaluation should proceed.
type RREvaluatorResultType string

// Enumeration of valid evaluator result types:
//
// - RREVALUATOR_RESULTTYPE_CONTINUE: Continue with standard rule evaluation.
// - RREVALUATOR_RESULTTYPE_ALLOW:    Immediately allow access, short-circuiting rule logic.
// - RREVALUATOR_RESULTTYPE_DENY:     Immediately deny access, short-circuiting rule logic.
const (
	RREVALUATOR_RESULTTYPE_CONTINUE RREvaluatorResultType = "continue"
	RREVALUATOR_RESULTTYPE_ALLOW    RREvaluatorResultType = "allow"
	RREVALUATOR_RESULTTYPE_DENY     RREvaluatorResultType = "deny"
)

// IsEmpty returns true if the result type is not set (i.e., an empty string).
func (t RREvaluatorResultType) IsEmpty() bool {
	return string(t) == ""
}

// String returns the lowercase string representation of the result type.
// Useful for logging or comparisons.
func (t RREvaluatorResultType) String() string {
	return strings.ToLower(string(t))
}
