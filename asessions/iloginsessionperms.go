package asessions

// ILoginSessionPerm extends ILoginSessionBase to include permission-related methods.
type ILoginSessionPerm interface {
	ILoginSessionBase // Embeds the base login session interface.

	// GetPerms retrieves the set of permissions associated with the session.
	GetPerms() PermSet

	// HasPerm checks if the user session has a specific permission.
	HasPerm(target Perm) bool

	// HasPermS checks if the user session has a specific permission.
	HasPermS(keyPermValue string) bool

	// HasPermSV checks if the user session has a specific permission value for a given key.
	HasPermSV(key string, permValue string) bool

	// HasPermB checks if the user session has a specific permission value for a given key.
	HasPermB(keyBits string) bool

	// HasPermBV checks if the user session has a specific permission value for a given key.
	HasPermBV(key string, bit int) bool

	// HasPermSet checks if the user session has on matching permission with the target PermSet.
	HasPermSet(target PermSet) bool
}
