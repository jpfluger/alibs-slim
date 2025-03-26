package aconns

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRecordIdentifierMeta(t *testing.T) {
	meta := NewRecordIdentifierMeta("email", "user@example.com")
	assert.NotNil(t, meta)
	assert.Equal(t, "email", string(meta.Type))
	assert.Equal(t, "user@example.com", meta.Id)

	meta = NewRecordIdentifierMeta("", "user@example.com")
	err := meta.Validate()
	assert.Error(t, err)

	meta = NewRecordIdentifierMeta("email", "")
	err = meta.Validate()
	assert.Error(t, err)
}

func TestRecordIdentifierMeta_Validate(t *testing.T) {
	meta := NewRecordIdentifierMeta("email", "user@example.com")
	assert.NoError(t, meta.Validate())

	meta = NewRecordIdentifierMeta("", "user@example.com")
	err := meta.Validate()
	assert.EqualError(t, err, "empty Type")

	meta = NewRecordIdentifierMeta("email", "")
	err = meta.Validate()
	assert.EqualError(t, err, "empty ID")
}

func TestRecordIdentifierMeta_HasMatch(t *testing.T) {
	meta1 := NewRecordIdentifierMeta("email", "user@example.com")
	meta2 := NewRecordIdentifierMeta("email", "user@example.com")
	assert.True(t, meta1.HasMatch(meta2))

	meta3 := NewRecordIdentifierMeta("phone", "123-456-7890")
	assert.False(t, meta1.HasMatch(meta3))
}

func TestRecordIdentifierMeta_HasMatchWithTypeId(t *testing.T) {
	meta := NewRecordIdentifierMeta("email", "user@example.com")
	assert.True(t, meta.HasMatchWithTypeId("email", "user@example.com"))
	assert.False(t, meta.HasMatchWithTypeId("phone", "user@example.com"))
}

func TestRecordIdentifierMetas_Set(t *testing.T) {
	metas := RecordIdentifierMetas{}

	// Add a new record
	metas, err := metas.Set(NewRecordIdentifierMeta("email", "user@example.com"))
	assert.NoError(t, err)
	assert.Len(t, metas, 1)
	assert.Equal(t, "email", string(metas[0].Type))
	assert.Equal(t, "user@example.com", metas[0].Id)

	// Add another record
	metas, err = metas.Set(NewRecordIdentifierMeta("phone", "123-456-7890"))
	assert.NoError(t, err)
	assert.Len(t, metas, 2)

	// Update the existing "email" record
	metas, err = metas.SetByType(NewRecordIdentifierMeta("email", "updated@example.com"))
	assert.NoError(t, err)
	assert.Len(t, metas, 2)                             // Length should remain the same
	assert.Equal(t, "updated@example.com", metas[0].Id) // Email updated
	assert.Equal(t, "phone", string(metas[1].Type))     // Phone remains unchanged

	// Add a new "fax" record
	metas, err = metas.Set(NewRecordIdentifierMeta("fax", "fax@example.com"))
	assert.NoError(t, err)
	assert.Len(t, metas, 3) // Length increases as it's a new Type
	assert.Equal(t, "fax", string(metas[2].Type))
	assert.Equal(t, "fax@example.com", metas[2].Id)
}

func TestRecordIdentifierMetas_RemoveExact(t *testing.T) {
	metas := RecordIdentifierMetas{
		NewRecordIdentifierMeta("email", "user@example.com"),
		NewRecordIdentifierMeta("phone", "123-456-7890"),
	}

	// Remove exact match
	target := NewRecordIdentifierMeta("email", "user@example.com")
	updated, err := metas.RemoveExact(target)
	assert.NoError(t, err)
	assert.Len(t, updated, 1)
	assert.Equal(t, "phone", string(updated[0].Type))

	// Attempt to remove non-existent exact match
	nonexistent := NewRecordIdentifierMeta("email", "nonexistent@example.com")
	updated, err = metas.RemoveExact(nonexistent)
	assert.NoError(t, err)
	assert.Len(t, updated, 2) // No changes
}

func TestRecordIdentifierMetas_RemoveById(t *testing.T) {
	metas := RecordIdentifierMetas{
		NewRecordIdentifierMeta("email", "user@example.com"),
		NewRecordIdentifierMeta("phone", "123-456-7890"),
	}

	// Remove by ID only
	updated, err := metas.RemoveById("user@example.com")
	assert.NoError(t, err)
	assert.Len(t, updated, 1)
	assert.Equal(t, "phone", string(updated[0].Type))

	// Attempt to remove non-existent ID
	updated, err = updated.RemoveById("nonexistent@example.com")
	assert.NoError(t, err)
	assert.Len(t, updated, 1) // No changes
}

func TestRecordIdentifierMetas_RemoveByTypeAndId(t *testing.T) {
	metas := RecordIdentifierMetas{
		NewRecordIdentifierMeta("email", "user@example.com"),
		NewRecordIdentifierMeta("phone", "123-456-7890"),
		NewRecordIdentifierMeta("email", "another@example.com"),
	}

	// Remove by Type and ID
	updated, err := metas.RemoveByTypeAndId("email", "user@example.com")
	assert.NoError(t, err)
	assert.Len(t, updated, 2)
	assert.Equal(t, "phone", string(updated[0].Type))
	assert.Equal(t, "email", string(updated[1].Type))
	assert.Equal(t, "another@example.com", updated[1].Id)

	// Attempt to remove non-existent Type and ID
	updated, err = updated.RemoveByTypeAndId("fax", "123-456-7890")
	assert.NoError(t, err)
	assert.Len(t, updated, 2) // No changes
}

func TestRecordIdentifierMetas_RemoveByType(t *testing.T) {
	metas := RecordIdentifierMetas{
		NewRecordIdentifierMeta("email", "user@example.com"),
		NewRecordIdentifierMeta("phone", "123-456-7890"),
		NewRecordIdentifierMeta("email", "another@example.com"),
	}

	// Remove by Type only
	updated, err := metas.RemoveByType("email")
	assert.NoError(t, err)
	assert.Len(t, updated, 1)
	assert.Equal(t, "phone", string(updated[0].Type))

	// Attempt to remove non-existent Type
	updated, err = updated.RemoveByType("fax")
	assert.NoError(t, err)
	assert.Len(t, updated, 1) // No changes
}

func TestRecordIdentifierMetas_FindWithMatchingOptions(t *testing.T) {
	metas := RecordIdentifierMetas{
		NewRecordIdentifierMeta("email", "user@example.com"),
		NewRecordIdentifierMeta("phone", "123-456-7890"),
		NewRecordIdentifierMeta("email", "another@example.com"),
	}

	// Test match by both Type and ID
	target := NewRecordIdentifierMeta("email", "user@example.com")
	found := metas.FindWithMatchingOptions(target, true, true)
	assert.NotNil(t, found)
	assert.Equal(t, "email", string(found.Type))
	assert.Equal(t, "user@example.com", found.Id)

	// Test match by ID only
	found = metas.FindWithMatchingOptions(&RecordIdentifierMeta{Id: "123-456-7890"}, false, true)
	assert.NotNil(t, found)
	assert.Equal(t, "phone", string(found.Type))

	// Test match by Type only
	found = metas.FindWithMatchingOptions(&RecordIdentifierMeta{Type: "email"}, true, false)
	assert.NotNil(t, found)
	assert.Equal(t, "email", string(found.Type))
	assert.Equal(t, "user@example.com", found.Id) // Returns the first match of Type

	// Test fallback general match
	found = metas.FindWithMatchingOptions(target, false, false)
	assert.NotNil(t, found)
	assert.Equal(t, "email", string(found.Type))
	assert.Equal(t, "user@example.com", found.Id)

	// Test no match
	nonexistent := NewRecordIdentifierMeta("email", "nonexistent@example.com")
	found = metas.FindWithMatchingOptions(nonexistent, true, true)
	assert.Nil(t, found)
}

func TestRecordIdentifierMetas_HasMatch(t *testing.T) {
	metas := RecordIdentifierMetas{
		NewRecordIdentifierMeta("email", "user@example.com"),
	}
	target := NewRecordIdentifierMeta("email", "user@example.com")
	assert.True(t, metas.HasMatch(target))

	target = NewRecordIdentifierMeta("phone", "123-456-7890")
	assert.False(t, metas.HasMatch(target))
}
