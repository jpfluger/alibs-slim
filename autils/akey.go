package autils

import (
	"fmt"
	"strings"
	"unicode"
)

// AKey is a structured string that:
// 1. Allows alpha A-Z, a-z, numbers 0-9, special characters (-_.)
// 2. Only alphas and numbers can begin or end the string
// 3. The period is the divider
// Derive your own custom "keys" from AKey, if desired.
type AKey string

// IsEmpty checks if the AKey is empty after trimming space.
func (rt AKey) IsEmpty() bool {
	rtNew := strings.TrimSpace(string(rt))
	return rtNew == ""
}

// TrimSpace trims leading and trailing white spaces from AKey.
func (rt AKey) TrimSpace() AKey {
	rtNew := strings.TrimSpace(string(rt))
	return AKey(rtNew)
}

// String returns the string representation of AKey.
func (rt AKey) String() string {
	return string(rt)
}

// HasMatch checks if the AKey matches the provided AKey.
func (rt AKey) HasMatch(rtType AKey) bool {
	return rt == rtType
}

// MatchesOne checks if the AKey matches any one of the provided AKeys.
func (rt AKey) MatchesOne(rtTypes ...AKey) bool {
	for _, rtType := range rtTypes {
		if rt == rtType {
			return true
		}
	}
	return false
}

// HasPrefix checks if the AKey has the provided AKey as a prefix.
func (rt AKey) HasPrefix(rtType AKey) bool {
	return strings.HasPrefix(rt.String(), rtType.String())
}

// HasSuffix checks if the AKey has the provided AKey as a suffix.
func (rt AKey) HasSuffix(rtType AKey) bool {
	return strings.HasSuffix(rt.String(), rtType.String())
}

// Validate checks if the AKey is valid according to the rules.
func (rt AKey) Validate() error {
	target := rt.String()
	for ii, r := range target {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '-' && r != '_' && r != '.' {
			return fmt.Errorf("invalid char '%s'; AKey allows alpha A-Z, a-z, numbers 0-9, special '-_.'", string(r))
		}
		if (ii == 0 || ii == len(target)-1) && (r == '-' || r == '_' || r == '.') {
			return fmt.Errorf("invalid char '%s' at position %d; only alphas and numbers can begin or end the AKey", string(r), ii)
		}
	}
	return nil
}

// AKeys is a slice of AKey.
type AKeys []AKey

// HasValues checks if the AKeys slice has any values.
func (rts AKeys) HasValues() bool {
	return len(rts) > 0
}

// HasMatch checks if any AKey in the slice matches the provided AKey.
func (rts AKeys) HasMatch(rType AKey) bool {
	if len(rts) == 0 || rType.IsEmpty() {
		return false
	}
	for _, rt := range rts {
		if rt == rType {
			return true
		}
	}
	return false
}

// HasPrefix checks if any AKey in the slice has the provided AKey as a prefix.
func (rts AKeys) HasPrefix(rType AKey) bool {
	if len(rts) == 0 || rType.IsEmpty() {
		return false
	}
	for _, rt := range rts {
		if rt.HasPrefix(rType) {
			return true
		}
	}
	return false
}

// Clone creates a copy of the AKeys slice.
func (rts AKeys) Clone() AKeys {
	var arr AKeys
	for _, rt := range rts {
		if !rt.IsEmpty() {
			arr = append(arr, rt)
		}
	}
	return arr
}

// ToArrStrings converts the AKeys slice to a slice of strings.
func (rts AKeys) ToArrStrings() []string {
	var arr []string
	for _, rt := range rts {
		if !rt.IsEmpty() {
			arr = append(arr, rt.String())
		}
	}
	return arr
}

// IncludeIfInTargets includes AKeys that match any of the target AKeys.
func (rts AKeys) IncludeIfInTargets(targets AKeys) AKeys {
	var arr AKeys
	for _, rt := range rts {
		if targets.HasMatch(rt) {
			arr = append(arr, rt)
		}
	}
	return arr
}

// Clean removes empty AKeys from the slice.
func (rts AKeys) Clean() AKeys {
	var arr AKeys
	for _, rt := range rts {
		if !rt.IsEmpty() {
			arr = append(arr, rt)
		}
	}
	return arr
}
