package auuids

import "strings"

// I think we can simplify `EndpointLookupParams`. Let's require a UUID. I have a package auuids.UUID that can handle this. But there might be multiple UUIDs.

// UUIDLabel represents the type of an operator.
type UUIDLabel string

// IsEmpty checks if the UUIDLabel is empty after trimming whitespace.
func (lab UUIDLabel) IsEmpty() bool {
	return strings.TrimSpace(string(lab)) == ""
}

// HasMatch checks if the UUIDLabel matches the given operator.
func (lab UUIDLabel) HasMatch(operatorType UUIDLabel) bool {
	return lab == operatorType
}

// String is the UUIDLabel as a string
func (lab UUIDLabel) String() string {
	return string(lab)
}

// MatchesOne checks if the UUIDLabel matches any of the given operators.
func (lab UUIDLabel) MatchesOne(operatorTypes ...UUIDLabel) bool {
	for _, operatorType := range operatorTypes {
		if lab == operatorType {
			return true
		}
	}
	return false
}

type UUIDLabels []UUIDLabel

// IsValid checks if all given operators exist in the UUIDLabels slice.
func (labs UUIDLabels) IsValid(operatorTypes ...UUIDLabel) bool {
	for _, operatorType := range operatorTypes {
		var isFound bool
		for _, lab := range labs {
			if operatorType == lab {
				isFound = true
				break
			}
		}
		if !isFound {
			return false
		}
	}
	return true
}
