package auuids

import (
	"encoding/json"
	"github.com/gofrs/uuid/v5"
	"github.com/jpfluger/alibs-slim/autils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestToUUIDFromUUID(t *testing.T) {
	t.Run("Valid UUID", func(t *testing.T) {
		validUUID := autils.NewUUID()
		uid := ToUUIDFromUUID(validUUID)

		assert.False(t, uid.IsNil(), "UID created from a valid UUID should not be nil")
		assert.True(t, uid.Valid, "UID created from a valid UUID should be valid")
		assert.Equal(t, validUUID.String(), uid.String(), "UID string should match the input UUID")
	})

	t.Run("Nil UUID", func(t *testing.T) {
		nilUUID := uuid.Nil
		uid := ToUUIDFromUUID(nilUUID)

		assert.True(t, uid.IsNil(), "UID created from a nil UUID should be nil")
		assert.False(t, uid.Valid, "UID created from a nil UUID should not be valid")
		assert.Equal(t, "", uid.String(), "UID string should be empty for a nil UUID")
	})
}

func TestUID(t *testing.T) {
	t.Run("NewUID", func(t *testing.T) {
		uid := NewUUID()
		assert.False(t, uid.IsNil(), "NewUID should not be nil")
		assert.True(t, uid.Valid, "NewUID should be valid")
	})

	t.Run("ParseUUID Valid", func(t *testing.T) {
		validUUID := autils.NewNullUUIDWithValue()
		uid := ParseUUID(validUUID.UUID.String())
		assert.False(t, uid.IsNil(), "Parsed UID should not be nil")
		assert.True(t, uid.Valid, "Parsed UID should be valid")
		assert.Equal(t, validUUID.UUID.String(), uid.String(), "Parsed UID string should match input")
	})

	t.Run("ParseUUID Invalid", func(t *testing.T) {
		uid := ParseUUID("invalid-uuid")
		assert.True(t, uid.IsNil(), "Invalid UID should be nil")
		assert.False(t, uid.Valid, "Invalid UID should not be valid")
	})

	t.Run("MarshalJSON", func(t *testing.T) {
		uid := NewUUID()
		data, err := json.Marshal(uid)
		assert.NoError(t, err, "Marshaling a valid UID should not error")
		assert.JSONEq(t, `"`+uid.String()+`"`, string(data), "Marshaled UID should match the expected string")

		nilUID := UUID{}
		data, err = json.Marshal(nilUID)
		assert.NoError(t, err, "Marshaling a nil UID should not error")
		assert.JSONEq(t, `null`, string(data), "Marshaled nil UID should be 'null'")
	})

	t.Run("UnmarshalJSON", func(t *testing.T) {
		validUUID := autils.NewNullUUIDWithValue()
		data := `"` + validUUID.UUID.String() + `"`

		var uid UUID
		err := json.Unmarshal([]byte(data), &uid)
		assert.NoError(t, err, "Unmarshaling a valid UID should not error")
		assert.False(t, uid.IsNil(), "Unmarshaled UID should not be nil")
		assert.True(t, uid.Valid, "Unmarshaled UID should be valid")
		assert.Equal(t, validUUID.UUID.String(), uid.String(), "Unmarshaled UID should match the input string")

		invalidData := `"invalid-uuid"`
		err = json.Unmarshal([]byte(invalidData), &uid)
		assert.Error(t, err, "Unmarshaling an invalid UID should return an error")

		nilData := `null`
		err = json.Unmarshal([]byte(nilData), &uid)
		assert.NoError(t, err, "Unmarshaling 'null' should not error")
		assert.True(t, uid.IsNil(), "Unmarshaled UID from 'null' should be nil")
	})
}

func TestUUIDs(t *testing.T) {
	t.Run("Contains", func(t *testing.T) {
		u1 := NewUUID()
		u2 := NewUUID()
		uuids := UUIDs{u1, u2}

		assert.True(t, uuids.Contains(u1), "Contains should return true for an existing UUID")
		assert.False(t, uuids.Contains(NewUUID()), "Contains should return false for a non-existent UUID")
	})

	t.Run("IsValid", func(t *testing.T) {
		u1 := NewUUID()
		u2 := NewUUID()
		uuids := UUIDs{u1, u2}

		assert.True(t, uuids.IsValid(u1, u2), "IsValid should return true for all existing UUIDs")
		assert.False(t, uuids.IsValid(u1, u2, NewUUID()), "IsValid should return false if any UUID is missing")
	})

	t.Run("Merge", func(t *testing.T) {
		u1 := NewUUID()
		u2 := NewUUID()
		u3 := NewUUID()

		uuids1 := UUIDs{u1, u2}
		uuids2 := UUIDs{u2, u3}

		merged := uuids1.Merge(uuids2)
		assert.Len(t, merged, 3, "Merged slice should contain unique elements")
		assert.Contains(t, merged, u1)
		assert.Contains(t, merged, u2)
		assert.Contains(t, merged, u3)
	})

	t.Run("Remove", func(t *testing.T) {
		u1 := NewUUID()
		u2 := NewUUID()
		u3 := NewUUID()

		uuids := UUIDs{u1, u2, u3}
		updated := uuids.Remove(u2)

		assert.Len(t, updated, 2, "Updated slice should have two elements after removal")
		assert.Contains(t, updated, u1)
		assert.Contains(t, updated, u3)
		assert.NotContains(t, updated, u2)
	})

	t.Run("Filter", func(t *testing.T) {
		u1 := NewUUID()
		u2 := NewUUID()
		u3 := NewUUID()

		uuids := UUIDs{u1, u2, u3}
		filtered := uuids.Filter(func(id UUID) bool {
			return id.UUID == u2.UUID
		})

		assert.Len(t, filtered, 1, "Filtered slice should contain only matching elements")
		assert.Contains(t, filtered, u2)
		assert.NotContains(t, filtered, u1)
		assert.NotContains(t, filtered, u3)
	})

	t.Run("Append Without Duplicates", func(t *testing.T) {
		u1 := NewUUID()
		uuids := UUIDs{}

		updated := uuids.Append(u1, false).Append(u1, false)
		assert.Len(t, updated, 1, "Slice should not contain duplicates when allowDuplicates is false")
	})

	t.Run("Append With Duplicates", func(t *testing.T) {
		u1 := NewUUID()
		uuids := UUIDs{}

		updated := uuids.Append(u1, true).Append(u1, true)
		assert.Len(t, updated, 2, "Slice should contain duplicates when allowDuplicates is true")
	})
}

func TestUID_JSONInStruct(t *testing.T) {
	type TestStruct struct {
		UserID UUID `json:"userId,omitempty"`
	}

	t.Run("Marshal Struct with Non-Nil UID", func(t *testing.T) {
		uid := NewUUID()
		ts := TestStruct{UserID: uid}

		data, err := json.Marshal(ts)
		assert.NoError(t, err, "Marshaling a struct with a non-nil UID should not return an error")
		expected := `{"userId":"` + uid.String() + `"}`
		assert.JSONEq(t, expected, string(data), "Marshaled JSON should include the non-nil UID")
	})

	t.Run("Marshal Struct with Nil UID", func(t *testing.T) {
		ts := TestStruct{}
		data, err := json.Marshal(ts)
		assert.NoError(t, err, "Marshaling a struct with a nil UID should not return an error")
		expected := `{"userId":null}`
		assert.JSONEq(t, expected, string(data), "Marshaled JSON should omit the nil UID field")
	})

	t.Run("Marshal Struct with pointer Nil UID", func(t *testing.T) {
		type TestStruct struct {
			UserID *UUID `json:"userId,omitempty"`
		}
		ts := TestStruct{}
		data, err := json.Marshal(ts)
		assert.NoError(t, err, "Marshaling a struct with a nil UID should not return an error")
		expected := `{}`
		assert.JSONEq(t, expected, string(data), "Marshaled JSON should omit the nil UID field")
	})

	t.Run("Unmarshal Struct with Non-Nil UID", func(t *testing.T) {
		uid := NewUUID()
		data := `{"userId":"` + uid.String() + `"}`

		var ts TestStruct
		err := json.Unmarshal([]byte(data), &ts)
		assert.NoError(t, err, "Unmarshaling a struct with a non-nil UID should not return an error")
		assert.Equal(t, uid, ts.UserID, "Unmarshaled UID should match the original")
	})

	t.Run("Unmarshal Struct with Nil UID", func(t *testing.T) {
		data := `{"userId":null}`

		var ts TestStruct
		err := json.Unmarshal([]byte(data), &ts)
		assert.NoError(t, err, "Unmarshaling a struct with a nil UID should not return an error")
		assert.True(t, ts.UserID.IsNil(), "Unmarshaled UID should be nil for 'null' JSON")
	})
}

func TestUID_FromString(t *testing.T) {
	t.Run("Valid UID String", func(t *testing.T) {
		var uid UUID
		err := uid.FromString("123e4567-e89b-12d3-a456-426614174000")
		assert.NoError(t, err)
		assert.Equal(t, "123e4567-e89b-12d3-a456-426614174000", uid.String())
	})

	t.Run("Invalid UID String", func(t *testing.T) {
		var uid UUID
		err := uid.FromString("invalid-uuid")
		assert.Error(t, err)
		assert.True(t, uid.IsNil())
	})
}

func TestUUIDsToString(t *testing.T) {
	t.Run("ToString with multiple UUIDs", func(t *testing.T) {
		id1 := NewUUID()
		id2 := NewUUID()

		uuids := UUIDs{id1, id2}
		expected := id1.String() + "," + id2.String()

		assert.Equal(t, expected, uuids.ToString(), "ToString should return comma-separated UUIDs")
	})

	t.Run("ToString with empty list", func(t *testing.T) {
		uuids := UUIDs{}
		assert.Equal(t, "", uuids.ToString(), "ToString should return empty string for empty slice")
	})
}

func TestUUIDsToStringArray(t *testing.T) {
	t.Run("ToStringArray with multiple UUIDs", func(t *testing.T) {
		id1 := NewUUID()
		id2 := NewUUID()

		ids := UUIDs{id1, id2}
		result := ids.ToStringArray()

		assert.Len(t, result, 2, "ToStringArray should return 2 elements")
		assert.Equal(t, id1.String(), result[0], "First UUID should match")
		assert.Equal(t, id2.String(), result[1], "Second UUID should match")
	})

	t.Run("ToStringArray with empty list", func(t *testing.T) {
		ids := UUIDs{}
		result := ids.ToStringArray()

		assert.Empty(t, result, "ToStringArray should return empty slice for empty UUIDs")
	})
}
