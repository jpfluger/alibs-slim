package asessions

import (
	"github.com/jpfluger/alibs-slim/atags"
	"github.com/jpfluger/alibs-slim/auser"
	"time"
)

// IUIActionKey defines the interface for UI action keys.
type IUIActionKey interface {
	GetAction() ActionKey
	GetUsername() auser.Username
	GetUrlPost() string
	GetFormDate() time.Time
	SetFormDate(formDate time.Time)
	GetTitle() string
	GetLead() string
	GetFooter() string
	GetTags() atags.TagMapString
}

// UIActionKey struct holds the details for a UI action key.
type UIActionKey struct {
	Action   ActionKey          `json:"actionKey"`          // The action key identifier.
	Username auser.Username     `json:"username,omitempty"` // The associated username.
	FormDate time.Time          `json:"formDate,omitempty"` // The date when the form was submitted.
	UrlPost  string             `json:"urlPost,omitempty"`  // The URL for POST submission.
	Title    string             `json:"title"`              // The title associated with the UI action.
	Lead     string             `json:"lead"`               // The lead text for the UI action.
	Footer   string             `json:"footer"`             // The footer text for the UI action.
	Tags     atags.TagMapString `json:"tags"`               // A map of tags associated with the UI action.
}

// NewUIActionKey creates a new UIActionKey with the provided details.
func NewUIActionKey(action ActionKey, username auser.Username, urlPost string, formDate time.Time, title string, leader string, footer string, tags atags.TagMapString) *UIActionKey {
	if tags == nil {
		tags = atags.TagMapString{}
	}
	if formDate.IsZero() {
		formDate = time.Now().UTC()
	}
	return &UIActionKey{
		Action:   action,
		Username: username,
		UrlPost:  urlPost,
		FormDate: formDate,
		Title:    title,
		Lead:     leader,
		Footer:   footer,
		Tags:     tags,
	}
}

// GetAction returns the action key.
func (u *UIActionKey) GetAction() ActionKey {
	return u.Action
}

// GetUsername returns the username.
func (u *UIActionKey) GetUsername() auser.Username {
	return u.Username
}

// GetUrlPost returns the URL for POST submission.
func (u *UIActionKey) GetUrlPost() string {
	return u.UrlPost
}

// GetFormDate returns the form submission date.
func (u *UIActionKey) GetFormDate() time.Time {
	return u.FormDate
}

// SetFormDate sets the form submission date.
func (u *UIActionKey) SetFormDate(formDate time.Time) {
	u.FormDate = formDate
}

// GetTitle returns the title.
func (u *UIActionKey) GetTitle() string {
	return u.Title
}

// GetLead returns the lead text.
func (u *UIActionKey) GetLead() string {
	return u.Lead
}

// GetFooter returns the footer text.
func (u *UIActionKey) GetFooter() string {
	return u.Footer
}

// GetTags returns the tags map.
func (u *UIActionKey) GetTags() atags.TagMapString {
	return u.Tags
}
