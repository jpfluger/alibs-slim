package autils

import (
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/stretchr/testify/assert"
)

// TestNewUUID checks if NewUUID function generates a non-nil UUID.
func TestNewUUID(t *testing.T) {
	id := NewUUID()
	assert.NotEqual(t, uuid.Nil, id, "NewUUID should generate a non-nil UUID")
}

// TestNewNullUUIDWithValue checks if NewNullUUIDWithValue function generates a valid NullUUID.
func TestNewNullUUIDWithValue(t *testing.T) {
	nullID := NewNullUUIDWithValue()
	assert.True(t, nullID.Valid, "NewNullUUIDWithValue should generate a valid NullUUID")
	assert.NotEqual(t, uuid.Nil, nullID.UUID, "NewNullUUIDWithValue should generate a non-nil UUID")
}

// TestNewUUIDAsString checks if NewUUIDAsString function returns a valid UUID string.
func TestNewUUIDAsString(t *testing.T) {
	idStr := NewUUIDAsString()
	assert.NotEmpty(t, idStr, "NewUUIDAsString should return a non-empty string")
	_, err := uuid.FromString(idStr)
	assert.NoError(t, err, "NewUUIDAsString should return a valid UUID string")
}

// TestParseUUID checks if ParseUUID function correctly parses valid and invalid UUID strings.
func TestParseUUID(t *testing.T) {
	validUUID := NewUUID().String()
	parsedValidUUID := ParseUUID(validUUID)
	assert.Equal(t, validUUID, parsedValidUUID.String(), "ParseUUID should correctly parse a valid UUID string")

	invalidUUID := "invalid-uuid-string"
	parsedInvalidUUID := ParseUUID(invalidUUID)
	assert.Equal(t, uuid.Nil, parsedInvalidUUID, "ParseUUID should return uuid.Nil for an invalid UUID string")
}

// TestParseNullUUID checks if ParseNullUUID function correctly parses valid and invalid UUID strings into NullUUID.
func TestParseNullUUID(t *testing.T) {
	validUUID := NewUUID().String()
	parsedValidNullUUID := ParseNullUUID(validUUID)
	assert.True(t, parsedValidNullUUID.Valid, "ParseNullUUID should return a valid NullUUID for a valid UUID string")
	assert.Equal(t, validUUID, parsedValidNullUUID.UUID.String(), "ParseNullUUID should correctly parse a valid UUID string")

	invalidUUID := "invalid-uuid-string"
	parsedInvalidNullUUID := ParseNullUUID(invalidUUID)
	assert.False(t, parsedInvalidNullUUID.Valid, "ParseNullUUID should return an invalid NullUUID for an invalid UUID string")
}

// TestUUIDToString checks if UUIDToString function correctly converts UUID to string.
func TestUUIDToString(t *testing.T) {
	validUUID := NewUUID()
	assert.Equal(t, validUUID.String(), UUIDToString(validUUID), "UUIDToString should convert UUID to its string representation")

	emptyUUID := uuid.Nil
	assert.Empty(t, UUIDToString(emptyUUID), "UUIDToString should return an empty string for uuid.Nil")
}

// TestUUIDToNullUUID checks if UUIDToNullUUID function correctly converts UUID to NullUUID.
func TestUUIDToNullUUID(t *testing.T) {
	validUUID := NewUUID()
	nullUUID := UUIDToNullUUID(validUUID)
	assert.True(t, nullUUID.Valid, "UUIDToNullUUID should return a valid NullUUID for a non-nil UUID")
	assert.Equal(t, validUUID, nullUUID.UUID, "UUIDToNullUUID should correctly convert UUID to NullUUID")

	emptyUUID := uuid.Nil
	nullUUIDEmpty := UUIDToNullUUID(emptyUUID)
	assert.False(t, nullUUIDEmpty.Valid, "UUIDToNullUUID should return an invalid NullUUID for uuid.Nil")
}

// TestUUIDsValidate checks if Validate method correctly validates a slice of UUIDs.
func TestUUIDsValidate(t *testing.T) {
	validUUIDs := UUIDs{NewUUID(), NewUUID()}
	assert.NoError(t, validUUIDs.Validate(), "Validate should not return an error for a slice of non-nil UUIDs")

	invalidUUIDs := UUIDs{uuid.Nil, NewUUID()}
	assert.Error(t, invalidUUIDs.Validate(), "Validate should return an error for a slice containing uuid.Nil")

	emptyUUIDs := UUIDs{}
	assert.NoError(t, emptyUUIDs.Validate(), "Validate should not return an error for an empty slice")
}

// TestUUIDsHas checks if Has method correctly identifies the presence of a UUID in a slice.
func TestUUIDsHas(t *testing.T) {
	targetUUID := NewUUID()
	uuids := UUIDs{targetUUID, NewUUID()}
	assert.True(t, uuids.Has(targetUUID), "Has should return true when the target UUID is in the slice")

	nonTargetUUID := NewUUID()
	assert.False(t, uuids.Has(nonTargetUUID), "Has should return false when the target UUID is not in the slice")
}

// TestUUIDsClean checks if Clean method correctly removes nil and duplicate UUIDs from a slice.
func TestUUIDsClean(t *testing.T) {
	duplicateUUID := NewUUID()
	uuids := UUIDs{duplicateUUID, duplicateUUID, uuid.Nil}
	cleanedUUIDs := uuids.Clean()
	assert.Len(t, cleanedUUIDs, 1, "Clean should return a slice with unique non-nil UUIDs")
	assert.Equal(t, duplicateUUID, cleanedUUIDs[0], "Clean should preserve the non-nil UUID")
}

// TestNullUUIDsValidate checks if Validate method correctly validates a slice of NullUUIDs.
func TestNullUUIDsValidate(t *testing.T) {
	validNullUUIDs := NullUUIDs{NewNullUUIDWithValue(), NewNullUUIDWithValue()}
	assert.NoError(t, validNullUUIDs.Validate(), "Validate should not return an error for a slice of valid NullUUIDs")

	invalidNullUUIDs := NullUUIDs{uuid.NullUUID{UUID: uuid.Nil, Valid: false}, NewNullUUIDWithValue()}
	assert.Error(t, invalidNullUUIDs.Validate(), "Validate should return an error for a slice containing an invalid NullUUID")

	emptyNullUUIDs := NullUUIDs{}
	assert.NoError(t, emptyNullUUIDs.Validate(), "Validate should not return an error for an empty slice")
}

// TestNullUUIDsHas checks if Has method correctly identifies the presence of a NullUUID in a slice.
func TestNullUUIDsHas(t *testing.T) {
	targetUUID := NewUUID()
	nullUUIDs := NullUUIDs{uuid.NullUUID{UUID: targetUUID, Valid: true}, NewNullUUIDWithValue()}
	assert.True(t, nullUUIDs.Has(targetUUID), "Has should return true when the target NullUUID is in the slice")

	nonTargetUUID := NewUUID()
	assert.False(t, nullUUIDs.Has(nonTargetUUID), "Has should return false when the target NullUUID is not in the slice")
}

// TestNullUUIDsClean checks if Clean method correctly removes invalid NullUUIDs from a slice.
func TestNullUUIDsClean(t *testing.T) {
	validNullUUID := NewNullUUIDWithValue()
	invalidNullUUID := uuid.NullUUID{UUID: uuid.Nil, Valid: false}
	nullUUIDs := NullUUIDs{validNullUUID, invalidNullUUID}
	cleanedNullUUIDs := nullUUIDs.Clean()
	assert.Len(t, cleanedNullUUIDs, 1, "Clean should return a slice with only valid NullUUIDs")
	assert.Equal(t, validNullUUID, cleanedNullUUIDs[0], "Clean should preserve the valid NullUUID")
}
