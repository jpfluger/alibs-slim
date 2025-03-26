package asessions

import (
	"github.com/jpfluger/alibs-slim/aimage"
	"github.com/jpfluger/alibs-slim/atags"
	"github.com/jpfluger/alibs-slim/auser"
	"github.com/jpfluger/alibs-slim/autils"
	"github.com/jpfluger/alibs-slim/auuids"
	"time"
)

// ILoginSessionBase defines the interface for a login session.
type ILoginSessionBase interface {
	// GetSID returns the session ID as a nullable UUID.
	GetSID() auuids.UUID

	// GetStatusType returns the current status of the login session.
	GetStatusType() LoginSessionStatusType

	// GetUsername returns the username associated with the session.
	GetUsername() auser.Username

	// GetLastLogin returns the timestamp of the last login.
	GetLastLogin() time.Time

	// GetLanguageType returns the language preference of the user.
	GetLanguageType() autils.LanguageType

	// GetActions returns the set of actions associated with the session.
	GetActions() ActionKeys

	// IsLoggedIn checks if the user is currently logged in.
	IsLoggedIn() bool

	// IsActionRequired checks if any action is required for the session.
	IsActionRequired() bool

	// FindAction searches for a specific action in the session.
	FindAction(target ActionKey) ActionKey

	// HasAction checks if a specific action exists in the session.
	HasAction(target ActionKey) bool

	// AddAction adds a new action to the session.
	AddAction(ActionKey)

	// SetActions sets the entire set of actions for the session.
	SetActions(ActionKeys)

	// RemoveAction removes an action from the session.
	RemoveAction(ActionKey)

	// GetRecordUserIdentity returns the user identity record from the database.
	GetRecordUserIdentity() auser.RecordUserIdentity

	// GetRUI is a convenience function pointing to GetRecordUserIdentity.
	GetRUI() auser.RecordUserIdentity

	// GetDisplayName returns the display name of the user.
	GetDisplayName() string

	// GetAvatar returns the avatar image of the user.
	GetAvatar() *aimage.Image

	// GetAvatarDataBase64 returns the avatar image in Base64 encoding.
	GetAvatarDataBase64() string

	// GetOptions returns the session options as a map of tags.
	GetOptions() atags.TagMapString

	// HasOption checks if a specific option is set in the session.
	HasOption(key atags.TagKey) bool

	// SetOptions sets the session options using a map of tags.
	SetOptions(options atags.TagMapString)

	// OptionValue returns the value of a specific option.
	OptionValue(key atags.TagKey) string

	// OptionSet sets the value of a specific option.
	OptionSet(key atags.TagKey, value string)

	// OptionsClear clears all options in the session.
	OptionsClear()

	// OptionRemove removes a specific option from the session.
	OptionRemove(key atags.TagKey)

	LoadMeta(target interface{}) error
	SaveMeta(target interface{}) error
}
