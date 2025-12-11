package asessions

import (
	"github.com/jpfluger/alibs-slim/auser"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUIDsPermSet_Validate(t *testing.T) {
	validUUID := auser.NewUID()
	validUIDs := auser.UIDs{validUUID}
	validPerms := MustNewPermSetByBits("read", PERM_R)
	validUIDsPermSet := &UIDsPermSet{
		UIDs:  validUIDs,
		Perms: validPerms,
	}
	assert.NoError(t, validUIDsPermSet.Validate())

	// Invalid UIDsPermSet with empty UIDs
	invalidUIDsPermSet := &UIDsPermSet{
		UIDs:  auser.UIDs{},
		Perms: validPerms,
	}
	assert.Error(t, invalidUIDsPermSet.Validate())

	// Invalid UIDsPermSet with nil Perms
	invalidUIDsPermSet = &UIDsPermSet{
		UIDs:  validUIDs,
		Perms: nil,
	}
	assert.Error(t, invalidUIDsPermSet.Validate())
}

func TestUIDsPermSet_GetUIDCount(t *testing.T) {
	uid1 := auser.NewUID()
	uid2 := auser.NewUID()
	uids := auser.UIDs{uid1, uid2}
	ups := &UIDsPermSet{UIDs: uids}
	assert.Equal(t, 2, ups.GetUIDCount())
}

func TestUIDsPermSet_HasUID(t *testing.T) {
	uid := auser.NewUID()
	ups := &UIDsPermSet{UIDs: auser.UIDs{uid}}
	assert.True(t, ups.HasUID(uid))

	nonExistentUID := auser.NewUID()
	assert.False(t, ups.HasUID(nonExistentUID))
}

func TestUIDsPermSet_SetUID(t *testing.T) {
	uid := auser.NewUID()
	ups := &UIDsPermSet{UIDs: auser.UIDs{}}
	ups.SetUID(uid)
	assert.Contains(t, ups.UIDs, uid)
}

func TestUIDsPermSet_RemoveUID(t *testing.T) {
	uid := auser.NewUID()
	ups := &UIDsPermSet{UIDs: auser.UIDs{uid}}
	ups.RemoveUID(uid)
	assert.NotContains(t, ups.UIDs, uid)
}

func TestUIDsPermSet_HasPerm(t *testing.T) {
	uid := auser.NewUID()
	perms := MustNewPermSetByBits("read", PERM_R)
	ups := &UIDsPermSet{
		UIDs:  auser.UIDs{uid},
		Perms: perms,
	}
	assert.True(t, ups.HasPerm("read", PERM_R))
	assert.False(t, ups.HasPerm("write", PERM_U))
}

func TestUIDsPermSets_Validate(t *testing.T) {
	uid := auser.NewUID()
	validPerms := MustNewPermSetByBits("read", PERM_R)
	validUIDsPermSet := &UIDsPermSet{
		UIDs:  auser.UIDs{uid},
		Perms: validPerms,
	}
	upss := UIDsPermSets{validUIDsPermSet}
	assert.NoError(t, upss.Validate())
}

func TestUIDsPermSets_GetUIDCountByPerm(t *testing.T) {
	uid := auser.NewUID()
	perms := MustNewPermSetByBits("read", PERM_R)
	ups := &UIDsPermSet{
		UIDs:  auser.UIDs{uid},
		Perms: perms,
	}
	upss := UIDsPermSets{ups}
	assert.Equal(t, 1, upss.GetUIDCountByPerm(*MustNewPermByPair("read", "R")))
	assert.Equal(t, 0, upss.GetUIDCountByPerm(*MustNewPermByPair("write", "W")))
}

func TestUIDsPermSets_HasUIDByPerm(t *testing.T) {
	uid := auser.NewUID()
	perms := MustNewPermSetByBits("read", PERM_R)
	ups := &UIDsPermSet{
		UIDs:  auser.UIDs{uid},
		Perms: perms,
	}
	upss := UIDsPermSets{ups}
	assert.True(t, upss.HasUIDByPerm(*MustNewPermByPair("read", "R"), uid))

	nonExistentUID := auser.NewUID()
	assert.False(t, upss.HasUIDByPerm(*MustNewPermByPair("read", "R"), nonExistentUID))
}

func TestUIDsPermSets_SetUIDByPerm(t *testing.T) {
	uid := auser.NewUID()
	perms := MustNewPermSetByBits("read", PERM_R)
	ups := &UIDsPermSet{
		UIDs:  auser.UIDs{},
		Perms: perms,
	}
	upss := UIDsPermSets{ups}
	upss.SetUIDByPerm(*MustNewPermByPair("read", "R"), uid)
	assert.Contains(t, upss[0].UIDs, uid)
}

func TestUIDsPermSets_RemoveUIDByPerm(t *testing.T) {
	uid := auser.NewUID()
	perms := MustNewPermSetByBits("read", PERM_R)
	ups := &UIDsPermSet{
		UIDs:  auser.UIDs{uid},
		Perms: perms,
	}
	upss := UIDsPermSets{ups}
	upss.RemoveUIDByPerm(*MustNewPermByPair("read", "R"), uid)
	assert.NotContains(t, upss[0].UIDs, uid)
}

func TestUIDsPermSets_Clean(t *testing.T) {
	uid := auser.NewUID()
	perms := MustNewPermSetByBits("read", PERM_R)
	ups := &UIDsPermSet{
		UIDs:  auser.UIDs{uid, auser.UID{}},
		Perms: perms,
	}
	upss := UIDsPermSets{ups}
	cleanedUpss := upss.Clean()
	assert.NotContains(t, cleanedUpss[0].UIDs, auser.UID{})
}

func TestCreateSingleUIDsPermSetsByKVString(t *testing.T) {
	uid := auser.NewUID()
	keyValue := "read:X"
	upss := CreateSingleUIDsPermSetsByKVString(keyValue, uid)

	assert.Len(t, upss, 1)
	assert.Contains(t, upss[0].UIDs, uid)
	assert.True(t, upss[0].Perms.HasPermS(keyValue))
}

func TestCreateSingleUIDsPermSetsByKVPair(t *testing.T) {
	uid := auser.NewUID()
	key := "read"
	value := "X"
	upss := CreateSingleUIDsPermSetsByKVPair(key, value, uid)

	assert.Len(t, upss, 1)
	assert.Contains(t, upss[0].UIDs, uid)
	assert.True(t, upss[0].Perms.HasPermSV(key, value))
}

func TestCreateSingleUIDsPermSets(t *testing.T) {
	uid := auser.NewUID()
	perms := MustNewPermSetByBits("read", PERM_R)
	upss := CreateSingleUIDsPermSets(perms, uid)

	assert.Len(t, upss, 1)
	assert.Contains(t, upss[0].UIDs, uid)
	assert.Equal(t, perms, upss[0].Perms)
}
