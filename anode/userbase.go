package anode

import (
	"github.com/jpfluger/alibs-slim/auser"
	"github.com/jpfluger/alibs-slim/autils"
	"time"
)

// UserBase is the base struct upon which a UserConfig can be created.
// Different sites/nodes/services will have different user requirements.
// The building blocks are spread around the User* structs (e.g., UserAccount/Vault/Profile/Credential).
// For an example of building a user config, see "userbase_test.go".
type UserBase struct {
	// UID is the primary user ID for the node.
	UID auser.UID `json:"uid,omitempty"`

	// LanguageType defines the language that should be displayed for this user.
	// LanguageType logic is as follows:
	// 1. If the user is unknown, check if the existing session has a LanguageType:
	//    a. If yes, then use it.
	//    b. If no, then use the LanguageType from the client-browser.
	// 2. If the user is known, then get the preferred language from the saved user:
	//    a. If the preferred language is empty, then use the existing session language.
	//    b. If both preferred language and existing session language are unknown, then use the LanguageType from the client-browser.
	// Ref: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Accept-Language
	LanguageType autils.LanguageType `json:"langType"`

	// LastAccess is optional and can be used in different ways.
	// For example, LastAccess may represent:
	// 1. The last "profile/domain switch" of the user (reasonable).
	// 2. It could be updated each time an action was performed on the account (heavier CPU).
	// 3. It is the same as the loginDate (not recommended as it could be confused with #1).
	LastAccess *time.Time `json:"lastAccess,omitempty"`
}

// GetUID returns the user's UID.
func (ub *UserBase) GetUID() auser.UID {
	return ub.UID
}

// GetLastAccess returns the last access time of the user.
func (ub *UserBase) GetLastAccess() *time.Time {
	return ub.LastAccess
}

// GetLanguageType returns the language type of the user.
func (ub *UserBase) GetLanguageType() autils.LanguageType {
	return ub.LanguageType
}
