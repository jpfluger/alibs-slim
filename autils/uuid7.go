package autils

import (
	"fmt"
	"github.com/gofrs/uuid/v5"
)

// Package autils provides utility functions for working with UUIDs.
// It uses the gofrs/uuid library to generate and parse UUIDs.
// Reference: https://github.com/gofrs/uuid

// NewUUID generates a new UUID version 7 and returns it.
func NewUUID() uuid.UUID {
	u7, _ := uuid.NewV7()
	return u7
}

// NewNullUUIDWithValue creates a NullUUID with a new UUID version 7.
func NewNullUUIDWithValue() uuid.NullUUID {
	return uuid.NullUUID{UUID: NewUUID(), Valid: true}
}

// NewUUIDAsString returns the string representation of a new UUID version 7.
func NewUUIDAsString() string {
	return NewUUID().String()
}

// ParseUUID parses a UUID from the given string.
func ParseUUID(target string) uuid.UUID {
	return uuid.FromStringOrNil(target)
}

// ParseNullUUID parses a NullUUID from the given string.
func ParseNullUUID(target string) uuid.NullUUID {
	u7 := ParseUUID(target)
	return uuid.NullUUID{UUID: u7, Valid: u7 != uuid.Nil}
}

// UUIDToString returns the string representation of the given UUID.
// If the UUID is nil, it returns an empty string.
func UUIDToString(target uuid.UUID) string {
	if target == uuid.Nil {
		return ""
	}
	return target.String()
}

// UUIDToStringEmpty is an alias for UUIDToString.
func UUIDToStringEmpty(target uuid.UUID) string {
	return UUIDToString(target)
}

// UUIDToNullUUID converts a UUID to a NullUUID.
func UUIDToNullUUID(target uuid.UUID) uuid.NullUUID {
	return uuid.NullUUID{
		UUID:  target,
		Valid: target != uuid.Nil,
	}
}

// UUIDs is a slice of UUIDs.
type UUIDs []uuid.UUID

// Validate checks if all UUIDs in the slice are non-nil.
func (ids UUIDs) Validate() error {
	return ids.ValidateWithOptions(false)
}

// ValidateWithOptions checks if all UUIDs in the slice are non-nil.
// If mustHaveCount is true, it also checks that the slice is not empty.
func (ids UUIDs) ValidateWithOptions(mustHaveCount bool) error {
	if len(ids) == 0 {
		if mustHaveCount {
			return fmt.Errorf("uids is empty")
		}
		return nil
	}
	for ii, uid := range ids {
		if uid == uuid.Nil {
			return fmt.Errorf("uid nil at index %d", ii)
		}
	}
	return nil
}

// Has checks if the slice contains the target UUID.
func (ids UUIDs) Has(target uuid.UUID) bool {
	for _, id := range ids {
		if id == target {
			return true
		}
	}
	return false
}

// Clean returns a new slice with duplicate and nil UUIDs removed.
func (ids UUIDs) Clean() UUIDs {
	var arr UUIDs
	seen := make(map[uuid.UUID]bool)
	for _, id := range ids {
		if id == uuid.Nil || seen[id] {
			continue
		}
		seen[id] = true
		arr = append(arr, id)
	}
	return arr
}

// NullUUIDs is a slice of NullUUIDs.
type NullUUIDs []uuid.NullUUID

// Validate checks if all NullUUIDs in the slice are non-nil.
func (ids NullUUIDs) Validate() error {
	return ids.ValidateWithOptions(false)
}

// ValidateWithOptions checks if all NullUUIDs in the slice are non-nil.
// If mustHaveCount is true, it also checks that the slice is not empty.
func (ids NullUUIDs) ValidateWithOptions(mustHaveCount bool) error {
	if len(ids) == 0 {
		if mustHaveCount {
			return fmt.Errorf("uids is empty")
		}
		return nil
	}
	for ii, uid := range ids {
		if uid.UUID == uuid.Nil {
			return fmt.Errorf("uid nil at index %d", ii)
		}
	}
	return nil
}

// Has checks if the slice contains the target NullUUID.
func (ids NullUUIDs) Has(target uuid.UUID) bool {
	for _, id := range ids {
		if id.UUID == target {
			return true
		}
	}
	return false
}

// Clean returns a new slice with nil NullUUIDs removed.
func (ids NullUUIDs) Clean() NullUUIDs {
	var arr NullUUIDs
	for _, id := range ids {
		if id.UUID == uuid.Nil {
			continue
		}
		arr = append(arr, id)
	}
	return arr
}
