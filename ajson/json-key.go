package ajson

import (
	"fmt"
	"strings"
)

// JsonKey represents a key in a JSON object, which may include a path separated by dots.
type JsonKey string

// JsonKeys is a slice of JsonKey, representing multiple keys in a JSON object.
type JsonKeys []JsonKey

// IsEmpty checks if the JsonKey is empty after trimming whitespace.
func (jk JsonKey) IsEmpty() bool {
	return strings.TrimSpace(string(jk)) == ""
}

// TrimSpace trims whitespace from the JsonKey and returns a new JsonKey.
func (jk JsonKey) TrimSpace() JsonKey {
	return JsonKey(strings.TrimSpace(string(jk)))
}

// String returns the JsonKey as a trimmed string.
func (jk JsonKey) String() string {
	return strings.TrimSpace(string(jk))
}

// IsRoot checks if the JsonKey is a root key (does not contain any dots).
func (jk JsonKey) IsRoot() bool {
	return !strings.Contains(jk.String(), ".")
}

// GetRoot extracts the root part of the JsonKey path.
func (jk JsonKey) GetRoot() JsonKey {
	parts := jk.GetPathParts()
	if len(parts) == 0 {
		return ""
	}
	return JsonKey(parts[0])
}

// GetPathLeaf extracts the last part of the JsonKey path.
func (jk JsonKey) GetPathLeaf() JsonKey {
	if jk.IsRoot() {
		return jk.TrimSpace()
	}
	parts := jk.GetPathParts()
	return JsonKey(parts[len(parts)-1])
}

// GetPathParts splits the JsonKey path into its constituent parts.
func (jk JsonKey) GetPathParts() []string {
	return strings.Split(jk.String(), ".")
}

// GetPathParent extracts the parent path of the JsonKey.
func (jk JsonKey) GetPathParent() JsonKey {
	if jk.IsRoot() {
		return ""
	}
	return JsonKey(jk.String()[:strings.LastIndex(jk.String(), ".")])
}

// Add appends a target JsonKey to the current JsonKey path.
func (jk *JsonKey) Add(target JsonKey) JsonKey {
	if target.IsEmpty() {
		return *jk
	}
	if jk.IsEmpty() {
		*jk = target
	} else {
		*jk = JsonKey(fmt.Sprintf("%s.%s", jk.String(), target.String()))
	}
	return *jk
}

// CopyPlusAdd creates a new JsonKey by appending a target JsonKey to the current JsonKey path.
func (jk JsonKey) CopyPlusAdd(target JsonKey) JsonKey {
	if target.IsEmpty() {
		return jk
	}
	if jk.IsEmpty() {
		return target
	}
	return JsonKey(fmt.Sprintf("%s.%s", jk.String(), target.String()))
}

// CopyPlusAddInt creates a new JsonKey by appending an integer to the current JsonKey path.
func (jk JsonKey) CopyPlusAddInt(target int) JsonKey {
	if target < 0 {
		return jk
	}
	if jk.IsEmpty() {
		return JsonKey(fmt.Sprintf("%d", target))
	}
	return JsonKey(fmt.Sprintf("%s.%d", jk.String(), target))
}
