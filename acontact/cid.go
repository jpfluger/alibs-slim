package acontact

import (
	"fmt"
	"github.com/gofrs/uuid/v5"
	"github.com/jpfluger/alibs-slim/autils"
	"strings"
)

// CID represents a unique contact identifier.
type CID struct {
	uuid.NullUUID
}

// NewCID creates a new CID with a non-nil value using NewNullUCIDWithValue.
func NewCID() CID {
	return CID{NullUUID: autils.NewNullUUIDWithValue()}
}

// ToCIDFromUUID creates a new CID with a non-nil value using NewNullUUIDWithValue.
func ToCIDFromUUID(id uuid.UUID) CID {
	return CID{NullUUID: uuid.NullUUID{UUID: id, Valid: id != uuid.Nil}}
}

// ParseCID parses a string into a CID using ParseNullUUID.
func ParseCID(target string) CID {
	return CID{NullUUID: autils.ParseNullUUID(target)}
}

// IsNil checks if the CID is nil or invalid.
func (sid CID) IsNil() bool {
	return !sid.Valid || sid.UUID == uuid.Nil
}

// String converts the CID to its string representation.
func (sid CID) String() string {
	if sid.IsNil() {
		return ""
	}
	return sid.UUID.String()
}

// FromString parses a string and assigns the CID to itself.
func (sid *CID) FromString(input string) error {
	if sid == nil {
		return fmt.Errorf("CID: FromString called on nil pointer")
	}

	*sid = ParseCID(input)
	if sid.IsNil() {
		return fmt.Errorf("invalid UUID string: %s", input)
	}
	return nil
}

// HasMatch checks if the CID matches the provided CID.
func (sid CID) HasMatch(target CID) bool {
	return sid.UUID == target.UUID && sid.Valid == target.Valid
}

// MatchesOne checks if the CID matches any of the provided CIDs.
func (sid CID) MatchesOne(targets ...CID) bool {
	for _, target := range targets {
		if sid.HasMatch(target) {
			return true
		}
	}
	return false
}

// MarshalJSON ensures that nil CIDs are marshaled as null.
func (sid CID) MarshalJSON() ([]byte, error) {
	return sid.NullUUID.MarshalJSON()
}

// UnmarshalJSON parses a CID from a JSON string, handling "null" or empty values.
func (sid *CID) UnmarshalJSON(data []byte) error {
	if sid == nil {
		return fmt.Errorf("CID: UnmarshalJSON on nil pointer")
	}

	str := strings.Trim(string(data), `"`)
	if str == "null" || str == "" {
		*sid = CID{} // Set to zero value
		return nil
	}

	return sid.FromString(str)
}

// MarshalText converts the CID to its string representation for text-based encoding.
func (sid CID) MarshalText() ([]byte, error) {
	if sid.IsNil() {
		return []byte(""), nil
	}
	return []byte(sid.UUID.String()), nil
}

// UnmarshalText parses a CID from a plain text string.
func (sid *CID) UnmarshalText(text []byte) error {
	if sid == nil {
		return fmt.Errorf("CID: UnmarshalText on nil pointer")
	}

	str := strings.TrimSpace(string(text))
	if str == "" {
		*sid = CID{} // Reset to zero value
		return nil
	}

	return sid.FromString(str)
}

// CIDs represents a collection of CID values.
type CIDs []CID

// Contains checks if the slice contains the given CID.
func (sids CIDs) Contains(sid CID) bool {
	for _, item := range sids {
		if item.HasMatch(sid) {
			return true
		}
	}
	return false
}

// IsValid checks if all given CIDs exist in the CIDs slice.
func (sids CIDs) IsValid(targets ...CID) bool {
	for _, target := range targets {
		if !sids.Contains(target) {
			return false
		}
	}
	return true
}

// Merge combines two CIDs slices, removing duplicates.
func (sids CIDs) Merge(other CIDs) CIDs {
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

// Remove creates a new slice that excludes the specified CID.
func (sids CIDs) Remove(target CID) CIDs {
	result := CIDs{}
	for _, id := range sids {
		if id.UUID != target.UUID {
			result = append(result, id)
		}
	}
	return result
}

// Filter creates a new slice containing only CIDs that match the given predicate.
func (sids CIDs) Filter(predicate func(CID) bool) CIDs {
	result := CIDs{}
	for _, id := range sids {
		if predicate(id) {
			result = append(result, id)
		}
	}
	return result
}

// Append adds a CID to the slice. Optionally prevents duplicates if allowDuplicates is false.
func (sids CIDs) Append(id CID, allowDuplicates bool) CIDs {
	if !allowDuplicates && sids.Contains(id) {
		return sids
	}
	return append(sids, id)
}

// Validate checks if all UUIDs in the slice are non-nil.
func (sids CIDs) Validate() error {
	return sids.ValidateWithOptions(false)
}

// ValidateWithOptions checks if all UUIDs in the slice are non-nil.
// If mustHaveCount is true, it also checks that the slice is not empty.
func (sids CIDs) ValidateWithOptions(mustHaveCount bool) error {
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
func (sids CIDs) Clean() CIDs {
	var arr CIDs
	seen := make(map[CID]bool)
	for _, id := range sids {
		if id.IsNil() || seen[id] {
			continue
		}
		seen[id] = true
		arr = append(arr, id)
	}
	return arr
}
