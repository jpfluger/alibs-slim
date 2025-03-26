package asessions

import (
	"encoding/json"
	"sort"
	"strings"
)

// PermSet is a map of string to Perm pointers
// The key must be lower-case.
type PermSet map[string]*Perm

// NewPermSetByPair creates a new PermSet containing a single Perm constructed from the given key and value.
func NewPermSetByPair(key string, value string) PermSet {
	ps := PermSet{}
	ps.SetPerm(NewPermByPair(key, value))
	return ps
}

// NewPermSetByString creates a new PermSet from the provided string slice
func NewPermSetByString(perms []string) PermSet {
	ps := PermSet{}

	for _, keyValue := range perms {
		ps.SetPerm(NewPerm(keyValue))
	}

	return ps
}

// NewPermSetByBits creates a new PermSet containing a single Perm with the given key and bitwise value.
func NewPermSetByBits(key string, bits int) PermSet {
	ps := PermSet{}

	key = strings.ToLower(strings.TrimSpace(key))
	if key == "" || bits == 0 {
		return ps // Return an empty PermSet for invalid key or bits
	}

	// Create a new Perm and add it to the PermSet
	perm := &Perm{key: key, value: &PermValue{value: bits}}
	ps[key] = perm

	return ps
}

// SetPermSetByBits creates or updates a PermSet with a single Perm constructed from the given key and bitwise value.
func SetPermSetByBits(ps PermSet, key string, bits int) PermSet {
	if ps == nil {
		ps = PermSet{}
	}

	key = strings.ToLower(strings.TrimSpace(key))
	if key == "" || bits == 0 {
		return ps // Ignore invalid key or empty bits
	}

	perm := ps[key]
	if perm == nil {
		// Create a new Perm if it doesn't exist
		perm = &Perm{key: key, value: &PermValue{value: bits}}
		ps[key] = perm
	} else {
		// Update the existing Perm with the new bits
		perm.value.SetByBit(bits)
	}

	return ps
}

// Validate ensures no nil values associated with a key
func (ps PermSet) Validate() {
	for key, perm := range ps {
		if perm == nil {
			delete(ps, key) // Remove nil entries
		}
	}
}

// SetPerm adds or updates a Perm in the PermSet.
func (ps PermSet) SetPerm(target *Perm) {
	if target == nil || !target.IsValid() {
		return
	}

	// Add or replace the permission in the PermSet
	ps[target.Key()] = target
}

// MergePerm merges a target Perm into the PermSet
func (ps PermSet) MergePerm(target *Perm) {
	if target == nil || !target.IsValid() {
		return
	}
	if existing, ok := ps[target.Key()]; ok {
		existing.value.MergePermsByBits(target.value.value)
	} else {
		ps.SetPerm(target.Clone())
	}
}

// MergeByPermSet merges another PermSet into this PermSet
func (ps PermSet) MergeByPermSet(target PermSet) {
	for _, perm := range target {
		ps.MergePerm(perm)
	}
}

// SubtractPerm subtracts a target Perm from the PermSet
func (ps PermSet) SubtractPerm(target *Perm) {
	if target == nil || !target.IsValid() {
		return
	}
	if existing, ok := ps[target.Key()]; ok {
		existing.value.SubtractPermsByBits(target.value.value)
		if existing.value.value == 0 {
			delete(ps, target.Key())
		}
	}
}

// Clone creates a deep copy of the PermSet.
func (ps PermSet) Clone() PermSet {
	cloned := PermSet{}
	if ps == nil || len(ps) == 0 {
		return cloned
	}
	for key, perm := range ps {
		if perm != nil {
			cloned[key] = perm.Clone() // Deep clone each Perm
		}
	}
	return cloned
}

// MatchesPerm checks if a PermSet contains a matching Perm
func (ps PermSet) MatchesPerm(target *Perm) bool {
	if target == nil || !target.IsValid() {
		return false
	}
	if existing, ok := ps[target.Key()]; ok {
		return existing.value.MatchOneByBit(target.value.value)
	}
	return false
}

// HasPerm checks if the PermSet has a specific permission.
func (ps PermSet) HasPerm(target Perm) bool {
	return ps.MatchesPerm(&target)
}

// HasPermS checks if the PermSet has a specific permission.
func (ps PermSet) HasPermS(keyPermValue string) bool {
	return ps.MatchesPerm(NewPerm(keyPermValue))
}

// HasPermSV checks if the PermSet has a specific permission value for a given key.
func (ps PermSet) HasPermSV(key string, permValue string) bool {
	return ps.MatchesPerm(NewPermByPair(key, permValue))
}

// HasPermB checks if the PermSet has a specific permission value for a given key.
func (ps PermSet) HasPermB(keyBits string) bool {
	return ps.MatchesPerm(NewPerm(keyBits))
}

// HasPermBV checks if the PermSet has a specific permission value for a given key.
func (ps PermSet) HasPermBV(key string, bit int) bool {
	return ps.MatchesPerm(NewPermByBitValue(key, bit))
}

// HasPermSet checks if the PermSet has a specific permission value for the target.
func (ps PermSet) HasPermSet(target PermSet) bool {
	if target == nil || len(target) == 0 {
		return false
	}
	for _, perm := range target {
		if existing, ok := ps[perm.Key()]; ok {
			return existing.value.MatchOneByBit(perm.value.value)
		}
	}
	return false
}

// ToStringArray converts a PermSet to a string array for serialization
func (ps PermSet) ToStringArray() []string {
	arr := make([]string, 0, len(ps))
	for _, perm := range ps {
		arr = append(arr, perm.Single())
	}
	sort.Strings(arr) // Ensure consistent order
	return arr
}

// FromStringArray populates a PermSet from a string array
func FromStringArray(perms []string) PermSet {
	ps := PermSet{}
	for _, keyValue := range perms {
		ps.SetPerm(NewPerm(keyValue))
	}
	return ps
}

// SubtractByPermSet subtracts another PermSet from this PermSet
func (ps PermSet) SubtractByPermSet(target PermSet) {
	for _, perm := range target {
		ps.SubtractPerm(perm)
	}
}

// IsSubsetOf checks if all permissions in the current PermSet are within the permissions of the target PermSet.
// Returns false if any permissions in the target PermSet exceed those in the ps PermSet.
func (ps PermSet) IsSubsetOf(target PermSet) bool {
	if ps == nil || target == nil {
		return false // A nil PermSet cannot be a subset
	}

	for key, perm := range ps {
		targetPerm, exists := target[key]
		if !exists {
			// If the target does not have this key, ps cannot be a subset
			return false
		}

		if perm != nil && targetPerm != nil {
			// Check if the current Perm exceeds the target Perm
			//			fmt.Println(perm.Single(), targetPerm.Single())
			if perm.value.value&^targetPerm.value.value != 0 {
				return false
			}
		} else if perm != nil && targetPerm == nil {
			// Current Perm has a value, but the target Perm is nil
			return false
		}
	}

	return true // All permissions in ps are within the corresponding permissions in target
}

// ReplaceExcessivePermSet ensures the permissions in the current PermSet do not exceed those in the master PermSet.
// It modifies the current PermSet in place to restrict any permissions that exceed the corresponding permissions in master.
func (ps PermSet) ReplaceExcessivePermSet(master PermSet) {
	if ps == nil || master == nil {
		return // Nothing to do if either PermSet is nil
	}

	for key, perm := range ps {
		// Get the corresponding permission from the master PermSet
		masterPerm, exists := master[key]
		if !exists || masterPerm == nil || masterPerm.value == nil {
			// If there is no corresponding master permission, remove the current permission entirely
			delete(ps, key)
			continue
		}

		// Ensure the current permission does not exceed the master permission
		if perm != nil {
			perm.ReplaceExcessiveBits(masterPerm.value.value)
			// If the permission becomes empty after replacement, remove it from the PermSet
			if perm.value.IsEmptyValue() {
				delete(ps, key)
			}
		}
	}
}

// MarshalJSON returns the ps in json format in bytes.
func (ps PermSet) MarshalJSON() ([]byte, error) {
	arr := []string{}

	for _, perm := range ps {
		arr = append(arr, perm.Single())
	}

	return json.Marshal(arr)
}

// UnmarshalJSON receives a json byte stream and converts it to a PermSet.
func (ps *PermSet) UnmarshalJSON(b []byte) error {
	if *ps == nil {
		*ps = PermSet{}
	}

	arr := &[]string{}

	if err := json.Unmarshal(b, arr); err != nil {
		return err
	}

	for _, keyValue := range *arr {
		ps.SetPerm(NewPerm(keyValue))
	}

	return nil
}

// MarshalAsInt converts the PermSet to JSON format with Perm values serialized as integers in an array of strings.
func (ps PermSet) MarshalAsInt() ([]byte, error) {
	arr := []string{}

	for _, perm := range ps {
		if perm == nil || perm.value == nil {
			// Skip nil entries to maintain consistency with MarshalJSON
			continue
		}

		// Add "key:intValue" format to the array
		arr = append(arr, perm.SingleAsInt())
	}

	// Sort the array to ensure consistent output
	sort.Strings(arr)

	// Marshal the array as JSON
	return json.Marshal(arr)
}

// GobEncode encodes a PermSet using the compact representation from MarshalAsInt.
func (ps PermSet) GobEncode() ([]byte, error) {
	// Use MarshalAsInt to generate a compact JSON representation
	return ps.MarshalAsInt()
}

// GobDecode decodes a PermSet from a byte slice using UnmarshalJSON.
func (ps *PermSet) GobDecode(data []byte) error {
	// Use UnmarshalJSON to parse the data into the PermSet
	return ps.UnmarshalJSON(data)
}
