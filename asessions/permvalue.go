package asessions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

// Permission constants as bit flags
const (
	PERM_X = 1 << 4 // Execute
	PERM_L = 1 << 5 // List
	PERM_C = 1 << 0 // Create
	PERM_R = 1 << 1 // Read
	PERM_U = 1 << 2 // Update
	PERM_D = 1 << 3 // Delete
)

// String equivalents
const (
	PERMS_X = "X" // Execute
	PERMS_L = "L" // List
	PERMS_C = "C" // Create
	PERMS_R = "R" // Read
	PERMS_U = "U" // Update
	PERMS_D = "D" // Delete
)

// Map string permissions to their bit values
var permsStringToBit = map[string]int{
	PERMS_X: PERM_X,
	PERMS_L: PERM_L,
	PERMS_C: PERM_C,
	PERMS_R: PERM_R,
	PERMS_U: PERM_U,
	PERMS_D: PERM_D,
}

// Map bit permissions to their string values
var permsBitToString = map[int]string{
	PERM_X: PERMS_X,
	PERM_L: PERMS_L,
	PERM_C: PERMS_C,
	PERM_R: PERMS_R,
	PERM_U: PERMS_U,
	PERM_D: PERMS_D,
}

// Convert a string representation of permissions to its bit equivalent
func PermissionsStringToBit(perms string) int {
	bitValue := 0
	for _, char := range strings.ToUpper(perms) {
		if bit, exists := permsStringToBit[string(char)]; exists {
			bitValue |= bit
		}
	}
	return bitValue
}

// Convert a bit representation of permissions to its string equivalent
func PermissionsBitToString(perms int) string {
	var result []string
	for bit, str := range permsBitToString {
		if perms&bit != 0 {
			result = append(result, str)
		}
	}
	return strings.Join(result, "")
}

// PermValue represents permissions as a single integer
type PermValue struct {
	value int
}

// MustNewPermValue creates a PermValue from a string representation
func MustNewPermValue(value string) *PermValue {
	pv := &PermValue{}
	pv.SetValues(value)
	return pv
}

// Values returns the string representation of the PermValue
func (pv *PermValue) Values() string {
	if pv == nil {
		return ""
	}
	var result []string
	if pv.value&PERM_X != 0 {
		result = append(result, "X")
	}
	if pv.value&PERM_L != 0 {
		result = append(result, "L")
	}
	if pv.value&PERM_C != 0 {
		result = append(result, "C")
	}
	if pv.value&PERM_R != 0 {
		result = append(result, "R")
	}
	if pv.value&PERM_U != 0 {
		result = append(result, "U")
	}
	if pv.value&PERM_D != 0 {
		result = append(result, "D")
	}
	return strings.Join(result, "")
}

// SetValues sets the permission flags based on the input string
func (pv *PermValue) SetValues(values string) {
	values = strings.ToUpper(values)
	pv.value = 0 // Reset all permissions
	if strings.Contains(values, "X") {
		pv.value |= PERM_X
	}
	if strings.Contains(values, "L") {
		pv.value |= PERM_L
	}
	if strings.Contains(values, "C") {
		pv.value |= PERM_C
	}
	if strings.Contains(values, "R") {
		pv.value |= PERM_R
	}
	if strings.Contains(values, "U") {
		pv.value |= PERM_U
	}
	if strings.Contains(values, "D") {
		pv.value |= PERM_D
	}
}

// IsPermValueAllowed checks if all characters in the input string are valid permission flags.
// If any character is not valid, it returns false.
func IsPermValueAllowed(values string) bool {
	values = strings.ToUpper(values) // Ensure case-insensitivity
	validChars := "XLCRUD"           // Define valid permission characters

	for _, char := range values {
		if !strings.ContainsRune(validChars, char) {
			return false // Found an invalid character
		}
	}
	return true // All characters are valid
}

// IsEmptyValue checks if no permissions are set
func (pv *PermValue) IsEmptyValue() bool {
	return pv.value == 0
}

// HasValue checks if at least one permission is set
func (pv *PermValue) HasValue() bool {
	return pv.value != 0
}

// MatchOneByBit checks if the PermValue matches at least one bit in the given bitwise format.
func (pv *PermValue) MatchOneByBit(bits int) bool {
	if pv == nil {
		return false
	}
	return pv.value&bits != 0
}

// SetByBit sets specific permissions based on the given bitwise format.
func (pv *PermValue) SetByBit(bits int) {
	if pv == nil {
		return
	}
	pv.value |= bits
}

// RemoveByBit removes specific permissions based on the given bitwise format.
func (pv *PermValue) RemoveByBit(bits int) {
	if pv == nil {
		return
	}
	pv.value &^= bits // Bitwise AND NOT to clear specific bits
}

// MatchOne checks if the PermValue matches at least one character in the input string
func (pv *PermValue) MatchOne(permChars string) bool {
	if pv == nil || permChars == "" {
		return false
	}
	target := MustNewPermValue(permChars)
	return pv.value&target.value != 0
}

// MatchOneByPerm checks if the PermValue matches at least one permission in the target PermValue
func (pv *PermValue) MatchOneByPerm(target *PermValue) bool {
	if pv == nil || target == nil {
		return false
	}
	return pv.value&target.value != 0
}

// MergePermsByBits merges the permissions with the given bitwise parameter
func (pv *PermValue) MergePermsByBits(bits int) {
	if pv == nil {
		return
	}
	pv.value |= bits
}

// MergePermsByChars merges the permissions with the characters in the input string
func (pv *PermValue) MergePermsByChars(permChars string) {
	if pv == nil {
		return
	}
	target := MustNewPermValue(permChars)
	pv.value |= target.value
}

// HasExcessiveBits checks if the PermValue has excessive permissions compared to the given bitwise parameter
func (pv *PermValue) HasExcessiveBits(bits int) bool {
	if pv == nil {
		return false
	}
	return pv.value&^bits != 0 // Retains bits in pv.value that are not in bits
}

// HasExcessiveChars checks if the PermValue has excessive permissions compared to the input string
func (pv *PermValue) HasExcessiveChars(targetPermChars string) bool {
	if pv == nil {
		return false
	}
	target := MustNewPermValue(targetPermChars)
	return pv.value&^target.value != 0
}

// ReplaceExcessiveBits keeps only the permissions that exist in the given bitwise parameter
func (pv *PermValue) ReplaceExcessiveBits(bits int) {
	if pv == nil {
		return
	}
	pv.value &= bits
}

// ReplaceExcessiveChars keeps only the permissions that exist in the input string
func (pv *PermValue) ReplaceExcessiveChars(permChars string) {
	if pv == nil {
		return
	}
	target := MustNewPermValue(permChars)
	pv.value &= target.value
}

// SubtractPermsByBits removes the permissions present in the given bitwise parameter
func (pv *PermValue) SubtractPermsByBits(bits int) {
	if pv == nil {
		return
	}
	pv.value &^= bits // Clears the bits in pv.value that are set in bits
}

// SubtractPermsByChars removes the permissions present in the input string
func (pv *PermValue) SubtractPermsByChars(permChars string) {
	if pv == nil {
		return
	}
	target := MustNewPermValue(permChars)
	pv.value &^= target.value
}

// Clone creates a copy of the PermValue
func (pv *PermValue) Clone() *PermValue {
	return &PermValue{value: pv.value}
}

// MarshalJSON converts PermValue to a human-readable string for JSON serialization
func (pv *PermValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(pv.Values())
}

// MarshalJSONAsInt converts PermValue to an integer representation for JSON serialization
func (pv *PermValue) MarshalJSONAsInt() ([]byte, error) {
	return json.Marshal(pv.value)
}

// UnmarshalJSON detects the type of input (string, int, or null) and unmarshals accordingly.
func (pv *PermValue) UnmarshalJSON(data []byte) error {
	// Trim leading and trailing whitespace
	trimmed := bytes.TrimSpace(data)

	// Handle null input
	if string(trimmed) == "null" {
		*pv = PermValue{} // Reset to default state
		return nil
	}

	// Check if the first character is a quote (indicating a string)
	if len(trimmed) > 0 && trimmed[0] == '"' {
		var stringValue string
		if err := json.Unmarshal(trimmed, &stringValue); err != nil {
			return fmt.Errorf("failed to unmarshal PermValue as string: %w", err)
		}
		pv.SetValues(stringValue)
		return nil
	}

	// Otherwise, attempt to unmarshal as an integer
	var intValue int
	if err := json.Unmarshal(trimmed, &intValue); err != nil {
		return fmt.Errorf("failed to unmarshal PermValue as int: %w", err)
	}
	pv.value = intValue
	return nil
}

//// Permission values represented as constants.
//const (
//	PERMVALUE_X = "X" // Execute permission
//	PERMVALUE_L = "L" // List permission
//	PERMVALUE_C = "C" // Create permission
//	PERMVALUE_R = "R" // Read permission
//	PERMVALUE_U = "U" // Update permission
//	PERMVALUE_D = "D" // Delete permission
//)
//
//// PermValue struct represents a set of permissions.
//type PermValue struct {
//	X bool // Execute permission flag
//	L bool // List permission flag
//	C bool // Create permission flag
//	R bool // Read permission flag
//	U bool // Update permission flag
//	D bool // Delete permission flag
//}
//
//// MustNewPermValue creates a new PermValue and sets its values based on the input string.
//func MustNewPermValue(value string) *PermValue {
//	pv := &PermValue{}
//	pv.SetValues(value)
//	return pv
//}
//
//// Values returns the string representation of the PermValue.
//func (pv *PermValue) Values() string {
//	if pv == nil {
//		return ""
//	}
//	var permChars []string
//	if pv.X {
//		permChars = append(permChars, PERMVALUE_X)
//	}
//	if pv.L {
//		permChars = append(permChars, PERMVALUE_L)
//	}
//	if pv.C {
//		permChars = append(permChars, PERMVALUE_C)
//	}
//	if pv.R {
//		permChars = append(permChars, PERMVALUE_R)
//	}
//	if pv.U {
//		permChars = append(permChars, PERMVALUE_U)
//	}
//	if pv.D {
//		permChars = append(permChars, PERMVALUE_D)
//	}
//	return strings.Join(permChars, "")
//}
//
//// SetValues sets the values of the PermValue based on the input string.
//func (pv *PermValue) SetValues(values string) {
//	values = strings.ToUpper(values)
//	pv.X = strings.Contains(values, PERMVALUE_X)
//	pv.L = strings.Contains(values, PERMVALUE_L)
//	pv.C = strings.Contains(values, PERMVALUE_C)
//	pv.R = strings.Contains(values, PERMVALUE_R)
//	pv.U = strings.Contains(values, PERMVALUE_U)
//	pv.D = strings.Contains(values, PERMVALUE_D)
//}
//
//// IsEmptyValue checks if all values inside PermValue are false.
//func (pv *PermValue) IsEmptyValue() bool {
//	//return !(pv.X && pv.L && pv.C && pv.R && pv.U && pv.D)
//	return !(pv.X || pv.L || pv.C || pv.R || pv.U || pv.D)
//}
//
//// HasValue checks if at least one values is true
//func (pv *PermValue) HasValue() bool {
//	return pv.X || pv.L || pv.C || pv.R || pv.U || pv.D
//}
//
//// MatchOne checks if the PermValue matches at least one character in the input string.
//func (pv *PermValue) MatchOne(permChars string) bool {
//	if pv == nil || permChars == "" {
//		return false
//	}
//	return (pv.X && strings.Contains(permChars, PERMVALUE_X)) ||
//		(pv.L && strings.Contains(permChars, PERMVALUE_L)) ||
//		(pv.C && strings.Contains(permChars, PERMVALUE_C)) ||
//		(pv.R && strings.Contains(permChars, PERMVALUE_R)) ||
//		(pv.U && strings.Contains(permChars, PERMVALUE_U)) ||
//		(pv.D && strings.Contains(permChars, PERMVALUE_D))
//}
//
//// MatchOneByPerm checks if the PermValue matches at least one permission in the target PermValue.
//func (pv *PermValue) MatchOneByPerm(target *PermValue) bool {
//	if pv == nil || target == nil {
//		return false
//	}
//	return (pv.X && target.X) ||
//		(pv.L && target.L) ||
//		(pv.C && target.C) ||
//		(pv.R && target.R) ||
//		(pv.U && target.U) ||
//		(pv.D && target.D)
//}
//
//// MergePermsByChars merges the permissions in the PermValue with the characters in the input string.
//func (pv *PermValue) MergePermsByChars(permChars string) {
//	if pv == nil {
//		return
//	}
//	permChars = strings.ToUpper(permChars)
//	pv.X = pv.X || strings.Contains(permChars, PERMVALUE_X)
//	pv.L = pv.L || strings.Contains(permChars, PERMVALUE_L)
//	pv.C = pv.C || strings.Contains(permChars, PERMVALUE_C)
//	pv.R = pv.R || strings.Contains(permChars, PERMVALUE_R)
//	pv.U = pv.U || strings.Contains(permChars, PERMVALUE_U)
//	pv.D = pv.D || strings.Contains(permChars, PERMVALUE_D)
//}
//
//// HasExcessiveChars checks if the PermValue has excessive permissions compared to the input string.
//// It returns true if the PermValue contains permissions not present in the target permission characters.
//func (pv *PermValue) HasExcessiveChars(targetPermChars string) bool {
//	if pv == nil {
//		return false
//	}
//	targetPermChars = strings.ToUpper(targetPermChars)
//	return (pv.X && !strings.Contains(targetPermChars, PERMVALUE_X)) ||
//		(pv.L && !strings.Contains(targetPermChars, PERMVALUE_L)) ||
//		(pv.C && !strings.Contains(targetPermChars, PERMVALUE_C)) ||
//		(pv.R && !strings.Contains(targetPermChars, PERMVALUE_R)) ||
//		(pv.U && !strings.Contains(targetPermChars, PERMVALUE_U)) ||
//		(pv.D && !strings.Contains(targetPermChars, PERMVALUE_D))
//}
//
//// ReplaceExcessiveChars replaces the permissions in the PermValue that are excessive compared to the input string.
//// It sets the permission flags to false if they are not present in the target permission characters.
//func (pv *PermValue) ReplaceExcessiveChars(permChars string) {
//	if pv == nil {
//		return
//	}
//	permChars = strings.ToUpper(permChars)
//	pv.X = pv.X && strings.Contains(permChars, PERMVALUE_X)
//	pv.L = pv.L && strings.Contains(permChars, PERMVALUE_L)
//	pv.C = pv.C && strings.Contains(permChars, PERMVALUE_C)
//	pv.R = pv.R && strings.Contains(permChars, PERMVALUE_R)
//	pv.U = pv.U && strings.Contains(permChars, PERMVALUE_U)
//	pv.D = pv.D && strings.Contains(permChars, PERMVALUE_D)
//}
//
//// SubtractPermsByChars removes the permissions in the PermValue that are present in the input string.
//func (pv *PermValue) SubtractPermsByChars(permChars string) {
//	if pv == nil {
//		return
//	}
//	permChars = strings.ToUpper(permChars)
//	pv.X = pv.X && !strings.Contains(permChars, PERMVALUE_X)
//	pv.L = pv.L && !strings.Contains(permChars, PERMVALUE_L)
//	pv.C = pv.C && !strings.Contains(permChars, PERMVALUE_C)
//	pv.R = pv.R && !strings.Contains(permChars, PERMVALUE_R)
//	pv.U = pv.U && !strings.Contains(permChars, PERMVALUE_U)
//	pv.D = pv.D && !strings.Contains(permChars, PERMVALUE_D)
//}
//
//// Clone creates a copy of the PermValue.
//func (pv *PermValue) Clone() *PermValue {
//	if pv == nil {
//		return &PermValue{}
//	}
//	return &PermValue{
//		X: pv.X,
//		L: pv.L,
//		C: pv.C,
//		R: pv.R,
//		U: pv.U,
//		D: pv.D,
//	}
//}
