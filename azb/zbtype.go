package azb

import (
	// Importing necessary packages for JSON and utility functions.
	"github.com/jpfluger/alibs-slim/ajson"
	"github.com/jpfluger/alibs-slim/autils"
	"strings"
)

// ZBType defines a custom type based on string to represent a specific type.
type ZBType string

// IsEmpty checks if the ZBType is empty after trimming space.
func (rt ZBType) IsEmpty() bool {
	rtNew := strings.TrimSpace(string(rt))
	return rtNew == ""
}

// TrimSpace trims the spaces from ZBType and returns a new ZBType.
func (rt ZBType) TrimSpace() ZBType {
	rtNew := strings.TrimSpace(string(rt))
	return ZBType(rtNew)
}

// String converts ZBType to a string.
func (rt ZBType) String() string {
	return string(rt)
}

// ToStringTrimLower trims the ZBType, converts it to lower case, and returns a string.
func (rt ZBType) ToStringTrimLower() string {
	return autils.ToStringTrimLower(string(rt))
}

// ToJsonKey converts ZBType to a JsonKey after trimming and lowering the string.
func (rt ZBType) ToJsonKey() ajson.JsonKey {
	return ajson.JsonKey(rt.ToStringTrimLower())
}

// HasMatch checks if the ZBType matches another ZBType.
func (rt ZBType) HasMatch(rtType ZBType) bool {
	return rt == rtType
}

// MatchesOne checks if the ZBType matches any one of the provided ZBTypes.
func (rt ZBType) MatchesOne(rtTypes ...ZBType) bool {
	for _, rtType := range rtTypes {
		if rt == rtType {
			return true
		}
	}
	return false
}

// HasPrefix checks if the ZBType has a specific prefix.
func (rt ZBType) HasPrefix(rtType ZBType) bool {
	return strings.HasPrefix(rt.String(), rtType.String())
}

// HasSuffix checks if the ZBType has a specific suffix.
func (rt ZBType) HasSuffix(rtType ZBType) bool {
	return strings.HasSuffix(rt.String(), rtType.String())
}

// GetLeaf extracts the last element in a JSON key path.
func (rt ZBType) GetLeaf() string {
	return rt.ToJsonKey().GetPathLeaf().String()
}

// GetLeafToUpper converts the last element in a JSON key path to upper case.
func (rt ZBType) GetLeafToUpper() string {
	return strings.ToUpper(rt.ToJsonKey().GetPathLeaf().String())
}

// ZBTypes defines a slice of ZBType.
type ZBTypes []ZBType

// HasValues checks if the slice of ZBType has any values.
func (rts ZBTypes) HasValues() bool {
	return rts != nil && len(rts) > 0
}

// HasMatch checks if any ZBType in the slice matches the provided ZBType.
func (rts ZBTypes) HasMatch(rType ZBType) bool {
	if rts == nil || len(rts) == 0 || rType.IsEmpty() {
		return false
	}
	for _, rt := range rts {
		if rt == rType {
			return true
		}
	}
	return false
}

// HasPrefix checks if any ZBType in the slice has the specified prefix.
func (rts ZBTypes) HasPrefix(rType ZBType) bool {
	if rts == nil || len(rts) == 0 || rType.IsEmpty() {
		return false
	}
	for _, rt := range rts {
		if rt.HasPrefix(rType) {
			return true
		}
	}
	return false
}

// Clone creates a copy of the ZBTypes slice with non-empty values.
func (rts ZBTypes) Clone() ZBTypes {
	arr := ZBTypes{}
	if rts == nil || len(rts) == 0 {
		return arr
	}
	for _, rt := range rts {
		if !rt.IsEmpty() {
			arr = append(arr, rt)
		}
	}
	return arr
}

// ToArrStrings converts the ZBTypes slice to a slice of strings.
func (rts ZBTypes) ToArrStrings() []string {
	arr := []string{}
	if rts == nil || len(rts) == 0 {
		return arr
	}
	for _, rt := range rts {
		if !rt.IsEmpty() {
			arr = append(arr, rt.String())
		}
	}
	return arr
}

// IncludeIfInTargets includes ZBTypes that match any of the target ZBTypes.
func (rts ZBTypes) IncludeIfInTargets(targets ...ZBType) ZBTypes {
	arr := ZBTypes{}
	// No need to check if targets is nil, as variadic params are always initialized.
	if rts == nil || len(rts) == 0 || len(targets) == 0 {
		return arr
	}
	for _, rt := range rts {
		for _, target := range targets {
			if rt == target {
				arr = append(arr, rt)
				break // Break the inner loop to avoid duplicates.
			}
		}
	}
	return arr
}

// Clean returns a new ZBTypes slice with empty values removed.
func (rts ZBTypes) Clean() ZBTypes {
	arr := ZBTypes{}
	if rts == nil || len(rts) == 0 {
		return arr
	}
	for _, rt := range rts {
		if rt.IsEmpty() {
			continue
		}
		arr = append(arr, rt)
	}
	return arr
}
