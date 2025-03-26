package alog

import (
	"github.com/jpfluger/alibs-slim/ajson"
	"github.com/jpfluger/alibs-slim/autils"
	"strings"
)

// WriterType defines a custom type for writer identifiers.
type WriterType string

// IsEmpty checks if the WriterType is empty after trimming spaces.
func (rt WriterType) IsEmpty() bool {
	return strings.TrimSpace(string(rt)) == ""
}

// TrimSpace trims spaces from the WriterType and returns a new WriterType.
func (rt WriterType) TrimSpace() WriterType {
	return WriterType(strings.TrimSpace(string(rt)))
}

// String converts the WriterType to a string.
func (rt WriterType) String() string {
	return string(rt)
}

// ToStringTrimLower trims spaces from the WriterType, converts it to lowercase, and returns it as a string.
func (rt WriterType) ToStringTrimLower() string {
	return autils.ToStringTrimLower(string(rt))
}

// ToJsonKey converts the WriterType to a JsonKey after trimming spaces and converting to lowercase.
func (rt WriterType) ToJsonKey() ajson.JsonKey {
	return ajson.JsonKey(rt.ToStringTrimLower())
}

// HasMatch checks if the WriterType matches the provided WriterType.
func (rt WriterType) HasMatch(rtType WriterType) bool {
	return rt == rtType
}

// MatchesOne checks if the WriterType matches any of the provided WriterTypes.
func (rt WriterType) MatchesOne(rtTypes ...WriterType) bool {
	for _, rtType := range rtTypes {
		if rt == rtType {
			return true
		}
	}
	return false
}

// HasPrefix checks if the WriterType has the specified prefix.
func (rt WriterType) HasPrefix(rtType WriterType) bool {
	return strings.HasPrefix(rt.String(), rtType.String())
}

// HasSuffix checks if the WriterType has the specified suffix.
func (rt WriterType) HasSuffix(rtType WriterType) bool {
	return strings.HasSuffix(rt.String(), rtType.String())
}

// GetLeaf extracts the last element in the path of the WriterType.
func (rt WriterType) GetLeaf() string {
	return rt.ToJsonKey().GetPathLeaf().String()
}

// GetLeafToUpper extracts the last element in the path of the WriterType and converts it to uppercase.
func (rt WriterType) GetLeafToUpper() string {
	return strings.ToUpper(rt.GetLeaf())
}

// ValidateAsLabel validates the WriterType as a label with specific rules.
func (rt WriterType) ValidateAsLabel() error {
	// Implementation is commented out; needs to be completed based on the rules.
	return nil
}

// WriterTypes defines a slice of WriterType.
type WriterTypes []WriterType

// HasValues checks if the WriterTypes slice has any values.
func (rts WriterTypes) HasValues() bool {
	return len(rts) > 0
}

// HasMatch checks if any WriterType in the slice matches the provided WriterType.
func (rts WriterTypes) HasMatch(rType WriterType) bool {
	for _, rt := range rts {
		if rt == rType {
			return true
		}
	}
	return false
}

// HasPrefix checks if any WriterType in the slice has the specified prefix.
func (rts WriterTypes) HasPrefix(rType WriterType) bool {
	for _, rt := range rts {
		if rt.HasPrefix(rType) {
			return true
		}
	}
	return false
}

// Clone creates a copy of the WriterTypes slice.
func (rts WriterTypes) Clone() WriterTypes {
	var arr WriterTypes
	for _, rt := range rts {
		if !rt.IsEmpty() {
			arr = append(arr, rt)
		}
	}
	return arr
}

// ToArrStrings converts the WriterTypes slice to a slice of strings.
func (rts WriterTypes) ToArrStrings() []string {
	var arr []string
	for _, rt := range rts {
		if !rt.IsEmpty() {
			arr = append(arr, rt.String())
		}
	}
	return arr
}

// IncludeIfInTargets includes WriterTypes that are present in the target slice.
func (rts WriterTypes) IncludeIfInTargets(targets WriterTypes) WriterTypes {
	var arr WriterTypes
	for _, rt := range rts {
		if targets.HasMatch(rt) {
			arr = append(arr, rt)
		}
	}
	return arr
}

// Clean removes empty WriterTypes from the slice.
func (rts WriterTypes) Clean() WriterTypes {
	var arr WriterTypes
	for _, rt := range rts {
		if !rt.IsEmpty() {
			arr = append(arr, rt)
		}
	}
	return arr
}
