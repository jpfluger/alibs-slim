package asessions

import (
	"encoding/json"
	"fmt"
	"github.com/jpfluger/alibs-slim/aimage"
	"github.com/jpfluger/alibs-slim/atags"
	"github.com/jpfluger/alibs-slim/auser"
	"github.com/jpfluger/alibs-slim/autils"
	"github.com/jpfluger/alibs-slim/auuids"
	"strings"
	"time"
)

// UserSessionBase holds the base information for a user session.
type UserSessionBase struct {
	// SID is a unique identifier used for logging purposes, not shared with the client.
	SID auuids.UUID `json:"sid"`

	// Status indicates the current status of the login session.
	Status LoginSessionStatusType `json:"status"`

	// Username is the human-readable identifier for the user.
	Username auser.Username `json:"username"`

	// RUI is the record user identity for tracking user identity details.
	RUI auser.RecordUserIdentity `json:"rui"`

	// LastLogin is the timestamp of the user's last login, displayable on the client.
	LastLogin time.Time `json:"lastLogin"`

	// LanguageType specifies the preferred language of the user.
	LanguageType autils.LanguageType `json:"langType"`

	// DisplayName is the primary name displayed for the user.
	DisplayName string `json:"displayName"`

	// Actions are tasks that need to be completed as part of the login process.
	// e.g. change password or OTP verify.
	Actions ActionKeys `json:"actions"`

	// Avatar is the user's profile image, stored in base64 format.
	Avatar *aimage.Image `json:"avatar,omitempty"`

	// Options are custom settings or preferences for the user session.
	// e.g. UI or Messages or Alerts of some type.
	Options atags.TagMapString `json:"options"`

	// Meta can hold additional dynamic metadata for the session.
	Meta json.RawMessage `json:"meta,omitempty"`
}

// NewUserSessionBase creates a new instance of UserSessionBase with default values.
func NewUserSessionBase() *UserSessionBase {
	return &UserSessionBase{
		SID:     auuids.NewUUID(),
		Actions: ActionKeys{},
		Options: atags.TagMapString{},
	}
}

// NewUserSessionBaseWithLoginStatus creates a new instance of UserSessionBase with a specified login status.
func NewUserSessionBaseWithLoginStatus(statusType LoginSessionStatusType) *UserSessionBase {
	return &UserSessionBase{
		SID:     auuids.NewUUID(),
		Actions: ActionKeys{},
		Status:  statusType,
		Options: atags.TagMapString{},
	}
}

// GetSID returns the session ID as a nullable UUID.
func (us *UserSessionBase) GetSID() auuids.UUID {
	return us.SID
}

// GetStatusType returns the current status of the login session.
func (us *UserSessionBase) GetStatusType() LoginSessionStatusType {
	return us.Status
}

// GetUsername returns the username or falls back to the RUI email.
func (us *UserSessionBase) GetUsername() auser.Username {
	if !us.Username.IsEmpty() {
		return us.Username
	}
	email := us.RUI.FindLabel("email")
	if email != "" {
		return auser.Username(email)
	}
	return ""
}

// GetLastLogin returns the timestamp of the last login.
func (us *UserSessionBase) GetLastLogin() time.Time {
	return us.LastLogin
}

// IsLoggedIn checks if the user is currently logged in.
func (us *UserSessionBase) IsLoggedIn() bool {
	return us != nil && us.Status > LOGIN_SESSION_STATUS_NONE && !us.GetUsername().IsEmpty()
}

// GetLanguageType returns the language preference of the user.
func (us *UserSessionBase) GetLanguageType() autils.LanguageType {
	return us.LanguageType
}

// GetActions returns the set of actions associated with the session.
func (us *UserSessionBase) GetActions() ActionKeys {
	return us.Actions
}

// IsActionRequired checks if any action is required for the session.
func (us *UserSessionBase) IsActionRequired() bool {
	return len(us.Actions) > 0
}

// FindAction searches for a specific action in the session.
func (us *UserSessionBase) FindAction(target ActionKey) ActionKey {
	return us.Actions.Find(target)
}

// HasAction checks if a specific action exists in the session.
func (us *UserSessionBase) HasAction(target ActionKey) bool {
	return !us.FindAction(target).IsEmpty()
}

// AddAction adds a new action to the session.
func (us *UserSessionBase) AddAction(target ActionKey) {
	us.Actions.Add(target)
}

// RemoveAction removes an action from the session.
func (us *UserSessionBase) RemoveAction(target ActionKey) {
	us.Actions = us.Actions.Remove(target)
}

// SetActions sets the entire set of actions for the session.
func (us *UserSessionBase) SetActions(actions ActionKeys) {
	if actions == nil {
		actions = ActionKeys{}
	}
	us.Actions = actions
}

// GetSIDString returns the string representation of SID or an empty string if SID is not set.
func (us *UserSessionBase) GetSIDString() string {
	return autils.UUIDToStringEmpty(us.SID.UUID)
}

// GetRecordUserIdentity returns the record user identity.
func (us *UserSessionBase) GetRecordUserIdentity() auser.RecordUserIdentity {
	return us.RUI
}

// GetRUI is a convenience function pointing to GetRecordUserIdentity.
func (us *UserSessionBase) GetRUI() auser.RecordUserIdentity {
	return us.RUI
}

// GetDisplayName returns the display name or username if display name is empty.
func (us *UserSessionBase) GetDisplayName() string {
	if us.DisplayName != "" {
		return us.DisplayName
	}
	return us.GetUsername().String()
}

// GetAvatar returns the avatar image of the user.
func (us *UserSessionBase) GetAvatar() *aimage.Image {
	return us.Avatar
}

// GetAvatarDataBase64 returns the base64-encoded avatar image or generates a default one.
func (us *UserSessionBase) GetAvatarDataBase64() string {
	if us.Avatar != nil {
		return us.Avatar.ToImageData()
	}
	firstChar := getFirstChar(us.DisplayName, us.Username)
	return aimage.CreateImageDataCircleAvatar(firstChar)
}

// GetOptions returns the session options as a map of tags.
func (us *UserSessionBase) GetOptions() atags.TagMapString {
	return us.Options
}

// SetOptions sets the session options using a map of tags.
func (us *UserSessionBase) SetOptions(options atags.TagMapString) {
	if options == nil {
		options = atags.TagMapString{}
	}
	us.Options = options
}

// OptionValue returns the value of a specific option.
func (us *UserSessionBase) OptionValue(key atags.TagKey) string {
	return us.Options.Value(key)
}

// HasOption checks if a specific option is set in the session.
func (us *UserSessionBase) HasOption(key atags.TagKey) bool {
	return strings.TrimSpace(us.Options.Value(key)) != ""
}

// OptionSet sets the value of a specific option.
func (us *UserSessionBase) OptionSet(key atags.TagKey, value string) {
	us.Options.Set(key, value)
}

// OptionRemove removes a specific option from the session.
func (us *UserSessionBase) OptionRemove(key atags.TagKey) {
	us.Options.Remove(key)
}

// OptionsClear clears all options in the session.
func (us *UserSessionBase) OptionsClear() {
	us.Options = atags.TagMapString{}
}

// LoadMeta unmarshals the Meta field into the provided target.
// The target must be a pointer to a struct or map.
func (us *UserSessionBase) LoadMeta(target interface{}) error {
	if us.Meta == nil {
		return nil // No metadata to load
	}

	if err := json.Unmarshal(us.Meta, target); err != nil {
		return fmt.Errorf("failed to unmarshal Meta: %w", err)
	}

	return nil
}

// SaveMeta marshals the provided target and saves it into the Meta field.
// The target must be serializable to JSON.
func (us *UserSessionBase) SaveMeta(target interface{}) error {
	if target == nil {
		us.Meta = nil // Clear Meta if target is nil
		return nil
	}

	metaBytes, err := json.Marshal(target)
	if err != nil {
		return fmt.Errorf("failed to marshal Meta: %w", err)
	}

	us.Meta = metaBytes
	return nil
}

// getFirstChar returns the first character of the display name or username.
func getFirstChar(displayName string, username auser.Username) string {
	if displayName != "" {
		return displayName[0:1]
	}
	if username.String() != "" {
		return username.String()[0:1]
	}
	return ""
}
