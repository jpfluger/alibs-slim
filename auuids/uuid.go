package auuids

import (
	"fmt"
	"github.com/gofrs/uuid/v5"
	"github.com/jpfluger/alibs-slim/autils"
	"strings"
)

// UUID Represents a unique identifier (e.g., UUID as a string).
// It uses NullUUID internally to handle nil or invalid cases.
type UUID struct {
	uuid.NullUUID
}

// NewUUID creates a new UUID with a non-nil value using NewNullUUIDWithValue.
func NewUUID() UUID {
	return UUID{NullUUID: autils.NewNullUUIDWithValue()}
}

// ToUUIDFromUUID creates a new UID with a non-nil value using NewNullUUIDWithValue.
func ToUUIDFromUUID(id uuid.UUID) UUID {
	return UUID{NullUUID: uuid.NullUUID{UUID: id, Valid: id != uuid.Nil}}
}

// ParseUUID parses a string into a UUID using ParseNullUUID.
func ParseUUID(target string) UUID {
	return UUID{NullUUID: autils.ParseNullUUID(target)}
}

// IsNil checks if the UUID is nil or invalid.
func (u UUID) IsNil() bool {
	return !u.Valid || u.UUID == uuid.Nil
}

// String converts the UUID to its string representation.
func (u UUID) String() string {
	if u.IsNil() {
		return ""
	}
	return u.UUID.String()
}

// FromString parses a string and assigns the UUID to itself.
func (u *UUID) FromString(input string) error {
	if u == nil {
		return fmt.Errorf("UUID: FromString called on nil pointer")
	}

	*u = ParseUUID(input)
	if u.IsNil() {
		return fmt.Errorf("invalid UUID string: %s", input)
	}
	return nil
}

// HasMatch checks if the UUID matches the provided UUID.
func (u UUID) HasMatch(target UUID) bool {
	return u.UUID == target.UUID && u.Valid == target.Valid
}

// MatchesOne checks if the UUID matches any of the provided UUIDs.
func (u UUID) MatchesOne(targets ...UUID) bool {
	for _, target := range targets {
		if u.HasMatch(target) {
			return true
		}
	}
	return false
}

// MarshalJSON ensures that nil UUIDs are marshaled as null.
func (u UUID) MarshalJSON() ([]byte, error) {
	return u.NullUUID.MarshalJSON()
}

// UnmarshalJSON parses a UUID from a JSON string, handling "null" or empty values.
func (u *UUID) UnmarshalJSON(data []byte) error {
	if u == nil {
		return fmt.Errorf("UUID: UnmarshalJSON on nil pointer")
	}

	str := strings.Trim(string(data), `"`)
	if str == "null" || str == "" {
		*u = UUID{} // Set to zero value
		return nil
	}

	return u.FromString(str)
}

// MarshalText converts the UUID to its string representation for text-based encoding.
func (u UUID) MarshalText() ([]byte, error) {
	if u.IsNil() {
		return []byte(""), nil
	}
	return []byte(u.UUID.String()), nil
}

// UnmarshalText parses a UUID from a plain text string.
func (u *UUID) UnmarshalText(text []byte) error {
	if u == nil {
		return fmt.Errorf("UUID: UnmarshalText on nil pointer")
	}

	str := strings.TrimSpace(string(text))
	if str == "" {
		*u = UUID{} // Reset to zero value
		return nil
	}

	return u.FromString(str)
}

// UUIDs represents a collection of UUID values.
type UUIDs []UUID

// Contains checks if the slice contains the given UUID.
func (uuids UUIDs) Contains(uuid UUID) bool {
	for _, item := range uuids {
		if item.HasMatch(uuid) {
			return true
		}
	}
	return false
}

// IsValid checks if all given UUIDs exist in the UUIDs slice.
func (uuids UUIDs) IsValid(targets ...UUID) bool {
	for _, target := range targets {
		if !uuids.Contains(target) {
			return false
		}
	}
	return true
}

// Merge combines two UUID slices, removing duplicates.
// This modifies the original slice and returns the updated version for chaining.
func (uuids UUIDs) Merge(other UUIDs) UUIDs {
	existing := make(map[uuid.UUID]bool)
	for _, id := range uuids {
		existing[id.UUID] = true
	}

	for _, id := range other {
		if !existing[id.UUID] {
			uuids = append(uuids, id)
		}
	}
	return uuids
}

// Remove creates a new slice that excludes the specified UUID.
func (uuids UUIDs) Remove(target UUID) UUIDs {
	result := UUIDs{}
	for _, id := range uuids {
		if id.UUID != target.UUID {
			result = append(result, id)
		}
	}
	return result
}

// Filter creates a new slice containing only UUIDs that match the given predicate.
func (uuids UUIDs) Filter(predicate func(UUID) bool) UUIDs {
	result := UUIDs{}
	for _, id := range uuids {
		if predicate(id) {
			result = append(result, id)
		}
	}
	return result
}

// Append adds a UUID to the slice. Optionally prevents duplicates if allowDuplicates is false.
func (uuids UUIDs) Append(id UUID, allowDuplicates bool) UUIDs {
	if !allowDuplicates && uuids.Contains(id) {
		return uuids
	}
	return append(uuids, id)
}

// Validate checks if all UUIDs in the slice are non-nil.
func (uuids UUIDs) Validate() error {
	return uuids.ValidateWithOptions(false)
}

// ValidateWithOptions checks if all UUIDs in the slice are non-nil.
// If mustHaveCount is true, it also checks that the slice is not empty.
func (uuids UUIDs) ValidateWithOptions(mustHaveCount bool) error {
	if len(uuids) == 0 {
		if mustHaveCount {
			return fmt.Errorf("uuids is empty")
		}
		return nil
	}
	for ii, uid := range uuids {
		if uid.IsNil() {
			return fmt.Errorf("uuid nil at index %d", ii)
		}
	}
	return nil
}

// Clean returns a new slice with duplicate and nil UUIDs removed.
func (uuids UUIDs) Clean() UUIDs {
	var arr UUIDs
	seen := make(map[UUID]bool)
	for _, id := range uuids {
		if id.IsNil() || seen[id] {
			continue
		}
		seen[id] = true
		arr = append(arr, id)
	}
	return arr
}

// ToString returns a comma-separated string of all UUIDs in the slice.
// Example: "uuid1,uuid2,uuid3"
func (uuids UUIDs) ToString() string {
	var sb strings.Builder
	for i, id := range uuids {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(id.String())
	}
	return sb.String()
}

// ToStringArray returns a slice of strings, where each element is the string
// representation of a UUID in the UUIDs slice.
// Example: []string{"uuid1", "uuid2", "uuid3"}
func (uuids UUIDs) ToStringArray() []string {
	out := make([]string, 0, len(uuids))
	for _, id := range uuids {
		out = append(out, id.String())
	}
	return out
}
