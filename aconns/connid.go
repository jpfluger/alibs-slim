package aconns

import (
	"fmt"
	"github.com/gofrs/uuid/v5"
	"github.com/jpfluger/alibs-slim/autils"
	"strings"
)

// ConnId represents a unique service identifier.
type ConnId struct {
	uuid.NullUUID
}

// NewConnId creates a new ConnId with a non-nil value using NewNullUUIDWithValue.
func NewConnId() ConnId {
	return ConnId{NullUUID: autils.NewNullUUIDWithValue()}
}

// ToConnIdFromUUID creates a new UID with a non-nil value using NewNullUUIDWithValue.
func ToConnIdFromUUID(id uuid.UUID) ConnId {
	return ConnId{NullUUID: uuid.NullUUID{UUID: id, Valid: id != uuid.Nil}}
}

// ParseConnId parses a string into a ConnId using ParseNullUUID.
func ParseConnId(target string) ConnId {
	return ConnId{NullUUID: autils.ParseNullUUID(target)}
}

// IsNil checks if the ConnId is nil or invalid.
func (sid ConnId) IsNil() bool {
	return !sid.Valid || sid.UUID == uuid.Nil
}

// String converts the ConnId to its string representation.
func (sid ConnId) String() string {
	if sid.IsNil() {
		return ""
	}
	return sid.UUID.String()
}

// FromString parses a string and assigns the ConnId to itself.
func (sid *ConnId) FromString(input string) error {
	if sid == nil {
		return fmt.Errorf("ConnId: FromString called on nil pointer")
	}

	*sid = ParseConnId(input)
	if sid.IsNil() {
		return fmt.Errorf("invalid UUID string: %s", input)
	}
	return nil
}

// HasMatch checks if the ConnId matches the provided ConnId.
func (sid ConnId) HasMatch(target ConnId) bool {
	return sid.UUID == target.UUID && sid.Valid == target.Valid
}

// MatchesOne checks if the ConnId matches any of the provided ConnIds.
func (sid ConnId) MatchesOne(targets ...ConnId) bool {
	for _, target := range targets {
		if sid.HasMatch(target) {
			return true
		}
	}
	return false
}

// MarshalJSON ensures that nil ConnIds are marshaled as null.
func (sid ConnId) MarshalJSON() ([]byte, error) {
	return sid.NullUUID.MarshalJSON()
}

// UnmarshalJSON parses a ConnId from a JSON string, handling "null" or empty values.
func (sid *ConnId) UnmarshalJSON(data []byte) error {
	if sid == nil {
		return fmt.Errorf("ConnId: UnmarshalJSON on nil pointer")
	}

	str := strings.Trim(string(data), `"`)
	if str == "null" || str == "" {
		*sid = ConnId{} // Set to zero value
		return nil
	}

	return sid.FromString(str)
}

// MarshalText converts the ConnId to its string representation for text-based encoding.
func (sid ConnId) MarshalText() ([]byte, error) {
	if sid.IsNil() {
		return []byte(""), nil
	}
	return []byte(sid.UUID.String()), nil
}

// UnmarshalText parses a ConnId from a plain text string.
func (sid *ConnId) UnmarshalText(text []byte) error {
	if sid == nil {
		return fmt.Errorf("ConnId: UnmarshalText on nil pointer")
	}

	str := strings.TrimSpace(string(text))
	if str == "" {
		*sid = ConnId{} // Reset to zero value
		return nil
	}

	return sid.FromString(str)
}

// ConnIds represents a collection of ConnId values.
type ConnIds []ConnId

// Contains checks if the slice contains the given ConnId.
func (sids ConnIds) Contains(sid ConnId) bool {
	for _, item := range sids {
		if item.HasMatch(sid) {
			return true
		}
	}
	return false
}

// IsValid checks if all given ConnIds exist in the ConnIds slice.
func (sids ConnIds) IsValid(targets ...ConnId) bool {
	for _, target := range targets {
		if !sids.Contains(target) {
			return false
		}
	}
	return true
}

// Merge combines two ConnIds slices, removing duplicates.
func (sids ConnIds) Merge(other ConnIds) ConnIds {
	existing := make(map[uuid.UUID]bool)
	for _, id := range sids {
		existing[id.UUID] = true
	}

	for _, id := range other {
		if !existing[id.UUID] {
			sids = append(sids, id)
		}
	}
	return sids
}

// Remove creates a new slice that excludes the specified ConnId.
func (sids ConnIds) Remove(target ConnId) ConnIds {
	result := ConnIds{}
	for _, id := range sids {
		if id.UUID != target.UUID {
			result = append(result, id)
		}
	}
	return result
}

// Filter creates a new slice containing only ConnIds that match the given predicate.
func (sids ConnIds) Filter(predicate func(ConnId) bool) ConnIds {
	result := ConnIds{}
	for _, id := range sids {
		if predicate(id) {
			result = append(result, id)
		}
	}
	return result
}

// Append adds a ConnId to the slice. Optionally prevents duplicates if allowDuplicates is false.
func (sids ConnIds) Append(id ConnId, allowDuplicates bool) ConnIds {
	if !allowDuplicates && sids.Contains(id) {
		return sids
	}
	return append(sids, id)
}

// Clean removes empty (nil) records from ConnIds.
func (sids ConnIds) Clean() ConnIds {
	result := ConnIds{}
	if sids == nil || len(sids) == 0 {
		return result
	}
	for _, item := range sids {
		if !item.IsNil() {
			result = append(result, item)
		}
	}
	return result
}
