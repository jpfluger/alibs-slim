package atags

import (
	"fmt"
	"strings"
)

// TagKey represents a composite key with a type and an ID, separated by a colon.
type TagKey string

// NewTagKey creates a new TagKey from a type and an ID. If either is empty,
// it returns a TagKey with the non-empty value. If both are provided, it joins them with a colon.
func NewTagKey(kType, kId string) TagKey {
	kType = strings.TrimSpace(kType)
	kId = strings.TrimSpace(kId)

	switch {
	case kType == "":
		return TagKey(kId)
	case kId == "":
		return TagKey(kType)
	default:
		return TagKey(fmt.Sprintf("%s:%s", kType, kId))
	}
}

// GetType extracts the type part from the TagKey.
func (key TagKey) GetType() string {
	parts := strings.SplitN(string(key), ":", 2)
	if len(parts) > 0 {
		return parts[0]
	}
	return ""
}

// GetUniqueId extracts the unique ID part from the TagKey.
func (key TagKey) GetUniqueId() string {
	parts := strings.SplitN(string(key), ":", 2)
	if len(parts) == 2 {
		return parts[1]
	}
	return ""
}

// IsEmpty checks if the TagKey is empty after trimming whitespace.
func (key TagKey) IsEmpty() bool {
	return strings.TrimSpace(string(key)) == ""
}

// TrimSpace returns a new TagKey with leading and trailing whitespace removed.
func (key TagKey) TrimSpace() TagKey {
	return TagKey(strings.ReplaceAll(string(key), " ", ""))
}

// String returns the string representation of the TagKey.
func (key TagKey) String() string {
	return string(key)
}
