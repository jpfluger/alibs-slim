package acron

import (
	"strings" // Importing the strings package for string manipulation.
)

// TaskType defines a custom type for job data identifiers.
type TaskType string

// IsEmpty checks if the TaskType is empty after trimming spaces.
func (nt TaskType) IsEmpty() bool {
	// Trim spaces from the TaskType and check if the result is an empty string.
	return strings.TrimSpace(string(nt)) == ""
}

// TrimSpace trims spaces from the TaskType and returns a new TaskType.
func (nt TaskType) TrimSpace() TaskType {
	// Trim spaces from the TaskType and return the result as a new TaskType.
	return TaskType(strings.TrimSpace(string(nt)))
}

// String converts TaskType to a string.
func (nt TaskType) String() string {
	// Convert the TaskType to a string and return it.
	return string(nt)
}

// ToStringTrimLower converts TaskType to a string, trims spaces, and makes it lowercase.
func (nt TaskType) ToStringTrimLower() string {
	// Convert the TaskType to a string, trim spaces, convert to lowercase, and return the result.
	return strings.ToLower(nt.TrimSpace().String())
}

// TaskTypes defines a slice of TaskType.
type TaskTypes []TaskType

// Contains checks if the TaskTypes slice contains a specific TaskType.
func (nts TaskTypes) Contains(nt TaskType) bool {
	// Iterate over the TaskTypes slice.
	for _, t := range nts {
		// Check if the current TaskType matches the specified TaskType.
		if t == nt {
			return true // Return true if a match is found.
		}
	}
	return false // Return false if no match is found.
}

// Add appends a new TaskType to the TaskTypes slice if it's not already present.
func (nts *TaskTypes) Add(nt TaskType) {
	// Check if the TaskType is not already in the slice.
	if !nts.Contains(nt) {
		*nts = append(*nts, nt) // Append the new TaskType to the slice.
	}
}

// Remove deletes a TaskType from the TaskTypes slice.
func (nts *TaskTypes) Remove(nt TaskType) {
	// Create a new slice to store the result.
	var result TaskTypes
	// Iterate over the TaskTypes slice.
	for _, t := range *nts {
		// If the current TaskType does not match the specified TaskType, add it to the result slice.
		if t != nt {
			result = append(result, t)
		}
	}
	*nts = result // Set the TaskTypes slice to the result.
}
