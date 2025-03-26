package auuids

import (
	"encoding/json"
	"errors"
)

type IDLabel struct {
	Label UUIDLabel `json:"label,omitempty" jsonschema:"type=string,title=Label,description=Human-readable label associated with the ID,example=PrimaryEndpoint"`
	Id    UUID      `json:"id,omitempty" jsonschema:"type=string,format=uuid,title=ID,description=Unique identifier (UUID) associated with the label,example=123e4567-e89b-12d3-a456-426614174000"`
}

// NewIDLabel creates a new IDLabel instance
func NewIDLabel(label UUIDLabel, id UUID) *IDLabel {
	if label.IsEmpty() {
		return nil
	}
	if id.IsNil() {
		return nil
	}
	return &IDLabel{
		Label: label,
		Id:    id,
	}
}

func (id *IDLabel) Validate() error {
	if id == nil {
		return errors.New("nil IDLabel")
	}
	if id.Label.IsEmpty() {
		return errors.New("empty Label")
	}
	if id.Id.IsNil() {
		return errors.New("nil ID")
	}
	return nil
}

func (id *IDLabel) HasMatch(other *IDLabel) bool {
	if id == nil || other == nil {
		return false
	}
	return id.Label == other.Label && id.Id == other.Id
}

func (id *IDLabel) HasMatchWithLabelId(label UUIDLabel, uuid UUID) bool {
	if id == nil {
		return false
	}
	if err := id.Validate(); err != nil {
		return false
	}
	return id.Label == label && id.Id == uuid
}

// MustToJson serializes IDLabel to JSON, returning an empty string on failure
func (id *IDLabel) MustToJson() string {
	if id == nil {
		return ""
	}
	data, err := json.Marshal(id)
	if err != nil {
		return ""
	}
	return string(data)
}

type IDLabels []*IDLabel

func (labels IDLabels) Set(target *IDLabel) (IDLabels, error) {
	return labels.SetWithOptions(target, false)
}

func (labels IDLabels) SetByLabelId(label UUIDLabel, id UUID) (IDLabels, error) {
	target := NewIDLabel(label, id)
	return labels.SetWithOptions(target, false)
}

func (labels IDLabels) SetByLabelIdMatchByLabel(label UUIDLabel, id UUID) (IDLabels, error) {
	target := NewIDLabel(label, id)
	return labels.SetWithOptions(target, true)
}

func (labels IDLabels) SetMatchByLabel(target *IDLabel) (IDLabels, error) {
	return labels.SetWithOptions(target, true)
}

// SetWithOptions adds or updates an IDLabel in the IDLabels array with flexible matching options.
func (labels IDLabels) SetWithOptions(target *IDLabel, matchByLabel bool) (IDLabels, error) {
	if target == nil {
		return labels, errors.New("target cannot be nil")
	}

	if err := target.Validate(); err != nil {
		return labels, err
	}

	if labels == nil || len(labels) == 0 {
		return IDLabels{target}, nil
	}

	for i, existingLabel := range labels {
		if existingLabel != nil {
			if matchByLabel && existingLabel.Label == target.Label {
				labels[i] = target
				return labels, nil
			} else if !matchByLabel && existingLabel.HasMatch(target) {
				labels[i] = target
				return labels, nil
			}
		}
	}

	labels = append(labels, target)
	return labels, nil
}

func (labels IDLabels) Remove(target *IDLabel) (IDLabels, error) {
	return labels.RemoveWithOptions(target, false)
}

func (labels IDLabels) RemoveByLabelId(label UUIDLabel, id UUID) (IDLabels, error) {
	target := NewIDLabel(label, id)
	return labels.RemoveWithOptions(target, false)
}

func (labels IDLabels) RemoveByLabelIdMatchByLabel(label UUIDLabel, id UUID) (IDLabels, error) {
	target := NewIDLabel(label, id)
	return labels.RemoveWithOptions(target, true)
}

func (labels IDLabels) RemoveMatchByLabel(target *IDLabel) (IDLabels, error) {
	return labels.RemoveWithOptions(target, true)
}

// RemoveWithOptions removes an IDLabel based on flexible matching criteria.
func (labels IDLabels) RemoveWithOptions(target *IDLabel, matchByLabel bool) (IDLabels, error) {
	if target == nil {
		return labels, errors.New("target cannot be nil")
	}

	if labels == nil || len(labels) == 0 {
		return labels, nil
	}

	for i, existingLabel := range labels {
		if existingLabel != nil {
			if matchByLabel && existingLabel.Label == target.Label {
				return append(labels[:i], labels[i+1:]...), nil
			} else if !matchByLabel && existingLabel.HasMatch(target) {
				return append(labels[:i], labels[i+1:]...), nil
			}
		}
	}

	return labels, nil
}

func (labels IDLabels) Find(id UUID) *IDLabel {
	target := &IDLabel{Id: id}
	return labels.FindWithOptions(target, false)
}

func (labels IDLabels) FindByLabelId(label UUIDLabel, id UUID) *IDLabel {
	target := NewIDLabel(label, id)
	return labels.FindWithOptions(target, false)
}

func (labels IDLabels) FindByLabelIdMatchByLabel(label UUIDLabel) *IDLabel {
	target := &IDLabel{Label: label}
	return labels.FindWithOptions(target, true)
}

func (labels IDLabels) FindMatchByLabel(target *IDLabel) *IDLabel {
	return labels.FindWithOptions(target, true)
}

// FindWithOptions finds an IDLabel based on flexible matching criteria.
func (labels IDLabels) FindWithOptions(target *IDLabel, matchByLabel bool) *IDLabel {
	if labels == nil || target == nil {
		return nil
	}

	for _, existingLabel := range labels {
		if existingLabel != nil {
			if matchByLabel && existingLabel.Label == target.Label {
				return existingLabel
			} else if !matchByLabel && existingLabel.HasMatch(target) {
				return existingLabel
			}
		}
	}
	return nil
}

func (labels IDLabels) HasMatch(targets ...*IDLabel) bool {
	if labels == nil || len(targets) == 0 {
		return false
	}
	for _, target := range targets {
		if labels.HasMatchWithOptions(target, false) {
			return true
		}
	}
	return false
}

func (labels IDLabels) HasMatchByLabelId(label UUIDLabel, id UUID) bool {
	target := NewIDLabel(label, id)
	return labels.HasMatchWithOptions(target, false)
}

func (labels IDLabels) HasMatchByLabelIdByLabel(label UUIDLabel, id UUID) bool {
	target := NewIDLabel(label, id)
	return labels.HasMatchWithOptions(target, true)
}

func (labels IDLabels) HasMatchByLabelOnly(label UUIDLabel) bool {
	target := &IDLabel{Label: label}
	return labels.HasMatchWithOptions(target, true)
}

func (labels IDLabels) HasExactMatch(target *IDLabel) bool {
	return labels.HasMatchWithOptions(target, false)
}

// HasMatchWithOptions checks if the IDLabels array contains a match based on flexible criteria.
func (labels IDLabels) HasMatchWithOptions(target *IDLabel, matchByLabel bool) bool {
	if labels == nil || target == nil {
		return false
	}

	for _, existingLabel := range labels {
		if existingLabel != nil {
			if matchByLabel && existingLabel.Label == target.Label {
				return true
			} else if !matchByLabel && existingLabel.HasMatch(target) {
				return true
			}
		}
	}
	return false
}

// Filter filters the IDLabels collection by a custom predicate.
func (labels IDLabels) Filter(predicate func(label *IDLabel) bool) IDLabels {
	filtered := IDLabels{}
	if labels == nil || len(labels) == 0 || predicate == nil {
		return filtered
	}
	for _, label := range labels {
		if predicate(label) {
			filtered = append(filtered, label)
		}
	}
	return filtered
}

// ToJSON marshals the IDLabels array to JSON
func (labels IDLabels) ToJSON() (string, error) {
	if labels == nil {
		return "", errors.New("IDLabels is nil")
	}

	jsonBytes, err := json.Marshal(labels)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

// FromJSON unmarshals a JSON string into the IDLabels array
func (labels *IDLabels) FromJSON(jsonString string) error {
	if labels == nil {
		return errors.New("IDLabels is nil")
	}

	return json.Unmarshal([]byte(jsonString), labels)
}
