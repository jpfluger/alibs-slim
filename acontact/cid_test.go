package acontact

import (
	"encoding/json"
	"github.com/gofrs/uuid/v5"
	"github.com/jpfluger/alibs-slim/autils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToCIDFromUUID(t *testing.T) {
	t.Run("Valid UUID", func(t *testing.T) {
		validUUID := autils.NewUUID()
		uid := ToCIDFromUUID(validUUID)

		assert.False(t, uid.IsNil(), "UID created from a valid UUID should not be nil")
		assert.True(t, uid.Valid, "UID created from a valid UUID should be valid")
		assert.Equal(t, validUUID.String(), uid.String(), "UID string should match the input UUID")
	})

	t.Run("Nil UUID", func(t *testing.T) {
		nilUUID := uuid.Nil
		uid := ToCIDFromUUID(nilUUID)

		assert.True(t, uid.IsNil(), "UID created from a nil UUID should be nil")
		assert.False(t, uid.Valid, "UID created from a nil UUID should not be valid")
		assert.Equal(t, "", uid.String(), "UID string should be empty for a nil UUID")
	})
}

func TestUID(t *testing.T) {
	t.Run("NewCID", func(t *testing.T) {
		did := NewCID()
		assert.False(t, did.IsNil(), "NewCID should not be nil")
		assert.True(t, did.Valid, "NewCID should be valid")
	})

	t.Run("ParseCID Valid", func(t *testing.T) {
		validUUID := autils.NewNullUUIDWithValue()
		did := ParseCID(validUUID.UUID.String())
		assert.False(t, did.IsNil(), "Parsed UID should not be nil")
		assert.True(t, did.Valid, "Parsed UID should be valid")
		assert.Equal(t, validUUID.UUID.String(), did.String(), "Parsed UID string should match input")
	})

	t.Run("ParseCID Invalid", func(t *testing.T) {
		did := ParseCID("invalid-uuid")
		assert.True(t, did.IsNil(), "Invalid UID should be nil")
		assert.False(t, did.Valid, "Invalid UID should not be valid")
	})

	t.Run("MarshalJSON", func(t *testing.T) {
		did := NewCID()
		data, err := json.Marshal(did)
		assert.NoError(t, err, "Marshaling a valid UID should not error")
		assert.JSONEq(t, `"`+did.String()+`"`, string(data), "Marshaled UID should match the expected string")

		nilDID := CID{}
		data, err = json.Marshal(nilDID)
		assert.NoError(t, err, "Marshaling a nil UID should not error")
		assert.JSONEq(t, `null`, string(data), "Marshaled nil UID should be 'null'")
	})

	t.Run("UnmarshalJSON", func(t *testing.T) {
		validUUID := autils.NewNullUUIDWithValue()
		data := `"` + validUUID.UUID.String() + `"`

		var did CID
		err := json.Unmarshal([]byte(data), &did)
		assert.NoError(t, err, "Unmarshaling a valid UID should not error")
		assert.False(t, did.IsNil(), "Unmarshaled UID should not be nil")
		assert.True(t, did.Valid, "Unmarshaled UID should be valid")
		assert.Equal(t, validUUID.UUID.String(), did.String(), "Unmarshaled UID should match the input string")

		invalidData := `"invalid-uuid"`
		err = json.Unmarshal([]byte(invalidData), &did)
		assert.Error(t, err, "Unmarshaling an invalid UID should return an error")

		nilData := `null`
		err = json.Unmarshal([]byte(nilData), &did)
		assert.NoError(t, err, "Unmarshaling 'null' should not error")
		assert.True(t, did.IsNil(), "Unmarshaled UID from 'null' should be nil")
	})
}

func TestUIDs(t *testing.T) {
	t.Run("Contains", func(t *testing.T) {
		did1 := NewCID()
		did2 := NewCID()
		dids := CIDs{did1, did2}

		assert.True(t, dids.Contains(did1), "UIDs should contain did1")
		assert.False(t, dids.Contains(NewCID()), "UIDs should not contain a different UID")
	})

	t.Run("IsValid", func(t *testing.T) {
		did1 := NewCID()
		did2 := NewCID()
		dids := CIDs{did1, did2}

		assert.True(t, dids.IsValid(did1, did2), "All UIDs in the input should exist in the UIDs slice")
		assert.False(t, dids.IsValid(did1, did2, NewCID()), "Should return false if any UID in the input is missing")
	})

	t.Run("Merge", func(t *testing.T) {
		did1 := NewCID()
		did2 := NewCID()
		did3 := NewCID()

		dids := CIDs{did1, did2}
		other := CIDs{did2, did3}

		merged := dids.Merge(other)
		assert.Len(t, merged, 3, "Merged UIDs should contain three unique IDs")
		assert.True(t, merged.Contains(did1), "Merged UIDs should contain did1")
		assert.True(t, merged.Contains(did2), "Merged UIDs should contain did2")
		assert.True(t, merged.Contains(did3), "Merged UIDs should contain did3")
	})

	t.Run("Remove", func(t *testing.T) {
		did1 := NewCID()
		did2 := NewCID()

		dids := CIDs{did1, did2}
		updated := dids.Remove(did1)

		assert.Len(t, updated, 1, "UIDs should contain one less after removal")
		assert.False(t, updated.Contains(did1), "Removed UID should not be in the slice")
		assert.True(t, updated.Contains(did2), "Non-removed UID should still be in the slice")
	})

	t.Run("Filter", func(t *testing.T) {
		did1 := NewCID()
		did2 := NewCID()

		dids := CIDs{did1, did2}
		filtered := dids.Filter(func(did CID) bool {
			return did == did1
		})

		assert.Len(t, filtered, 1, "Filtered UIDs should contain only matching IDs")
		assert.True(t, filtered.Contains(did1), "Filtered UIDs should contain did1")
		assert.False(t, filtered.Contains(did2), "Filtered UIDs should not contain did2")
	})

	t.Run("Append", func(t *testing.T) {
		did1 := NewCID()
		dids := CIDs{}

		dids = dids.Append(did1, false)
		assert.Len(t, dids, 1, "Appended UIDs should contain one item")
		assert.True(t, dids.Contains(did1), "Appended UIDs should contain did1")

		dids = dids.Append(did1, false)
		assert.Len(t, dids, 1, "Appended UIDs should not allow duplicates when allowDuplicates is false")

		dids = dids.Append(did1, true)
		assert.Len(t, dids, 2, "Appended UIDs should allow duplicates when allowDuplicates is true")
	})
}
