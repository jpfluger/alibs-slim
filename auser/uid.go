package auser

import (
	"fmt"
	"strings"

	"github.com/gofrs/uuid/v5"
	"github.com/jpfluger/alibs-slim/autils"
	"github.com/jpfluger/alibs-slim/auuids"
)

// UID represents a unique user identifier.
type UID struct {
	uuid.NullUUID
}

// NewDeterministicUID generates a Version 5 UUID based on the DNS namespace and a custom string key.
// This guarantees that the same key (e.g., "ebot") always results in the exact same UID.
func NewDeterministicUID(key string) UID {
	// Uses uuid.NamespaceDNS as the base, consistent with standard V5 generation
	return ToUIDFromUUID(uuid.NewV5(uuid.NamespaceDNS, key))
}

// NewUID creates a new UID with a non-nil value using NewNullUUIDWithValue.
func NewUID() UID {
	return UID{NullUUID: autils.NewNullUUIDWithValue()}
}

// ToUIDFromUUID creates a new UID with a non-nil value using NewNullUUIDWithValue.
func ToUIDFromUUID(id uuid.UUID) UID {
	return UID{NullUUID: uuid.NullUUID{UUID: id, Valid: id != uuid.Nil}}
}

// ParseUID parses a string into a UID using ParseNullUUID.
func ParseUID(target string) UID {
	return UID{NullUUID: autils.ParseNullUUID(target)}
}

// IsNil checks if the UID is nil or invalid.
func (sid UID) IsNil() bool {
	return !sid.Valid || sid.UUID == uuid.Nil
}

// String converts the UID to its string representation.
func (sid UID) String() string {
	if sid.IsNil() {
		return ""
	}
	return sid.UUID.String()
}

// ToUUID transforms the type to UUID
func (sid UID) ToUUID() auuids.UUID {
	return auuids.UUID{sid.NullUUID}
}

// FromString parses a string and assigns the UID to itself.
func (sid *UID) FromString(input string) error {
	if sid == nil {
		return fmt.Errorf("UID: FromString called on nil pointer")
	}

	*sid = ParseUID(input)
	if sid.IsNil() {
		return fmt.Errorf("invalid UUID string: %s", input)
	}
	return nil
}

// HasMatch checks if the UID matches the provided UID.
func (sid UID) HasMatch(target UID) bool {
	return sid.UUID == target.UUID && sid.Valid == target.Valid
}

// MatchesOne checks if the UID matches any of the provided UIDs.
func (sid UID) MatchesOne(targets ...UID) bool {
	for _, target := range targets {
		if sid.HasMatch(target) {
			return true
		}
	}
	return false
}

// MarshalJSON ensures that nil UIDs are marshaled as null.
func (sid UID) MarshalJSON() ([]byte, error) {
	return sid.NullUUID.MarshalJSON()
}

// UnmarshalJSON parses a UID from a JSON string, handling "null" or empty values.
func (sid *UID) UnmarshalJSON(data []byte) error {
	if sid == nil {
		return fmt.Errorf("UID: UnmarshalJSON on nil pointer")
	}

	str := strings.Trim(string(data), `"`)
	if str == "null" || str == "" {
		*sid = UID{} // Set to zero value
		return nil
	}

	return sid.FromString(str)
}

// MarshalText converts the UID to its string representation for text-based encoding.
func (sid UID) MarshalText() ([]byte, error) {
	if sid.IsNil() {
		return []byte(""), nil
	}
	return []byte(sid.UUID.String()), nil
}

// UnmarshalText parses a UID from a plain text string.
func (sid *UID) UnmarshalText(text []byte) error {
	if sid == nil {
		return fmt.Errorf("UID: UnmarshalText on nil pointer")
	}

	str := strings.TrimSpace(string(text))
	if str == "" {
		*sid = UID{} // Reset to zero value
		return nil
	}

	return sid.FromString(str)
}

// UIDs represents a collection of UID values.
type UIDs []UID

// Contains checks if the slice contains the given UID.
func (sids UIDs) Contains(sid UID) bool {
	for _, item := range sids {
		if item.HasMatch(sid) {
			return true
		}
	}
	return false
}

// IsValid checks if all given UIDs exist in the UIDs slice.
func (sids UIDs) IsValid(targets ...UID) bool {
	for _, target := range targets {
		if !sids.Contains(target) {
			return false
		}
	}
	return true
}

// Merge combines two UIDs slices, removing duplicates.
func (sids UIDs) Merge(other UIDs) UIDs {
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

// Remove creates a new slice that excludes the specified UID.
func (sids UIDs) Remove(target UID) UIDs {
	result := UIDs{}
	for _, id := range sids {
		if id.UUID != target.UUID {
			result = append(result, id)
		}
	}
	return result
}

// Filter creates a new slice containing only UIDs that match the given predicate.
func (sids UIDs) Filter(predicate func(UID) bool) UIDs {
	result := UIDs{}
	for _, id := range sids {
		if predicate(id) {
			result = append(result, id)
		}
	}
	return result
}

// Append adds a UID to the slice. Optionally prevents duplicates if allowDuplicates is false.
func (sids UIDs) Append(id UID, allowDuplicates bool) UIDs {
	if !allowDuplicates && sids.Contains(id) {
		return sids
	}
	return append(sids, id)
}

// Validate checks if all UUIDs in the slice are non-nil.
func (sids UIDs) Validate() error {
	return sids.ValidateWithOptions(false)
}

// ValidateWithOptions checks if all UUIDs in the slice are non-nil.
// If mustHaveCount is true, it also checks that the slice is not empty.
func (sids UIDs) ValidateWithOptions(mustHaveCount bool) error {
	if len(sids) == 0 {
		if mustHaveCount {
			return fmt.Errorf("uids is empty")
		}
		return nil
	}
	for ii, uid := range sids {
		if uid.IsNil() {
			return fmt.Errorf("uid nil at index %d", ii)
		}
	}
	return nil
}

// Clean returns a new slice with duplicate and nil UUIDs removed.
func (sids UIDs) Clean() UIDs {
	var arr UIDs
	seen := make(map[UID]bool)
	for _, id := range sids {
		if id.IsNil() || seen[id] {
			continue
		}
		seen[id] = true
		arr = append(arr, id)
	}
	return arr
}
