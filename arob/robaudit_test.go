package arob

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestROBAudit_HasErrors(t *testing.T) {
	tests := []struct {
		name     string
		points   ROBAuditPoints
		expected bool
	}{
		{
			name:     "No points - no errors",
			points:   ROBAuditPoints{},
			expected: false,
		},
		{
			name: "Only non-error points",
			points: ROBAuditPoints{
				NewROBAuditPoint(ROBTYPE_INFO, "info message"),
				NewROBAuditPoint(ROBTYPE_WARNING, "warning message"),
			},
			expected: false,
		},
		{
			name: "Contains error point",
			points: ROBAuditPoints{
				NewROBAuditPoint(ROBTYPE_INFO, "info"),
				NewROBAuditPoint(ROBTYPE_ERROR, "error occurred"),
			},
			expected: true,
		},
		{
			name: "Contains critical point",
			points: ROBAuditPoints{
				NewROBAuditPoint(ROBTYPE_CRITICAL, "system down"),
			},
			expected: true,
		},
		{
			name: "Contains emergency point",
			points: ROBAuditPoints{
				NewROBAuditPoint(ROBTYPE_EMERGENCY, "catastrophic failure"),
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			audit := &ROBAudit{
				Points: tt.points,
			}
			assert.Equal(t, tt.expected, audit.HasErrors())
		})
	}
}

func TestNewROBAuditPoint(t *testing.T) {
	tests := []struct {
		name       string
		inputType  ROBType
		inputMsg   string
		expectType ROBType
		expectMsg  string
	}{
		{
			name:       "Standard type - no normalization",
			inputType:  ROBTYPE_INFO,
			inputMsg:   "  system ready  ",
			expectType: ROBTYPE_INFO,
			expectMsg:  "system ready",
		},
		{
			name:       "Shorthand 'err' normalizes to error",
			inputType:  "err",
			inputMsg:   "invalid input",
			expectType: ROBTYPE_ERROR,
			expectMsg:  "invalid input",
		},
		{
			name:       "Empty type defaults to debug",
			inputType:  "",
			inputMsg:   "fallback case",
			expectType: ROBTYPE_DEBUG,
			expectMsg:  "fallback case",
		},
		{
			name:       "Whitespace trimmed from message",
			inputType:  ROBTYPE_WARNING,
			inputMsg:   "  deprecated API  ",
			expectType: ROBTYPE_WARNING,
			expectMsg:  "deprecated API",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			point := NewROBAuditPoint(tt.inputType, tt.inputMsg)
			assert.Equal(t, tt.expectType, point.Type)
			assert.Equal(t, tt.expectMsg, point.Message)
			assert.Empty(t, point.Field)
		})
	}
}

func TestNewROBAuditPointWithField(t *testing.T) {
	point := NewROBAuditPointWithField("crit", "DB_CONN", "connection timeout")

	assert.Equal(t, ROBTYPE_CRITICAL, point.Type)
	assert.Equal(t, "connection timeout", point.Message)
	assert.Equal(t, "DB_CONN", point.Field)
}

func TestNewROBAuditPointWithFieldf(t *testing.T) {
	point := NewROBAuditPointWithFieldf(ROBTYPE_ERROR, "USER_ID", "user %d not found", 404)

	assert.Equal(t, ROBTYPE_ERROR, point.Type)
	assert.Equal(t, "user 404 not found", point.Message)
	assert.Equal(t, "USER_ID", point.Field)
}

func TestROBAuditPoint_IsErr(t *testing.T) {
	tests := []struct {
		name     string
		point    *ROBAuditPoint
		expected bool
	}{
		{
			name:     "Error type returns true",
			point:    NewROBAuditPoint(ROBTYPE_ERROR, "fail"),
			expected: true,
		},
		{
			name:     "Critical type returns true",
			point:    NewROBAuditPoint(ROBTYPE_CRITICAL, "crash"),
			expected: true,
		},
		{
			name:     "Emergency type returns true",
			point:    NewROBAuditPoint(ROBTYPE_EMERGENCY, "panic"),
			expected: true,
		},
		{
			name:     "Warning type returns false",
			point:    NewROBAuditPoint(ROBTYPE_WARNING, "slow"),
			expected: false,
		},
		{
			name:     "Info type returns false",
			point:    NewROBAuditPoint(ROBTYPE_INFO, "ok"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.point.IsErr())
		})
	}
}

func TestROBAuditPoints_ToStringArray(t *testing.T) {
	points := ROBAuditPoints{
		NewROBAuditPoint(ROBTYPE_INFO, "starting"),
		NewROBAuditPoint(ROBTYPE_ERROR, "failed to connect"),
		NewROBAuditPoint(ROBTYPE_ERROR, "invalid credentials"),
		NewROBAuditPoint(ROBTYPE_WARNING, "high latency"),
	}

	t.Run("Filter errors only", func(t *testing.T) {
		result := points.ToStringArray(ROBTYPE_ERROR)
		expected := []string{"failed to connect", "invalid credentials"}
		assert.Equal(t, expected, result)
	})

	t.Run("Filter warnings only", func(t *testing.T) {
		result := points.ToStringArray(ROBTYPE_WARNING)
		expected := []string{"high latency"}
		assert.Equal(t, expected, result)
	})

	t.Run("No matches returns empty", func(t *testing.T) {
		result := points.ToStringArray(ROBTYPE_DEBUG)
		assert.Empty(t, result)
	})
}

func TestROBAuditLabel_String(t *testing.T) {
	label := ROBAuditLabel("config")
	assert.Equal(t, "config", label.String())
}

func TestROBAuditLabel_IsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		label    ROBAuditLabel
		expected bool
	}{
		{
			name:     "Empty string",
			label:    "",
			expected: true,
		},
		{
			name:     "Whitespace only",
			label:    "   ",
			expected: true,
		},
		{
			name:     "Non-empty",
			label:    "paths",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.label.IsEmpty())
		})
	}
}

func TestNewROBAudit(t *testing.T) {
	label := ROBAuditLabel("schema")
	audit := NewROBAudit(label)

	assert.Equal(t, label, audit.Label)
	assert.False(t, audit.HasPassed)
	assert.False(t, audit.TriedCheck)
	assert.Empty(t, audit.Points)
}

func TestROBAudit_FilterErrors(t *testing.T) {
	tests := []struct {
		name     string
		points   ROBAuditPoints
		expected ROBAuditPoints
	}{
		{
			name:     "No errors",
			points:   ROBAuditPoints{NewROBAuditPoint(ROBTYPE_INFO, "ok"), NewROBAuditPoint(ROBTYPE_WARNING, "slow")},
			expected: ROBAuditPoints{},
		},
		{
			name:     "Mixed points with errors",
			points:   ROBAuditPoints{NewROBAuditPoint(ROBTYPE_ERROR, "fail"), NewROBAuditPoint(ROBTYPE_INFO, "info"), NewROBAuditPoint(ROBTYPE_CRITICAL, "crash")},
			expected: ROBAuditPoints{NewROBAuditPoint(ROBTYPE_ERROR, "fail"), NewROBAuditPoint(ROBTYPE_CRITICAL, "crash")},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			audit := &ROBAudit{Points: tt.points}
			result := audit.FilterErrors()
			assert.Len(t, result, len(tt.expected))
			for i, exp := range tt.expected {
				assert.Equal(t, exp.Type, result[i].Type)
				assert.Equal(t, exp.Message, result[i].Message)
			}
		})
	}
}

func TestROBAudit_ErrorSummary(t *testing.T) {
	tests := []struct {
		name     string
		audit    *ROBAudit
		expected string
	}{
		{
			name:     "No errors returns empty",
			audit:    &ROBAudit{Label: "test", Points: ROBAuditPoints{NewROBAuditPoint(ROBTYPE_INFO, "ok")}},
			expected: "",
		},
		{
			name: "With errors and fields",
			audit: &ROBAudit{
				Label: "config",
				Points: ROBAuditPoints{
					NewROBAuditPointWithField(ROBTYPE_ERROR, "FIELD1", "invalid value"),
					NewROBAuditPoint(ROBTYPE_CRITICAL, "system failure"),
				},
			},
			expected: "Errors in section 'config':\n  1. invalid value (Type: error, Field: FIELD1)\n  2. system failure (Type: critical, Field: )\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.audit.ErrorSummary()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestROBAudits_FilterByErrors(t *testing.T) {
	audits := ROBAudits{
		{Label: "noerr", Points: ROBAuditPoints{NewROBAuditPoint(ROBTYPE_INFO, "ok")}},
		{Label: "haserr", Points: ROBAuditPoints{NewROBAuditPoint(ROBTYPE_ERROR, "fail")}},
		{Label: "another", Points: ROBAuditPoints{NewROBAuditPoint(ROBTYPE_CRITICAL, "crash")}},
	}

	result := audits.FilterByErrors()
	assert.Len(t, result, 2)
	assert.Equal(t, "haserr", result[0].Label.String())
	assert.Equal(t, "another", result[1].Label.String())
}

func TestROBAudits_HasErrors(t *testing.T) {
	tests := []struct {
		name     string
		audits   ROBAudits
		expected bool
	}{
		{
			name:     "No audits - false",
			audits:   ROBAudits{},
			expected: false,
		},
		{
			name: "No errors - false",
			audits: ROBAudits{
				{Points: ROBAuditPoints{NewROBAuditPoint(ROBTYPE_INFO, "ok")}},
			},
			expected: false,
		},
		{
			name: "Has errors - true",
			audits: ROBAudits{
				{Points: ROBAuditPoints{NewROBAuditPoint(ROBTYPE_ERROR, "fail")}},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.audits.HasErrors())
		})
	}
}

func TestROBAudits_ErrorSummary(t *testing.T) {
	tests := []struct {
		name     string
		audits   ROBAudits
		expected string
	}{
		{
			name:     "No errors",
			audits:   ROBAudits{},
			expected: "No errors found.",
		},
		{
			name: "With errors in multiple sections",
			audits: ROBAudits{
				{
					Label: "section1",
					Points: ROBAuditPoints{
						NewROBAuditPoint(ROBTYPE_ERROR, "error1"),
						NewROBAuditPoint(ROBTYPE_INFO, "info"),
					},
				},
				{
					Label: "section2",
					Points: ROBAuditPoints{
						NewROBAuditPointWithField(ROBTYPE_CRITICAL, "FIELD", "critical error"),
					},
				},
			},
			expected: "Found errors in 2 sections:\n\nErrors in section 'section1':\n  1. error1 (Type: error, Field: )\n\nErrors in section 'section2':\n  1. critical error (Type: critical, Field: FIELD)\n\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.audits.ErrorSummary()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestROBAudits_AllErrors(t *testing.T) {
	audits := ROBAudits{
		{
			Label: "section1",
			Points: ROBAuditPoints{
				NewROBAuditPoint(ROBTYPE_ERROR, "error1"),
				NewROBAuditPoint(ROBTYPE_INFO, "info"),
				NewROBAuditPointWithField(ROBTYPE_CRITICAL, "FIELD1", "critical1"),
			},
		},
		{
			Label: "section2",
			Points: ROBAuditPoints{
				NewROBAuditPoint(ROBTYPE_WARNING, "warn"),
				NewROBAuditPoint(ROBTYPE_ERROR, "error2"),
			},
		},
	}

	result := audits.AllErrors()
	expected := []string{
		"[section1] error1 (Type: error, Field: )",
		"[section1] critical1 (Type: critical, Field: FIELD1)",
		"[section2] error2 (Type: error, Field: )",
	}
	assert.Equal(t, expected, result)
}
