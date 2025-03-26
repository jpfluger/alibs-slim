package asessions

import (
	"fmt"
	"github.com/jpfluger/alibs-slim/auser"
)

// UIDsPermSet associates a group of UUIDs with a set of permissions (PermSet).
type UIDsPermSet struct {
	UIDs  auser.UIDs `json:"uids,omitempty"`  // Slice of UUIDs.
	Perms PermSet    `json:"perms,omitempty"` // Associated permissions.
}

// Validate checks the integrity of the UIDsPermSet.
func (ups *UIDsPermSet) Validate() error {
	if ups == nil {
		return fmt.Errorf("UIDsPermSet is nil")
	}
	if err := ups.UIDs.ValidateWithOptions(true); err != nil {
		return err
	}
	if ups.Perms == nil || len(ups.Perms) == 0 {
		return fmt.Errorf("permissions in UIDsPermSet are empty")
	}
	return nil
}

// GetUIDCount returns the count of UUIDs in the set.
func (ups *UIDsPermSet) GetUIDCount() int {
	return len(ups.UIDs)
}

// HasUID checks if the specified UUID is in the set.
func (ups *UIDsPermSet) HasUID(target auser.UID) bool {
	return ups.UIDs.IsValid(target)
}

// SetUID adds a UUID to the set if it's not already present.
func (ups *UIDsPermSet) SetUID(target auser.UID) {
	if target.IsNil() {
		return
	}
	if !ups.HasUID(target) {
		ups.UIDs = append(ups.UIDs, target)
	}
}

// RemoveUID removes a UUID from the set.
func (ups *UIDsPermSet) RemoveUID(target auser.UID) {
	if target.IsNil() {
		return
	}
	var uidsNew auser.UIDs
	for _, uid := range ups.UIDs {
		if uid != target {
			uidsNew = append(uidsNew, uid)
		}
	}
	ups.UIDs = uidsNew
}

// HasPerm checks if a specific permission is in the set by its bit value.
func (ups *UIDsPermSet) HasPerm(key string, bit int) bool {
	perm := NewPermByBitValue(key, bit)
	return ups.Perms.MatchesPerm(perm)
}

// UIDsPermSets is a slice of UIDsPermSet pointers.
type UIDsPermSets []*UIDsPermSet

// Validate checks the integrity of each UIDsPermSet in the slice.
func (upss UIDsPermSets) Validate() error {
	return upss.ValidateWithOptions(false)
}

// ValidateWithOptions checks the integrity of each UIDsPermSet in the slice with an option to require a count.
func (upss UIDsPermSets) ValidateWithOptions(mustHaveCount bool) error {
	if upss == nil || len(upss) == 0 {
		if mustHaveCount {
			return fmt.Errorf("UIDsPermSets is empty")
		}
		return nil
	}
	for ii, ups := range upss {
		if err := ups.Validate(); err != nil {
			return fmt.Errorf("UIDsPermSet is nil at index %d", ii)
		}
	}
	return nil
}

// GetUIDCountByPerm returns the count of UUIDs associated with a specific permission.
func (upss UIDsPermSets) GetUIDCountByPerm(perm Perm) int {
	for _, ups := range upss {
		if ups.Perms.HasPerm(perm) {
			return ups.GetUIDCount()
		}
	}
	return 0
}

// GetUIDCountByPermBit returns the count of UUIDs associated with a specific permission bit.
func (upss UIDsPermSets) GetUIDCountByPermBV(key string, bit int) int {
	return upss.GetUIDCountByPerm(*NewPermByBitValue(key, bit))
}

// GetUIDCountByPermString returns the count of UUIDs associated with a specific permission value string.
func (upss UIDsPermSets) GetUIDCountByPermSV(key string, permValue string) int {
	return upss.GetUIDCountByPerm(*NewPermByPair(key, permValue))
}

// HasUIDByPerm checks if a specific UUID is associated with a specific permission.
func (upss UIDsPermSets) HasUIDByPerm(perm Perm, target auser.UID) bool {
	for _, ups := range upss {
		if ups.Perms.HasPerm(perm) {
			return ups.HasUID(target)
		}
	}
	return false
}

// HasUIDByPermBit checks if a specific UUID is associated with a specific permission bit.
func (upss UIDsPermSets) HasUIDByPermBit(key string, bit int, target auser.UID) bool {
	return upss.HasUIDByPerm(*NewPermByBitValue(key, bit), target)
}

// HasUIDByPermString checks if a specific UUID is associated with a specific permission value string.
func (upss UIDsPermSets) HasUIDByPermString(key string, permValue string, target auser.UID) bool {
	return upss.HasUIDByPerm(*NewPermByPair(key, permValue), target)
}

// SetUIDByPerm adds a UUID to the set associated with a specific permission.
func (upss UIDsPermSets) SetUIDByPerm(perm Perm, target auser.UID) {
	for _, ups := range upss {
		if ups.Perms.HasPerm(perm) {
			ups.SetUID(target)
			return
		}
	}
}

// SetUIDByPermBit adds a UUID to the set associated with a specific permission bit.
func (upss UIDsPermSets) SetUIDByPermBit(key string, bit int, target auser.UID) {
	upss.SetUIDByPerm(*NewPermByBitValue(key, bit), target)
}

// SetUIDByPermString adds a UUID to the set associated with a specific permission value string.
func (upss UIDsPermSets) SetUIDByPermString(key string, permValue string, target auser.UID) {
	upss.SetUIDByPerm(*NewPermByPair(key, permValue), target)
}

// RemoveUIDByPerm removes a UUID from the set associated with a specific permission.
func (upss UIDsPermSets) RemoveUIDByPerm(perm Perm, target auser.UID) {
	for _, ups := range upss {
		if ups.Perms.HasPerm(perm) {
			ups.RemoveUID(target)
		}
	}
}

// RemoveUIDByPermBit removes a UUID to the set associated with a specific permission bit.
func (upss UIDsPermSets) RemoveUIDByPermBit(key string, bit int, target auser.UID) {
	upss.RemoveUIDByPerm(*NewPermByBitValue(key, bit), target)
}

// RemoveUIDByPermString removes a UUID to the set associated with a specific permission value string.
func (upss UIDsPermSets) RemoveUIDByPermString(key string, permValue string, target auser.UID) {
	upss.RemoveUIDByPerm(*NewPermByPair(key, permValue), target)
}

// Clean removes any nil UUIDs from the UIDsPermSets.
func (upss UIDsPermSets) Clean() UIDsPermSets {
	var arr UIDsPermSets
	for _, ups := range upss {
		t := ups.UIDs.Clean()
		if t != nil {
			ups.UIDs = t
			arr = append(arr, ups)
		}
	}
	return arr
}

// CreateSingleUIDsPermSetsByKVString creates a UIDsPermSets with a single set of permissions defined by a key-value string.
func CreateSingleUIDsPermSetsByKVString(keyValue string, uids ...auser.UID) UIDsPermSets {
	return UIDsPermSets{
		{
			UIDs:  uids,
			Perms: NewPermSetByString([]string{keyValue}),
		},
	}
}

// CreateSingleUIDsPermSetsByKVPair creates a UIDsPermSets with a single set of permissions defined by a key-value pair.
func CreateSingleUIDsPermSetsByKVPair(key string, value string, uids ...auser.UID) UIDsPermSets {
	return UIDsPermSets{
		{
			UIDs:  uids,
			Perms: NewPermSetByPair(key, value),
		},
	}
}

// CreateSingleUIDsPermSets creates a UIDsPermSets with a single set of permissions.
func CreateSingleUIDsPermSets(perms PermSet, uids ...auser.UID) UIDsPermSets {
	return UIDsPermSets{
		{
			UIDs:  uids,
			Perms: perms,
		},
	}
}
