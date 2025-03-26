package asessions

type RoleFactory map[RoleName]PermSet

// BuildPermSet receives a Role and returns a merged PermSet based on
// the RoleFactory definitions along with the Role's PermsPlus and PermsMinus.
func (rf RoleFactory) BuildPermSet(role *Role) PermSet {
	// Validate the Role
	if role == nil {
		return PermSet{}
	}

	// Initialize an empty PermSet to collect merged permissions
	mergedPermSet := PermSet{}

	// Merge permissions from RoleFactory for each role name in role.Names
	if perms, ok := rf[role.Name]; ok {
		mergedPermSet.MergeByPermSet(perms)
	}

	// Apply PermsPlus to add additional permissions
	mergedPermSet.MergeByPermSet(role.PermsPlus)

	// Apply PermsMinus to remove specified permissions
	mergedPermSet.SubtractByPermSet(role.PermsMinus)

	return mergedPermSet
}

// BuildPermSetWithLimit constructs a PermSet based on RoleFactory and the provided Role,
// but ensures that the permissions do not exceed those defined in defaultPermSet.
func (rf RoleFactory) BuildPermSetWithLimit(role *Role, defaultPermSet PermSet) PermSet {
	// Validate inputs
	if role == nil || defaultPermSet == nil {
		return PermSet{}
	}

	// Step 1: Build the initial PermSet for the Role
	mergedPermSet := rf.BuildPermSet(role)

	// Step 2: Restrict permissions based on the defaultPermSet
	limitedPermSet := PermSet{}
	for key, perm := range mergedPermSet {
		if defaultPerm, ok := defaultPermSet[key]; ok {
			// Limit permissions to what's allowed by defaultPermSet
			limitedPerm := perm.Clone()
			limitedPerm.ReplaceExcessivePermsByChars(defaultPerm.Value()) // Restrict permissions
			limitedPermSet[key] = limitedPerm
		}
	}

	return limitedPermSet
}

// BuildPermSetMulti receives a Role and returns a merged PermSet based on
// the RoleFactory definitions along with the Role's PermsPlus and PermsMinus.
func (rf RoleFactory) BuildPermSetMulti(role *RoleMulti) PermSet {
	// Validate the Role
	if role == nil {
		return PermSet{}
	}

	// Initialize an empty PermSet to collect merged permissions
	mergedPermSet := PermSet{}

	// Merge permissions from RoleFactory for each role name in role.Names
	for _, roleName := range role.Names {
		if perms, ok := rf[roleName]; ok {
			mergedPermSet.MergeByPermSet(perms)
		}
	}

	// Apply PermsPlus to add additional permissions
	mergedPermSet.MergeByPermSet(role.PermsPlus)

	// Apply PermsMinus to remove specified permissions
	mergedPermSet.SubtractByPermSet(role.PermsMinus)

	return mergedPermSet
}

// BuildPermSetWithLimitMulti constructs a PermSet based on RoleFactory and the provided RoleMulti,
// but ensures that the permissions do not exceed those defined in defaultPermSet.
func (rf RoleFactory) BuildPermSetWithLimitMulti(role *RoleMulti, defaultPermSet PermSet) PermSet {
	// Validate inputs
	if role == nil || defaultPermSet == nil {
		return PermSet{}
	}

	// Step 1: Build the initial PermSet for the Role
	mergedPermSet := rf.BuildPermSetMulti(role)

	// Step 2: Restrict permissions based on the defaultPermSet
	limitedPermSet := PermSet{}
	for key, perm := range mergedPermSet {
		if defaultPerm, ok := defaultPermSet[key]; ok {
			// Limit permissions to what's allowed by defaultPermSet
			limitedPerm := perm.Clone()
			limitedPerm.ReplaceExcessivePermsByChars(defaultPerm.Value()) // Restrict permissions
			limitedPermSet[key] = limitedPerm
		}
	}

	return limitedPermSet
}
