package acron

import (
	"reflect"
	"testing"
)

// TestTaskType_IsEmpty checks the IsEmpty method for various TaskType values.
func TestTaskType_IsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		jobType  TaskType
		expected bool
	}{
		{"Empty", "", true},
		{"Whitespace", "   ", true},
		{"NotEmpty", "job-data", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.jobType.IsEmpty(); got != test.expected {
				t.Errorf("TaskType.IsEmpty() = %v, want %v", got, test.expected)
			}
		})
	}
}

// TestTaskType_TrimSpace checks the TrimSpace method for various TaskType values.
func TestTaskType_TrimSpace(t *testing.T) {
	tests := []struct {
		name     string
		jobType  TaskType
		expected TaskType
	}{
		{"LeadingSpace", " job-data", "job-data"},
		{"TrailingSpace", "job-data ", "job-data"},
		{"BothSpaces", " job-data ", "job-data"},
		{"NoSpaces", "job-data", "job-data"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.jobType.TrimSpace(); got != test.expected {
				t.Errorf("TaskType.TrimSpace() = %v, want %v", got, test.expected)
			}
		})
	}
}

// TestTaskType_ToStringTrimLower checks the ToStringTrimLower method for various TaskType values.
func TestTaskType_ToStringTrimLower(t *testing.T) {
	tests := []struct {
		name     string
		jobType  TaskType
		expected string
	}{
		{"AllCaps", "JOB-DATA", "job-data"},
		{"MixedCase", "Job-Data", "job-data"},
		{"LowerCase", "job-data", "job-data"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.jobType.ToStringTrimLower(); got != test.expected {
				t.Errorf("TaskType.ToStringTrimLower() = %v, want %v", got, test.expected)
			}
		})
	}
}

// TestTaskType_TaskTypesContains checks the Contains method for the TaskTypes slice.
func TestTaskTypesContains(t *testing.T) {
	jobTypes := TaskTypes{"job-data", "script-python", "task-xml"}

	tests := []struct {
		name     string
		jobType  TaskType
		expected bool
	}{
		{"Contains", "job-data", true},
		{"DoesNotContain", "unknown-type", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := jobTypes.Contains(test.jobType); got != test.expected {
				t.Errorf("TaskTypes.Contains() = %v, want %v", got, test.expected)
			}
		})
	}
}

// TestTaskTypesAdd checks the Add method for the TaskTypes slice.
func TestTaskTypesAdd(t *testing.T) {
	tests := []struct {
		name     string
		initial  TaskTypes
		jobType  TaskType
		expected TaskTypes
	}{
		{"AddNew", TaskTypes{"job-data", "script-python"}, "task-xml", TaskTypes{"job-data", "script-python", "task-xml"}},
		{"AddExisting", TaskTypes{"job-data", "script-python"}, "job-data", TaskTypes{"job-data", "script-python"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Create a new slice for each test to avoid modifying the same slice.
			jobTypes := make(TaskTypes, len(test.initial))
			copy(jobTypes, test.initial)

			jobTypes.Add(test.jobType)
			if !reflect.DeepEqual(jobTypes, test.expected) {
				t.Errorf("TaskTypes.Add() = %v, want %v", jobTypes, test.expected)
			}
		})
	}
}

// TestTaskTypesRemove checks the Remove method for the TaskTypes slice.
func TestTaskTypesRemove(t *testing.T) {
	tests := []struct {
		name     string
		initial  TaskTypes
		jobType  TaskType
		expected TaskTypes
	}{
		{"RemoveExisting", TaskTypes{"job-data", "script-python", "task-xml"}, "script-python", TaskTypes{"job-data", "task-xml"}},
		{"RemoveNonExisting", TaskTypes{"job-data", "script-python", "task-xml"}, "unknown-type", TaskTypes{"job-data", "script-python", "task-xml"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Create a new slice for each test to avoid modifying the same slice.
			jobTypes := make(TaskTypes, len(test.initial))
			copy(jobTypes, test.initial)

			jobTypes.Remove(test.jobType)
			if !reflect.DeepEqual(jobTypes, test.expected) {
				t.Errorf("TaskTypes.Remove() = %v, want %v", jobTypes, test.expected)
			}
		})
	}
}
