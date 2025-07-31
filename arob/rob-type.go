package arob

import (
	"fmt"
	"github.com/jpfluger/alibs-slim/ajson"
	"strings"
	"unicode"
)

// ROBType represents the type of a ROB (Return/Request/Response Object) message.
type ROBType string

// Predefined constants for various ROB message types.
const (
	ROBTYPE_DEBUG     ROBType = "debug"
	ROBTYPE_INFO      ROBType = "info"
	ROBTYPE_NOTICE    ROBType = "notice"
	ROBTYPE_WARNING   ROBType = "warning"
	ROBTYPE_ERROR     ROBType = "error"
	ROBTYPE_CRITICAL  ROBType = "critical"
	ROBTYPE_EMERGENCY ROBType = "emergency"
)

// NormalizeROBType maps shorthand or non-standard ROBType values to predefined constants.
// This ensures consistent internal usage of ROB message types across the system.
//
// For example:
//
//	"emerg" → ROBTYPE_EMERGENCY
//	"crit"  → ROBTYPE_CRITICAL
//	"err"   → ROBTYPE_ERROR
//	"warn"  → ROBTYPE_WARNING
//	""      → ROBTYPE_DEBUG (default fallback)
//
// Values that already match predefined constants are returned as-is.
func NormalizeROBType(robType ROBType) ROBType {
	switch robType {
	case "emerg":
		return ROBTYPE_EMERGENCY
	case "crit":
		return ROBTYPE_CRITICAL
	case "err":
		return ROBTYPE_ERROR
	case "warn":
		return ROBTYPE_WARNING
	case "":
		return ROBTYPE_DEBUG
	default:
		break
	}

	if robType != ROBTYPE_EMERGENCY &&
		robType != ROBTYPE_CRITICAL &&
		robType != ROBTYPE_WARNING &&
		robType != ROBTYPE_INFO &&
		robType != ROBTYPE_NOTICE &&
		robType != ROBTYPE_DEBUG &&
		robType != ROBTYPE_ERROR {
		return ROBTYPE_ERROR
	}

	return robType
}

// IsEmpty checks if the ROBType is empty after trimming whitespace.
func (rt ROBType) IsEmpty() bool {
	return strings.TrimSpace(string(rt)) == ""
}

// TrimSpace trims leading and trailing whitespace from the ROBType.
func (rt ROBType) TrimSpace() ROBType {
	return ROBType(strings.TrimSpace(string(rt)))
}

// String returns the string representation of the ROBType.
func (rt ROBType) String() string {
	return string(rt)
}

// ToJsonKey converts the ROBType to a JsonKey.
func (rt ROBType) ToJsonKey() ajson.JsonKey {
	return ajson.JsonKey(rt.String())
}

// HasMatch checks if the ROBType matches another ROBType.
func (rt ROBType) HasMatch(rtType ROBType) bool {
	return rt == rtType
}

// MatchesOne checks if the ROBType matches any one of the provided ROBTypes.
func (rt ROBType) MatchesOne(rtTypes ...ROBType) bool {
	for _, rtType := range rtTypes {
		if rt == rtType {
			return true
		}
	}
	return false
}

// HasPrefix checks if the ROBType has a prefix that matches another ROBType.
func (rt ROBType) HasPrefix(rtType ROBType) bool {
	return strings.HasPrefix(rt.String(), rtType.String())
}

// Validate ensures the ROBType adheres to the specified format rules.
func (rt ROBType) Validate() error {
	target := rt.String()
	for i, r := range target {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && !strings.ContainsRune("-_.", r) {
			return fmt.Errorf("invalid character '%s'; ROBType allows alpha A-Z, a-z, numbers 0-9, special '-_.'", string(r))
		}
		if i == 0 || i == len(target)-1 {
			if strings.ContainsRune("-_.", r) {
				return fmt.Errorf("invalid character '%s' at position %d; only alphas and numbers can begin or end the ROBType", string(r), i)
			}
		}
	}
	return nil
}

// ROBTypes is a slice of ROBType, representing a collection of ROB message types.
type ROBTypes []ROBType

// HasValues checks if the ROBTypes slice contains any values.
func (rts ROBTypes) HasValues() bool {
	return len(rts) > 0
}

// HasMatch checks if any ROBType in the slice matches the provided ROBType.
func (rts ROBTypes) HasMatch(rType ROBType) bool {
	for _, rt := range rts {
		if rt == rType {
			return true
		}
	}
	return false
}

// HasPrefix checks if any ROBType in the slice has a prefix that matches the provided ROBType.
func (rts ROBTypes) HasPrefix(rType ROBType) bool {
	for _, rt := range rts {
		if rt.HasPrefix(rType) {
			return true
		}
	}
	return false
}

// Clone creates a copy of the ROBTypes slice.
func (rts ROBTypes) Clone() ROBTypes {
	var arr ROBTypes
	for _, rt := range rts {
		if !rt.IsEmpty() {
			arr = append(arr, rt)
		}
	}
	return arr
}

// ToArrStrings converts the ROBTypes slice to a slice of strings.
func (rts ROBTypes) ToArrStrings() []string {
	var arr []string
	for _, rt := range rts {
		if !rt.IsEmpty() {
			arr = append(arr, rt.String())
		}
	}
	return arr
}

// IncludeIfInTargets includes ROBTypes that are present in the provided target ROBTypes slice.
func (rts ROBTypes) IncludeIfInTargets(targets ROBTypes) ROBTypes {
	var arr ROBTypes
	for _, rt := range rts {
		if targets.HasMatch(rt) {
			arr = append(arr, rt)
		}
	}
	return arr
}

// Clean removes empty ROBTypes from the slice.
func (rts ROBTypes) Clean() ROBTypes {
	var arr ROBTypes
	for _, rt := range rts {
		if !rt.IsEmpty() {
			arr = append(arr, rt)
		}
	}
	return arr
}
