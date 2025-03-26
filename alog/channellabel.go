package alog

import (
	"github.com/jpfluger/alibs-slim/ajson"
	"github.com/jpfluger/alibs-slim/autils"
	"strings"
)

// ChannelLabel represents a label for a logging channel.
type ChannelLabel string

// IsEmpty checks if the ChannelLabel is empty after trimming whitespace.
func (cl ChannelLabel) IsEmpty() bool {
	return strings.TrimSpace(string(cl)) == ""
}

// TrimSpace trims whitespace from the ChannelLabel and returns a new ChannelLabel.
func (cl ChannelLabel) TrimSpace() ChannelLabel {
	return ChannelLabel(strings.TrimSpace(string(cl)))
}

// String converts the ChannelLabel to a string.
func (cl ChannelLabel) String() string {
	return string(cl)
}

// ToStringTrimLower trims whitespace and converts the ChannelLabel to lowercase.
func (cl ChannelLabel) ToStringTrimLower() string {
	return autils.ToStringTrimLower(string(cl))
}

// ToJsonKey converts the ChannelLabel to a JsonKey after trimming and lowering case.
func (cl ChannelLabel) ToJsonKey() ajson.JsonKey {
	return ajson.JsonKey(cl.ToStringTrimLower())
}

// HasMatch checks if the ChannelLabel matches the provided ChannelLabel.
func (cl ChannelLabel) HasMatch(clType ChannelLabel) bool {
	return cl == clType
}

// MatchesOne checks if the ChannelLabel matches any of the provided ChannelLabels.
func (cl ChannelLabel) MatchesOne(clTypes ...ChannelLabel) bool {
	for _, clType := range clTypes {
		if cl == clType {
			return true
		}
	}
	return false
}

// HasPrefix checks if the ChannelLabel has the specified prefix.
func (cl ChannelLabel) HasPrefix(clType ChannelLabel) bool {
	return strings.HasPrefix(cl.String(), clType.String())
}

// HasSuffix checks if the ChannelLabel has the specified suffix.
func (cl ChannelLabel) HasSuffix(clType ChannelLabel) bool {
	return strings.HasSuffix(cl.String(), clType.String())
}

// GetLeaf extracts the last element in the path of the ChannelLabel.
func (cl ChannelLabel) GetLeaf() string {
	return cl.ToJsonKey().GetPathLeaf().String()
}

// GetLeafToUpper extracts the last element in the path of the ChannelLabel and converts it to uppercase.
func (cl ChannelLabel) GetLeafToUpper() string {
	return strings.ToUpper(cl.GetLeaf())
}

// ValidateAsLabel validates the ChannelLabel against specific rules.
func (cl ChannelLabel) ValidateAsLabel() error {
	// Implementation should validate the ChannelLabel according to the rules.
	// Currently, this method is not implemented.
	return nil
}

// ChannelLabels represents a slice of ChannelLabel.
type ChannelLabels []ChannelLabel

// HasValues checks if the ChannelLabels slice has any values.
func (cls ChannelLabels) HasValues() bool {
	return len(cls) > 0
}

// HasMatch checks if any ChannelLabel in the slice matches the provided ChannelLabel.
func (cls ChannelLabels) HasMatch(clType ChannelLabel) bool {
	for _, cl := range cls {
		if cl == clType {
			return true
		}
	}
	return false
}

// HasPrefix checks if any ChannelLabel in the slice has the specified prefix.
func (cls ChannelLabels) HasPrefix(clType ChannelLabel) bool {
	for _, cl := range cls {
		if cl.HasPrefix(clType) {
			return true
		}
	}
	return false
}

// Clone creates a copy of the ChannelLabels slice.
func (cls ChannelLabels) Clone() ChannelLabels {
	var arr ChannelLabels
	for _, cl := range cls {
		if !cl.IsEmpty() {
			arr = append(arr, cl)
		}
	}
	return arr
}

// ToArrStrings converts the ChannelLabels slice to a slice of strings.
func (cls ChannelLabels) ToArrStrings() []string {
	var arr []string
	for _, cl := range cls {
		if !cl.IsEmpty() {
			arr = append(arr, cl.String())
		}
	}
	return arr
}

// IncludeIfInTargets includes ChannelLabels that are present in the target slice.
func (cls ChannelLabels) IncludeIfInTargets(targets ChannelLabels) ChannelLabels {
	var arr ChannelLabels
	for _, cl := range cls {
		if targets.HasMatch(cl) {
			arr = append(arr, cl)
		}
	}
	return arr
}

// Clean removes empty ChannelLabels from the slice.
func (cls ChannelLabels) Clean() ChannelLabels {
	var arr ChannelLabels
	for _, cl := range cls {
		if !cl.IsEmpty() {
			arr = append(arr, cl)
		}
	}
	return arr
}
