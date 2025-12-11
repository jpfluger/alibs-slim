package asessions

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jpfluger/alibs-slim/autils"
)

// IPerm defines an interface for permission handling.
type IPerm interface {
	IsValid() bool
	Value() string
	Key() string
	LastChange() *time.Time
	SetValue(string)
	SetKey(string)
	SetLastChange(*time.Time)
	SetKeyValue(string, string)
	Single() string
	CanCreate() bool
	CanRead() bool
	CanUpdate() bool
	CanDelete() bool
	CanExecute() bool
	MatchOne(string) bool
	MergePermsByChars(string)
	SubtractPermsByChars(string)

	MatchOneByBit(bit int) bool
	MergePermsByBits(bits int)
	HasExcessiveBits(bits int) bool
	ReplaceExcessiveBits(bits int)
	SubtractPermsByBits(bits int)
}

// Perm represents a permission with a key, value, and optional category.
type Perm struct {
	key      string     // Key is the name of the key to which this value attaches.
	value    *PermValue // Value holds the bitwise XCRUD permissions.
	category string     // Category is optional and may be used for grouping permissions.
}

// NewPerm creates a new Perm from a colon-separated key-value string.
// It parses the input using ParsePerm and returns the Perm if successful.
// If parsing fails (e.g., invalid format or permission values), it returns nil and the error.
func NewPerm(keyValue string) (*Perm, error) {
	perm, err := ParsePerm(keyValue)
	if err != nil {
		return nil, err
	}
	return perm, nil
}

// NewPermSetByString creates a new PermSet from a slice of colon-separated key-value strings.
// It attempts to parse each string using NewPerm and adds valid permissions to the set.
// If any string is invalid, it returns nil and an error wrapping the first failure.
// Successful calls return the populated PermSet with no error.
func NewPermSetByString(perms []string) (PermSet, error) {
	ps := PermSet{}
	for _, keyValue := range perms {
		perm, err := NewPerm(keyValue)
		if err != nil {
			return nil, fmt.Errorf("invalid perm string %q: %w", keyValue, err)
		}
		ps.SetPerm(perm)
	}
	return ps, nil
}

// MustNewPerm creates a new Perm from a colon-separated key-value string.
func MustNewPerm(keyValue string) *Perm {
	perm, err := ParsePerm(keyValue)
	if err != nil {
		// Handle invalid input gracefully
		return MustNewPermByPair("", "")
	}
	return perm
}

// MustNewPermByBitValue creates a new Perm from a key and bit pair.
func MustNewPermByBitValue(key string, bit int) *Perm {
	key = strings.ToLower(strings.TrimSpace(key))
	if bit < 0 {
		return MustNewPermByPair(key, "")
	}

	category, keyClean := autils.ExtractPrefixBrackets(key)
	return &Perm{key: keyClean, value: &PermValue{bit}, category: category}
}

// MustNewPermByPair creates a new Perm from a key and value pair.
func MustNewPermByPair(key string, value string) *Perm {
	key = strings.ToLower(strings.TrimSpace(key))
	value = strings.ToUpper(strings.TrimSpace(value))

	pv := MustNewPermValue(value)

	category, keyClean := autils.ExtractPrefixBrackets(key)
	return &Perm{key: keyClean, value: pv, category: category}
}

// MustNewPermByPairCategory creates a new Perm with a specified category.
func MustNewPermByPairCategory(key string, value string, category string) *Perm {
	key = strings.ToLower(strings.TrimSpace(key))
	value = strings.ToUpper(strings.TrimSpace(value))
	category = strings.TrimSpace(category)

	pv := &PermValue{}
	pv.SetValues(value)

	return &Perm{key: key, value: pv, category: category}
}

// IsValid checks if the Perm is valid by ensuring it has a non-empty key and value.
func (perm *Perm) IsValid() bool {
	return !(perm.value == nil || strings.TrimSpace(perm.key) == "")
}

// IsValueEmpty checks if the PermValue is empty.
func (perm *Perm) IsValueEmpty() bool {
	return perm.value == nil || perm.value.IsEmptyValue()
}

// RawValue returns the raw PermValue.
func (perm *Perm) RawValue() *PermValue {
	return perm.value
}

// Value returns the string representation of the PermValue.
func (perm *Perm) Value() string {
	return perm.value.Values()
}

// SetValue sets the PermValue based on the input string.
func (perm *Perm) SetValue(value string) {
	perm.value.SetValues(value)
}

// Key returns the key of the Perm.
func (perm *Perm) Key() string {
	return perm.key
}

// SetKey sets the key of the Perm.
func (perm *Perm) SetKey(key string) {
	perm.key = key
}

// SetKeyValue sets both the key and value of the Perm.
func (perm *Perm) SetKeyValue(key string, value string) {
	perm.value.SetValues(value)
	perm.key = key
}

// Category returns the category of the Perm.
func (perm *Perm) Category() string {
	return perm.category
}

// SetCategory sets the category of the Perm.
func (perm *Perm) SetCategory(category string) {
	perm.category = category
}

// Single returns a colon-separated string of the key and value.
func (perm *Perm) Single() string {
	return fmt.Sprintf("%s:%s", perm.key, perm.Value())
}

// SingleAsInt returns a colon-separated string of the key and value as an integer.
func (perm *Perm) SingleAsInt() string {
	if perm.value == nil {
		return fmt.Sprintf("%s:0", perm.key)
	}
	return fmt.Sprintf("%s:%d", perm.key, perm.value.value)
}

// SingleWithCategory returns a string representation of the Perm with the category included.
func (perm *Perm) SingleWithCategory() string {
	if perm.category == "" {
		return fmt.Sprintf("%s:%s", perm.key, perm.Value())
	}
	return fmt.Sprintf("[%s]%s:%s", perm.category, perm.key, perm.Value())
}

// CanCreate checks if the create permission is set.
func (perm *Perm) CanCreate() bool {
	return perm.value.value&PERM_C != 0
}

// CanRead checks if the read permission is set.
func (perm *Perm) CanRead() bool {
	return perm.value.value&PERM_R != 0
}

// CanUpdate checks if the update permission is set.
func (perm *Perm) CanUpdate() bool {
	return perm.value.value&PERM_U != 0
}

// CanDelete checks if the delete permission is set.
func (perm *Perm) CanDelete() bool {
	return perm.value.value&PERM_D != 0
}

// CanExecute checks if the execute permission is set.
func (perm *Perm) CanExecute() bool {
	return perm.value.value&PERM_X != 0
}

// CanList checks if the list permission is set.
func (perm *Perm) CanList() bool {
	return perm.value.value&PERM_L != 0
}

// SetByBit sets specific permissions on the Perm based on the given bitwise format.
func (perm *Perm) SetByBit(bits int) {
	if perm == nil || perm.value == nil {
		return
	}
	perm.value.SetByBit(bits)
}

// RemoveByBit removes specific permissions from the Perm based on the given bitwise format.
func (perm *Perm) RemoveByBit(bits int) {
	if perm == nil || perm.value == nil {
		return
	}
	perm.value.RemoveByBit(bits)
}

// MatchOneByBit checks if the Perm matches at least one bit in the given bitwise format.
func (perm *Perm) MatchOneByBit(bits int) bool {
	if perm == nil || perm.value == nil {
		return false
	}
	return perm.value.MatchOneByBit(bits)
}

// MatchOne checks if at least one permission matches the input string.
func (perm *Perm) MatchOne(permChars string) bool {
	return perm.value.MatchOne(permChars)
}

// MatchOneByPerm checks if at least one permission matches another Perm.
func (perm *Perm) MatchOneByPerm(target *Perm) bool {
	if target == nil {
		return false
	}
	return perm.value.MatchOneByPerm(target.value)
}

func (perm *Perm) MergePermsByBits(bits int) {
	if perm == nil || perm.value == nil {
		return
	}
	perm.value.MergePermsByBits(bits)
}

// MergePermsByChars merges permissions with those represented in the input string.
func (perm *Perm) MergePermsByChars(permChars string) {
	perm.value.MergePermsByChars(permChars)
}

func (perm *Perm) SubtractPermsByBits(bits int) {
	if perm == nil || perm.value == nil {
		return
	}
	perm.value.SubtractPermsByBits(bits)
}

func (perm *Perm) HasExcessivePerm(target *Perm) bool {
	var val int
	if target != nil || target.value != nil {
		val = perm.value.value
	}
	return perm.value.HasExcessiveBits(val)
}

// SubtractPermsByChars removes permissions that are present in the input string.
func (perm *Perm) SubtractPermsByChars(permChars string) {
	perm.value.SubtractPermsByChars(permChars)
}

func (perm *Perm) HasExcessiveBits(bits int) bool {
	if perm == nil || perm.value == nil {
		return false
	}
	return perm.value.HasExcessiveBits(bits)
}

// HasExcessivePermsByChars checks for permissions not present in the input string.
func (perm *Perm) HasExcessivePermsByChars(permChars string) bool {
	return perm.value.HasExcessiveChars(permChars)
}

func (perm *Perm) ReplaceExcessiveBits(bits int) {
	if perm == nil || perm.value == nil {
		return
	}
	perm.value.ReplaceExcessiveBits(bits)
}

// ReplaceExcessivePermsByChars replaces excessive permissions with those in the input string.
func (perm *Perm) ReplaceExcessivePermsByChars(permChars string) {
	perm.value.ReplaceExcessiveChars(permChars)
}

// Clone creates a deep copy of the Perm.
func (perm *Perm) Clone() *Perm {
	if perm == nil {
		return nil
	}

	// Deep copy of PermValue if it exists
	var clonedValue *PermValue
	if perm.value != nil {
		clonedValue = &PermValue{
			value: perm.value.value,
		}
	}

	return &Perm{
		key:      perm.key,      // Direct string copy (strings are immutable in Go)
		value:    clonedValue,   // Use the deep-copied value
		category: perm.category, // Direct string copy
	}
}

// MarshalJSON converts the Perm to JSON format.
func (perm Perm) MarshalJSON() ([]byte, error) {
	var s string
	if perm.IsValid() {
		s = perm.Single()
	}
	return json.Marshal(s)
}

// MarshalJSONAsInt converts the Perm to JSON format with the value as an integer.
func (perm Perm) MarshalJSONAsInt() ([]byte, error) {
	var s string
	if perm.IsValid() {
		s = perm.SingleAsInt()
	}
	return json.Marshal(s)
}

func (p *Perm) GobEncode() ([]byte, error) {
	// Prepare a serializable representation of Perm
	data := map[string]interface{}{
		"key":      p.key,
		"value":    p.value, // Ensure PermValue is also serializable
		"category": p.category,
	}

	// Serialize using JSON as an intermediate format
	return json.Marshal(data)
}

func (p *Perm) GobDecode(data []byte) error {
	// Deserialize using JSON as an intermediate format
	temp := map[string]interface{}{}
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	// Extract fields and populate the Perm struct
	p.key = temp["key"].(string)
	p.category = temp["category"].(string)

	// Handle PermValue deserialization if necessary
	if temp["value"] != nil {
		p.value = &PermValue{} // Initialize the PermValue
		// Deserialize the PermValue from its representation (if needed)
	}

	return nil
}

// ParsePerm parses a key-value string or a bit-based representation to create a new Perm.
func ParsePerm(input string) (*Perm, error) {
	input = strings.TrimSpace(input)
	if input == "" || input == "null" {
		return nil, fmt.Errorf("invalid format for Perm: %q", input)
	}

	// Split the input into key and value parts
	parts := strings.SplitN(input, ":", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format for Perm: %q", input)
	}

	key := strings.TrimSpace(parts[0])
	if key == "" {
		return nil, fmt.Errorf("key missing for Perm: %q", input)
	}

	valuePart := strings.TrimSpace(parts[1])

	// Handle "null" or empty value
	if valuePart == "null" || valuePart == "" {
		return &Perm{
			key:   key,
			value: &PermValue{},
		}, nil
	}

	// Attempt to parse the value as an integer
	if intValue, err := strconv.Atoi(valuePart); err == nil {
		return &Perm{
			key: key,
			value: &PermValue{
				value: intValue,
			},
		}, nil
	}

	// If not an integer, validate and parse as a string representation
	if !IsPermValueAllowed(valuePart) {
		return nil, fmt.Errorf("invalid value for Perm: %q", valuePart)
	}

	perm := &Perm{
		key:   key,
		value: &PermValue{},
	}
	perm.value.SetValues(valuePart)

	// Ensure the parsed value is valid (non-empty)
	if perm.value.IsEmptyValue() {
		return nil, fmt.Errorf("invalid value '%s' for Perm: %q", valuePart, input)
	}

	return perm, nil
}

// UnmarshalJSON converts a JSON byte stream to a Perm, supporting both string and int representations for value.
func (perm *Perm) UnmarshalJSON(b []byte) error {
	if perm == nil {
		return fmt.Errorf("UnmarshalJSON called on nil Perm")
	}

	input := strings.Trim(strings.TrimSpace(string(b)), `"`)
	parsedPerm, err := ParsePerm(input)
	if err != nil {
		return err
	}

	*perm = *parsedPerm
	return nil
}
