package aconns

import (
	"errors"
	"strings"
)

// RecordIdentifierMeta represents a single identifier for a user.
type RecordIdentifierMeta struct {
	Type RecordIdentifierType `json:"type" jsonschema:"type=string,title=Identifier Type,description=The type of identifier,example=Email"`
	Id   string               `json:"id" jsonschema:"type=string,title=Identifier Value,description=The value of the identifier,example=bob@example.com"`
}

// NewRecordIdentifierMeta creates a new RecordIdentifierMeta instance.
func NewRecordIdentifierMeta(recIdType RecordIdentifierType, id string) *RecordIdentifierMeta {
	return &RecordIdentifierMeta{
		Type: recIdType,
		Id:   strings.TrimSpace(id),
	}
}

// Validate validates the RecordIdentifierMeta fields.
func (rid *RecordIdentifierMeta) Validate() error {
	if rid == nil {
		return errors.New("nil RecordIdentifierMeta")
	}
	if rid.Type.IsEmpty() {
		return errors.New("empty Type")
	}
	if strings.TrimSpace(rid.Id) == "" {
		return errors.New("empty ID")
	}
	return nil
}

// HasMatch checks if the current RecordIdentifierMeta matches another.
func (rid *RecordIdentifierMeta) HasMatch(other *RecordIdentifierMeta) bool {
	if rid == nil || other == nil {
		return false
	}
	return rid.Type == other.Type && rid.Id == other.Id
}

// HasMatchWithTypeId checks if the RecordIdentifierMeta matches a given type and ID.
func (rid *RecordIdentifierMeta) HasMatchWithTypeId(recIdType RecordIdentifierType, id string) bool {
	if rid == nil {
		return false
	}
	recIdType = recIdType.TrimSpace()
	id = strings.TrimSpace(id)
	if recIdType.IsEmpty() || id == "" {
		return false
	}
	return rid.Type == recIdType && rid.Id == id
}

// RecordIdentifierMetas represents a list of RecordIdentifierMeta objects.
type RecordIdentifierMetas []*RecordIdentifierMeta

// Set adds or updates a RecordIdentifierMeta without considering type-based or ID-based matching.
func (rids RecordIdentifierMetas) Set(target *RecordIdentifierMeta) (RecordIdentifierMetas, error) {
	return rids.SetWithOptions(target, false, false)
}

// SetById adds or updates a RecordIdentifierMeta by ID only.
func (rids RecordIdentifierMetas) SetById(id string) (RecordIdentifierMetas, error) {
	target := NewRecordIdentifierMeta("", id)
	return rids.SetWithOptions(target, false, true)
}

// SetByTypeAndId adds or updates a RecordIdentifierMeta by type and ID, matching both.
func (rids RecordIdentifierMetas) SetByTypeAndId(recIdType RecordIdentifierType, id string) (RecordIdentifierMetas, error) {
	target := NewRecordIdentifierMeta(recIdType, id)
	return rids.SetWithOptions(target, true, true)
}

// SetByType adds or updates a RecordIdentifierMeta by matching based on type only.
func (rids RecordIdentifierMetas) SetByType(target *RecordIdentifierMeta) (RecordIdentifierMetas, error) {
	return rids.SetWithOptions(target, true, false)
}

// SetWithOptions adds or updates a RecordIdentifierMeta in the RecordIdentifierMetas array with flexible matching options.
func (rids RecordIdentifierMetas) SetWithOptions(target *RecordIdentifierMeta, matchByType bool, matchById bool) (RecordIdentifierMetas, error) {
	if target == nil {
		return rids, errors.New("target cannot be nil")
	}
	if err := target.Validate(); err != nil {
		return rids, err
	}
	for i, existing := range rids {
		if existing != nil {
			// Match by both Type and ID if required
			if matchByType && matchById {
				if existing.Type == target.Type &&
					strings.TrimSpace(existing.Id) == strings.TrimSpace(target.Id) {
					rids[i] = target
					return rids, nil
				}
			} else if matchById {
				// Match by ID only
				if strings.TrimSpace(existing.Id) == strings.TrimSpace(target.Id) {
					rids[i] = target
					return rids, nil
				}
			} else if matchByType {
				// Match by Type only
				if existing.Type == target.Type {
					rids[i] = target
					return rids, nil
				}
			}
		}
	}
	return append(rids, target), nil
}

// RemoveExact removes a RecordIdentifierMeta by exact match (Type and ID).
func (rids RecordIdentifierMetas) RemoveExact(target *RecordIdentifierMeta) (RecordIdentifierMetas, error) {
	return rids.RemoveWithOptions(target, true, true)
}

// RemoveById removes a RecordIdentifierMeta by ID only.
func (rids RecordIdentifierMetas) RemoveById(id string) (RecordIdentifierMetas, error) {
	target := NewRecordIdentifierMeta("", id)
	return rids.RemoveWithOptions(target, false, true)
}

// RemoveByTypeAndId removes a RecordIdentifierMeta by Type and ID.
func (rids RecordIdentifierMetas) RemoveByTypeAndId(recIdType RecordIdentifierType, id string) (RecordIdentifierMetas, error) {
	target := NewRecordIdentifierMeta(recIdType, id)
	return rids.RemoveWithOptions(target, true, true)
}

// RemoveByType removes a RecordIdentifierMeta by Type only.
func (rids RecordIdentifierMetas) RemoveByType(recIdType RecordIdentifierType) (RecordIdentifierMetas, error) {
	target := &RecordIdentifierMeta{Type: recIdType}
	return rids.RemoveWithOptions(target, true, false)
}

// RemoveWithOptions removes a RecordIdentifierMeta with flexible matching options.
func (rids RecordIdentifierMetas) RemoveWithOptions(target *RecordIdentifierMeta, matchByType bool, matchById bool) (RecordIdentifierMetas, error) {
	if target == nil {
		return rids, errors.New("target cannot be nil")
	}
	filtered := RecordIdentifierMetas{}
	for _, existing := range rids {
		if existing == nil {
			continue
		}

		// Match by both Type and ID if required
		if matchByType && matchById {
			if existing.Type == target.Type &&
				strings.TrimSpace(existing.Id) == strings.TrimSpace(target.Id) {
				continue
			}
		} else if matchById {
			// Match by ID only
			if strings.TrimSpace(existing.Id) == strings.TrimSpace(target.Id) {
				continue
			}
		} else if matchByType {
			// Match by Type only
			if existing.Type == target.Type {
				continue
			}
		}

		// Retain the record if no match
		filtered = append(filtered, existing)
	}
	return filtered, nil
}

// FindById finds a RecordIdentifierMeta strictly by its ID.
func (rids RecordIdentifierMetas) FindById(id string) *RecordIdentifierMeta {
	target := NewRecordIdentifierMeta("", id)
	return rids.FindWithMatchingOptions(target, false, true)
}

// FindByTypeAndId finds a RecordIdentifierMeta by its Type and ID.
func (rids RecordIdentifierMetas) FindByTypeAndId(recIdType RecordIdentifierType, id string) *RecordIdentifierMeta {
	target := NewRecordIdentifierMeta(recIdType, id)
	return rids.FindWithMatchingOptions(target, true, true)
}

// FindByType finds a RecordIdentifierMeta by its Type, ignoring ID.
func (rids RecordIdentifierMetas) FindByType(recIdType RecordIdentifierType) *RecordIdentifierMeta {
	target := &RecordIdentifierMeta{Type: recIdType.TrimSpace()}
	return rids.FindWithMatchingOptions(target, true, false)
}

// FindExactMatch finds a RecordIdentifierMeta that matches exactly (Type and ID).
func (rids RecordIdentifierMetas) FindExactMatch(target *RecordIdentifierMeta) *RecordIdentifierMeta {
	return rids.FindWithMatchingOptions(target, true, true)
}

// FindWithMatchingOptions finds a RecordIdentifierMeta using flexible matching options.
func (rids RecordIdentifierMetas) FindWithMatchingOptions(target *RecordIdentifierMeta, matchByType bool, matchById bool) *RecordIdentifierMeta {
	if target == nil {
		return nil
	}
	for _, existing := range rids {
		if existing != nil {
			// Prioritize matching both Type and ID
			if matchById && matchByType {
				if strings.TrimSpace(existing.Id) == strings.TrimSpace(target.Id) &&
					existing.Type == target.Type {
					return existing
				}
			} else if matchById {
				// Match by ID only if requested
				if strings.TrimSpace(existing.Id) == strings.TrimSpace(target.Id) {
					return existing
				}
			} else if matchByType {
				// Match by Type only if requested
				if existing.Type == target.Type {
					return existing
				}
			} else {
				// Fallback to general match if neither matchByType nor matchById
				if existing.HasMatch(target) {
					return existing
				}
			}
		}
	}
	return nil
}

func (rids RecordIdentifierMetas) HasMatch(targets ...*RecordIdentifierMeta) bool {
	if rids == nil || len(targets) == 0 {
		return false
	}
	for _, target := range targets {
		if rids.HasMatchWithOptions(target, false) {
			return true
		}
	}
	return false
}

func (rids RecordIdentifierMetas) HasMatchByTypeId(recIdType RecordIdentifierType, id string) bool {
	target := NewRecordIdentifierMeta(recIdType, id)
	return rids.HasMatchWithOptions(target, false)
}

func (rids RecordIdentifierMetas) HasMatchByTypeIdByType(recIdType RecordIdentifierType, id string) bool {
	target := NewRecordIdentifierMeta(recIdType, id)
	return rids.HasMatchWithOptions(target, true)
}

func (rids RecordIdentifierMetas) HasMatchByTypeOnly(recIdType RecordIdentifierType) bool {
	target := &RecordIdentifierMeta{Type: recIdType}
	return rids.HasMatchWithOptions(target, true)
}

func (rids RecordIdentifierMetas) HasExactMatch(target *RecordIdentifierMeta) bool {
	return rids.HasMatchWithOptions(target, false)
}

// HasMatchWithOptions checks if the RecordIdentifierMetas array contains a match based on flexible criteria.
func (rids RecordIdentifierMetas) HasMatchWithOptions(target *RecordIdentifierMeta, matchByLabel bool) bool {
	if rids == nil || target == nil {
		return false
	}

	for _, existingLabel := range rids {
		if existingLabel != nil {
			if matchByLabel && existingLabel.Type == target.Type {
				return true
			} else if !matchByLabel && existingLabel.HasMatch(target) {
				return true
			}
		}
	}
	return false
}

// Filter filters the RecordIdentifierMetas collection by a custom predicate.
func (rids RecordIdentifierMetas) Filter(predicate func(label *RecordIdentifierMeta) bool) RecordIdentifierMetas {
	filtered := RecordIdentifierMetas{}
	if rids == nil || len(rids) == 0 || predicate == nil {
		return filtered
	}
	for _, label := range rids {
		if predicate(label) {
			filtered = append(filtered, label)
		}
	}
	return filtered
}
