package asessions

import (
	"github.com/labstack/echo/v4"
)

// ECHOSCS_OBJECTKEY_USER_SESSION is the key used to store and retrieve the UserSessionPerm object from Echo's context.
const ECHOSCS_OBJECTKEY_USER_SESSION = "us"

// CastLoginSessionPermFromEchoContext attempts to retrieve an ILoginSessionPerm from Echo's context.
func CastLoginSessionPermFromEchoContext(c echo.Context) ILoginSessionPerm {
	if us, ok := c.Get(ECHOSCS_OBJECTKEY_USER_SESSION).(ILoginSessionPerm); ok {
		return us
	}
	return nil
}

// CastUserSessionPermFromEchoContext attempts to retrieve a *UserSessionPerm from Echo's context.
func CastUserSessionPermFromEchoContext(c echo.Context) *UserSessionPerm {
	if us, ok := c.Get(ECHOSCS_OBJECTKEY_USER_SESSION).(*UserSessionPerm); ok {
		return us
	}
	if us, ok := c.Get(ECHOSCS_OBJECTKEY_USER_SESSION).(UserSessionPerm); ok {
		return &us
	}
	return nil
}

// CastUserSessionPermFromILoginSession attempts to cast an ILoginSessionPerm to a *UserSessionPerm.
func CastUserSessionPermFromILoginSession(loginSession ILoginSessionPerm) *UserSessionPerm {
	if us, ok := loginSession.(*UserSessionPerm); ok {
		return us
	}
	return nil
}

// UserSessionPerm extends UserSessionBase with permission-related functionality.
type UserSessionPerm struct {
	UserSessionBase

	// Perms holds the set of permissions associated with the user session.
	Perms PermSet `json:"perms"`
}

// NewUserSessionPerm creates a new UserSessionPerm with default values.
func NewUserSessionPerm() *UserSessionPerm {
	return &UserSessionPerm{
		UserSessionBase: *NewUserSessionBase(),
		Perms:           PermSet{},
	}
}

// NewUserSessionPermWithLoginStatus creates a new UserSessionPerm with a specified login status.
func NewUserSessionPermWithLoginStatus(statusType LoginSessionStatusType) *UserSessionPerm {
	return &UserSessionPerm{
		UserSessionBase: *NewUserSessionBaseWithLoginStatus(statusType),
		Perms:           PermSet{},
	}
}

// GetPerms returns the set of permissions associated with the user session.
func (us *UserSessionPerm) GetPerms() PermSet {
	return us.Perms
}

// HasPerm checks if the user session has a specific permission.
func (us *UserSessionPerm) HasPerm(target Perm) bool {
	return us.Perms.MatchesPerm(&target)
}

// HasPermS checks if the user session has a specific permission.
func (us *UserSessionPerm) HasPermS(keyPermValue string) bool {
	return us.Perms.MatchesPerm(NewPerm(keyPermValue))
}

// HasPermSV checks if the user session has a specific permission value for a given key.
func (us *UserSessionPerm) HasPermSV(key string, permValue string) bool {
	return us.Perms.MatchesPerm(NewPermByPair(key, permValue))
}

// HasPermB checks if the user session has a specific permission value for a given key.
func (us *UserSessionPerm) HasPermB(keyBits string) bool {
	return us.Perms.MatchesPerm(NewPerm(keyBits))
}

// HasPermBV checks if the user session has a specific permission value for a given key.
func (us *UserSessionPerm) HasPermBV(key string, bit int) bool {
	return us.Perms.MatchesPerm(NewPermByBitValue(key, bit))
}

// HasPermSet checks if the user session has on matching permission with the target PermSet.
func (us *UserSessionPerm) HasPermSet(target PermSet) bool {
	return us.Perms.HasPermSet(target)
}
